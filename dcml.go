package main

import (
    "fmt"
    "log"
)

// contribution machnism
func ContributionMechanism() {
	// FIXME: after generate 100 blocks, the delegate nodes' cv and cl only update once, but we should update them every 10 round
	//UpdateCv(con, unvalid)
	//Upcv(con)
	//UpdateCl()

	fmt.Print("\n-------------Updating contibution value and contribution level...--------------\n")
	waitingTime()
	ShowCvCl(round)
	fmt.Print("\n--------------Contibution value and contribution level updated!------------\n")

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

// check consensus whether end
func Check(ch bool) bool {
	if ch {
		fmt.Println("\n--------------Current consensus is over!--------------\n")
		return true
	}
	fmt.Println("\n--------------Current consensus si continuing...!---------------\n")
	return false
}

// DCML algorithm
// notified other nodes whether voted
func getNotify() {
	for i := 0; i < nodeNum; i++ {
		if delePool[i].isDelete {
			log.Warnf("\nThis node alread deleted!\n")
		} else if delePool[i].isGood {
			log.Warnf("\nThis node was good node, you can vote it more!\n")
		} else {
			log.Warnf("\nVote it by your heart!\n")
		}
	}
}

// candidate monitor
// create counter and buffer while consensus initilizing
func CandidateMonitor() {
	fmt.Println("\n---------------Start to monitor candidate nodes.----------------\n")
	waitingTime()

	fmt.Println("\n--------Creating the global counter and alternative buffer.---------\n")
	waitingTime()

	counter = 0
	buffer = nodeNum / 2
	localCnt := 0

	for i := 0; i < delegateNum; i++ {
		//counter[i] = 0
		//buffer[i]  = 0

		// read data from state table if the data exist
		// if not, we statistic it automatically
		if delePool[i].isDelete && freezePool[i].isDelete {
			counter = len(deletePool) + len(freezePool)
		} else if deletePool[i].isDelete {
			counter = len(deletePool)
		} else if freezePool[i].isDelete {
			counter = len(freezePool)
		} else {
			localCnt++
		}
	}

	// set counter to zero while counter is too big every consensus end
	if Check(check) {
		counter = 0
	}
}

// local density of one sample point
func LocalDensity() []float64 {
	var neighbor = make([]D, candidateNum)
	var density [candidateNum]float64
	var dis [candidateNum]float64

	// k-th distance, specific randomly
	// the argument we talk about in the paper
	k := 10

	// wk is weight to avoid the distance become zero
	wk := 0.01

	// calculate the distance between other sample point and point p
	// sort them then take kth distance as k-distance
	for i := 0; i < candidateNum; i++ {
		//var kd [candidateNum]float64

		// calculate the distance, but with this moment, we just simulate the distance instead of calculate the real distance
		dis[i] = float64(rand.Intn(candidateNum))
		if dis[i] == 0 {
			dis[i] += wk
		}

		// sort distance
		if dis[i] < dis[i+1] {
			dis[i] = dis[i+1]
		}
		kth := dis[k]

		if dis[i] <= kth {
			neighbor[i] = candPool[i]
		}

		// the local density of point p
		density[i] = float64(1) / kth
	}
	return density[:candidateNum]
}

// density mean
func DensityMean() float64 {
	/*
	   // find the n point which is p's neighbor
	   var neighbor  = make([]D, candidateNum)
	   for i := 0; i < candidateNum; i++ {
	       if dis[i] <= kth {
	           neighbor[i] = candPool[i]
	       }
	   }
	*/

	densityMean := 0.0
	sum := 0.0
	density := LocalDensity()

	// calculate the mean of density
	for i := 0; i < candidateNum; i++ {
		sum += density[i]
	}
	densityMean = 1 / delegateNum * sum

	return densityMean
}

// LOF score
func LOFScore() []float64 {
	density := LocalDensity()
	densityMean := DensityMean()
	var score [delegateNum]float64

	for i := 0; i < delegateNum; i++ {
		score[i] = densityMean / density[i]
	}
	return score[:delegateNum]
}

// local outlier factor, which is based density
// core thought is a node whether is abnormal or normal which depends local enviornment
// calculation steps
//      1. dis(p, kth)
//      2. k-d(p)
//      3. r-d_k(p, o) = max{k-d(o), d(p, o)} + w_k
//      4. lrd_k(p)
//      5. LOF_k(p)
// in other hand, we should to process dataset into several segments to reduce the time complexity, but here we haven't dataset
//
func LOF() {
	fmt.Println("\n---------------Start to run LOF algorithm.----------------\n")
	waitingTime()

	var anomaly [delegateNum]float64
	var normal [delegateNum]float64

	fmt.Println("\n-----------------LOF algorithm runing...-----------------\n")
	waitingTime()

	score := LOFScore()
	for i := 0; i < delegateNum; i++ {
		if score[i] <= 1.0 {
			normal[i] = score[i]
			log.Warnf("\nThe data point is Ok, not like an anomaly!\n")
		} else if score[i] > 1.0 && score[i] <= 1.3 {
			normal[i] = score[i]
			log.Warnf("\nIt seems that this point closer to others, not like an anomaly!\n")
		} else {
			anomaly[i] = score[i]
			log.Warnf("\nThis point far away other nodes, just like an anomaly!\n")
		}
	}

	fmt.Println("\n-----------------LOF Algorithm run over!-----------------\n")
	waitingTime()
}

// calculate the features average
// one-d
func FeatureAverage1D() ([]float64, []float64, []float64) {
	// initialize arguments
	var miu [candidateNum]float64
	var sigma [candidateNum]float64
	var delta [candidateNum]float64
	var x [candidateNum]float64
	miuSum := 0.0
	sigmaSum := 0.0
	deltaSum := 0.0

	for i := 1; i < candidateNum; i++ {
		for j := 0; j < candidateNum; j++ {
			// simulate the data point randomlynp.random.multivariate_normal
			x[j] = float64(rand.Intn(nodeNum))
			miu[j] = float64(rand.Intn(nodeNum))

			miuSum += x[j]
			transpose := (float64(i)*x[j] + float64(j)) - (float64(i)*miu[j] + float64(j))
			sigmaSum += (x[j] - miu[j]) * transpose
			deltaSum += (x[j] - miu[j]) * (x[j] - miu[j])

			miu[j] = 1 / candidateNum * miuSum
			sigma[j] = 1 / candidateNum * sigmaSum
			delta[j] = math.Sqrt(1 / candidateNum * deltaSum)
		}
	}

	return miu[:], sigma[:], delta[:]
}

// two-d
func FeatureAverage2D() ([][candidateNum]float64, [][candidateNum]float64, [][candidateNum]float64) {
	var miu [candidateNum][candidateNum]float64
	var sigma [candidateNum][candidateNum]float64
	var delta [candidateNum][candidateNum]float64
	var x [candidateNum][candidateNum]float64
	miuSum := 0.0
	sigmaSum := 0.0
	deltaSum := 0.0

	for i := 1; i < candidateNum; i++ {
		for j := 0; j < candidateNum; j++ {
			// simulate the data point randomly
			x[i][j] = float64(rand.Intn(nodeNum))
			miu[i][j] = float64(rand.Intn(nodeNum))

			miuSum += x[i][j]
			transpose := (float64(i)*x[i][j] + float64(j)) - (float64(i)*miu[i][j] + float64(j))
			sigmaSum += (x[i][j] - miu[i][j]) * transpose
			deltaSum += (x[i][j] - miu[i][j]) * (x[i][j] - miu[i][j])

			miu[i][j] = 1 / candidateNum * miuSum
			sigma[i][j] = 1 / candidateNum * sigmaSum
			delta[i][j] = math.Sqrt(1 / candidateNum * deltaSum)
		}
	}

	return miu[:][:], sigma[:][:], delta[:][:]
}

// calculate the probility
// one-d
func Probility1D() []float64 {
	miu, _, delta := FeatureAverage1D()
	var x [candidateNum]float64
	var pre [candidateNum]float64
	var exp [candidateNum]float64
	var prob [candidateNum]float64

	for j := 1; j < candidateNum; j++ {
		x[j] = float64(rand.Intn(nodeNum))

		pre[j] = 1 / (math.Sqrt(2*3.14) * delta[j])
		exp[j] = -1 / 2 * math.Pow(x[j]-miu[j], 2) / (math.Pow(delta[j], 2))
		prob[j] = pre[j] * math.Exp(exp[j])
	}

	return prob[:]
}

// two-d
func Probility2D() [][candidateNum]float64 {
	miu, _, delta := FeatureAverage2D()
	var x [candidateNum][candidateNum]float64
	var pre [candidateNum][candidateNum]float64
	var exp [candidateNum][candidateNum]float64
	var prob [candidateNum][candidateNum]float64

	for i := 1; i < candidateNum; i++ {
		for j := 0; j < candidateNum; j++ {
			x[i][j] = float64(rand.Intn(nodeNum))

			pre[i][j] = 1 / (math.Sqrt(2*3.14) * delta[i][j])
			exp[i][j] = -1 / 2 * math.Pow(x[i][j]-miu[i][j], 2) / (math.Pow(delta[i][j], 2))
			prob[i][j] = pre[i][j] * math.Exp(exp[i][j])
		}
	}

	return prob[:][:]
}

// judge one point whether a normal or anomaly one
// one-d
func JudgeIt1D() {
	prob := Probility1D()
	for i := 0; i < candidateNum; i++ {
		if prob[i] < epsilon1 {
			log.Warnf("\nTHIS POINT IS ANOMALY!!!!\n")
			waitingTime()
		} else {
			log.Warnf("\nTHIS POINT IS NORMAL!!!!\n")
			waitingTime()
		}
	}
}

// two-d
func JudgeIt2D() {
	prob := Probility2D()
	for i := 0; i < candidateNum; i++ {
		for j := 0; j < candidateNum; j++ {
			if prob[i][j] < epsilon2 {
				log.Warnf("\nTHIS POINT IS ANOMALY!!!!\n")
				waitingTime()
			} else {
				log.Warnf("\nTHIS POINT IS NORMAL!!!!\n")
				waitingTime()
			}
		}
	}
}

// multivariate guassian model(Multivariate normal distribution)
// calculation steps
//      1. calculate the average of every feature
//      2. calculate the model prob(x) use new samples
//      3. compare prob(x) with epsilon
func MGM() {
	fmt.Println("\n-----------------MGM algorithm runing...-----------------\n")
	waitingTime()

	if dimension == 1 {
		FeatureAverage1D()
		Probility1D()
		JudgeIt1D()
	}
	if dimension == 2 {
		FeatureAverage2D()
		Probility2D()
		JudgeIt2D()
	}

	fmt.Println("\n-----------------MGM Algorithm run over!-----------------\n")
	waitingTime()
}

// abnormal detection
func AbnormalDetection(cand D) bool {
	fmt.Println("\n---------------Start to abnormal detect...----------------\n")
	waitingTime()

	// FIXME: MGM need to filter the nodes checked after LOF instead of use MGM after LOF directly
	LOF()
	MGM()

	fmt.Println("\n----------------Abnormal detect have done!----------------\n")
	waitingTime()

	return true
}

// three alternative strategies
// alternate on time
func TimingAlternate(number int) {
	timer := threshold
	//alternateNum := candidateNum
	alternateNum := number

	// timer counter
	for i := timer; i > 0; i-- {
		for j := 0; j < alternateNum; j++ {
			// validate abnormal detection result
			if AbnormalDetection(candPool[j]) {
				alterPool[j] = candPool[j]
			} else {
				log.Warnf("\nAbnormal detection failed, AD next one!\n")
			}
		}
	}
}

// alternate smally
// should satisfied n >= 3f + 1, f is abnormal node, n is total delegate number
func MinimumAlternate() {
	tmpCnt := 0

	for i := 0; i < delegateNum; i++ {
		if deletePool[i].isDelete {
			tmpCnt++
		}
		if freezePool[i].isDelete {
			tmpCnt++
		}

		// judge the condtion whether satisfied
		if delegateNum-tmpCnt < mini {
			if AbnormalDetection(candPool[i]) {
				for j := 0; j < alternateNum; j++ {
					alterPool[j] = candPool[j]
				}
			}
		}
	}
}

// alternated accroding interval
func alternateInterval() {
	// alternated accroding interval
	for i := interval; i > 0; i-- {
		for i := 0; i < delegateNum; i++ {
			if AbnormalDetection(candPool[i]) {
				for j := 0; j < alternateNum; j++ {
					alterPool[j] = candPool[j]
				}
			}
		}
	}
}

// alternated accroding full load
func alternateFullLoad() {
	cnt := 0
	for i := 0; i < delegateNum; i++ {
		if AbnormalDetection(candPool[i]) {
			cnt++
		}
	}

	if cnt == alternateNum {
		for i := 0; i < delegateNum; i++ {
			if AbnormalDetection(candPool[i]) {
				for j := 0; j < alternateNum; j++ {
					alterPool[j] = candPool[j]
				}
			}
		}
	} else if cnt > alternateNum {
		log.Warnf("\nThe buffer is full, please store qualified candidate nodes after nodes alternated in the buffer!\n")
		waitingTime()
	} else {
		log.Warnf("\nThe buffer need to fill, please going on...\n")
		waitingTime()
	}
}

// statistic bad nodes' number
func Statistic() int {
	tmpCnt := 0
	for i := 0; i < delegateNum; i++ {
		if deletePool[i].isDelete {
			tmpCnt++
		}
		if freezePool[i].isDelete {
			tmpCnt++
		}
	}
	return tmpCnt
}

// alternate regularly
func RegularAlternate() {
	// tmpCnt means the number of need to alternate
	tmpCnt := 0
	for i := 0; i < delegateNum; i++ {
		if deletePool[i].isDelete {
			tmpCnt++
		}
		if freezePool[i].isDelete {
			tmpCnt++
		}
	}

	// accroding tmpCnt to select alternative mode
	if tmpCnt >= 3 {
		alternateInterval()
	} else {
		alternateFullLoad()
	}
}

// alternative strategy
func SelectAlternativeStrategy() {
	fmt.Println("\n---------------Start to select alternative strategy...----------------\n")
	waitingTime()

	tmpCnt := 0
	for i := 0; i < delegateNum; i++ {
		if deletePool[i].isDelete {
			tmpCnt++
		}
		if freezePool[i].isDelete {
			tmpCnt++
		}
	}

	// we use timing alternate when the delete node and freeze node more than half of delegate number
	// if not, we use mininum alternate when the rest node is equals mini number
	// else, we use regular alternate mode
	if tmpCnt > delegateNum/2 {
		TimingAlternate(candidateNum)
	} else if delegateNum-tmpCnt == mini {
		MinimumAlternate()
	} else {
		RegularAlternate()
	}

	fmt.Println("\n---------------Alternative strategy already selected!----------------\n")
	waitingTime()
}

// alternate dynamicly
// what we called alternate dynamicly means that we select nodes from alterPool to participate the consensus
func DynamicAlternate() {
	fmt.Println("\n-------------------Start to dynamic alternating...---------------------\n")
	waitingTime()

	CandidateMonitor()

	//This step, abnormal detection was executed in the SelectAlternativeStrategy()
	SelectAlternativeStrategy()

	fmt.Println("\n-------------------Dynamic alternating have done!----------------------\n")
	waitingTime()
}

// DCML algorithm
// check nodes' attributions then notify other nodes whether this one is ok with broadcast
// and then, alternate dynamicly
func DCML() {
	fmt.Println("\n----------------------Runing DCML algorithm...----------------------\n")
	waitingTime()

	CheckAttr()
	getNotify()
	DynamicAlternate()

	fmt.Println("\n----------------------DCML algorithm run over!----------------------\n")
	waitingTime()
}
