## 项目介绍

​	刚注册t00ls的时候发现论坛中有自动签到脚本，在这篇文章中:https://www.t00ls.com/articles-70458.html,但由于这个脚本是用python写的，而我自己的机子上只是跑一些简单的docker容器，没有安装python环境，就想用go来编写一个然后编译为可执行文件，并写入计划任务（python也可以编译，但想锻炼一下go编写脚本的能力）

项目代码：https://github.com/Wink541/T00ls

运行的配置文件如下：

	{
		"proxy": "",
		"accountBase64Text":[
			"",
			""
		]
	}

`proxy`：代理配置，例：http://127.0.0.1:8080

`accountBase64Text`：用户信息`json`格式的`base64`编码字符串（避免尴尬，假设在需要更改配置的时候，直接显示为明文），例：`{"username":"xxxxx","password":"xxxxx","questionId":"xxxxx","answer":"xxxxx"}` 的编码结果为`eyJ1c2VybmFtZSI6Inh4eHh4IiwicGFzc3dvcmQiOiJ4eHh4eCIsInF1ZXN0aW9uSWQiOiJ4eHh4eCIsImFuc3dlciI6Inh4eHh4In0=`

```
username: 用户名
password: 密码
questionId: 安全问题
    # 0 = 没有安全提问
    # 1 = 母亲的名字
    # 2 = 爷爷的名字
    # 3 = 父亲出生的城市
    # 4 = 您其中一位老师的名字
    # 5 = 您个人计算机的型号
    # 6 = 您最喜欢的餐馆名称
    # 7 = 驾驶执照的最后四位数字
answer: 安全问题答案
```



## 运行方式

编译好`t00ls`可执行文件（在Windows中编译方式），再将编译好的文件上传至服务器

```
go env -w GOOS=linux
go build -o t00ls main.go
```



通过指定配置文件路径即可运行，会在`t00ls`所在路径产生日志文件`app.log`，示例

```
$ ./t00ls
Usage: ./t00ls <filename>

$ ./t00ls /opt/t00ls/account.json
```



## 计划任务

计划每12小时执行一次（晚上12点整执行有概率签到失败），也可以根据自己的需求来进行更改

```
0 */12 * * * root /root/t00ls/t00ls /root/t00ls/account.json
```

运行日志如下：

```
[INFO]    2025/01/31 00:00:01 开始读取配置文件: /root/t00ls/account.json
[INFO]    2025/01/31 00:00:01 未发现代理配置,将正常运行
[INFO]    2025/01/31 00:00:01 用户 duoduo 开始登录...
[SUCCESS] 2025/01/31 00:00:04 用户 duoduo 登录成功: 登录成功~
[SUCCESS] 2025/01/31 00:00:05 用户 duoduo 签到成功: 签到成功~
[INFO]    2025/01/31 12:00:01 开始读取配置文件: /root/t00ls/account.json
[INFO]    2025/01/31 12:00:01 未发现代理配置,将正常运行
[INFO]    2025/01/31 12:00:01 用户 duoduo 开始登录...
[SUCCESS] 2025/01/31 12:00:03 用户 duoduo 登录成功: 登录成功~
[WARNING] 2025/01/31 12:00:04 用户 duoduo 签到失败: 今日已签到~
```

