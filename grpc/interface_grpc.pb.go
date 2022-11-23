// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package skrr

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AuctionClient is the client API for Auction service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuctionClient interface {
	Result(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Outcome, error)
	Bid(ctx context.Context, in *BidMessage, opts ...grpc.CallOption) (*Ack, error)
	BackUp(ctx context.Context, in *BackUpMessage, opts ...grpc.CallOption) (*Ack, error)
	Ping(ctx context.Context, in *PingMessage, opts ...grpc.CallOption) (*Ack, error)
	// front-end methods
	GetPrimary(ctx context.Context, in *Void, opts ...grpc.CallOption) (*ElectionResultMessage, error)
	CheckHeartbeat(ctx context.Context, in *Void, opts ...grpc.CallOption) (Auction_CheckHeartbeatClient, error)
}

type auctionClient struct {
	cc grpc.ClientConnInterface
}

func NewAuctionClient(cc grpc.ClientConnInterface) AuctionClient {
	return &auctionClient{cc}
}

func (c *auctionClient) Result(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Outcome, error) {
	out := new(Outcome)
	err := c.cc.Invoke(ctx, "/auction.Auction/Result", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionClient) Bid(ctx context.Context, in *BidMessage, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/auction.Auction/Bid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionClient) BackUp(ctx context.Context, in *BackUpMessage, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/auction.Auction/BackUp", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionClient) Ping(ctx context.Context, in *PingMessage, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/auction.Auction/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionClient) GetPrimary(ctx context.Context, in *Void, opts ...grpc.CallOption) (*ElectionResultMessage, error) {
	out := new(ElectionResultMessage)
	err := c.cc.Invoke(ctx, "/auction.Auction/GetPrimary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionClient) CheckHeartbeat(ctx context.Context, in *Void, opts ...grpc.CallOption) (Auction_CheckHeartbeatClient, error) {
	stream, err := c.cc.NewStream(ctx, &Auction_ServiceDesc.Streams[0], "/auction.Auction/CheckHeartbeat", opts...)
	if err != nil {
		return nil, err
	}
	x := &auctionCheckHeartbeatClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Auction_CheckHeartbeatClient interface {
	Recv() (*PingMessage, error)
	grpc.ClientStream
}

type auctionCheckHeartbeatClient struct {
	grpc.ClientStream
}

func (x *auctionCheckHeartbeatClient) Recv() (*PingMessage, error) {
	m := new(PingMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AuctionServer is the server API for Auction service.
// All implementations must embed UnimplementedAuctionServer
// for forward compatibility
type AuctionServer interface {
	Result(context.Context, *Void) (*Outcome, error)
	Bid(context.Context, *BidMessage) (*Ack, error)
	BackUp(context.Context, *BackUpMessage) (*Ack, error)
	Ping(context.Context, *PingMessage) (*Ack, error)
	// front-end methods
	GetPrimary(context.Context, *Void) (*ElectionResultMessage, error)
	CheckHeartbeat(*Void, Auction_CheckHeartbeatServer) error
	mustEmbedUnimplementedAuctionServer()
}

// UnimplementedAuctionServer must be embedded to have forward compatible implementations.
type UnimplementedAuctionServer struct {
}

func (UnimplementedAuctionServer) Result(context.Context, *Void) (*Outcome, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Result not implemented")
}
func (UnimplementedAuctionServer) Bid(context.Context, *BidMessage) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bid not implemented")
}
func (UnimplementedAuctionServer) BackUp(context.Context, *BackUpMessage) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BackUp not implemented")
}
func (UnimplementedAuctionServer) Ping(context.Context, *PingMessage) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedAuctionServer) GetPrimary(context.Context, *Void) (*ElectionResultMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPrimary not implemented")
}
func (UnimplementedAuctionServer) CheckHeartbeat(*Void, Auction_CheckHeartbeatServer) error {
	return status.Errorf(codes.Unimplemented, "method CheckHeartbeat not implemented")
}
func (UnimplementedAuctionServer) mustEmbedUnimplementedAuctionServer() {}

// UnsafeAuctionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuctionServer will
// result in compilation errors.
type UnsafeAuctionServer interface {
	mustEmbedUnimplementedAuctionServer()
}

func RegisterAuctionServer(s grpc.ServiceRegistrar, srv AuctionServer) {
	s.RegisterService(&Auction_ServiceDesc, srv)
}

func _Auction_Result_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServer).Result(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auction.Auction/Result",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServer).Result(ctx, req.(*Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auction_Bid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BidMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServer).Bid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auction.Auction/Bid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServer).Bid(ctx, req.(*BidMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auction_BackUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BackUpMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServer).BackUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auction.Auction/BackUp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServer).BackUp(ctx, req.(*BackUpMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auction_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auction.Auction/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServer).Ping(ctx, req.(*PingMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auction_GetPrimary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServer).GetPrimary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auction.Auction/GetPrimary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServer).GetPrimary(ctx, req.(*Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auction_CheckHeartbeat_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Void)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AuctionServer).CheckHeartbeat(m, &auctionCheckHeartbeatServer{stream})
}

type Auction_CheckHeartbeatServer interface {
	Send(*PingMessage) error
	grpc.ServerStream
}

type auctionCheckHeartbeatServer struct {
	grpc.ServerStream
}

func (x *auctionCheckHeartbeatServer) Send(m *PingMessage) error {
	return x.ServerStream.SendMsg(m)
}

// Auction_ServiceDesc is the grpc.ServiceDesc for Auction service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Auction_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auction.Auction",
	HandlerType: (*AuctionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Result",
			Handler:    _Auction_Result_Handler,
		},
		{
			MethodName: "Bid",
			Handler:    _Auction_Bid_Handler,
		},
		{
			MethodName: "BackUp",
			Handler:    _Auction_BackUp_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _Auction_Ping_Handler,
		},
		{
			MethodName: "GetPrimary",
			Handler:    _Auction_GetPrimary_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CheckHeartbeat",
			Handler:       _Auction_CheckHeartbeat_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "grpc/interface.proto",
}
