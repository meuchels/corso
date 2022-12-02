package graph

import (
	"context"
	nethttp "net/http"
	"net/http/httputil"
	"strings"
	"time"

	az "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	ka "github.com/microsoft/kiota-authentication-azure-go"
	khttp "github.com/microsoft/kiota-http-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphgocore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/pkg/errors"

	"github.com/alcionai/corso/src/pkg/logger"
	"github.com/alcionai/corso/src/pkg/path"
)

const (
	logGraphRequestsEnvKey = "LOG_GRAPH_REQUESTS"
)

// CreateAdapter uses provided credentials to log into M365 using Kiota Azure Library
// with Azure identity package. An adapter object is a necessary to component
// to create  *msgraphsdk.GraphServiceClient
func CreateAdapter(tenant, client, secret string) (*msgraphsdk.GraphRequestAdapter, error) {
	// Client Provider: Uses Secret for access to tenant-level data
	cred, err := az.NewClientSecretCredential(tenant, client, secret, nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating m365 client secret credentials")
	}

	auth, err := ka.NewAzureIdentityAuthenticationProviderWithScopes(
		cred,
		[]string{"https://graph.microsoft.com/.default"},
	)
	if err != nil {
		return nil, errors.Wrap(err, "creating new AzureIdentityAuthentication")
	}

	clientOptions := msgraphsdk.GetDefaultClientOptions()
	middlewares := msgraphgocore.GetDefaultMiddlewaresWithOptions(&clientOptions)

	// When true, additional logging middleware support added for http request
	// if os.Getenv(logGraphRequestsEnvKey) != "" {
	middlewares = append(middlewares, &LoggingMiddleware{})
	// }

	httpClient := msgraphgocore.GetDefaultClient(&clientOptions, middlewares...)
	httpClient.Timeout = time.Second * 90

	return msgraphsdk.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		auth, nil, nil, httpClient)
}

// LoggingMiddleware can be used to log the http request sent by the graph client
type LoggingMiddleware struct{}

// Intercept implements the RequestInterceptor interface and decodes the parameters name
func (handler *LoggingMiddleware) Intercept(
	pipeline khttp.Pipeline, middlewareIndex int, req *nethttp.Request,
) (*nethttp.Response, error) {
	requestDump, _ := httputil.DumpRequest(req, true)

	resp, err := pipeline.Next(req, middlewareIndex)

	if resp != nil && (resp.StatusCode/100) != 2 {
		dump, _ := httputil.DumpResponse(resp, true)
		logger.Ctx(context.TODO()).Infof("\n-----\nNEW RESP ERR\n-----\n")
		logger.Ctx(context.TODO()).
			Infof("-----\n> %v %v %v\n> %v %v\n> %v %v\n\n> %v %v\n-----\n",
				"url", req.Method, req.URL,
				"reqlen", req.ContentLength,
				"code", resp.Status,
				"RESP:", string(dump))
		if resp.StatusCode == 400 {
			logger.Ctx(context.TODO()).Info("REQUEST:", string(requestDump))
		}
		logger.Ctx(context.TODO()).Infof("-----\nFIN\n-----\n")
	}

	return resp, err
}

func StringToPathCategory(input string) path.CategoryType {
	param := strings.ToLower(input)

	switch param {
	case "email":
		return path.EmailCategory
	case "contacts":
		return path.ContactsCategory
	case "events":
		return path.EventsCategory
	case "files":
		return path.FilesCategory
	case "libraries":
		return path.LibrariesCategory
	default:
		return path.UnknownCategory
	}
}
