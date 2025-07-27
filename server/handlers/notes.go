package handlers

import (
    "context"
    "encoding/json"
    "go-notes-lambda/models"
    "github.com/aws/aws-lambda-go/events"
    "net/http"
)

var notes = []models.Note{}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    switch req.HTTPMethod {
    case "GET":
        return handleGet()
    case "POST":
        return handlePost(req.Body)
    default:
        return events.APIGatewayProxyResponse{
            StatusCode: http.StatusMethodNotAllowed,
            Body:       `{"error": "Method not allowed"}`,
        }, nil
    }
}

func handleGet() (events.APIGatewayProxyResponse, error) {
    body, _ := json.Marshal(notes)
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusOK,
        Body:       string(body),
    }, nil
}

func handlePost(body string) (events.APIGatewayProxyResponse, error) {
    var newNote models.Note
    json.Unmarshal([]byte(body), &newNote)
    notes = append(notes, newNote)
    return events.APIGatewayProxyResponse{
        StatusCode: http.StatusCreated,
        Body:       `{"message": "Note created"}`,
    }, nil
}
