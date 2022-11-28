package main

import (
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

var isFirstPrimary bool // = flag.Bool("first", false, "is this the first primary?")

func main() {
	//flag.Parse()
	
	f, err := os.OpenFile("output.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()
    log.SetOutput(f)

	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 5000

	if len(os.Args) == 3 && os.Args[2] == "first" {
		isFirstPrimary = true
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:               ownPort,
		previousAuctions: make([]auction, 0),
		peers:            make(map[int32]skrr.AuctionClient),
		ctx:              ctx,
		isPrimary:        isFirstPrimary,
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

	if p.isPrimary {
		go heartBeat(p)
		p.currentPrimary = p.id
	} else {
		p.lastPing = time.Now()
		go ListenForHeartBeat(p)
	}

	for {
		if !p.isPrimary {
			time.Sleep(5 * time.Second)
			log.Printf("Auction ongoing: %t. highest bid: %d, with id: %d",
				!p.isOver, p.highestBid, p.highestBidderID)
			continue
		}
		log.Printf("Started a new aution!\n")
		fmt.Println("Started a new aution!")

		// Reset the current auction for the primary and all backups!
		p.bidders = make([]int32, 0)
		p.highestBid = 0
		p.highestBidderID = 0
		p.isOver = false
		p.TalkToBoisAndRemoveBadBois(&skrr.BidMessage{
			Amount:   0,
			BidderID: 0,			
		}, p.isOver)
		
		p.auctionStart = time.Now()
		time.Sleep(20 * time.Second)

		p.isOver = true

		p.previousAuctions = append(p.previousAuctions, auction{
			highestBid: p.highestBid,
			bidderID:   p.highestBidderID,
			isOver:     p.isOver,
		})

		// Inform backup replicas that the auction is over.
		// These are the final values!
		p.TalkToBoisAndRemoveBadBois(&skrr.BidMessage{
			Amount:   p.highestBid,
			BidderID: p.highestBidderID,
		}, p.isOver)

		log.Printf("Aution finished. Bidder%d won with a Highest bid: %d", p.highestBidderID, p.highestBid)

		log.Printf("Starting new auction in %d sec", 10)

		time.Sleep(10 * time.Second)
	}
}

type auction struct {
	highestBid int32
	bidderID   int32
	isOver     bool
}

type peer struct {
	skrr.UnimplementedAuctionServer
	id               int32
	highestBid       int32
	highestBidderID  int32
	isOver           bool
	previousAuctions []auction
	isPrimary        bool
	peers            map[int32]skrr.AuctionClient
	bidders          []int32
	ctx              context.Context
	lastPing         time.Time

	auctionStart 	 time.Time
	currentPrimary 	 int32
}

func heartBeat(p *peer) {
	for {
		time.Sleep(5 * time.Second)
		for _, peer := range p.peers {
			peer.Ping(context.Background(), &skrr.PingMessage{
				PingerID: int32(p.id),
			})
		}
		fmt.Println("Sent heartbeat")
	}
}

func (p *peer) Bid(ctx context.Context, in *skrr.BidMessage) (*skrr.Ack, error) {
	log.Printf("Received bid from %d, bid: %d", in.BidderID, in.Amount)

	if !contains(p.bidders, in.BidderID) {
		p.bidders = append(p.bidders, in.BidderID)
	}
	if p.highestBid < int32(in.Amount) && p.TalkToBoisAndRemoveBadBois(in, p.isOver) {
		p.highestBid = int32(in.Amount)
		p.highestBidderID = int32(in.BidderID)

		return &skrr.Ack{
			Output: "successful",
		}, nil
	} else if p.highestBid >= int32(in.Amount) {
		return &skrr.Ack{
			Output: "failure",
		},nil
	} else {
		return &skrr.Ack{
			Output: "exception",
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

func (p *peer) TalkToBoisAndRemoveBadBois(bm *skrr.BidMessage, isOver bool) bool {
	errs := make([]int32, 0)
	for port, peer := range p.peers {
		//timeoutctx, cancel := context.WithTimeout(p.ctx, 5*time.Second)
		//defer cancel()
		_, err := peer.BackUp(p.ctx, &skrr.BackUpMessage{
			HighestBid:      bm.Amount,
			Bidders:         p.bidders,
			HighestBidderID: bm.BidderID,
			IsOver:          isOver,
			PrimaryId: 		 p.id,
		})
		if err != nil {
			errs = append(errs, port)
		}
	}

	if len(errs) > 0 {
		for _, v := range errs {
			delete(p.peers, v)
			log.Printf("Removed Inactive backup with id: %d\n", v)
		}
		return false
	}
	return true
}

func (p *peer) Ping(ctx context.Context, in *skrr.PingMessage) (*skrr.Ack, error) {
	p.lastPing = time.Now()

	// if the first primary isn't set to serverid = 0, meaning it is not on port 5000.
	delete(p.peers, in.PingerID)

	return &skrr.Ack{Output: "big success"}, nil
}

func ListenForHeartBeat(p *peer) {
	for {
		time.Sleep(5 * time.Second)
		fmt.Printf("Heartbeat received: %s\n", p.lastPing.Local().String())
		if time.Since(p.lastPing) > 8*time.Second {
			log.Println("Primary died :( Time for a bully sesh")
			fmt.Println("Bully time")
			p.BullyTheBois()
		}

		if p.isPrimary {
			return
		}
	}
}

func (p *peer) BullyTheBois() {
	IAmTheBiggestBoss := true
	var biggestBoss int32

	for scrubId, _ := range p.peers {
		if scrubId > biggestBoss{
			biggestBoss = scrubId
		}
		if scrubId > p.id {
			IAmTheBiggestBoss = false
		}
	}

	p.currentPrimary = biggestBoss

	if IAmTheBiggestBoss {
		p.TalkToBoisAndRemoveBadBois(&skrr.BidMessage{
			Amount:   p.highestBid,
			BidderID: p.highestBidderID,
		}, p.isOver)
		p.isPrimary = true
		go heartBeat(p)
		log.Printf("This: Server %d is now in control! and its heart is beatin!", p.id)
		fmt.Printf("This: Server %d is now in control! and its heart is beatin!\n", p.id)
	}
}

func (p *peer) Result(ctx context.Context, in *skrr.Void) (*skrr.Outcome, error) {
	var timeLeft string

	if !p.isOver {
		// calculating time left of the current ongoing auction.
		timeLeft = time.Until(p.auctionStart.Add(20 * time.Second)).String()
	} else {
		// calculating time until the next auction
		timeLeft = time.Until(p.auctionStart.Add(30 * time.Second)).String()
	}
	
	return &skrr.Outcome{
		HighestBid: p.highestBid, 
		BidderId: p.highestBidderID, 
		IsOver: p.isOver,
		TimeLeft: timeLeft }, nil
}

func (p *peer) BackUp(ctx context.Context, in *skrr.BackUpMessage) (*skrr.Ack, error) {
	if !contains(p.bidders, in.HighestBidderID) {
		p.bidders = append(p.bidders, in.HighestBidderID)
	}
	fmt.Printf("Backupmsg: %v\n", in)
	if p.highestBid < int32(in.HighestBid) || in.IsOver || (!in.IsOver && p.isOver) {
		p.highestBid = int32(in.HighestBid)
		p.highestBidderID = in.HighestBidderID
		p.isOver = in.IsOver
	}
	return &skrr.Ack{Output: "succes"}, nil
}

func (p *peer) GetPrimary(ctx context.Context, in *skrr.Void) (*skrr.ElectionResultMessage, error) {
	return &skrr.ElectionResultMessage{
		PrimaryServerPort: p.currentPrimary,
		}, nil
}

func (p *peer) CheckHeartbeat(in *skrr.Void, stream skrr.Auction_CheckHeartbeatServer) error {
	for {
		stream.Send(&skrr.PingMessage{PingerID: p.id})
		time.Sleep(1 * time.Second)
	}
}
