package main

import (
	"embed"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"
)

//go:embed facts.txt template.html
var content embed.FS

type PageData struct {
	Fact string
}

func loadFacts() ([]string, error) {
	data, err := content.ReadFile("facts.txt")
	if err != nil {
		return nil, fmt.Errorf("error reading facts file: %w", err)
	}
	lines := string(data)
	return splitLines(lines), nil
}

func splitLines(data string) []string {
	return strings.Split(strings.TrimSpace(data), "\n")
}

func randomFactHandler(facts []string, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(facts) == 0 {
			http.Error(w, "No facts available", http.StatusInternalServerError)
			return
		}

		index := rand.Intn(len(facts))
		fact := facts[index]

		userAgent := r.Header.Get("User-Agent")
		isCurl := strings.HasPrefix(strings.ToLower(userAgent), "curl")

		if isCurl {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintln(w, fact)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		err := tmpl.Execute(w, PageData{Fact: fact})
		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	facts, err := loadFacts()
	if err != nil {
		panic(fmt.Sprintf("Failed to load facts: %v", err))
	}

	tmpl, err := template.ParseFS(content, "template.html")
	if err != nil {
		panic(fmt.Sprintf("Failed to parse template: %v", err))
	}

	http.HandleFunc("/", randomFactHandler(facts, tmpl))
	fmt.Println("Server is running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(fmt.Sprintf("Server failed to start: %v", err))
	}
}

