package managedtenants

import (
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
    ifda19816f54f079134d70c11e75d6b26799300cf72079e282f1d3bb9a6750354 "github.com/alcionai/corso/src/internal/connector/graph/betasdk/models"
)

// ManagedTenantable 
type ManagedTenantable interface {
    ifda19816f54f079134d70c11e75d6b26799300cf72079e282f1d3bb9a6750354.Entityable
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable
    GetAggregatedPolicyCompliances()([]AggregatedPolicyComplianceable)
    GetAuditEvents()([]AuditEventable)
    GetCloudPcConnections()([]CloudPcConnectionable)
    GetCloudPcDevices()([]CloudPcDeviceable)
    GetCloudPcsOverview()([]CloudPcOverviewable)
    GetConditionalAccessPolicyCoverages()([]ConditionalAccessPolicyCoverageable)
    GetCredentialUserRegistrationsSummaries()([]CredentialUserRegistrationsSummaryable)
    GetDeviceCompliancePolicySettingStateSummaries()([]DeviceCompliancePolicySettingStateSummaryable)
    GetManagedDeviceCompliances()([]ManagedDeviceComplianceable)
    GetManagedDeviceComplianceTrends()([]ManagedDeviceComplianceTrendable)
    GetManagedTenantAlertLogs()([]ManagedTenantAlertLogable)
    GetManagedTenantAlertRuleDefinitions()([]ManagedTenantAlertRuleDefinitionable)
    GetManagedTenantAlertRules()([]ManagedTenantAlertRuleable)
    GetManagedTenantAlerts()([]ManagedTenantAlertable)
    GetManagedTenantApiNotifications()([]ManagedTenantApiNotificationable)
    GetManagedTenantEmailNotifications()([]ManagedTenantEmailNotificationable)
    GetManagedTenantTicketingEndpoints()([]ManagedTenantTicketingEndpointable)
    GetManagementActions()([]ManagementActionable)
    GetManagementActionTenantDeploymentStatuses()([]ManagementActionTenantDeploymentStatusable)
    GetManagementIntents()([]ManagementIntentable)
    GetManagementTemplateCollections()([]ManagementTemplateCollectionable)
    GetManagementTemplateCollectionTenantSummaries()([]ManagementTemplateCollectionTenantSummaryable)
    GetManagementTemplates()([]ManagementTemplateable)
    GetManagementTemplateSteps()([]ManagementTemplateStepable)
    GetManagementTemplateStepTenantSummaries()([]ManagementTemplateStepTenantSummaryable)
    GetManagementTemplateStepVersions()([]ManagementTemplateStepVersionable)
    GetMyRoles()([]MyRoleable)
    GetTenantGroups()([]TenantGroupable)
    GetTenants()([]Tenantable)
    GetTenantsCustomizedInformation()([]TenantCustomizedInformationable)
    GetTenantsDetailedInformation()([]TenantDetailedInformationable)
    GetTenantTags()([]TenantTagable)
    GetWindowsDeviceMalwareStates()([]WindowsDeviceMalwareStateable)
    GetWindowsProtectionStates()([]WindowsProtectionStateable)
    SetAggregatedPolicyCompliances(value []AggregatedPolicyComplianceable)()
    SetAuditEvents(value []AuditEventable)()
    SetCloudPcConnections(value []CloudPcConnectionable)()
    SetCloudPcDevices(value []CloudPcDeviceable)()
    SetCloudPcsOverview(value []CloudPcOverviewable)()
    SetConditionalAccessPolicyCoverages(value []ConditionalAccessPolicyCoverageable)()
    SetCredentialUserRegistrationsSummaries(value []CredentialUserRegistrationsSummaryable)()
    SetDeviceCompliancePolicySettingStateSummaries(value []DeviceCompliancePolicySettingStateSummaryable)()
    SetManagedDeviceCompliances(value []ManagedDeviceComplianceable)()
    SetManagedDeviceComplianceTrends(value []ManagedDeviceComplianceTrendable)()
    SetManagedTenantAlertLogs(value []ManagedTenantAlertLogable)()
    SetManagedTenantAlertRuleDefinitions(value []ManagedTenantAlertRuleDefinitionable)()
    SetManagedTenantAlertRules(value []ManagedTenantAlertRuleable)()
    SetManagedTenantAlerts(value []ManagedTenantAlertable)()
    SetManagedTenantApiNotifications(value []ManagedTenantApiNotificationable)()
    SetManagedTenantEmailNotifications(value []ManagedTenantEmailNotificationable)()
    SetManagedTenantTicketingEndpoints(value []ManagedTenantTicketingEndpointable)()
    SetManagementActions(value []ManagementActionable)()
    SetManagementActionTenantDeploymentStatuses(value []ManagementActionTenantDeploymentStatusable)()
    SetManagementIntents(value []ManagementIntentable)()
    SetManagementTemplateCollections(value []ManagementTemplateCollectionable)()
    SetManagementTemplateCollectionTenantSummaries(value []ManagementTemplateCollectionTenantSummaryable)()
    SetManagementTemplates(value []ManagementTemplateable)()
    SetManagementTemplateSteps(value []ManagementTemplateStepable)()
    SetManagementTemplateStepTenantSummaries(value []ManagementTemplateStepTenantSummaryable)()
    SetManagementTemplateStepVersions(value []ManagementTemplateStepVersionable)()
    SetMyRoles(value []MyRoleable)()
    SetTenantGroups(value []TenantGroupable)()
    SetTenants(value []Tenantable)()
    SetTenantsCustomizedInformation(value []TenantCustomizedInformationable)()
    SetTenantsDetailedInformation(value []TenantDetailedInformationable)()
    SetTenantTags(value []TenantTagable)()
    SetWindowsDeviceMalwareStates(value []WindowsDeviceMalwareStateable)()
    SetWindowsProtectionStates(value []WindowsProtectionStateable)()
}
