package refund

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"math/big"
)

func NewGasRefundHandler(ak auth.AccountKeeper, sk types.SupplyKeeper) sdk.GasRefundHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, sim bool,
	) (err error) {
		gasRefundHandler := NewGasRefundDecorator(ak, sk)
		return gasRefundHandler(ctx, tx, sim)
	}
}


type GasRefundHandler struct {
	ak           keeper.AccountKeeper
	supplyKeeper types.SupplyKeeper
}

func (grh GasRefundHandler) GasRefundHandle(ctx sdk.Context, tx sdk.Tx, sim bool) (err error) {

	currentGasMeter := ctx.GasMeter()
	TempGasMeter := sdk.NewInfiniteGasMeter()
	ctx = ctx.WithGasMeter(TempGasMeter)

	defer func() {
		ctx = ctx.WithGasMeter(currentGasMeter)
	}()

	gasLimit := currentGasMeter.Limit()
	gasUsed := currentGasMeter.GasConsumed()

	if gasUsed >= gasLimit {
		return nil
	}

	feeTx, ok := tx.(ante.FeeTx)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feePayer := feeTx.FeePayer()
	feePayerAcc := grh.ak.GetAccount(ctx, feePayer)
	if feePayerAcc == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	gas := feeTx.GetGas()
	fees := feeTx.GetFee()
	gasFees := make(sdk.Coins, len(fees))

	for i, fee := range fees {
		gasPrice := new(big.Int).Div(fee.Amount.BigInt(), new(big.Int).SetUint64(gas))
		gasConsumed := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gasUsed))
		gasCost := sdk.NewCoin(fee.Denom, sdk.NewDecFromBigIntWithPrec(gasConsumed, sdk.Precision))
		gasRefund := fee.Sub(gasCost)

		gasFees[i] = gasRefund
	}

	err = RefundFees(grh.supplyKeeper, ctx, feePayerAcc.GetAddress(), gasFees)
	if err != nil {
		return err
	}

	return nil
}

func NewGasRefundDecorator(ak auth.AccountKeeper, sk types.SupplyKeeper) sdk.GasRefundHandler {
	cgrh := GasRefundHandler{
		ak:           ak,
		supplyKeeper: sk,
	}

	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (err error) {
		return cgrh.GasRefundHandle(ctx, tx, simulate)
	}
}


func RefundFees(supplyKeeper types.SupplyKeeper, ctx sdk.Context, acc sdk.AccAddress, refundFees sdk.Coins) error {
	blockTime := ctx.BlockHeader().Time
	feeCollector := supplyKeeper.GetModuleAccount(ctx, types.FeeCollectorName)
	coins := feeCollector.GetCoins()

	if !refundFees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid refund fee amount: %s", refundFees)
	}

	// verify the account has enough funds to pay for fees
	_, hasNeg := coins.SafeSub(refundFees)
	if hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to refund for fees; %s < %s", coins, refundFees)
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	spendableCoins := feeCollector.SpendableCoins(blockTime)
	if _, hasNeg := spendableCoins.SafeSub(refundFees); hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for refund fees; %s < %s", spendableCoins, refundFees)
	}

	err := supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.FeeCollectorName, acc, refundFees)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}

	return nil
}