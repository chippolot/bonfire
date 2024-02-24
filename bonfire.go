package bonfire

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/sashabaranov/go-openai"
)

//go:embed res/usr_prompt.txt
//go:embed res/sys_prompt.txt
var ResourcesFS embed.FS

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

	sysPrompt, err := readEmbeddedFileText("res/sys_prompt.txt")
	if err != nil {
		return Result{}, err
	}
	validEntityTypes, _ := json.Marshal(AllEntityTypes)
	sysPrompt = fmt.Sprintf(sysPrompt, validEntityTypes)

	usrPrompt, err := readEmbeddedFileText("res/usr_prompt.txt")
	if err != nil {
		return Result{}, err
	}
	usrPrompt = fmt.Sprintf(usrPrompt, RandomEntityType())

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

	return Result{Entity: e}, nil
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
