package eml

import (
	"bytes"
	"io"
	"log"
	"net/mail"
	"os"
)

type Mail struct {
	From    string
	To      string
	Subject string
	Body    string
}

func ParseFile(filePath string) (Mail, error) {
	r, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
		return (Mail{}), err
	}

	b := bytes.NewReader(r)

	m, err := mail.ReadMessage(b)
	if err != nil {
		log.Fatal(err)
	}

	header := m.Header

	body, err := io.ReadAll(m.Body)
	if err != nil {
		log.Fatal(err)
		return (Mail{}), err
	}

	return (Mail{
		From:    header.Get("From"),
		To:      header.Get("To"),
		Subject: header.Get("Subject"),
		Body:    string(body),
	}), nil

}
