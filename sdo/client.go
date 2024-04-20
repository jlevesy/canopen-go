package sdo

import (
	"context"
	"io"

	"github.com/jlevesy/canopen-go"
)

type Client struct {
	tx     canopen.Transmitter
	frames chan canopen.Frame
	nodeID uint8
}

func NewClient(nodeID uint8, tx canopen.Transmitter) *Client {
	return &Client{
		tx:     tx,
		frames: make(chan canopen.Frame),
		nodeID: nodeID,
	}
}

type Request struct {
	NodeID uint8

	Index    uint16
	SubIndex uint8

	PayloadLength int
	Payload       io.Reader

	BlockTransfer bool
}

type Response struct {
}

func (c *Client) Do(ctx context.Context, req Request) (*Response, error) {
	// Check if we are expedited, segmented, or block.
	// Blah blah blah form the correct frame

}
