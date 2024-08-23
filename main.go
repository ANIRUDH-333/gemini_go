package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)


func main() {
	ExampleChatSession_history()
}

func ExampleGenerativeModel_GenerateContent_textOnly() {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    model := client.GenerativeModel("gemini-1.5-flash")
    resp, err := model.GenerateContent(ctx, genai.Text("Do you know Sentinel in Guidewire"))
    if err != nil {
        log.Fatal(err)
    }

    printResponse(resp)
}

func ExampleGenerativeModel_GenerateContent_systemInstruction() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SystemInstruction = genai.NewUserContent(genai.Text("You are a cat. Your name is Neko."))
	resp, err := model.GenerateContent(ctx, genai.Text("Good morning! Who are you?"))
	if err != nil {
		log.Fatal(err)
	}
	printResponse(resp)

}

func ExampleGenerativeModel_GenerateContent_config() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro-latest")
	model.SetTemperature(0.9)
	model.SetTopP(0.5)
	model.SetTopK(20)
	model.SetMaxOutputTokens(100)
	model.SystemInstruction = genai.NewUserContent(genai.Text("You are a hera in the bollywood. You are very good with dance"))
	model.ResponseMIMEType = "application/json"
	resp, err := model.GenerateContent(ctx, genai.Text("Can you tell me the different types of dances"))
	if err != nil {
		log.Fatal(err)
	}
	printResponse(resp)

}

func ExampleGenerativeModel_GenerateContent_safetySettingMulti() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockLowAndAbove,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockMediumAndAbove,
		},
	}
	resp, err := model.GenerateContent(ctx, genai.Text("I support Martians Soccer Club and I think Jupiterians Football Club sucks! Write a ironic phrase about them."))
	if err != nil {
		log.Fatal(err)
	}
	printResponse(resp)

}

func ExampleGenerativeModel_GenerateContentStream() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	iter := model.GenerateContentStream(ctx, genai.Text("Write a story about a magic backpack."))
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		printResponse(resp)
	}

}

func ExampleGenerativeModel_jSONSchema() {
	// This example shows how to get a JSON response that conforms to a schema.
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro-latest")
	// Ask the model to respond with JSON.
	model.ResponseMIMEType = "application/json"
	// Specify the schema.
	model.ResponseSchema = &genai.Schema{
		Type:  genai.TypeArray,
		Items: &genai.Schema{Type: genai.TypeString},
	}
	resp, err := model.GenerateContent(ctx, genai.Text("List a few popular cookie recipes using this JSON schema."))
	if err != nil {
		log.Fatal(err)
	}
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			var recipes []string
			if err := json.Unmarshal([]byte(txt), &recipes); err != nil {
				log.Fatal(err)
			}
			fmt.Println(recipes)
		}
	}

}

func ExampleGenerativeModel_jSONNoSchema() {
	// This example shows how to get a JSON response without requestin a specific
	// schema.
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-pro-latest")
	// Ask the model to respond with JSON.
	model.ResponseMIMEType = "application/json"
	resp, err := model.GenerateContent(ctx, genai.Text("List a few popular cookie recipes."))
	if err != nil {
		log.Fatal(err)
	}

	printResponse(resp)

}

// This example shows how to set the History field on ChatSession explicitly.
func ExampleChatSession_history() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	cs := model.StartChat()

	cs.History = []*genai.Content{
		{
			Parts: []genai.Part{
				genai.Text("Hello, I have 2 dogs in my house."),
			},
			Role: "user",
		},
		{
			Parts: []genai.Part{
				genai.Text("Great to meet you. What would you like to know?"),
			},
			Role: "model",
		},
	}

	res, err := cs.SendMessage(ctx, genai.Text("How many paws are in my house?"))
	if err != nil {
		log.Fatal(err)
	}
	printResponse(res)

}

func ExampleChatSession_streaming() {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	cs := model.StartChat()

	cs.History = []*genai.Content{
		{
			Parts: []genai.Part{
				genai.Text("Hello, I have 2 dogs in my house."),
			},
			Role: "user",
		},
		{
			Parts: []genai.Part{
				genai.Text("Great to meet you. What would you like to know?"),
			},
			Role: "model",
		},
	}

	iter := cs.SendMessageStream(ctx, genai.Text("How many paws are in my house?"))
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		printResponse(resp)
	}

}



func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}