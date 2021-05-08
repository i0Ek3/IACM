package main

// fuse machnism
func FuseMachnism() float64 {
	var thres float64
	var fuseFactor float64
	totalVotes := rand.Intn(nodeNum)

	// fusing machnism
	if attention == 1 {
		// pay more attention to against vote
		fuseFactor = 0.1
		thres = fuseFactor * (float64)(totalVotes)
	} else {
		fuseFactor = 0.5
		thres = fuseFactor * (float64)(totalVotes)
	}
	return thres
}

// credit machnism
func CreditMachnism() float64 {
	var result float64
	supportVotes := rand.Intn(nodeNum)
	againstVotes := rand.Intn(candidateNum)
	result = alpha*(float64)(supportVotes) - beta*(float64)(againstVotes)
	return result
}

// standby witness machnism
// select n+m nodes to be delegate, n to consensus then alternate m nodes
func StandbyWitnessMachnism() {
	m := 5
	SelectDelegate(NUMBER + m)
	InitialDelegate(NUMBER + m)
	Consensus()
	// TODO: alternate these m nodes first
	TimingAlternate(candidateNum)

}
