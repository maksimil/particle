package api

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/snowflake/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"lukechampine.com/blake3"
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

	// sending discord webhook
	SendArticle(body, files)

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

func SendArticle(article Article, files ArticleFiles) {
	client := webhook.New(
		snowflake.GetEnv("WEBHOOK_ID"),
		os.Getenv("WEBHOOK_TOKEN"))
	defer client.Close(context.TODO())

	log.Info().Msg("Initialized the client")

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	tarfiles := []struct{ Name, Body string }{
		{"index.md", files.Index},
	}

	for _, chapter := range files.Chapters {
		tarfiles = append(tarfiles, struct{ Name, Body string }{
			chapter.Id + ".md", chapter.Contents})
	}

	for _, file := range tarfiles {
		hdr := &tar.Header{
			Name: file.Name,
			Mode: 0600,
			Size: int64(len(file.Body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatal().Msgf("Error on tar file %s:\n%s", file.Name, err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			log.Fatal().Msgf("Error on tar file %s:\n%s", file.Name, err)
		}
	}

	if err := tw.Close(); err != nil {
		log.Fatal().Msgf("Error on closing tar archive:\n%s", err)
	}

	hash := blake3.Sum256(buf.Bytes())
	fname := strings.ToLower(base32.StdEncoding.EncodeToString(hash[:])[:8])

	message := discord.NewWebhookMessageCreateBuilder().
		SetContentf(
			"%s\nNew article submission\n\n*title*: %s\n*author*: %s",
			os.Getenv("MENTIONS"), article.Title, article.Author).
		AddFiles(discord.NewFile(fname+".tar", "", bytes.NewReader(buf.Bytes())))

	client.CreateMessage(message.Build())

	log.Info().Msg("Sent the notification")
}
