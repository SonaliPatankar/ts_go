// package main

// import (
// 	"log"
// 	"net/http"
// )

// func main() {
// 	http.HandleFunc("/github/user/", githubUserHandler) // from github.go
// 	log.Println("Server running on http://localhost:8080")
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/notes", notesHandler)
	http.HandleFunc("/notes/", noteByIDHandler)

	log.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
