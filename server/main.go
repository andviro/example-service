package main

import (
	"log"
)

func main() {
	app := AppFactory(":8080")
	log.Fatal(app.Run())
}
