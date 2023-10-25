package testnode

import (
	"fmt"
	"time"

	"github.com/celestiaorg/celestia-app/app"
	"github.com/celestiaorg/celestia-app/node"
	"github.com/celestiaorg/celestia-app/pkg/appconsts"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	srvtypes "github.com/cosmos/cosmos-sdk/server/types"
	tmconfig "github.com/tendermint/tendermint/config"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"github.com/tendermint/tendermint/types"
)

// Config is the configuration of a test node.
type Config struct {
	// ChainID is the chain ID of the network.
	ChainID string
	// TmConfig is the Tendermint configuration used for the network.
	TmConfig *tmconfig.Config
	// AppConfig is the application configuration of the test node.
	AppConfig *srvconfig.Config
	// ConsensusParams are the consensus parameters of the test node.
	ConsensusParams *tmproto.ConsensusParams
	// AppOptions are the application options of the test node. Portions of the
	// app config will automatically be set into the app option when the app
	// config is set.
	AppOptions *app.KVAppOptions
	// GenesisOptions are the genesis options of the test node.
	GenesisOptions []GenesisOption
	// Accounts are the accounts of the test node.
	Accounts []string
	// AppCreator is used to create the application for the testnode.
	AppCreator srvtypes.AppCreator
	// SupressLogs
	SupressLogs bool
}

// WithChainID sets the ChainID and returns the Config.
func (c *Config) WithChainID(s string) *Config {
	c.ChainID = s
	return c
}

// WithTendermintConfig sets the TmConfig and returns the *Config.
func (c *Config) WithTendermintConfig(conf *tmconfig.Config) *Config {
	c.TmConfig = conf
	return c
}

// WithAppConfig sets the AppConfig and returns the Config.
//
// Warning: This method will also overwrite relevant portions of the app config
// to the app options. See the SetFromAppConfig method for more information on
// which values are overwritten.
func (c *Config) WithAppConfig(conf *srvconfig.Config) *Config {
	c.AppConfig = conf
	c.AppOptions.SetFromAppConfig(conf)
	return c
}

// WithConsensusParams sets the ConsensusParams and returns the Config.
func (c *Config) WithConsensusParams(params *tmproto.ConsensusParams) *Config {
	c.ConsensusParams = params
	return c
}

// WithAppOptions sets the AppOptions and returns the Config.
//
// Warning: If the app config is set after this, it could overwrite some values.
// See SetFromAppConfig for more information on which values are overwritten.
func (c *Config) WithAppOptions(opts *app.KVAppOptions) *Config {
	c.AppOptions = opts
	return c
}

// WithGenesisOptions sets the GenesisOptions and returns the Config.
func (c *Config) WithGenesisOptions(opts ...GenesisOption) *Config {
	c.GenesisOptions = opts
	return c
}

// WithAccounts sets the Accounts and returns the Config.
func (c *Config) WithAccounts(accs []string) *Config {
	c.Accounts = accs
	return c
}

// WithAppCreator sets the AppCreator and returns the Config.
func (c *Config) WithAppCreator(creator srvtypes.AppCreator) *Config {
	c.AppCreator = creator
	return c
}

// WithSupressLogs sets the SupressLogs and returns the Config.
func (c *Config) WithSupressLogs(sl bool) *Config {
	c.SupressLogs = sl
	return c
}

// WithTimeoutCommit sets the TimeoutCommit and returns the Config.
func (c *Config) WithTimeoutCommit(d time.Duration) *Config {
	c.TmConfig.Consensus.TimeoutCommit = d
	return c
}

func DefaultConfig() *Config {
	tmcfg := DefaultTendermintConfig()
	tmcfg.Consensus.TimeoutCommit = 1 * time.Millisecond
	cfg := &Config{}
	return cfg.
		WithAccounts([]string{}).
		WithChainID(tmrand.Str(6)).
		WithTendermintConfig(DefaultTendermintConfig()).
		WithConsensusParams(DefaultParams()).
		WithAppOptions(app.DefaultAppOptions()).
		WithAppConfig(DefaultAppConfig()).
		WithGenesisOptions().
		WithAppCreator(node.NewAppServer).
		WithSupressLogs(true)
}

func DefaultParams() *tmproto.ConsensusParams {
	cparams := types.DefaultConsensusParams()
	cparams.Block.TimeIotaMs = 1
	cparams.Block.MaxBytes = appconsts.DefaultMaxBytes
	return cparams
}

func DefaultTendermintConfig() *tmconfig.Config {
	tmCfg := tmconfig.DefaultConfig()
	// Reduce the target height duration so that blocks are produced faster
	// during tests.
	tmCfg.Consensus.TimeoutCommit = 100 * time.Millisecond
	tmCfg.Consensus.TimeoutPropose = 200 * time.Millisecond

	// set the mempool's MaxTxBytes to allow the testnode to accept a
	// transaction that fills the entire square. Any blob transaction larger
	// than the square size will still fail no matter what.
	tmCfg.Mempool.MaxTxBytes = appconsts.DefaultMaxBytes

	// remove all barriers from the testnode being able to accept very large
	// transactions and respond to very queries with large responses (~200MB was
	// chosen only as an arbitrary large number).
	tmCfg.RPC.MaxBodyBytes = 200_000_000

	// set all the ports to random open ones
	tmCfg.RPC.ListenAddress = fmt.Sprintf("tcp://127.0.0.1:%d", GetFreePort())
	tmCfg.P2P.ListenAddress = fmt.Sprintf("tcp://127.0.0.1:%d", GetFreePort())
	tmCfg.RPC.GRPCListenAddress = fmt.Sprintf("tcp://127.0.0.1:%d", GetFreePort())

	return tmCfg
}
