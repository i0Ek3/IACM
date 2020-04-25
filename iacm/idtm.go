// IACM-IDTM
// Name: idtm.go
// 
//

package main

import (
    "crypto/sha256"
    "encoding/hex"
    "math/rand"
    "strconv"
    "math"
    "time"
    "fmt"
)

// the struct of Node
type Node struct {
    Info        string
    Id          int
    Votes       int
    Auth        int 
    d           int
    Cl          int
    Cv          float64
    Con         int
    isDelete    bool
    isGood      bool
    f           string
}

// the struct of Block
type Block struct {
    Index       int
    Timestamp   string
    Hash        string
    Prehash     string
    Data        []byte
    delegate    *Node
}

var (
    nodeNum     = 15   // sum of nodes
    delegateNum = 5    // sum of delegate nodes
    alpha       = 0.8      
    beta        = 0.2
    NodeAddr    = make([]Node, nodeNum)
    deletePool  = make([]Node, delegateNum) // store cl=4's nodes
    freezePool  = make([]Node, delegateNum) // store cl=3's nodes
    commonPool  = make([]Node, delegateNum) // store cl=1|2's nodes
    input       string     // the anwser to next loop
    blockchain  []Block    // the blockchain
    curCv       [10]float64  // current contribution value
    rewardTimes float64    // times of reward
    punishTimes float64    // times of punish
)

// first block
func geneBlock() Block {
    gene := Block{0, time.Now().String(), "", "", []byte("I'm the first block"), nil}
    blockchain = append(blockchain, gene)
    gene.Hash = string(calHash(gene))
    return gene
}

// calculate block hash
func calHash(block Block) []byte {
    hash := strconv.Itoa(block.Index) + block.Timestamp + block.Prehash + hex.EncodeToString(block.Data)
    h := sha256.New()
    h.Write([]byte(hash))
    hashed := h.Sum(nil)
    return hashed
}

// create nodes then initialize it
func CreateNode() {
    for i := 0; i < nodeNum; i++ {
        info := fmt.Sprintf("node")
        name := i
        vote := 0
        auth := 0 
        d    := 0
        cl   := 0
        cv   := 0.0
        con  := 0
        bad  := false
        good := false
        f    := "\n"
        NodeAddr[i] = Node{info, name, vote, auth, d, cl, cv, con, bad, good, f}
    }
}

// vote simulation
func Vote() {
    for i := 0; i < nodeNum; i++ {
        rand.Seed(time.Now().UnixNano())
        time.Sleep(100000)
        vote := rand.Intn(nodeNum * 10) // every node have 10 tickets
        // TODO: vote normalization
        //v := (int)(1 / (1 + math.Exp((float64)(vote))))
        if vote > 140 {
            vote -= 10
        } 
        if vote <= 10 {
            vote += 10
        }
        NodeAddr[i].Votes = vote
        fmt.Printf("node %d votes: %d\n", i, vote)
    }
}

// random authencation
func Auth() {
    for i := 0; i < nodeNum; i++ {
        rand.Seed(time.Now().UnixNano())
        time.Sleep(100000)
        auth := rand.Intn(2) // output 0 or 1 randomly
        NodeAddr[i].Auth = auth
        fmt.Printf("node %d auth: %d\n", i, auth)
    }
}

// calculate support degree, also call it sd
func CalSD() {
    for i := 0; i < nodeNum; i++ {
        if NodeAddr[i].Auth == 1 {
            NodeAddr[i].d = (int)(alpha * (float64)(NodeAddr[i].Auth) + beta * (float64)(NodeAddr[i].Votes))
            fmt.Printf("node %d support degree: %d\n", i, NodeAddr[i].d)
        } else {
            NodeAddr[i].d = (int)((1-alpha) * (float64)(NodeAddr[i].Auth) + beta * (float64)(NodeAddr[i].Votes))
            fmt.Printf("node %d support degree: %d\n", i, NodeAddr[i].d)
        }
    }
}

// select delegate node
func SelectDelegate() []Node {
    n := NodeAddr
    for i := 0; i < len(n); i++ {
        for j := 0; j < len(n)-1; j++ {
            if n[j].d < n[j+1].d {
                n[j], n[j+1] = n[j+1], n[j]
            }
        }
    }
    return n[:delegateNum]
}

// delegate nodes consensus
// generate block
func (node *Node) GenerateNewBlock(lastBlock Block, data []byte) Block {
    time.Sleep(3 * time.Second) // for easy, every 3sec generate a block
    var newBlock = Block{lastBlock.Index+1, time.Now().String(), lastBlock.Hash, "", data, nil}
    newBlock.Hash = hex.EncodeToString(calHash(newBlock))
    newBlock.delegate = node
    return newBlock
}

func waitingTime() {
    if delegateNum <= 10 {
        time.Sleep(5 * time.Second)
    } else {
        time.Sleep(10 * time.Second)
    }
}

// validate the block
func isBlockValid(newBlock, oldBlock Block) bool{
    fmt.Println("\n------------------Validating the block...--------------------\n")
    waitingTime()
    if oldBlock.Index + 1 != newBlock.Index{
   	    fmt.Println("\nValidation failed! Wrong index!\n")
   	    return false
    }
    if newBlock.Prehash != oldBlock.Hash{
   	    fmt.Println("\nValidation failed! Wrong Prehash!\n")
   	    return false
    }
   fmt.Println("\n-------------------Validation Successful!---------------------\n")
   return true
}

// initialize consensus
func Initial() {
    for i := 0; i < delegateNum; i++ {
        NodeAddr[i].Cl = 2
        NodeAddr[i].Cv = 0.5
        NodeAddr[i].Con = 0
        if NodeAddr[i].Auth == 1 {
            fmt.Printf("node %d votes %d authencated cl %d cv %f con %d.\n", NodeAddr[i].Id, NodeAddr[i].Votes, NodeAddr[i].Cl, NodeAddr[i].Cv, NodeAddr[i].Con)
        } else {
            fmt.Printf("node %d votes %d unauthencated cl %d cv %f con %d.\n", NodeAddr[i].Id, NodeAddr[i].Votes, NodeAddr[i].Cl, NodeAddr[i].Cv, NodeAddr[i].Con)
        }
    }
}

// random shuffle
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
	tmp := NodeAddr[idx1]
	NodeAddr[idx1] = NodeAddr[idx2]
	NodeAddr[idx2] = tmp
}

// process
func Process() {
    // create then print nodes
    CreateNode()
    fmt.Print("\n---------------Initialize node list-----------------\n")
    fmt.Println("\tid votes auth d cl cv con bad good\n")
    waitingTime()
    fmt.Println(NodeAddr)

    // vote then auth and calculate support degree
    fmt.Print("\n--------------------Node voting--------------------\n")
    waitingTime()
    Vote()
    fmt.Print("\n---------------Node authenticating--------------------\n")
    waitingTime()
    Auth() // auth
    fmt.Print("\n--------Calculate the node's support degree---------\n")
    waitingTime()
    CalSD() // calculate sd

    // selection delegate nodes then print
    nodes := SelectDelegate()
    fmt.Print("\n-------------------Delegate list--------------------\n")
    waitingTime()
    fmt.Println(nodes)

    // initial consensus
    fmt.Print("\n-------------Initializing consensus...---------------\n")
    waitingTime()
    Initial()

    // create genesis block
    first := geneBlock() 
    lastBlock := first

    // delegate nodes start to generating block
    fmt.Print("\n-------------------Generating block...--------------------\n")
    waitingTime()
    for i := 0; i < len(nodes); i++ {
        if nodes[i].Auth == 1 {
            fmt.Printf("Block generated by node %d, which have %d votes and authencated!\n", nodes[i].Id, nodes[i].Votes)
        } else {
            fmt.Printf("Block generated by node %d, which have %d votes and unauthencated!\n", nodes[i].Id, nodes[i].Votes)
        }
        NodeAddr[i].Con++
        lastBlock = nodes[i].GenerateNewBlock(lastBlock, []byte(fmt.Sprintf("new block [%d]", i)))
    }

    // validating the block then add into blockchain
    blockHeight := len(blockchain)
    oldBlock := blockchain[blockHeight-1]
    if isBlockValid(lastBlock, oldBlock) {
        blockchain = append(blockchain, lastBlock)
    }
}

// update delegate's cv
func UpdateCv() {
    fmt.Print("\n-------------------Updating contibution value...--------------------\n")
    waitingTime()
    
    reward      := 0.05
    punish      := -0.2
    lambda1     := 0.05
    lambda2     := 0.25
    sum1        := 0.0
    sum2        := 0.0

    newBlock := geneBlock()
    blockHeight := len(blockchain)
    oldBlock := blockchain[blockHeight-1]
    
    for i := 0; i < delegateNum; i++ {
        delta := -(float64)(NodeAddr[i].Con)
        if isBlockValid(newBlock, oldBlock) { 
            NodeAddr[i].Cv += 0.05 // reward
            rewardTimes++
            sum1 += rewardTimes * reward
            if NodeAddr[i].Con == 0 {
                curCv[i] = NodeAddr[i].Cv + rewardTimes * reward
            } else {
                if NodeAddr[i].Con >= 3 {
                    curCv[i] = NodeAddr[i].Cv + lambda1 * 1/(math.Exp(delta * sum1))
                } else {
                    curCv[i] = NodeAddr[i].Cv + lambda2 * 1/(math.Exp(delta * sum1))
                }
            }
        } else {
            NodeAddr[i].Cv -= 0.2 // punish
            punishTimes++
            sum2 += punishTimes * punish
            if NodeAddr[i].Con == 0 {
                curCv[i] = NodeAddr[i].Cv + punishTimes * punish 
            } else {
                if NodeAddr[i].Con >= 3 {
                    curCv[i] = NodeAddr[i].Cv + lambda2 * 1/(math.Exp(delta * sum2))
                } else {
                    curCv[i] = NodeAddr[i].Cv + lambda1 * 1/(math.Exp(delta * sum2))
                }
            }
        }
    }

    fmt.Print("\n-------------------Contibution value updated!--------------------\n")
}

// update delegate's cl
func UpdateCl() {
    fmt.Print("\n-------------------Updating contibution level...--------------------\n")
    waitingTime()

    for i := 0; i < delegateNum; i++ {
        if NodeAddr[i].Cv >= 0.5 && NodeAddr[i].Cv < 0.75 {
            NodeAddr[i].Cl = 2
        } else if NodeAddr[i].Cv >= 0.75 && NodeAddr[i].Cv <= 1 {
            NodeAddr[i].Cl = 1
        } else if NodeAddr[i].Cv >= 0.25 && NodeAddr[i].Cv < 0.5 {
            NodeAddr[i].Cl = 3
        } else {
            NodeAddr[i].Cl = 4
        }
    }
    fmt.Print("\n-------------------Contibution level updated!--------------------\n")
}

func ShowCvCl() {
    for i := 0; i < delegateNum; i++ {
        if NodeAddr[i].Auth == 1 {
            fmt.Printf("node %d votes %d authencated cl %d cv %f con %d.\n", NodeAddr[i].Id, NodeAddr[i].Votes, NodeAddr[i].Cl, NodeAddr[i].Cv, NodeAddr[i].Con)
        } else {
            fmt.Printf("node %d votes %d unauthencated cl %d cv %f con %d.\n", NodeAddr[i].Id, NodeAddr[i].Votes, NodeAddr[i].Cl, NodeAddr[i].Cv, NodeAddr[i].Con)
        }
    }
}

// call feedback system, reward good nodes and punish bad nodes
func Feedback() {
    fmt.Print("\n----------Call feedback system to reward and punish nodes...----------\n")
    waitingTime()
    for i := 0; i < delegateNum; i++ {
        if NodeAddr[i].Cl == 3 || NodeAddr[i].Cl == 4 {
            NodeAddr[i].isDelete = true
            fmt.Printf("Delegate %d is bad node, already taken down!\n", NodeAddr[i].Id)
        } else if NodeAddr[i].Cl == 1 {
            NodeAddr[i].isGood = true
            fmt.Printf("Delegate %d is good node, already got reward.\n", NodeAddr[i].Id)
        } else {
            fmt.Printf("Delegate %d is common node, already back to pool.\n", NodeAddr[i].Id)
        }
    }
    fmt.Print("\n------------------Reward and punish have done!--------------------\n")
}

// check nodes' attribution back from last round
func CheckAttr() {
    fmt.Print("-------------------Checking nodes' attribution...--------------------\n")
    waitingTime()
    for i := 0; i < nodeNum; i++ {
        if NodeAddr[i].isDelete && NodeAddr[i].Cl == 4 {
            deletePool = append(deletePool, NodeAddr[i])
        } else if NodeAddr[i].isDelete && NodeAddr[i].Cl == 3 {
            freezePool = append(freezePool, NodeAddr[i])
            // TODO: here should use thread to run follows code!
            // freezen on round then add these cl=3's node into commonPool
            // time.Sleep(3 * time.Second)
            // commonPool = append(commonPool, NodeAddr[i])
        } else if NodeAddr[i].isGood && NodeAddr[i].Cl == 1 {
            commonPool = append(commonPool, NodeAddr[i])
        } else {
            commonPool = append(commonPool, NodeAddr[i])
        }
    }
    fmt.Print("-------------------Nodes' attribution checked!--------------------\n")
}

func main() {
LOOP:
    Process()
    UpdateCv()
    UpdateCl()
    ShowCvCl()
    Feedback()
    Shuffle()
    CheckAttr()

    //NextLoop()
    fmt.Println("\nCurrent consensus round have done, would you like to start next round? y to contine, n to stop:")
    fmt.Scanln(&input)
    if input == "y" || input == "Y" {
        goto LOOP
    } else {
        fmt.Println("Consensus endup, see you next time!")
        return
    }
}
