package utils

import (
    "fmt"
    "time"

    log "github.com/sirupsen/logrus"
)

// time simulation
func waitingTime() {
	if candidateNum <= (int)(1/3*nodeNum) {
		time.Sleep(10 * time.Second)
	} else {
		time.Sleep(6 * time.Second)
	}
}

// check nodes' attribution in a new round
func CheckAttr() {
	fmt.Print("-------------------Checking nodes' attribution...--------------------\n")
	waitingTime()

	// TODO: think about develop dcml algorithm from here or start the new selection
	// if one node's cv = 4, we don't use it again, or it's cv = 3, we will freezen it one round
	// and others, back to work

	for i := 0; i < nodeNum; i++ {
		if delePool[i].isDelete && delePool[i].Cl == 4 {
			deletePool = append(deletePool, delePool[i])
		} else if delePool[i].isDelete && delePool[i].Cl == 3 {
			freezePool = append(freezePool, delePool[i])
			// TODO: here should use thread to run follows code!
			// time.Sleep(3 * time.Second)
			// commonPool = append(commonPool, nodePool[i])
		} else if delePool[i].isGood && delePool[i].Cl == 1 {
			premiumPool = append(premiumPool, delePool[i])
		} else {
			commonPool = append(commonPool, delePool[i])
		}
	}
	fmt.Print("--------------------Nodes' attribution checked!---------------------\n")
}

func Debug() {
	ShowErr()
	log.Debugf("---> marked! <---")
}

// ShowErr shows the error msg
func ShowErr() error {
	return showErr
}
