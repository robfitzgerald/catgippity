package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
)

// reads a pre-recorded welcome message from a cat
func welcome(config ConfigFile, cat_id int) (string, error) {
	if cat_id < 0 || cat_id >= len(config.Cats) {
		return "", fmt.Errorf("cat_id %d is invalid", cat_id)
	}
	record := config.Cats[cat_id]
	return record.Welcome, nil
}

// submits a query to Gemini to generate a custom welcome message to the user
func welcome_workflow(config ConfigFile, cwd string, cat_id int, model *genai.GenerativeModel, ctx context.Context) (string, error) {

	record := config.Cats[cat_id]
	img_filepath := filepath.Join(cwd, config.Metadata.ImageDirectory, record.File)
	img_file, file_err := os.Open(img_filepath)
	if file_err != nil {
		return "", errors.Join(file_err, fmt.Errorf("error finding image file %s", img_filepath))
	}
	img_bytes, bytes_err := io.ReadAll(img_file)
	if bytes_err != nil {
		return "", errors.Join(bytes_err, fmt.Errorf("error reading image file %s", img_filepath))
	}
	prompt := cat_welcome_prompt(record, img_bytes)

	res, model_err := model.GenerateContent(ctx, prompt...)
	if model_err != nil {
		return "", errors.Join(model_err, fmt.Errorf("invocation of Gemini model failed"))
	}

	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no parts in first candidate response")
	}
	return fmt.Sprintf("%s", res.Candidates[0].Content.Parts[0]), nil
}

// submits a query to Gemini based on a user input
func query_workflow(config ConfigFile, cwd string, cat_id int, question string, history string, model *genai.GenerativeModel, ctx context.Context) (string, error) {
	record := config.Cats[cat_id]
	img_filepath := filepath.Join(cwd, config.Metadata.ImageDirectory, record.File)
	img_file, file_err := os.Open(img_filepath)
	if file_err != nil {
		return "", errors.Join(file_err, fmt.Errorf("error finding image file %s", img_filepath))
	}
	img_bytes, bytes_err := io.ReadAll(img_file)
	if bytes_err != nil {
		return "", errors.Join(bytes_err, fmt.Errorf("error reading image file %s", img_filepath))
	}
	prompt := cat_query_prompt(img_bytes, question, history)
	// fmt.Printf("prompt: \n%s", prompt)

	res, model_err := model.GenerateContent(ctx, prompt...)
	// fmt.Printf("response parts: \n%s", res.Candidates[0].Content.Parts)
	if model_err != nil {
		return "", errors.Join(model_err, fmt.Errorf("invocation of Gemini model failed"))
	}

	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no parts in first candidate response")
	}
	return fmt.Sprintf("%s", res.Candidates[0].Content.Parts[0]), nil
}

// track the Count of API calls since the beginning of the last time window
type RateLimiterState struct {
	mu        sync.Mutex
	Count     int
	LastReset time.Time
}

// start a new time window with a count of zero
func (state *RateLimiterState) Reset() {
	state.mu.Lock()
	defer state.mu.Unlock()

	state.Count = 0
	state.LastReset = time.Now()
}

// add one to the count of API calls
func (state *RateLimiterState) Increment() {
	state.mu.Lock()
	defer state.mu.Unlock()

	state.Count += 1
}

// keeps track of the rate of usage for the Gemini API. this rate limiting is designed
// so that, if the number of API endpoints increases, it could be tracked by a map[string]RateLimiterState
func rate_limit_middleware(state *RateLimiterState, timeWindowSeconds int, rateLimit int) gin.HandlerFunc {
	return func(c *gin.Context) {

		if time.Since(state.LastReset).Seconds() > float64(timeWindowSeconds) {
			state.Reset()
		}

		if state.Count >= rateLimit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"cat_talk": "Thanks for using catgippity! Our server is currently rate-limited, please try again later. MEOW!"})
			return
		}
		state.Increment()

		c.Next()
	}
}
