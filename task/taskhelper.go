/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/04/21        Liu Jiachang
 */

package task

import (
	"errors"
	"gopkg.in/mgo.v2"
	"time"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

var (
	ErrHandlerExists = errors.New("This type handler is exists!")
	ErrMaxWorker     = errors.New("Max worker!")
	ErrNoTask        = errors.New("There is no task!")
)

type Handler func(t Task) error

type Server struct {
	mutex       sync.RWMutex
	workersize  uint32
	cursize     uint32
	regularchan chan struct{}
	notifychan  chan struct{}
	taskSession *mgo.Session
	mux         map[int]Handler
}

func NewServer(url string, size uint32) (*Server, error) {
	Session, err := mgo.DialWithTimeout(url, 5 * time.Second)

	if err != nil {
		return nil, err
	}

	s := Server{
		workersize:       size,
		cursize:          0,
		regularchan:      make(chan struct{}),
		notifychan:       make(chan struct{}),
		taskSession:      Session,
		mux:              make(map[int]Handler),
	}

	return &s, nil
}

func (this *Server) Register(t int, h Handler) error {
	_, ok := this.mux[t]

	if !ok {
		this.mux[t] = h

		return nil
	}

	return ErrHandlerExists
}

func (this *Server) NewWorker() (Worker, error) {
	var t Task

	c := this.taskSession.DB(MDbName).C(MDColl)
	c.Find(bson.M{}).One(&t)
	c.RemoveId(t.Id)

	if t.Type == 0 {
		return nil, ErrNoTask
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	if this.cursize < this.workersize {
		this.cursize++

		w := Worker{
			s:         this,
			t:         t,
			h:         this.mux[int(t.Type)],
		}

		return w, nil
	}

	return nil, ErrMaxWorker
}

func (this *Server) Notify() {
	this.notifychan <- struct {}{}
}

func (this *Server) Start() {
	go func() {
		for {
			time.Sleep(3 * time.Second)
			this.regularchan <- struct {}{}
		}

	}()

	for {
		select {
		case <- this.notifychan:
			w, err := this.NewWorker()

			if err == nil {
				w.Run()
			}
		case <- this.regularchan:
			w, err := this.NewWorker()

			if err == nil {
				w.Run()
			}
		}
	}
}