package search

import (
	"encoding/json"
	"errors"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/logutils"
	"net"
	"os"
	"strconv"
)

type Request struct {
	Type string `json:"type"`
}

type Response struct {
	Result bool `json:"result"`
	Error string `json:"error"`
	Data interface{} `json:"data"`
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
			logutils.Trace("ReadFromUDP. ", string(data[:count]))
			var rsp Response
			rsp.Data, err = request(data[:count], rAddr)
			if err != nil {
				rsp.Result = false
				rsp.Error = err.Error()
			} else {
				rsp.Result = true
			}
			buffer, err := json.Marshal(&rsp)
			if err == nil {
				_, err := conn.WriteToUDP(buffer, rAddr)
				if err != nil {
					logutils.Error("Failed to WriteToUDP. error: ", err)
					break
				}
			}
		}
	}()
}

func request(data []byte, rAddr *net.UDPAddr) (interface{}, error) {
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}
	switch req.Type {
	default:
		return nil, errors.New("Unknown request type: " + req.Type)
	}
}
