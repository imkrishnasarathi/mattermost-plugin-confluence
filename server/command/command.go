package command

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"

	"github.com/Brightscout/mattermost-plugin-confluence/server/config"
	"github.com/Brightscout/mattermost-plugin-confluence/server/serializer"
	"github.com/Brightscout/mattermost-plugin-confluence/server/service"
)

type HandlerFunc func(context *model.CommandArgs, args ...string) *model.CommandResponse

type Handler struct {
	handlers       map[string]HandlerFunc
	defaultHandler HandlerFunc
}

const (
	specifyAlias              = "Please specify alias."
	subscriptionDeleteSuccess = "Subscription with alias **%s** deleted successfully."
	noChannelSubscription     = "No subscription found for this channel."
)

var ConfluenceCommandHandler = Handler{
	handlers: map[string]HandlerFunc{
		"list":        listChannelSubscription,
		"unsubscribe": deleteSubscription,
		"edit":        editSubscription,
	},
	defaultHandler: executeConfluenceDefault,
}

func GetCommand() *model.Command {
	return &model.Command{
		Trigger:          "confluence",
		DisplayName:      "Confluence",
		Description:      "Integration with Confluence.",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: subscribe, list, unsubscribe \"<alias>\"",
		AutoCompleteHint: "[command]",
	}
}

// TODO : Show help text instead of invalid command.
func executeConfluenceDefault(context *model.CommandArgs, args ...string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         "Invalid command",
	}
}

func postCommandResponse(context *model.CommandArgs, text string) {
	post := &model.Post{
		UserId:    config.BotUserID,
		ChannelId: context.ChannelId,
		Message:   text,
	}
	_ = config.Mattermost.SendEphemeralPost(context.UserId, post)
}

func (ch Handler) Handle(context *model.CommandArgs, args ...string) *model.CommandResponse {
	for n := len(args); n > 0; n-- {
		h := ch.handlers[strings.Join(args[:n], "/")]
		if h != nil {
			return h(context, args[n:]...)
		}
	}
	return ch.defaultHandler(context, args...)
}

func deleteSubscription(context *model.CommandArgs, args ...string) *model.CommandResponse {
	if len(args) == 0 {
		postCommandResponse(context, specifyAlias)
		return &model.CommandResponse{}
	}
	alias := args[0]
	if err := service.DeleteSubscription(context.ChannelId, args[0]); err != nil {
		postCommandResponse(context, err.Error())
		return &model.CommandResponse{}
	}
	postCommandResponse(context, fmt.Sprintf(subscriptionDeleteSuccess, alias))
	return &model.CommandResponse{}
}

func listChannelSubscription(context *model.CommandArgs, args ...string) *model.CommandResponse {
	channelSubscriptions, _, gErr := service.GetChannelSubscriptions(context.ChannelId)
	if gErr != nil {
		postCommandResponse(context, gErr.Error())
		return &model.CommandResponse{}
	}

	if len(channelSubscriptions) == 0 {
		postCommandResponse(context, noChannelSubscription)
		return &model.CommandResponse{}
	}
	list := serializer.FormattedSubscriptionList(channelSubscriptions)
	postCommandResponse(context, list)
	return &model.CommandResponse{}
}

func editSubscription(context *model.CommandArgs, args ...string) *model.CommandResponse {
	if len(args) == 0 {
		postCommandResponse(context, specifyAlias)
		return &model.CommandResponse{}
	}
	alias := args[0]
	if err := service.OpenSubscriptionEditModal(context.ChannelId, context.UserId, alias); err != nil {
		postCommandResponse(context, err.Error())
	}
	return &model.CommandResponse{}
}