package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "ts_go/server/models"
    "log"
    "github.com/aws/aws-lambda-go/events"
    "strconv"
    "math/rand"
	"time"
)

var notes = []models.Note{}
func generateUniqueID() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(1000000000) // random number up to 1 billion
}
func NotesHttpHandler(w http.ResponseWriter, r *http.Request) {
    // ctx := context.Background()
    req := events.APIGatewayProxyRequest{
        HTTPMethod: r.Method,
    }
    if r.Method == "OPTIONS" {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.WriteHeader(http.StatusOK)
        return
    }
    // Always set CORS headers for all responses
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

    if r.Method != "GET" && r.Method != "POST" {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" || VerifyJWT(authHeader) != nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte(`{"error": "unauthorized"}`))
            return
        }
    }
    if r.Method == "GET" {
        req.Path = "/notes"
    }
    if r.Method == "POST" {
        req.Path = "/notes"
        var newNote models.Note
        if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"error": "Invalid request body"}`))
            return
        }
        notes = append(notes, newNote)
        w.WriteHeader(http.StatusCreated)
        w.Write([]byte(`{"message": "Note created"}`))
        return
    }
    if r.Method == "PUT" {
        req.Path = "/notes"
        var updatedNote models.Note
        if err := json.NewDecoder(r.Body).Decode(&updatedNote); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"error": "Invalid request body"}`))
            return
        }
        for i, note := range notes {
            if note.ID == updatedNote.ID {
                notes[i] = updatedNote
                w.WriteHeader(http.StatusOK)
                w.Write([]byte(`{"message": "Note updated"}`))
                return
            }
        }
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"error": "Note not found"}`))
        return
    }
    if r.Method == "DELETE" {
        req.Path = "/notes"
        idStr := r.URL.Query().Get("id")
        if idStr == "" {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"error": "ID is required"}`))
            return
        }
        id, err := strconv.Atoi(idStr)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"error": "Invalid ID"}`))
            return
        }
        for i, note := range notes {
            if note.ID == id {
                notes = append(notes[:i], notes[i+1:]...)
                w.WriteHeader(http.StatusOK)
                w.Write([]byte(`{"message": "Note deleted"}`))
                return
            }
        }
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"error": "Note not found"}`))
        return
    }

}


func getCORSHeaders() map[string]string {
    return map[string]string{
        "Access-Control-Allow-Origin":  "*",
        "Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
        "Access-Control-Allow-Headers": "Content-Type,Authorization",
        "Content-Type":                 "application/json",
    }
}

func NotesHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    if req.HTTPMethod == "OPTIONS" {
        return events.APIGatewayProxyResponse{
            StatusCode: http.StatusOK,
            Headers:    getCORSHeaders(),
            Body:       "",
        }, nil
    }

    // Skip JWT check on /login
    if req.Path != "/login" {
        if err := VerifyJWT(req.Headers["Authorization"]); err != nil {
            return events.APIGatewayProxyResponse{StatusCode: http.StatusUnauthorized, Headers: getCORSHeaders(), Body: `{"error": "unauthorized"}`}, nil
        }
    }

    switch req.HTTPMethod {
    case "GET":
        return handleGet()
    case "POST":
        if req.Path == "/login" {
            return LoginHandler(req)
        }
        return handlePost(req.Body)
    case "PUT":
        return handleUpdate(req.Body)
    case "DELETE":
        return handleDelete(req.PathParameters["id"], req.Body)
    default:
        return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed, Headers: getCORSHeaders(), Body: `{"error": "method not allowed"}`}, nil
    }
}

// GET /notes
func handleGet() (events.APIGatewayProxyResponse, error) {
    notes, err := getAllNotesFromDynamoDB()
    if err != nil {
        return events.APIGatewayProxyResponse{StatusCode: 500, Headers: getCORSHeaders(), Body: `{"error":"failed to fetch notes"}`}, nil
    }
    body, _ := json.Marshal(notes)
    return events.APIGatewayProxyResponse{StatusCode: 200, Headers: getCORSHeaders(), Body: string(body)}, nil
}

// POST /notes
func handlePost(body string) (events.APIGatewayProxyResponse, error) {
    log.Println("handlePost called with body:", body)
    var newNote models.Note
    if err := json.Unmarshal([]byte(body), &newNote); err != nil {
        return events.APIGatewayProxyResponse{StatusCode: 400, Headers: getCORSHeaders(), Body: `{"error":"invalid request body"}`}, nil
    }

    if newNote.ID == 0 {
		newNote.ID = generateUniqueID()
	}
    log.Printf("Saving note to DynamoDB: %+v", newNote)
    if err := saveNoteToDynamoDB(newNote); err != nil {
        log.Printf("Error saving note to DynamoDB: %v", err)
        return events.APIGatewayProxyResponse{StatusCode: 500, Headers: getCORSHeaders(), Body: `{"error":"failed to save note"}`}, nil
    }
    // Return the created note object
    noteJson, _ := json.Marshal(newNote)
    return events.APIGatewayProxyResponse{StatusCode: 201, Headers: getCORSHeaders(), Body: string(noteJson)}, nil
}

// PUT /notes
func handleUpdate(body string) (events.APIGatewayProxyResponse, error) {
    var note models.Note
    if err := json.Unmarshal([]byte(body), &note); err != nil {
        return events.APIGatewayProxyResponse{StatusCode: 400, Headers: getCORSHeaders(), Body: `{"error":"invalid request body"}`}, nil
    }
    if err := updateNoteInDynamoDB(note); err != nil {
        return events.APIGatewayProxyResponse{StatusCode: 500, Headers: getCORSHeaders(), Body: `{"error":"failed to update note"}`}, nil
    }
    return events.APIGatewayProxyResponse{StatusCode: 200, Headers: getCORSHeaders(), Body: `{"message":"note updated"}`}, nil
}

// DELETE /notes/{id}
func handleDelete(idStr string, body string) (events.APIGatewayProxyResponse, error) {
    var id int
    var err error

    if idStr != "" {
        id, err = strconv.Atoi(idStr)
        if err != nil {
            return events.APIGatewayProxyResponse{StatusCode: 400, Headers: getCORSHeaders(), Body: `{"error":"invalid id"}`}, nil
        }
    } else {
        var reqData struct{ ID int `json:"id"` }
        if err := json.Unmarshal([]byte(body), &reqData); err != nil {
            return events.APIGatewayProxyResponse{StatusCode: 400, Headers: getCORSHeaders(), Body: `{"error":"invalid request body"}`}, nil
        }
        id = reqData.ID
    }

    if err := deleteNoteFromDynamoDB(id); err != nil {
        return events.APIGatewayProxyResponse{StatusCode: 500, Headers: getCORSHeaders(), Body: `{"error":"failed to delete note"}`}, nil
    }
    return events.APIGatewayProxyResponse{StatusCode: 200, Headers: getCORSHeaders(), Body: `{"message":"note deleted"}`}, nil
}
