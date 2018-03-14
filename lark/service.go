package lark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Service struct {
	userTokenKey string
	botTokenKey  string
	channelID    string
	url          string
}

type Param struct {
	Text    string `json:"text"`
	Token   string `json:"token"`
	Channel string `json:"channel"`
	Url     string `json:"url"`
}

func NewService(c Config) *Service {
	return &Service{
		userTokenKey: c.UserTokenKey,
		botTokenKey:  c.BotTokenKey,
		channelID:    c.ChannelID,
		url:          c.Url,
	}
}

func (s *Service) Open() error {
	return nil
}

func (s *Service) Close() error {
	return nil
}
func (s *Service) Send(message string) {
	go s.post(message)
}

func (s *Service) post(message string) {

	p := Param{}
	p.Text = message
	p.Token = s.botTokenKey
	p.Channel = s.channelID
	jsonStr, err := json.Marshal(p)
	if err != nil {
		fmt.Println("error:", err)
	}
	req, err := http.NewRequest("POST", s.url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("lark post error ", err)
		return
	}
	defer resp.Body.Close()
}
