package canopen

import (
	"go.einride.tech/can"
)

type COBID struct {
	NodeID       uint8
	FunctionCode uint8
}

type Frame struct {
	ID         COBID
	Length     uint8
	Data       can.Data
	IsRemote   bool
	IsExtended bool
}

func NewFrame(fr can.Frame) Frame {
	return Frame{
		ID:         parseCOBID(fr.ID),
		Length:     fr.Length,
		Data:       fr.Data,
		IsRemote:   fr.IsRemote,
		IsExtended: fr.IsExtended,
	}
}

func (f *Frame) ToCANFrame() can.Frame {
	return can.Frame{
		ID:         encodeCOBID(f.ID),
		Length:     f.Length,
		Data:       f.Data,
		IsRemote:   f.IsRemote,
		IsExtended: f.IsExtended,
	}
}

const cobIDMask uint32 = 0b11110000000

func parseCOBID(cobID uint32) COBID {
	return COBID{
		FunctionCode: uint8((cobID & cobIDMask) >> 7),
		NodeID:       uint8(cobID & (^cobIDMask)),
	}
}

func encodeCOBID(id COBID) uint32 {
	res := uint32(id.FunctionCode)
	res = res<<7 | uint32(id.NodeID)
	return res
}
