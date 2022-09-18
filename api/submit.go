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

		var body interface{}
		err = json.Unmarshal(rawbody, &body)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error in json parsing: %s", err))
			log.Error().Msg(fmt.Sprintf("Errorous json: %s", rawbody))
			http.Error(w, "Json parsing error", http.StatusBadRequest)
			return
		}

		log.Info().Msg(fmt.Sprintf("Post request with body %s", body))
		fmt.Fprint(w, "hi")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
