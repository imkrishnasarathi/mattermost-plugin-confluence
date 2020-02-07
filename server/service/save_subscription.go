package service

import (
	"net/http"

	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/store"
)

const (
	generalSaveError        = "An error occurred attempting to save a subscription."
	aliasAlreadyExist       = "A subscription with the same alias already exists."
	urlSpaceKeyAlreadyExist = "A subscription with the same url and space key already exists."
	subscriptionSaveSuccess = "Your subscription has been saved."
)

// TODO : Check if we can refactor this.
func SaveNewSubscription(subscription serializer.Subscription, userID string) (int, error) {
	channelSubscriptions, cKey, gErr := GetChannelSubscriptions(subscription.ChannelID)
	if gErr != nil {
		return http.StatusInternalServerError, errors.New(generalSaveError)
	}
	if _, ok := channelSubscriptions[subscription.Alias]; ok {
		return http.StatusBadRequest, errors.New(aliasAlreadyExist)
	}
	keySubscriptions, key, keyErr := GetURLSpaceKeyCombinationSubscriptions(subscription.BaseURL, subscription.SpaceKey)
	if keyErr != nil {
		return http.StatusInternalServerError, keyErr
	}
	if _, ok := keySubscriptions[subscription.ChannelID]; ok {
		return http.StatusBadRequest, errors.New(urlSpaceKeyAlreadyExist)
	}

	keySubscriptions[subscription.ChannelID] = subscription.Events
	channelSubscriptions[subscription.Alias] = subscription
	if err := store.Set(key, keySubscriptions); err != nil {
		return http.StatusInternalServerError, errors.New(generalSaveError)
	}
	if err := store.Set(cKey, channelSubscriptions); err != nil {
		return http.StatusInternalServerError, errors.New(generalSaveError)
	}

	post := &model.Post{
		UserId:    config.BotUserID,
		ChannelId: subscription.ChannelID,
		Message:   subscriptionSaveSuccess,
	}
	_ = config.Mattermost.SendEphemeralPost(userID, post)

	return http.StatusOK, nil
}
