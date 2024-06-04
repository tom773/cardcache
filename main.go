package main

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/tom773/cardcache/cache"
	"github.com/tom773/cardcache/protocol"
)

// Mostly boilerplate from the Go God AnthonyGG
const defaultListenAddr = ":42069"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitCh    chan struct{}
	msgCh     chan Message
	cache     *cache.Cache
}

const (
	red        = "\033[31m"
	green      = "\033[32m"
	yellow     = "\033[33m"
	blue       = "\033[34m"
	magenta    = "\033[35m"
	colorReset = "\033[0m"
)

func NewServer(config Config) *Server {
	if len(config.ListenAddr) == 0 {
		config.ListenAddr = defaultListenAddr
	}
	s := &Server{
		Config:    config,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan Message),
		cache:     cache.NewCache(),
	}
	return s
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.loop()

	slog.Info("Listening", "addr", s.ListenAddr)

	return s.acceptLoop()
}
func (s *Server) handleMsg(rawMsg []byte) ([]byte, error) {
	buf := make([]byte, len(rawMsg))

	copy(buf, rawMsg)
	r := protocol.Praw(buf)
	cmd, kv, err := protocol.Pcmd(r)
	var response []byte
	if err != nil {
		return []byte(""), err
	}
	switch cmd {
	case protocol.CmdSet:
		err := s.cache.Set(kv.Key, kv.Value)
		if err != nil {
			return []byte(""), err
		}
		msg := fmt.Sprintf("%sSet key %s with value %s", green, string(kv.Key), string(kv.Value))
		response = []byte(msg)
	case protocol.CmdGet:
		value, err := s.cache.Get(kv.Key)
		if err != nil {
			return []byte(""), err
		}
		msg := fmt.Sprintf("%sValue for key %s: %s", blue, string(kv.Key), string(value))
		response = []byte(msg)
	case protocol.CmdDel:
		err := s.cache.Del(kv.Key)
		if err != nil {
			return []byte(""), err
		}
		msg := fmt.Sprintf("%sDeleted key %s", red, string(kv.Key))
		response = []byte(msg)
	default:
		return []byte(""), fmt.Errorf("Unknown command")
	}

	return response, nil
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		case msg := <-s.msgCh:
			response, err := s.handleMsg(msg.Content)
			if err != nil {
				slog.Error("Handle Msg Err", "err", err)
				continue
			}
			if len(response) == 0 {
				continue
			}
			msg.peer.outCh <- response
		case <-s.quitCh:
			return
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("Accept Err", "err", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer
	slog.Info("New Peer", "addr", conn.RemoteAddr())
	go peer.readLoop()
	go peer.writeLoop()
}

func main() {
	s := NewServer(Config{ListenAddr: ":42069"})
	if err := s.Start(); err != nil {
		slog.Error("Start Err", "err", err)
	}
}
