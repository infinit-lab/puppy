package search

import (
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/logutils"
	"net"
	"os"
	"strconv"
)

type Request struct {
	Command string `json:"command"`
	Session int `json:"session"`
}

type Response struct {
	Request
	Result bool `json:"result"`
	Error string `json:"error"`
}

type udpFrameHandler struct {

}

type udpServer struct {
}

func init() {
	port := config.GetInt("search.port")
	logutils.Trace("Get search.port is ", port)
	if port == 0 {
		port = 5254
		logutils.Trace("Reset search.port to ", port)
	}
	address := "0.0.0.0:" + strconv.Itoa(port)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		logutils.Error("Failed to ResolveUDPAddr. error: ", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		logutils.Error("Failed to ListenUDP. error: ", err)
		os.Exit(1)
	}

	bufferSize := config.GetInt("search.bufferSize")
	logutils.Trace("Get search.bufferSize is ", bufferSize)
	if bufferSize == 0 {
		bufferSize = 1024
		logutils.Trace("Reset search.bufferSize to ", bufferSize)
	}

	go func() {
		defer func() {
			_ = conn.Close()
		}()

		for {
			data := make([]byte, bufferSize)
			count, rAddr, err := conn.ReadFromUDP(data)
			if err != nil {
				logutils.Error("Failed to ReadFromUDP")
				break
			}
			logutils.Trace("ReadFromUDP. ", string(data[:count]), string(rAddr.IP))
		}
	}()
}

func (h *udpFrameHandler) onGetFrame(buffer []byte) {

}
