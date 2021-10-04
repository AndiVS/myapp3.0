package main

import (
	"fmt"
	"myapp3.0/internal/start"
)

func main() {
	fmt.Println("Welcome to the webserver")

	e := start.New()

	e.Start(":8000")
}
