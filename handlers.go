package main

import (
	"encoding/json"
	"log"
	"net/http"
)

//envDB := os.Getenv("DB_FILE")
//envDBPath := os.Getenv("DB_PATH")

// api/nextdate
func handlerApiNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	str, err := NextDate(dataParse(now), date, repeat)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(str))
}

// api/task
func handlerApiTask(w http.ResponseWriter, r *http.Request) {
	storage, err := New(EnvDBPath, EnvDB)
	if err != nil {
		panic(err)
	}

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Println("Ошибка при декодировании JSON:", err)
		return
	}
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
