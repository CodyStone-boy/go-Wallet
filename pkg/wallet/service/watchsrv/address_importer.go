package watchsrv

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	models "github.com/hiromaily/go-bitcoin/pkg/models/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/repository/watchrepo"
	"github.com/hiromaily/go-bitcoin/pkg/wallet"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/api"
	"github.com/hiromaily/go-bitcoin/pkg/wallet/coin"
)

// AddressImporter is AddressImporter interface
type AddressImporter interface {
	ImportAddress(fileName string, accountType account.AccountType, isRescan bool) error
}

// AddressImport type
type AddressImport struct {
	btc          api.Bitcoiner
	logger       *zap.Logger
	dbConn       *sql.DB
	addrRepo     watchrepo.AddressRepositorier
	addrFileRepo address.FileRepositorier
	coinTypeCode coin.CoinTypeCode
	wtype        wallet.WalletType
}

// NewAddressImport returns AddressImport object
func NewAddressImport(
	btc api.Bitcoiner,
	logger *zap.Logger,
	dbConn *sql.DB,
	addrRepo watchrepo.AddressRepositorier,
	addrFileRepo address.FileRepositorier,
	coinTypeCode coin.CoinTypeCode,
	wtype wallet.WalletType) *AddressImport {

	return &AddressImport{
		btc:          btc,
		logger:       logger,
		dbConn:       dbConn,
		addrRepo:     addrRepo,
		addrFileRepo: addrFileRepo,
		coinTypeCode: coinTypeCode,
		wtype:        wtype,
	}
}

// ImportAddress import PubKey from csv filecsv into database,
//  - if account is client, which doesn't have account ??
func (a *AddressImport) ImportAddress(fileName string, accountType account.AccountType, isRescan bool) error {
	// read file for public key
	pubKeys, err := a.addrFileRepo.ImportAddress(fileName)
	if err != nil {
		return errors.Wrap(err, "fail to call key.ImportPubKey()")
	}

	pubKeyData := make([]*models.Address, 0, len(pubKeys))
	for _, key := range pubKeys {
		inner := strings.Split(key, ",")
		//kind of required address is different according to account
		var addr string
		if accountType == account.AccountTypeClient {
			switch a.btc.CoinTypeCode() {
			case coin.BTC:
				addr = inner[2] //p2sh_segwit_address
			case coin.BCH:
				addr = inner[1] //p2pkh_address
			default:
				return errors.Errorf("coinTypeCode is out of range: %s", a.btc.CoinTypeCode().String())
			}
		} else {
			addr = inner[4] //multisig_address
		}

		//call bitcoin API `importaddress` with account(label)
		//Note: Error would occur when using only 1 bitcoin core server under development
		// because address is already imported
		//isRescan would be `false` usually
		err := a.btc.ImportAddressWithLabel(addr, accountType.String(), isRescan)
		if err != nil {
			//-4: The wallet already contains the private key for this address or script
			a.logger.Warn(
				"fail to call btc.ImportAddressWithLabel() but continue following addresses",
				zap.String("address", addr),
				zap.String("account_type", accountType.String()),
				zap.Error(err))
			continue
		}

		pubKeyData = append(pubKeyData, &models.Address{
			Coin:          a.coinTypeCode.String(),
			Account:       accountType.String(),
			WalletAddress: addr,
		})

		//confirm pubkey is added as watch only wallet
		a.checkImportedPubKey(addr)
	}

	//insert imported pubKey
	err = a.addrRepo.InsertBulk(pubKeyData)
	if err != nil {
		return errors.Wrap(err, "fail to call repo.Pubkey().InsertBulk()")
		//TODO:What if this inserting is failed, how it can be recovered to keep consistancy
		// pubkey is added in wallet, but database doesn't have records
		// try to run this func again
	}

	return nil
}

// checkImportedPubKey confirm pubkey is added as watch only wallet
func (a *AddressImport) checkImportedPubKey(addr string) {
	addrInfo, err := a.btc.GetAddressInfo(addr)
	if err != nil {
		a.logger.Error(
			"fail to call btc.GetAddressInfo()",
			zap.String("address", addr),
			zap.Error(err))
		return
	}
	a.logger.Debug("account is found",
		zap.String("account", addrInfo.GetLabelName()),
		zap.String("address", addr))

	//`watch only wallet` is expected
	//TODO: if wallet,keygen,sign is working on only one bitcoin core server,
	// result would be `iswatchonly=false`
	if !addrInfo.Iswatchonly {
		a.logger.Warn("this address must be watch only wallet")
		//grok.Value(addrInfo)
	}
}
