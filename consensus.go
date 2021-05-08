package main

import "fmt"

// runing consensus
func Consensus() {
	Process()

	// TODO: add interactive interface
	//fmt.Println("Please input node number always be odd:")
	//fmt.Scanln(&nodeNum)
	//delegateNum = (int)(nodeNum / 3)
	//fmt.Printf("This round we have %d nodes then select %d nodes be a delegate node.", nodeNum, delegateNum)

	gene := Block{time.Now().String(), "0000000000000000000000000000000000000000000000000000000000000000", "", "I'm the genesis block", 1, "0x000"}
	gene.calHash()
	blockchain = append(blockchain, gene)
	//genesisBlock()

	fmt.Print("\n---------------Generating the genesis block...-----------------\n")
	waitingTime()
	// the genesis block cannot validate any more
	fmt.Println("\n")
	fmt.Println(blockchain[0])

	// run 30 round for consensus
	// we use delegateNum*i blocks with i rounds to simulate dpos consensus
	j := 0
	i := 5
	for k := 0; k < delegateNum*i-1; k++ {
		round++
		waitingTime()

		//newBlock := generateBlock(gene, "block content", gene.Address)
		newBlock := generateBlock(blockchain[k], "block content", delePool[j].Address)

		// validate the block then statistic block's coniunity
		// but in real blockchain, one block generated should have validated by six delegate, we just simulate it by one
		//nodes := SelectDelegate()

		if isBlockValid(newBlock, blockchain[k]) {
			if k >= delegateNum {
				k %= 10
			}
			candPool[k].Con++
			blockchain = append(blockchain, newBlock)
		} else {
			candPool[k].Unvalid++
		}

		// broadcast message to other nodes
		Broadcast()

		// every block generate, we update cv and cl
		// FIXME: when we add UpdateCv/Cl() here, something looks so wired, while generated the block 3, rest of them generated are same and with cl/cv
		// but when we noted UpdateCv/Cl(), everything is ok except cv cannot update

		// FIXME: contribution value explode while use UpdateCv() here
		//UpdateCv(candPool[k].Con, candPool[k].Unvalid)

		Upcv(candPool[k].Con)
		UpdateCl()
		fmt.Print("\n-----------Updating contibution value and contribution level...----------\n")
		waitingTime()
		ShowCvCl(k)
		fmt.Print("\n------------Contibution value and contribution level updated!------------\n")

		// print the next block
		fmt.Print("\n------------------------Generating block...------------------------------\n")
		waitingTime()
		fmt.Println("\n")
		fmt.Println(blockchain[k+1])

		// simulate randomly delegate block genreate
		j++
		j = j % len(delePool)

		// FIXME: we want to shows cv and cl information every 10 blocks after generated, but it's so wierd that follows code just shown us chaos here.
		// we cannot validate the delegate's con right now after it generated one block, but we can validate it's valid.
		// so every 10 round we just output the information of cv and cl,
		// while all the process done, we update cv and cl then show again
		if round%10 == 0 {
			// FIXME: chaos here
			//ContributionMechanism()
		}
	}
	//ContributionMechanism()
	check := true
	Check(check)

	// call dcml algorithm
	DCML()
}
