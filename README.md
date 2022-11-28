# Dsys - Assignment 5 - Distributed Auction

## How to start

1. Run the passive replication by starting replica nodes using ```go run . x```
   1. Note: this implementation doesn't use flags. The first argument following the ., is the id of the replica. We have hardcoded the system to work for 3 replicas.
      1. ```go run . 0 first```
      2. ```go run . 1```
      3. `go run . 2`
      4. We have hardcoded the system to work for 3 replicas. So, the only 3 ids that will work are 0,1,2.
   2. One node has to be marked as the first primary replica. This is marked as an argument following the *x* using a ```first```. This should always be the one with the lowest id, meaning id of 0.
      1. ```go run . x first```
2. Run auction clients using `go run Client/client.go`
   1. Important: Remember using the `-cID` flag to give a client node an ID.
      1. This ID is used inside the auction system, so remember to increment it for each new client started.

### Copy paste to separate terminals

First the replication nodes:

1. `go run . 0 first`
2. `go run . 1`
3. `go run . 2`

Then the client(s)
1.`go run Client/client.go -cID 1`

2.`go run Client/client.go -cID 2`

## How to use (Client)
There is a single auction running at a time. It has duration of 20sec, if you wish to bid write in the console: `bid [amount]` the bid amount most be higher than the current highest bid to be accepted.

To see the result of the latest auction write `result` in the console, this will also tell you the time left of the auction or when the next auction will start.