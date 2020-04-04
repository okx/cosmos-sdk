package types

// votes from the voting convert to the consensus power
// reduce
func VotesToConsensusPower(votes Dec) int64 {
	return votes.QuoInt(PowerReduction).Int64()
}

// convert consensus power to votes
func VotesFromConsensusPower(power int64) Dec {
	return NewDec(power).MulInt(PowerReduction)
}
