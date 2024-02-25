package bonfire

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/chippolot/bonfire/internal/util"
	"github.com/mroth/weightedrand/v2"
)

const (
	PromptWorld PromptType = iota
	PromptCatalyst
	PromptEntity
)

type PromptType int

var promptTypeStringMapping = util.NewStringMapping[PromptType](map[PromptType]string{
	PromptWorld:    "world",
	PromptCatalyst: "catalyst",
	PromptEntity:   "entity",
})

func ParsePromptType(str string) (PromptType, error) {
	if val, ok := promptTypeStringMapping.ToValue[str]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("unknown prompt type: %s", str)
}

func (pt PromptType) ToString() (string, error) {
	if str, ok := promptTypeStringMapping.ToString[pt]; ok {
		return str, nil
	}
	return "", fmt.Errorf("unknown prompt type: %v", pt)
}

func generatePrompts(promptType PromptType, style string, res *embed.FS) (string, string, error) {

	switch promptType {
	case PromptWorld:
		return generateWorldPrompts(style, res)
	case PromptCatalyst:
		return generateCatalystPrompts(style, res)
	case PromptEntity:
		return generateEntityPrompts(style, res)
	}
	return "", "", fmt.Errorf("unknown prompt type %v", promptType)
}

func generateWorldPrompts(style string, res *embed.FS) (string, string, error) {
	prompts, err := util.ReadEmbeddedFiles(res, "res/prompts/system.txt", "res/prompts/user_world.txt")
	if err != nil {
		return "", "", err
	}

	maxLength := 800

	validEntityTypes, _ := json.Marshal(AllEntityTypes)
	sysPrompt := fmt.Sprintf(prompts[0], validEntityTypes, maxLength)
	usrPrompt := fmt.Sprintf(prompts[1], style)

	return sysPrompt, usrPrompt, nil
}

func generateCatalystPrompts(style string, res *embed.FS) (string, string, error) {
	prompts, err := util.ReadEmbeddedFiles(res, "res/prompts/system.txt", "res/prompts/user_catalyst.txt")
	if err != nil {
		return "", "", err
	}

	maxLength := 800

	validEntityTypes, _ := json.Marshal(AllEntityTypes)
	sysPrompt := fmt.Sprintf(prompts[0], validEntityTypes, maxLength)
	usrPrompt := fmt.Sprintf(prompts[1], style)

	return sysPrompt, usrPrompt, nil
}

func generateEntityPrompts(style string, res *embed.FS) (string, string, error) {
	prompts, err := util.ReadEmbeddedFiles(res, "res/prompts/system.txt", "res/prompts/user_entity.txt")
	if err != nil {
		return "", "", err
	}

	maxLength := 500

	validEntityTypes, _ := json.Marshal(AllNonSingletonEntityTypes)
	sysPrompt := fmt.Sprintf(prompts[0], validEntityTypes, maxLength)
	usrPrompt := fmt.Sprintf(prompts[1], RandomNonSingletonEntityType(), style)

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
