package groups

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// ItemSitesItemListsItemContentTypesItemIsPublishedResponse provides operations to call the isPublished method.
type ItemSitesItemListsItemContentTypesItemIsPublishedResponse struct {
    // Stores model information.
    backingStore BackingStore
}
// NewItemSitesItemListsItemContentTypesItemIsPublishedResponse instantiates a new ItemSitesItemListsItemContentTypesItemIsPublishedResponse and sets the default values.
func NewItemSitesItemListsItemContentTypesItemIsPublishedResponse()(*ItemSitesItemListsItemContentTypesItemIsPublishedResponse) {
    m := &ItemSitesItemListsItemContentTypesItemIsPublishedResponse{
    }
    m._backingStore = BackingStoreFactorySingleton.Instance.CreateBackingStore();
    m.SetAdditionalData(make(map[string]interface{}));
    return m
}
// CreateItemSitesItemListsItemContentTypesItemIsPublishedResponseFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateItemSitesItemListsItemContentTypesItemIsPublishedResponseFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewItemSitesItemListsItemContentTypesItemIsPublishedResponse(), nil
}
// GetAdditionalData gets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ItemSitesItemListsItemContentTypesItemIsPublishedResponse) GetAdditionalData()(map[string]interface{}) {
    map[string]interface{} value = m._backingStore.Get("additionalData")
    if value == nil {
        value = make(map[string]interface{});
        m.SetAdditionalData(value);
    }
    return value;
}
// GetBackingStore gets the backingStore property value. Stores model information.
func (m *ItemSitesItemListsItemContentTypesItemIsPublishedResponse) GetBackingStore()(BackingStore) {
    return m.backingStore
}
// GetFieldDeserializers the deserialization information for the current model
func (m *ItemSitesItemListsItemContentTypesItemIsPublishedResponse) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := make(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error))
    res["value"] = func (n i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode) error {
        val, err := n.GetBoolValue()
        if err != nil {
            return err
        }
        if val != nil {
            m.SetValue(val)
        }
        return nil
    }
    return res
}
// GetValue gets the value property value. The value property
func (m *ItemSitesItemListsItemContentTypesItemIsPublishedResponse) GetValue()(*bool) {
    return m.GetBackingStore().Get("value");
}
// Serialize serializes information the current object
func (m *ItemSitesItemListsItemContentTypesItemIsPublishedResponse) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    {
        err := writer.WriteBoolValue("value", m.GetValue())
        if err != nil {
            return err
        }
    }
    return nil
}
// SetAdditionalData sets the additionalData property value. Stores additional data not described in the OpenAPI description found when deserializing. Can be used for serialization as well.
func (m *ItemSitesItemListsItemContentTypesItemIsPublishedResponse) SetAdditionalData(value map[string]interface{})() {
    m.GetBackingStore().Set("additionalData", value)
}
// SetBackingStore sets the backingStore property value. Stores model information.
func (m *ItemSitesItemListsItemContentTypesItemIsPublishedResponse) SetBackingStore(value BackingStore)() {
    m.GetBackingStore().Set("backingStore", value)
}
// SetValue sets the value property value. The value property
func (m *ItemSitesItemListsItemContentTypesItemIsPublishedResponse) SetValue(value *bool)() {
    m.GetBackingStore().Set("value", value)
}
