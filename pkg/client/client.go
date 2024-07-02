package client

import (
	"fmt"

	"github.com/elleqt/gotron-sdk/pkg/proto/api"
	"google.golang.org/grpc"
)

// Client controller structure
type Client struct {
	Address string
	Conn    *grpc.ClientConn
	Client  api.WalletClient
	opts    []grpc.DialOption
}

// New create grpc controller
func New(address string) *Client {
	client := &Client{
		Address: address,
	}
	return client
}

// Start initiate grpc  connection
func (g *Client) Start(opts ...grpc.DialOption) error {
	var err error
	if len(g.Address) == 0 {
		g.Address = "grpc.trongrid.io:50051"
	}
	g.opts = opts
	g.Conn, err = grpc.NewClient(g.Address, opts...)

	if err != nil {
		return fmt.Errorf("connecting GRPC Client: %v", err)
	}
	g.Client = api.NewWalletClient(g.Conn)
	return nil
}

// Stop GRPC Connection
func (g *Client) Stop() {
	if g.Conn != nil {
		g.Conn.Close()
	}
}

// Reconnect GRPC
func (g *Client) Reconnect(url string) error {
	g.Stop()
	if len(url) > 0 {
		g.Address = url
	}
	g.Start(g.opts...)
	return nil
}

// GetMessageBytes return grpc message from bytes
func GetMessageBytes(m []byte) *api.BytesMessage {
	message := new(api.BytesMessage)
	message.Value = m
	return message
}

// GetMessageNumber return grpc message number
func GetMessageNumber(n int64) *api.NumberMessage {
	message := new(api.NumberMessage)
	message.Num = n
	return message
}

// GetPaginatedMessage return grpc message number
func GetPaginatedMessage(offset int64, limit int64) *api.PaginatedMessage {
	return &api.PaginatedMessage{
		Offset: offset,
		Limit:  limit,
	}
}
