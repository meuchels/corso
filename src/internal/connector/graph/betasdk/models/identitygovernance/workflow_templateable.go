package identitygovernance

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ifda19816f54f079134d70c11e75d6b26799300cf72079e282f1d3bb9a6750354 "github.com/alcionai/corso/src/internal/connector/graph/betasdk/models"
)

// WorkflowTemplateable 
type WorkflowTemplateable interface {
    ifda19816f54f079134d70c11e75d6b26799300cf72079e282f1d3bb9a6750354.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetCategory()(*LifecycleWorkflowCategory)
    GetDescription()(*string)
    GetDisplayName()(*string)
    GetExecutionConditions()(WorkflowExecutionConditionsable)
    GetTasks()([]Taskable)
    SetCategory(value *LifecycleWorkflowCategory)()
    SetDescription(value *string)()
    SetDisplayName(value *string)()
    SetExecutionConditions(value WorkflowExecutionConditionsable)()
    SetTasks(value []Taskable)()
}
