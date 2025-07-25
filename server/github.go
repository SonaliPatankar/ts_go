package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type GitHubUser struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Location  string `json:"location"`
	PublicRepos int  `json:"public_repos"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}

func githubUserHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)

	username := strings.TrimPrefix(r.URL.Path, "/github/user/")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	resp, err := http.Get("https://api.github.com/users/" + username)
	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var user GitHubUser
	json.Unmarshal(body, &user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
