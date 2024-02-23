package bonfire

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type Result struct {
	Entity Entity
}

func Generate(openAIToken string) (Result, error) {
	const promptFormat string = "Generate a dark souls-like %s json entity. " +
		"Json fields: 'name', 'id' (snake case), 'type' (snake case), 'lore' (500 chars max). " +
		"Valid types: %s." +
		"In the lore wrap any references to other entities in <ref id=''></ref> tags ."

	validEntityTypes, _ := json.Marshal(AllEntityTypes)
	prompt := fmt.Sprintf(promptFormat, RandomEntityType(), validEntityTypes)

	// Generate entity
	jsonData, err := queryLLM(openAIToken, prompt)
	if err != nil {
		return Result{}, err
	}

	// Parse entity
	var e Entity
	err = json.Unmarshal([]byte(jsonData), &e)
	if err != nil {
		return Result{}, err
	}

	// Validate entity
	if err = e.validate(); err != nil {
		return Result{}, err
	}

	return Result{Entity: e}, nil
}

func queryLLM(token string, prompt string) (string, error) {
	client := openai.NewClient(token)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
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
