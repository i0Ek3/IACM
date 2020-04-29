# IACM

> **[WIP]**

IACM is an improved DPoS consensus, but this repo just simplized, thanks for those people who contribute to the Golang DPoS implementation.

Of course, all stuff is simulation instead of real trail environment, so be careful!

## IDTM(WIP)

Main point:

Node authencation and support degree calculation, then select the delegate nodes. While consensus start, use contribution meachism to update contirbution value and it's level, take measures to reward and punish node. 

## DCML(Not done yet!)

Main point:

When bad nodes have taken down, the blank position should affect consensus process, so we use some methods such as candidate monitor and machine learning algorithms to detect latent bad nodes and candidate proper nodes into consensus process to consensus well.

## Update Log

- v0429: Fix and make some bugs.
- v0428: Modelize the code, fix some issues but not all.
- v0425: Finished basic function of IDTM.

## TODO

- Add consensus delay
- Statistic and validate round times 
- Add P2P function
- Add interactive interface
- Fix same Cv/Cl value of every delegate node
- Implement DCML algorithm
- ...

## Credit 

Big thanks to GitHub with open source code.
