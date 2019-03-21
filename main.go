package main

import "github.com/adeo/turbine-go-api-skeleton/cmd"

//go:generate go run scripts/includeopenapi.go

func main() {
	cmd.Execute()
}
