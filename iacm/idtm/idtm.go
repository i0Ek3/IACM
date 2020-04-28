// IACM-IDTM
//      Name: idtm.go
//      Author: mcy
//      An improved DPoS consensus implementation. 
//

package main // this one can be execute by command: go run 
//package idtm // but this one only can use command: go build

import (
    "crypto/sha256"
    "encoding/hex"
    "math/rand"
    "strconv"
    "math"
    "time"
    "fmt"
)

const (
    nodeNum         = 101   // sum of nodes
    candidateNum    = 30    // candidate number
    delegateNum     = 10    // delegate number
    alpha           = 0.8   // factor of auth   
    beta            = 0.2   // factor of vote
)

// the struct of common node
type Node struct {
    Info        string
    Id          int
    Votes       int
    f           string
}

// the struct of candidate node
type Candidate struct {
    Node
    Auth        int 
    d           int
    fm          string
}

// the struct of delegate node
type Delegate struct {
    Candidate
    Cl          int
    Cv          float64
    Con         int
    isDelete    bool
    isGood      bool
    fmm         string
}

// the struct of Block
type Block struct {
    Height      int
    Timestamp   string
    Hash        string
    Prehash     string
    Data        []byte
    Address     string
    delegate    *Node
}

var (
    nodePool    = make([]Node, nodeNum)     // common node pool
    candPool    = make([]Candidate, candidateNum)// candidate node pool
    delePool    = make([]Delegate, delegateNum) // delegate node pool
    
    deletePool  = make([]Delegate, delegateNum) // state table: store cl=4's nodes
    freezePool  = make([]Delegate, delegateNum) // state table: store cl=3's nodes
    commonPool  = make([]Delegate, delegateNum) // state table: store cl=2's nodes
    premiumPool = make([]Delegate, delegateNum) // state table: store cl=1's nodes
    
    round       = 0          // times of round 
    input       string       // the anwser to next loop
    blockchain  []Block      // the blockchain
    curCv       [delegateNum]float64  // current contribution value
    rewardTimes float64      // times of reward
    punishTimes float64      // times of punish
    vote        int
)

// first block
func geneBlock() Block {
    gene := Block{0, time.Now().String(), "", "", []byte("I'm the first block"), "", nil}
    blockchain = append(blockchain, gene)
    gene.Hash = string(calHash(gene))
    return gene
}

// delegate nodes consensus
// generate the block
func (node *Node) GenerateNewBlock(lastBlock Block, data []byte, addr string) Block {
    time.Sleep(3 * time.Second) // for easy use, every 3sec generate a block
    var newBlock = Block{lastBlock.Height+1, time.Now().String(), lastBlock.Hash, "", data, addr, nil}
    newBlock.Hash = hex.EncodeToString(calHash(newBlock))
    newBlock.delegate = node
    return newBlock
}

// calculate the block hash
func calHash(block Block) []byte {
    hash := strconv.Itoa(block.Height) + block.Timestamp + block.Prehash + hex.EncodeToString(block.Data) + block.Address
    h := sha256.New()
    h.Write([]byte(hash))
    hashed := h.Sum(nil)
    return hashed
}

// create nodes then initialize it
func CreateNode() {
    for i := 0; i < nodeNum; i++ {
        info := fmt.Sprintf("common node:") // node information
        id   := i       // node id number
        vote := rand.Intn(nodeNum)       // the number of vote
        //auth := 0       // authencation
        //d    := 0       // support degree
        //cl   := 0       // contribution level
        //cv   := 0.0     // contribution value
        //con  := 0       // the contious of generate block
        //bad  := false   // bad node?
        //good := false   // good node?
        f      := "\n"    // format control
        nodePool[i] = Node{info, id, vote, f}
    }
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

    // candPool store top 30 common nodes
    //candPool = append(candPool, nodePool[:candidateNum])
    return nodePool[:candidateNum]
}

// initial candidate node
func initCandidate() {
    for i := 0; i < candidateNum; i++ {
        info := fmt.Sprintf("candidate")
        id   := i
        vote := 0
        auth := 0
        d    := 0
        f    := ""
        fm   := "\n"

        // TODO: for information shows well, we can forbid the struct inheritance 
        candPool = append(candPool, Candidate{Node{"candidate", i, 0, f}, 0, 0, fm})
        candPool[i] = Candidate{Node{info, id, vote, f}, auth, d, fm}
        fmt.Printf("initilizing... %v\n",  candPool[i])
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
        } else {
            vote = rand.Intn(50) 
        }

        //vote := rand.Intn(nodeNum * 10) // every node have 10 tickets
        // TODO: vote normalization
        //v := (int)(1 / (1 + math.Exp((float64)(vote))))
        if vote >= 130 {
            vote -= 10
        } 
        if vote <= 10 {
            vote += 10
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
            candPool[i].d = (int)(alpha * (float64)(candPool[i].Auth) + beta * (float64)(candPool[i].Votes))
            fmt.Printf("candidate %d have %d votes and authencated, the support degree is %d.\n", candPool[i].Id, candPool[i].Votes, candPool[i].d)
        } else {
            candPool[i].d = (int)((1-alpha) * (float64)(candPool[i].Auth) + beta * (float64)(candPool[i].Votes))
            fmt.Printf("candidate %d have %d votes and unauthencated, the support degree is %d.\n", candPool[i].Id, candPool[i].Votes, candPool[i].d)
        }
    }
}

// select delegate node from candidate
func SelectDelegate() []Candidate {
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
    //delePool = append(delePool, candPool)
    return n[:delegateNum]
}

// initialize consensus: initialize delegate node
func InitialDelegate() {
    nodes := SelectDelegate()
    for i := 0; i < delegateNum; i++ {
        delePool[i].Cl       = 2
        delePool[i].Cv       = 0.5
        delePool[i].Con      = 0
        delePool[i].isDelete = false
        delePool[i].isGood   = false
        
        // FIXME: cannot read the votes from delegate node
        if candPool[i].Auth == 1 {
            fmt.Printf("candidate %d(delegate %d) have authencated, and cl = %d cv = %f con = %d.\n", nodes[i].Id, i, delePool[i].Cl, delePool[i].Cv, delePool[i].Con)
        } else {
            fmt.Printf("candidate %d(delegate %d) have unauthencated, and cl = %d cv = %f con = %d.\n", nodes[i].Id, i, delePool[i].Cl, delePool[i].Cv, delePool[i].Con)
        }
    }
}

// time simulation
func waitingTime() {
    if candidateNum <= (int)(1/3 * nodeNum) {
        time.Sleep(3 * time.Second)
    } else {
        time.Sleep(5 * time.Second)
    }
}

// validate the block generated by delegate
func isBlockValid(newBlock, oldBlock Block) bool{
    //fmt.Println("\n------------------Validating the block...--------------------\n")
    waitingTime()
    
    // TODO: validate delegateNum delegate nodes
    for i := 0; i < delegateNum; i++ {
        //
    }
    
    //if oldBlock.Height + 1 != newBlock.Height{
   	//    //fmt.Println("\nValidation failed! Wrong Height!\n")
   	//    return false
    //}
    //if newBlock.Prehash != oldBlock.Hash{
   	//    //fmt.Println("\nValidation failed! Wrong Prehash!\n")
   	//    return false
    //}
    //fmt.Println("\n-------------------Validation Successful!---------------------\n")
    return true
}

// random shuffle delegate node
func Shuffle() {
	rand.Seed(time.Now().Unix())
	idx1 := rand.Intn(delegateNum)
	idx2 := rand.Intn(delegateNum)

	if idx1 == idx2 {
	Label :
		idx2 = rand.Intn(delegateNum)
		if idx1 == idx2 {
			goto Label
		}
	}

    SelectDelegate()
	tmp := delePool[idx1]
	delePool[idx1] = delePool[idx2]
	delePool[idx2] = tmp
}

// process
func Process() {
    // statistic round number
    round++
    
    // create common nodes then print them information
    CreateNode()
    fmt.Print("\n----------------Initializing nodes---------------\n")
    fmt.Println("\t\tinfo id votes \n")
    waitingTime()
    fmt.Println(nodePool)
   
    // select candidate from common node which have more votes
    fmt.Print("\n-------------Select candidate nodes...------------\n")
    fmt.Println("\t\tinfo id votes\n")
    waitingTime()
    SelectCandidate()
    fmt.Print(nodePool[:candidateNum])

    // initial candidate node list
    fmt.Print("\n----------Initializing candidate nodes...----------\n")
    fmt.Println("\t\tinfo id votes auth d \n")
    waitingTime()
    initCandidate()
    //fmt.Println(candPool)

    // simulate the auth
    fmt.Print("\n-------------------Authenticating-----------------\n")
    waitingTime()
    Auth()

    // simulate the vote
    fmt.Print("\n-----------------------Voting--------------------\n")
    waitingTime()
    Vote()
    
    // calculate the support degree
    fmt.Print("\n--------Calculate the candidates' support degree---------\n")
    waitingTime()
    CalSD()

    // selection delegate from candidate
    nodes := SelectDelegate()
    fmt.Print("\n-----------------Select delegate nodes-----------------\n")
    fmt.Println("\t\tinfo id votes auth d \n")
    waitingTime()
    fmt.Println(nodes)

    // initial consensus
    fmt.Print("\n---------------Initializing consensus...---------------\n")
    waitingTime()
    InitialDelegate()

    // create genesis block
    first := geneBlock() 
    lastBlock := first

    // generate the block
    fmt.Print("\n-------------------Generating block...--------------------\n")
    waitingTime()
    for i := 0; i < len(nodes); i++ {
        if nodes[i].Auth == 1 {
            // TODO: think add block address with hash code
            fmt.Printf("Block generated by candidate %d(delegate %d), which have authencated!\n", nodes[i].Id, i)
        } else {
            fmt.Printf("Block generated by candidate %d(delegate %d), which have unauthencated!\n", nodes[i].Id, i)
        }

        // statistic continuity
        delePool[i].Con++
        lastBlock = nodes[i].GenerateNewBlock(lastBlock, []byte(fmt.Sprintf("new block [%d]", i)), "")
    }

    // validate the block then add it into the blockchain
    // FIXME: Is this one have necessity? 
    //blockHeight := len(blockchain)
    //oldBlock := blockchain[blockHeight-1]
    //if isBlockValid(lastBlock, oldBlock) {
    //    blockchain = append(blockchain, lastBlock)
    //}
}

// update delegate's cv
func UpdateCv() {
    waitingTime()
    
    // the arguments of calculate the contribution value
    reward      := 0.05
    punish      := -0.2
    lambda1     := 0.05
    lambda2     := 0.25
    sum1        := 0.0
    sum2        := 0.0

    // FIXME: something wrong here
    newBlock := geneBlock()
    blockHeight := len(blockchain)
    oldBlock := blockchain[blockHeight-1]
    
    for i := 0; i < delegateNum; i++ {
        delta := -(float64)(delePool[i].Con)
        // FIXME: validation always be the same 
        if isBlockValid(newBlock, oldBlock) { 
            delePool[i].Cv += 0.05 // reward 0.05 cv
            rewardTimes++
            sum1 += rewardTimes * reward
            if delePool[i].Con == 0 {
                curCv[i] = delePool[i].Cv + rewardTimes * reward
            } else {
                if delePool[i].Con >= 3 {
                    curCv[i] = delePool[i].Cv + lambda1 * 1/(math.Exp(delta * sum1))
                } else {
                    curCv[i] = delePool[i].Cv + lambda2 * 1/(math.Exp(delta * sum1))
                }
            }
        } else {
            delePool[i].Cv -= 0.2 // punish cv
            punishTimes++
            sum2 += punishTimes * punish
            if delePool[i].Con == 0 {
                curCv[i] = delePool[i].Cv + punishTimes * punish 
            } else {
                if delePool[i].Con >= 3 {
                    curCv[i] = delePool[i].Cv + lambda2 * 1/(math.Exp(delta * sum2))
                } else {
                    curCv[i] = delePool[i].Cv + lambda1 * 1/(math.Exp(delta * sum2))
                }
            }
        }
    }
}

// update delegate's cl
func UpdateCl() {
    //fmt.Print("\n-------------------Updating contibution level...--------------------\n")
    waitingTime()

    for i := 0; i < delegateNum; i++ {
        if delePool[i].Cv >= 0.75 && delePool[i].Cv <= 1 {
            delePool[i].Cl = 1
        } else if delePool[i].Cv >= 0.5 && delePool[i].Cv < 0.75 {
            delePool[i].Cl = 2
        } else if delePool[i].Cv >= 0.25 && delePool[i].Cv < 0.5 {
            delePool[i].Cl = 3
        } else {
            delePool[i].Cl = 4
        }
    }
}

// display the information of delegates' cv and cl
func ShowCvCl() {
    nodes := SelectDelegate()
    for i := 0; i < delegateNum; i++ {
        if delePool[i].Auth == 1 {
            fmt.Printf("candidate %d(delegate %d) have authencated and cl = %d cv = %f.\n", nodes[i].Id, i, delePool[i].Cl, delePool[i].Cv)
        } else {
            fmt.Printf("candidate %d(delegate %d) have unauthencated and cl = %d cv = %f.\n", nodes[i].Id, i, delePool[i].Cl, delePool[i].Cv)
        }
    }
}

// call feedback system
func Feedback() {
    fmt.Print("\n----------Call feedback system to reward and punish nodes...----------\n")
    waitingTime()
    
    // according the nodes' contribution level to put them back into conresponding pool
    nodes := SelectDelegate()
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
        }else if delePool[i].Cl == 1 {
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
    fmt.Print("\nGood node we reward it 0.05 contribution value every consensus round!\n")
    fmt.Print("\nBad node we punish it 0.2 contribution value every consensus round!\n")
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
    fmt.Print("-------------------Nodes' attribution checked!--------------------\n")
}

func main() {
    // TODO: add interactive interface
    //fmt.Println("Please input node number always be odd:")
    //fmt.Scanln(&nodeNum)
    //delegateNum = (int)(nodeNum / 3)
    //fmt.Printf("This round we have %d nodes then select %d nodes be a delegate node.", nodeNum, delegateNum)
LOOP:
    // TODO: add more loop to generate the block instead of one round only
    Process()
    fmt.Print("\n-------------Updating contibution value and contribution level...--------------\n")
    UpdateCv()
    UpdateCl()
    ShowCvCl()
    fmt.Print("\n------------Contibution value and contribution level updated!------------\n")
    Feedback()
    Shuffle()
    //CheckAttr()

    //NextLoop()
    fmt.Println("\n----------------------Next loop?-----------------------\n")
    fmt.Println("Current consensus round have done, would you like to start next round? y to contine, n to stop:")
    fmt.Scanln(&input)
    if input == "y" || input == "Y" {
        goto LOOP
    } else {
        fmt.Println("Consensus endup, see you next time!")
        return
    }
}
