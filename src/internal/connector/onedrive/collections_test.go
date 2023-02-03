package onedrive

import (
	"context"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/maps"

	"github.com/alcionai/corso/src/internal/connector/graph"
	gapi "github.com/alcionai/corso/src/internal/connector/graph/api"
	"github.com/alcionai/corso/src/internal/connector/support"
	"github.com/alcionai/corso/src/internal/data"
	"github.com/alcionai/corso/src/internal/tester"
	"github.com/alcionai/corso/src/pkg/control"
	"github.com/alcionai/corso/src/pkg/path"
	"github.com/alcionai/corso/src/pkg/selectors"
)

const (
	testBaseDrivePath = "drive/driveID1/root:"
)

func expectedPathAsSlice(t *testing.T, tenant, user string, rest ...string) []string {
	res := make([]string, 0, len(rest))

	for _, r := range rest {
		p, err := GetCanonicalPath(r, tenant, user, OneDriveSource)
		require.NoError(t, err)

		res = append(res, p.String())
	}

	return res
}

type OneDriveCollectionsSuite struct {
	suite.Suite
}

func TestOneDriveCollectionsSuite(t *testing.T) {
	suite.Run(t, new(OneDriveCollectionsSuite))
}

func (suite *OneDriveCollectionsSuite) TestGetCanonicalPath() {
	tenant, resourceOwner := "tenant", "resourceOwner"

	table := []struct {
		name      string
		source    driveSource
		dir       []string
		expect    string
		expectErr assert.ErrorAssertionFunc
	}{
		{
			name:      "onedrive",
			source:    OneDriveSource,
			dir:       []string{"onedrive"},
			expect:    "tenant/onedrive/resourceOwner/files/onedrive",
			expectErr: assert.NoError,
		},
		{
			name:      "sharepoint",
			source:    SharePointSource,
			dir:       []string{"sharepoint"},
			expect:    "tenant/sharepoint/resourceOwner/libraries/sharepoint",
			expectErr: assert.NoError,
		},
		{
			name:      "unknown",
			source:    unknownDriveSource,
			dir:       []string{"unknown"},
			expectErr: assert.Error,
		},
	}
	for _, test := range table {
		suite.T().Run(test.name, func(t *testing.T) {
			p := strings.Join(test.dir, "/")
			result, err := GetCanonicalPath(p, tenant, resourceOwner, test.source)
			test.expectErr(t, err)
			if result != nil {
				assert.Equal(t, test.expect, result.String())
			}
		})
	}
}

func (suite *OneDriveCollectionsSuite) TestUpdateCollections() {
	anyFolder := (&selectors.OneDriveBackup{}).Folders(selectors.Any())[0]

	const (
		tenant    = "tenant"
		user      = "user"
		folder    = "/folder"
		folderSub = "/folder/subfolder"
		pkg       = "/package"
	)

	tests := []struct {
		testCase                string
		items                   []models.DriveItemable
		inputFolderMap          map[string]string
		scope                   selectors.OneDriveScope
		expect                  assert.ErrorAssertionFunc
		expectedCollectionPaths []string
		expectedItemCount       int
		expectedContainerCount  int
		expectedFileCount       int
		expectedMetadataPaths   map[string]string
		expectedExcludes        map[string]struct{}
	}{
		{
			testCase: "Invalid item",
			items: []models.DriveItemable{
				driveItem("item", "item", testBaseDrivePath, false, false, false),
			},
			inputFolderMap:        map[string]string{},
			scope:                 anyFolder,
			expect:                assert.Error,
			expectedMetadataPaths: map[string]string{},
			expectedExcludes:      map[string]struct{}{},
		},
		{
			testCase: "Single File",
			items: []models.DriveItemable{
				driveItem("file", "file", testBaseDrivePath, true, false, false),
			},
			inputFolderMap: map[string]string{},
			scope:          anyFolder,
			expect:         assert.NoError,
			expectedCollectionPaths: expectedPathAsSlice(
				suite.T(),
				tenant,
				user,
				testBaseDrivePath,
			),
			expectedItemCount:      1,
			expectedFileCount:      1,
			expectedContainerCount: 1,
			// Root folder is skipped since it's always present.
			expectedMetadataPaths: map[string]string{},
			expectedExcludes:      map[string]struct{}{},
		},
		{
			testCase: "Single Folder",
			items: []models.DriveItemable{
				driveItem("folder", "folder", testBaseDrivePath, false, true, false),
			},
			inputFolderMap: map[string]string{},
			scope:          anyFolder,
			expect:         assert.NoError,
			expectedCollectionPaths: expectedPathAsSlice(
				suite.T(),
				tenant,
				user,
				testBaseDrivePath,
			),
			expectedMetadataPaths: map[string]string{
				"folder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder",
				)[0],
			},
			expectedItemCount:      1,
			expectedContainerCount: 1,
			expectedExcludes:       map[string]struct{}{},
		},
		{
			testCase: "Single Package",
			items: []models.DriveItemable{
				driveItem("package", "package", testBaseDrivePath, false, false, true),
			},
			inputFolderMap: map[string]string{},
			scope:          anyFolder,
			expect:         assert.NoError,
			expectedCollectionPaths: expectedPathAsSlice(
				suite.T(),
				tenant,
				user,
				testBaseDrivePath,
			),
			expectedMetadataPaths: map[string]string{
				"package": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/package",
				)[0],
			},
			expectedItemCount:      1,
			expectedContainerCount: 1,
			expectedExcludes:       map[string]struct{}{},
		},
		{
			testCase: "1 root file, 1 folder, 1 package, 2 files, 3 collections",
			items: []models.DriveItemable{
				driveItem("fileInRoot", "fileInRoot", testBaseDrivePath, true, false, false),
				driveItem("folder", "folder", testBaseDrivePath, false, true, false),
				driveItem("package", "package", testBaseDrivePath, false, false, true),
				driveItem("fileInFolder", "fileInFolder", testBaseDrivePath+folder, true, false, false),
				driveItem("fileInPackage", "fileInPackage", testBaseDrivePath+pkg, true, false, false),
			},
			inputFolderMap: map[string]string{},
			scope:          anyFolder,
			expect:         assert.NoError,
			expectedCollectionPaths: expectedPathAsSlice(
				suite.T(),
				tenant,
				user,
				testBaseDrivePath,
				testBaseDrivePath+folder,
				testBaseDrivePath+pkg,
			),
			expectedItemCount:      5,
			expectedFileCount:      3,
			expectedContainerCount: 3,
			expectedMetadataPaths: map[string]string{
				"folder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder",
				)[0],
				"package": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/package",
				)[0],
			},
			expectedExcludes: map[string]struct{}{},
		},
		{
			testCase: "contains folder selector",
			items: []models.DriveItemable{
				driveItem("fileInRoot", "fileInRoot", testBaseDrivePath, true, false, false),
				driveItem("folder", "folder", testBaseDrivePath, false, true, false),
				driveItem("subfolder", "subfolder", testBaseDrivePath+folder, false, true, false),
				driveItem("folder2", "folder", testBaseDrivePath+folderSub, false, true, false),
				driveItem("package", "package", testBaseDrivePath, false, false, true),
				driveItem("fileInFolder", "fileInFolder", testBaseDrivePath+folder, true, false, false),
				driveItem("fileInFolder2", "fileInFolder2", testBaseDrivePath+folderSub+folder, true, false, false),
				driveItem("fileInFolderPackage", "fileInPackage", testBaseDrivePath+pkg, true, false, false),
			},
			inputFolderMap: map[string]string{},
			scope:          (&selectors.OneDriveBackup{}).Folders([]string{"folder"})[0],
			expect:         assert.NoError,
			expectedCollectionPaths: expectedPathAsSlice(
				suite.T(),
				tenant,
				user,
				testBaseDrivePath+"/folder",
				testBaseDrivePath+folderSub,
				testBaseDrivePath+folderSub+folder,
			),
			expectedItemCount:      4,
			expectedFileCount:      2,
			expectedContainerCount: 3,
			// just "folder" isn't added here because the include check is done on the
			// parent path since we only check later if something is a folder or not.
			expectedMetadataPaths: map[string]string{
				"subfolder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder/subfolder",
				)[0],
				"folder2": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder/subfolder/folder",
				)[0],
			},
			expectedExcludes: map[string]struct{}{},
		},
		{
			testCase: "prefix subfolder selector",
			items: []models.DriveItemable{
				driveItem("fileInRoot", "fileInRoot", testBaseDrivePath, true, false, false),
				driveItem("folder", "folder", testBaseDrivePath, false, true, false),
				driveItem("subfolder", "subfolder", testBaseDrivePath+folder, false, true, false),
				driveItem("folder2", "folder", testBaseDrivePath+folderSub, false, true, false),
				driveItem("package", "package", testBaseDrivePath, false, false, true),
				driveItem("fileInFolder", "fileInFolder", testBaseDrivePath+folder, true, false, false),
				driveItem("fileInFolder2", "fileInFolder2", testBaseDrivePath+folderSub+folder, true, false, false),
				driveItem("fileInPackage", "fileInPackage", testBaseDrivePath+pkg, true, false, false),
			},
			inputFolderMap: map[string]string{},
			scope: (&selectors.OneDriveBackup{}).
				Folders([]string{"/folder/subfolder"}, selectors.PrefixMatch())[0],
			expect: assert.NoError,
			expectedCollectionPaths: expectedPathAsSlice(
				suite.T(),
				tenant,
				user,
				testBaseDrivePath+folderSub,
				testBaseDrivePath+folderSub+folder,
			),
			expectedItemCount:      2,
			expectedFileCount:      1,
			expectedContainerCount: 2,
			expectedMetadataPaths: map[string]string{
				"folder2": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder/subfolder/folder",
				)[0],
			},
			expectedExcludes: map[string]struct{}{},
		},
		{
			testCase: "match subfolder selector",
			items: []models.DriveItemable{
				driveItem("fileInRoot", "fileInRoot", testBaseDrivePath, true, false, false),
				driveItem("folder", "folder", testBaseDrivePath, false, true, false),
				driveItem("subfolder", "subfolder", testBaseDrivePath+folder, false, true, false),
				driveItem("package", "package", testBaseDrivePath, false, false, true),
				driveItem("fileInFolder", "fileInFolder", testBaseDrivePath+folder, true, false, false),
				driveItem("fileInSubfolder", "fileInSubfolder", testBaseDrivePath+folderSub, true, false, false),
				driveItem("fileInPackage", "fileInPackage", testBaseDrivePath+pkg, true, false, false),
			},
			inputFolderMap: map[string]string{},
			scope:          (&selectors.OneDriveBackup{}).Folders([]string{"folder/subfolder"})[0],
			expect:         assert.NoError,
			expectedCollectionPaths: expectedPathAsSlice(
				suite.T(),
				tenant,
				user,
				testBaseDrivePath+folderSub,
			),
			expectedItemCount:      1,
			expectedFileCount:      1,
			expectedContainerCount: 1,
			// No child folders for subfolder so nothing here.
			expectedMetadataPaths: map[string]string{},
			expectedExcludes:      map[string]struct{}{},
		},
		{
			testCase: "not moved folder tree",
			items: []models.DriveItemable{
				driveItem("folder", "folder", testBaseDrivePath, false, true, false),
			},
			inputFolderMap: map[string]string{
				"folder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder",
				)[0],
				"subfolder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder/subfolder",
				)[0],
			},
			scope:  anyFolder,
			expect: assert.NoError,
			expectedCollectionPaths: expectedPathAsSlice(
				suite.T(),
				tenant,
				user,
				testBaseDrivePath,
			),
			expectedItemCount:      1,
			expectedFileCount:      0,
			expectedContainerCount: 1,
			expectedMetadataPaths: map[string]string{
				"folder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder",
				)[0],
				"subfolder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder/subfolder",
				)[0],
			},
			expectedExcludes: map[string]struct{}{},
		},
		{
			testCase: "moved folder tree",
			items: []models.DriveItemable{
				driveItem("folder", "folder", testBaseDrivePath, false, true, false),
			},
			inputFolderMap: map[string]string{
				"folder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/a-folder",
				)[0],
				"subfolder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/a-folder/subfolder",
				)[0],
			},
			scope:  anyFolder,
			expect: assert.NoError,
			expectedCollectionPaths: expectedPathAsSlice(
				suite.T(),
				tenant,
				user,
				testBaseDrivePath,
			),
			expectedItemCount:      1,
			expectedFileCount:      0,
			expectedContainerCount: 1,
			expectedMetadataPaths: map[string]string{
				"folder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder",
				)[0],
				"subfolder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder/subfolder",
				)[0],
			},
			expectedExcludes: map[string]struct{}{},
		},
		{
			testCase: "moved folder tree and subfolder 1",
			items: []models.DriveItemable{
				driveItem("folder", "folder", testBaseDrivePath, false, true, false),
				driveItem("subfolder", "subfolder", testBaseDrivePath, false, true, false),
			},
			inputFolderMap: map[string]string{
				"folder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/a-folder",
				)[0],
				"subfolder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/a-folder/subfolder",
				)[0],
			},
			scope:  anyFolder,
			expect: assert.NoError,
			expectedCollectionPaths: expectedPathAsSlice(
				suite.T(),
				tenant,
				user,
				testBaseDrivePath,
			),
			expectedItemCount:      2,
			expectedFileCount:      0,
			expectedContainerCount: 1,
			expectedMetadataPaths: map[string]string{
				"folder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder",
				)[0],
				"subfolder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/subfolder",
				)[0],
			},
			expectedExcludes: map[string]struct{}{},
		},
		{
			testCase: "moved folder tree and subfolder 2",
			items: []models.DriveItemable{
				driveItem("subfolder", "subfolder", testBaseDrivePath, false, true, false),
				driveItem("folder", "folder", testBaseDrivePath, false, true, false),
			},
			inputFolderMap: map[string]string{
				"folder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/a-folder",
				)[0],
				"subfolder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/a-folder/subfolder",
				)[0],
			},
			scope:  anyFolder,
			expect: assert.NoError,
			expectedCollectionPaths: expectedPathAsSlice(
				suite.T(),
				tenant,
				user,
				testBaseDrivePath,
			),
			expectedItemCount:      2,
			expectedFileCount:      0,
			expectedContainerCount: 1,
			expectedMetadataPaths: map[string]string{
				"folder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder",
				)[0],
				"subfolder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/subfolder",
				)[0],
			},
			expectedExcludes: map[string]struct{}{},
		},
		{
			testCase: "deleted folder and package",
			items: []models.DriveItemable{
				delItem("folder", testBaseDrivePath, false, true, false),
				delItem("package", testBaseDrivePath, false, false, true),
			},
			inputFolderMap: map[string]string{
				"folder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder",
				)[0],
				"package": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/package",
				)[0],
			},
			scope:                   anyFolder,
			expect:                  assert.NoError,
			expectedCollectionPaths: []string{},
			expectedItemCount:       0,
			expectedFileCount:       0,
			expectedContainerCount:  0,
			expectedMetadataPaths:   map[string]string{},
			expectedExcludes:        map[string]struct{}{},
		},
		{
			testCase: "delete folder tree move subfolder",
			items: []models.DriveItemable{
				delItem("folder", testBaseDrivePath, false, true, false),
				driveItem("subfolder", "subfolder", testBaseDrivePath, false, true, false),
			},
			inputFolderMap: map[string]string{
				"folder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder",
				)[0],
				"subfolder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/folder/subfolder",
				)[0],
			},
			scope:  anyFolder,
			expect: assert.NoError,
			expectedCollectionPaths: expectedPathAsSlice(
				suite.T(),
				tenant,
				user,
				testBaseDrivePath,
			),
			expectedItemCount:      1,
			expectedFileCount:      0,
			expectedContainerCount: 1,
			expectedMetadataPaths: map[string]string{
				"subfolder": expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath+"/subfolder",
				)[0],
			},
			expectedExcludes: map[string]struct{}{},
		},
		{
			testCase: "delete file",
			items: []models.DriveItemable{
				delItem("item", testBaseDrivePath, true, false, false),
			},
			inputFolderMap:          map[string]string{},
			scope:                   anyFolder,
			expect:                  assert.NoError,
			expectedCollectionPaths: []string{},
			expectedItemCount:       1,
			expectedFileCount:       1,
			expectedContainerCount:  0,
			expectedMetadataPaths:   map[string]string{},
			expectedExcludes: map[string]struct{}{
				"item": {},
			},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.testCase, func(t *testing.T) {
			ctx, flush := tester.NewContext()
			defer flush()

			excludes := map[string]struct{}{}
			outputFolderMap := map[string]string{}
			maps.Copy(outputFolderMap, tt.inputFolderMap)
			c := NewCollections(
				graph.HTTPClient(graph.NoTimeout()),
				tenant,
				user,
				OneDriveSource,
				testFolderMatcher{tt.scope},
				&MockGraphService{},
				nil,
				control.Options{ToggleFeatures: control.Toggles{EnablePermissionsBackup: true}})

			err := c.UpdateCollections(
				ctx,
				"driveID",
				"General",
				tt.items,
				tt.inputFolderMap,
				outputFolderMap,
				excludes,
			)
			tt.expect(t, err)
			assert.Equal(t, len(tt.expectedCollectionPaths), len(c.CollectionMap), "collection paths")
			assert.Equal(t, tt.expectedItemCount, c.NumItems, "item count")
			assert.Equal(t, tt.expectedFileCount, c.NumFiles, "file count")
			assert.Equal(t, tt.expectedContainerCount, c.NumContainers, "container count")
			for _, collPath := range tt.expectedCollectionPaths {
				assert.Contains(t, c.CollectionMap, collPath)
			}

			assert.Equal(t, tt.expectedMetadataPaths, outputFolderMap)
			assert.Equal(t, tt.expectedExcludes, excludes)
		})
	}
}

func (suite *OneDriveCollectionsSuite) TestDeserializeMetadata() {
	tenant := "a-tenant"
	user := "a-user"
	driveID1 := "1"
	driveID2 := "2"
	deltaURL1 := "url/1"
	deltaURL2 := "url/2"

	folderID1 := "folder1"
	folderID2 := "folder2"
	path1 := "folder1/path"
	path2 := "folder2/path"

	table := []struct {
		name string
		// Each function returns the set of files for a single data.Collection.
		cols           []func() []graph.MetadataCollectionEntry
		expectedDeltas map[string]string
		expectedPaths  map[string]map[string]string
		errCheck       assert.ErrorAssertionFunc
	}{
		{
			name: "SuccessOneDriveAllOneCollection",
			cols: []func() []graph.MetadataCollectionEntry{
				func() []graph.MetadataCollectionEntry {
					return []graph.MetadataCollectionEntry{
						graph.NewMetadataEntry(
							graph.DeltaURLsFileName,
							map[string]string{driveID1: deltaURL1},
						),
						graph.NewMetadataEntry(
							graph.PreviousPathFileName,
							map[string]map[string]string{
								driveID1: {
									folderID1: path1,
								},
							},
						),
					}
				},
			},
			expectedDeltas: map[string]string{
				driveID1: deltaURL1,
			},
			expectedPaths: map[string]map[string]string{
				driveID1: {
					folderID1: path1,
				},
			},
			errCheck: assert.NoError,
		},
		{
			name: "MissingPaths",
			cols: []func() []graph.MetadataCollectionEntry{
				func() []graph.MetadataCollectionEntry {
					return []graph.MetadataCollectionEntry{
						graph.NewMetadataEntry(
							graph.DeltaURLsFileName,
							map[string]string{driveID1: deltaURL1},
						),
					}
				},
			},
			expectedDeltas: map[string]string{},
			expectedPaths:  map[string]map[string]string{},
			errCheck:       assert.NoError,
		},
		{
			name: "MissingDeltas",
			cols: []func() []graph.MetadataCollectionEntry{
				func() []graph.MetadataCollectionEntry {
					return []graph.MetadataCollectionEntry{
						graph.NewMetadataEntry(
							graph.PreviousPathFileName,
							map[string]map[string]string{
								driveID1: {
									folderID1: path1,
								},
							},
						),
					}
				},
			},
			expectedDeltas: map[string]string{},
			expectedPaths:  map[string]map[string]string{},
			errCheck:       assert.NoError,
		},
		{
			// An empty path map but valid delta results in metadata being returned
			// since it's possible to have a drive with no folders other than the
			// root.
			name: "EmptyPaths",
			cols: []func() []graph.MetadataCollectionEntry{
				func() []graph.MetadataCollectionEntry {
					return []graph.MetadataCollectionEntry{
						graph.NewMetadataEntry(
							graph.DeltaURLsFileName,
							map[string]string{driveID1: deltaURL1},
						),
						graph.NewMetadataEntry(
							graph.PreviousPathFileName,
							map[string]map[string]string{
								driveID1: {},
							},
						),
					}
				},
			},
			expectedDeltas: map[string]string{driveID1: deltaURL1},
			expectedPaths:  map[string]map[string]string{driveID1: {}},
			errCheck:       assert.NoError,
		},
		{
			// An empty delta map but valid path results in no metadata for that drive
			// being returned since the path map is only useful if we have a valid
			// delta.
			name: "EmptyDeltas",
			cols: []func() []graph.MetadataCollectionEntry{
				func() []graph.MetadataCollectionEntry {
					return []graph.MetadataCollectionEntry{
						graph.NewMetadataEntry(
							graph.DeltaURLsFileName,
							map[string]string{
								driveID1: "",
							},
						),
						graph.NewMetadataEntry(
							graph.PreviousPathFileName,
							map[string]map[string]string{
								driveID1: {
									folderID1: path1,
								},
							},
						),
					}
				},
			},
			expectedDeltas: map[string]string{},
			expectedPaths:  map[string]map[string]string{},
			errCheck:       assert.NoError,
		},
		{
			name: "SuccessTwoDrivesTwoCollections",
			cols: []func() []graph.MetadataCollectionEntry{
				func() []graph.MetadataCollectionEntry {
					return []graph.MetadataCollectionEntry{
						graph.NewMetadataEntry(
							graph.DeltaURLsFileName,
							map[string]string{driveID1: deltaURL1},
						),
						graph.NewMetadataEntry(
							graph.PreviousPathFileName,
							map[string]map[string]string{
								driveID1: {
									folderID1: path1,
								},
							},
						),
					}
				},
				func() []graph.MetadataCollectionEntry {
					return []graph.MetadataCollectionEntry{
						graph.NewMetadataEntry(
							graph.DeltaURLsFileName,
							map[string]string{driveID2: deltaURL2},
						),
						graph.NewMetadataEntry(
							graph.PreviousPathFileName,
							map[string]map[string]string{
								driveID2: {
									folderID2: path2,
								},
							},
						),
					}
				},
			},
			expectedDeltas: map[string]string{
				driveID1: deltaURL1,
				driveID2: deltaURL2,
			},
			expectedPaths: map[string]map[string]string{
				driveID1: {
					folderID1: path1,
				},
				driveID2: {
					folderID2: path2,
				},
			},
			errCheck: assert.NoError,
		},
		{
			// Bad formats are logged but skip adding entries to the maps and don't
			// return an error.
			name: "BadFormat",
			cols: []func() []graph.MetadataCollectionEntry{
				func() []graph.MetadataCollectionEntry {
					return []graph.MetadataCollectionEntry{
						graph.NewMetadataEntry(
							graph.PreviousPathFileName,
							map[string]string{driveID1: deltaURL1},
						),
					}
				},
			},
			expectedDeltas: map[string]string{},
			expectedPaths:  map[string]map[string]string{},
			errCheck:       assert.NoError,
		},
		{
			// Unexpected files are logged and skipped. They don't cause an error to
			// be returned.
			name: "BadFileName",
			cols: []func() []graph.MetadataCollectionEntry{
				func() []graph.MetadataCollectionEntry {
					return []graph.MetadataCollectionEntry{
						graph.NewMetadataEntry(
							graph.DeltaURLsFileName,
							map[string]string{driveID1: deltaURL1},
						),
						graph.NewMetadataEntry(
							graph.PreviousPathFileName,
							map[string]map[string]string{
								driveID1: {
									folderID1: path1,
								},
							},
						),
						graph.NewMetadataEntry(
							"foo",
							map[string]string{driveID1: deltaURL1},
						),
					}
				},
			},
			expectedDeltas: map[string]string{
				driveID1: deltaURL1,
			},
			expectedPaths: map[string]map[string]string{
				driveID1: {
					folderID1: path1,
				},
			},
			errCheck: assert.NoError,
		},
		{
			name: "DriveAlreadyFound_Paths",
			cols: []func() []graph.MetadataCollectionEntry{
				func() []graph.MetadataCollectionEntry {
					return []graph.MetadataCollectionEntry{
						graph.NewMetadataEntry(
							graph.DeltaURLsFileName,
							map[string]string{driveID1: deltaURL1},
						),
						graph.NewMetadataEntry(
							graph.PreviousPathFileName,
							map[string]map[string]string{
								driveID1: {
									folderID1: path1,
								},
							},
						),
					}
				},
				func() []graph.MetadataCollectionEntry {
					return []graph.MetadataCollectionEntry{
						graph.NewMetadataEntry(
							graph.PreviousPathFileName,
							map[string]map[string]string{
								driveID1: {
									folderID2: path2,
								},
							},
						),
					}
				},
			},
			expectedDeltas: nil,
			expectedPaths:  nil,
			errCheck:       assert.Error,
		},
		{
			name: "DriveAlreadyFound_Deltas",
			cols: []func() []graph.MetadataCollectionEntry{
				func() []graph.MetadataCollectionEntry {
					return []graph.MetadataCollectionEntry{
						graph.NewMetadataEntry(
							graph.DeltaURLsFileName,
							map[string]string{driveID1: deltaURL1},
						),
						graph.NewMetadataEntry(
							graph.PreviousPathFileName,
							map[string]map[string]string{
								driveID1: {
									folderID1: path1,
								},
							},
						),
					}
				},
				func() []graph.MetadataCollectionEntry {
					return []graph.MetadataCollectionEntry{
						graph.NewMetadataEntry(
							graph.DeltaURLsFileName,
							map[string]string{driveID1: deltaURL2},
						),
					}
				},
			},
			expectedDeltas: nil,
			expectedPaths:  nil,
			errCheck:       assert.Error,
		},
	}

	for _, test := range table {
		suite.T().Run(test.name, func(t *testing.T) {
			ctx, flush := tester.NewContext()
			defer flush()

			cols := []data.Collection{}

			for _, c := range test.cols {
				mc, err := graph.MakeMetadataCollection(
					tenant,
					user,
					path.OneDriveService,
					path.FilesCategory,
					c(),
					func(*support.ConnectorOperationStatus) {},
				)
				require.NoError(t, err)

				cols = append(cols, mc)
			}

			deltas, paths, err := deserializeMetadata(ctx, cols)
			test.errCheck(t, err)

			assert.Equal(t, test.expectedDeltas, deltas)
			assert.Equal(t, test.expectedPaths, paths)
		})
	}
}

type mockDeltaPageLinker struct {
	link  *string
	delta *string
}

func (pl *mockDeltaPageLinker) GetOdataNextLink() *string {
	return pl.link
}

func (pl *mockDeltaPageLinker) GetOdataDeltaLink() *string {
	return pl.delta
}

type deltaPagerResult struct {
	items     []models.DriveItemable
	nextLink  *string
	deltaLink *string
	err       error
}

type mockItemPager struct {
	// DriveID -> set of return values for queries for that drive.
	toReturn []deltaPagerResult
	getIdx   int
}

func (p *mockItemPager) GetPage(context.Context) (gapi.DeltaPageLinker, error) {
	if len(p.toReturn) <= p.getIdx {
		return nil, assert.AnError
	}

	idx := p.getIdx
	p.getIdx++

	return &mockDeltaPageLinker{
		p.toReturn[idx].nextLink,
		p.toReturn[idx].deltaLink,
	}, p.toReturn[idx].err
}

func (p *mockItemPager) SetNext(string) {}

func (p *mockItemPager) ValuesIn(gapi.DeltaPageLinker) ([]models.DriveItemable, error) {
	idx := p.getIdx
	if idx > 0 {
		// Return values lag by one since we increment in GetPage().
		idx--
	}

	if len(p.toReturn) <= idx {
		return nil, assert.AnError
	}

	return p.toReturn[idx].items, nil
}

func (suite *OneDriveCollectionsSuite) TestGet() {
	anyFolder := (&selectors.OneDriveBackup{}).Folders(selectors.Any())[0]

	tenant := "a-tenant"
	user := "a-user"

	metadataPath, err := path.Builder{}.ToServiceCategoryMetadataPath(
		tenant,
		user,
		path.OneDriveService,
		path.FilesCategory,
		false,
	)
	require.NoError(suite.T(), err, "making metadata path")

	rootFolderPath := expectedPathAsSlice(
		suite.T(),
		tenant,
		user,
		testBaseDrivePath,
	)[0]
	folderPath := expectedPathAsSlice(
		suite.T(),
		tenant,
		user,
		testBaseDrivePath+"/folder",
	)[0]

	empty := ""
	next := "next"
	delta := "delta1"
	delta2 := "delta2"

	driveID1 := uuid.NewString()
	drive1 := models.NewDrive()
	drive1.SetId(&driveID1)
	drive1.SetName(&driveID1)

	driveID2 := uuid.NewString()
	drive2 := models.NewDrive()
	drive2.SetId(&driveID2)
	drive2.SetName(&driveID2)

	driveBasePath2 := "drive/driveID2/root:"

	rootFolderPath2 := expectedPathAsSlice(
		suite.T(),
		tenant,
		user,
		driveBasePath2,
	)[0]
	folderPath2 := expectedPathAsSlice(
		suite.T(),
		tenant,
		user,
		driveBasePath2+"/folder",
	)[0]

	table := []struct {
		name     string
		drives   []models.Driveable
		items    map[string][]deltaPagerResult
		errCheck assert.ErrorAssertionFunc
		// Collection name -> set of item IDs. We can't check item data because
		// that's not mocked out. Metadata is checked separately.
		expectedCollections map[string][]string
		expectedDeltaURLs   map[string]string
		expectedFolderPaths map[string]map[string]string
		expectedDelList     map[string]struct{}
	}{
		{
			name:   "OneDrive_OneItemPage_DelFileOnly_NoFolders_NoErrors",
			drives: []models.Driveable{drive1},
			items: map[string][]deltaPagerResult{
				driveID1: {
					{
						items: []models.DriveItemable{
							delItem("file", testBaseDrivePath, true, false, false),
						},
						deltaLink: &delta,
					},
				},
			},
			errCheck:            assert.NoError,
			expectedCollections: map[string][]string{},
			expectedDeltaURLs: map[string]string{
				driveID1: delta,
			},
			expectedFolderPaths: map[string]map[string]string{
				// We need an empty map here so deserializing metadata knows the delta
				// token for this drive is valid.
				driveID1: {},
			},
			expectedDelList: map[string]struct{}{
				"file": {},
			},
		},
		{
			name:   "OneDrive_OneItemPage_NoFolders_NoErrors",
			drives: []models.Driveable{drive1},
			items: map[string][]deltaPagerResult{
				driveID1: {
					{
						items: []models.DriveItemable{
							driveItem("file", "file", testBaseDrivePath, true, false, false),
						},
						deltaLink: &delta,
					},
				},
			},
			errCheck: assert.NoError,
			expectedCollections: map[string][]string{
				expectedPathAsSlice(
					suite.T(),
					tenant,
					user,
					testBaseDrivePath,
				)[0]: {"file"},
			},
			expectedDeltaURLs: map[string]string{
				driveID1: delta,
			},
			expectedFolderPaths: map[string]map[string]string{
				// We need an empty map here so deserializing metadata knows the delta
				// token for this drive is valid.
				driveID1: {},
			},
			expectedDelList: map[string]struct{}{},
		},
		{
			name:   "OneDrive_OneItemPage_NoErrors",
			drives: []models.Driveable{drive1},
			items: map[string][]deltaPagerResult{
				driveID1: {
					{
						items: []models.DriveItemable{
							driveItem("folder", "folder", testBaseDrivePath, false, true, false),
							driveItem("file", "file", testBaseDrivePath+"/folder", true, false, false),
						},
						deltaLink: &delta,
					},
				},
			},
			errCheck: assert.NoError,
			expectedCollections: map[string][]string{
				folderPath:     {"file"},
				rootFolderPath: {"folder"},
			},
			expectedDeltaURLs: map[string]string{
				driveID1: delta,
			},
			expectedFolderPaths: map[string]map[string]string{
				driveID1: {
					"folder": folderPath,
				},
			},
			expectedDelList: map[string]struct{}{},
		},
		{
			name:   "OneDrive_OneItemPage_EmptyDelta_NoErrors",
			drives: []models.Driveable{drive1},
			items: map[string][]deltaPagerResult{
				driveID1: {
					{
						items: []models.DriveItemable{
							driveItem("folder", "folder", testBaseDrivePath, false, true, false),
							driveItem("file", "file", testBaseDrivePath+"/folder", true, false, false),
						},
						deltaLink: &empty,
					},
				},
			},
			errCheck: assert.NoError,
			expectedCollections: map[string][]string{
				folderPath:     {"file"},
				rootFolderPath: {"folder"},
			},
			expectedDeltaURLs:   map[string]string{},
			expectedFolderPaths: map[string]map[string]string{},
			expectedDelList:     map[string]struct{}{},
		},
		{
			name:   "OneDrive_TwoItemPages_NoErrors",
			drives: []models.Driveable{drive1},
			items: map[string][]deltaPagerResult{
				driveID1: {
					{
						items: []models.DriveItemable{
							driveItem("folder", "folder", testBaseDrivePath, false, true, false),
							driveItem("file", "file", testBaseDrivePath+"/folder", true, false, false),
						},
						nextLink: &next,
					},
					{
						items: []models.DriveItemable{
							driveItem("folder", "folder", testBaseDrivePath, false, true, false),
							driveItem("file2", "file2", testBaseDrivePath+"/folder", true, false, false),
						},
						deltaLink: &delta,
					},
				},
			},
			errCheck: assert.NoError,
			expectedCollections: map[string][]string{
				folderPath:     {"file", "file2"},
				rootFolderPath: {"folder"},
			},
			expectedDeltaURLs: map[string]string{
				driveID1: delta,
			},
			expectedFolderPaths: map[string]map[string]string{
				driveID1: {
					"folder": folderPath,
				},
			},
			expectedDelList: map[string]struct{}{},
		},
		{
			name: "TwoDrives_OneItemPageEach_NoErrors",
			drives: []models.Driveable{
				drive1,
				drive2,
			},
			items: map[string][]deltaPagerResult{
				driveID1: {
					{
						items: []models.DriveItemable{
							driveItem("folder", "folder", testBaseDrivePath, false, true, false),
							driveItem("file", "file", testBaseDrivePath+"/folder", true, false, false),
						},
						deltaLink: &delta,
					},
				},
				driveID2: {
					{
						items: []models.DriveItemable{
							driveItem("folder", "folder", driveBasePath2, false, true, false),
							driveItem("file", "file", driveBasePath2+"/folder", true, false, false),
						},
						deltaLink: &delta2,
					},
				},
			},
			errCheck: assert.NoError,
			expectedCollections: map[string][]string{
				folderPath:      {"file"},
				folderPath2:     {"file"},
				rootFolderPath:  {"folder"},
				rootFolderPath2: {"folder"},
			},
			expectedDeltaURLs: map[string]string{
				driveID1: delta,
				driveID2: delta2,
			},
			expectedFolderPaths: map[string]map[string]string{
				driveID1: {
					"folder": folderPath,
				},
				driveID2: {
					"folder": folderPath2,
				},
			},
			expectedDelList: map[string]struct{}{},
		},
		{
			name:   "OneDrive_OneItemPage_Errors",
			drives: []models.Driveable{drive1},
			items: map[string][]deltaPagerResult{
				driveID1: {
					{
						err: assert.AnError,
					},
				},
			},
			errCheck:            assert.Error,
			expectedCollections: nil,
			expectedDeltaURLs:   nil,
			expectedFolderPaths: nil,
			expectedDelList:     nil,
		},
	}
	for _, test := range table {
		suite.T().Run(test.name, func(t *testing.T) {
			ctx, flush := tester.NewContext()
			defer flush()

			drivePagerFunc := func(
				source driveSource,
				servicer graph.Servicer,
				resourceOwner string,
				fields []string,
			) (drivePager, error) {
				return &mockDrivePager{
					toReturn: []pagerResult{
						{
							drives: test.drives,
						},
					},
				}, nil
			}

			itemPagerFunc := func(
				servicer graph.Servicer,
				driveID, link string,
			) itemPager {
				return &mockItemPager{
					toReturn: test.items[driveID],
				}
			}

			c := NewCollections(
				graph.HTTPClient(graph.NoTimeout()),
				tenant,
				user,
				OneDriveSource,
				testFolderMatcher{anyFolder},
				&MockGraphService{},
				func(*support.ConnectorOperationStatus) {},
				control.Options{ToggleFeatures: control.Toggles{EnablePermissionsBackup: true}},
			)
			c.drivePagerFunc = drivePagerFunc
			c.itemPagerFunc = itemPagerFunc

			// TODO(ashmrtn): Allow passing previous metadata.
			cols, _, err := c.Get(ctx, nil)
			test.errCheck(t, err)

			if err != nil {
				return
			}

			for _, baseCol := range cols {
				folderPath := baseCol.FullPath().String()
				if folderPath == metadataPath.String() {
					deltas, paths, err := deserializeMetadata(ctx, []data.Collection{baseCol})
					if !assert.NoError(t, err, "deserializing metadata") {
						continue
					}

					assert.Equal(t, test.expectedDeltaURLs, deltas)
					assert.Equal(t, test.expectedFolderPaths, paths)

					continue
				}

				// TODO(ashmrtn): We should really be getting items in the collection
				// via the Items() channel, but we don't have a way to mock out the
				// actual item fetch yet (mostly wiring issues). The lack of that makes
				// this check a bit more bittle since internal details can change.
				col, ok := baseCol.(*Collection)
				require.True(t, ok, "getting onedrive.Collection handle")

				itemIDs := make([]string, 0, len(col.driveItems))

				for id := range col.driveItems {
					itemIDs = append(itemIDs, id)
				}

				assert.ElementsMatch(t, test.expectedCollections[folderPath], itemIDs)
			}

			// TODO(ashmrtn): Uncomment this when we begin return the set of items to
			// remove from the upcoming backup.
			// assert.Equal(t, test.expectedDelList, delList)
		})
	}
}

func driveItem(
	id string,
	name string,
	parentPath string,
	isFile, isFolder, isPackage bool,
) models.DriveItemable {
	item := models.NewDriveItem()
	item.SetName(&name)
	item.SetId(&id)

	parentReference := models.NewItemReference()
	parentReference.SetPath(&parentPath)
	item.SetParentReference(parentReference)

	switch {
	case isFile:
		item.SetFile(models.NewFile())
	case isFolder:
		item.SetFolder(models.NewFolder())
	case isPackage:
		item.SetPackage(models.NewPackage_escaped())
	}

	return item
}

// delItem creates a DriveItemable that is marked as deleted. path must be set
// to the base drive path.
func delItem(
	id string,
	parentPath string,
	isFile, isFolder, isPackage bool,
) models.DriveItemable {
	item := models.NewDriveItem()
	item.SetId(&id)
	item.SetDeleted(models.NewDeleted())

	parentReference := models.NewItemReference()
	parentReference.SetPath(&parentPath)
	item.SetParentReference(parentReference)

	switch {
	case isFile:
		item.SetFile(models.NewFile())
	case isFolder:
		item.SetFolder(models.NewFolder())
	case isPackage:
		item.SetPackage(models.NewPackage_escaped())
	}

	return item
}
