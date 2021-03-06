package main

import (
	"bytes"
	"testing"
	"time"

	pinba "github.com/olegfedoseev/pinba-server/client"
	"github.com/stretchr/testify/assert"
)

// Valid Pinba packet with test request and couple of timers
var pinbaPacket = []byte{0xa, 0x8, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x7, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x72, 0x75, 0x1a, 0x9, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x2e, 0x70, 0x68, 0x70, 0x20, 0x0, 0x28, 0xa6, 0x2, 0x30,
	0x80, 0x80, 0x40, 0x3d, 0x7, 0x9b, 0xba, 0x3c, 0x45, 0x0, 0x0, 0x0, 0x0,
	0x4d, 0xa, 0xd7, 0x23, 0x3c, 0x50, 0x1, 0x50, 0x1, 0x5d, 0x9e, 0xd2, 0xc1,
	0x3b, 0x5d, 0x4a, 0x96, 0x13, 0x3a, 0x60, 0x3, 0x60, 0x1, 0x68, 0x4, 0x68,
	0x6, 0x68, 0x8, 0x68, 0xa, 0x70, 0x5, 0x70, 0x7, 0x70, 0x9, 0x70, 0xb, 0x7a,
	0x8, 0x72, 0x65, 0x71, 0x5f, 0x76, 0x61, 0x6c, 0x31, 0x7a, 0x8, 0x72, 0x65,
	0x71, 0x5f, 0x74, 0x61, 0x67, 0x31, 0x7a, 0x8, 0x72, 0x65, 0x71, 0x5f, 0x76,
	0x61, 0x6c, 0x32, 0x7a, 0x8, 0x72, 0x65, 0x71, 0x5f, 0x74, 0x61, 0x67, 0x32,
	0x7a, 0x4, 0x6b, 0x65, 0x79, 0x31, 0x7a, 0x4, 0x76, 0x61, 0x6c, 0x31, 0x7a,
	0x4, 0x6b, 0x65, 0x79, 0x32, 0x7a, 0x4, 0x76, 0x61, 0x6c, 0x32, 0x7a, 0x4,
	0x6b, 0x65, 0x79, 0x33, 0x7a, 0x4, 0x76, 0x61, 0x6c, 0x33, 0x7a, 0x4, 0x6b,
	0x65, 0x79, 0x34, 0x7a, 0x4, 0x76, 0x61, 0x6c, 0x34, 0x80, 0x1, 0xc8, 0x1,
	0x88, 0x1, 0x80, 0xc0, 0x85, 0x3, 0xa0, 0x1, 0x1, 0xa0, 0x1, 0x3, 0xa8, 0x1,
	0x0, 0xa8, 0x1, 0x2, 0xb5, 0x1, 0x0, 0x0, 0x0, 0x0, 0xb5, 0x1, 0x0, 0x0, 0x0,
	0x0, 0xbd, 0x1, 0x0, 0x0, 0x0, 0x0, 0xbd, 0x1, 0x0, 0x0, 0x0, 0x0}

var testTimestamp int64 = 1452146656

// Ugly but verbose test
// TODO: cleanup
func TestPacketPrepare(t *testing.T) {
	packet := Packet{}

	packet.AddRequest(pinbaPacket)
	packet.AddRequest(pinbaPacket)
	packet.AddRequest(pinbaPacket)
	packet.AddRequest(pinbaPacket)
	packet.AddRequest(pinbaPacket)

	data, err := packet.Get(time.Unix(testTimestamp, 0))
	assert.NoError(t, err)

	reader := bytes.NewReader(data)
	message := pinba.ServerMessage{}
	err = message.ReadFrom(reader)
	assert.NoError(t, err)
	assert.EqualValues(t, message.Timestamp, testTimestamp)

	requests, _ := pinba.NewPinbaRequests(message.Timestamp, &message.Data)
	assert.EqualValues(t, testTimestamp, requests.Timestamp)
	assert.EqualValues(t, 5, len(requests.Requests))
	assert.Equal(t, "hostname", requests.Requests[0].Hostname)

	packet.Reset()

	packet.AddRequest(pinbaPacket)
	packet.AddRequest(pinbaPacket)

	data, err = packet.Get(time.Unix(testTimestamp, 0))
	assert.NoError(t, err)

	reader = bytes.NewReader(data)
	message = pinba.ServerMessage{}
	err = message.ReadFrom(reader)
	assert.NoError(t, err)
	assert.EqualValues(t, message.Timestamp, testTimestamp)

	requests, _ = pinba.NewPinbaRequests(message.Timestamp, &message.Data)
	assert.EqualValues(t, testTimestamp, requests.Timestamp)
	assert.EqualValues(t, 2, len(requests.Requests))
	assert.Equal(t, "hostname", requests.Requests[0].Hostname)
}
