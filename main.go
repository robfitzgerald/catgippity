package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"
)

func main() {

	ctx := context.Background()

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}
	cats_filename := filepath.Join(dir, "config.yaml")
	cats_file, err := os.Open(cats_filename)
	if err != nil {
		fmt.Println("Error finding file data/yaml/cats.yaml:", err)
		return
	}
	defer cats_file.Close() // Close the file on exit
	cats_str, err := io.ReadAll(cats_file)
	if err != nil {
		fmt.Println("Error reading file data/yaml/cats.yaml:", err)
		return
	}

	config := ConfigFile{}
	err2 := yaml.Unmarshal(cats_str, &config)
	if err2 != nil {
		log.Fatalf("error: %v", err)
		return
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Close()

	// The Gemini 1.5 models are versatile and work with most use cases
	model := client.GenerativeModel("gemini-1.5-flash")
	// model := client.GenerativeModel("gemini-1.5-pro")

	r := gin.Default()

	default_conf := cors.DefaultConfig()
	server_port := os.Getenv("PORT")
	if server_port != "8080" {
		default_conf.AllowOrigins = []string{fmt.Sprintf("http://localhost:%s", server_port)} // testing
	}

	r.Use(cors.New(default_conf))

	r.LoadHTMLFiles("index.html")
	r.Static("/js", "./pub/js")
	r.Static("/img", "./pub/img")
	r.StaticFile("/favicon.ico", "./favicon.ico")

	r.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": config,
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/custom-welcome/*cat_id", func(c *gin.Context) {
		cat_id_str := strings.ReplaceAll(c.Param("cat_id"), "/", "")
		if cat_id_str == "" {
			cat_id_str = "0"
		}
		cat_id, cat_id_err := strconv.Atoi(cat_id_str)
		if cat_id_err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": cat_id_err.Error(),
			})
		} else {
			res, welcome_err := welcome_workflow(config, dir, cat_id, model, ctx)
			img_url, img_err := get_cat_image_url(config, dir, cat_id)
			if welcome_err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": welcome_err.Error(),
				})
			} else if img_err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": img_err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"cat_talk":  res,
					"image_url": img_url,
				})
			}
		}
	})

	r.GET("/welcome/*cat_id", func(c *gin.Context) {
		cat_id_str := strings.ReplaceAll(c.Param("cat_id"), "/", "")
		if cat_id_str == "" {
			cat_id_str = "0"
		}
		cat_id, cat_id_err := strconv.Atoi(cat_id_str)
		if cat_id_err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": cat_id_err.Error(),
			})
		} else {
			welcome_str, welcome_err := welcome(config, dir, cat_id)
			img_url, img_err := get_cat_image_url(config, dir, cat_id)
			if welcome_err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": welcome_err.Error(),
				})
			} else if img_err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": img_err.Error(),
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"cat_talk":  welcome_str,
					"image_url": img_url,
				})
			}
		}
	})

	r.POST("/query", func(c *gin.Context) {

		type QueryRequestBody struct {
			Question string `json:"question" binding:"required"`
			History  string `json:"history"`
		}

		var body QueryRequestBody
		bind_err := c.BindJSON(&body)
		if bind_err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": bind_err.Error()})
		}
		res, query_err := query_workflow(config, dir, 1, body.Question, body.History, model, ctx)

		if query_err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": query_err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"cat_talk": res,
			})
		}
	})

	r.Run() // http://0.0.0.0:8080/config

}
