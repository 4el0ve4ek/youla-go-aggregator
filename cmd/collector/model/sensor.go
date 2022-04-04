package model

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Время между обращениями к датчику
const waitTime = time.Millisecond * 500

// Добавлен для таймаута, возможного прокси
var client *http.Client

func init() {
	client = &http.Client{
		Timeout: time.Millisecond * 50,
	}
}

// Sensor отвечает за один из датчиков
// по url он будет получать значение value, если что-то пошло не так, то error будет true
// cancel нужен, чтобы запросы к датчику.
type Sensor struct {
	sync.Mutex

	url    string
	value  int
	error  bool
	cancel context.CancelFunc
}

// New создает новый сенсор, объект который будет спрашивать у счетчика его значение.
func New(url string) *Sensor {
	ctx, cancel := context.WithCancel(context.Background())
	s := &Sensor{
		url:    url,
		error:  true,
		cancel: cancel,
	}

	go s.listen(ctx)
	return s
}

// listen с периодичностью waitTime пробует получить от датчика результат.
func (s *Sensor) listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(waitTime):
			s.tryGet()
		}
	}
}

// tryGet делает запрос по url, пытается из тела считать значение.
func (s *Sensor) tryGet() {
	s.reset()

	resp, err := client.Get(s.url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	if value, err := strconv.Atoi(string(body)); err == nil {
		if !Valid(value) {
			log.Println("invalid value in body")
			return
		}
		s.Lock()
		defer s.Unlock()
		s.value = value
		s.error = false
	}
}

// Обнуляет значения
func (s *Sensor) reset() {
	s.Lock()
	defer s.Unlock()
	s.value = 0
	s.error = true
}

// Get возвращает значения.
func (s *Sensor) Get() (val int, err bool) {
	s.Lock()
	defer s.Unlock()

	return s.value, s.error
}

// Stop перестает обновлять значения.
func (s *Sensor) Stop() {
	s.cancel()
}

// Valid проверяет, что значение лежит в диапазоне [0, 255]
func Valid(value int) bool {
	return 0 <= value && value <= 255
}
