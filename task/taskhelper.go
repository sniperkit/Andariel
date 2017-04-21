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
	"sync"
	"errors"
)

var (
	ErrTasksLenNotEnough = errors.New("tasks has not enough items for pop")
	ErrTasksCapNotEnough = errors.New("tasks has not enough space for push")
)

var (
	t *taskHelper
	o sync.Once
)

type taskHelper struct {
	tasks       []Task
	head        int
	tail        int
	size        int
}

func GetTaskHelper() *taskHelper {
	o.Do(func() {
		t = &taskHelper{
			tasks: make([]Task, TaskSize + 1),
			size: TaskSize,
		}
	})

	return t
}

func (this *taskHelper) Len() int {
	if this.head == this.tail {

		return 0
	} else if this.tail > this.head {

		return this.tail - this.head
	} else {

		return this.tail + this.size + 1 - this.head
	}
}

func (this *taskHelper) Cap() int {
	return this.size - this.Len()
}

func (this *taskHelper) Pop() (Task, error) {
	if this.head == this.tail {
		return nil, ErrTasksLenNotEnough
	}

	t := this.tasks[this.head]
	this.tasks[this.head] = nil
	this.head ++

	return t, nil
}

func (this *taskHelper) Push(t Task) error {
	if this.head == this.tail {
		return ErrTasksCapNotEnough
	}

	this.tasks[this.tail] = t
	this.tail = (this.tail + 1) % (this.size + 1)

	return nil
}

func (this *taskHelper) IsFull() bool {
	return this.Cap() == 0
}

func (this *taskHelper) IsEmpty() bool {
	return this.Len() == 0
}

func (this *taskHelper) sendToHandler(task Task) {

}

func (this *taskHelper) GetTaskFromMD() {

}
