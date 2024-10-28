# KQueuey

 KQueuey is a distributed message queue, with data replication for fault tolerance. Kqueuey uses the Raft Consensus Algorithm 
 to allow messages to be replicated to other machines, to maintain high availability and reliability if a leader is down. To ensure
 messages can still be processed. Kqueuey uses a similar queuing model as AWS SQS, with visibility timeouts, 
 and ability to send and receive messages between software components. Standard queue type is used. Messages are guaranteed 
 to be delivered at least once.






![Screenshot 2024-10-27 at 10 53 51 PM](https://github.com/user-attachments/assets/842e1d72-d0c6-40ac-9c0e-5d48a0287581)
