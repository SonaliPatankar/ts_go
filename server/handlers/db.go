package handlers

import (
	"context"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"ts_go/server/models"
)

var dbClient *dynamodb.Client
var tableName = "Notes"

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}
	dbClient = dynamodb.NewFromConfig(cfg)
	log.Println("✅ DynamoDB client initialized")
}

// Save a new note
func saveNoteToDynamoDB(note models.Note) error {
	_, err := dbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &tableName,
		Item: map[string]types.AttributeValue{
			"id":      &types.AttributeValueMemberN{Value: strconv.Itoa(note.ID)},
			"content": &types.AttributeValueMemberS{Value: note.Content},
		},
	})
	return err
}

// Get all notes
func getAllNotesFromDynamoDB() ([]models.Note, error) {
	var notes []models.Note
	out, err := dbClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: &tableName,
	})
	if err != nil {
		return nil, err
	}

	for _, item := range out.Items {
		idStr := item["id"].(*types.AttributeValueMemberN).Value
		id, _ := strconv.Atoi(idStr)
		content := item["content"].(*types.AttributeValueMemberS).Value
		notes = append(notes, models.Note{ID: id, Content: content})
	}
	return notes, nil
}

// Update a note
func updateNoteInDynamoDB(note models.Note) error {
	_, err := dbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: strconv.Itoa(note.ID)},
		},
		UpdateExpression:          awsString("SET content = :c"),
		ExpressionAttributeValues: map[string]types.AttributeValue{":c": &types.AttributeValueMemberS{Value: note.Content}},
	})
	return err
}

// Delete a note
func deleteNoteFromDynamoDB(id int) error {
	_, err := dbClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberN{Value: strconv.Itoa(id)},
		},
	})
	return err
}

func awsString(s string) *string {
	return &s
}
