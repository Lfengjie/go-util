package util

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

func sendEmailBySmtp(subj string, body string, fromEmail string, password string, toEmail []string) {
	// Connect to the SMTP Server
	var toEmailStr string = ""
	for i, v := range toEmail {
		if i == 0 {
			toEmailStr = toEmail[i]
			continue
		}
		toEmailStr = fmt.Sprintf("%s;%s", toEmailStr, v)
	}

	headers := make(map[string]string)
	headers["From"] = fromEmail
	headers["To"] = toEmailStr
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	servername := "mail.didichuxing.com:587"
	host, _, _ := net.SplitHostPort(servername)
	auth := LoginAuth(fromEmail, password)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	c, err := smtp.Dial(servername)
	if err != nil {
		log.Panic(err)
	}

	c.StartTLS(tlsconfig)

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// To && From
	if err = c.Mail(fromEmail); err != nil {
		log.Panic(err)
	}

	for _, v := range toEmail {
		if err = c.Rcpt(v); err != nil {
			log.Panic(err)
		}
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	// var message string = "dumy test message"
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}
