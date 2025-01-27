package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	baseURL       = "http://localhost:8081"
	createPostfix = "/notes"
	getPostfix    = "/notes/%d"
)

type NoteInfo struct {
	Title    string `json:"title"`
	Context  string `json:"context"`
	Author   string `json:"author"`
	IsPublic bool   `json:"is_public"`
}

type Note struct {
	Id        int64     `json:"id"`
	Info      NoteInfo  `json:"info"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func createNote() (Note, error) {
	note := NoteInfo{
		Title:    gofakeit.BeerName(),
		Context:  gofakeit.IPv4Address(),
		Author:   gofakeit.Name(),
		IsPublic: gofakeit.Bool(),
	}

	data, err := json.Marshal(note)
	if err != nil {
		return Note{}, err
	}

	resp, err := http.Post(baseURL+createPostfix, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return Note{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return Note{}, err
	}

	var createNote Note
	if err := json.NewDecoder(resp.Body).Decode(&createNote); err != nil {
		return Note{}, err
	}

	return createNote, nil
}

func getNote(id int64) (Note, error) {
	resp, err := http.Get(baseURL + fmt.Sprintf(getPostfix, id))
	if err != nil {
		log.Fatal("Failed to get note:", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return Note{}, err
	}

	var note Note
	if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
		return Note{}, err
	}

	return note, nil
}

func main() {
	note, err := createNote()
	if err != nil {
		log.Fatal("Failed to create note:", err)
	}

	log.Printf(color.RedString("Note created:\n"), color.GreenString("%+v", note))

	note, err = getNote(note.Id)
	if err != nil {
		log.Fatal("Failed to get note:", err)
	}

	log.Printf(color.RedString("Note got:\n"), color.GreenString("%+v", note))
}
