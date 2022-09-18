package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	output := zerolog.ConsoleWriter{}
	output.Out = os.Stderr
	output.TimeFormat = time.RFC3339

	log.Logger = log.Output(output)
}

func spretty(data interface{}) string {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Warn().Msg("Pretty-printing failed")
		log.Warn().Msg(fmt.Sprintf("Data: %s", data))
		return ""
	}
	return (string)(out)
}

type Article struct {
	Title       string
	Author      string
	Description string
	Chapters    []Chapter
}

type Chapter struct {
	Title    string
	Contents string
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()
		rawbody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error on post %s", err))
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		var body Article
		err = json.Unmarshal(rawbody, &body)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error in json parsing: %s", err))
			log.Error().Msg(fmt.Sprintf("Errorous json: %s", rawbody))
			http.Error(w, "Json parsing error", http.StatusBadRequest)
			return
		}

		log.Info().Msg(fmt.Sprintf("Post request with body:\n%s", spretty(body)))
		fmt.Fprint(w, "hi")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
