package api

import (
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
)

// CreateNewAddress アカウント名から新しいアドレスを生成する
// これによって作成されたアカウントはbitcoin core側のwalletで管理される
func (b *Bitcoin) CreateNewAddress(accountName string) (btcutil.Address, error) {
	addr, err := b.client.GetNewAddress(accountName)
	if err != nil {
		return nil, errors.Errorf("GetNewAddress(%s): error: %v", accountName, err)
	}

	return addr, nil
}

// GetAddressesByAccount アカウント名から紐づくすべてのアドレスを取得する
func (b *Bitcoin) GetAddressesByAccount(accountName string) ([]btcutil.Address, error) {
	addrs, err := b.client.GetAddressesByAccount(accountName)
	if err != nil {
		return nil, errors.Errorf("GetAddressesByAccount(%s): error: %v", accountName, err)
	}

	return addrs, nil
}

// ValidateAddress 渡されたアドレスの整合性をチェックする
// TODO: こちらの機能はCayenne側でも必要だが、Cayenneの場合、Bitcoin Coreの機能を単独で使うことは難くはないが、煩雑になってしまう
// TODO: 動作未検証、address_test.goを書いて検証すること
func (b *Bitcoin) ValidateAddress(addr string) error {
	//func (c *Client) ValidateAddress(address btcutil.Address) (*btcjson.ValidateAddressWalletResult, error) {
	address, err := b.decodeAddress(addr)
	if err != nil {
		return err
	}
	_, err = b.client.ValidateAddress(address)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bitcoin) decodeAddress(addr string) (btcutil.Address, error) {
	address, err := btcutil.DecodeAddress(addr, b.chainConf)
	if err != nil {
		return nil, errors.Errorf("btcutil.DecodeAddress() error: %v", err)
	}
	return address, nil
}
