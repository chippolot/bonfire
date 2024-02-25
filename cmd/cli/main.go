package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chippolot/bonfire"
	"github.com/chippolot/bonfire/internal/util"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	app := &cli.App{
		Name:  "bonfire",
		Usage: "generates soulslike lore",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "token",
				Aliases:  []string{"t"},
				Value:    "",
				Usage:    "OpenAI token",
				Required: false,
				EnvVars:  []string{"OPENAI_API_KEY"},
			},
			&cli.StringFlag{
				Name:     "prompt-type",
				Aliases:  []string{"pt"},
				Value:    "entity",
				Usage:    "Type of prompt to generate",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "show-prompts",
				Aliases:  []string{"p"},
				Value:    false,
				Usage:    "If true, prompts will be displayed in output",
				Required: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			token := ctx.String("token")
			promptTypeStr := ctx.String("prompt-type")
			promptType, err := bonfire.ParsePromptType(promptTypeStr)
			if err != nil {
				panic(err)
			}
			showPrompts := ctx.Bool("show-prompts")
			dataStore := bonfire.MakeSqliteDataStore()
			opts := bonfire.Options{}

			result, err := bonfire.Generate(promptType, token, dataStore, opts)
			if err != nil {
				panic(err)
			}

			if showPrompts {
				fmt.Println("System Prompt:")
				fmt.Println(result.Prompts[0])
				fmt.Println()
				fmt.Println("User Prompt:")
				fmt.Println(result.Prompts[1])
				fmt.Println()
			}

			jsonString, err := util.MarshalJsonNoHtmlEscape(result.Response)
			if err != nil {
				panic(err)
			}

			fmt.Println("Response:")
			fmt.Println(jsonString)
			fmt.Println()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
