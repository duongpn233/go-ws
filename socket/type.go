package socket

import (
	"encoding/hex"
	"encoding/json"
)

type ContentType json.RawMessage

func (c *ContentType) UnmarshalJSON(b []byte) error {
	*c = ContentType(b)
	return nil
}

func (c ContentType) MarshalJSON() ([]byte, error) {
	return json.RawMessage(c).MarshalJSON()
}

type Message struct {
	Event   string      `json:"event"`
	Content ContentType `json:"content"`
}

type HexData []byte

func (h *HexData) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	decoded, err := hex.DecodeString(s)
	if err != nil {
		return err
	}
	*h = HexData(decoded)
	return nil
}
