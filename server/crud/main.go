package main

import (
	"fmt"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type ToDo struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Completed bool `json:"completed"`
}
func performPostRequest() {
	myURL := "https://jsonplaceholder.typicode.com/todos"
	todo := ToDo{
		UserID:    1,
		ID:        1,
		Title:     "Sample ToDo",
		Completed: false,
	}

	//convert the ToDo struct to JSON
	jsonMarshal, err := json.Marshal(todo)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	//convert the JSON to a string
	jsonString := string(jsonMarshal)

	//convert the string to reader
	reader := strings.NewReader(jsonString)

	res, err := http.Post(myURL, "application/json", reader)
	if err != nil {
		fmt.Println("Error performing POST request:", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		fmt.Println("Error: Status code is", res.StatusCode)
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response Body:", string(resBody))

}
func performPutRequest() {
	myURL := "https://jsonplaceholder.typicode.com/todos/1" // Use a specific resource for PUT
	todo := ToDo{
		UserID:    1,
		ID:        1,
		Title:     "Sample ToDo",
		Completed: false,
	}

	jsonMarshal, err := json.Marshal(todo)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	reader := strings.NewReader(string(jsonMarshal))

	req, err := http.NewRequest(http.MethodPut, myURL, reader)
	if err != nil {
		fmt.Println("Error creating PUT request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error performing PUT request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: Status code is", resp.StatusCode)
		return
	}
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("PUT Response Body:", string(resBody))
}
func main(){
	
	performPostRequest()
	performPutRequest()
	res,err:= http.Get("https://jsonplaceholder.typicode.com/todos/1")
	if err!=nil{
		fmt.Println("Error fetching data:",err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Error: Status code is", res.StatusCode)
		return
	}

	// body,err := io.ReadAll(res.Body)
	// if err!=nil{
	// 	fmt.Println("Error reading response body:",err)
	// 	return
	// }

	// fmt.Println("Response Body:",string(body))
	var todo ToDo
	if err := json.NewDecoder(res.Body).Decode(&todo); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Println("ToDo:", todo)
}