package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Url struct {
	Id          string    `json:"id"`
	OriginalUrl string    `json:"original_url"`
	ShortUrl    string    `json:"short_url"`
	CreatedAt   time.Time `json:"created_at"`
}

var UrlDb = make(map[string]Url)

func generateShortUrl(OriginalUrl string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalUrl))
	data := hasher.Sum(nil)
	shortUrl := hex.EncodeToString(data)
	return shortUrl[:8]
}
func createUrl(originalUrl string) string {
	shortUrl := generateShortUrl(originalUrl)
	url := Url{
		Id:          shortUrl,
		OriginalUrl: originalUrl,
		ShortUrl:    shortUrl,
		CreatedAt:   time.Now(),
	}
	UrlDb[shortUrl] = url
	return url.ShortUrl
}
func getUrl(id string) (Url, error) {
	url, exists := UrlDb[id]
	if !exists {
		return Url{}, errors.New("URL not found")
	}
	return url, nil
}
func handle(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Url string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println("Received URL:", data)
	shortUrl := createUrl(data.Url)
	shortenUrl := Url{
		Id:          shortUrl,
		OriginalUrl: data.Url,
		ShortUrl:    shortUrl,
		CreatedAt:   time.Now(),
	}
	parsedResult, err := json.Marshal(shortenUrl)
	if err != nil {
		http.Error(w, "Error creating short URL", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(parsedResult)
	// if r.Method == http.MethodPost {
	// 	originalUrl := r.Body
	// 	fmt.Println("Original URL:", originalUrl)
	// }
}
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len("/redirect/"):]
	originalUrl, err := getUrl(url)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}
	fmt.Println("Redirecting to:", originalUrl.OriginalUrl)
	http.Redirect(w, r, originalUrl.OriginalUrl, http.StatusFound)
}
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the URL Shortener Service!")
}
func main() {
	fmt.Println("Hello, World!")
	fmt.Println("Server starting at port 3000 ....")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/shorten", handle)
	http.HandleFunc("/redirect/", redirectHandler)
	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

}
