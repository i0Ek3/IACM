package utils

import "errors"

const (
	nodeNum      = 101      // sum of nodes
	candidateNum = 30       // candidate number
	delegateNum  = 10       // delegate number
	alternateNum = 30
	NUMBER       = 10
	alpha        = 0.8      // factor of auth
	beta         = 0.2      // factor of vote
	reward       = 0.005
	punish       = -0.02
	lambda       = 0.5
	lambda1      = 0.25
	lambda2      = 0.05
	check        = false
	threshold    = 5        // the length of timer
	mini         = 6        // mininum delegate
	interval     = 10       // the factor of interval
	dimension    = 1        // use in MGM
	epsilon1     = 0.6
	epsilon2     = 0.8
	attention    = 1        // for FCSW use
)

var (
	nodePool = make([]Node, nodeNum)   // common node pool
	candPool = make([]D, candidateNum) // candidate node pool
	delePool = make([]D, nodeNum)      // delegate node pool

	deletePool  = make([]D, delegateNum) // state table: store cl=4's nodes
	freezePool  = make([]D, delegateNum) // state table: store cl=3's nodes
	commonPool  = make([]D, delegateNum) // state table: store cl=2's nodes
	premiumPool = make([]D, delegateNum) // state table: store cl=1's nodes

	round       = 0                  // times of round
	input       string               // the anwser to next loop
	blockchain  []Block              // the blockchain
	curCv       [delegateNum]float64 // current contribution value
	rewardTimes float64              // times of reward
	punishTimes float64              // times of punish
	vote        int

	counter int // maybe we needn't change them size dynamticly
	buffer  int // cause of them size is cleared
	//counter [nodeNum/2]int // for dcml algo use
	//buffer  [nodeNum/2]int // same as last

    alterPool = make([]D, candidateNum)

    showErr = errors.New("somthing wrong here.")
)
