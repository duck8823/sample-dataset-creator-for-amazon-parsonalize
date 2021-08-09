package main

import (
	"github.com/duck8823/sample-dataset-creator-for-amazon-personalize/creator"
	"log"
)

//go:generate sh -c "mkdir -p models && npx quicktype -s schema -l go -o models/models.go --package models --just-types-and-package ./schema/*.json"
func main() {
	c := &creator.CsvCreator{
		Output: "./output",
	}
	if err := c.Create(); err != nil {
		log.Fatalf("失敗しました: %#v", err)
	}
}
