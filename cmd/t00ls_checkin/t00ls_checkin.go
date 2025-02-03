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
	configPath := flag.String("config", "account.json", "log file path")
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
		util.RunTask(*configPath)
		fmt.Println("Sleep 24 hours for next checkin.")
		time.Sleep(24 * time.Hour)
	}
}
