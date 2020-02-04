package serializer

import (
	"fmt"
	url2 "net/url"
	"strings"

	"github.com/pkg/errors"
)

type Subscription struct {
	Alias     string   `json:"alias"`
	BaseURL   string   `json:"baseURL"`
	SpaceKey  string   `json:"spaceKey"`
	Events    []string `json:"events"`
	ChannelID string   `json:"channelID"`
}

var eventTypes = map[string]string{
	"comment_created": "Comment Create",
	"comment_updated": "Comment Update",
	"comment_removed": "Comment Remove",
	"page_created":    "Page Create",
	"page_updated":    "Page Update",
	"page_trashed":    "Page Trash",
	"page_restored":   "Page Restore",
	"page_removed":    "Page Remove",
}

func (s *Subscription) IsValid() error {
	// TODO : Clean subscription data
	if s.Alias == "" {
		return errors.New("alias can not be empty")
	}
	if s.BaseURL == "" {
		return errors.New("base url can not be empty")
	}
	if _, err := url2.Parse(s.BaseURL); err != nil {
		return errors.New("enter a valid url")
	}
	if s.SpaceKey == "" {
		return errors.New("space key can not be empty")
	}
	return nil
}

func FormattedSubscriptionList(channelSubscriptions map[string]Subscription) string {
	list := fmt.Sprintf("| Alias | Base Url | Space Key | Events|\n| :----: |:--------:| :--------:| :-----:|")
	for _, subscription := range channelSubscriptions {
		var events []string
		for _, event := range subscription.Events {
			events = append(events, eventTypes[event])
		}
		list += fmt.Sprintf("\n|%s|%s|%s|%s|", subscription.Alias, subscription.BaseURL, subscription.SpaceKey, strings.Join(events, ", "))
	}
	return list
}