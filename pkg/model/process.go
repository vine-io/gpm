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

import (
	"os"
)

type ProcessExec string

const (
	// ForkExec 二进制可执行文件
	ForkExec ProcessExec = "fork"
	// JavaExec Java 执行器
	JavaExec ProcessExec = "java"
	// PythonExec python 执行器
	PythonExec ProcessExec = "python"
	// NodeExec node 执行器
	NodeExec ProcessExec = "node"
	// ShellExec unix,linux 下为 /bin/sh, windows 下则为 cmd
	ShellExec ProcessExec = "shell"
)

// ProcessStatus 进程状态
type ProcessStatus string

const (
	StatusInit      ProcessStatus = "init"      // 进程初始化, 还未启动过
	StatusRunning   ProcessStatus = "running"   // 进程正在运行中
	StatusStopped   ProcessStatus = "stopped"   // 进程停止
	StatusFailed    ProcessStatus = "failed"    // 进程执行失败
	StatusReload    ProcessStatus = "reload"    // 进程重启
	StatusUpgrading ProcessStatus = "upgrading" // 进程升级中
)

// Process 进程信息
type Process struct {
	// 进程 id
	ID int64 `gorm:"primaryKey" json:"id,omitempty" yaml:"-"`
	// 进程名称
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// 进程执行器，默认为 fork
	Exec ProcessExec `json:"exec,omitempty" yaml:"exec,omitempty"`
	// 执行器路径, Exec 为 fork 时, Bin 为二进制文件,否则 Bin 为执行器路径
	Bin string `json:"bin,omitempty" yaml:"bin,omitempty"`
	// 进程启动参数
	Args string `json:"args,omitempty" yaml:"args,omitempty"`
	// 进程执行时的 pid
	PID int64 `json:"pid,omitempty" yaml:"pid,omitempty"`
	// 执行命令时所在的目录
	Chroot string `json:"chroot,omitempty" yaml:"chroot,omitempty"`
	// 执行进程时所用的 uid (windows 无效)
	UID int32 `json:"uid,omitempty" yaml:"uid,omitempty"`
	// uid 对应的用户名称
	User string `json:"user,omitempty" yaml:"user,omitempty"`
	// 执行进程时所有的 gid (windows 无效)
	GID int32 `json:"gid,omitempty" yaml:"gid,omitempty"`
	// gid 对应的用户组
	Group string `json:"group,omitempty" yaml:"group,omitempty"`
	// 进程版本
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	// 是否自启动, 默认为 false
	AutoRestart bool `json:"autoRestart,omitempty" yaml:"autoRestart,omitempty"`

	// 创建时间
	CreationTimestamp int64 `gorm:"autoUpdateTime:nano" json:"creationTimestamp" yaml:"-"`
	// 修改时间
	UpdateTimestamp int64 `gorm:"autoUpdateTime:nano" json:"updateTimestamp,omitempty" yaml:"-"`
	// 删除时间
	DeletionTimestamp int64 `json:"deletionTimestamp,omitempty" yaml:"-"`
	// 进程状态
	Status ProcessStatus `json:"status,omitempty" yaml:"-"`
	// 进程状态为 failed 的错误信息
	Msg string `json:"msg,omitempty" yaml:"-"`

	// point to os.Process
	Pointer *os.Process `gorm:"-" json:"-" yaml:"-"`
	// cpu, 内存信息
	Stat *Stat `gorm:"-" json:"stat,omitempty" yaml:"-"`
}

func (p Process) TableName() string {
	return "process"
}

type Stat struct {
	// cpu 占用百分比
	CPUPercent float32 `json:"cpuPercent"`
	// 内存占用, rss
	Memory uint64 `json:"memory"`
	// 内存占用百分比
	MemPercent float32 `json:"memPercent"`
}
