# KQueuey

KQueuey is a distributed message queue, with data replication for fault tolerance. Kqueuey uses the Raft Consensus Algorithm 
to allow messages to be replicated to other machines, to maintain high availability and reliability if a leader is down. To ensure
messages can still be processed. Kqueuey uses a similar queuing model as AWS SQS, with visibility timeouts, 
and ability to send and receive messages between software components. Standard queue type is uses, messages are guaranteed 
to be delivered at least once, but sometimes data may be sent again and out of order, but provides the best 
effort ordering to ensure data is close to the order it was sent. Uses BaderDB (an embedded KV Store), as the database.