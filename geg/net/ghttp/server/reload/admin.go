package main

import (
    "gitee.com/johng/gf/g"
)

func main() {
    s := g.Server()
    s.EnableAdmin("/admin")
    //s.BindHookHandler("/admin/*any", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
    //    if !r.BasicAuth("admin", "123", "") {
    //        r.Exit()
    //    }
    //})
    s.SetPort(8199)
    s.Run()
}