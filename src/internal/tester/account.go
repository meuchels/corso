package tester

import (
	"testing"

	"github.com/alcionai/corso/src/internal/tester/aw"
	"github.com/alcionai/corso/src/pkg/account"
	"github.com/alcionai/corso/src/pkg/credentials"
)

var M365AcctCredEnvs = []string{
	credentials.AzureClientID,
	credentials.AzureClientSecret,
}

// NewM365Account returns an account.Account object initialized with environment
// variables used for integration tests that use Graph Connector.
func NewM365Account(t *testing.T) account.Account {
	cfg, err := readTestConfig()
	aw.MustNoErr(t, err, "configuring m365 account from test configuration")

	acc, err := account.NewAccount(
		account.ProviderM365,
		account.M365Config{
			M365:          credentials.GetM365(),
			AzureTenantID: cfg[TestCfgAzureTenantID],
		},
	)
	aw.MustNoErr(t, err, "initializing account")

	return acc
}

func NewMockM365Account(t *testing.T) account.Account {
	acc, err := account.NewAccount(
		account.ProviderM365,
		account.M365Config{
			M365: credentials.M365{
				AzureClientID:     "12345",
				AzureClientSecret: "abcde",
			},
			AzureTenantID: "09876",
		},
	)
	aw.MustNoErr(t, err, "initializing mock account")

	return acc
}
