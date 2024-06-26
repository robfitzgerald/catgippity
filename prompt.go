package main

import (
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

func cat_welcome_prompt(
	record CatRecord, png []byte) []genai.Part {

	personality := strings.Join(record.Personality, ", ")

	query := fmt.Sprintf(`
	you are the cat in the attached image file, and your job is to provide advice to people.
		your name is %s,
		and you can be described as %s
		say something welcoming to your new customer. 
		be sure to find ways to inject a few cat noises into your statements.
		please limit the response to 2 sentences.`,
		record.Name,
		personality,
	)

	prompt := []genai.Part{
		genai.ImageData("png", png),
		genai.Text(query),
	}

	return prompt
}

func cat_query_prompt(
	png []byte, question string, history string) []genai.Part {

	history_statement := ""
	if len(history) == 0 {
		history_statement = ""
	} else {
		history_statement = fmt.Sprintf("your current customer has had the following chat history with you: %s", history)
	}

	query := fmt.Sprintf(`
	you are the cat in the attached image file, and your job is to provide advice to people.
		%s
		your current customer has now asked the following:
		"%s"
		provide some advice to the user's statements.
		be sure to find ways to inject a few cat noises into your statements.
		please limit the response to a paragraph.`,
		history_statement,
		question,
	)

	prompt := []genai.Part{
		genai.ImageData("png", png),
		genai.Text(query),
	}

	return prompt
}
