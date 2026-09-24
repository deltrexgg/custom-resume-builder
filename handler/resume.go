package handler

import (
	"custom-resume-builder/docx"
	"custom-resume-builder/pdf"
	"custom-resume-builder/prompt"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)


type RequestBody struct {
	JobDescription string        `json:"job_description"`
	PersonalData   prompt.Resume `json:"personal_data"`
	Format         string        `json:"format"`

	APIURL		   string		 `json:"api_url"`
	AIModel		   string		 `json:"ai_model"`
	EncryptedKey   string		 `json:"encrypted_key"`
}

func GenerateResumeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody RequestBody

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&requestBody); err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(requestBody.JobDescription) == "" {
		http.Error(w, "job_description is required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(requestBody.PersonalData.Profile.Name) == "" {
		http.Error(w, "personal_data.profile.name is required", http.StatusBadRequest)
		return
	}

	format := strings.ToLower(strings.TrimSpace(requestBody.Format))

	if format == "" {
		format = "pdf"
	}

	if format != "pdf" && format != "docx" {
		http.Error(w, "format must be either pdf or docx", http.StatusBadRequest)
		return
	}

	data, err := prompt.TailorResume(
		requestBody.JobDescription,
		requestBody.PersonalData,
		requestBody.APIURL,
		requestBody.AIModel,
		requestBody.EncryptedKey,
	)
	if err != nil {
		log.Printf("Failed to tailor resume: %v", err)
		http.Error(w, "Failed to generate resume", http.StatusInternalServerError)
		return
	}

	name := sanitizeFileName(requestBody.PersonalData.Profile.Name)

	if name == "" {
		name = "resume"
	}

	baseName := name + "_resume"
	docxPath := filepath.Join(os.TempDir(), baseName+".docx")

	if err := docx.GenerateResumeDocx(*data, docxPath); err != nil {
		log.Printf("Failed to generate DOCX: %v", err)
		http.Error(w, "Failed to generate DOCX", http.StatusInternalServerError)
		return
	}

	defer os.Remove(docxPath)

	var filePath string
	var contentType string
	var fileName string

	if format == "docx" {
		filePath = docxPath
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		fileName = baseName + ".docx"
	} else {
		outputDir := os.TempDir()

		if err := pdf.ConvertToPDF(docxPath, outputDir); err != nil {
			log.Printf("Failed to convert DOCX to PDF: %v", err)
			http.Error(w, "Failed to convert resume to PDF", http.StatusInternalServerError)
			return
		}

		filePath = filepath.Join(
			outputDir,
			baseName+".pdf",
		)

		contentType = "application/pdf"
		fileName = baseName + ".pdf"

		defer os.Remove(filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open generated file: %v", err)
		http.Error(w, "Failed to open generated resume", http.StatusInternalServerError)
		return
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "Failed to read generated file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, fileName),
	)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))

	http.ServeContent(
		w,
		r,
		fileName,
		stat.ModTime(),
		file,
	)
}

func sanitizeFileName(name string) string {
	name = strings.TrimSpace(name)

	name = regexp.MustCompile(`[^a-zA-Z0-9_-]+`).ReplaceAllString(name, "_")

	name = strings.Trim(name, "_")

	return name
}
