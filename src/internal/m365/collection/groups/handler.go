package groups

import (
	"context"

	"github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/alcionai/corso/src/pkg/services/m365/api"
)

type BackupMessagesHandler interface {
	GetMessage(ctx context.Context, teamID, channelID, itemID string) (models.ChatMessageable, error)
	NewMessagePager(teamID, channelID string) api.ChannelMessageDeltaEnumerator
	GetChannel(ctx context.Context, teamID, channelID string) (models.Channelable, error)
	GetReply(ctx context.Context, teamID, channelID, messageID string) (serialization.Parsable, error)
}
