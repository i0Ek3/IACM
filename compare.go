package main

import "fmt"

// comparison algorithm
// original DPoS
func ComparisonDPoS() {
	// statistic round number
	round++

	fmt.Print("\n---------Runing comparison algorithm DPoS...---------\n")
	waitingTime()

	// create common nodes then print them information
	CreateNode()
	waitingTime()

	// select candidate from common node which have more votes
	fmt.Print("\n-------------Select candidate nodes...------------\n")
	fmt.Println("\t\tinfo id votes\n")
	waitingTime()
	SelectCandidate()

	// initial candidate node list
	fmt.Print("\n----------Initializing candidate nodes...----------\n")
	fmt.Println("\tinfo id votes auth d cl cv con bad good \n")
	waitingTime()
	InitCandidate()

	// simulate the vote
	fmt.Print("\n-----------------------Voting--------------------\n")
	waitingTime()
	Vote()

	// selection delegate from candidate
	fmt.Print("\n-----------------Select delegate nodes-----------------\n")
	fmt.Println("\tinfo id votes auth d cl cv con bad good \n")
	waitingTime()
	SelectDelegate(NUMBER)
	fmt.Println("\n")
	//fmt.Println(nodes)

	// initial consensus
	fmt.Print("\n---------------Initializing consensus...---------------\n")
	waitingTime()
	InitialDelegate(NUMBER)

	fmt.Print("\n------------Comparison algorithm DPoS over!---------------------\n")
	waitingTime()
}

// FCSW--an improved DPoS
func ComparisonFCSW() {
	fmt.Print("\n---------Runing comparison algorithm FCSW...---------\n")
	waitingTime()

	FuseMachnism()
	CreditMachnism()
	StandbyWitnessMachnism()

	fmt.Print("\n------------Comparison algorithm FCSW over!---------------------\n")
	waitingTime()
}

