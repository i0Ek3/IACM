package process

import "fmt"

// call feedback system
func Feedback() {
	fmt.Print("\n----------Call feedback system to reward and punish nodes...----------\n")
	waitingTime()

	// according the nodes' contribution level to put them back into conresponding pool
	nodes := SelectDelegate(NUMBER)
	for i := 0; i < delegateNum; i++ {
		if delePool[i].Cl == 4 {
			delePool[i].isDelete = true
			deletePool = append(deletePool, delePool[i])
			fmt.Printf("candidate %d(delegate %d) is bad node, already taken down and put it back into delete pool!\n", nodes[i].Id, i)
		} else if delePool[i].Cl == 3 {
			delePool[i].isDelete = true
			freezePool = append(freezePool, delePool[i])
			// TODO: after freezen node one round then put it back into common pool to start the new progress
			// that means you should add delay function for algorithm
			fmt.Printf("candidate %d(delegate %d) is abnormal node, record it and put it back into freeze pool!\n", nodes[i].Id, i)
		} else if delePool[i].Cl == 1 {
			delePool[i].isGood = true
			// after record then put it back into common pool
			premiumPool = append(premiumPool, delePool[i])
			commonPool = append(commonPool, delePool[i])
			fmt.Printf("candidate %d(delegate %d) is good node, already got reward and put it back into premium pool.\n", nodes[i].Id, i)
		} else {
			delePool[i].isGood = true
			commonPool = append(commonPool, delePool[i])
			fmt.Printf("candidate %d(delegate %d) is common node, already put it back into common pool.\n", nodes[i].Id, i)
		}
	}
	fmt.Print("\n------------------Reward and punish have done!--------------------\n")
	fmt.Printf("\nGood node we reward it %f contribution value every consensus round!\n", reward)
	fmt.Printf("\nBad node we punish it %f contribution value every consensus round!\n", punish)
	fmt.Print("\n-------------------------------------------------------------------\n")
}
