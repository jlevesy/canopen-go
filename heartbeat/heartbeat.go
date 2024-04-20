package heartbeat

import (
	"context"

	"github.com/jlevesy/canopen-go"
	"go.einride.tech/can"
)

type State uint8

const (
	StateInitialising   State = 0x01
	StateStopped        State = 0x04
	StateOperational    State = 0x05
	StatePreOperational State = 0x7f
)

const functionCode = 0b1110

type Heartbeat struct {
	NodeID uint8
	State  State
}

type Sender struct {
	tx canopen.Transmitter
}

func NewSender(tx canopen.Transmitter) *Sender {
	return &Sender{tx: tx}
}

func (s *Sender) SendHeartbeat(ctx context.Context, beat Heartbeat) error {
	return s.tx.Transmit(
		ctx,
		canopen.Frame{
			ID: canopen.COBID{
				FunctionCode: functionCode,
				NodeID:       beat.NodeID,
			},
			Length: 1,
			Data:   can.Data{byte(beat.State)},
		},
	)
}

type Handler interface {
	HandleHeartbeat(Heartbeat)
}

type HandlerFunc func(Heartbeat)

func (h HandlerFunc) HandleHeartbeat(beat Heartbeat) {
	h(beat)
}

type receiver struct {
	handler Handler
}

func (r *receiver) Receive(frame canopen.Frame) {
	if frame.Length != 1 {
		return
	}

	r.handler.HandleHeartbeat(Heartbeat{NodeID: frame.ID.NodeID, State: State(frame.Data[0])})
}

type GlobalReceiver struct {
	receiver
}

func NewGlobalReceiver(handler Handler) *GlobalReceiver {
	return &GlobalReceiver{receiver: receiver{handler: handler}}
}

func (g *GlobalReceiver) Matches() canopen.FrameMatcher {
	return canopen.MatchFunctionCode(functionCode)
}

type NodeReceiver struct {
	receiver

	nodeID uint8
}

func NewNodeReceiver(nodeID uint8, handler Handler) *NodeReceiver {
	return &NodeReceiver{nodeID: nodeID, receiver: receiver{handler: handler}}
}

func (n *NodeReceiver) Matches() canopen.FrameMatcher {
	return canopen.MatchCOBID(
		canopen.COBID{
			FunctionCode: functionCode,
			NodeID:       n.nodeID,
		},
	)
}
