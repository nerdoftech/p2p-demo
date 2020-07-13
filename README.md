# p2p-demo

This repo has a demo of the basic parts of the neighbor discovery potocol for the Ethereum devp2p networking.

Code coverage average for the files in pkg is 93.26% (100, 89.3, 90.5).

## P2P test

To fulfill the requirements, an automated test of a setup with a bootstrap node and 2 regular nodes. It does not use the commands in this project, but takes the code packages from this repo that are wrappers for the go-ethereum p2p packages. It starts up a boot node, then starts up two p2p nodes and waits for an `add` `PeerEvent` and runs assertions to ensure that the IDs match.

From this repo just run the following command:

```text
go test -v test/p2p_suite_test.go
=== RUN   TestTest
Running Suite: P2P Suite
========================
Random Seed: 1594613433
Will run 1 of 1 specs

4:10AM INF allowing bootnode to start pkg=test url=enode://0d2bbd8fcd655b7b165500464f7c9dbaaf2b66dc693a0829f842e6f33d041aaebe76965477ff3ae3f4ff658f3cac6ba35dcdc76a5c6c88bd0b52bd0287956868@127.0.0.1:0?discport=40677
4:10AM INF New local node record ctx=["seq",1,"id","c8fdb6baef45727b64359ff269c97d7c0b46d14135bfa0ce88c880f0eb0d0ca2","ip","127.0.0.1","udp",47515,"tcp",47515] node_name=node1 pkg=p2p-server
4:10AM INF Started P2P networking ctx=["self","enode://66de5bdf5a5f5912786e99acac057e0dac72242cced0bd5ce81b9e6ac971a637e7d66c8b651c7a4bc16b6fd31ba42b51e54745b311cdb8fab94c63975ae77660@127.0.0.1:47515"] node_name=node1 pkg=p2p-server
4:10AM INF New local node record ctx=["seq",1,"id","0c5c9219f79fbdb8765d095d21745d293581ecbb28d49ada307fe908d31f78bf","ip","127.0.0.1","udp",36074,"tcp",36074] node_name=node2 pkg=p2p-server
4:10AM INF Started P2P networking ctx=["self","enode://d9578d7cb1057b38d9f78de0755c627ecc724c7de5bb77671e9e37e36e6bb88091f848b59d247906c7e6013c4269dbaf555f6a3ce6148e69d5d174806332615a@127.0.0.1:36074"] node_name=node2 pkg=p2p-server
• [SLOW TEST:7.021 seconds]
test p2p node interaction
/home/eric/code/go/src/github.com/nerdoftech/p2p-demo/test/p2p_suite_test.go:40
  run nodes
  /home/eric/code/go/src/github.com/nerdoftech/p2p-demo/test/p2p_suite_test.go:68
    should work
    /home/eric/code/go/src/github.com/nerdoftech/p2p-demo/test/p2p_suite_test.go:70
------------------------------

Ran 1 of 1 Specs in 17.026 seconds
SUCCESS! -- 1 Passed | 0 Failed | 0 Pending | 0 Skipped
--- PASS: TestTest (17.03s)
PASS
ok    command-line-arguments  17.074s
```

__Note:__ 

For reasons I have not yet determined, sometimes the peer connections can take a long time. Since this protocol is brand new to me, I have not been able to fully understand why this is. I have noticed that unless the bootnode is running for at least 5-6 seconds before the nodes connect, no peering will happen at all. Then, sometimes full peer meshing can be as fast as 2 seconds and other times take upwards of 30 seconds. 

Looking through the logs at the `trace` level (example of the test running at trace level in file `trace.txt`), there is quite a bit of "Findnode failed" messages which then lead to the node logging "Too many findnode failures, dropping". Following this, the node starts attempting the bonding process all over again (ping/pong). I notice that when running the commands and the boot node and `node1` are already running, that bringing up `node2` happens much faster.

I will continue to look at it to see if I can figure this out.

## Commands

This repo has 2 command programs: `bootnode` and `node`.

### Command `bootnode`

Runs a p2p bootnode that only has the UDP port for the node discovery protocol.

### Command `node`

Runs a full p2p node, both the UDP port for the discovery protocol and the TCP port for the RLPx protocol. It prints out `PeerEvent`s such as `add` and `drop`.

However, using any of the RLPx protocol features was not implemented.

### Command operation

Running the commands manually will log output that shows the nodes discovering each other via the bootnode. For the sake of brevity, the logs are run at info level.

First, we start the bootnode:

```text
go run bootnode.go
ID: 42e1758985f74310bb40e4bcf67cc6001c3ca774e23cad92afc19d89f9b7e641
URL: enode://1306201778f9bd38f7b49115e242347770059f8c0ca874079399bb122708cd1e76c30dd6bbc2f0fabc7b5adc08d832286995454aafc0fe242d3c9b11c739d4b7@127.0.0.1:0?discport=30303
```

Next, we start `node1` pointing at our boot node it waits for peer (note the node ID). We can also see that it has 0 for the `peercount`:

```text
go run node.go -random -name node1  -bootnode enode://1306201778f9bd38f7b49115e242347770059f8c0ca874079399bb122708cd1e76c30dd6bbc2f0fabc7b5adc08d832286995454aafc0fe242d3c9b11c739d4b7@127.0.0.1:0?discport=30303
10:41PM INF New local node record ctx=["seq",1,"id","5944a692b630d7a3c4ec10f2cc62e93246f3b47e492891ef4f51f66dd9888ef9","ip","127.0.0.1","udp",39267,"tcp",39267] node_name=node1 pkg=p2p-server
10:41PM INF Started P2P networking ctx=["self","enode://bea605a66249a2729c185f85a0dbd1fa5a17feb5cc52d8402d5c0049d48b189db30ea155289af3cf539ed548c20a868d8202624633214e52f115bb743e796835@127.0.0.1:39267"] node_name=node1 pkg=p2p-server
10:41PM INF Looking for peers ctx=["peercount",0,"tried",0,"static",0] node_name=node1 pkg=p2p-server
```

Then we start `node2` and it discovers `node1` quickly as shown by the node ID in the `add` PeerEvent and subsequently we see that the `peercount` goes to 1:

```text
go run node.go -random -name node2  -bootnode enode://1306201778f9bd38f7b49115e242347770059f8c0ca874079399bb122708cd1e76c30dd6bbc2f0fabc7b5adc08d832286995454aafc0fe242d3c9b11c739d4b7@127.0.0.1:0?discport=30303
10:42PM INF New local node record ctx=["seq",1,"id","be8a7c277ffc9100ca71d851fcf94ed95bb6421b84c968b4d58c1db949a6c312","ip","127.0.0.1","udp",43990,"tcp",43990] node_name=node2 pkg=p2p-server
10:42PM INF Started P2P networking ctx=["self","enode://fa44ba9bab662ce427b22dae6ac1e76a7fa633aa20c85b99de457cae56169fb1509be8e9e89cfe73311938fb06e7ae12996dca7abe884e896cea505ddba5ac85@127.0.0.1:43990"] node_name=node2 pkg=p2p-server
10:42PM INF received new event local_addr=127.0.0.1:43990 peer=5944a692b630d7a3c4ec10f2cc62e93246f3b47e492891ef4f51f66dd9888ef9 remote_addr=127.0.0.1:42130 type=add
10:42PM INF Looking for peers ctx=["peercount",1,"tried",0,"static",0] node_name=node2 pkg=p2p-server
```

Looking back at the logs of `node1`, we can see that it likewise added `node2` as indicated by the node ID and the `peercount` is now 1:

```text
10:42PM INF received new event local_addr=127.0.0.1:42130 peer=be8a7c277ffc9100ca71d851fcf94ed95bb6421b84c968b4d58c1db949a6c312 remote_addr=127.0.0.1:43990 type=add
10:42PM INF Looking for peers ctx=["peercount",1,"tried",1,"static",0] node_name=node1 pkg=p2p-server
```

We kill the `node2` process and then `node1` shows it getting removed as a peer:

```text
10:43PM INF received new event error="client quitting" local_addr=127.0.0.1:42130 peer=be8a7c277ffc9100ca71d851fcf94ed95bb6421b84c968b4d58c1db949a6c312 remote_addr=127.0.0.1:43990 type=drop
10:43PM INF Looking for peers ctx=["peercount",0,"tried",1,"static",0] node_name=node1 pkg=p2p-server

```

## Deeper look at the protocol

Running `node2` at trace level affords us the opportunity to really see the devp2p protocol in action. We will truncate many of the log messages since it is extremely verbose.

We start `node2`:

```text
go run node.go -random -log trace -name node2  -bootnode enode://1306201778f9bd38f7b49115e242347770059f8c0ca874079399bb122708cd1e76c30dd6bbc2f0fabc7b5adc08d832286995454aafc0fe242d3c9b11c739d4b7@127.0.0.1:0?discport=30303
11:31PM DBG setting allowed nets allowed nets=[{"IP":"127.0.0.0","Mask":"/wAAAA=="}] pkg=node
11:31PM DBG Parsed boot node ID=42e1758985f74310bb40e4bcf67cc6001c3ca774e23cad92afc19d89f9b7e641 pkg=node url=enode://1306201778f9bd38f7b49115e242347770059f8c0ca874079399bb122708cd1e76c30dd6bbc2f0fabc7b5adc08d832286995454aafc0fe242d3c9b11c739d4b7@127.0.0.1:0?discport=30303
{"level":"debug","pkg":"util","time":"2020-07-12T23:31:01-05:00","message":"generating node key"}
{"level":"debug","pkg":"util","time":"2020-07-12T23:31:01-05:00","message":"key generated"}
11:31PM DBG starting server
11:31PM DBG UDP listener up ctx=["addr",{"IP":"127.0.0.1","Port":57243,"Zone":""}] node_name=node2 pkg=p2p-server
11:31PM DBG TCP listener up ctx=["addr",{"IP":"127.0.0.1","Port":57243,"Zone":""}] node_name=node2 pkg=p2p-server
11:31PM TRC Found seed node in database ctx="marshaling error: json: unsupported type: func() interface {}" node_name=node2 pkg=p2p-server
11:31PM INF New local node record ctx=["seq",1,"id","31d5a6e18619ca5c78a6b7ce458f38bbac9925d00394b4544c1e493e5553d1f4","ip","127.0.0.1","udp",57243,"tcp",57243] node_name=node2 pkg=p2p-server
11:31PM DBG server started, waiting for signal to shutdown
11:31PM INF Started P2P networking ctx=["self","enode://577ef7c67120c22c4be2776d73cbb7ca0afd1f92ce55cdcf2799c2852452a060a335963195c9d528d9c773ff073624879b98847dc9bc40a4932b5401fb9fb44d@127.0.0.1:57243"] node_name=node2 pkg=p2p-server

```

Remember that from starting the boot node earlier, we have a bootnode ID of `42e1758985f74310bb40e4bcf67cc6001c3ca774e23cad92afc19d89f9b7e641`

We can see the node trying to bond (PING/PONG) with the boot node. This bonding process must happen before we can send a `FINDNODE` packet as this was designed to mitigate DDoS attacks. It does this about 3-4 times before sending a `FINDNODE`:

```text
11:31PM TRC >> PING/v4 ctx=["id","42e1758985f74310bb40e4bcf67cc6001c3ca774e23cad92afc19d89f9b7e641","addr",{"IP":"127.0.0.1","Port":30303,"Zone":""},"err",null] node_name=node2 pkg=p2p-server
11:31PM TRC << PONG/v4 ctx=["id","42e1758985f74310bb40e4bcf67cc6001c3ca774e23cad92afc19d89f9b7e641","addr",{"IP":"127.0.0.1","Port":30303,"Zone":""},"err",null] node_name=node2 pkg=p2p-server
```

From what we can see of these next few log entries, it appears to log the id of the node it is sending to, not the payload of the packet. I am making this assumption because according to the spec, the first `FINDNODE` to the boot node is supposed to be for self, as in asking the boot node for neighbors close in distance (XOR of node IDs) to itself.

```text
11:31PM TRC >> FINDNODE/v4 ctx=["id","42e1758985f74310bb40e4bcf67cc6001c3ca774e23cad92afc19d89f9b7e641","addr",{"IP":"127.0.0.1","Port":30303,"Zone":""},"err",null] node_name=node2 pkg=p2p-server
11:31PM TRC >> FINDNODE/v4 ctx=["id","42e1758985f74310bb40e4bcf67cc6001c3ca774e23cad92afc19d89f9b7e641","addr",{"IP":"127.0.0.1","Port":30303,"Zone":""},"err",null] node_name=node2 pkg=p2p-server
11:31PM TRC << NEIGHBORS/v4 ctx=["id","42e1758985f74310bb40e4bcf67cc6001c3ca774e23cad92afc19d89f9b7e641","addr",{"IP":"127.0.0.1","Port":30303,"Zone":""},"err",null] node_name=node2 pkg=p2p-server
11:31PM TRC << NEIGHBORS/v4 ctx=["id","42e1758985f74310bb40e4bcf67cc6001c3ca774e23cad92afc19d89f9b7e641","addr",{"IP":"127.0.0.1","Port":30303,"Zone":""},"err",null] node_name=node2 pkg=p2p-server
```

This goes on quite a few times with other messages and some of the failure messages mentioned earlier. Eventually, you see `node2` attempt to reach out to `node1` and form a peer relationship:

```text
11:31PM TRC Starting p2p dial ctx=["id","5944a692b630d7a3c4ec10f2cc62e93246f3b47e492891ef4f51f66dd9888ef9","ip","127.0.0.1","flag",1] node_name=node2 pkg=p2p-server
11:31PM TRC Discarding dial candidate ctx=["id","31d5a6e18619ca5c78a6b7ce458f38bbac9925d00394b4544c1e493e5553d1f4","ip","127.0.0.1","reason",{}] node_name=node2 pkg=p2p-server
11:31PM TRC >> PING/v4 ctx=["id","5944a692b630d7a3c4ec10f2cc62e93246f3b47e492891ef4f51f66dd9888ef9","addr",{"IP":"127.0.0.1","Port":39267,"Zone":""},"err",null] node_name=node2 pkg=p2p-server
11:31PM TRC << PONG/v4 ctx=["id","5944a692b630d7a3c4ec10f2cc62e93246f3b47e492891ef4f51f66dd9888ef9","addr",{"IP":"127.0.0.1","Port":39267,"Zone":""},"err",null] node_name=node2 pkg=p2p-server
11:31PM TRC << PING/v4 ctx=["id","5944a692b630d7a3c4ec10f2cc62e93246f3b47e492891ef4f51f66dd9888ef9","addr",{"IP":"127.0.0.1","Port":39267,"Zone":""},"err",null] node_name=node2 pkg=p2p-server
11:31PM TRC >> PONG/v4 ctx=["id","5944a692b630d7a3c4ec10f2cc62e93246f3b47e492891ef4f51f66dd9888ef9","addr",{"IP":"127.0.0.1","Port":39267,"Zone":""},"err",null] node_name=node2 pkg=p2p-server
11:31PM INF received new event local_addr=127.0.0.1:58798 peer=5944a692b630d7a3c4ec10f2cc62e93246f3b47e492891ef4f51f66dd9888ef9 remote_addr=127.0.0.1:39267 type=add
11:31PM DBG Adding p2p peer ctx=["peercount",1,"id","5944a692b630d7a3c4ec10f2cc62e93246f3b47e492891ef4f51f66dd9888ef9","conn",1,"addr",{"IP":"127.0.0.1","Port":39267,"Zone":""},"name","node1"] node_name=node2 pkg=p2p-server
11:31PM TRC << FINDNODE/v4 ctx=["id","5944a692b630d7a3c4ec10f2cc62e93246f3b47e492891ef4f51f66dd9888ef9","addr",{"IP":"127.0.0.1","Port":39267,"Zone":""},"err",null] node_name=node2 pkg=p2p-server
11:31PM TRC >> NEIGHBORS/v4 ctx=["id","5944a692b630d7a3c4ec10f2cc62e93246f3b47e492891ef4f51f66dd9888ef9","addr",{"IP":"127.0.0.1","Port":39267,"Zone":""},"err",null] node_name=node2 pkg=p2p-server
....
11:31PM INF Looking for peers ctx=["peercount",1,"tried",1,"static",0] node_name=node2 pkg=p2p-server
```

## Other thoughts

It appears there are wireshark dissectors for devp2p: https://github.com/ConsenSys/ethereum-dissectors

It would have been really insightful to watch this in wireshark, but it required building wireshark from source, so I did not pursue it in the interest of time.
