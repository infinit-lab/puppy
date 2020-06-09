package search

import (
	"encoding/json"
	"github.com/infinit-lab/yolanda/logutils"
	"net"
	"testing"
	"time"
)

var conn net.Conn

type searchHandler struct {

}

func (h *searchHandler) onGetFrame(buffer []byte) {
	logutils.Trace(string(buffer))
	_ = conn.Close()
}

func sendWait(data []byte) {
	var err error
	frameList := packBuffer(data, 1)
	conn, err = net.Dial("udp", "127.0.0.1:5254")
	if err != nil {
		logutils.Error("Failed to Dial. error: ", err)
		return
	}
	logutils.Trace("Dial success.")
	c := newCache(new(searchHandler))

	for _, frame := range frameList {
		_, _ = conn.Write(frame)
	}

	for {
		buffer := make([]byte, 1024)
		recvLen, err := conn.Read(buffer)
		if err != nil {
			logutils.Error("Failed to Read. error: ", err)
			break
		}
		c.push(buffer[:recvLen])
	}
	conn = nil
}

func TestSearch(t *testing.T) {
	var request Request
	request.Command = cmdSearch
	request.Session = 0

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatal("Failed to Marshal. error: ", err)
	}
	sendWait(data)

	time.Sleep(100 * time.Millisecond)

	request.Command = cmdNetList
	data, err = json.Marshal(request)
	if err != nil {
		t.Fatal("Failed to Marshal. error: ", err)
	}
	sendWait(data)
}


