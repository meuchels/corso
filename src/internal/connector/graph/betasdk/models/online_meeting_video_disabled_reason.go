package models
import (
    "errors"
)
// Provides operations to manage the sites property of the microsoft.graph.browserSiteList entity.
type OnlineMeetingVideoDisabledReason int

const (
    WATERMARKPROTECTION_ONLINEMEETINGVIDEODISABLEDREASON OnlineMeetingVideoDisabledReason = iota
    UNKNOWNFUTUREVALUE_ONLINEMEETINGVIDEODISABLEDREASON
)

func (i OnlineMeetingVideoDisabledReason) String() string {
    return []string{"watermarkProtection", "unknownFutureValue"}[i]
}
func ParseOnlineMeetingVideoDisabledReason(v string) (interface{}, error) {
    result := WATERMARKPROTECTION_ONLINEMEETINGVIDEODISABLEDREASON
    switch v {
        case "watermarkProtection":
            result = WATERMARKPROTECTION_ONLINEMEETINGVIDEODISABLEDREASON
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_ONLINEMEETINGVIDEODISABLEDREASON
        default:
            return 0, errors.New("Unknown OnlineMeetingVideoDisabledReason value: " + v)
    }
    return &result, nil
}
func SerializeOnlineMeetingVideoDisabledReason(values []OnlineMeetingVideoDisabledReason) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
