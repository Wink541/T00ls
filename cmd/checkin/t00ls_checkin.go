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
		promptConfig(configPassword)
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

func promptConfig(configPassword []byte) {
	secQuestions := `
    # 安全提问参考
    # 0 = 没有安全提问
    # 1 = 母亲的名字
    # 2 = 爷爷的名字
    # 3 = 父亲出生的城市
    # 4 = 您其中一位老师的名字
    # 5 = 您个人计算机的型号
    # 6 = 您最喜欢的餐馆名称
    # 7 = 驾驶执照的最后四位数字`
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
	secQuestionID, err := util.ReadBytesFromStdIN(secQuestions + "\n请输入密保问题ID: ")
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
