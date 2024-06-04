package main

import (
	"net"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
	outCh chan []byte
}

type Message struct {
	peer    *Peer
	Content []byte
}

func NewPeer(conn net.Conn, msgCh chan Message) *Peer {
	return &Peer{conn: conn, msgCh: msgCh, outCh: make(chan []byte)}
}

func (p *Peer) readLoop() error {
	for {
		buf := make([]byte, 1024)
		n, err := p.conn.Read(buf)
		if err != nil {
			return err
		}

		msg := Message{
			peer:    p,
			Content: make([]byte, n),
		}

		copy(msg.Content, buf[:n])
		p.msgCh <- msg
	}
}

func (p *Peer) writeLoop() error {
	for {
		select {
		case msg := <-p.outCh:
			finalmsg := append(msg, []byte("\n")...)
			_, err := p.conn.Write(finalmsg)
			if err != nil {
				return err
			}
		}
	}
}
