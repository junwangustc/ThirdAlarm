package alert

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/junwangustc/ThirdAlarm/slack"
	"github.com/junwangustc/ThirdAlarm/sms"
)

type Service struct {
	Cfg          *Config
	SlackService *slack.Service
	SMSService   *sms.Service
	UrlStatus    map[string]*URL
	AlarmStatus  map[string]int
	ExitChan     chan int
}

func NewService(cfg Config) *Service {
	return &Service{Cfg: &cfg, UrlStatus: make(map[string]*URL), AlarmStatus: make(map[string]int), ExitChan: make(chan int, 1)}
}

func (s *Service) Open() error {
	for index, url := range s.Cfg.LiveUrl {
		s.UrlStatus[url.Name] = &s.Cfg.LiveUrl[index]
		s.AlarmStatus[url.Name] = -1
	}
	go s.Run()
	return nil
}

func (s *Service) Close() error {
	s.ExitChan <- 1
	return nil
}
func (s *Service) Run() {
	tick := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-tick.C:
			for k, v := range s.UrlStatus {
				if err := s.Ping(v.Url); err != nil {
					if s.AlarmStatus[k] == -1 {
						//发告警
						s.AlarmStatus[k] = 1
						go s.SendSlack("存活监控异常", k+"监控连接失败", false, v.CareSlack)
						go s.SendSMS(k+"监控连接失败", v.CareSMS)
						log.Println("SMS debug", v.CareSMS)

					}
				} else {
					if s.AlarmStatus[k] == 1 {
						//发告警
						s.AlarmStatus[k] = -1
						go s.SendSlack("存活监控恢复正常", k+"监控连接恢复正常", true, v.CareSlack)
						go s.SendSMS(k+"监控连接恢复正常", v.CareSMS)
						log.Println("SMS debug", v.CareSMS)
					}

				}

			}

		case <-s.ExitChan:
			return
		}
	}
}
func (s *Service) Ping(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
		log.Println(err)
		return err
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body) //请求数据进行读取
	_ = err
	return nil
}

func (s *Service) SendSlack(msg, title string, status bool, emails []string) {
	s.SlackService.SendSlack(msg, title, status, emails)

}
func (s *Service) SendSMS(msg string, smsNumer []string) {
	s.SMSService.Send(smsNumer, msg)

}
