package main

//go:generate sh -c "mkdir -p models && npx quicktype -s schema -l go -o models/models.go --package models --just-types-and-package ./schema/*.json"
func main() {

}
