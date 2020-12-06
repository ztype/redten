package game

import "encoding/json"

type Msg struct {
	Cmd  string          `json:"cmd"`
	Data json.RawMessage `json:"data"`
	mp   map[string]interface{}
}

func NewMsg(cmd string) *Msg {
	return &Msg{Cmd: cmd, mp: make(map[string]interface{})}
}

func (m *Msg) Set(key string, v interface{}) *Msg {
	m.mp[key] = v
	return m
}

func (m *Msg) JsonBytes() []byte {
	if len(m.Data) == 0 {
		md, _ := json.Marshal(m.mp)
		m.Data = md
	}
	bs, _ := json.Marshal(m)
	return bs
}
