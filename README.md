# IACM

IACM is an improved DPoS consensus implement by Go, which includes algorithm IDTM and DCML. Of course, all stuff here just simulation instead of real ones, so if you want to use it, please be careful and YOU OWN RISK. After all, this repo just for personal academic research, BE KNOWN!!!

## Build & Run

You can build and run this file in macOS/Windows/Linux platform with same command:

- Build: `$ go build iacm.go`
- Run:   `$ go run iacm.go`

## IDTM

Main point:

Node authencation and support degree calculation, then select the delegate nodes. While consensus start, use contribution meachism to update contirbution value and it's level, take measures to reward and punish node. 

## DCML

Main point:

When bad nodes have taken down, the blank position should affect consensus process, so we use some methods such as candidate monitor and machine learning algorithms to detect latent bad nodes and candidate proper nodes into consensus process to consensus well.

## Update Log

- 20201025: Discard refactor cause of I cannot understand it yet, bye~
- 20201012: Support go mod and remove unuseful lines.
- 20200917: Add error debug information partially.
- 20200510: Add and upgrade comparison algorithm.
- 20200509: Finished algorithm MGM but we got nothing about the data. 
- 20200508: Implement part of alternative strategies and LOF algorithm in DCML.
- 20200430: Messed up, more bugs on the way. 
- 20200429: Fix and make some bugs.
- 20200428: Modelize the code, fix some issues but not all.
- 20200425: Finished basic function of IDTM.

## TODO

> If I have time I will fix follows issues and todos.

- [x] Fix same block issue after 3 block generated
- [x] Implement DCML algorithm
- [x] Finished FCSW comparison algorithm
- [ ] Fix same Cv/Cl value of every delegate node
- [ ] Start to block 2 after every 10 blocks generated
- [ ] Add consensus delay
- [ ] Statistic and validate round times 
- [ ] Add P2P function
- [ ] Add interactive interface
- [ ] Output data required
- [ ] Bad nodes fails to simulate

## Credit 

Big thanks to GitHub with open source code.
