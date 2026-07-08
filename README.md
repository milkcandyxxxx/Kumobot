# Kumobot

## 起源

long long ago，一位传奇的少女...；好的不扯了。奶糖第一次解出的机器人是[Yunzai-Bot](https://github.com/yoimiya-kokomi/miao-plugin)，中文名是云崽（一开始是专注于原神的相关功能，后续逐步扩展业务）。Kumo在日语中有云的意思。所以直译过来是云机器人，不过我更想称之为“云宝”。

## 用途

本项目采用go语言编写(支持1.24及以上低版本未测试)，用于接入各种协议的实现，进行各类聊天软件的插件开发及使用。

在尽可能保持功能完整的情况下尽可能的使插件的开发变的简单（目前仅初步实现onebot11和12）

> **注意事项**
>
> 本人也是go初学者，bug，漏洞等可能会比较多。所以也希望大家在使用后，可以多提出一点修改建议，奶糖会尽可能修复！

## 快速开始

- 方法一（推荐）

  直接到发布页下载[脚手架工具](https://github.com/milkcandyxxxx/Kumobot/releases)

  运行cli工具，创建目录结构

  ```bash
  ./Kumobot.exe
  cd {{你的项目名}}
  ```

  ```text
  目录结构
  name
  │  config.yaml
  │  go.mod
  │  main.go
  │
  └─plugins
      │  plugins.go
      │
      └─ping
              ping.go
  ```

  更新依赖库

  ```bash
  go mod tidy
  ```

  然后就可以直接开始书写了！

- 方法二：直接git拉取整个项目库（不推荐不介绍了）

## 协议选择

目前完善的协议仅为onebot11，后续会更新更多协议的支持

## 连接方式

目前仅支持正向的Websocket

## 配置文件

```yaml
# 云宝配置文件
# 连接配置
connect:
  #  WebSocket 地址
  ws_url : ""
  token : ""
```

## 插件的创建

plugins/下，一个文件夹就是一个插件。

```text
└─plugins
    │  plugins.go
    │
    └─ping
    │       ping.go
    │
    └─新插件
            插件文件.go
```

## 插件的基础格式

```go
package ping

import (
    "github.com/milkcandyxxxx/Kumobot/bot" // 主控包
    "github.com/milkcandyxxxx/Kumobot/bot/rule" // 规则匹配包
)

func init() {
    bot.OnPlugin("ping", "无", "milk", "无", "2", "false")// 注册插件，目前是前置写满5个参数。依次：名字，版本，作者，帮助，优先级，独家。全部都是字符串类型就行
    bot.OnMessage(func(ctx *bot.Ctx) {// 注册注册触发词
       ctx.ON11.Send("pong")// 执行的函数
    }, rule.OnCommand("ping"))// 匹配的规则
}
```

> **注**：
>
> 每个插件需要显示在plugins.go中导入
>
> 插件的注册需全部写入init()里进行初始化

```go
package plugins
import _ "项目名称/plugins/ping"
```

## 消息的基础格式

```go
type Event struct { // 共用字段
    Time       int64            // 事件时间戳
    SubType    string           // 事件子类型
    MessageID  string           // 消息唯一ID
    Message    []MessageSegment // 消息段数组
    UserID     string           // 发送者账号
    GroupID    string           // 群组账号
    AltMessage string           // 纯文本
    DetailType string           // 事件细分类型,OB11为MessageType
    // OneBot11 专属字段
    PostType    string     // 事件大类
    MessageType string     // 消息类型
    SelfID      int64      // 机器人自身账号
    Anonymous   any        // 匿名发言信息
    Sender      OB11Sender // 发送者详情
    // OneBot12 专属字段
    ID   string  // 事件标识
    Self BotSelf // 自身账号信息
    Type string  // 事件主类型
    GuildID string // 频道ID
}
```

## 具体方法介绍

### 回调函数 Handler

触发后，所执行的函数

```go
func(ctx *bot.Ctx）{}
```

### 规则 Rule

判断是否执行的规则

```go
type Rule func(ctx *Ctx) bool
```

- `Rule.OnlyGroup() //指定只能群组触发`
- `Rule.Group(groupID string) // 指定什么群触发`
- `Rule.User(userID string) // 指定用户触发`
- `Rule.OnCommand(cmd string) // 指定前置词触发`

### bot

赋值控制机器人的连接终端启动注册插件等

具体方法：

#### `bot.OnPlugin(info ...string)` 注册插件

使用案例：`bot.OnPlugin("ping", "无", "milk", "无", "2", "false")`

#### `bot.OnMessage(h Handler,rules ...Rule)` 注册功能

使用案例：

```go
bot.OnMessage(func(ctx *bot.Ctx) {
    ctx.ON11.Send("pong")
}, rule.OnCommand("ping"))
```

#### `bot.OnCron(spec string, h Handler)` 注册定时任务

使用案例：

#### `bot.OnChat(h Handler, rules ...Rule)` 注册会话，用于多轮会话。

```go
bot.OnChat(func(ctx *bot.Ctx) {
		rand.Seed(time.Now().UnixNano()) //设置随机数种子 
		randNum := rand.Intn(100) // 生成随机数
		ctx.ON11.Send("请输入一个数字猜大小,范围1-100")
		for {
			numberEvent := <-ctx.Ch
			number, _ := strconv.Atoi(numberEvent.AltMessage)
			if number > randNum {
				ctx.ON11.Send("大")
			} else if number < randNum {
				ctx.ON11.Send("小")
			} else {
				ctx.ON11.Send("猜对了！")
				return
			}
		}
	}, rule.OnCommand("猜数字"))
```

### Ctx

消息上下文以及各种api操作的执行者

#### `Ctx.ON11`：onebot11协议的api执行

实现了所有[onebot11协议](https://github.com/botuniverse/onebot-11/blob/master/api/public.md)的公开api，同时，也会对于各自协议添加一些快捷方法。

> **注**：
>
> 本框架中的所有数据存储格式均为string，但onebot11的绝大部分数据类型都是int64，请使用api时转换后再传入。

#### `Ctx.Event`：消息事件（具体内容见上）

例子：Ctx.Event.AltMessage // 当前消息的纯文本

#### `Ctx.Ch`：多轮会话的通道（需要搭配OnChat使用）

例子：numberEvent := <-ctx.Ch

#### 快捷方法

- `Ctx.ON11.Send(msg string) (int64, error) //  一键发送默认为回复（在哪触发的哪里回复，默认at）`
