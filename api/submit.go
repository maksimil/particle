package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
		handler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// reading raw body
	defer r.Body.Close()
	rawbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error on post %s", err))
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// parsing json data
	var body Article
	err = json.Unmarshal(rawbody, &body)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error in json parsing: %s", err))
		log.Error().Msg(fmt.Sprintf("Errorous json: %s", rawbody))
		http.Error(w, "Json parsing error", http.StatusBadRequest)
		return
	}

	log.Info().Msg(fmt.Sprintf("Post request with body:\n%s", Spretty(body)))

	// creating file data
	files := makeFiles(body)
	log.Info().Msg(fmt.Sprintf("ArticleFiles:\n%s", Spretty(files)))

	fmt.Fprint(w, "hi")
}

func makeFiles(data Article) ArticleFiles {
	files := ArticleFiles{}

	files.Index = fmt.Sprintf(
		"---\ntitle: %s\nauthor: %s\nchapters:\n%s\n---\n\n%s",
		data.Title, data.Author,
		strings.Join(Collect(data.Chapters, func(k Chapter) string {
			return fmt.Sprintf("  - %s", k.Id)
		}), "\n"), data.Description)

	files.Chapters = make([]ChapterFile, 0)

	for i := 0; i < len(data.Chapters); i++ {
		chapter := data.Chapters[i]
		files.Chapters = append(files.Chapters, ChapterFile{
			Id: chapter.Id,
			Contents: fmt.Sprintf(
				"---\ntitle: %s\n---\n\n%s",
				chapter.Title, chapter.Contents,
			),
		})
	}

	return files
}
