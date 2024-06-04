package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Command byte

const (
	CmdSet Command = iota
	CmdGet
	CmdDel
)

type ICommand interface {
	Bytes() []byte
}

type CommandSet struct {
	Key   []byte
	Value []byte
}

type CommandGet struct {
	Key []byte
}

type CommandDel struct {
	Key []byte
}

type KV struct {
	Key   []byte
	Value []byte
}

func (c Command) String() string {
	switch c {
	case CmdSet:
		return "SET"
	case CmdGet:
		return "GET"
	case CmdDel:
		return "DEL"
	default:
		return "Unknown"
	}
}

func (c *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, CmdSet)
	binary.Write(buf, binary.LittleEndian, uint16(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)

	binary.Write(buf, binary.LittleEndian, uint16(len(c.Value)))
	binary.Write(buf, binary.LittleEndian, c.Value)

	return buf.Bytes()

}

func (g *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, CmdGet)
	binary.Write(buf, binary.LittleEndian, uint16(len(g.Key)))
	binary.Write(buf, binary.LittleEndian, g.Key)

	return buf.Bytes()
}

func (d *CommandDel) Bytes() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, CmdDel)
	binary.Write(buf, binary.LittleEndian, uint16(len(d.Key)))
	binary.Write(buf, binary.LittleEndian, d.Key)

	return buf.Bytes()
}

func Praw(rawMsg []byte) ICommand {

	parts := bytes.Split(rawMsg, []byte(" "))

	cmd := parts[0]
	key := parts[1]

	switch string(cmd) {
	case "SET":
		if len(parts) < 3 {
			fmt.Println(red, "Error parsing message: SET command requires a value", reset)
			return nil
		}
		value := parts[2]
		return &CommandSet{
			Key:   key,
			Value: value,
		}
	case "GET":
		return &CommandGet{
			Key: key,
		}
	case "DEL":
		return &CommandDel{
			Key: key,
		}
	default:
		fmt.Println(red, "Error parsing message: Unknown command. Use SET, GET, or DEL", reset)
		return nil
	}
}
func Pcmd(r ICommand) (Command, *KV, error) {

	var cmd Command
	kv := &KV{}
	Ibytes := bytes.NewReader(r.Bytes())

	err := binary.Read(Ibytes, binary.LittleEndian, &cmd)
	if err != nil {
		fmt.Println(red, "Error reading command", reset)
		return cmd, kv, err
	}

	switch cmd {
	case CmdSet:
		c, err := parseSetCommand(Ibytes)
		if err != nil {
			fmt.Println(red, "Error parsing set command", reset)
			return cmd, kv, err
		}
		kv.Key = bytes.TrimSuffix(c.Key, []byte("\n"))
		kv.Value = bytes.TrimSuffix(c.Value, []byte("\n"))
		return cmd, kv, nil
	case CmdGet:
		g, err := parseGetCommand(Ibytes)
		if err != nil {
			fmt.Println(red, "Error parsing get command", reset)
			return cmd, kv, err
		}
		kv.Key = bytes.TrimSuffix(g.Key, []byte("\n"))
		return cmd, kv, nil
	case CmdDel:
		d, err := parseDelCommand(Ibytes)
		if err != nil {
			fmt.Println(red, "Error parsing del command", reset)
			return cmd, kv, err
		}
		kv.Key = bytes.TrimSuffix(d.Key, []byte("\n"))
		return cmd, kv, nil
	default:
		return cmd, kv, fmt.Errorf("Unknown command: %d", cmd)
	}

}

func parseSetCommand(r io.Reader) (*CommandSet, error) {
	c := &CommandSet{}

	var keyLen uint16
	binary.Read(r, binary.LittleEndian, &keyLen)

	c.Key = make([]byte, keyLen)
	_, err := r.Read(c.Key)
	if err != nil {
		return nil, err
	}

	var valueLen uint16
	binary.Read(r, binary.LittleEndian, &valueLen)

	c.Value = make([]byte, valueLen)
	_, err = r.Read(c.Value)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func parseDelCommand(r io.Reader) (*CommandDel, error) {
	d := &CommandDel{}

	var keyLen uint16
	binary.Read(r, binary.LittleEndian, &keyLen)

	d.Key = make([]byte, keyLen)
	_, err := r.Read(d.Key)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func parseGetCommand(r io.Reader) (*CommandGet, error) {
	g := &CommandGet{}

	var keyLen uint16
	binary.Read(r, binary.LittleEndian, &keyLen)

	g.Key = make([]byte, keyLen)
	_, err := r.Read(g.Key)
	if err != nil {
		return nil, err
	}

	return g, nil
}

const (
	green   = "\033[32m"
	red     = "\033[31m"
	cyan    = "\033[36m"
	magenta = "\033[35m"
	reset   = "\033[0m"
)