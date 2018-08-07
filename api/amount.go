package api

import (
	"github.com/btcsuite/btcutil"
)

// ConvertToSatoshi bitcoinをSatoshiに変換する
func ConvertToSatoshi(f float64) (btcutil.Amount, error) {
	// Amount
	// Satoshiに変換しないといけない
	// 1Satoshi＝0.00000001BTC
	return btcutil.NewAmount(f)
}

// ListAccounts これは単純にアカウントの資産一覧が表示されるだけ
func (b *Bitcoin) ListAccounts() (map[string]btcutil.Amount, error) {
	return b.Client.ListAccounts()
}

// GetBalanceByAccount アカウントに対してのBalanceを取得する
func (b *Bitcoin) GetBalanceByAccount(name string) (btcutil.Amount, error) {
	return b.Client.GetBalance("name")
}
