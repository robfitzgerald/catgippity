package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func main() {

	dir, err := os.Getwd()
	if err != nil {
		// do
		// return err
		fmt.Println("Error getting current working directory:", err)
		return
	}
	cats_filename := filepath.Join(dir, "config.yaml")
	// Open the file
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

	cats_yaml := ConfigFile{}
	err2 := yaml.Unmarshal(cats_str, &cats_yaml)
	if err2 != nil {
		log.Fatalf("error: %v", err)
		return
	}
	// fmt.Printf("--- CatsFile:\n%v\n\n", cats_yaml)

	// d, err := yaml.Marshal(&cats_yaml)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// fmt.Printf("--- t dump:\n%s\n\n", string(d))

	r := gin.Default()
	r.GET("/config", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": cats_yaml,
		})
	})
	r.Run() // http://0.0.0.0:8080/config

}
