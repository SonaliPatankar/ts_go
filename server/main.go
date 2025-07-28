package main

import (
    "github.com/aws/aws-lambda-go/lambda"
    "ts_go/server/handlers"
)

func main() {
    lambda.Start(handlers.NotesHandler)
}
