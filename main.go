package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/vmorsell/openai-gpt-sdk-go/gpt"
)

func main() {
	app := &cli.App{
		Name: "gpt",
		Action: func(cCtx *cli.Context) error {
			return call(cCtx.Args().Get(0))
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

const (
	apiKey = ""
)

func call(message string) error {
	config := gpt.NewConfig().WithAPIKey(apiKey)
	client := gpt.NewClient(config)

	ch := make(chan *gpt.ChatCompletionChunkEvent)
	go func() {
		err := client.ChatCompletionStream(gpt.ChatCompletionInput{
			Model: gpt.GPT35Turbo,
			Messages: []gpt.Message{
				{
					Role:    gpt.RoleSystem,
					Content: "You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible.",
				},
				{
					Role:    gpt.RoleUser,
					Content: message,
				},
			},
			Stream: true,
		}, ch)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for {
		ev, ok := <-ch
		if !ok {
			break
		}

		if ev.Choices == nil {
			continue
		}

		if ev.Choices[0].Delta.Content == nil {
			continue
		}

		fmt.Print(*ev.Choices[0].Delta.Content)
	}
	fmt.Println()
	return nil
}
