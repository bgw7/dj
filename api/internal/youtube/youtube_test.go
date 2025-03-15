package youtube

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestSMSList(t *testing.T) {
	s := `[{"threadid": 232,"type":"msg"}]`
	out := bytes.NewReader([]byte(s))
	var msgs []TextMessage
	if err := json.NewDecoder(out).Decode(&msgs); err != nil {
		print(err)
	}
	print(msgs)
}
