package handler

import (
	"encoding/json"
	"net/http"
)

type AIServices struct {
	Service string   `json:"service"`
	APIURL  string   `json:"api_url"`
	Models  []string `json:"models"`
}

var dropData = []AIServices{
	{
		Service: "Google Gemini",
		APIURL:  "https://generativelanguage.googleapis.com/v1beta/openai",
		Models: []string{
			"gemini-3.5-flash-lite",
			"gemini-2.5-flash",
			"gemini-2.5-flash-lite",
			"gemini-2.5-pro",
			"gemini-2.0-flash",
			"gemini-2.0-flash-lite",
		},
	},
	{
		Service: "OpenAI",
		APIURL:  "https://api.openai.com/v1",
		Models: []string{
			"gpt-5",
			"gpt-5-mini",
			"gpt-5-nano",
			"gpt-4.1",
			"gpt-4.1-mini",
			"gpt-4.1-nano",
		},
	},
	{
		Service: "Anthropic Claude",
		APIURL:  "https://api.anthropic.com/v1",
		Models: []string{
			"claude-sonnet-4",
			"claude-opus-4",
			"claude-3-7-sonnet",
			"claude-3-5-haiku",
		},
	},
}

func AIServiceDropDown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(
			w,
			"Only GET requests are allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(dropData); err != nil {
		http.Error(
			w,
			"Failed to encode AI services",
			http.StatusInternalServerError,
		)
		return
	}
}