package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func init_img_prompt(
	color string,
	hairstyle string,
	accessory string,
	personality string,
	vargs ...string) string {

	additional := ""
	if len(vargs) > 0 {
		additional += "\n                " + strings.Join(vargs, ",") + "."
	}

	prompt := fmt.Sprintf(`create ascii art for a cat that runs an advice column in a retro geocities style. 
	the cat's 
		hair color is %s, 
		hairstyle is a %s,
		and has %s. 
		the cat's personality is %s.%s
		limit the response to 20 lines, and do not include any additional text beyond the ascii art.`, color, hairstyle, accessory, personality, additional)
	return prompt
}

func main_xx() {
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// The Gemini 1.5 models are versatile and work with most use cases
	model := client.GenerativeModel("gemini-1.5-flash")
	// model := client.GenerativeModel("gemini-1.5-pro")

	p1 := init_img_prompt("red", "bangs", "long hair", "moody", "the cat is responding to the question 'why does gemini ascii art cats look like triangles with donuts?'")
	fmt.Printf("Prompt:\n%s", p1)

	resp, err := model.GenerateContent(ctx, genai.Text(p1))

	if err != nil {
		fmt.Println("Error opening file:", err)
		return // Exit the program
	}

	has_parts := len(resp.Candidates[len(resp.Candidates)-1].Content.Parts) > 0
	if !has_parts {
		fmt.Printf("No text parts found in first candidate. response content:")
		fmt.Printf("Response:\n%s", resp.Candidates[0].Content)
		return
	}
	fmt.Printf("Response:\n%s", resp.Candidates[0].Content.Parts[0])
	// resp, err := model.GenerateContent(
	// 	  ctx,
	// 	  genai.Text("What's is in this photo?"),
	// 	  genai.ImageData("jpeg", imgData))
}
