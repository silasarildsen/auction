package main

import (
	"bufio"
	"context"
	"flag"
	"math/rand"

	"os"
	"regexp"
	skrr "skrr/grpc"
	"strconv"
	"strings"
	"time"

	"fmt"
	"log"

	grpc "google.golang.org/grpc"
)

var clientPort = flag.Int("ownPort", 5005, "brr")
var clientId = flag.Int("cID", 1, "id of the client, used for auctions")
var serverPort = flag.Int("sPort", 5000, "the server port")

var server skrr.AuctionClient
var serverConn *grpc.ClientConn
var front *Frontend

type Frontend struct {
	replicas map[int32]skrr.AuctionClient
}

func main() {
	flag.Parse()

	front = &Frontend{
		replicas: make(map[int32]skrr.AuctionClient),
	}

	// connect our "frontend" to the replicas
	for i := 5000; i <= 5002; i++ {
		s, sConn := ConnectToServer(i)
		if i == *serverPort {
			server = s
			serverConn = sConn
		}
		front.replicas[int32(i)] = s
		fmt.Println(s.GetPrimary(context.Background(), &skrr.Void{}))
	}
	go ListenToServerHeartBeat()

	listen()
}

func ListenToServerHeartBeat() {
	stream, _ := server.CheckHeartbeat(context.Background(), &skrr.Void{})
	lastBeat := time.Now()

	for {
		_, err := stream.Recv()

		//log.Println("received heartbeat")

		if time.Since(lastBeat) > 11*time.Second {
			res, _ := front.replicas[5000 + rand.Int31n(2)].GetPrimary(context.Background(), &skrr.Void{})
			fmt.Printf("New Primary port: %d", res.PrimaryServerPort)
			server = front.replicas[int32(res.PrimaryServerPort)]
			stream, _ = server.CheckHeartbeat(context.Background(), &skrr.Void{})
		}
		if err == nil {
			lastBeat = time.Now()
		}
	}
}

func listen() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		if input == "result" {
			result()
		} else if ok, _ := regexp.MatchString("bid \\d+", input); ok {
			sesses := strings.Split(input, " ")
			bidAmount, _ := strconv.Atoi(sesses[1])
			bid(bidAmount)
		}
	}
}

func bid(bidAmount int) {
	res, _ := server.Bid(context.Background(), &skrr.BidMessage{
		Amount:   int32(bidAmount),
		BidderID: int32(*clientId),
	})
	if res.Output == "successful" {
		log.Printf("Client %d, succesfully made a new highest bid", *clientId)
		fmt.Printf("Client %d, Congratz you are now the highest bidder\r\n", *clientId)
	} else if res.Output == "failure" {
		log.Printf("Client %d bid too low", *clientId)
		fmt.Printf("Your bid was too low\n", *clientId)
	}else if res.Output == "exception" {
		log.Printf("Client %d had something go wrong", *clientId)
		fmt.Printf("Client %d, Something went wrong\r\n", *clientId)
	}

}

func result() {
	result, err := server.Result(context.Background(), &skrr.Void{})
	if err != nil {
		log.Fatal("Died during result. 💀")
	}
	if result.IsOver {
		log.Printf("Auction has ended with result: %d by BidderID: %d\r\n", result.HighestBid, result.BidderId)
		fmt.Printf("Auction has ended with result: %d by BidderID: %d\r\n", result.HighestBid, result.BidderId)
	} else {
		log.Printf("Auction is ongoing with current highest bid: %d by BidderID: %d. Time left: %s\r\n", result.HighestBid, result.BidderId, result.TimeLeft)
		fmt.Printf("Auction is ongoing with current highest bid: %d by BidderID: %d\r\n", result.HighestBid, result.BidderId)
	}
}

func ConnectToServer(serverPortP int) (skrr.AuctionClient, *grpc.ClientConn) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithInsecure())

	//dial the server, with the flag "server", to get a connection to it
	log.Printf("client %d: Attempts to dial on port %d\n", *clientId, serverPortP)
	conn, err := grpc.Dial(fmt.Sprintf(":%d", serverPortP), opts...)
	if err != nil {
		log.Printf("Fail to Dial : %v", err)
		log.Fatal()
	}

	// makes a client from the server connection and saves the connection
	// and prints rather or not the connection was is READY
	serverLocalVar := skrr.NewAuctionClient(conn)
	serverConnLocalVar := conn
	log.Printf("the connection to %d is: %s", serverPortP, conn.GetState().String())

	return serverLocalVar, serverConnLocalVar
}
