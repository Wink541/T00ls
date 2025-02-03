package main

import (
	"T00ls/util"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var ConfigPath string

func main() {
	logPath := "app.log"
	AbsPath := util.GetAbsPath()
	configPath := flag.String("config", "account.json", "config file path")
	flag.Parse()

	if *configPath == "" {
		flag.Usage()
		return
	}

	ConfigPath = *configPath
	logFile, err := os.OpenFile(filepath.Join(AbsPath, logPath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		util.Error.Println(err)
		fmt.Printf("日志文件打开失败: %v", err)
		os.Exit(1)
	}
	defer logFile.Close()
	util.LogKeep(logFile)

	configPassword, err := util.ReadBytesFromStdIN("请输入密码: ")
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		username, err := util.ReadBytesFromStdIN("请输入用户名: ")
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		password, err := util.ReadBytesFromStdIN("请输入密码: ")
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		secQuestionID, err := util.ReadBytesFromStdIN("请输入密保问题ID: ")
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		secQuestionAnswer, err := util.ReadBytesFromStdIN("请输入密保问题答案: ")
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		err = makeConfig(string(username),
			string(password),
			string(secQuestionID),
			string(secQuestionAnswer),
			configPassword)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
	}

	for {
		err = util.RunTask(ConfigPath, configPassword)
		if err != nil {
			util.Error.Println(err)
			fmt.Printf("用户签到失败: %v", err)
			os.Exit(1)
		}
		fmt.Println("用户签到完成, 等待下一次签到")
		time.Sleep(24 * time.Hour)
	}
}

func makeConfig(username, password, secQuestionID, secQuestionAnswer string, configPassword []byte) error {
	accountInfo := util.AccountInfo{
		Username:   username,
		Password:   password,
		QuestionId: secQuestionID,
		Answer:     secQuestionAnswer,
	}
	encryptedAccountInfo := accountInfo.ToBase64Text(configPassword)
	config := util.ConfigInfo{
		Proxy:             "",
		AccountBase64Text: []util.Base64UserConfig{encryptedAccountInfo},
	}

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		err = fmt.Errorf("生成配置文件失败: %s%s%s", util.Red, err, util.Reset)
		return err
	}

	return os.WriteFile(ConfigPath, data, 0600)
}
