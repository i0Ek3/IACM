package node

import (
    "fmt"
    "time"
    "math"
    "math/rand"
)

// the struct of candidate node
type Candidate struct {
	Node
	Auth int
	d    int
	fm   string
}

// select candidate from common node
func SelectCandidate() []Node {
	for i := 0; i < nodeNum; i++ {
		for j := 0; j < nodeNum-1; j++ {
			if nodePool[j].Votes < nodePool[j+1].Votes {
				nodePool[j], nodePool[j+1] = nodePool[j+1], nodePool[j]
			}
		}
	}
	for i := 0; i < candidateNum; i++ {
		fmt.Println(nodePool[i])
	}

	// candPool store top 30 common nodes
	//candPool = append(candPool, nodePool[:candidateNum])
	//fmt.Println(nodePool[:candidateNum])
	return nodePool[:candidateNum]
}

// initial candidate node
func InitCandidate() {
	for i := 0; i < candidateNum; i++ {
		//fmt.Printf("candidate %d ", nodes[i].Id)
		info := fmt.Sprintf("candidate")
		id := i
		vote := 0
		f := ""
		auth := 0
		d := 0
		cl := 0
		cv := 0.0
		con := 0
		un := 0
		bad := false
		good := false
		addr := ""
		fmm := ""

		// TODO: for information shows well, we can forbid the struct inheritance
		candPool = append(candPool, D{Node{info, id, vote, f}, auth, d, cl, cv, con, un, bad, good, addr, fmm})
		candPool[i] = D{Node{info, id, vote, f}, auth, d, cl, cv, con, un, bad, good, addr, fmm}

		// use key:value method of struct but shows error here
		//candPool[i] = D{Node{Info:"candidate", Id:i, Votes:0}, Auth:0, d:0, fmm:"\n"}
		fmt.Print("\ninitilizing...", candPool[i])
	}
}

// simulate random authencation
func Auth() {
	for i := 0; i < candidateNum; i++ {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(100000)
		auth := rand.Intn(2) // output 0 or 1 randomly
		candPool[i].Auth = auth
		if candPool[i].Auth == 1 {
			fmt.Printf("candidate %d have authencated.\n", candPool[i].Id)
		} else {
			fmt.Printf("candidate %d have unauthencated.\n", candPool[i].Id)
		}
	}
}

// simulate vote
func Vote() {
	for i := 0; i < candidateNum; i++ {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(100000)

		// TODO: here should simulate voting dynamticly
		// reduce the votes of unauthenticated node
		if candPool[i].Auth == 1 {
			vote = rand.Intn(nodeNum)
			if vote <= 50 {
				vote += 25
			}
		} else {
			vote = rand.Intn(50)
		}

		//vote := rand.Intn(nodeNum * 10) // every node have 10 tickets
		// TODO: vote normalization
		//v := (int)(1 / (1 + math.Exp((float64)(vote))))
		if vote >= 100 {
			vote -= 10
		}
		if vote <= 5 {
			vote += 3
		}
		candPool[i].Votes = vote

		if candPool[i].Auth == 1 {
			fmt.Printf("candidate %d have %d votes and authencated.\n", candPool[i].Id, candPool[i].Votes)
		} else {
			fmt.Printf("candidate %d have %d votes and unauthencated.\n", candPool[i].Id, candPool[i].Votes)
		}
	}
}

// calculate the candidate node support degree, also we call it sd
func CalSD() {
	for i := 0; i < candidateNum; i++ {
		if candPool[i].Auth == 1 {
			candPool[i].d = (int)(alpha*(float64)(candPool[i].Auth) + beta*(float64)(candPool[i].Votes))
			fmt.Printf("candidate %d have %d votes and authencated, the support degree is %d.\n", candPool[i].Id, candPool[i].Votes, candPool[i].d)
		} else {
			candPool[i].d = (int)((1-alpha)*(float64)(candPool[i].Auth) + beta*(float64)(candPool[i].Votes))
			fmt.Printf("candidate %d have %d votes and unauthencated, the support degree is %d.\n", candPool[i].Id, candPool[i].Votes, candPool[i].d)
		}
	}
}
