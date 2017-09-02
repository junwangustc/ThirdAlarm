package main

import (
	"github.com/junwangustc/ThirdAlarm/alert"
	"github.com/junwangustc/ThirdAlarm/slack"
	"github.com/junwangustc/ThirdAlarm/sms"
	log "github.com/junwangustc/ustclog"
)

type Service interface {
	Open() error
	Close() error
}
type Server struct {
	SlackService *slack.Service
	AlertService *alert.Service
	SMSService   *sms.Service
	services     []Service
	cfg          *Config
}

func (s *Server) appendSMSService(c sms.Config) {
	if c.Enabled {
		smsService := sms.NewService(c)
		s.SMSService = smsService
		s.services = append(s.services, smsService)
	}
}

func (s *Server) appendSlackService(c slack.Config) {
	if c.Enabled {
		slackService := slack.NewService(c)
		s.SlackService = slackService
		s.services = append(s.services, slackService)
	}
}

func (s *Server) appendAlertService(c alert.Config) {
	alertService := alert.NewService(c)
	alertService.SlackService = s.SlackService
	alertService.SMSService = s.SMSService
	s.AlertService = alertService
	s.services = append(s.services, alertService)
}
func (s *Server) Open() error {
	if err := func() error {
		for _, service := range s.services {
			log.Println("start Opening service %T", service)
			if err := service.Open(); err != nil {
				return err
			}
			log.Println("opened Service %T:%s", service)
		}
		return nil
	}(); err != nil {
		s.Close()
		return err
	}
	return nil
}

func (s *Server) Close() {
	for _, service := range s.services {
		log.Printf("D! closing service: %T", service)
		err := service.Close()
		if err != nil {
			log.Printf("E! error closing service %T: %v", service, err)
		}
		log.Printf("D! closed service: %T", service)
	}
}

func NewServer(cfg *Config) (srv *Server) {
	srv = &Server{cfg: cfg}
	return srv
}
func (s *Server) Run() (err error) {

	s.appendSlackService(s.cfg.Slack)
	s.appendSMSService(s.cfg.SMS)
	s.appendAlertService(s.cfg.Alert)
	if err = s.Open(); err != nil {
		log.Println("[Error]open server fail", err)
		return err
	}
	return nil
}
