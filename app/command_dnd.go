// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	goi18n "github.com/mattermost/go-i18n/i18n"
	"github.com/telluria-bruno-gouvea/mattermost-server/v5/model"
)

type DndProvider struct {
}

const (
	CMD_DND = "dnd"
)

func init() {
	RegisterCommandProvider(&DndProvider{})
}

func (me *DndProvider) GetTrigger() string {
	return CMD_DND
}

func (me *DndProvider) GetCommand(a *App, T goi18n.TranslateFunc) *model.Command {
	return &model.Command{
		Trigger:          CMD_DND,
		AutoComplete:     true,
		AutoCompleteDesc: T("api.command_dnd.desc"),
		DisplayName:      T("api.command_dnd.name"),
	}
}

func (me *DndProvider) DoCommand(a *App, args *model.CommandArgs, message string) *model.CommandResponse {
	status, err := a.GetStatus(args.UserId)
	if err != nil {
		return &model.CommandResponse{ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL, Text: args.T("api.command_dnd.error")}
	} else {
		if status.Status == "dnd" {
			a.SetStatusOnline(args.UserId, true)
			return &model.CommandResponse{ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL, Text: args.T("api.command_dnd.disabled")}
		}
	}

	a.SetStatusDoNotDisturb(args.UserId)

	return &model.CommandResponse{ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL, Text: args.T("api.command_dnd.success")}
}
