# Kqueuey

# highly distributed, simple message queue.

 KQueuey is a simple distributed message queue buily in Go, with data replication using Raft Consensus Algorithm.
 Ability to queue messages and poll messages between software components with deadlines to ensure a message is finished or re-processed.

### Overview
   1. Client distributes messages to leader node.
   2. Leader node recieves messages, and sends batch messages back to what ever client is polling messages.
   4. Clients process messages, and sends tombstones back to prevent any messages being resent, or if deadline is reached, message is re-received to any client polling messages.


![kqueuey](https://github.com/user-attachments/assets/099b646a-a415-49e9-b055-8c8ccf611d9d)
