package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// RelatedContactable 
type RelatedContactable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.AdditionalDataHolder
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAccessConsent()(*bool)
    GetDisplayName()(*string)
    GetEmailAddress()(*string)
    GetId()(*string)
    GetMobilePhone()(*string)
    GetOdataType()(*string)
    GetRelationship()(*ContactRelationship)
    SetAccessConsent(value *bool)()
    SetDisplayName(value *string)()
    SetEmailAddress(value *string)()
    SetId(value *string)()
    SetMobilePhone(value *string)()
    SetOdataType(value *string)()
    SetRelationship(value *ContactRelationship)()
}
