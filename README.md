# IACM

IACM is an improved DPoS consensus implement by Go, which includes algorithm IDTM and DCML. Of course, all stuff here just simulation instead of real ones, so if you want to use it, please be careful and YOU OWN RISK. After all, this repo just for personal academic research, BE KNOWN!!!

## IDTM

Main point:

Node authencation and support degree calculation, then select the delegate nodes. While consensus start, use contribution meachism to update contirbution value and it's level, take measures to reward and punish node. 

## DCML

Main point:

When bad nodes have taken down, the blank position should affect consensus process, so we use some methods such as candidate monitor and machine learning algorithms to detect latent bad nodes and candidate proper nodes into consensus process to consensus well.

## Update Log

- v0510: Add and upgrade comparison algorithm.
- v0509: Finished algorithm MGM but we got nothing about the data. 
- v0508: Implement part of alternative strategies and LOF algorithm in DCML.
- v0430: Messed up, more bugs on the way. >_<
- v0429: Fix and make some bugs.
- v0428: Modelize the code, fix some issues but not all.
- v0425: Finished basic function of IDTM.

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

## Credit 

Big thanks to GitHub with open source code.

ps: Freshbird write the F**king code, you can take a look then leave, good for you.
