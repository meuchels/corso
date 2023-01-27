package models

import (
	i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
	msmodel "github.com/microsoftgraph/msgraph-sdk-go/models"
)

// CanvasLayoutable
type CanvasLayoutable interface {
	msmodel.Entityable
	i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
	GetHorizontalSections() []HorizontalSectionable
	GetVerticalSection() VerticalSectionable
	SetHorizontalSections(value []HorizontalSectionable)
	SetVerticalSection(value VerticalSectionable)
}
