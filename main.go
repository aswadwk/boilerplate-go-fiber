package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type OpenAi struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	// Temperature      int64          `json:"temperature"`
	MaxTokens int64 `json:"max_tokens"`
	// TopP             int64          `json:"top_p"`
	// FrequencyPenalty int64          `json:"frequency_penalty"`
	// PresencePenalty  int64          `json:"presence_penalty"`
	ResponseFormat ResponseFormat `json:"response_format"`
	Stream         bool           `json:"stream"`
}

type Message struct {
	Role    string    `json:"role"`
	Content []Content `json:"content"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {

		agent := fiber.Post("https://api.openai.com/v1/chat/completions")
		agent.JSON(OpenAi{
			Model: "gpt-3.5-turbo",
			Messages: []Message{
				{
					Role: "user",
					Content: []Content{
						{
							Type: "text",
							Text: "What is the meaning of life?",
						},
					},
				},
			},
			ResponseFormat: ResponseFormat{
				Type: "text",
			},
			MaxTokens: 100,
			Stream:    false,
		})

		// set headers
		agent.ContentType("application/json")
		agent.Add("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
		agent.Debug()

		statusCode, body, errs := agent.Bytes()
		if len(errs) > 0 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"errs": errs,
			})
		}

		var something fiber.Map
		err = json.Unmarshal(body, &something)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"err": err,
			})
		}

		return c.Status(statusCode).JSON(something)
	})

	app.Listen(":3000")
}
