package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chippolot/bonfire"
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
			showPrompts := ctx.Bool("show-prompts")
			dataStore := bonfire.MakeSqliteDataStore()
			opts := bonfire.Options{}

			result, err := bonfire.Generate(token, dataStore, opts)
			if err != nil {
				panic(err)
			}

			jsonData, err := result.Entity.JSON()
			if err != nil {
				panic(err)
			}
			jsonString := string(jsonData)

			if showPrompts {
				fmt.Println("System Prompt:")
				fmt.Println(result.Prompts[0])
				fmt.Println()
				fmt.Println("User Prompt:")
				fmt.Println(result.Prompts[1])
				fmt.Println()
			}
			fmt.Println("Entity:")
			fmt.Println(jsonString)
			fmt.Println()

			fmt.Println("Found References:")
			result.Entity.ParseReferences(func(match string, id string, inner string) string {
				fmt.Printf("id: %s, text: %s, full: %s\n", id, inner, match)
				return match
			})
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
