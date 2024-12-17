package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Storage struct {
	db *sql.DB
}

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Id struct {
	Id int `json:"id"`
}

type Error struct {
	Error string `json:"error"`
}

// конструктор
func New(storagePath string, dbFileName string) (*Storage, error) {
	absPath, err := filepath.Abs(storagePath)
	if err != nil {
		log.Fatal(err)
	}
	//dbFile := filepath.Join(filepath.Dir(appPath), dbFileName)
	//наличие файла
	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		// Если файл не существует, создаем его
		file, err := os.Create(absPath)
		if err != nil {
			log.Println("Ошибка создания файла:", err)
			return nil, err
		}
		defer file.Close()
	}

	db, err := sql.Open("sqlite3", dbFileName) // Подключаемся к БД
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	//defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8),
		title VARCHAR(256) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT "");
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`)

	// Создаем таблицу, если ее еще нет
	/*stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS scheduler(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8),
		title VARCHAR(256) NOT NULL DEFAULT "",
		comment TEXT NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT "");
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
	`)*/
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	/*_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}*/

	return &Storage{db: db}, nil

}

// Add
func (s Storage) Add(t Task) (int, error) {
	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))

	if err != nil {
		return 0, err
	}
	// идентификатор последней добавленной записи
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(lastId), nil
}
