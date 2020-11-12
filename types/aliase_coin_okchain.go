package types

import (
	"math/big"
)

type Coin = DecCoin
type Coins = DecCoins

type SysCoin = DecCoin
type SysCoins = DecCoins

//var ZeroInt2 = ZeroDec
//var OneInt2 = OneDec
//
//var MinInt2 = MinDec
//var MaxInt2 = MaxDec

func NewCoin(denom string, amount interface{}) DecCoin {
	switch amount := amount.(type) {
	case Int:
		return NewDecCoin(denom, amount)
	case Dec:
		return NewDecCoinFromDec(denom, amount)
	default:
		panic("Invalid amount")
	}
}


func (d Dec) BigInt() *big.Int {
	return d.Int
}

func NewDecCoinsFromDec(denom string, amount Dec) DecCoins {
	return DecCoins{NewDecCoinFromDec(denom, amount)}
}

func (dec DecCoin) ToCoins() Coins {
	return NewCoins(dec)
}

func newDecFromInt(i Int) Dec {
	return newDecFromIntWithPrec(i, 0)
}

func newDecFromIntWithPrec(i Int, prec int64) Dec {
	return Dec{
		new(big.Int).Mul(i.BigInt(), precisionMultiplier(prec)),
	}
}
// Round a decimal with precision, perform bankers rounding (gaussian rounding)
func (d Dec) RoundDecimal(precision int64) Dec {
	precisionMul := NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(precision), nil))
	return newDecFromInt(d.MulInt(precisionMul).RoundInt()).QuoInt(precisionMul)
}


func MustParseCoins(denom, amount string) Coins {
	coins, err := ParseCoins(amount + denom)
	if err != nil {
		panic(err)
	}
	return coins
}


func GetSystemFee() Coin {
	return NewDecCoinFromDec(DefaultBondDenom, NewDecWithPrec(125, 4))
}

func ZeroFee() Coin {
	return NewCoin(DefaultBondDenom, ZeroInt())
}
