package sites

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ifda19816f54f079134d70c11e75d6b26799300cf72079e282f1d3bb9a6750354 "github.com/alcionai/corso/src/internal/connector/graph/betasdk/models"
)

// ItemOnenotePagesItemOnenotePatchContentPostRequestBodyable 
type ItemOnenotePagesItemOnenotePatchContentPostRequestBodyable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCommands()([]ifda19816f54f079134d70c11e75d6b26799300cf72079e282f1d3bb9a6750354.OnenotePatchContentCommandable)
    SetCommands(value []ifda19816f54f079134d70c11e75d6b26799300cf72079e282f1d3bb9a6750354.OnenotePatchContentCommandable)()
}
