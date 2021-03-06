// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package interfaces

import (
	"github.com/telluria-bruno-gouvea/mattermost-server/v5/model"
)

type IndexerJobInterface interface {
	MakeWorker() model.Worker
}
