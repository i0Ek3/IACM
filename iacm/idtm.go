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
    "sort"
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
    reward          = 0.005
    punish          = -0.02
    lambda          = 0.5
    lambda1         = 0.25
    lambda2         = 0.05
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

// simplize the struct, just use Node and D struct
type D struct {
    Node
    Auth        int
    d           int
    Cl          int
    Cv          float64
    Con         int
    Unvalid     int
    isDelete    bool
    isGood      bool
    Address     string
    fmm         string
}

// the struct of Block
type Block struct {
    Timestamp   string
    Prehash     string
    Hash        string
    Data        string 
    Height      int
    Address     string
    //delegate    *Node
}

var (
    nodePool    = make([]Node, nodeNum)  // common node pool
    candPool    = make([]D, candidateNum)// candidate node pool
    delePool    = make([]D, nodeNum) // delegate node pool
    
    deletePool  = make([]D, delegateNum) // state table: store cl=4's nodes
    freezePool  = make([]D, delegateNum) // state table: store cl=3's nodes
    commonPool  = make([]D, delegateNum) // state table: store cl=2's nodes
    premiumPool = make([]D, delegateNum) // state table: store cl=1's nodes
    
    round       = 0          // times of round 
    input       string       // the anwser to next loop
    blockchain  []Block      // the blockchain
    curCv       [delegateNum]float64  // current contribution value
    rewardTimes float64      // times of reward
    punishTimes float64      // times of punish
    vote        int
)

// first block: genesis block
func genesisBlock() Block {
    // Prehash have 64 bit, address have 8 bit
    gene := Block{time.Now().String(), "0000000000000000000000000000000000000000000000000000000000000000", "", "I'm the genesis block", 1, "0x0000"}
    blockchain = append(blockchain, gene)
    //gene.Hash = string(gene.calHash())
    gene.calHash()
    return gene
}

// the new version of generate the block
func generateBlock(oldBlock Block, data string, addr string) Block {
    newBlock := Block{}
    newBlock.Timestamp = time.Now().String()//Format("2020-01-01 00:00:00")
	newBlock.Prehash = oldBlock.Hash
	newBlock.calHash()
	newBlock.Data = data
	newBlock.Height = oldBlock.Height + 1
	newBlock.Address = addr
	return newBlock
}

// generate the block
func (node *Node) GenerateNewBlock(lastBlock Block, data string, addr string) Block {
    time.Sleep(3 * time.Second) // for easy use, every 3sec generate a block
    //Block{lastBlock.Height+1, time.Now().String(), lastBlock.Hash, "", data, addr, nil}
    newBlock := Block{}
    newBlock.Timestamp = time.Now().String()
    newBlock.Prehash = lastBlock.Hash 
    newBlock.Data = data
    newBlock.Height = lastBlock.Height + 1
    newBlock.Address = addr
    //newBlock.Hash = hex.EncodeToString(newBlock.calHash())
    //newBlock.delegate = node
    return newBlock
}

// calculate the block hash
func (block *Block) calHash()  {
    hashstr := strconv.Itoa(block.Height) + block.Timestamp + block.Prehash + block.Data + block.Address
    hash := sha256.Sum256([]byte(hashstr))
	block.Hash = hex.EncodeToString(hash[:])
    //h := sha256.New()
    //h.Write(hash)
    //hashed := h.Sum(nil)
    //return hashed
}

// create nodes then initialize it
func CreateNode() {
    fmt.Print("\n------------------Initializing common nodes------------------\n")
    fmt.Println("\t\t\tinfo id votes \n")
    waitingTime()
    for i := 0; i < nodeNum; i++ {
        info := fmt.Sprintf("common node:") // node information
        id   := i       // node id number
        vote := rand.Intn(nodeNum)          // the number of vote
        f    := ""    // format control
        nodePool[i] = Node{info, id, vote, f}
        fmt.Println("initializing...", nodePool[i])
        //fmt.Println("\n")
    }
}

// use lambda to sort nodes
func sortNodes() {
	sort.Slice(candPool, func(i, j int) bool {
		return candPool[i].Votes > candPool[j].Votes
	})
	delePool = candPool[:delegateNum]
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
        id   := i
        vote := 0
        f    := ""
        auth := 0
        d    := 0
        cl   := 0
        cv   := 0.0
        con  := 0
        un   := 0
        bad  := false
        good := false
        addr := ""
        fmm  := ""

        // TODO: for information shows well, we can forbid the struct inheritance 
        candPool = append(candPool, D{Node{info, id, vote, f}, auth, d, cl, cv, con, un, bad, good, addr, fmm})
        candPool[i] = D{Node{info, id, vote, f}, auth, d, cl, cv, con, un, bad, good, addr, fmm}
        
        // use key:value method of struct but shows error here
        //candPool[i] = D{Node{Info:"candidate", Id:i, Votes:0}, Auth:0, d:0, fmm:"\n"}
        fmt.Print("\ninitilizing...",  candPool[i])
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
            candPool[i].d = (int)(alpha * (float64)(candPool[i].Auth) + beta * (float64)(candPool[i].Votes))
            fmt.Printf("candidate %d have %d votes and authencated, the support degree is %d.\n", candPool[i].Id, candPool[i].Votes, candPool[i].d)
        } else {
            candPool[i].d = (int)((1-alpha) * (float64)(candPool[i].Auth) + beta * (float64)(candPool[i].Votes))
            fmt.Printf("candidate %d have %d votes and unauthencated, the support degree is %d.\n", candPool[i].Id, candPool[i].Votes, candPool[i].d)
        }
    }
}

// select delegate node from candidate
func SelectDelegate() []D {
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

    for i := 0; i < delegateNum; i++ {
        fmt.Println(candPool[i])
    }

    //delePool = append(delePool, candPool)
    return n[:delegateNum]
}

// initialize consensus: initialize delegate node
func InitialDelegate() {
    //nodes := SelectDelegate()
    for i := 0; i < delegateNum; i++ {
        delePool[i].Cl       = 2
        delePool[i].Cv       = 0.05
        delePool[i].Con      = 0
        delePool[i].Unvalid  = 0
        delePool[i].isDelete = false
        delePool[i].isGood   = false
        
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
        id   := i
        vote := 0
        f    := ""
        auth := candPool[i].Auth
        d    := candPool[i].d
        cl   := 2
        cv   := 0.5
        con  := 0
        un   := 0
        bad  := false
        good := false
        addr := "0x00" + strconv.Itoa(i+1)
        fmm  := ""

        delePool = append(delePool, D{Node{info, id, vote, f}, auth, d, cl, cv, con, un, bad, good, addr, fmm})
        delePool[i] = D{Node{info, id, vote, f}, auth, d, cl, cv, con, un, bad, good, addr, fmm}
        fmt.Println("initialized to delegate %d:", candPool[i], i)
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
func isBlockValid(newBlock, oldBlock Block) bool {
    fmt.Println("\n------------------Validating the block...--------------------\n")
    waitingTime()
    
    // TODO: validate delegateNum delegate nodes
    for i := 0; i < delegateNum; i++ {
        //
    }
    
    // validate height and prehash
    if oldBlock.Height + 1 != newBlock.Height {
   	    fmt.Println("\nValidation failed! Wrong Height!\n")
   	    return false
    } else {
        fmt.Println("\nBlock Height validating successful!\n")
    }
    if newBlock.Prehash != oldBlock.Hash {
   	    fmt.Println("\nValidation failed! Wrong Prehash!")
   	    return false
    } else {
        fmt.Println("Block Prehash validating successful!\n")
    }
    fmt.Println("\n-------------------Block validated!--------------------\n")
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

// process to initializing until before consensus
func Process() {
    // statistic round number
    round++
    
    // create common nodes then print them information
    CreateNode()
    //fmt.Print("\n----------------Initializing nodes---------------\n")
    //fmt.Println("\t\tinfo id votes \n")
    //waitingTime()
    //fmt.Println(nodePool)
   
    // select candidate from common node which have more votes
    fmt.Print("\n-------------Select candidate nodes...------------\n")
    fmt.Println("\t\tinfo id votes\n")
    waitingTime()
    SelectCandidate()
    //fmt.Print(nodePool[:candidateNum])

    // initial candidate node list
    fmt.Print("\n----------Initializing candidate nodes...----------\n")
    fmt.Println("\tinfo id votes auth d cl cv con bad good \n")
    waitingTime()
    InitCandidate()
    //fmt.Println(candPool)

    // simulate the auth
    fmt.Println("\n")
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
    //nodes := SelectDelegate()
    fmt.Print("\n-----------------Select delegate nodes-----------------\n")
    fmt.Println("\tinfo id votes auth d cl cv con bad good \n")
    waitingTime()
    SelectDelegate()
    fmt.Println("\n")
    //fmt.Println(nodes)

    // initial consensus
    fmt.Print("\n---------------Initializing consensus...---------------\n")
    waitingTime()
    InitialDelegate()
    
    // old generate block version
    //waitingTime()
    //nodes := SelectDelegate()
    //for i := 0; i < len(nodes); i++ {
    //    if nodes[i].Auth == 1 {
    //        fmt.Printf("Block generated by candidate %d(delegate %d), which have authencated!\n", nodes[i].Id, i)
    //    } else {
    //        fmt.Printf("Block generated by candidate %d(delegate %d), which have unauthencated!\n", nodes[i].Id, i)
    //    }
}

// generate the block recurrently
func genLoop() {
    // create genesis block
    //first := genesisBlock() 
    //lastBlock := first

    // generate the block
    fmt.Print("\n------------------------Generating block...-------------------------\n")
    waitingTime()
    nodes := SelectDelegate()
    for i := 0; i < len(nodes); i++ {
        if nodes[i].Auth == 1 {
            // TODO: think add block address with hash code
            fmt.Printf("Block generated by candidate %d(delegate %d), which have authencated!\n", nodes[i].Id, i)
        } else {
            fmt.Printf("Block generated by candidate %d(delegate %d), which have unauthencated!\n", nodes[i].Id, i)
        }

        // statistic continuity
        delePool[i].Con++
        //lastBlock = generateBlock(blockchain[i], fmt.Sprintf("block content"), "")
    }

    // validate the block then add it into the blockchain
    // FIXME: Is this one have necessity? 
    blockHeight := len(blockchain)
    oldBlock := blockchain[blockHeight-1]
    newBlock := generateBlock(oldBlock, "block content", "")
    if isBlockValid(newBlock, oldBlock) {
        blockchain = append(blockchain, newBlock)
    }
}

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
            fmt.Println("\nsomething wrong here!\n")
        }
    }
}

// update delegate's cv
// yeah, we messed up this one!
func UpdateCv(con, unvalid int) { 
    // the arguments of calculate the contribution value
    sum1        := 0.0
    sum2        := 0.0

    blockHeight := len(blockchain)
    oldBlock := blockchain[blockHeight-1]
    newBlock := generateBlock(oldBlock, "block content", "")
    
    // first, judge the block whether valid, if it is, then judge the con's number 
    // if con is big more and more than threshold, we reward this delegate much more
    // if not, con > 0 and con < threshold, we reward it a little
    // in other words, we needn't validate the block whether it valid or not
    // cause of con, if con > 0 that means the block must be validated and valid
    // else, while con = 0 that means block is unvalid
    
    for i := 0; i < delegateNum; i++ {
        delta := -(float64)(delePool[i].Con)
        
        if oldBlock.Height+1 == newBlock.Height || newBlock.Prehash == oldBlock.Hash { 
            delePool[i].Con++
            if con > 0 {
                //rewardTimes += con
            } else if con == 0 {
                punishTimes++
            } else {
                //
            }
            blockchain = append(blockchain, newBlock)
            
            if con >= 3 {
                delePool[i].Cv += rewardTimes * reward // reward cv
            }
            sum1 += rewardTimes * reward

            // con = 0 means this delegate not product the block
            if con == 0 {
                curCv[i] = delePool[i].Cv + punishTimes * punish
            } else {
                // that means this delegate is good one
                if con >= 3 {
                    curCv[i] = delePool[i].Cv + rewardTimes * reward
                    //curCv[i] = delePool[i].Cv + lambda1 * 1/(math.Exp(delta * sum1))
                } else {
                //else if con > 0 && con < 3 {
                    l := (lambda1 + lambda2) / 2
                    curCv[i] = delePool[i].Cv + l * 1/(math.Exp(delta * sum1))
                }
            }
        } else {
            delePool[i].Cv += punish // punish cv
            punishTimes++
            sum2 += punishTimes * punish
           
            if oldBlock.Height+1 != newBlock.Height || newBlock.Prehash != oldBlock.Hash { 
                if delePool[i].Unvalid >= 3 {
                    curCv[i] = delePool[i].Cv + punishTimes * punish 
                }
            } else {
                if con >= 3 {
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
    //nodes := SelectDelegate()

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
    fmt.Printf("\nGood node we reward it %f contribution value every consensus round!\n", reward)
    fmt.Printf("\nBad node we punish it %f contribution value every consensus round!\n", punish)
    fmt.Print("\n-------------------------------------------------------------------\n")
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

// contribution machnism
func ContributionMechanism() {
    // FIXME: after generate 100 blocks, the delegate nodes' cv and cl only update once, but we should update them every 10 round
    //UpdateCv(con, unvalid)
    //Upcv(con)
    //UpdateCl()
    
    fmt.Print("\n-------------Updating contibution value and contribution level...--------------\n")
    waitingTime()
    ShowCvCl(round)
    fmt.Print("\n------------Contibution value and contribution level updated!------------\n")
    
    Feedback()
    Shuffle()
}

// when block validated successful, we should broadcast it to other nodes
// but here we just simulate it
func Broadcast() {
    fmt.Println("\n--------------------Broadcasting...---------------------\n")
    waitingTime()
    fmt.Println("\nSend message to other nodes: block generated and validated, please record it!\n")
    fmt.Println("\n--------------------Broadcast done!---------------------\n")
}

// DCML algorithm 
// get votes notify
func getNotify() {
    for i := 0; i < nodeNum; i++ {
        if delePool[i].isDelete {
            fmt.Println("\nThis node alread deleted!\n")
        } else if delePool[i].isGood {
            fmt.Println("\nThis node was good node, you can vote it more!\n")
        } else {
            fmt.Println("\nVote it by your heart!\n")
        }
    }
}

// candidate monitor
func CandidateMonitor() {

}

// local outlier factor
func LOF() {

}

// multivariable guassian model
func MGM() {

}

// abnormal detection
func AbnormalDetection() {
    LOF()
    MGM()
}

// three alternative strategies
// alternate on time
func TimingAlternate() {

}

// alternate smally
func MinimumAlternate() {

}

// alternate regularly 
func RegualrAlternate() {

}

// alternative strategy
func SelectAlternativeStrategy() {

}

// alternate dynamicly
func DynamicAlternate() {
    CandidateMonitor()
    AbnormalDetection()
    SelectAlternativeStrategy()
}

// DCML algorithm calling
func DCML() {
    //CandidateMonitor() // dynamic alternataion algo step
    //AbnormalDetection()
    //SelectAlternativeStrategy() // includes three alternative stragtegies
    CheckAttr()
    getNotify()
    DynamicAlternate() 
}

// main function
func main() {
LOOP:
    Process()
    
    // TODO: add interactive interface
    //fmt.Println("Please input node number always be odd:")
    //fmt.Scanln(&nodeNum)
    //delegateNum = (int)(nodeNum / 3)
    //fmt.Printf("This round we have %d nodes then select %d nodes be a delegate node.", nodeNum, delegateNum)

    gene := Block{time.Now().String(), "0000000000000000000000000000000000000000000000000000000000000000", "","I'm the genesis block", 1, "0x0000"}
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
    for k := 0; k < delegateNum * i - 1; k++ {
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
        
        //UpdateCv(nodes[k].Con, nodes[k].Unvalid)
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
        if round % 10 == 0 {
            // FIXME: chaos here
            //ContributionMechanism()
        } 
    }
    //ContributionMechanism()
    //DCML()

    // run the next loop
    fmt.Println("\n---------------------------Next loop?----------------------------\n")
    fmt.Println("Current consensus round have done, would you like to start next round? y to contine, n to stop:")
    
    // interaction in the end
    fmt.Scanln(&input)
    if input == "y" || input == "Y" {
        goto LOOP
    } else {
        fmt.Println("Consensus endup, see you next time!")
        return
    }
}
