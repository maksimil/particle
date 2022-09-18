package api

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

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

type ArticleFiles struct {
	Index    string
	Chapters []ChapterFile
}

type ChapterFile struct {
	Id       string
	Contents string
}

func Spretty(data interface{}) string {
	out, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Warn().Msg("Pretty-printing failed")
		log.Warn().Msg(fmt.Sprintf("Data: %s", data))
		return ""
	}
	return (string)(out)
}

func Collect[K any, V any](arr []K, fn func(v K) V) []V {
	out := make([]V, 0)

	for i := 0; i < len(arr); i++ {
		out = append(out, fn(arr[i]))
	}

	return out
}
