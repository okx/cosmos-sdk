package types

import "math/big"

// Round a decimal with precision, perform bankers rounding (gaussian rounding)
func (d Dec) RoundDecimal(precision int64) Dec {
	precisionMul := NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(precision), nil))
	return newDecFromInt(d.MulInt(precisionMul).RoundInt()).QuoInt(precisionMul)
}
