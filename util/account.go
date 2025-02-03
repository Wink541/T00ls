package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type Base64 string

func (Base64 Base64) Decode() []byte {
	data, _ := base64.StdEncoding.DecodeString(string(Base64))
	return data
}

func (Base64 Base64) ToAccountInfo() AccountInfo {
	var accountInfo AccountInfo
	_ = json.Unmarshal(Base64.Decode(), &accountInfo)
	return accountInfo
}

type FileInfo struct {
	Proxy             string
	AccountBase64Text []Base64 `json:"accountBase64Text"`
}

type AccountInfo struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	QuestionId string `json:"questionId"`
	Answer     string `json:"answer"`
}

type LoginResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type LoginRespSuccess struct {
	LoginResp
	Formhash string `json:"formhash"`
	Cookie   struct {
		Auth string `json:"auth"`
	} `json:"cookie"`
}

type SignInResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
)

var (
	Success *log.Logger
	Error   *log.Logger
	Warning *log.Logger
)

func RunTask(filename string) {
	fileInfo, proxyUrl := LoadFile(filename)
	var transport *http.Transport
	proxy, err := url.Parse(proxyUrl)
	if err != nil || proxy.Scheme == "" || proxy.Host == "" || proxy.Port() == "" {
		log.Printf("未发现代理配置,将正常运行")
		transport = &http.Transport{
			Proxy: nil,
		}
	} else {
		log.Printf("配置代理: %s%s%s", Cyan, proxyUrl, Reset)
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}

	for _, base64Text := range fileInfo.AccountBase64Text {
		accountInfo := base64Text.ToAccountInfo()
		AccountSignIn(accountInfo, transport)
	}
}

func LoadFile(filename string) (FileInfo, string) {
	log.Printf("开始读取配置文件: %s%s%s", Cyan, filename, Reset)
	file, err := os.ReadFile(filename)
	if err != nil {
		Error.Printf("配置文件读取错误: %s%s%s, 请检查配置文件后再次运行", Red, err, Reset)
		os.Exit(-1)
	}
	var fileInfo FileInfo
	err = json.Unmarshal(file, &fileInfo)
	if err != nil {
		Error.Printf("配置文件格式错误: %s%s%s, 请检查配置文件后再次运行", Red, err, Reset)
		os.Exit(-1)
	}
	return fileInfo, fileInfo.Proxy
}

func AccountSignIn(accountInfo AccountInfo, transport *http.Transport) {
	log.Printf("用户 %s%s%s 开始登录...", Cyan, accountInfo.Username, Reset)
	loginUrl := "https://www.t00ls.com/login.json"
	loginData := url.Values{
		"action":     []string{"login"},
		"username":   {accountInfo.Username},
		"password":   {accountInfo.Password},
		"questionid": {accountInfo.QuestionId},
		"answer":     {accountInfo.Answer},
	}
	loginReq := CreateReq(loginUrl, loginData, []*http.Cookie{})
	loginBody, loginCookie := POSTRequest(loginReq, transport)
	var loginResp LoginResp
	_ = json.Unmarshal(loginBody, &loginResp)
	if loginResp.Status == "success" {
		Success.Printf("用户 %s%s%s 登录成功: %s登录成功~%s", Green, accountInfo.Username, Reset, Green, Reset)
		var loginRespSuccess LoginRespSuccess
		_ = json.Unmarshal(loginBody, &loginRespSuccess)
		signInUrl := "https://www.t00ls.com/ajax-sign.json"
		signInData := url.Values{
			"formhash":   {loginRespSuccess.Formhash},
			"signsubmit": []string{"true"},
		}
		signInReq := CreateReq(signInUrl, signInData, loginCookie)
		signInBody, _ := POSTRequest(signInReq, transport)
		var signInResp SignInResp
		_ = json.Unmarshal(signInBody, &signInResp)
		if signInResp.Status == "success" && signInResp.Message == "success" {
			Success.Printf("用户 %s%s%s 签到成功: %s%s%s", Green, accountInfo.Username, Reset, Green, "签到成功~", Reset)
			return
		}

		if signInResp.Status == "fail" {
			if signInResp.Message == "alreadysign" {
				Warning.Printf("用户 %s%s%s 签到失败: %s%s%s", Yellow, accountInfo.Username, Reset, Yellow, "今日已签到~", Reset)
				return
			}
			Warning.Printf("用户 %s%s%s 签到失败: %s%s%s", Yellow, accountInfo.Username, Reset, Yellow, signInResp.Message, Reset)
			return
		}
		return
	} else {
		Warning.Printf("用户 %s%s%s 登录失败: %s%s%s ", Yellow, accountInfo.Username, Reset, Yellow, loginResp.Message, Reset)
	}
}

func CreateReq(reqUrl string, data url.Values, cookie []*http.Cookie) *http.Request {
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBufferString(data.Encode()))
	if err != nil {
		Error.Printf("创建请求时发生错误: %s%s%s", Red, err, Reset)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, c := range cookie {
		req.AddCookie(c)
	}
	return req
}

func POSTRequest(req *http.Request, transport *http.Transport) ([]byte, []*http.Cookie) {
	client := &http.Client{
		Transport:     transport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 5,
	}

	resp, err := client.Do(req)
	cookie := resp.Cookies()
	if err != nil {
		Error.Printf("客户端请求出错: %s%s%s", Red, err, Reset)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, cookie
}

func LogKeep(logFile *os.File) {
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetPrefix(fmt.Sprintf("%s[INFO]%s    ", Cyan, Reset))

	Success = log.New(logFile, fmt.Sprintf("%s[SUCCESS]%s ", Green, Reset), log.Ldate|log.Ltime)
	Error = log.New(logFile, fmt.Sprintf("%s[ERROR]%s   ", Red, Reset), log.Ldate|log.Ltime)
	Warning = log.New(logFile, fmt.Sprintf("%s[WARNING]%s ", Yellow, Reset), log.Ldate|log.Ltime)
}

func GetAbsPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}
