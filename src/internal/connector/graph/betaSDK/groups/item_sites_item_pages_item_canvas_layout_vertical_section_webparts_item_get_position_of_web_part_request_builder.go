package groups

import (
    "context"
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    ifda19816f54f079134d70c11e75d6b26799300cf72079e282f1d3bb9a6750354 "github.com/alcionai/corso/src/internal/connector/graph/betasdk/models"
    i7ad325c11fbf3db4d761c429267362d8b24daa1eda0081f914ebc3cdc85181a0 "github.com/alcionai/corso/src/internal/connector/graph/betasdk/models/odataerrors"
)

// ItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilder provides operations to call the getPositionOfWebPart method.
type ItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilder struct {
    // Path parameters for the request
    pathParameters map[string]string
    // The request adapter to use to execute the requests.
    requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter
    // Url template to use to build the URL for the current request builder
    urlTemplate string
}
// ItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilderPostRequestConfiguration configuration for the request such as headers, query parameters, and middleware options.
type ItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilderPostRequestConfiguration struct {
    // Request headers
    Headers *i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestHeaders
    // Request options
    Options []i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestOption
}
// NewItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilderInternal instantiates a new GetPositionOfWebPartRequestBuilder and sets the default values.
func NewItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilderInternal(pathParameters map[string]string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilder) {
    m := &ItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilder{
    }
    m.urlTemplate = "{+baseurl}/groups/{group%2Did}/sites/{site%2Did}/pages/{sitePage%2Did}/canvasLayout/verticalSection/webparts/{webPart%2Did}/microsoft.graph.getPositionOfWebPart";
    urlTplParams := make(map[string]string)
    for idx, item := range pathParameters {
        urlTplParams[idx] = item
    }
    m.pathParameters = urlTplParams;
    m.requestAdapter = requestAdapter;
    return m
}
// NewItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilder instantiates a new GetPositionOfWebPartRequestBuilder and sets the default values.
func NewItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilder(rawUrl string, requestAdapter i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestAdapter)(*ItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilder) {
    urlParams := make(map[string]string)
    urlParams["request-raw-url"] = rawUrl
    return NewItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilderInternal(urlParams, requestAdapter)
}
// CreatePostRequestInformation invoke action getPositionOfWebPart
func (m *ItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilder) CreatePostRequestInformation(ctx context.Context, requestConfiguration *ItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilderPostRequestConfiguration)(*i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.RequestInformation, error) {
    requestInfo := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.NewRequestInformation()
    requestInfo.UrlTemplate = m.urlTemplate
    requestInfo.PathParameters = m.pathParameters
    requestInfo.Method = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.POST
    requestInfo.Headers.Add("Accept", "application/json")
    if requestConfiguration != nil {
        requestInfo.Headers.AddAll(requestConfiguration.Headers)
        requestInfo.AddRequestOptions(requestConfiguration.Options)
    }
    return requestInfo, nil
}
// Post invoke action getPositionOfWebPart
// [Find more info here]
// 
// [Find more info here]: https://docs.microsoft.com/graph/api/webpart-getposition?view=graph-rest-1.0
func (m *ItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilder) Post(ctx context.Context, requestConfiguration *ItemSitesItemPagesItemCanvasLayoutVerticalSectionWebpartsItemGetPositionOfWebPartRequestBuilderPostRequestConfiguration)(ifda19816f54f079134d70c11e75d6b26799300cf72079e282f1d3bb9a6750354.WebPartPositionable, error) {
    requestInfo, err := m.CreatePostRequestInformation(ctx, requestConfiguration);
    if err != nil {
        return nil, err
    }
    errorMapping := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.ErrorMappings {
        "4XX": i7ad325c11fbf3db4d761c429267362d8b24daa1eda0081f914ebc3cdc85181a0.CreateODataErrorFromDiscriminatorValue,
        "5XX": i7ad325c11fbf3db4d761c429267362d8b24daa1eda0081f914ebc3cdc85181a0.CreateODataErrorFromDiscriminatorValue,
    }
    res, err := m.requestAdapter.SendAsync(ctx, requestInfo, ifda19816f54f079134d70c11e75d6b26799300cf72079e282f1d3bb9a6750354.CreateWebPartPositionFromDiscriminatorValue, errorMapping)
    if err != nil {
        return nil, err
    }
    if res == nil {
        return nil, nil
    }
    return res.(ifda19816f54f079134d70c11e75d6b26799300cf72079e282f1d3bb9a6750354.WebPartPositionable), nil
}
