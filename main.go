// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/wenchangshou2/zui/pkg/logging"
	"github.com/wenchangshou2/zui/pkg/setting"
	"github.com/wenchangshou2/zutil"
	"log"
)

func main() {
	var (
		exit_chan chan bool
		err       error
	)
	exit_chan = make(chan bool)
	confPath, _ := zutil.GetFullPath("conf/app.ini")
	if err = setting.InitSetting(confPath); err != nil {
		log.Fatal("读取配置文件失败:" + err.Error())
		return
	}
	InitComputerControl()
	logPath, _ := zutil.GetFullPath(setting.AppSetting.LogSavePath)
	if err = logging.InitLogging(logPath, setting.AppSetting.LogLevel); err != nil {
		log.Fatalf("创建日志模块失败")
		return
	}
	if err = InitSchedule(setting.ServerSetting.Ip, setting.ServerSetting.Port); err != nil {
		log.Fatalf("初始化调度失败:"+err.Error())
		return
	}
	fmt.Println("init ui")
	go InitUI(exit_chan)
	select {
	case e := <-exit_chan:
		fmt.Println("exit", e)

	}
}
