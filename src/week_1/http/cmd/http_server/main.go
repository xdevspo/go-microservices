package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	baseURL       = "localhost:8081"
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

type SyncMap struct {
	elems map[int64]*Note
	m     sync.RWMutex
}

var notes = &SyncMap{
	elems: make(map[int64]*Note),
}

func createNoteHandler(w http.ResponseWriter, r *http.Request) {
	info := &NoteInfo{}
	if err := json.NewDecoder(r.Body).Decode(info); err != nil {
		http.Error(w, "Failed to encode note data", http.StatusBadRequest)
		return
	}

	source := rand.NewSource(time.Now().UnixNano())
	rr := rand.New(source)
	now := time.Now()

	note := &Note{
		Id:        rr.Int63(),
		Info:      *info,
		CreatedAt: now,
		UpdatedAt: now,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "Failed to encode note data", http.StatusInternalServerError)
		return
	}

	notes.m.Lock()
	defer notes.m.Unlock()

	notes.elems[note.Id] = note
}

func getNoteHandler(w http.ResponseWriter, r *http.Request) {
	noteId := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(noteId, 10, 64)
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	notes.m.RLock()
	defer notes.m.RUnlock()

	note, ok := notes.elems[id]
	if !ok {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(note); err != nil {
		http.Error(w, "Failed to encode note data", http.StatusInternalServerError)
		return
	}
}

func main() {
	r := chi.NewRouter()

	r.Post(createPostfix, createNoteHandler)
	r.Get(getPostfix, getNoteHandler)

	fmt.Println("Server started on", baseURL)

	err := http.ListenAndServe(baseURL, r)
	if err != nil {
		panic(err)
	}

}
