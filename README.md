# Kqueuey

# highly distrubuted, simple message queue.

 KQueuey is a simple distributed message queue, with data replication using Raft Consensus Algorithm.
 Ability to queue messages and poll messages between software components with message deadline.

### Overview
   1. Client distributes messages to leader node.
   2. Leader node recieves messages, and sends batch messages back to what ever worker is polling messages with deadline.
   3. Client recieves messages from leader node.
   4. Clients process messages, and sends tombstones back to prevent any messages being resent, or if deadline is reach, message is re-received.


![skqueuey](https://github.com/user-attachments/assets/efab74b1-67ac-4abe-8b21-dddf8cf2db8c)
