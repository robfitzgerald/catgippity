package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	ctx := context.Background()

	// static website configuration
	dir, err := os.Getwd()
	config, conf_err := load_config(dir)
	if conf_err != nil {
		log.Fatalf("error: %v", err)
		return
	}

	// load Gemini Client and Model
	apiKey := os.Getenv("API_KEY")
	client, err := load_client(ctx, apiKey)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Close()
	model := client.GenerativeModel(config.Metadata.ModelName)

	// configure host/port for server
	r := gin.Default()
	server_cors := cors.DefaultConfig()
	server_port := os.Getenv("PORT")
	server_host, server_host_exists := os.LookupEnv("HOST")
	if server_host_exists {
		log.Printf("running on %s:%s\n", server_host, server_port)
		server_cors.AllowOrigins = []string{server_host}
	} else if server_port != "8080" {
		log.Printf("running on http://localhost:%s\n", server_port)
		server_cors.AllowOrigins = []string{fmt.Sprintf("http://localhost:%s", server_port)} // testing
	} else {
		log.Printf("running on http://localhost:8080\n")
	}
	r.Use(cors.New(server_cors))

	// apply in-memory rate limiting
	state := RateLimiterState{}
	test_rate_limiting := rate_limit_middleware(&state, config.Metadata.TimeWindowSeconds, config.Metadata.RateLimit)

	// server file access patterns
	r.LoadHTMLFiles("index.html")
	r.Static("/js", "./pub/js")
	r.Static("/img", "./pub/img")
	r.StaticFile("/favicon.ico", "./favicon.ico")

	// routes
	r.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": config,
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/custom-welcome/*cat_id", test_rate_limiting, func(c *gin.Context) {
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
			welcome_str, welcome_err := welcome(config, cat_id)
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

	r.POST("/query", test_rate_limiting, func(c *gin.Context) {
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

	r.Run()

}
