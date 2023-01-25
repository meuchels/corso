package models

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// WindowsPhone81AppXCollectionResponse 
type WindowsPhone81AppXCollectionResponse struct {
    BaseCollectionPaginationCountResponse
    // The value property
    value []WindowsPhone81AppXable
}
// NewWindowsPhone81AppXCollectionResponse instantiates a new WindowsPhone81AppXCollectionResponse and sets the default values.
func NewWindowsPhone81AppXCollectionResponse()(*WindowsPhone81AppXCollectionResponse) {
    m := &WindowsPhone81AppXCollectionResponse{
        BaseCollectionPaginationCountResponse: *NewBaseCollectionPaginationCountResponse(),
    }
    return m
}
// CreateWindowsPhone81AppXCollectionResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateWindowsPhone81AppXCollectionResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewWindowsPhone81AppXCollectionResponse(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *WindowsPhone81AppXCollectionResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.BaseCollectionPaginationCountResponse.GetFieldDeserializers()
    res["value"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetCollectionOfObjectValues(CreateWindowsPhone81AppXFromDiscriminatorValue)
        if err != nil {
            return err
        }
        if val != nil {
            res := make([]WindowsPhone81AppXable, len(val))
            for i, v := range val {
                res[i] = v.(WindowsPhone81AppXable)
            }
            m.SetValue(res)
        }
        return nil
    }
    return res
}
// GetValue gets the value property value. The value property
func (m *WindowsPhone81AppXCollectionResponse) GetValue()([]WindowsPhone81AppXable) {
    return m.value
}
// Serialize serializes information the current object
func (m *WindowsPhone81AppXCollectionResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.BaseCollectionPaginationCountResponse.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetValue() != nil {
        cast := make([]i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, len(m.GetValue()))
        for i, v := range m.GetValue() {
            cast[i] = v.(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable)
        }
        err = writer.WriteCollectionOfObjectValues("value", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetValue sets the value property value. The value property
func (m *WindowsPhone81AppXCollectionResponse) SetValue(value []WindowsPhone81AppXable)() {
    m.value = value
}
