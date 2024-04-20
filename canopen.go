package canopen

import "context"

type Transmitter interface {
	Transmit(context.Context, Frame) error
}

type Receiver interface {
	Matches() FrameMatcher
	Receive(Frame)
}

type FrameMuxer interface {
	RegisterReceiver(Receiver) error
}

type FrameMatcher struct {
	id                COBID
	functionCode      uint8
	nodeID            uint8
	matchFunctionCode bool
	matchID           bool
	matchNodeID       bool
	matchAll          bool
}

func MatchCOBID(id COBID) FrameMatcher {
	return FrameMatcher{matchID: true, id: id}
}

func MatchFunctionCode(code uint8) FrameMatcher {
	return FrameMatcher{matchFunctionCode: true, functionCode: code}
}

func MatchNodeID(nodeID uint8) FrameMatcher {
	return FrameMatcher{matchNodeID: true, nodeID: nodeID}
}

var MatchAllFrames = FrameMatcher{matchAll: true}
