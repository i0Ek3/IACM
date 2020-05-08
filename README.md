# IACM

> **[WIP]**

IACM is an improved DPoS consensus, but this repo just simplized, thanks for those people who contribute to the Golang DPoS implementation.

Of course, all stuff is simulation instead of real trail environment, so be careful!

## IDTM

Main point:

Node authencation and support degree calculation, then select the delegate nodes. While consensus start, use contribution meachism to update contirbution value and it's level, take measures to reward and punish node. 

## DCML

Main point:

When bad nodes have taken down, the blank position should affect consensus process, so we use some methods such as candidate monitor and machine learning algorithms to detect latent bad nodes and candidate proper nodes into consensus process to consensus well.

## Update Log

- v0508: Implement part of alternative strategies and LOF algorithm in DCML.
- v0430: Messed up, more bugs on the way. >_<
- v0429: Fix and make some bugs.
- v0428: Modelize the code, fix some issues but not all.
- v0425: Finished basic function of IDTM.

## TODO

- [x] Fix same block issue after 3 block generated
- [x] Implement DCML algorithm(WIP)
- [ ] Fix same Cv/Cl value of every delegate node
- [ ] Start to block 2 after every 10 blocks generated
- [ ] Add consensus delay
- [ ] Statistic and validate round times 
- [ ] Add P2P function
- [ ] Add interactive interface


## Credit 

Big thanks to GitHub with open source code.
