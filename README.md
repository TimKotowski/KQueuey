# Kqueuey

# highly distrubuted, simple message queue.

 KQueuey is a simple distributed message queue, with data replication using Raft Consensus Algorithm.
 Ability to queue messages and poll messages between software components with message deadline.

### Overview
   1. Applicaiton compontents send messages to client to distribute messages to leader node.
   2. Leader node recieves messages, and sends message back what ever compontent of polling messages with deadline.
   3. Client recieve message from leader node and sends to application compontents polling messages.
   4. Applicaiton compontents process messages, send tombstone back to send messages are not resent.
   5. Applicaiton compontents process message passed deadline, message is resent to compontents.


![skqueuey](https://github.com/user-attachments/assets/efab74b1-67ac-4abe-8b21-dddf8cf2db8c)
