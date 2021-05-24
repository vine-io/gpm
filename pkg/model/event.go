// MIT License
//
// Copyright (c) 2021 Lack
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package model

// EventAction 事件执行类型
type EventAction string

const (
	CreateProcess  EventAction = "create-process"  // 创建进程
	DeleteProcess  EventAction = "delete-process"  // 删除进程
	UpgradeProcess EventAction = "upgrade-process" // 升级进程
	StartProcess   EventAction = "start-process"   // 启动进程
	StopProcess    EventAction = "stop-process"    // 停止进程
	ReloadProcess  EventAction = "reload-process"  // 重启进程
)

// Event 事件信息
type Event struct {
	// 时间 id
	UUID string `json:"uuid,omitempty"`
	// 事件类型
	Action EventAction `json:"action,omitempty"`
	// Target 事件目标, 如 process id
	Target string `json:"target,omitempty"`
	// 开始时间
	StartTimestamp int64 `json:"startTimestamp,omitempty"`
	// 结束时间
	EndTimestamp int64 `json:"endTimestamp,omitempty"`
	// 详细信息
	Msg string `json:"msg,omitempty"`
}
