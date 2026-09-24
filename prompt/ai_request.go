package prompt

// OpenAI-compatible AI client and resume tailoring logic.

import (
	"bytes"
	"custom-resume-builder/utils"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Resume struct {
	Profile    Profile      `json:"profile"`
	Contact    Contact      `json:"contact"`
	Summary    string       `json:"summary"`
	Skills     Skills       `json:"skills"`
	Experience []Experience `json:"experience"`
	Education  []Education  `json:"education"`
	Projects   []Project    `json:"projects"`
}

type Profile struct {
	Name        string `json:"name"`
	Designation string `json:"designation"`
}

type Contact struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Github   string `json:"github"`
	Linkedin string `json:"linkedin"`
	Website  string	`json:"website"`
	Location string `json:"location"`
}

type Skills struct {
	Languages            []string `json:"languages"`
	FrameworksLibraries   []string `json:"frameworks_libraries"`
	DatabasesCaching     []string `json:"databases_caching"`
	DevOpsInfrastructure []string `json:"devops_infrastructure"`
}

type Experience struct {
	Company    string   `json:"company"`
	Type       string   `json:"type"`
	Date       string   `json:"date"`
	Role       string   `json:"role"`
	Highlights []string `json:"highlights"`
}

type Education struct {
	Degree      string `json:"degree"`
	Institution string `json:"institution"`
	University  string `json:"university"`
	Period      string `json:"period"`
}

type Project struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Highlights []string `json:"highlights"`
}


type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func TailorResume(
	jobDescription string,
	masterResume Resume,
	apiURL string,
	aiModel string,
	encryptedKey string,
) (*Resume, error) {

	resumeJSON, err := json.Marshal(masterResume)
	if err != nil {
		return nil, fmt.Errorf("marshal resume: %w", err)
	}

	userPrompt := fmt.Sprintf(
		"JOB DESCRIPTION:\n%s\n\nCANDIDATE RESUME DATA:\n%s",
		jobDescription,
		resumeJSON,
	)

	body := ChatRequest{
		Model: aiModel,
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: PROMPT_TO_TAILOR_RESUME,
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
		Temperature: 0,
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		apiURL+"/chat/completions",
		bytes.NewReader(bodyJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	apiKey, err := utils.DecryptAPIKey(encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("API encryption: %w", err)
	}

	if apiKey == "" {
		return nil, fmt.Errorf("API key is empty")
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("AI request: %w", err)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"AI returned status %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	var aiResponse ChatResponse

	if err := json.Unmarshal(responseBody, &aiResponse); err != nil {
		return nil, fmt.Errorf("parse AI response: %w", err)
	}

	if len(aiResponse.Choices) == 0 {
		return nil, fmt.Errorf("AI returned no choices")
	}

	content := cleanJSON(aiResponse.Choices[0].Message.Content)

	var tailoredResume Resume

	if err := json.Unmarshal([]byte(content), &tailoredResume); err != nil {
		return nil, fmt.Errorf(
			"invalid resume JSON: %w\nresponse: %s",
			err,
			content,
		)
	}

	return &tailoredResume, nil
}

func cleanJSON(content string) string {
	content = strings.TrimSpace(content)

	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")

	return strings.TrimSpace(content)
}