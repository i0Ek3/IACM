package node

import (
    "fmt"
    "time"
    "math/rand"
)

// the struct of delegate node
type Delegate struct {
	Candidate
	Cl       int
	Cv       float64
	Con      int
	isDelete bool
	isGood   bool
	fmm      string
}

// select delegate node from candidate
func SelectDelegate(number int) []D {
	n := candPool
	for i := 0; i < len(n); i++ {
		// check node wheather isDelete from round 2
		//if round > 1 && candPool[i].isDelete {
		// TODO: delete this node
		//}
		for j := 0; j < len(n)-1; j++ {
			if n[j].d < n[j+1].d {
				n[j], n[j+1] = n[j+1], n[j]
			}
		}
	}

	for i := 0; i < number; i++ {
		fmt.Println(candPool[i])
	}

	//delePool = append(delePool, candPool)
	return n[:number]
}

// initialize consensus: initialize delegate node
func InitialDelegate(number int) {
	//nodes := SelectDelegate()
	for i := 0; i < number; i++ {
		delePool[i].Cl = 2
		delePool[i].Cv = 0.05
		delePool[i].Con = 0
		delePool[i].Unvalid = 0
		delePool[i].isDelete = false
		delePool[i].isGood = false

		/*
		   if candPool[i].Auth == 1 {
		       fmt.Printf("candidate %d(delegate %d) have authencated, and cl = %d cv = %f con = %d\n", nodes[i].Id, i, delePool[i].Cl, delePool[i].Cv, delePool[i].Con)
		   } else {
		       fmt.Printf("candidate %d(delegate %d) have unauthencated, and cl = %d cv = %f con = %d\n", nodes[i].Id, i, delePool[i].Cl, delePool[i].Cv, delePool[i].Con)
		   }
		*/

		fmt.Printf("candidate %d ", candPool[i].Id)
		// use struct to initial delegate
		info := fmt.Sprintf("delegate")
		id := i
		vote := 0
		f := ""
		auth := candPool[i].Auth
		d := candPool[i].d
		cl := 2
		cv := 0.5
		con := 0
		un := 0
		bad := false
		good := false
		addr := "0x00" + strconv.Itoa(i+1)
		fmm := ""

		delePool = append(delePool, D{Node{info, id, vote, f}, auth, d, cl, cv, con, un, bad, good, addr, fmm})
		delePool[i] = D{Node{info, id, vote, f}, auth, d, cl, cv, con, un, bad, good, addr, fmm}
		fmt.Println("initialized to delegate", i, delePool[i])
	}
}

// random shuffle delegate node
func Shuffle() {
	rand.Seed(time.Now().Unix())
	idx1 := rand.Intn(delegateNum)
	idx2 := rand.Intn(delegateNum)

	if idx1 == idx2 {
	Label:
		idx2 = rand.Intn(delegateNum)
		if idx1 == idx2 {
			goto Label
		}
	}

	SelectDelegate(NUMBER)
	tmp := delePool[idx1]
	delePool[idx1] = delePool[idx2]
	delePool[idx2] = tmp
}
