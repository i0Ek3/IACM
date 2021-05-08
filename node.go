package main

import "fmt"

// the struct of common node
type Node struct {
	Info  string
	Id    int
	Votes int
	f     string
}

// create nodes then initialize it
func CreateNode() error {
	fmt.Print("\n------------------Initializing common nodes------------------\n")
	fmt.Println("\t\t\tinfo id votes \n")
	waitingTime()
	for i := 0; i < nodeNum; i++ {
		info := fmt.Sprintf("common node:") // node information
		id := i                             // node id number
		vote := rand.Intn(nodeNum)          // the number of vote
		f := ""                             // format control
		nodePool[i] = Node{info, id, vote, f}
		fmt.Println("initializing...", nodePool[i])
	}
	return nil
}

// use lambda to sort nodes
func sortNodes() {
	sort.Slice(candPool, func(i, j int) bool {
		return candPool[i].Votes > candPool[j].Votes
	})
	delePool = candPool[:delegateNum]
}
