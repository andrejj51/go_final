package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// repeat пустой
func repeatNill(date time.Time) (string, error) {
	if date.After(time.Now()) {
		return date.Format("20060102"), nil
	}
	return time.Now().Format("20060102"), errors.New("так как указанная дата меньше текущей, то в задачу будет записана теекущая дата")
}

// d
func repeatD(now time.Time, date time.Time, repeat string) (string, error) {
	repeatList := strings.Split(repeat, " ")
	if len(repeatList) == 2 {
		i, err := strconv.ParseInt(repeatList[1], 10, 64)
		if err != nil {
			return "", errors.New("указан неверный формат repeat d err (не указан интервал в днях)")
		}
		if i < 1 || i > 400 {
			return "", errors.New("указан неверный формат repeat d err (выход за рамки интервала)")
		}
		fmt.Println("d", i)
		t, err := d(now, date, i)
		if err != nil {
			return "", err
		}
		return t, nil // d n , где 1 <= n <= 400
	}
	return "", errors.New("в колонке repeat — d указан неверный формат") // d || d n n+1 ...
}

// y
func repeatY(now time.Time, date time.Time, repeat string) (string, error) {
	t, err := y(now, date)
	if err != nil {
		return "", err
	}
	return t, nil
}

// repeat всё остальное
func repeatDefault() (string, error) {
	return "", errors.New("неверный формат")
}

// data parse
func dataParse(date string) (time.Time, error) {
	v, err := time.Parse("20060102", date)
	if err != nil {
		log.Println(err)
		return time.Now(), err
	}
	return v, nil
}

// 1 <= d <= 7
func d(now time.Time, date time.Time, dateInt int64) (string, error) {
	next := date.AddDate(0, 0, int(dateInt))
	if next.After(now) {
		return next.Format("20060102"), nil
	}

	return "", errors.New("возвращаемая дата меньше или равна текущей")
}

// у
func y(now time.Time, date time.Time) (string, error) {
	next := date.AddDate(1, 0, 0)
	if next.After(now) {
		return next.Format("20060102"), nil
	}

	return "", errors.New("возвращаемая дата меньше или равна текущей")
}

// repeat rules
func NextDate(now time.Time, date string, repeat string) (string, error) {
	_, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("время в переменной date не может быть преобразовано в корректную дату")
	}

	switch {
	case repeat == "":
		v, err := dataParse(date)
		if err != nil {
			log.Println(err)
		}
		return repeatNill(v)
	case []rune(repeat)[0] == 'd' && []rune(repeat)[1] == ' ':
		v, err := dataParse(date)
		if err != nil {
			log.Println(err)
		}
		return repeatD(now, v, repeat)
	case []rune(repeat)[0] == 'y' && len([]rune(repeat)) == 1:
		v, err := dataParse(date)
		if err != nil {
			log.Println(err)
		}
		return repeatY(now, v, repeat)
	case []rune(repeat)[0] == 'w' && []rune(repeat)[1] == ' ':
		return "", errors.New("в разработке ")
		//return repeatW(now, dataParse(date), repeat)
	case []rune(repeat)[0] == 'm' && []rune(repeat)[1] == ' ':
		return "", errors.New("в разработке ")
		//return repeatM(repeat)
	}
	return "", errors.New("непредвиденная ошиька")
}
