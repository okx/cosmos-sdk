package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidDeflation CodeType = 101
)

// ErrInvalidDeflation returns an error when a deflation amount is larger than the current supply in the ledger
func ErrInvalidDeflation(codespace sdk.CodespaceType, deflationAmount, currentSupplyAmount sdk.Dec, tokenSymbol string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDeflation,
		fmt.Sprintf("failed. the deflation %s is larger than the current supply %s",
			fmt.Sprintf("%s%s", deflationAmount, tokenSymbol),
			fmt.Sprintf("%s%s", currentSupplyAmount, tokenSymbol),
		),
	)
}
