package main

import "github.com/adeo/go-api-skeleton/cmd"

//go:generate go run scripts/includeopenapi.go

func main() {
	cmd.Execute()
}
