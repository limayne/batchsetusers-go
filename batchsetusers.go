package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "time"
    "regexp"
)

func main() {
    if os.Geteuid() != 0 {
        fmt.Println("This script must be run as root")
        os.Exit(1)
    }

    fmt.Print("请输入需要创建的用户，每个用户名后请用\\分隔，无需空格，以回车键结束：\n")
    var userlist string
    fmt.Scanln(&userlist)

    if userlist == "" {
        fmt.Println("请输入至少一个用户名")
        os.Exit(1)
    }

    users := strings.Split(userlist, "\\")
    for _, user := range users {
        if !isValidUsername(user) {
            fmt.Printf("用户名 %s 不符合命名规范，请使用小写字母、数字、下划线和连字符\n", user)
            os.Exit(1)
        }

        if _, err := exec.Command("id", user).Output(); err == nil {
            fmt.Printf("用户 %s 已经存在\n", user)
            continue
        }

        pass, _ := exec.Command("openssl", "rand", "-base64", "8").Output()
        pass = []byte(strings.TrimRight(string(pass), "\n"))

        if err := exec.Command("useradd", "-m", "-p", string(pass), "-d", "/home/"+user, user).Run(); err != nil {
            fmt.Printf("创建用户 %s 失败\n", user)
            os.Exit(1)
        }

        userfile, err := os.OpenFile("./user.info", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
        if err != nil {
            fmt.Printf("打开用户信息文件失败: %v\n", err)
            os.Exit(1)
        }
        fmt.Fprintf(userfile, "%s:%s\n", user, pass)
        userfile.Close()
        fmt.Printf("用户 %s 创建成功\n", user)
    }

    logfile, err := os.OpenFile("./user.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        fmt.Printf("打开日志文件失败: %v\n", err)
        os.Exit(1)
    }
    fmt.Fprintf(logfile, "%s 用户创建脚本执行成功\n", time.Now().Format("2006-01-02 15:04:05"))
    logfile.Close()
}

func isValidUsername(username string) bool {
    match, err := regexp.MatchString("^[a-z_][a-z0-9_-]{0,30}$", username)
    return err == nil && match
}