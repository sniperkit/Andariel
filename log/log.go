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
 *     Initial: 2017/04/28        Yusan Kurban
 */

package log

import (
	"github.com/Sirupsen/logrus"
)

type AndarielLogLevel logrus.Level

const (
	AndarielLogLevelDebug = AndarielLogLevel(logrus.DebugLevel)
	AndarielLogLevelInfo = AndarielLogLevel(logrus.InfoLevel)
	AndarielLogLevelWarn = AndarielLogLevel(logrus.WarnLevel)
	AndarielLogLevelError = AndarielLogLevel(logrus.ErrorLevel)

	// 设置全局日志设定
	AndarielLogLevelDefault = AndarielLogLevelDebug
)

// 日志通用接口
type AndarielLoggerInf interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type AndarielLogger struct {
	entry *logrus.Entry
}

type AndarielLogTag logrus.Fields

// 创建日志管理工具
func AndarielCreateLogger(tags *AndarielLogTag, level AndarielLogLevel) *AndarielLogger {
	entry := logrus.WithFields(logrus.Fields(*tags))
	entry.Logger.Level = logrus.Level(level)

	return &AndarielLogger{entry: entry}
}

func (this *AndarielLogger) Debug(args ...interface{}) {
	this.entry.Debug(args)
}

func (this *AndarielLogger) Info(args ...interface{}) {
	this.entry.Info(args)
}

func (this *AndarielLogger) Warn(args ...interface{}) {
	this.entry.Warn(args)
}

func (this *AndarielLogger) Error(args ...interface{}) {
	this.entry.Error(args)
}
