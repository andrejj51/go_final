package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

//envDB := os.Getenv("DB_FILE")
//envDBPath := os.Getenv("DB_PATH")

// api/nextdate
func handlerApiNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	v, err := dataParse(now)
	if err != nil {
		log.Println(err)
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	str, err := NextDate(v, date, repeat)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(str))
}

// api/task
func handlerApiTaskPost(w http.ResponseWriter, r *http.Request) {
	storage, err := New(EnvDBPath, EnvDB)
	if err != nil {
		panic(err)
	}

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Println("Ошибка при декодировании JSON:", err)
		return
	}

	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	v, err := time.Parse("20060102", task.Date)
	if err != nil {
		log.Println(err)
	}
	if v.Before(time.Now()) {
		if task.Repeat == "" {
			task.Date = time.Now().Format("20060102")
		}
		if task.Repeat != "" {
			d, err := NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				log.Println(err)
			}
			task.Date = d
		}
	}
	d, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		log.Println(err)
		e := Error{Error: err.Error()}
		// Кодирование ошибки в JSON
		encodedJSON, err := json.Marshal(e)
		if err != nil {
			log.Println("Ошибка при кодировании JSON:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//w.WriteHeader(http.StatusCreated)
		//json.NewEncoder(w).Encode(encodedJSON)
		w.Write(encodedJSON)
		return
	}
	task.Date = d

	id, err := storage.Add(task)

	if err != nil {
		e := Error{Error: err.Error()}
		// Кодирование ошибки в JSON
		encodedJSON, err := json.Marshal(e)
		if err != nil {
			log.Println("Ошибка при кодировании JSON:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//w.WriteHeader(http.StatusCreated)
		//json.NewEncoder(w).Encode(encodedJSON)
		w.Write(encodedJSON)
		return
	}
	i := Id{Id: id}
	// Кодирование объекта в JSON
	encodedJSON, err := json.Marshal(i)
	if err != nil {
		log.Println("Ошибка при кодировании JSON:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(encodedJSON)
	w.Write(encodedJSON)

}

func handlerApiTaskGet(w http.ResponseWriter, r *http.Request) {
	storage, err := New(EnvDBPath, EnvDB)
	if err != nil {
		panic(err)
	}

	tasks, err := storage.Get()
	if err != nil {
		e := Error{Error: err.Error()}
		// Кодирование ошибки в JSON
		encodedJSON, err := json.Marshal(e)
		if err != nil {
			log.Println("Ошибка при кодировании JSON:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		//w.WriteHeader(http.StatusCreated)
		//json.NewEncoder(w).Encode(encodedJSON)
		w.Write(encodedJSON)
		return
	}
	// Кодирование объекта в JSON
	encodedJSON, err := json.Marshal(tasks)
	if err != nil {
		log.Println("Ошибка при кодировании JSON:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(encodedJSON)
	w.Write(encodedJSON)
}

func handlerApiTaskEdit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	storage, err := New(EnvDBPath, EnvDB)
	if err != nil {
		panic(err)
	}

	task, err := storage.GetId(id)
	if err != nil {
		e := Error{Error: errors.New("Задача не найдена").Error()}
		// Кодирование ошибки в JSON
		encodedJSON, err := json.Marshal(e)
		if err != nil {
			log.Println("Ошибка при кодировании JSON:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(encodedJSON)
		return
	}
	// Кодирование объекта в JSON
	encodedJSON, err := json.Marshal(task)
	if err != nil {
		log.Println("Ошибка при кодировании JSON:", err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(encodedJSON)
}
