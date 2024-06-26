package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type CatRecord struct {
	Name        string   `yaml:"name"`
	Created     string   `yaml:"created"`
	File        string   `yaml:"file"`
	Color       string   `yaml:"color"`
	Variety     string   `yaml:"variety"`
	Features    []string `yaml:"features"`
	Personality []string `yaml:"personality"`
	Prompt      string   `yaml:"prompt"`
	Welcome     string   `yaml:"welcome"`
	Comment     string   `yaml:"comment"`
}

type ConfigFile struct {
	Metadata struct {
		BaseUrl        string `yaml:"base_url"`
		ImageDirectory string `yaml:"image_directory"`
	} `yaml:"metadata"`
	Cats []CatRecord
}

func load_config(cwd string) (ConfigFile, error) {
	config := ConfigFile{}
	cats_filename := filepath.Join(cwd, "config.yaml")
	cats_file, err := os.Open(cats_filename)
	if err != nil {
		return config, errors.Join(fmt.Errorf("failure finding server config file"), err)
	}
	defer cats_file.Close() // Close the file on exit
	cats_str, err := io.ReadAll(cats_file)
	if err != nil {
		return config, errors.Join(fmt.Errorf("failure reading server config file"), err)
	}

	err2 := yaml.Unmarshal(cats_str, &config)
	if err2 != nil {
		return config, errors.Join(fmt.Errorf("failure deserializing config file"), err2)
	}
	return config, nil
}

func get_cat_image_url(config ConfigFile, cwd string, cat_id int) (string, error) {
	if cat_id < 0 || cat_id >= len(config.Cats) {
		return "", fmt.Errorf("invalid cat_id %d", cat_id)
	}
	record := config.Cats[cat_id]
	url := fmt.Sprintf("/img/cat/%s", record.File)
	return url, nil
}

func get_cat_image_bytes(config ConfigFile, cwd string, cat_id int) ([]byte, error) {
	record := config.Cats[cat_id]
	img_filepath := filepath.Join(cwd, config.Metadata.ImageDirectory, record.File)
	img_file, file_err := os.Open(img_filepath)
	if file_err != nil {
		return nil, errors.Join(file_err, fmt.Errorf("error finding image file %s", img_filepath))
	}
	img_bytes, bytes_err := io.ReadAll(img_file)
	if bytes_err != nil {
		return nil, errors.Join(bytes_err, fmt.Errorf("error reading image file %s", img_filepath))
	}
	return img_bytes, nil
}
