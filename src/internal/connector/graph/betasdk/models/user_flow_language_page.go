package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// UserFlowLanguagePage provides operations to manage the sites property of the microsoft.graph.browserSiteList entity.
type UserFlowLanguagePage struct {
    Entity
}
// NewUserFlowLanguagePage instantiates a new userFlowLanguagePage and sets the default values.
func NewUserFlowLanguagePage()(*UserFlowLanguagePage) {
    m := &UserFlowLanguagePage{
        Entity: *NewEntity(),
    }
    return m
}
// CreateUserFlowLanguagePageFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateUserFlowLanguagePageFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewUserFlowLanguagePage(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *UserFlowLanguagePage) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    return res
}
// Serialize serializes information the current object
func (m *UserFlowLanguagePage) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    return nil
}
