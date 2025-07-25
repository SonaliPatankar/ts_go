package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)


var notes []Note
var nextID = 1



func notesHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(notes)

	case "POST":
		var note Note
		if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		note.ID = nextID
		nextID++
		notes = append(notes, note)
		json.NewEncoder(w).Encode(note)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func noteByIDHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "PUT":
		var updated Note
		if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}
		for i, note := range notes {
			if note.ID == id {
				notes[i].Content = updated.Content
				json.NewEncoder(w).Encode(notes[i])
				return
			}
		}
		http.Error(w, "Note not found", http.StatusNotFound)

	case "DELETE":
		for i, note := range notes {
			if note.ID == id {
				notes = append(notes[:i], notes[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, "Note not found", http.StatusNotFound)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	(*w).Header().Set("Content-Type", "application/json")
}
