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
		panic(err)
	}
	defer logFile.Close()
	util.LogKeep(logFile)

	for {
		err = util.RunTask(*configPath)
		if err != nil {
			util.Error.Println(err)
			fmt.Println("用户签到失败, 请查看日志")
			os.Exit(1)
		}
		fmt.Println("用户签到完成, 等待下一次签到")
		time.Sleep(24 * time.Hour)
	}
}
