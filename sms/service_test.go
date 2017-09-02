package sms

import "testing"

func fakeConfig() Config {
	c := NewConfig()
	c.SendUrl = "fake-url"
	c.SenderKey = "fake-sender-key"
	c.Token = "faker-token"
	return c
}

func TestSendSMS(t *testing.T) {
	receiver := "phone"
	message := "alarm测试"
	config := fakeConfig()
	service := NewService(config)
	err := service.sendMessage(receiver, message)
	if err != nil {
		t.Errorf("send sms fail: %v", err)
	}
}
