package main

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/tom773/cardcache/cache"
	"github.com/tom773/cardcache/peer"
	"github.com/tom773/cardcache/protocol"
)

const defaultListenAddr = ":42069"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	peers     map[*peer.Peer]bool
	ln        net.Listener
	addPeerCh chan *peer.Peer
	quitCh    chan struct{}
	msgCh     chan peer.Message
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
		peers:     make(map[*peer.Peer]bool),
		addPeerCh: make(chan *peer.Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan peer.Message),
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
func (s *Server) handleMsg(p *peer.Peer, rawMsg []byte) ([]byte, error) {
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
		msg := fmt.Sprintf("%sSet key %s with value %s%s", green, string(kv.Key), string(kv.Value), colorReset)
		response = []byte(msg)

	case protocol.CmdGet:
		value, error := s.cache.Get(kv.Key)
		if error {
			msg := fmt.Sprintf("%s%s%s", red, value, colorReset)
			response = []byte(msg)
			break
		}
		msg := fmt.Sprintf("%sValue for key %s: %s%s", blue, string(kv.Key), string(value), colorReset)
		response = []byte(msg)

	case protocol.CmdDel:
		err := s.cache.Del(kv.Key)
		if err != nil {
			return []byte(""), err
		}
		msg := fmt.Sprintf("%sDeleted key %s%s", red, string(kv.Key), colorReset)
		response = []byte(msg)
	case protocol.CmdSub:
		go func() {
			ch := s.cache.PubSub.Subscribe(p, string(kv.Key))
			for {
				select {
				case msg := <-ch:
					p.OutCh <- msg
				}
			}
		}()
		msg := fmt.Sprintf("%sSubscribed to key %s%s", yellow, string(kv.Key), colorReset)
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
			response, err := s.handleMsg(msg.Peer, msg.Content)
			if err != nil {
				slog.Error("Handle Msg Err", "err", err)
				continue
			}
			if len(response) == 0 {
				continue
			}
			msg.Peer.OutCh <- response
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
	peer := peer.NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer
	slog.Info("New Peer", "addr", conn.RemoteAddr())
	go peer.ReadLoop()
	go peer.WriteLoop()
}

func main() {
	s := NewServer(Config{ListenAddr: ":42069"})
	if err := s.Start(); err != nil {
		slog.Error("Start Err", "err", err)
	}
}
