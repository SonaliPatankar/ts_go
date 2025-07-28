package handlers

import (
    "context"
    "encoding/json"
    "net/http"
    "ts_go/server/models"

    "github.com/aws/aws-lambda-go/events"
    "strconv"
)

var notes = []models.Note{}

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

}


func NotesHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    if req.HTTPMethod == "OPTIONS" {
        return events.APIGatewayProxyResponse{
            StatusCode: http.StatusOK,
            Headers:    getCORSHeaders(),
            Body:       "",
        }, nil
    }

    // ⛔ Require JWT for all methods except /login
    if req.Path != "/login" {
        if err := VerifyJWT(req.Headers["Authorization"]); err != nil {
            return events.APIGatewayProxyResponse{
                StatusCode: http.StatusUnauthorized,
                Headers:    getCORSHeaders(),
                Body:       `{"error": "unauthorized"}`,
            }, nil
        }
    }

    switch req.HTTPMethod {
    case "GET":
        return handleGet()
    case "PUT":
        return handleUpdate(req.Body)
    case "DELETE":
        return handleDelete(req.PathParameters["id"], req.Body)
    case "POST":
        if req.Path == "/login" {
            return LoginHandler(req)
        }
        return handlePost(req.Body)
    default:
        return events.APIGatewayProxyResponse{
            StatusCode: http.StatusMethodNotAllowed,
            Headers:    getCORSHeaders(),
            Body:       `{"error": "method not allowed"}`,
        }, nil
    }
}

// PUT /notes
func handleUpdate(body string) (events.APIGatewayProxyResponse, error) {
    var updatedNote models.Note
    if err := json.Unmarshal([]byte(body), &updatedNote); err != nil {
        return events.APIGatewayProxyResponse{
            StatusCode: http.StatusBadRequest,
            Headers:    getCORSHeaders(),
            Body:       `{"error": "Invalid request body"}`,
        }, nil
    }
    for i, note := range notes {
        if note.ID == updatedNote.ID {
            notes[i] = updatedNote
            return events.APIGatewayProxyResponse{
                StatusCode: http.StatusOK,
                Headers:    getCORSHeaders(),
                Body:       `{"message": "Note updated"}`,
            }, nil
        }
    }
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusNotFound,
        Headers:    getCORSHeaders(),
        Body:       `{"error": "Note not found"}`,
    }, nil
}

// DELETE /notes/{id} or /notes (with body)
func handleDelete(idStr string, body string) (events.APIGatewayProxyResponse, error) {
    var id int
    var err error
    if idStr != "" {
        id, err = strconv.Atoi(idStr)
        if err != nil {
            return events.APIGatewayProxyResponse{
                StatusCode: http.StatusBadRequest,
                Headers:    getCORSHeaders(),
                Body:       `{"error": "Invalid ID in path"}`,
            }, nil
        }
    } else {
        // fallback: try to get ID from body for backward compatibility
        var reqData struct {
            ID int `json:"id"`
        }
        if err := json.Unmarshal([]byte(body), &reqData); err != nil {
            return events.APIGatewayProxyResponse{
                StatusCode: http.StatusBadRequest,
                Headers:    getCORSHeaders(),
                Body:       `{"error": "Invalid request body"}`,
            }, nil
        }
        id = reqData.ID
    }
    for i, note := range notes {
        if note.ID == id {
            notes = append(notes[:i], notes[i+1:]...)
            return events.APIGatewayProxyResponse{
                StatusCode: http.StatusOK,
                Headers:    getCORSHeaders(),
                Body:       `{"message": "Note deleted"}`,
            }, nil
        }
    }
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusNotFound,
        Headers:    getCORSHeaders(),
        Body:       `{"error": "Note not found"}`,
    }, nil
}

// GET /notes
func handleGet() (events.APIGatewayProxyResponse, error) {
    body, _ := json.Marshal(notes)
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Headers:    getCORSHeaders(),
        Body:       string(body),
    }, nil
}

// POST /notes
func handlePost(body string) (events.APIGatewayProxyResponse, error) {
    var newNote models.Note
    json.Unmarshal([]byte(body), &newNote)
    notes = append(notes, newNote)

    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusCreated,
        Headers:    getCORSHeaders(),
        Body:       `{"message": "Note created"}`,
    }, nil
}

// CORS headers for all responses
func getCORSHeaders() map[string]string {
    return map[string]string{
        "Access-Control-Allow-Origin":  "*", // 🔐 use exact S3 domain in prod
        "Access-Control-Allow-Methods": "GET,POST,PUT,DELETE,OPTIONS",
        "Access-Control-Allow-Headers": "Content-Type,Authorization",
        "Content-Type":                 "application/json",
    }
}
