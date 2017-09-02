package slack

import (
	"fmt"
	"log"
	"strings"
	"time"

	slackPkg "github.com/nlopes/slack"
)

type Service struct {
	tokenKey   string
	userMap    map[string]string
	channelMap map[string]string
}

func NewService(c Config) *Service {
	return &Service{
		tokenKey:   c.TokenKey,
		userMap:    make(map[string]string),
		channelMap: make(map[string]string),
	}
}

func (s *Service) Open() error {
	api := slackPkg.New(s.tokenKey)
	Users, err := api.GetUsers()
	if err != nil {
		return fmt.Errorf("Get Slack users info failed: %v", err)
	}
	for _, user := range Users {
		s.userMap[user.Profile.Email] = user.ID
	}
	Channels, err := api.GetChannels(true)
	if err != nil {
		return fmt.Errorf("Get Slack channels info failed: %v", err)
	}
	for _, channel := range Channels {
		s.channelMap[channel.Name] = channel.ID
	}
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) SendSlack(title, msg string, recovery bool, emails []string) {
	if len(emails) == 0 {
		log.Printf("[WARNING] this message has no slack user care")
		return
	}

	varUIDMap := make(map[string]string)
	varCIDMap := make(map[string]string)
	for _, vEmail := range emails {
		if strings.Contains(vEmail, "#") {
			if vCID, ok := s.channelMap[vEmail[1:]]; ok {
				varCIDMap[vEmail] = vCID
			} else {
				varCIDMap[vEmail] = ""
				log.Printf("[WARNING] no this channel id: %v in team", vEmail)
			}
		} else {
			if vUID, ok := s.userMap[vEmail]; ok {
				varUIDMap[vEmail] = vUID
			} else {
				varUIDMap[vEmail] = ""
				log.Printf("[WARNING] no user id: %v in channel", vEmail)
			}
		}
	}

	api := slackPkg.New(s.tokenKey)
	for user, v := range varUIDMap {
		if v != "" {
			go s.post(api, title, user, v, msg, recovery)
		}
	}
	for _, channelID := range varCIDMap {
		go s.PostMessage(api, channelID, title, msg, recovery)
	}
}

func (s *Service) post(api *slackPkg.Client, title, user, userID, msg string, recovery bool) {
	var tempDelay time.Duration
	var channelID string
	var err error
	for {
		_, _, channelID, err = api.OpenIMChannel(userID)
		if err != nil {
			if tempDelay == 0 {
				tempDelay = time.Second
			} else {
				tempDelay *= 2
			}
			if tempDelay > 5*time.Second {
				log.Printf("[ERROR] user[%s]open IM Channel fail finally", user)
				break
			}
			log.Printf("[ERROR] user[%s] open IM Channel fail: %v, retrying in %v", user, err, tempDelay)
			time.Sleep(tempDelay)
			continue
		}
		s.PostMessage(api, channelID, title, msg, recovery)
		break
	}
}

func (s *Service) PostMessage(api *slackPkg.Client, channelID, title, msg string, recovery bool) {
	warningType := "danger"
	if recovery {
		warningType = "good"
	}
	params := slackPkg.PostMessageParameters{
		Markdown: true,
	}
	var attachment []slackPkg.Attachment

	attachment = append(attachment, slackPkg.Attachment{
		Text:       msg,
		Color:      warningType,
		MarkdownIn: []string{"text"},
	})

	params.Attachments = attachment

	var tempDelay time.Duration
	var err error
	for {
		_, _, err = api.PostMessage(channelID, title, params)
		if err != nil {
			if tempDelay == 0 {
				tempDelay = 100 * time.Millisecond
			} else {
				tempDelay *= 2
			}
			if tempDelay > 1*time.Second {
				log.Printf("[ERROR] post message[%s] fail finally", title)
				break
			}
			log.Printf("[ERROR] post message fail: %v, retrying in %v", err, tempDelay)
			time.Sleep(tempDelay)
			continue
		}
		break
	}
}
