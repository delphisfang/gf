package main


import (
    "fmt"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/gfile"
)

func main() {
    v := g.View()
    // 设置模板目录为当前main.go所在目录下的template目录
    v.AddPath(gfile.MainPkgPath() + gfile.Separator + "template")
    b, err := v.Parse("index.html", map[string]interface{} {
        "k" : "v",
    })
    fmt.Println(err)
    fmt.Println(string(b))
}