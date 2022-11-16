package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	skrr "skrr/grpc"

	"google.golang.org/grpc"
)

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 5000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:               ownPort,
		previousAuctions: make([]auction, 0),
		peers:            make(map[int32]skrr.AuctionClient),
		ctx:              ctx,
	}

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	grpcServer := grpc.NewServer()
	skrr.RegisterAuctionServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i)

		if port == ownPort {
			continue
		}

		var conn *grpc.ClientConn
		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		c := skrr.NewAuctionClient(conn)
		p.peers[port] = c
	}

	for {
		time.Sleep(20 * time.Second)

		p.isOver = true

		p.previousAuctions = append(p.previousAuctions, auction{
			highestBid: p.highestBid,
			bidderID: p.highestBidderID,
			isOver: p.isOver,
		})

		log.Printf("Aution finished. Bidder%d won with a Highest bid: %d",p.highestBidderID, p.highestBid)

		log.Printf("Starting new auction in %d secc", 10 )

		time.Sleep(10 * time.Second)
		p.bidders = make([]int32, 0)
		p.highestBid = 0
		p.highestBidderID = 0
		p.isOver = false
	}
}

type auction struct {
	highestBid int32
	bidderID   int32
	isOver bool
}

type peer struct {
	skrr.UnimplementedAuctionServer
	id               int32
	highestBid       int32
	highestBidderID  int32
	isOver			 bool
	previousAuctions []auction
	isPrimary        bool
	peers            map[int32]skrr.AuctionClient
	bidders			 []int32
	ctx              context.Context
}

func (p *peer) Bid(ctx context.Context, in *skrr.BidMessage) (*skrr.Ack, error) {
	if !contains(p.bidders, in.BidderID) {
		p.bidders = append(p.bidders, in.BidderID)
	}
	if p.highestBid < int32(in.Amount) && p.TalkToBoisAndRemoveBadBois(*in) {
		p.highestBid = int32(in.Amount)
		return &skrr.Ack{
			Output: "successful",
		}, nil
	} else {
		return &skrr.Ack{
			Output: "failure",
		}, nil
	}
}

func contains(s []int32, e int32) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func (p *peer) TalkToBoisAndRemoveBadBois(BidMessage skrr.BidMessage) bool {
	errs := make([]int32, 0)
	for port, peer := range p.peers {
		timeoutctx, cancel := context.WithTimeout(p.ctx, 5*time.Second)
		defer cancel()
		_, err := peer.BackUp(timeoutctx, &skrr.BackUpMessage{
			HighestBid: p.highestBid,
			Bidders: p.bidders,
			HighestBidderID: p.highestBidderID,
		})
		if err != nil {
			errs = append(errs, port)
		}
	}
	
	if len(errs) > 0 {
		for _, v := range errs {
			delete(p.peers, v)
		}
		return false
	}
	return true
}

func (p *peer) Result(ctx context.Context, in *skrr.Void) (*skrr.Outcome, error) {
	return &skrr.Outcome{HighestBid: p.highestBid, BidderId: p.highestBidderID, IsOver: p.isOver},nil
}

func (p *peer) BackUp(ctx context.Context, in *skrr.BackUpMessage) (*skrr.Ack, error) {
	if !contains(p.bidders, in.HighestBidderID) {
		p.bidders = append(p.bidders, in.HighestBidderID)
	}
	if p.highestBid < int32(in.HighestBid) {
		p.highestBid = int32(in.HighestBid)
		p.highestBidderID = in.HighestBidderID
	}
	return &skrr.Ack{Output: "succes"},nil
}
