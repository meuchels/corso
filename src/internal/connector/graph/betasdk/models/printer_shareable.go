package models

import (
    i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e "time"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// PrinterShareable 
type PrinterShareable interface {
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    PrinterBaseable
    GetAllowAllUsers()(*bool)
    GetAllowedGroups()([]Groupable)
    GetAllowedUsers()([]Userable)
    GetCreatedDateTime()(*i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)
    GetPrinter()(Printerable)
    GetViewPoint()(PrinterShareViewpointable)
    SetAllowAllUsers(value *bool)()
    SetAllowedGroups(value []Groupable)()
    SetAllowedUsers(value []Userable)()
    SetCreatedDateTime(value *i336074805fc853987abe6f7fe3ad97a6a6f3077a16391fec744f671a015fbd7e.Time)()
    SetPrinter(value Printerable)()
    SetViewPoint(value PrinterShareViewpointable)()
}
