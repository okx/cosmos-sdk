package types

import (
	"fmt"
	"encoding/json"


	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto"
)

type BaseAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.Coins      `json:"coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber uint64         `json:"account_number"`
	Sequence      uint64         `json:"sequence"`
}

type DecAccount struct {
	Address       sdk.AccAddress `json:"address"`
	Coins         sdk.DecCoins   `json:"coins"`
	PubKey        crypto.PubKey  `json:"public_key"`
	AccountNumber uint64         `json:"account_number"`
	Sequence      uint64         `json:"sequence"`
}

// String implements fmt.Stringer
func (acc DecAccount) String() string {
	var pubkey string

	if acc.PubKey != nil {
		pubkey = sdk.MustBech32ifyAccPub(acc.PubKey)
	}

	return fmt.Sprintf(`Account:
 Address:       %s
 Pubkey:        %s
 Coins:         %v
 AccountNumber: %d
 Sequence:      %d`,
		acc.Address, pubkey, acc.Coins, acc.AccountNumber, acc.Sequence,
	)
}

func ValidOriginalSymbol(name string) bool {
	return false
}

func ValidSymbol(name string) bool {
	return false
}


func BaseAccountToDecAccount(account auth.BaseAccount) DecAccount {
	var decCoins sdk.DecCoins
	for _, coin := range account.Coins {
		dec := sdk.NewDecFromBigIntWithPrec(coin.Amount.BigInt(), sdk.Precision)
		decCoin := sdk.NewDecCoinFromDec(coin.Denom, dec)
		decCoins = append(decCoins, decCoin)
	}
	decAccount := DecAccount{
		Address:       account.Address,
		PubKey:        account.PubKey,
		Coins:         decCoins,
		AccountNumber: account.AccountNumber,
		Sequence:      account.Sequence,
	}
	return decAccount
}


type Currency struct {
	Description string  `json:"description"`
	Symbol      string  `json:"symbol"`
	TotalSupply sdk.Dec `json:"total_supply"`
}

func (currency Currency) String() string {
	b, err := json.Marshal(currency)
	if err != nil {
		return "[{}]"
	}
	return string(b)
}

//type ByDenom sdk.Coins
//
//func (d ByDenom) Len() int           { return len(d) }
//func (d ByDenom) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
//func (d ByDenom) Less(i, j int) bool { return d[i].Denom < d[j].Denom }

type Transfer struct {
	To     string `json:"to"`
	Amount string `json:"amount"`
}

type TransferUnit struct {
	To    sdk.AccAddress `json:"to"`
	Coins sdk.Coins      `json:"coins"`
}

type CoinsInfo []CoinInfo

func (d CoinsInfo) Len() int           { return len(d) }
func (d CoinsInfo) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d CoinsInfo) Less(i, j int) bool { return d[i].Symbol < d[j].Symbol }

type AccountResponse struct {
	Address    string    `json:"address"`
	Currencies CoinsInfo `json:"currencies"`
}

func NewAccountResponse(addr string) AccountResponse {
	var accountResponse AccountResponse
	accountResponse.Address = addr
	accountResponse.Currencies = []CoinInfo{}
	return accountResponse
}

type CoinInfo struct {
	Symbol    string `json:"symbol" v2:"currency"`
	Available string `json:"available" v2:"available"`
	Freeze    string `json:"freeze" v2:"freeze"`
	Locked    string `json:"locked" v2:"locked"`
}

func NewCoinInfo(symbol, available, freeze, locked string) *CoinInfo {
	return &CoinInfo{
		Symbol:    symbol,
		Available: available,
		Freeze:    freeze,
		Locked:    locked,
	}
}

//type QueryPage struct {
//	Page    int `json:"page"`
//	PerPage int `json:"per_page"`
//}

type AccountParam struct {
	Symbol string `json:"symbol"`
	Show   string `json:"show"`
	//QueryPage
}

type AccountParamV2 struct {
	Currency string `json:"currency"`
	HideZero string `json:"hide_zero"`
}

type AccCoins struct {
	Acc   sdk.AccAddress `json:"address"`
	Coins sdk.Coins      `json:"coins"`
}

