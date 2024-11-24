package main

import (
	"embed"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
)

//go:embed facts.txt
var facts embed.FS

func loadFacts() ([]string, error) {
	// Read the embedded file
	data, err := facts.ReadFile("facts.txt")
	if err != nil {
		return nil, err
	}
	// Split the file into lines and return as a slice
	lines := string(data)
	return splitLines(lines), nil
}

func splitLines(data string) []string {
	return strings.Split(strings.TrimSpace(data), "\n")
}

func randomFactHandler(facts []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get a random index
		index := rand.Intn(len(facts))

		// Check if the request is from curl by looking at the User-Agent
		userAgent := r.Header.Get("User-Agent")
		isCurl := strings.HasPrefix(strings.ToLower(userAgent), "curl")

		if isCurl {
			// For curl requests, just return the fact as plain text
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintln(w, facts[index])
			return
		}

		// For browser requests, return the HTML version
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html lang="en">
			<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Random Fact</title>
			<style>
			body {
			background-color: #24273a;
			color: #8aadf4;
			display: flex;
			justify-content: center;
			align-items: center;
			height: 100vh;
			margin: 0;
			font-size: 24px;
			font-weight: bold;
			font-family: Arial, sans-serif;
			}
			</style>
			</head>
			<body>
			<p>%s</p>
			</body>
			</html>
			`, facts[index])
	}
}

func main() {
	// Load facts from the embedded file
	facts, err := loadFacts()
	if err != nil {
		panic(err)
	}

	// Set up the HTTP server
	http.HandleFunc("/", randomFactHandler(facts))
	fmt.Println("Server is running on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

