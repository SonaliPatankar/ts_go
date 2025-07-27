package main

import (
    "github.com/aws/aws-lambda-go/lambda"
    "go-notes-lambda/handlers"
)

func main() {
    lambda.Start(handlers.Handler)
}
