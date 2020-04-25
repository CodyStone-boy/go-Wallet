package main

import (
	"database/sql"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/config"
	mysql "github.com/hiromaily/go-bitcoin/pkg/db/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/logger"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/walletrepo"
	"github.com/hiromaily/go-bitcoin/pkg/repository"
	"github.com/hiromaily/go-bitcoin/pkg/tracer"
	"github.com/hiromaily/go-bitcoin/pkg/tx"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/types"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/wallets/wallet"
)

// Registry is for registry interface
type Registry interface {
	NewWalleter() wallets.Walleter
}

type registry struct {
	conf         *config.Config
	mysqlXClient *sqlx.DB
	mysqlClient  *sql.DB
	logger       *zap.Logger
	rpcClient    *rpcclient.Client
	walletType   types.WalletType
}

// NewRegistry is to register registry interface
func NewRegistry(conf *config.Config, walletType types.WalletType) Registry {
	return &registry{
		conf:       conf,
		walletType: walletType,
	}
}

// NewWalleter is to register for walleter interface
func (r *registry) NewWalleter() wallets.Walleter {
	return wallet.NewWallet(
		r.newBTC(),
		r.newLogger(),
		r.newTracer(),
		r.newStorager(),
		r.newTxRepo(),
		r.newTxInputRepo(),
		r.newTxOutputRepo(),
		r.newAddressFileStorager(),
		r.newTxFileStorager(),
		r.walletType,
	)
}

func (r *registry) newRPCClient() *rpcclient.Client {
	var err error
	if r.rpcClient == nil {
		r.rpcClient, err = api.NewRPCClient(&r.conf.Bitcoin)
	}
	if err != nil {
		panic(err)
	}
	return r.rpcClient
}

func (r *registry) newBTC() api.Bitcoiner {
	bit, err := api.NewBitcoin(r.newRPCClient(), &r.conf.Bitcoin, r.newLogger(), r.conf.CoinTypeCode)
	if err != nil {
		panic(err)
	}
	return bit
}

func (r *registry) newLogger() *zap.Logger {
	if r.logger == nil {
		r.logger = logger.NewZapLogger(&r.conf.Logger)
	}
	return r.logger
}

func (r *registry) newTracer() opentracing.Tracer {
	return tracer.NewTracer(r.conf.Tracer)
}

func (r *registry) newStorager() rdb.WalletStorager {
	// if there are multiple options, set proper one
	// storager interface as MySQL
	return walletrepo.NewWalletRepository(
		r.newMySQLXClient(),
		r.newLogger(),
	)
}

func (r *registry) newTxRepo() repository.TxRepository {
	return repository.NewTxRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newTxInputRepo() repository.TxInputRepository {
	return repository.NewTxInputRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newTxOutputRepo() repository.TxOutputRepository {
	return repository.NewTxOutputRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newMySQLClient() *sql.DB {
	if r.mysqlClient == nil {
		dbConn, err := mysql.NewMySQL(&r.conf.MySQL)
		if err != nil {
			panic(err)
		}
		r.mysqlClient = dbConn
	}
	return r.mysqlClient
}

func (r *registry) newMySQLXClient() *sqlx.DB {
	if r.mysqlXClient == nil {
		dbConn, err := mysql.NewMySQLX(&r.conf.MySQL)
		if err != nil {
			panic(err)
		}
		r.mysqlXClient = dbConn
	}
	return r.mysqlXClient
}

func (r *registry) newAddressFileStorager() address.Storager {
	return address.NewFileRepository(
		r.conf.PubkeyFile.BasePath,
		r.newLogger(),
	)
}

func (r *registry) newTxFileStorager() tx.FileStorager {
	return tx.NewFileRepository(
		r.conf.TxFile.BasePath,
		r.newLogger(),
	)
}
