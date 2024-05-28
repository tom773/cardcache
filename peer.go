package main

import (
	"net"
)

const green = "\033[32m"
const colorReset = "\033[0m"

type Peer struct {
	conn  net.Conn
	msgCh chan []byte
}

type Message struct {
	From    *Peer
	Content []byte
}

func NewPeer(conn net.Conn, msgCh chan []byte) *Peer {
	return &Peer{conn: conn, msgCh: msgCh}
}

func (p *Peer) readLoop() error {
	for {
		buf := Message{
			From:    p,
			Content: make([]byte, 1024),
		}
		n, err := p.conn.Read(buf.Content)
		if err != nil {
			return err
		}
		msgBuf := make([]byte, n)
		copy(msgBuf, buf.Content[:n])
		p.msgCh <- msgBuf
	}
}
