package bonfire

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/mroth/weightedrand/v2"
	"github.com/sashabaranov/go-openai"
)

//go:embed res/usr_prompt.txt
//go:embed res/sys_prompt.txt
var ResourcesFS embed.FS

type Options struct {
}

type Result struct {
	Prompts []string
	Entity  Entity
}

func Generate(openAIToken string, dataStore DataStore, opts Options) (Result, error) {
	// Generare prompts
	sysPrompt, usrPrompt, err := generatePrompts()
	if err != nil {
		return Result{}, err
	}

	// Generate entity
	jsonData, err := queryLLM(openAIToken, sysPrompt, usrPrompt)
	if err != nil {
		return Result{}, err
	}

	// Parse entity
	var e Entity
	err = json.Unmarshal([]byte(jsonData), &e)
	if err != nil {
		return Result{}, err
	}
	e.CreatedAt = time.Now().UTC()

	// Validate entity
	if err = e.validate(); err != nil {
		return Result{}, err
	}

	return Result{
		Prompts: []string{sysPrompt, usrPrompt},
		Entity:  e,
	}, nil
}

func generatePrompts() (string, string, error) {
	sysPrompt, err := readEmbeddedFileText("res/sys_prompt.txt")
	if err != nil {
		return "", "", err
	}
	validEntityTypes, _ := json.Marshal(AllEntityTypes)
	sysPrompt = fmt.Sprintf(sysPrompt, validEntityTypes)

	usrPrompt, err := readEmbeddedFileText("res/usr_prompt.txt")
	if err != nil {
		return "", "", err
	}
	usrPrompt = fmt.Sprintf(usrPrompt, RandomEntityType())

	// Determine number of references
	chooser, _ := weightedrand.NewChooser(
		weightedrand.NewChoice(1, 3),
		weightedrand.NewChoice(2, 2),
		weightedrand.NewChoice(3, 1),
	)
	numReferences := chooser.Pick()
	usrPrompt += fmt.Sprintf("\nInclude %d references to other entities.", numReferences)

	return sysPrompt, usrPrompt, nil
}

func queryLLM(token string, systemPrompt string, userPrompt string) (string, error) {
	client := openai.NewClient(token)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{Type: openai.ChatCompletionResponseFormatTypeJSONObject},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func readEmbeddedFileText(path string) (string, error) {
	file, err := ResourcesFS.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
