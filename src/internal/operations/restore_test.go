package operations

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/alcionai/corso/src/internal/connector/exchange"
	"github.com/alcionai/corso/src/internal/connector/mockconnector"
	"github.com/alcionai/corso/src/internal/connector/support"
	"github.com/alcionai/corso/src/internal/data"
	"github.com/alcionai/corso/src/internal/events"
	evmock "github.com/alcionai/corso/src/internal/events/mock"
	"github.com/alcionai/corso/src/internal/kopia"
	"github.com/alcionai/corso/src/internal/model"
	"github.com/alcionai/corso/src/internal/stats"
	"github.com/alcionai/corso/src/internal/tester"
	"github.com/alcionai/corso/src/internal/tester/aw"
	"github.com/alcionai/corso/src/pkg/account"
	"github.com/alcionai/corso/src/pkg/control"
	"github.com/alcionai/corso/src/pkg/selectors"
	"github.com/alcionai/corso/src/pkg/store"
)

// ---------------------------------------------------------------------------
// unit
// ---------------------------------------------------------------------------

type RestoreOpSuite struct {
	suite.Suite
}

func TestRestoreOpSuite(t *testing.T) {
	suite.Run(t, new(RestoreOpSuite))
}

func (suite *RestoreOpSuite) TestRestoreOperation_PersistResults() {
	ctx, flush := tester.NewContext()
	defer flush()

	var (
		kw   = &kopia.Wrapper{}
		sw   = &store.Wrapper{}
		acct = account.Account{}
		now  = time.Now()
		dest = tester.DefaultTestRestoreDestination()
	)

	table := []struct {
		expectStatus opStatus
		expectErr    assert.ErrorAssertionFunc
		stats        restoreStats
		fail         error
	}{
		{
			expectStatus: Completed,
			expectErr:    aw.NoErr,
			stats: restoreStats{
				resourceCount: 1,
				bytesRead: &stats.ByteCounter{
					NumBytes: 42,
				},
				cs: []data.RestoreCollection{
					data.NotFoundRestoreCollection{
						Collection: &mockconnector.MockExchangeDataCollection{},
					},
				},
				gc: &support.ConnectorOperationStatus{
					Metrics: support.CollectionMetrics{
						Objects:   1,
						Successes: 1,
					},
				},
			},
		},
		{
			expectStatus: Failed,
			expectErr:    aw.Err,
			fail:         assert.AnError,
			stats: restoreStats{
				bytesRead: &stats.ByteCounter{},
				gc:        &support.ConnectorOperationStatus{},
			},
		},
		{
			expectStatus: NoData,
			expectErr:    aw.NoErr,
			stats: restoreStats{
				bytesRead: &stats.ByteCounter{},
				cs:        []data.RestoreCollection{},
				gc:        &support.ConnectorOperationStatus{},
			},
		},
	}
	for _, test := range table {
		suite.T().Run(test.expectStatus.String(), func(t *testing.T) {
			op, err := NewRestoreOperation(
				ctx,
				control.Options{},
				kw,
				sw,
				acct,
				"foo",
				selectors.Selector{DiscreteOwner: "test"},
				dest,
				evmock.NewBus())
			aw.MustNoErr(t, err)

			op.Errors.Fail(test.fail)

			test.expectErr(t, op.persistResults(ctx, now, &test.stats))

			assert.Equal(t, test.expectStatus.String(), op.Status.String(), "status")
			assert.Equal(t, len(test.stats.cs), op.Results.ItemsRead, "items read")
			assert.Equal(t, test.stats.gc.Metrics.Successes, op.Results.ItemsWritten, "items written")
			assert.Equal(t, test.stats.bytesRead.NumBytes, op.Results.BytesRead, "resource owners")
			assert.Equal(t, test.stats.resourceCount, op.Results.ResourceOwners, "resource owners")
			assert.Equal(t, now, op.Results.StartedAt, "started at")
			assert.Less(t, now, op.Results.CompletedAt, "completed at")
		})
	}
}

// ---------------------------------------------------------------------------
// integration
// ---------------------------------------------------------------------------

type RestoreOpIntegrationSuite struct {
	suite.Suite

	backupID    model.StableID
	numItems    int
	kopiaCloser func(ctx context.Context)
	kw          *kopia.Wrapper
	sw          *store.Wrapper
	ms          *kopia.ModelStore
}

func TestRestoreOpIntegrationSuite(t *testing.T) {
	tester.RunOnAny(
		t,
		tester.CorsoCITests,
		tester.CorsoOperationTests)

	suite.Run(t, new(RestoreOpIntegrationSuite))
}

func (suite *RestoreOpIntegrationSuite) SetupSuite() {
	ctx, flush := tester.NewContext()
	defer flush()

	tester.MustGetEnvSets(suite.T(), tester.M365AcctCredEnvs)

	t := suite.T()

	m365UserID := tester.M365UserID(t)
	acct := tester.NewM365Account(t)

	// need to initialize the repository before we can test connecting to it.
	st := tester.NewPrefixedS3Storage(t)

	k := kopia.NewConn(st)
	aw.MustNoErr(t, k.Initialize(ctx))

	suite.kopiaCloser = func(ctx context.Context) {
		k.Close(ctx)
	}

	kw, err := kopia.NewWrapper(k)
	aw.MustNoErr(t, err)

	suite.kw = kw

	ms, err := kopia.NewModelStore(k)
	aw.MustNoErr(t, err)

	suite.ms = ms

	sw := store.NewKopiaStore(ms)
	suite.sw = sw

	users := []string{m365UserID}

	bsel := selectors.NewExchangeBackup(users)
	bsel.DiscreteOwner = m365UserID
	bsel.Include(
		bsel.MailFolders([]string{exchange.DefaultMailFolder}, selectors.PrefixMatch()),
		bsel.ContactFolders([]string{exchange.DefaultContactFolder}, selectors.PrefixMatch()),
		bsel.EventCalendars([]string{exchange.DefaultCalendar}, selectors.PrefixMatch()),
	)

	bo, err := NewBackupOperation(
		ctx,
		control.Options{},
		kw,
		sw,
		acct,
		bsel.Selector,
		evmock.NewBus())
	aw.MustNoErr(t, err)
	aw.MustNoErr(t, bo.Run(ctx))
	require.NotEmpty(t, bo.Results.BackupID)

	suite.backupID = bo.Results.BackupID
	// Discount metadata files (3 paths, 3 deltas) as
	// they are not part of the data restored.
	suite.numItems = bo.Results.ItemsWritten - 6
}

func (suite *RestoreOpIntegrationSuite) TearDownSuite() {
	ctx, flush := tester.NewContext()
	defer flush()

	if suite.ms != nil {
		suite.ms.Close(ctx)
	}

	if suite.kw != nil {
		suite.kw.Close(ctx)
	}

	if suite.kopiaCloser != nil {
		suite.kopiaCloser(ctx)
	}
}

func (suite *RestoreOpIntegrationSuite) TestNewRestoreOperation() {
	kw := &kopia.Wrapper{}
	sw := &store.Wrapper{}
	acct := tester.NewM365Account(suite.T())
	dest := tester.DefaultTestRestoreDestination()

	table := []struct {
		name     string
		opts     control.Options
		kw       *kopia.Wrapper
		sw       *store.Wrapper
		acct     account.Account
		targets  []string
		errCheck assert.ErrorAssertionFunc
	}{
		{"good", control.Options{}, kw, sw, acct, nil, aw.NoErr},
		{"missing kopia", control.Options{}, nil, sw, acct, nil, aw.Err},
		{"missing modelstore", control.Options{}, kw, nil, acct, nil, aw.Err},
	}
	for _, test := range table {
		suite.T().Run(test.name, func(t *testing.T) {
			ctx, flush := tester.NewContext()
			defer flush()

			_, err := NewRestoreOperation(
				ctx,
				test.opts,
				test.kw,
				test.sw,
				test.acct,
				"backup-id",
				selectors.Selector{DiscreteOwner: "test"},
				dest,
				evmock.NewBus())
			test.errCheck(t, err)
		})
	}
}

func (suite *RestoreOpIntegrationSuite) TestRestore_Run() {
	ctx, flush := tester.NewContext()
	defer flush()

	t := suite.T()
	users := []string{tester.M365UserID(t)}

	rsel := selectors.NewExchangeRestore(users)
	rsel.Include(rsel.AllData())

	dest := tester.DefaultTestRestoreDestination()
	mb := evmock.NewBus()

	ro, err := NewRestoreOperation(
		ctx,
		control.Options{},
		suite.kw,
		suite.sw,
		tester.NewM365Account(t),
		suite.backupID,
		rsel.Selector,
		dest,
		mb)
	aw.MustNoErr(t, err)

	ds, err := ro.Run(ctx)

	aw.MustNoErr(t, err, "restoreOp.Run()")
	require.NotEmpty(t, ro.Results, "restoreOp results")
	require.NotNil(t, ds, "restored details")
	assert.Equal(t, ro.Status, Completed, "restoreOp status")
	assert.Equal(t, ro.Results.ItemsWritten, len(ds.Entries), "count of items written matches restored entries in details")
	assert.Less(t, 0, ro.Results.ItemsRead, "restore items read")
	assert.Less(t, 0, ro.Results.ItemsWritten, "restored items written")
	assert.Less(t, int64(0), ro.Results.BytesRead, "bytes read")
	assert.Equal(t, 1, ro.Results.ResourceOwners, "resource Owners")
	aw.NoErr(t, ro.Errors.Failure(), "non-recoverable error")
	assert.Empty(t, ro.Errors.Recovered(), "recoverable errors")
	assert.Equal(t, suite.numItems, ro.Results.ItemsWritten, "backup and restore wrote the same num of items")
	assert.Equal(t, 1, mb.TimesCalled[events.RestoreStart], "restore-start events")
	assert.Equal(t, 1, mb.TimesCalled[events.RestoreEnd], "restore-end events")
}

func (suite *RestoreOpIntegrationSuite) TestRestore_Run_ErrorNoResults() {
	ctx, flush := tester.NewContext()
	defer flush()

	t := suite.T()

	rsel := selectors.NewExchangeRestore(selectors.None())
	rsel.Include(rsel.AllData())

	dest := tester.DefaultTestRestoreDestination()
	mb := evmock.NewBus()

	ro, err := NewRestoreOperation(
		ctx,
		control.Options{},
		suite.kw,
		suite.sw,
		tester.NewM365Account(t),
		suite.backupID,
		rsel.Selector,
		dest,
		mb)
	aw.MustNoErr(t, err)

	ds, err := ro.Run(ctx)
	aw.MustErr(t, err, "restoreOp.Run() should have errored")
	require.Nil(t, ds, "restoreOp.Run() should not produce details")
	assert.Zero(t, ro.Results.ResourceOwners, "resource owners")
	assert.Zero(t, ro.Results.BytesRead, "bytes read")
	assert.Equal(t, 1, mb.TimesCalled[events.RestoreStart], "restore-start events")
	assert.Zero(t, mb.TimesCalled[events.RestoreEnd], "restore-end events")
}
