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
		},
		Action: func(ctx *cli.Context) error {
			token := ctx.String("token")
			result, err := bonfire.Generate(token)
			if err != nil {
				panic(err)
			}
			fmt.Println(result.Lore)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
