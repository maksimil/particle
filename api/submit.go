package handler

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
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Chapters    []Chapter `json:"chapters"`
}

type Chapter struct {
	Title    string `json:"title"`
	Id       string `json:"id"`
	Contents string `json:"contents"`
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

		files := makeFiles(body)
		log.Info().Msg(fmt.Sprintf("ArticleFiles:\n%s\n%s", spretty(files), files.Index))

		fmt.Fprint(w, "hi")
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type ArticleFiles struct {
	Index    string
	Chapters []ChapterFile
}

type ChapterFile struct {
	Id       string
	Contents string
}

func c[K any](v K) K {
	return v
}

func collect[K any, V any](arr []K, fn func(v K) V) []V {
	out := make([]V, 0)

	for i := 0; i < len(arr); i++ {
		out = append(out, fn(arr[i]))
	}

	return out
}

func makeFiles(data Article) ArticleFiles {
	files := ArticleFiles{}

	files.Index = fmt.Sprintf(
		"---\ntitle: %s\nauthor: %s\nchapters:\n%s\n---\n\n%s",
		data.Title, data.Author,
		strings.Join(collect(data.Chapters, func(k Chapter) string {
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
