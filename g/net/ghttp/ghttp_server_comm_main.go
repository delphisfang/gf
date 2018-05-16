// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// Web Server进程间通信 - 主进程.
// 管理子进程按照规则听话玩，不然有一百种方法让子进程在本地混不下去.

package ghttp

import (
    "os"
    "time"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/os/gproc"
)

// (主进程)主进程与子进程上一次活跃时间映射map
var procUpdateTimeMap = gmap.NewIntIntMap()

// 开启服务
func onCommMainStart(pid int, data []byte) {
    p := procManager.NewProcess(os.Args[0], os.Args, os.Environ())
    p.Start()
    sendProcessMsg(p.Pid(), gMSG_START, nil)
}

// 心跳处理(方法为空，逻辑放到公共通信switch中进行处理)
func onCommMainHeartbeat(pid int, data []byte) {

}

// 平滑重启服务
func onCommMainReload(pid int, data []byte) {
    procManager.Send(formatMsgBuffer(gMSG_RELOAD, nil))
}

// 完整重启服务
func onCommMainRestart(pid int, data []byte) {
    // 如果是父进程接收到重启指令，那么通知所有子进程重启
    if pid == gproc.Pid() {
        procManager.Send(formatMsgBuffer(gMSG_RESTART, nil))
        return
    }
    // 否则杀掉子进程，然后新建一个完整的子进程
    if p, err := os.FindProcess(pid); err == nil && p != nil {
        p.Kill()
        p.Wait()
    }
    sendProcessMsg(gproc.Pid(), gMSG_START, nil)
}

// 新建子进程通知
func onCommMainNewFork(pid int, data []byte) {
    procManager.AddProcess(pid)
    checkHeartbeat.Set(true)
}

// 关闭服务，通知所有子进程退出(Kill强制性退出)
func onCommMainShutdown(pid int, data []byte) {
    procManager.Send(formatMsgBuffer(gMSG_CLOSE, nil))
    procManager.KillAll()
    procManager.WaitAll()
}

// 更新指定进程的通信时间记录
func updateProcessCommTime(pid int) {
    procUpdateTimeMap.Set(pid, int(gtime.Millisecond()))
}

// 主进程与子进程相互异步方式发送心跳信息，保持活跃状态
func handleMainProcessHeartbeat() {
    for {
        time.Sleep(gPROC_HEARTBEAT_INTERVAL*time.Millisecond)
        procManager.Send(formatMsgBuffer(gMSG_HEARTBEAT, nil))
        // 清理过期进程
        if checkHeartbeat.Val() {
            for _, pid := range procManager.Pids() {
                if int(gtime.Millisecond()) - procUpdateTimeMap.Get(pid) > gPROC_HEARTBEAT_TIMEOUT {
                    // 这里需要手动从进程管理器中去掉该进程
                    procManager.RemoveProcess(pid)
                    sendProcessMsg(pid, gMSG_CLOSE, nil)
                }
            }
        }
        // (双保险)如果所有子进程都退出，并且主进程未活动达到超时时间，那么主进程也没存在的必要
        if procManager.Size() == 0 && int(gtime.Millisecond()) - lastUpdateTime.Val() > gPROC_HEARTBEAT_TIMEOUT{
            os.Exit(0)
        }
    }
}