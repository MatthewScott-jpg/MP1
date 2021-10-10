# MP1
This project was created by Matthew Scott, Liam Potts, and Roland Afzelius. It creates structs resembling nodes in a network, then
implements push, pull, and push-pull gossip to evaluate their performance.

## How to run
1. Download the repository and navigate to the folder with `main.go`, then enter `go run main.go`.

2. The total number of runs required to complete each algorithm are printed. Print statements for intermediate steps are 
commented out but can be uncommented to see what nodes are doing at each step. Each total run print takes between 5-10
seconds to be called.

## Design
Although a central repository is used, this was to simplify the reusability of functions. Each node carries its own copy
of the channels used for contacting nodes, and none know if the nodes they are contacting are infected or not, so they 
are all operating independently. The channel maps were implemented in a way that the key associated with a map is the 
name of the node that receives messages through it.
#### Push
Push protocol is implemented by having infected nodes push their message through the contactChannel map. After a brief 
delay to maintain synchronicity, all nodes check their channel for messages, and uninfected nodes update their message
if they received a new one.
#### Pull
Pull protocol is implemented by having uninfected nodes send their name through the nameChannel map to a random node. 
Nodes then check their nameChannel for requests, and if they have received any, send their message through the 
contactChannel using the name received as the key in the contactChannel map. Finally, all nodes check their contactChannel
for messages and uninfected nodes update if they receive a new message.
#### Push/Pull
Push/pull protocol runs push protocol until half of the nodes in system are infected, then switches to pull protocol.

## References
read channel with default case: https://stackoverflow.com/questions/3398490/checking-if-a-channel-has-a-ready-to-read-value-using-go