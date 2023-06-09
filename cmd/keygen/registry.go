package main

import (
	"database/sql"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/hiromaily/go-crypto-wallet/pkg/account"
	"github.com/hiromaily/go-crypto-wallet/pkg/address"
	"github.com/hiromaily/go-crypto-wallet/pkg/config"
	mysql "github.com/hiromaily/go-crypto-wallet/pkg/db/rdb"
	"github.com/hiromaily/go-crypto-wallet/pkg/logger"
	"github.com/hiromaily/go-crypto-wallet/pkg/repository/coldrepo"
	"github.com/hiromaily/go-crypto-wallet/pkg/tx"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/btcgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/ethgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/api/xrpgrp/xrp"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/coin"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/key"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/btc/coldsrv"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/btc/coldsrv/keygensrv"
	commonsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/coldsrv"
	ethsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/eth/keygensrv"
	xrpsrv "github.com/hiromaily/go-crypto-wallet/pkg/wallet/service/xrp/keygensrv"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/btcwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/ethwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/wallet/wallets/xrpwallet"
	"github.com/hiromaily/go-crypto-wallet/pkg/ws"
)

// Registry is for registry interface
type Registry interface {
	NewKeygener() wallets.Keygener
}

type registry struct {
	conf         *config.WalletRoot
	accountConf  *account.AccountRoot
	walletType   wallet.WalletType
	logger       *zap.Logger
	btc          btcgrp.Bitcoiner
	eth          ethgrp.Ethereumer
	xrp          xrpgrp.Rippler
	rpcClient    *rpcclient.Client
	rpcEthClient *ethrpc.Client
	wsXrpPublic  *ws.WS
	wsXrpAdmin   *ws.WS
	grpcConn     *grpc.ClientConn
	rippleAPI    *xrp.RippleAPI
	mysqlClient  *sql.DB
	multisig     account.MultisigAccounter
}

// NewRegistry is to register registry interface
func NewRegistry(conf *config.WalletRoot, accountConf *account.AccountRoot, walletType wallet.WalletType) Registry {
	return &registry{
		conf:        conf,
		accountConf: accountConf,
		walletType:  walletType,
	}
}

// NewKeygener is to register for keygener interface
// Which is better ?
// - create each interface getter to difine as interface
// - return struct itself
// func (r *registry) NewKeygener() *keygen.Keygen {
func (r *registry) NewKeygener() wallets.Keygener {
	switch {
	case coin.IsBTCGroup(r.conf.CoinTypeCode):
		return r.newBTCKeygener()
	case coin.IsETHGroup(r.conf.CoinTypeCode):
		return r.newETHKeygener()
	case r.conf.CoinTypeCode == coin.XRP:
		return r.newXRPKeygener()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", r.conf.CoinTypeCode))
	}
}

func (r *registry) newBTCKeygener() wallets.Keygener {
	return btcwallet.NewBTCKeygen(
		r.newBTC(),
		r.newMySQLClient(),
		r.conf.AddressType,
		r.newSeeder(),
		r.newHdWallter(),
		r.newPrivKeyer(),
		r.newFullPubKeyImporter(),
		r.newMultisiger(),
		r.newAddressExporter(),
		r.newSigner(),
		r.walletType,
	)
}

func (r *registry) newETHKeygener() wallets.Keygener {
	return ethwallet.NewETHKeygen(
		r.newETH(),
		r.newMySQLClient(),
		r.newLogger(),
		r.walletType,
		r.newSeeder(),
		r.newHdWallter(),
		r.newPrivKeyer(),
		r.newAddressExporter(),
		r.newETHSigner(),
	)
}

func (r *registry) newXRPKeygener() wallets.Keygener {
	return xrpwallet.NewXRPKeygen(
		r.newXRP(),
		r.newMySQLClient(),
		r.newLogger(),
		r.walletType,
		r.newSeeder(),
		r.newHdWallter(),
		r.newXRPKeyGenerator(),
		r.newAddressExporter(),
		r.newXRPSigner(),
	)
}

func (r *registry) newSeeder() service.Seeder {
	return commonsrv.NewSeed(
		r.newLogger(),
		r.newSeedRepo(),
		r.walletType,
	)
}

func (r *registry) newHdWallter() service.HDWalleter {
	return commonsrv.NewHDWallet(
		r.newLogger(),
		r.newHdWalletRepo(),
		r.newKeyGenerator(),
		r.conf.CoinTypeCode,
		r.walletType,
	)
}

func (r *registry) newHdWalletRepo() commonsrv.HDWalletRepo {
	return commonsrv.NewAccountHDWalletRepo(
		r.newAccountKeyRepo(),
	)
}

func (r *registry) newPrivKeyer() service.PrivKeyer {
	switch {
	case coin.IsBTCGroup(r.conf.CoinTypeCode):
		return keygensrv.NewPrivKey(
			r.newBTC(),
			r.newLogger(),
			r.newAccountKeyRepo(),
			r.walletType,
		)
	case coin.IsETHGroup(r.conf.CoinTypeCode):
		return ethsrv.NewPrivKey(
			r.newETH(),
			r.newLogger(),
			r.newAccountKeyRepo(),
			r.walletType,
		)
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", r.conf.CoinTypeCode))
	}
}

func (r *registry) newFullPubKeyImporter() service.FullPubKeyImporter {
	return keygensrv.NewFullPubkeyImport(
		r.newBTC(),
		r.newLogger(),
		r.newAuthFullPubKeyRepo(),
		r.newPubkeyFileStorager(),
		r.walletType,
	)
}

func (r *registry) newMultisiger() service.Multisiger {
	return keygensrv.NewMultisig(
		r.newBTC(),
		r.newLogger(),
		r.newAuthFullPubKeyRepo(),
		r.newAccountKeyRepo(),
		r.newMultiAccount(),
		r.walletType,
	)
}

func (r *registry) newAddressExporter() service.AddressExporter {
	return commonsrv.NewAddressExport(
		r.newLogger(),
		r.newAccountKeyRepo(),
		r.newAddressFileStorager(),
		r.newMultiAccount(),
		r.conf.CoinTypeCode,
		r.walletType,
	)
}

func (r *registry) newSigner() service.Signer {
	return coldsrv.NewSign(
		r.newBTC(),
		r.newLogger(),
		r.newAccountKeyRepo(),
		r.newAuthKeyRepo(),
		r.newTxFileStorager(),
		r.newMultiAccount(),
		r.walletType,
	)
}

func (r *registry) newETHSigner() service.Signer {
	return ethsrv.NewSign(
		r.newETH(),
		r.newLogger(),
		r.newTxFileStorager(),
		r.walletType,
	)
}

func (r *registry) newXRPSigner() service.Signer {
	return xrpsrv.NewSign(
		r.newXRP(),
		r.newLogger(),
		r.newXRPAccountKeyRepo(),
		r.newTxFileStorager(),
		r.walletType,
	)
}

func (r *registry) newXRPKeyGenerator() xrpsrv.XRPKeyGenerator {
	return xrpsrv.NewXRPKeyGenerate(
		r.newXRP(),
		r.newLogger(),
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.walletType,
		r.newAccountKeyRepo(),
		r.newXRPAccountKeyRepo(),
	)
}

func (r *registry) newRPCClient() *rpcclient.Client {
	var err error
	if r.rpcClient == nil {
		r.rpcClient, err = btcgrp.NewRPCClient(&r.conf.Bitcoin)
	}
	if err != nil {
		panic(err)
	}
	return r.rpcClient
}

func (r *registry) newBTC() btcgrp.Bitcoiner {
	if r.btc == nil {
		var err error
		r.btc, err = btcgrp.NewBitcoin(
			r.newRPCClient(),
			&r.conf.Bitcoin,
			r.newLogger(),
			r.conf.CoinTypeCode,
		)
		if err != nil {
			panic(err)
		}
	}
	return r.btc
}

func (r *registry) newEthRPCClient() *ethrpc.Client {
	if r.rpcEthClient == nil {
		var err error
		r.rpcEthClient, err = ethgrp.NewRPCClient(&r.conf.Ethereum)
		if err != nil {
			panic(err)
		}
	}
	return r.rpcEthClient
}

func (r *registry) newETH() ethgrp.Ethereumer {
	if r.eth == nil {
		var err error
		r.eth, err = ethgrp.NewEthereum(
			r.newEthRPCClient(),
			&r.conf.Ethereum,
			r.newLogger(),
			r.conf.CoinTypeCode,
		)
		if err != nil {
			panic(err)
		}
	}
	return r.eth
}

func (r *registry) newXRPWSClient() (*ws.WS, *ws.WS) {
	if r.wsXrpPublic == nil {
		var err error
		r.wsXrpPublic, r.wsXrpAdmin, err = xrpgrp.NewWSClient(&r.conf.Ripple)
		if err != nil {
			panic(err)
		}
	}
	return r.wsXrpPublic, r.wsXrpAdmin
}

func (r *registry) newRippleAPI() *xrp.RippleAPI {
	if r.rippleAPI == nil {
		r.rippleAPI = xrp.NewRippleAPI(r.newGRPCConn(), r.newLogger())
	}
	return r.rippleAPI
}

func (r *registry) newGRPCConn() *grpc.ClientConn {
	if r.grpcConn == nil {
		var err error
		r.grpcConn, err = xrpgrp.NewGRPCClient(&r.conf.Ripple.API)
		if err != nil {
			panic(err)
		}
	}
	return r.grpcConn
}

func (r *registry) newXRP() xrpgrp.Rippler {
	if r.xrp == nil {
		var err error
		wsPublic, wsAdmin := r.newXRPWSClient()
		r.xrp, err = xrpgrp.NewRipple(
			wsPublic,
			wsAdmin,
			r.newRippleAPI(),
			&r.conf.Ripple,
			r.newLogger(),
			r.conf.CoinTypeCode,
		)
		if err != nil {
			panic(err)
		}
	}
	return r.xrp
}

func (r *registry) newLogger() *zap.Logger {
	if r.logger == nil {
		r.logger = logger.NewZapLogger(&r.conf.Logger)
	}
	return r.logger
}

func (r *registry) newSeedRepo() coldrepo.SeedRepositorier {
	return coldrepo.NewSeedRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newAccountKeyRepo() coldrepo.AccountKeyRepositorier {
	return coldrepo.NewAccountKeyRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newXRPAccountKeyRepo() coldrepo.XRPAccountKeyRepositorier {
	return coldrepo.NewXRPAccountKeyRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newAuthFullPubKeyRepo() coldrepo.AuthFullPubkeyRepositorier {
	return coldrepo.NewAuthFullPubkeyRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		r.newLogger(),
	)
}

func (r *registry) newAuthKeyRepo() coldrepo.AuthAccountKeyRepositorier {
	return coldrepo.NewAuthAccountKeyRepository(
		r.newMySQLClient(),
		r.conf.CoinTypeCode,
		"",
		r.newLogger(),
	)
}

func (r *registry) newKeyGenerator() key.Generator {
	var chainConf *chaincfg.Params
	switch {
	case coin.IsBTCGroup(r.conf.CoinTypeCode):
		chainConf = r.newBTC().GetChainConf()
	case coin.IsETHGroup(r.conf.CoinTypeCode):
		chainConf = r.newETH().GetChainConf()
	case r.conf.CoinTypeCode == coin.XRP:
		chainConf = r.newXRP().GetChainConf()
	default:
		panic(fmt.Sprintf("coinType[%s] is not implemented yet.", r.conf.CoinTypeCode))
	}

	return key.NewHDKey(
		key.PurposeTypeBIP44,
		r.conf.CoinTypeCode,
		chainConf,
		r.newLogger())
}

func (r *registry) newMySQLClient() *sql.DB {
	if r.mysqlClient == nil {
		dbConn, err := mysql.NewMySQL(&r.conf.MySQL)
		if err != nil {
			panic(err)
		}
		r.mysqlClient = dbConn
	}
	if r.conf.MySQL.Debug {
		boil.DebugMode = true
	}

	return r.mysqlClient
}

func (r *registry) newAddressFileStorager() address.FileRepositorier {
	return address.NewFileRepository(
		r.conf.FilePath.Address,
		r.newLogger(),
	)
}

func (r *registry) newPubkeyFileStorager() address.FileRepositorier {
	return address.NewFileRepository(
		r.conf.FilePath.FullPubKey,
		r.newLogger(),
	)
}

func (r *registry) newTxFileStorager() tx.FileRepositorier {
	return tx.NewFileRepository(
		r.conf.FilePath.Tx,
		r.newLogger(),
	)
}

func (r *registry) newMultiAccount() account.MultisigAccounter {
	if r.multisig == nil {
		if r.accountConf == nil && r.accountConf.Multisigs == nil {
			return account.NewMultisigAccounts(nil)
		}
		r.multisig = account.NewMultisigAccounts(r.accountConf.Multisigs)
	}
	return r.multisig
}
