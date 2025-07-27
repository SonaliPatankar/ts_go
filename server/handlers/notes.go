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

// Main Lambda handler
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    // ✅ Handle CORS preflight request
    if req.HTTPMethod == "OPTIONS" {
        return events.APIGatewayProxyResponse{
            StatusCode: http.StatusOK,
            Headers:    getCORSHeaders(),
            Body:       "",
        }, nil
    }

    switch req.HTTPMethod {
    case "GET":
        return handleGet()
    case "POST":
        return handlePost(req.Body)
    case "PUT":
        return handleUpdate(req.Body)
    case "DELETE":
        // Try to get ID from path parameters (e.g., /notes/{id})
        idStr := req.PathParameters["id"]
        return handleDelete(idStr, req.Body)
    default:
        return events.APIGatewayProxyResponse{
            StatusCode: http.StatusMethodNotAllowed,
            Headers:    getCORSHeaders(),
            Body:       `{"error": "Method not allowed"}`,
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
