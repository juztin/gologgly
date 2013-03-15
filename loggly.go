package gologgly

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"bitbucket.org/juztin/wombat/config"
)

type Logger struct {
	url    string
	client *http.Client
}

var (
	url    = "https://logs.loggly.com/inputs/"
	logger *Logger
)

func init() {
	key, _ := config.GroupString("loggly", "key")
	logger = New(key)
}

func New(key string) *Logger {
	u := url
	if key != "" {
		u = u + key
	}
	return &Logger{u, &http.Client{}}
}

func (l *Logger) send(o interface{}) error {
	var b []byte
	if s, ok := o.(string); ok {
		b = []byte(s)
	} else if j, err := json.Marshal(o); err != nil {
		return err
	} else {
		b = j
	}

	if l.url == "" {
		log.Printf("Loggly:%v\n", string(b))
		return nil
	}

	body := bytes.NewBuffer(b)
	resp, err := l.client.Post(l.url, "application/json", body)
	if err != nil {
		log.Println("[ERROR] Loggly:", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Println("[ERROR] Loggly: Invalid status,", resp.StatusCode, ":", string(b))
		return errors.New(fmt.Sprintf("Failed to log to Loggly, %d", resp.StatusCode))
	}

	return nil
}

func (l *Logger) Log(d interface{}) error {
	return l.send(d)
}

func (l *Logger) Error(err error) error {
	return l.send(err.Error())
}

func Log(d interface{}) error {
	return logger.send(d)
}

func Error(err error) error {
	return logger.send(err.Error())
}

