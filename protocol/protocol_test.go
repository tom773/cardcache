package protocol

import (
	"testing"
)

type testCase struct {
	raw string
	exp string
}

func TestPcmd(t *testing.T) {

	testCases := []struct {
		raw string
	}{
		{raw: "SET foo bar"},
		{raw: "GET foo"},
		{raw: "DEL foo"},
	}

	for _, tc := range testCases {
		t.Run(tc.raw, func(t *testing.T) {
			raw := []byte(tc.raw)
			pcmd := Praw(raw)
			if pcmd == nil {
				t.Errorf("Parse(%s) returned nil", tc.raw)
			}
			cmd, data, err := Pcmd(pcmd)
			if err != nil {
				t.Errorf("Pcmd(%s) returned error: %v", tc.raw, err)
			}
			if cmd < 0 || cmd > 2 {
				t.Errorf("Pcmd(%s) returned invalid command: %d", tc.raw, cmd)
			}
			if data == nil {
				t.Errorf("Pcmd(%s) returned nil data", tc.raw)
			}
		})
	}
}
