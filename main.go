package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

var EnvDB string
var EnvDBPath string

func main() {
	// загрузка окружения из .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	envPort := os.Getenv("TODO_PORT")
	EnvDB = os.Getenv("DB_FILE")
	EnvDBPath = os.Getenv("DB_PATH")
	_, err = New(EnvDBPath, EnvDB)
	if err != nil {
		panic(err)
	}

	// маршрутизатор
	r := chi.NewRouter()

	// загрузка файлов сервера
	r.Mount("/", http.FileServer(http.Dir("./web/")))

	// api
	r.Get("/api/nextdate", handlerApiNextDate)
	r.Get("/api/nextdate?now={now}", handlerApiNextDate)
	r.Post("/api/task", handlerApiTaskPost)
	r.Get("/api/tasks", handlerApiTaskGet)
	r.Get("/api/task", handlerApiTaskEdit)

	// запуск сервера
	port := fmt.Sprintf(":%s", envPort)
	err = http.ListenAndServe(port, r)
	if err != nil {
		panic(err)
	}
}
