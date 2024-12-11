package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

// база
func db_file(envDB string) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), envDB)
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	if install == true {
		db, err := sql.Open("sqlite", envDB)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer db.Close()

		type task struct {
			id      int
			date    string
			title   string
			comment string
			repeat  string
		}

		table := `CREATE TABLE IF NOT EXISTS scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date CHAR(8),
				title VARCHAR(256) NOT NULL DEFAULT "",
				comment TEXT NOT NULL DEFAULT "",
				repeat VARCHAR(128) NOT NULL DEFAULT ""
			)`

		//index := "CREATE INDEX date ON scheduler (column date)"

		_, err = db.Exec(table)
		if err != nil {
			log.Fatal(err)
		}
		/*
			CREATE TABLE scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date CHAR(8),
				title VARCHAR(256) NOT NULL DEFAULT "",
				comment TEXT NOT NULL DEFAULT "",
				repeat VARCHAR(128) NOT NULL DEFAULT ""
			);
			CREATE INDEX date ON scheduler (column date);*/
	}
}

func main() {
	// загрузка окружения из .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	envPort := os.Getenv("TODO_PORT")
	envDB := os.Getenv("DB_FILE")

	db_file(envDB)

	// база
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), envDB)
	_, err = os.Stat(dbFile)

	// если install равен true, после открытия БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX

	// загрузка файлов сервера
	http.Handle("/", http.FileServer(http.Dir("./web/")))

	// запуск сервера
	port := fmt.Sprintf(":%s", envPort)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
