package main

import (
	"fmt"
	"encoding/json"

)
type Person struct {
	Name string
	Age  int
	IsAdult bool
}

func main() {
	person := Person{
		Name:    "John Doe",
		Age:    30,
		IsAdult: true,
	}

	fmt.Println("Person:", person)

	person_json, err := json.Marshal(person)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	fmt.Println("JSON:", string(person_json))

	var person2 Person
	err = json.Unmarshal(person_json, &person2)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	fmt.Println("Person 2:", person2)
}