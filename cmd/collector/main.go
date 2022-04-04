package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"collector/model"
	"github.com/gorilla/mux"
)

// url для http запросов на датчики
var servers = []string{
	`http://localhost:8081/`,
	`http://localhost:8082/`,
	`http://localhost:8083/`,
	`http://localhost:8084/`,
}

// Массив объектов, которые слушают датчики
var sensors = make([]*model.Sensor, 0, len(servers))

func main() {
	if s := os.Getenv("server_urls"); s != "" {
		servers = strings.Split(s, " ")
	}

	port := ":8080"
	if s := os.Getenv("port"); s != "" {
		port = s
	}

	for _, url := range servers {
		sensors = append(sensors, model.New(url))
	}

	r := mux.NewRouter()
	SetRoutes(r)
	err := http.ListenAndServe(port, r)
	if err != nil {
		fmt.Println(err)
	}
}

// SetRoutes устанавливает "ручки" на роутер
func SetRoutes(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		v, e := Evaluate(sensors)
		_, err := fmt.Fprintf(w, "%d/%d", v, e)
		if err != nil {
			fmt.Println(err)
		}
	})
}

// Evaluate считает среднее по всем значениям датчиков, до которых получилось достучаться.
// Возвращает значение, которое получилось, и количество сенсоров с ошибкой.
func Evaluate(sensors []*model.Sensor) (value int, errors int) {
	valid := 0
	for _, sensor := range sensors {
		val, err := sensor.Get()
		if err {
			errors++
		} else {
			value += val
			valid++
		}
	}
	if valid != 0 {
		value /= valid
	}
	return
}
