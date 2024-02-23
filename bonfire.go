package bonfire

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type Result struct {
	Lore string
}

func Generate(openAIToken string) (Result, error) {
	const prompt string = "Generate lore for a demi-god from a game like dark souls or elden ring. " +
		"The description should just be a fragment and leave the reader with more questions than it answers. " +
		"It may hint at or reference locations or characters which are part of the broader mythos."

	// Generate lore
	lore, err := queryLLM(openAIToken, prompt)
	if err != nil {
		return Result{}, err
	}

	return Result{Lore: lore}, nil
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
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
