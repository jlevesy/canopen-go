package canopen_test

import (
	"testing"

	"github.com/jlevesy/canopen-go"
	"github.com/stretchr/testify/assert"
	"go.einride.tech/can"
)

func TestNewFrame(t *testing.T) {
	for _, testCase := range []struct {
		input can.Frame
		want  canopen.Frame
	}{
		{
			input: can.Frame{
				ID:         0b01001101010,
				Length:     2,
				Data:       can.Data{1, 2, 0, 0, 0, 0, 0},
				IsExtended: true,
				IsRemote:   false,
			},
			want: canopen.Frame{
				ID: canopen.COBID{
					FunctionCode: 0b0100,
					NodeID:       0b1101010,
				},
				Length:     2,
				Data:       can.Data{1, 2, 0, 0, 0, 0, 0},
				IsExtended: true,
				IsRemote:   false,
			},
		},
	} {
		t.Run("", func(t *testing.T) {
			fr := canopen.NewFrame(testCase.input)
			assert.Equal(t, testCase.want, fr)
			assert.Equal(t, testCase.input, fr.ToCANFrame())
		})
	}
}
