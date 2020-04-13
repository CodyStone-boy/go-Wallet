package coldwallet

//Cold wallet

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hiromaily/go-bitcoin/pkg/account"
	"github.com/hiromaily/go-bitcoin/pkg/model/rdb/coldrepo"
	ctype "github.com/hiromaily/go-bitcoin/pkg/wallets/api/types"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/key"
	"github.com/hiromaily/go-bitcoin/pkg/wallets/types"
)

//1.Seedの生成+DBに登録
//2.Multisig Keyの生成+DBに登録(承認用は端末を分けて管理しないと意味がないかも)

//CreateMultiSig(addmultisigaddress)にwalletにmultisig用のprivate keyを登録する
//これのパラメータには、multisigしないと送金許可しないアドレス(receipt, payment)+承認用のアドレスをセット
//これによって、生成されたアドレスから送金する場合、パラメータにセットしたアドレスに紐づく秘密鍵が必要
//payment,receiptのアドレスは、実際には、addmultisigaddressによって生成されたアドレスに置き換えられる。

//3.Client Keyの生成+DBに登録
//4.Receipt Keyの生成 + Multisig対応 + DBに登録 (1日1Key消費するイメージ)
//5.Payment Keyの生成+ Multisig + DBに登録 (1日1Key消費するイメージ)

// GenerateSeed seedを生成する
func (w *ColdWallet) GenerateSeed() ([]byte, error) {

	bSeed, err := w.retrieveSeed()
	if err == nil {
		return bSeed, nil
	}

	// seed生成
	// set default seed
	var strSeed string

	//TODO: envは削除するので、一旦本番モードで実装
	bSeed, err = key.GenerateSeed()
	if err != nil {
		return nil, errors.Errorf("key.GenerateSeed() error: %s", err)
	}
	strSeed = key.SeedToString(bSeed)
	//if w.Env == ctype.EnvDev && w.Seed != "" {
	//	strSeed = w.Seed
	//	bSeed, err = key.SeedToByte(strSeed)
	//	if err != nil {
	//		return nil, errors.Errorf("key.SeedToByte() error: %s", err)
	//	}
	//} else {
	//	bSeed, err = key.GenerateSeed()
	//	if err != nil {
	//		return nil, errors.Errorf("key.GenerateSeed() error: %s", err)
	//	}
	//	strSeed = key.SeedToString(bSeed)
	//}

	// DBにseed情報を登録
	_, err = w.storager.InsertSeed(strSeed, nil, true)
	if err != nil {
		return nil, errors.Errorf("DB.InsertSeed() error: %s", err)
	}

	return bSeed, nil
}

func (w *ColdWallet) retrieveSeed() ([]byte, error) {
	// DBからseed情報を登録
	seed, err := w.storager.GetSeedOne()
	if err == nil && seed.Seed != "" {
		w.logger.Info("seed have already been generated")
		return key.SeedToByte(seed.Seed)
	}

	return nil, errors.Errorf("DB.GetSeedOne() error: %s", err)
}

// GenerateAccountKey AccountType属性のアカウントKeyを生成する
// TODO:AccountTypeAuthorizationのときは、レコードがある場合は追加できないようにしたほうがいい？？
func (w *ColdWallet) GenerateAccountKey(accountType account.AccountType, coinType ctype.CoinType, seed []byte, count uint32) ([]key.WalletKey, error) {
	if w.wtype == types.WalletTypeWatchOnly {
		return nil, errors.New("it's available on Coldwallet1, Coldwallet2")
	}

	//現在のindexを取得
	idx, err := w.storager.GetMaxIndexOnAccountKeyTable(accountType)
	if err != nil {
		idx = 0
	} else {
		idx++
	}
	w.logger.Info(
		"call storager.GetMaxIndexOnAccountKeyTable() current index",
		zap.Int64("idx", idx))

	return w.generateAccountKey(accountType, coinType, seed, uint32(idx), count)
}

// generateKey AccountType属性のアカウントKeyを生成する
func (w *ColdWallet) generateAccountKey(accountType account.AccountType, coinType ctype.CoinType, seed []byte, idxFrom, count uint32) ([]key.WalletKey, error) {
	// HDウォレットのkeyを生成する
	walletKeys, err := w.generateAccountKeyData(accountType, coinType, seed, idxFrom, count)
	if err != nil {
		return nil, errors.Errorf("key.generateAccountKeyData(AccountTypeClient) error: %s", err)
	}

	// Account
	//var account string
	//if accountType != ctype.AccountTypeClient {
	//	account = string(accountType)
	//}
	account := string(accountType)

	// DBにClientAccountのKey情報を登録
	accountKeyClients := make([]coldrepo.AccountKeyTable, len(walletKeys))
	for idx, key := range walletKeys {
		accountKeyClients[idx] = coldrepo.AccountKeyTable{
			WalletAddress:         key.Address,
			P2shSegwitAddress:     key.P2shSegwit,
			FullPublicKey:         key.FullPubKey,
			WalletMultisigAddress: "",
			RedeemScript:          key.RedeemScript,
			WalletImportFormat:    key.WIF,
			Account:               account,
			Idx:                   idxFrom,
		}
		idxFrom++
	}
	err = w.storager.InsertAccountKeyTable(accountType, accountKeyClients, nil, true)
	if err != nil {
		return nil, errors.Errorf("DB.InsertAccountKeyTable() error: %s", err)
	}

	return walletKeys, err
}

// generateKeyData AccountType属性のアカウントKeyを生成する
func (w *ColdWallet) generateAccountKeyData(accountType account.AccountType, coinType ctype.CoinType, seed []byte, idxFrom, count uint32) ([]key.WalletKey, error) {
	// Keyオブジェクト
	keyData := key.NewKey(coinType, w.btc.GetChainConf(), w.logger)

	// key生成
	priv, _, err := keyData.CreateAccount(seed, accountType)
	if err != nil {
		return nil, errors.Errorf("key.CreateAccount() error: %s", err)
	}

	walletKeys, err := keyData.CreateKeysWithIndex(priv, idxFrom, count)
	if err != nil {
		return nil, errors.Errorf("key.CreateKeysWithIndex() error: %s", err)
	}

	return walletKeys, nil
}