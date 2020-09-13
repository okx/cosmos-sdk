package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	KeyMintDenom     = []byte("MintDenom")
	KeyDeflationRate = []byte("DeflationRate")
	//KeyInflationMax  = []byte("InflationMax")
	//KeyInflationMin  = []byte("InflationMin")
	//KeyGoalBonded    = []byte("GoalBonded")
	KeyBlocksPerYear      = []byte("BlocksPerYear")
	KeyDeflationYears     = []byte("DeflationYears")
	KeyInitTokensPerBlock = []byte("InitTokensPerBlock")
)

// mint parameters
type Params struct {
	MintDenom     string  `json:"mint_denom" yaml:"mint_denom"`         // type of coin to mint
	DeflationRate sdk.Dec `json:"inflation_rate" yaml:"inflation_rate"` // maximum annual change in inflation rate
	//InflationMax  sdk.Dec `json:"inflation_max" yaml:"inflation_max"`     // maximum inflation rate
	//InflationMin  sdk.Dec `json:"inflation_min" yaml:"inflation_min"`     // minimum inflation rate
	//GoalBonded    sdk.Dec `json:"goal_bonded" yaml:"goal_bonded"`         // goal of percent bonded atoms
	BlocksPerYear      uint64  `json:"blocks_per_year" yaml:"blocks_per_year"` // expected blocks per year
	DeflationYears     uint64  `json:"inflation_years" yaml:"inflation_years"`
	InitTokensPerBlock sdk.Dec `json:"init_tokens_per_block" yaml:"init_tokens_per_block"`
}

// ParamTable for minting module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(mintDenom string, deflationRateChange, inflationMax,
	inflationMin, goalBonded sdk.Dec, blocksPerYear, deflationYears uint64, initTokensPerBlock sdk.Dec) Params {

	return Params{
		MintDenom:     mintDenom,
		DeflationRate: deflationRateChange,
		//InflationMax:  inflationMax,
		//InflationMin:  inflationMin,
		//GoalBonded:    goalBonded,
		BlocksPerYear:      blocksPerYear,
		DeflationYears:     deflationYears,
		InitTokensPerBlock: initTokensPerBlock,
	}
}

// default minting module parameters
func DefaultParams() Params {
	return Params{
		MintDenom:     sdk.DefaultBondDenom,
		DeflationRate: sdk.NewDecWithPrec(50, 2),
		//InflationMax:  sdk.NewDecWithPrec(20, 2),
		//InflationMin:  sdk.NewDecWithPrec(7, 2),
		//GoalBonded:    sdk.NewDecWithPrec(67, 2),
		BlocksPerYear:      uint64(60 * 60 * 8766 / 3), // assuming 5 second block times
		DeflationYears:     4,
		InitTokensPerBlock: sdk.NewDec(50),
	}
}

// validate params
func ValidateParams(params Params) error {
	//if params.GoalBonded.LT(sdk.ZeroDec()) {
	//	return fmt.Errorf("mint parameter GoalBonded should be positive, is %s ", params.GoalBonded.String())
	//}
	//if params.GoalBonded.GT(sdk.OneDec()) {
	//	return fmt.Errorf("mint parameter GoalBonded must be <= 1, is %s", params.GoalBonded.String())
	//}
	//if params.InflationMax.LT(params.InflationMin) {
	//	return fmt.Errorf("mint parameter Max inflation must be greater than or equal to min inflation")
	//}
	if params.MintDenom == "" {
		return fmt.Errorf("mint parameter MintDenom can't be an empty string")
	}
	return nil
}

func (p Params) String() string {
	return fmt.Sprintf(`Minting Params:
  Mint Denom:             %s
  Deflation Rate:         %s
  Blocks Per Year:        %d
  Deflation Years:        %d
  Init Tokens Per Block:  %s
`,
		p.MintDenom, p.DeflationRate, p.BlocksPerYear, p.DeflationYears, p.InitTokensPerBlock,
	)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyMintDenom, &p.MintDenom},
		{KeyDeflationRate, &p.DeflationRate},
		//{KeyInflationMax, &p.InflationMax},
		//{KeyInflationMin, &p.InflationMin},
		//{KeyGoalBonded, &p.GoalBonded},
		{KeyBlocksPerYear, &p.BlocksPerYear},
		{KeyDeflationYears, &p.DeflationYears},
		{KeyInitTokensPerBlock, &p.InitTokensPerBlock},
	}
}
