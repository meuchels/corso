package sharepoint

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/alcionai/corso/src/internal/connector/discovery/api"
	"github.com/alcionai/corso/src/internal/connector/graph"
	"github.com/alcionai/corso/src/internal/connector/onedrive"
	sapi "github.com/alcionai/corso/src/internal/connector/sharepoint/api"
	"github.com/alcionai/corso/src/internal/connector/support"
	"github.com/alcionai/corso/src/internal/data"
	"github.com/alcionai/corso/src/internal/observe"
	"github.com/alcionai/corso/src/pkg/account"
	"github.com/alcionai/corso/src/pkg/control"
	"github.com/alcionai/corso/src/pkg/logger"
	"github.com/alcionai/corso/src/pkg/path"
	"github.com/alcionai/corso/src/pkg/selectors"
)

type statusUpdater interface {
	UpdateStatus(status *support.ConnectorOperationStatus)
}

// DataCollections returns a set of DataCollection which represents the SharePoint data
// for the specified user
func DataCollections(
	ctx context.Context,
	itemClient *http.Client,
	selector selectors.Selector,
	metadata []data.RestoreCollection,
	creds account.M365Config,
	serv graph.Servicer,
	su statusUpdater,
	ctrlOpts control.Options,
) ([]data.BackupCollection, map[string]struct{}, error) {
	b, err := selector.ToSharePointBackup()
	if err != nil {
		return nil, nil, errors.Wrap(err, "sharePointDataCollection: parsing selector")
	}

	var (
		site        = b.DiscreteOwner
		collections = []data.BackupCollection{}
		errs        error
	)

	for _, scope := range b.Scopes() {
		foldersComplete, closer := observe.MessageWithCompletion(ctx, observe.Bulletf(
			"%s - %s",
			observe.Safe(scope.Category().PathType().String()),
			observe.PII(site)))
		defer closer()
		defer close(foldersComplete)

		var spcs []data.BackupCollection

		switch scope.Category().PathType() {
		case path.ListsCategory:
			spcs, err = collectLists(
				ctx,
				serv,
				creds.AzureTenantID,
				site,
				su,
				ctrlOpts)
			if err != nil {
				return nil, nil, support.WrapAndAppend(site, err, errs)
			}

		case path.LibrariesCategory:
			spcs, _, err = collectLibraries(
				ctx,
				itemClient,
				metadata,
				serv,
				creds.AzureTenantID,
				site,
				scope,
				su,
				ctrlOpts)
			if err != nil {
				return nil, nil, support.WrapAndAppend(site, err, errs)
			}
		case path.PagesCategory:
			spcs, err = collectPages(
				ctx,
				creds,
				serv,
				site,
				su,
				ctrlOpts)
			if err != nil {
				return nil, nil, support.WrapAndAppend(site, err, errs)
			}
		}

		collections = append(collections, spcs...)
		foldersComplete <- struct{}{}
	}

	return collections, nil, errs
}

func collectLists(
	ctx context.Context,
	serv graph.Servicer,
	tenantID, siteID string,
	updater statusUpdater,
	ctrlOpts control.Options,
) ([]data.BackupCollection, error) {
	logger.Ctx(ctx).With("site", siteID).Debug("Creating SharePoint List Collections")

	spcs := make([]data.BackupCollection, 0)

	tuples, err := preFetchLists(ctx, serv, siteID)
	if err != nil {
		return nil, err
	}

	for _, tuple := range tuples {
		dir, err := path.Builder{}.Append(tuple.name).
			ToDataLayerSharePointPath(
				tenantID,
				siteID,
				path.ListsCategory,
				false)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create collection path for site: %s", siteID)
		}

		collection := NewCollection(dir, serv, List, updater.UpdateStatus, ctrlOpts)
		collection.AddJob(tuple.id)

		spcs = append(spcs, collection)
	}

	return spcs, nil
}

// collectLibraries constructs a onedrive Collections struct and Get()s
// all the drives associated with the site.
func collectLibraries(
	ctx context.Context,
	itemClient *http.Client,
	metadata []data.RestoreCollection,
	serv graph.Servicer,
	tenantID, siteID string,
	scope selectors.SharePointScope,
	updater statusUpdater,
	ctrlOpts control.Options,
) ([]data.BackupCollection, map[string]struct{}, error) {
	var (
		collections = []data.BackupCollection{}
		errs        error
	)

	logger.Ctx(ctx).With("site", siteID).Debug("Creating SharePoint Library collections")

	odcs, excludes, err := onedrive.NewCollections(
		itemClient,
		tenantID,
		siteID,
		onedrive.SharePointSource,
		folderMatcher{scope},
		serv,
		updater.UpdateStatus,
		ctrlOpts).Get(ctx, metadata)
	if err != nil {
		return nil, nil, support.WrapAndAppend(siteID, err, errs)
	}

	return append(collections, odcs...), excludes, errs
}

// collectPages constructs a sharepoint Collections struct and Get()s the associated
// M365 IDs for the associated Pages.
func collectPages(
	ctx context.Context,
	creds account.M365Config,
	serv graph.Servicer,
	siteID string,
	updater statusUpdater,
	ctrlOpts control.Options,
) ([]data.BackupCollection, error) {
	logger.Ctx(ctx).With("site", siteID).Debug("Creating SharePoint Pages collections")

	spcs := make([]data.BackupCollection, 0)

	// make the betaClient
	// Need to receive From DataCollection Call
	adpt, err := graph.CreateAdapter(creds.AzureTenantID, creds.AzureClientID, creds.AzureClientSecret)
	if err != nil {
		return nil, errors.New("unable to create adapter w/ env credentials")
	}

	betaService := api.NewBetaService(adpt)

	tuples, err := sapi.FetchPages(ctx, betaService, siteID)
	if err != nil {
		return nil, err
	}

	for _, tuple := range tuples {
		dir, err := path.Builder{}.Append(tuple.Name).
			ToDataLayerSharePointPath(
				creds.AzureTenantID,
				siteID,
				path.PagesCategory,
				false)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create collection path for site: %s", siteID)
		}

		collection := NewCollection(dir, serv, Pages, updater.UpdateStatus, ctrlOpts)
		collection.betaService = betaService
		collection.AddJob(tuple.ID)

		spcs = append(spcs, collection)
	}

	return spcs, nil
}

type folderMatcher struct {
	scope selectors.SharePointScope
}

func (fm folderMatcher) IsAny() bool {
	return fm.scope.IsAny(selectors.SharePointLibrary)
}

func (fm folderMatcher) Matches(dir string) bool {
	return fm.scope.Matches(selectors.SharePointLibrary, dir)
}
