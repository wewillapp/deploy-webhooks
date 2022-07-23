package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Request struct {
	ProjectName string `json:"projectName"`
	Token       string `json:"token"`
}
type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

func Main(req Request) (*Response, error) {
	if os.Getenv("AUTH_TOKEN") == "" || os.Getenv("BOT_TOKEN") == "" {
		return nil, errors.New("cannot load config files")
	}
	if !validateToken(req.Token) {
		return nil, errors.New("invalid token")
	}
	if err := sendNotifyMessage(req); err != nil {
		return nil, err
	}
	return &Response{
		Body: "✅",
	}, nil
}

func sendNotifyMessage(request Request) error {
	form := url.Values{}
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return err
	}
	t := time.Now().In(loc).Format("2006-01-02 15:04:05")
	replyMsg := fmt.Sprintf("\n\n🚀 Project %s has been deployed at %s", request.ProjectName, t)
	form.Add("message", replyMsg)
	req, _ := http.NewRequest(
		http.MethodPost,
		"https://notify-api.line.me/api/notify",
		strings.NewReader(form.Encode()),
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer "+os.Getenv("BOT_TOKEN"))
	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func validateToken(token string) bool {
	authToken := os.Getenv("AUTH_TOKEN")
	return token == authToken
}
