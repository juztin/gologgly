// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gologgly

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"code.minty.io/config"
)

type Logger struct {
	url    string
	client *http.Client
}

var (
	logger *Logger
)

func init() {
	url := config.RequiredGroupString("loggly", "url")
	logger = New(url)
}

func New(url string) *Logger {
	return &Logger{url, &http.Client{}}
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

	body := bytes.NewBuffer(b)
	resp, err := l.client.Post(l.url, "application/json", body)
	if err != nil {
		return fmt.Errorf("failed posting to Loggly: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status from Loggly:", resp.StatusCode, ":", string(b))
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
	if logger.url == "" {
		return errors.New("loggly URL doesn't exist in 'config.json'")
	}
	return logger.send(d)
}

func Error(err error) error {
	if logger.url == "" {
		return errors.New("loggly URL doesn't exist in 'config.json'")
	}
	return logger.send(err.Error())
}
