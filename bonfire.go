package bonfire

import (
	"embed"
	"encoding/json"
	"time"

	"github.com/chippolot/bonfire/internal/llm"
)

//go:embed res/prompts/system.txt
//go:embed res/prompts/user_world.txt
//go:embed res/prompts/user_catalyst.txt
//go:embed res/prompts/user_entity.txt
var ResourcesFS embed.FS

type Options struct {
	Style string
}

type GenerateResult struct {
	Prompts  []string
	Response *Response
}

type Response struct {
	Entity     *Entity                `json:"entity"`
	References []*EntityReferenceHint `json:"references"`
}

func (r *Response) validate() error {
	return r.Entity.validate()
}

func Generate(promptType PromptType, openAIToken string, dataStore DataStore, opts Options) (GenerateResult, error) {
	style := opts.Style
	if style == "" {
		style = "a dark souls-like game"
	}

	// Generare prompts
	sysPrompt, usrPrompt, err := generatePrompts(promptType, style, &ResourcesFS)
	if err != nil {
		return GenerateResult{}, err
	}

	// Generate response
	jsonData, err := llm.Query(openAIToken, sysPrompt, usrPrompt)
	if err != nil {
		return GenerateResult{}, err
	}

	// Parse response
	var r Response
	err = json.Unmarshal([]byte(jsonData), &r)
	if err != nil {
		return GenerateResult{}, err
	}
	r.Entity.CreatedAt = time.Now().UTC()

	// Validate response
	if err = r.validate(); err != nil {
		return GenerateResult{}, err
	}

	return GenerateResult{
		Prompts:  []string{sysPrompt, usrPrompt},
		Response: &r,
	}, nil
}
