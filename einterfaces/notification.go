// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package einterfaces

import (
	"github.com/telluria-bruno-gouvea/mattermost-server/v5/model"
)

type NotificationInterface interface {
	GetNotificationMessage(ack *model.PushNotificationAck, userId string) (*model.PushNotification, *model.AppError)
	CheckLicense() *model.AppError
}
