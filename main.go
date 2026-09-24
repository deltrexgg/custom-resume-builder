package main

import (
	"custom-resume-builder/handler"

	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
    log.Println("No .env file found, using environment variables")
}

	port := "8030"

	
	
	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.MainPage)
	mux.HandleFunc("/generate-resume", handler.GenerateResumeHandler)
	mux.HandleFunc("/ai-services", handler.AIServiceDropDown)
	mux.HandleFunc("/encrypt-apikey", handler.APIEncrypt)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	println("Server running :", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
