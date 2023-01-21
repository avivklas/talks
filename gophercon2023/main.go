package main

import (
	"bufio"
	"encoding/json"
	"io"
	"sync"
)

type Event struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
	Data string `json:"data"`
}

func reader(rd io.Reader, onLine func([]byte) error) (err error) {
	mu := sync.RWMutex{}
	mu.
	d := json.NewDecoder(rd)
	for d.Decode() {
		err = onLine(s.Bytes())
		if err != nil {
			return
		}
	}

	return s.Err()
}

func processor1(rd io.Reader, onEvent func(event *Event)) error {
	return reader(rd, func(b []byte) (err error) {
		ev := new(Event)
		err = json.Unmarshal(b, ev)
		if err != nil {
			return
		}

		onEvent(ev)

		return
	})
}

func processor2(rd io.Reader, onEvent func(event *Event)) error {
	pool := sync.Pool{New: func() any {
		return new(Event)
	}}
	
	return reader(rd, func(b []byte) (err error) {
		ev := pool.Get().(*Event)

		err = json.Unmarshal(b, ev)
		if err != nil {
			pool.Put(ev)
			return
		}

		onEvent(ev)

		pool.Put(ev)

		return
	})
}
