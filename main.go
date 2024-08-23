package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)


func main() {
    ExampleGenerativeModel_GenerateContent_textOnly()
}

func ExampleGenerativeModel_GenerateContent_textOnly() {
    ctx := context.Background()
    client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyA7Djpfu1U2UQjk1pdfSl99ZE0j7OnF_zg"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    model := client.GenerativeModel("gemini-1.5-flash")
    resp, err := model.GenerateContent(ctx, genai.Text("Write a story about a magic backpack."))
    if err != nil {
        log.Fatal(err)
    }

    printResponse(resp)
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