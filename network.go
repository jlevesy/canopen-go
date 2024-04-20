package canopen

import (
	"context"
	"errors"
	"net"
	"sync/atomic"

	"go.einride.tech/can/pkg/socketcan"
)

var (
	errReceiverNotAdressable = errors.New("receiver must implement one of [COBIDReceiver, FunctionReceiver NodeIDReceiver])")
	errRegisterWhileRunning  = errors.New("unable to register a receiver when the network is running")
)

type Network struct {
	tx *socketcan.Transmitter
	rx *socketcan.Receiver

	receivers map[FrameMatcher][]Receiver

	running atomic.Bool
}

func NewNetwork(conn net.Conn) *Network {
	cl := Network{
		tx:        socketcan.NewTransmitter(conn),
		rx:        socketcan.NewReceiver(conn),
		receivers: make(map[FrameMatcher][]Receiver),
	}

	return &cl
}

func (b *Network) RegisterReceiver(rec Receiver) error {
	if b.running.Load() {
		return errRegisterWhileRunning
	}

	var (
		matcher = rec.Matches()
		recs    = b.receivers[matcher]
	)

	b.receivers[matcher] = append(recs, rec)

	return nil
}

func (b *Network) Run() {
	// Note: no proper termination here as receive will exit if the conn is not readable anymore...
	// That being said, error feedback might not be ideal?
	go b.receive()
}

func (b *Network) Transmit(ctx context.Context, fr Frame) error {
	return b.tx.TransmitFrame(ctx, fr.ToCANFrame())
}

func (b *Network) receive() {
	b.running.Store(true)

	for b.rx.Receive() {
		fr := NewFrame(b.rx.Frame())

		for _, matcher := range expandPotentialMatchers(fr) {
			for _, r := range b.receivers[matcher] {
				go r.Receive(fr)
			}
		}
	}
}

func expandPotentialMatchers(fr Frame) []FrameMatcher {
	return []FrameMatcher{
		MatchCOBID(fr.ID),
		MatchFunctionCode(fr.ID.FunctionCode),
		MatchNodeID(fr.ID.NodeID),
		MatchAllFrames,
	}
}
