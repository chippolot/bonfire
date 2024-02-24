package bonfire

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"
)

type UnknownReference struct {
	Id                  string
	ReferencingEntityId string
}

type DataStore interface {
	GetEntityById(id string) (*Entity, error)
	GetEntitiesByType(entityType EntityType) ([]*Entity, error)
	GetReferencedEntities(id string) ([]*Entity, error)
	GetUnknownReferences() ([]*UnknownReference, error)

	AddEntity(e *Entity) error

	Close() error
}

type Result struct {
	Entity Entity
}

func Generate(openAIToken string, dataStore DataStore) (Result, error) {
	const promptFormat string = "Generate a unique dark souls-like %s json entity. " +
		"Json fields: 'name', 'id' (snake case), 'type' (snake case), 'lore' (500 chars max). " +
		"Valid types: %s." +
		"Wrap any references to other entities in <ref id=''/> tags ."

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
	e.CreatedAt = time.Now().UTC()

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
