package eth

import (
	"context"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

// Ethereum includes client to call Json-RPC
type Ethereum struct {
	client       *ethrpc.Client
	coinTypeCode coin.CoinTypeCode //eth
	logger       *zap.Logger
	ctx          context.Context
	netID        uint16
	version      string
	gethDir      string
}

// NewEthereum creates ethereum object
func NewEthereum(
	ctx context.Context,
	client *ethrpc.Client,
	coinTypeCode coin.CoinTypeCode,
	conf *config.Ethereum,
	logger *zap.Logger) (*Ethereum, error) {

	eth := &Ethereum{
		client:       client,
		coinTypeCode: coinTypeCode,
		logger:       logger,
		ctx:          ctx,
		gethDir:      conf.Geth.DirName,
	}

	// get NetID
	netID, err := eth.NetVersion()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call eth.NetVersion()")
	}
	eth.netID = netID

	// get version
	clientVer, err := eth.ClientVersion()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call eth.ClientVersion()")
	}
	eth.version = clientVer

	// check sync progress
	res, isDone, err := eth.Syncing()
	if err != nil {
		return nil, errors.Wrap(err, "fail to call eth.Syncing()")
	}
	if !isDone {
		logger.Warn("sync is not completed yet")
	}
	logger.Info("result_syncing",
		zap.Int64("knownStates", res.KnownStates),
		zap.Int64("pulledStates", res.PulledStates),
		zap.Int64("startingBlock", res.StartingBlock),
		zap.Int64("currentBlock", res.CurrentBlock),
		zap.Int64("highestBlock", res.HighestBlock),
	)

	return eth, nil
}

// Close disconnect to server
func (e *Ethereum) Close() {
	if e.client != nil {
		e.client.Close()
	}
}
