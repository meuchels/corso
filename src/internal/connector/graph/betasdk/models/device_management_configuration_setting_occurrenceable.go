package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// DeviceManagementConfigurationSettingOccurrenceable 
type DeviceManagementConfigurationSettingOccurrenceable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetMaxDeviceOccurrence()(*int32)
    GetMinDeviceOccurrence()(*int32)
    GetOdataType()(*string)
    SetMaxDeviceOccurrence(value *int32)()
    SetMinDeviceOccurrence(value *int32)()
    SetOdataType(value *string)()
}
