package nmt

import (
	"context"
	"errors"

	"github.com/jlevesy/canopen-go"
	"go.einride.tech/can"
)

var nmtCOBID = canopen.COBID{}

var errMalformedFrame = errors.New("malformed frame")

const AllNodes = 0

type State uint8

const (
	StateOperational       State = 0x01
	StateStopped           State = 0x02
	StatePreOperational    State = 0x80
	StateResetNode         State = 0x81
	StateResetComunication State = 0x82
)

type StateManager struct {
	tx canopen.Transmitter
}

func NewStateManager(tx canopen.Transmitter) *StateManager {
	return &StateManager{tx: tx}
}

func (nm *StateManager) SetNetworkState(ctx context.Context, nodeID uint8, st State) error {
	return nm.tx.Transmit(
		ctx,
		canopen.Frame{
			ID:     nmtCOBID,
			Length: 2,
			Data:   can.Data{byte(nodeID), byte(st)},
		},
	)
}

type StateChangeHandler interface {
	OnNetworkStateChanged(nodeID uint8, st State)
}

type StateChangeHandlerFunc func(uint8, State)

func (f StateChangeHandlerFunc) OnNetworkStateChanged(nodeID uint8, st State) {
	f(nodeID, st)
}

type StateReceiver struct {
	handler StateChangeHandler
}

func NewStateReceiver(hdlr StateChangeHandler) *StateReceiver {
	return &StateReceiver{handler: hdlr}
}

func (n *StateReceiver) Matches() canopen.FrameMatcher {
	return canopen.MatchCOBID(nmtCOBID)
}

func (n *StateReceiver) Receive(frame canopen.Frame) {
	if frame.Length != 2 {
		return
	}

	n.handler.OnNetworkStateChanged(uint8(frame.Data[0]), State(frame.Data[1]))
}
