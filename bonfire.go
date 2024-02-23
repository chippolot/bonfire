package bonfire

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type Result struct {
	Lore string
}

var entityTypes = []string{
	"character",
	"item",
	"location",
	"faction",
	"event",
}

func Generate(openAIToken string) (Result, error) {
	const promptFormat string = "Generate a dark souls-like %s json entry. " +
		"Json fields: 'name', 'id' (snake case), 'type' (snake case), 'lore' (500 chars max)." +
		"Valid types are: %s." +
		"In the lore wrap any references to other entities in <ref id=''></ref> tags."

	validEntityTypes := "'" + strings.Join(entityTypes, "', '") + "'"
	prompt := fmt.Sprintf(promptFormat, randomEntityType(), validEntityTypes)

	// Generate lore
	lore, err := queryLLM(openAIToken, prompt)
	if err != nil {
		return Result{}, err
	}

	return Result{Lore: lore}, nil
}

func randomEntityType() string {
	entityIdx := rand.Intn(len(entityTypes))
	return entityTypes[entityIdx]
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
