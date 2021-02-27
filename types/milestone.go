package types

import (
	"strconv"
	"sync"
)

var (
	MILESTONE_MERCURY_HEIGHT     string
	MILESTONE_MERCURY_HEIGHT_NUM int64
	once                         sync.Once
)

func string2number(input string) int64 {
	if len(input) == 0 {
		input = "0"
	}
	res, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		panic(err)
	}
	return res
}

func initVersionBlockHeight() {
	once.Do(func() {
		MILESTONE_MERCURY_HEIGHT_NUM = string2number(MILESTONE_MERCURY_HEIGHT)
	})
}

func init() {
	initVersionBlockHeight()
}

//disable transfer tokens to contract address by cli
func IsDisableTransferToContractBlock(height int64) bool {
	return height >= MILESTONE_MERCURY_HEIGHT_NUM
}

//disable change the param EvmDenom by proposal
func IsDisableChangeEvmDenomByProposal(height int64) bool {
	return height >= MILESTONE_MERCURY_HEIGHT_NUM
}

//disable transfer tokens by module of cosmos-sdk/bank
func IsDisableBankTransferBlock(height int64) bool {
	return height >= MILESTONE_MERCURY_HEIGHT_NUM
}
