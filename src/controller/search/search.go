package search

import (
	"encoding/json"
	"fmt"
	n "github.com/infinit-lab/taiji/src/controller/net"
	"github.com/infinit-lab/taiji/src/controller/system"
	"github.com/infinit-lab/yolanda/config"
	"github.com/infinit-lab/yolanda/logutils"
	"github.com/infinit-lab/yolanda/utils"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type Request struct {
	Command string `json:"command"`
	Session int    `json:"session"`
}

type Response struct {
	Request
	Result bool   `json:"result"`
	Error  string `json:"error"`
}

type udpFrameHandler struct {
	key string
}

type udpClient struct {
	cache   *cache
	timer   *time.Timer
	handler *udpFrameHandler
	addr    *net.UDPAddr
	frameIndex uint16
}

type udpServer struct {
	conn       *net.UDPConn
	cacheMap   map[string]*udpClient
	cacheMutex sync.Mutex
}

const (
	cmdSearch       string = "search"
	cmdNetList      string = "net_list"
	cmdSetNet       string = "set_net"
	cmdUpdate       string = "update"
	cmdUpdateNotify string = "update_notify"
)

var server *udpServer

func init() {
	server = new(udpServer)
	server.cacheMap = make(map[string]*udpClient)

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

	server.conn, err = net.ListenUDP("udp", addr)
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
			_ = server.conn.Close()
		}()

		for {
			data := make([]byte, bufferSize)
			recvLen, rAddr, err := server.conn.ReadFromUDP(data)
			if err != nil {
				logutils.Error("Failed to ReadFromUDP")
				break
			}
			key := fmt.Sprintf("%d.%d.%d.%d:%d", rAddr.IP[0], rAddr.IP[1],
				rAddr.IP[2], rAddr.IP[3], rAddr.Port)
			server.cacheMutex.Lock()
			cache, ok := server.cacheMap[key]
			if !ok {
				cache = new(udpClient)
				cache.handler = new(udpFrameHandler)
				cache.handler.key = key
				cache.cache = newCache(cache.handler)
				cache.timer = time.NewTimer(time.Second * 30)
				cache.addr = rAddr
				go func() {
					select {
					case <-cache.timer.C:
						cache.cache.close()
						server.cacheMutex.Lock()
						delete(server.cacheMap, key)
						server.cacheMutex.Unlock()
					}
				}()
				server.cacheMap[key] = cache
			}
			server.cacheMutex.Unlock()
			if cache != nil {
				cache.timer.Reset(time.Second * 30)
				cache.cache.push(data[:recvLen])
			}
		}
	}()
}

type searchResponse struct {
	Response
	Data struct {
		FingerPrint string `json:"fingerprint"`
		Version     system.Version
	} `json:"data"`
}

type netListResponse struct {
	Response
	Data []*utils.Adapter `json:"data"`
}

type setNetRequest struct {
	Request
	Data utils.Adapter `json:"data"`
}

func (h *udpFrameHandler) onGetFrame(buffer []byte) {
	var request Request
	err := json.Unmarshal(buffer, &request)
	if err != nil {
		logutils.Error("Failed to Unmarshal. error: ", err)
		return
	}
	switch request.Command {
	case cmdSearch:
		var response searchResponse
		response.Request = request
		var err error
		response.Data.FingerPrint, err = utils.GetMachineFingerprint()
		if err != nil {
			logutils.Error("Failed to GetMachineFingerprint. error: ", err)
			h.responseError(request, err.Error())
			return
		}
		response.Data.Version = system.GetVersion()
		response.Result = true
		h.response(response)
	case cmdNetList:
		var response netListResponse
		response.Request = request
		var err error
		response.Data, err = utils.GetNetworkInfo()
		if err != nil {
			logutils.Error("Failed to GetNetworkInfo. error: ", err)
			h.responseError(request, err.Error())
			return
		}
		response.Result = true
		h.response(response)
	case cmdSetNet:
		var req setNetRequest
		err := n.SetAdapter(&req.Data)
		if err != nil {
			logutils.Error("Failed to SetAdapter. error: ", err)
			h.responseError(request, err.Error())
			return
		}
		var response Response
		response.Request = request
		response.Result = true
		h.response(response)
	case cmdUpdate:
	default:
		return
	}
}

func (s *udpServer) getClientByHandler(h *udpFrameHandler) *udpClient {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	client, ok := s.cacheMap[h.key]
	if !ok {
		return nil
	}
	return client
}

func (h *udpFrameHandler) responseError(request Request, err string) {
	var response Response
	response.Request = request
	response.Result = false
	response.Error = err
	h.response(response)
}

func (h *udpFrameHandler) response(response interface{}) {
	data, err := json.Marshal(response)
	if err != nil {
		logutils.Error("Failed to Marshal. error: ", err)
		return
	}
	client := server.getClientByHandler(h)
	if client == nil {
		logutils.Error("Failed to getAddrByHandler")
		return
	}
	client.frameIndex++
	frameList := packBuffer(data, client.frameIndex)
	for _, frame := range frameList {
		_, _ = server.conn.WriteToUDP(frame, client.addr)
	}
}
