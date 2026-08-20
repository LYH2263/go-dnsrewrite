package upstream

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

// Encode 长度前缀 JSON 帧。
func Encode(v any) ([]byte, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 4+len(body))
	binary.BigEndian.PutUint32(out[:4], uint32(len(body)))
	copy(out[4:], body)
	return out, nil
}

// Decode 读一帧。
func Decode(r io.Reader, v any) error {
	var hdr [4]byte
	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		return err
	}
	n := binary.BigEndian.Uint32(hdr[:])
	if n > 1<<20 {
		return errFrame
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return err
	}
	return json.Unmarshal(buf, v)
}

type upError string

func (e upError) Error() string { return string(e) }

const errFrame upError = "upstream: frame too large"
