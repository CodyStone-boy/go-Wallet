package main

import (
	"fmt"

	"github.com/hiromaily/go-bitcoin/pkg/config"
	"github.com/hiromaily/go-bitcoin/pkg/db/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model"
	"github.com/hiromaily/go-bitcoin/pkg/txfile"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/key"
)

// Registry is for registry interface
type Registry interface {
	NewSigner() wallet.Signer
}

type registry struct {
	conf       *config.Config
	walletType wallet.WalletType
}

// NewRegistry is to register regstry interface
func NewRegistry(conf *config.Config, walletType wallet.WalletType) Registry {
	return &registry{
		conf:       conf,
		walletType: walletType,
	}
}

// NewSigner is to register for Signer interface
func (r *registry) NewSigner() wallet.Signer {
	r.newLogger()
	r.setFilePath()

	return wallet.NewWallet(
		r.newBTC(),
		r.newStorager(),
		r.walletType,
	)
}

func (r *registry) newBTC() api.Bitcoiner {
	// Connection to Bitcoin core
	bit, err := api.Connection(&r.conf.Bitcoin, enum.CoinType(r.conf.CoinType))
	if err != nil {
		panic(fmt.Sprintf("btc.Connection error: %s", err))
	}
	return bit
}

//TODO: change to interface
func (r *registry) newStorager() *model.DB {
	// MySQL
	rds, err := rdb.NewMySQL(&r.conf.MySQL)
	if err != nil {
		panic(fmt.Sprintf("rds.Connection() error: %s", err))
	}
	return model.NewDB(rds)
}

//TODO: change to interface
func (r *registry) newLogger() {
	// logger
	logger.NewLogger(&r.conf.Logger)
}

//TODO: move to somewhere
func (r *registry) setFilePath() {
	// TxFile
	if r.conf.TxFile.BasePath != "" {
		txfile.SetFilePath(r.conf.TxFile.BasePath)
	}

	// PubkeyCSV
	if r.conf.PubkeyFile.BasePath != "" {
		key.SetFilePath(r.conf.PubkeyFile.BasePath)
	}
}
