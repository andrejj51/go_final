package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// загрузка окружения из .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	envPort := os.Getenv("TODO_PORT")

	// загрузка файлов сервера
	http.Handle("/", http.FileServer(http.Dir("./web/")))

	// запуск сервера
	port := fmt.Sprintf(":%s", envPort)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
