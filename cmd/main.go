package main

import (
	"fmt"
	"goCounter/app"
	"goCounter/db"
	"log"
	"os"
	"path/filepath"
)

func getExecPath() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal("Error getting exec path")
	}
	exPath := filepath.Dir(ex)
	return exPath
}

func run() error {
	// path := getExecPath()
	// db := db.Init(path + "/gocounter/db.txt")
	db := db.Init("./db.txt")
	counterApp := app.CreateApp(db)
	if _, err := counterApp.Run(); err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Print("\033[H\033[2J")
	if err := run(); err != nil {
		log.Fatal("Error excecuting app", err)
	}
}
