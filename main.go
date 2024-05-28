package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
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
	msgCh     chan []byte
}

func NewServer(config Config) *Server {
	if len(config.ListenAddr) == 0 {
		config.ListenAddr = defaultListenAddr
	}
	s := &Server{
		Config:    config,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan []byte),
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
func (s *Server) handleMsg(rawMsg []byte) error {
	fmt.Fprintf(os.Stdout, "%s%s", green, rawMsg)
	return nil

}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		case rawMsg := <-s.msgCh:
			if err := s.handleMsg(rawMsg); err != nil {
				slog.Error("Handle Msg Err", "err", err)
			}
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
}

func main() {
	s := NewServer(Config{ListenAddr: ":42069"})
	if err := s.Start(); err != nil {
		slog.Error("Start Err", "err", err)
	}
}
