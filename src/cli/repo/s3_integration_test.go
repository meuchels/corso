package repo_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/alcionai/corso/src/cli"
	"github.com/alcionai/corso/src/cli/config"
	"github.com/alcionai/corso/src/internal/tester"
	"github.com/alcionai/corso/src/internal/tester/aw"
	"github.com/alcionai/corso/src/pkg/account"
	"github.com/alcionai/corso/src/pkg/control"
	"github.com/alcionai/corso/src/pkg/repository"
)

type S3IntegrationSuite struct {
	suite.Suite
}

func TestS3IntegrationSuite(t *testing.T) {
	tester.RunOnAny(
		t,
		tester.CorsoCITests,
		tester.CorsoCLITests,
		tester.CorsoCLIRepoTests)

	suite.Run(t, new(S3IntegrationSuite))
}

func (suite *S3IntegrationSuite) SetupSuite() {
	tester.MustGetEnvSets(suite.T(), tester.AWSStorageCredEnvs, tester.M365AcctCredEnvs)
}

func (suite *S3IntegrationSuite) TestInitS3Cmd() {
	table := []struct {
		name         string
		bucketPrefix string
	}{
		{
			name:         "NoPrefix",
			bucketPrefix: "",
		},
		{
			name:         "S3Prefix",
			bucketPrefix: "s3://",
		},
	}

	for _, test := range table {
		suite.T().Run(test.name, func(t *testing.T) {
			ctx, flush := tester.NewContext()
			defer flush()

			st := tester.NewPrefixedS3Storage(t)
			cfg, err := st.S3Config()
			aw.MustNoErr(t, err)

			vpr, configFP := tester.MakeTempTestConfigClone(t, nil)

			ctx = config.SetViper(ctx, vpr)

			cmd := tester.StubRootCmd(
				"repo", "init", "s3",
				"--config-file", configFP,
				"--bucket", test.bucketPrefix+cfg.Bucket,
				"--prefix", cfg.Prefix)
			cli.BuildCommandTree(cmd)

			// run the command
			aw.MustNoErr(t, cmd.ExecuteContext(ctx))

			// a second initialization should result in an error
			err = cmd.ExecuteContext(ctx)
			aw.Err(t, err)
			aw.ErrIs(t, err, repository.ErrorRepoAlreadyExists)
		})
	}
}

func (suite *S3IntegrationSuite) TestInitMultipleTimes() {
	t := suite.T()
	ctx, flush := tester.NewContext()

	defer flush()

	st := tester.NewPrefixedS3Storage(t)
	cfg, err := st.S3Config()
	aw.MustNoErr(t, err)

	vpr, configFP := tester.MakeTempTestConfigClone(t, nil)

	ctx = config.SetViper(ctx, vpr)

	for i := 0; i < 2; i++ {
		cmd := tester.StubRootCmd(
			"repo", "init", "s3",
			"--config-file", configFP,
			"--bucket", cfg.Bucket,
			"--prefix", cfg.Prefix,
			"--succeed-if-exists",
		)
		cli.BuildCommandTree(cmd)

		// run the command
		aw.MustNoErr(t, cmd.ExecuteContext(ctx))
	}
}

func (suite *S3IntegrationSuite) TestInitS3Cmd_missingBucket() {
	t := suite.T()
	ctx, flush := tester.NewContext()

	defer flush()

	st := tester.NewPrefixedS3Storage(t)
	cfg, err := st.S3Config()
	aw.MustNoErr(t, err)

	vpr, configFP := tester.MakeTempTestConfigClone(t, nil)

	ctx = config.SetViper(ctx, vpr)

	cmd := tester.StubRootCmd(
		"repo", "init", "s3",
		"--config-file", configFP,
		"--prefix", cfg.Prefix)
	cli.BuildCommandTree(cmd)

	// run the command
	aw.MustErr(t, cmd.ExecuteContext(ctx))
}

func (suite *S3IntegrationSuite) TestConnectS3Cmd() {
	table := []struct {
		name         string
		bucketPrefix string
	}{
		{
			name:         "NoPrefix",
			bucketPrefix: "",
		},
		{
			name:         "S3Prefix",
			bucketPrefix: "s3://",
		},
	}

	for _, test := range table {
		suite.T().Run(test.name, func(t *testing.T) {
			ctx, flush := tester.NewContext()
			defer flush()

			st := tester.NewPrefixedS3Storage(t)
			cfg, err := st.S3Config()
			aw.MustNoErr(t, err)

			force := map[string]string{
				tester.TestCfgAccountProvider: "M365",
				tester.TestCfgStorageProvider: "S3",
				tester.TestCfgPrefix:          cfg.Prefix,
			}
			vpr, configFP := tester.MakeTempTestConfigClone(t, force)

			ctx = config.SetViper(ctx, vpr)

			// init the repo first
			_, err = repository.Initialize(ctx, account.Account{}, st, control.Options{})
			aw.MustNoErr(t, err)

			// then test it
			cmd := tester.StubRootCmd(
				"repo", "connect", "s3",
				"--config-file", configFP,
				"--bucket", test.bucketPrefix+cfg.Bucket,
				"--prefix", cfg.Prefix,
			)
			cli.BuildCommandTree(cmd)

			// run the command
			aw.NoErr(t, cmd.ExecuteContext(ctx))
		})
	}
}

func (suite *S3IntegrationSuite) TestConnectS3Cmd_BadBucket() {
	t := suite.T()
	ctx, flush := tester.NewContext()

	defer flush()

	st := tester.NewPrefixedS3Storage(t)
	cfg, err := st.S3Config()
	aw.MustNoErr(t, err)

	vpr, configFP := tester.MakeTempTestConfigClone(t, nil)

	ctx = config.SetViper(ctx, vpr)

	cmd := tester.StubRootCmd(
		"repo", "connect", "s3",
		"--config-file", configFP,
		"--bucket", "wrong",
		"--prefix", cfg.Prefix)
	cli.BuildCommandTree(cmd)

	// run the command
	aw.MustErr(t, cmd.ExecuteContext(ctx))
}

func (suite *S3IntegrationSuite) TestConnectS3Cmd_BadPrefix() {
	t := suite.T()
	ctx, flush := tester.NewContext()

	defer flush()

	st := tester.NewPrefixedS3Storage(t)
	cfg, err := st.S3Config()
	aw.MustNoErr(t, err)

	vpr, configFP := tester.MakeTempTestConfigClone(t, nil)

	ctx = config.SetViper(ctx, vpr)

	cmd := tester.StubRootCmd(
		"repo", "connect", "s3",
		"--config-file", configFP,
		"--bucket", cfg.Bucket,
		"--prefix", "wrong")
	cli.BuildCommandTree(cmd)

	// run the command
	aw.MustErr(t, cmd.ExecuteContext(ctx))
}
