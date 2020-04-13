package wallets

import (
	"github.com/btcsuite/btcutil"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/address"
	"github.com/hiromaily/go-bitcoin/pkg/keystatus"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/api"
	ctype "github.com/hiromaily/go-bitcoin/pkg/wallets/api/types"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/key"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

// Walleter is for watch only wallet service interface
type Walleter interface {
	ImportPublicKey(fileName string, accountType account.AccountType, isRescan bool) error
	DetectReceivedCoin(adjustmentFee float64) (string, string, error)
	CreateUnsignedPaymentTx(adjustmentFee float64) (string, string, error)
	SendToAccount(from, to account.AccountType, amount btcutil.Amount) (string, string, error)
	SendFromFile(filePath string) (string, error)
	UpdateStatus() error
	Done()
	GetDB() rdb.WalletStorager
	GetBTC() api.Bitcoiner
	GetType() types.WalletType
}

// Coldwalleter may be Not used anywhere
type Coldwalleter interface {
	KeySigner
	KeygenExclusiver
	SignatureExclusiver

	Done()
	GetDB() rdb.ColdStorager
	GetBTC() api.Bitcoiner
	GetType() types.WalletType
}

// common interface for keygen/signature
type KeySigner interface {
	SignatureFromFile(filePath string) (string, bool, string, error)
	GenerateSeed() ([]byte, error)
	GenerateAccountKey(accountType account.AccountType, coinType ctype.CoinType, seed []byte, count uint32) ([]key.WalletKey, error)
	ImportPrivateKey(accountType account.AccountType) error
}

// Keygener is for keygen wallet service interface
type Keygener interface {
	KeySigner

	KeygenExclusiver

	Done()
	GetDB() rdb.ColdStorager
	GetBTC() api.Bitcoiner
	GetType() types.WalletType
}

type KeygenExclusiver interface {
	ExportAccountKey(accountType account.AccountType, keyStatus keystatus.KeyStatus) (string, error)
	ImportMultisigAddrForColdWallet1(fileName string, accountType account.AccountType) error
}

// Signer is for signature wallet service interface
type Signer interface {
	KeySigner

	SignatureExclusiver

	Done()
	GetDB() rdb.ColdStorager
	GetBTC() api.Bitcoiner
	GetType() types.WalletType
}

type SignatureExclusiver interface {
	ImportPublicKeyForColdWallet2(fileName string, accountType account.AccountType) error
	AddMultisigAddressByAuthorization(accountType account.AccountType, addressType address.AddressType) error
	ExportAddedPubkeyHistory(accountType account.AccountType) (string, error)
}