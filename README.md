# p2p-demo

This repo has a demo of the basics parts pf the discovery portion for the Ethereum devp2p networking protocol.

Code coverage average is at 91.9% (100, 89.2, 90.5).

## P2P test

To fulfill the requirements, an automated test of a setup with a bootstrap node and 2 regular nodes.

## Commands

This repo has 2 command programs: `bootnode` and `node`.

### bootnode

Runs a p2p bootnode that only has the UDP port for the node discovery protocol.

### node

Runs a full p2p node, both the UDP port for the discovery protocol and the TCP port for the RLPx protocol. It prints out `PeerEvent`s such as `add` and `drop`.

However, none of this is implemented RLPx messagig is implemented.

