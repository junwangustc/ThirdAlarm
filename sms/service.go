package sms

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/junwangustc/ustclog"
)

type Service struct {
	client    http.Client
	senderKey string
	expandKey string
	token     string
	sendUrl   string
}

func NewService(c Config) *Service {
	return &Service{
		client:    http.Client{},
		senderKey: c.SenderKey,
		token:     c.Token,
		sendUrl:   c.SendUrl,
		expandKey: c.Expandkey,
	}
}

func (s *Service) Open() error {
	s.client.Timeout = time.Duration(30) * time.Second
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) getAuthSign(ts, receiver, message string) string {
	raw := []byte(fmt.Sprintf("%s%s%s%s", s.expandKey, ts, receiver, message))
	h := sha256.New()
	h.Write(raw)
	src := h.Sum(nil)
	sha := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(sha, src)
	n := len(s.token)
	for i := 0; i < len(sha); i++ {
		sha[i] ^= byte(s.token[i%n])
	}
	return string(base64.StdEncoding.EncodeToString(sha))
}

func (s *Service) getSmsRequest(receiver, message string) (*http.Request, error) {
	urlEncodeMsg := url.QueryEscape(message)
	now := strconv.Itoa(int(time.Now().Unix()))
	auth := s.getAuthSign(now, receiver, urlEncodeMsg)
	urlValue := url.Values{
		"receiver":   {receiver},
		"message":    {urlEncodeMsg},
		"sender_key": {s.senderKey},
		"timestamp":  {now},
		"auth":       {auth},
	}

	request, err := http.NewRequest("POST", s.sendUrl, strings.NewReader(urlValue.Encode()))
	if err != nil {
		return request, fmt.Errorf("get sms request fail: %v", err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", fmt.Sprintf("%v", len(urlValue)))

	return request, nil
}

func (s *Service) sendMessage(receiver, message string) error {
	request, err := s.getSmsRequest(receiver, message)
	if err != nil {
		return fmt.Errorf("new requests for %v fail: %v", receiver, err)
	}

	response, err := s.client.Do(request)
	if err != nil {
		return fmt.Errorf("do request fail: %v", err)
	}
	defer response.Body.Close()

	byteBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Failed on reading response body: %s", err)
	}
	var f interface{}
	err = json.Unmarshal(byteBody, &f)
	if err != nil {
		return fmt.Errorf("Unmarshal sms response body %s fail: %v", string(byteBody), err)
	}
	responseMap := f.(map[string]interface{})

	if response.StatusCode != 200 {
		return fmt.Errorf("Response status: %v", response.Status)
	}

	if responseMap["msg"] != "success" {
		return fmt.Errorf("Send message to %s fail", receiver)
	}
	return nil
}

func (s *Service) Send(receivers []string, message string) {
	if len(receivers) == 0 {
		log.Warn("this message has no sms user care")
		return
	}
	for _, receiver := range receivers {
		err := s.sendMessage(receiver, message)
		if err != nil {
			log.Errorf("send sms to %s fail: %v", receiver, err)
		} else {
			log.Printf("send sms to %s success", receiver)
		}
	}
}
