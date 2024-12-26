package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
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
	Id string `json:"id"`
}

type TaskAndId struct {
	Id      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type Tasks struct {
	Tasks []TaskAndId `json:"tasks"`
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

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &Storage{db: db}, nil

}

// Add
func (s Storage) Add(t Task) (string, error) {
	res, err := s.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))

	if err != nil {
		return "", err
	}
	// идентификатор последней добавленной записи
	lastId, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	return string(lastId), nil
}

// Get
func (s Storage) Get() (Tasks, error) {
	rows, err := s.db.Query("SELECT * FROM scheduler WHERE date >= :now ORDER BY date ASC LIMIT 20", sql.Named("now", time.Now().Format("20060102")))

	var task TaskAndId
	var tasks []TaskAndId

	if err != nil {
		return Tasks{Tasks: []TaskAndId{}}, err
	}
	defer rows.Close()

	for rows.Next() {
		task = TaskAndId{}
		err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return Tasks{Tasks: []TaskAndId{}}, err
		}
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		res := Tasks{Tasks: []TaskAndId{}}
		return res, nil
	}
	res := Tasks{Tasks: tasks}
	return res, nil
}

// GetId
func (s Storage) GetId(id string) (TaskAndId, error) {
	row := s.db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", id))

	var task TaskAndId = TaskAndId{}

	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		return TaskAndId{}, err
	}
	return task, nil
}
