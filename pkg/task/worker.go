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
 *     Initial: 2017/04/30        Liu Jiachang
 */

package task

import (
    "fmt"
)

type Worker struct {
    server      *Server
    tchan       chan Task
    task        Task
}

func (this *Worker) Run() {
    go func() {
        for {
            this.server.distChan <- this.tchan
            this.task = <- this.tchan
            handler, ok := this.server.mux[int(this.task.Type)]

            if ok {
                err := handler(&this.task)

                if err != nil {
                    fmt.Println("[ERROR]:", this.task.Id, " Handler err")
                    this.server.ChangeActive(this.task.Id, TaskUnexecuted)
                } else {
                    this.server.DelTask(this.task.Id)
                }
            }
        }
    }()
}