package check_sum

import "encoding/binary"

const (
	maxByte  = 0xff // maximal byte value
)

type CheckSum struct {}

func (s CheckSum) Sum(line uint32, data []byte) byte {
	var sum byte
	var bLine = make([]byte, 4)
	binary.LittleEndian.PutUint32(bLine, line)
	for _, b := range bLine {
		sum += b
	}
	for _, b := range data {
		sum += b
	}
	return maxByte - sum
}
