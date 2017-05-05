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
	"time"
	"sync"
)

// 定义错误类型
var (
	ErrHandlerExists = errors.New("This type handler is exists!")
	ErrMaxWorker     = errors.New("Max worker!")
)

// 定义处理方法类型
type Handler func(t Task) error

// Server 配置参数
type ServerOption struct {
	Workersize  uint32
	queueengine QueueEngine
}

// Server 结构
type Server struct {
	mutex       sync.RWMutex
	cursize     uint32
	regularchan chan struct{}
	notifychan  chan struct{}
	mux         map[int]Handler
	option      ServerOption
}

// 创建 Server
func NewServer(option ServerOption) *Server {

	s := Server{
		cursize:         0,
		regularchan:     make(chan struct{}),
		notifychan:      make(chan struct{}),
		mux:             make(map[int]Handler),
		option:          option,
	}

	return &s
}

// 注册 type 对应处理方法
func (this *Server) Register(t int, h Handler) error {
	_, ok := this.mux[t]

	if !ok {
		this.mux[t] = h

		return nil
	}

	return ErrHandlerExists
}

// 创建人物执行者
func (this *Server) NewWorker() (Worker, error) {
	var w Worker

	this.mutex.Lock()
	defer this.mutex.Unlock()

	t, err := this.option.queueengine.FetchTask()

	if err != nil {
		return w, err
	}

	if this.cursize < this.option.Workersize {
		err = this.option.queueengine.DelTask(t.Id)

		if err != nil {
			return w, err
		}

		this.cursize++

		w = Worker{
			s:         this,
			t:         *t,
			h:         this.mux[int(t.Type)],
		}

		return w, nil
	}

	return w, ErrMaxWorker
}

// 通知有新的任务
func (this *Server) Notify() {
	this.notifychan <- struct {}{}
}

// 启动 Server
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