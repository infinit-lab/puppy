package search

import (
	"github.com/infinit-lab/yolanda/logutils"
	"testing"
	"time"
)

func TestCacheByteToUint16(t *testing.T) {
	i := byteToUint16(0xAA, 0x55)
	if i != 0x55AA {
		t.Error("Failed to byteToUint16")
	}
	logutils.TraceF("%04x", i)
	i = 0x55AA
	bytes := uint16ToByte(i)
	if len(bytes) != 2 {
		t.Fatal("the length of bytes should be 2")
	}
	if bytes[0] != 0xAA || bytes[1] != 0x55 {
		t.Error("Failed to uint16ToByte")
	}
	logutils.TraceF("%02x%02x", bytes[1], bytes[0])
}

var buffer []byte

type testFrameHandler struct {

}

func (h *testFrameHandler) onGetFrame(b []byte) {
	logutils.TraceF("onGetFrame, %x", b)
	if len(b) != len(buffer) {
		logutils.Error("size is not same.")
		return
	}
	for i, temp := range b {
		if temp != buffer[i] {
			logutils.Error("data is not same.")
			return
		}
	}
}

func TestCachePackBuffer(t *testing.T) {
	for i := 0; i < 0xFF * 3 + 20; i++ {
		buffer = append(buffer, byte(i))
	}
	frameList := packBuffer(buffer, 1)
	c := newCache(new(testFrameHandler))
	for _, frame := range frameList {
		c.push(frame)
	}
	c.close()
	time.Sleep(100 * time.Millisecond)
}
