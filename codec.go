package main

type CatRecord struct {
	Name        string   `yaml:"name"`
	Created     string   `yaml:"created"`
	File        string   `yaml:"file"`
	Color       string   `yaml:"color"`
	Variety     string   `yaml:"variety"`
	Features    []string `yaml:"features"`
	Personality []string `yaml:"personality"`
	Prompt      string   `yaml:"prompt"`
	Comment     string   `yaml:"comment"`
}

type ConfigFile struct {
	Metadata struct {
		ImageDirectory string `yaml:"image_directory"`
	} `yaml:"metadata"`
	Cats []CatRecord
}
