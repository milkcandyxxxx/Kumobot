/**
 * @author milkcandy
 * @date 2026/6/15
 * @description TODO
 */

package ON11

import (
	"fmt"
	"github.com/milkcandyxxxx/Kumobot/adapter/onebot11"
	"github.com/milkcandyxxxx/Kumobot/bot"
	"strconv"
)

type Ctx struct {
	*CronCtx
}
type Plugin struct {
	*bot.PluginInfo
}
type CronCtx struct {
	*bot.Ctx
	*onebot11.Adapter
}

// OnMessage 注册消息
// 注意：这里的*Ctx是ON11的ctx

func OnMessage(p *bot.PluginInfo, h func(*Ctx)) *bot.Responder {
	r := &bot.Responder{
		Type: "Message",
		Handler: func(base *bot.Ctx) {
			ctx := &Ctx{
				&CronCtx{
					Ctx:     base,
					Adapter: base.Bot.Adapter.(*onebot11.Adapter),
				},
				// Ctx:     base,
				// Adapter: base.Bot.Adapter.(*onebot11.Adapter),
			}
			h(ctx)
			// 这段代码存在的意义是，Handler里需要存入的是func(bot.ctx)但是这里只有ON11的，所以只能自行将bot.ctx转化为on11.ctx
		},
	}
	p.Respond = append(p.Respond, r)
	// 写入对应插件
	return r
}

// OnCron 定时任务注册（6位cron表达式：秒 分 时 日 月 周）
func OnCron(p *bot.PluginInfo, spec string, h func(ctx *CronCtx)) *bot.Responder {
	r := &bot.Responder{
		Type: "Cron",
		Handler: func(base *bot.Ctx) {
			ctx :=
				&CronCtx{
					Ctx:     base,
					Adapter: base.Bot.Adapter.(*onebot11.Adapter),
				}
			h(ctx)
			// 这段代码存在的意义是，Handler里需要存入的是func(bot.ctx)但是这里只有ON11的，所以只能自行将bot.ctx转化为on11.ctx
		},
		Timing: spec,
	}
	p.Respond = append(p.Respond, r)
	fmt.Printf("%+v\n", r)
	return r
}

// SendGroupMessageAt  发送群聊信息，并at回复
func (o *CronCtx) SendGroupMessageAt(groupID string, msg string) (int64, error) {
	// 构造at消息段
	UserID, _ := strconv.ParseInt(o.Event.UserID, 10, 64)
	msgAt := []map[string]interface{}{
		{
			"type": "at",
			"data": map[string]interface{}{
				"qq": UserID,
			},
		},
		{
			"type": "text",
			"data": map[string]interface{}{
				"text": msg,
			},
		},
	}
	return o.SendGroupMessage(groupID, msgAt)
}

// Send 一键发送默认为回复（在哪触发的哪里回复，默认at）目前只适用于qq，部分平台id不为用户名。例如telegram直接写@+名字即可
func (o *Ctx) Send(msg string) (int64, error) {
	if o.Event.DetailType == "private" {
		id, _ := strconv.ParseInt(o.Event.UserID, 10, 64)
		return o.SendPrivateMessage(id, msg)
	}
	if o.Event.DetailType == "channel" || o.Event.DetailType == "group" {
		return o.SendGroupMessageAt(o.Event.GroupID, msg)
	}
	return 0, nil
}
