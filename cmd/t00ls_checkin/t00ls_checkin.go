package main

import (
	"T00ls/util"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	logPath := "app.log"
	AbsPath := util.GetAbsPath()
	configPath := flag.String("config", "account.json", "config file path")
	flag.Parse()

	if *configPath == "" {
		flag.Usage()
		return
	}

	logFile, err := os.OpenFile(filepath.Join(AbsPath, logPath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		util.Error.Println(err)
		fmt.Printf("日志文件打开失败: %v", err)
		os.Exit(1)
	}
	defer logFile.Close()
	util.LogKeep(logFile)

	for {
		err = util.RunTask(*configPath)
		if err != nil {
			util.Error.Println(err)
			fmt.Printf("用户签到失败: %v", err)
			os.Exit(1)
		}
		fmt.Println("用户签到完成, 等待下一次签到")
		time.Sleep(24 * time.Hour)
	}
}
