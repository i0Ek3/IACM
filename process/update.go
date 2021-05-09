package process

import (
    "fmt"
    "math"

    log "github.com/sirupsen/logrus"
)

// yet another update contribution value version
func Upcv(delePoolCon int) {
	//nodes := SelectDelegate()
	for i := 0; i < delegateNum; i++ {
		if delePoolCon == 0 {
			delePool[i].Cv += punish // punish cv
		} else if delePoolCon > 0 && delePoolCon < 3 {
			delePool[i].Cv += reward
		} else if delePoolCon >= 3 {
			delePool[i].Cv += lambda * reward
		} else {
			log.Warnf("\nsomething wrong here!\n")
		}
	}
}

// update delegate's cv
// yeah, we messed up this one!
func UpdateCv(con, unvalid int) {
	// the arguments of calculate the contribution value
	sum1 := 0.0
	sum2 := 0.0

	blockHeight := len(blockchain)
	oldBlock := blockchain[blockHeight-1]
	newBlock := generateBlock(oldBlock, "block content", "")

	// first, judge the block whether valid, if it is, then judge the con's number
	// if con is big more and more than threshold, we reward this delegate much more
	// if not, con > 0 and con < threshold, we reward it a little
	// in other words, we needn't validate the block whether it valid or not
	// cause of con, if con > 0 that means the block must be validated and valid
	// else, while con = 0 that means block is unvalid

	// FIXME: messed up here
	for i := 0; i < delegateNum; i++ {
		delta := -(float64)(delePool[i].Con)

		if oldBlock.Height+1 == newBlock.Height || newBlock.Prehash == oldBlock.Hash {
			delePool[i].Con++
			if con > 0 {
				rewardTimes += (float64)(con)
			} else {
				punishTimes++
			}

			// add into blockchain
			blockchain = append(blockchain, newBlock)

			// reward
			if con >= 3 {
				delePool[i].Cv += rewardTimes * reward // reward cv
			}
			sum1 += rewardTimes * reward

			// con = 0 means this delegate didn't product the block
			if con == 0 {
				curCv[i] = delePool[i].Cv + punishTimes*punish
			} else {
				// that means this delegate is good one
				if con >= 3 {
					curCv[i] = delePool[i].Cv + rewardTimes*reward
					//curCv[i] = delePool[i].Cv + lambda1 * 1/(math.Exp(delta * sum1))
				} else {
					//else if con > 0 && con < 3 {
					l := (lambda1 + lambda2) / 2
					curCv[i] = delePool[i].Cv + l*1/(math.Exp(delta*sum1))
				}
			}
		} else {
			delePool[i].Cv += punish // punish cv
			punishTimes++
			sum2 += punishTimes * punish

			if oldBlock.Height+1 != newBlock.Height || newBlock.Prehash != oldBlock.Hash {
				if delePool[i].Unvalid >= 3 {
					curCv[i] = delePool[i].Cv + punishTimes*punish
				}
			} else {
				if con >= 3 {
					curCv[i] = delePool[i].Cv + lambda2*1/(math.Exp(delta*sum2))
				} else {
					curCv[i] = delePool[i].Cv + lambda1*1/(math.Exp(delta*sum2))
				}
			}
		}
	}
}

// update delegate's cl
func UpdateCl() {
	// TODO: use select key word to substitude if loop

	// set cl = 2||3 if the target out of contribution value range
	// for simulate well, we shrink the default contribution value to 0.05
	for i := 0; i < delegateNum; i++ {
		if delePool[i].Cv >= 0.75 && delePool[i].Cv < 1 {
			delePool[i].Cl = 1
		} else if delePool[i].Cv >= 0.5 && delePool[i].Cv < 0.75 {
			delePool[i].Cl = 2
		} else if delePool[i].Cv >= 0.25 && delePool[i].Cv < 0.5 {
			delePool[i].Cl = 3
		} else if delePool[i].Cv >= 0 && delePool[i].Cv < 0.25 {
			delePool[i].Cl = 4
		} else if delePool[i].Cv >= 1 { // out of range, set to 2
			// FIXME: the contribution value out of range
			delePool[i].Cl = 2
		} else if delePool[i].Cv < 0 { // out of range, set to 3
			delePool[i].Cl = 3
		} else {
			continue
		}
	}
}

// display the information of delegates' cv and cl
func ShowCvCl(round int) {
	// FIXME: every delegate after generate block, we update it's conrespoding cl and cv, one by one instead of update all
	for i := 0; i < round; i++ {
		//fmt.Println("\n")
		if delePool[i].Auth == 1 {
			fmt.Printf("candidate %d(delegate %d) have authencated and cl = %d cv = %f.\n", candPool[i].Id, i, delePool[i].Cl, delePool[i].Cv)
		} else {
			fmt.Printf("candidate %d(delegate %d) have unauthencated and cl = %d cv = %f.\n", candPool[i].Id, i, delePool[i].Cl, delePool[i].Cv)
		}
	}
}
