/**
 * @author milkcandy
 * @date 2026/5/22
 * @description TODO
 */

package onebot11

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/milkcandyxxxx/Kumobot/adapter"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ################################ 公共方法 #####################################

// Adapter OneBot11Adapter 适配器结构体
type Adapter struct {
	*adapter.AdapterInfo
}

// Connect 连接ws
func (a *Adapter) Connect() error {
	// 检测是ws还是http
	if a.WsUrl == "" {
		return nil
	}
	// 检测地址合法性
	wslAddr, err := url.Parse(a.WsUrl)
	if err != nil {

		log.Println("地址格式错误")
		return err
	}
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Token))
	conn, _, err := websocket.DefaultDialer.Dial(wslAddr.String(), header)
	if err != nil {
		log.Println("连接失败")
		return err
	}
	a.Conn = conn
	log.Println("ws连接成功")
	return nil
}

// Disconnect 断开连接
func (a *Adapter) Disconnect() error {
	// 避免未连接就断开
	if a.Conn != nil {
		return a.Conn.Close()
	}
	return nil
}

// ReadMessage 读取信息
func (a *Adapter) ReadMessage(ch chan *adapter.Event) error {
	for {
		_, message, err := a.Conn.ReadMessage()
		if err != nil {
			log.Println(err, "ws连接出错，自动尝试重新连接")
			a.Connect()

		}
		log.Println(string(message))
		if err != nil {
			log.Println("获取消息失败")
		}
		// 先放入万能字典，查看是事件还是动作响应
		var check map[string]interface{}
		err = json.Unmarshal(message, &check)
		if err != nil {
			log.Println("解析消息失败", err)
		}
		if _, exists := check["self_id"]; exists {
			// 看事件类型
			var oneBotEvent OneBotEvent
			err = json.Unmarshal(message, &oneBotEvent)
			if err != nil {
				log.Println("解析消息失败")
				continue
			}
			// 根据类型决策处置方法
			switch oneBotEvent.PostType {
			// 为消息，走插件流程
			case "message":
				var oneBotEvent OneBotEvent
				err = json.Unmarshal(message, &oneBotEvent)
				var oB11BaseMessageEvent OB11BaseMessageEvent
				err = json.Unmarshal(message, &oB11BaseMessageEvent)
				if err != nil {
					log.Println("解析消息失败")
					continue
				}
				// 判断消息类型
				switch oB11BaseMessageEvent.MessageType {
				// 为私聊
				case "private":
					var oB11PrivateMessageEvent OB11PrivateMessageEvent
					err = json.Unmarshal(message, &oB11PrivateMessageEvent)
					event := oB11PrivateMessageEvent.ToAdapterEvent()
					fmt.Println(event)
					ch <- event
					continue
				// 为群组
				case "group":
					var oB11GroupMessageEvent OB11GroupMessageEvent
					err = json.Unmarshal(message, &oB11GroupMessageEvent)
					event := oB11GroupMessageEvent.ToAdapterEvent()
					fmt.Println(event)
					ch <- event
					continue
				}
				// case "meta_event":
				// 	atomic.StoreInt64(&a.HeartbeatTime, time.Now().Unix())
				// 	continue
			}

		} else {
			var res adapter.Response
			err = json.Unmarshal(message, &res)
			a.Mu.Lock()
			ch, ok := a.Echo[res.Echo]
			a.Mu.Unlock()
			if ok {
				ch <- res
				continue
			}
			continue

		}
		continue
	}
}

// 原本想加，但是发现有点多余，先暂存
// // CheckHeartbeat 心跳检查
// func (a *Adapter) CheckHeartbeat(timeOut int64, checkFrequency time.Duration) {
// 	ticker := time.NewTicker(checkFrequency)
// 	defer ticker.Stop()
// 	for range ticker.C {
// 		heartTime := atomic.LoadInt64(&a.HeartbeatTime)
// 		if time.Now().Unix()-heartTime > timeOut {
// 			fmt.Println("超时")
// 		} else {
// 			fmt.Println("不超时")
// 		}
// 	}
// }

// #################################################私有方法##########################################

// CallAction 执行特定方法用于处理某些协议特定的api
func (a *Adapter) CallAction(action string, params map[string]interface{}) (adapter.Response, error) {
	payload := map[string]interface{}{
		"action": action,
		"params": params,
	}
	echo := time.Now().UnixNano()
	echoStr := strconv.FormatInt(echo, 10)
	payload["echo"] = echoStr
	ch := make(chan adapter.Response, 1)
	a.Mu.Lock()
	a.Echo[echoStr] = ch
	a.Mu.Unlock()
	defer func() {
		a.Mu.Lock()
		delete(a.Echo, echoStr)
		a.Mu.Unlock()
	}()
	// 序列化
	body, _ := json.Marshal(payload)
	// if a.WsUrl == "" {
	// 	req, err := http.NewRequest("POST", a.WsUrl, bytes.NewBuffer(body))
	// }

	err := a.Conn.WriteMessage(websocket.TextMessage, body)
	if err != nil {
		return adapter.Response{}, err
	}
	select {
	case resp := <-ch:
		fmt.Println("返还值", resp)
		return resp, nil
	}
}

// SendPrivateMessage 发送私聊信息
func (a *Adapter) SendPrivateMessage(userID int64, msg interface{}) (int64, error) {
	// 构建消息段
	params := map[string]interface{}{
		"user_id": userID,
		"message": msg,
	}
	res, err := a.CallAction("send_private_msg", params)
	if err != nil {
		return 0, err
	}
	var Data struct {
		MsgID int64 `json:"message_id"`
	}
	if err := json.Unmarshal(res.Data, &Data); err != nil {
		return 0, err
	}
	return Data.MsgID, nil
}

// SendGroupMessage  发送群组信息
func (a *Adapter) SendGroupMessage(groupID string, msg interface{}) (int64, error) {
	// 构建消息段,因为总event获取的都是string，但是onebot11协议数字类型都是int

	group, _ := strconv.ParseInt(groupID, 10, 64)
	params := map[string]interface{}{
		"group_id": group,
		"message":  msg,
	}

	res, err := a.CallAction("send_group_msg", params)
	if err != nil {
		return 0, err
	}
	var Data struct {
		MsgID int64 `json:"message_id"`
	}
	if err := json.Unmarshal(res.Data, &Data); err != nil {
		return 0, err
	}
	return Data.MsgID, nil
}

// DeleteMsg  撤回消息
func (a *Adapter) DeleteMsg(msgID int64) error {
	params := map[string]interface{}{
		"message_id": msgID,
	}
	_, err := a.CallAction("delete_msg", params)
	if err != nil {
		return err
	}
	return nil
}

// GetMsgData GetMsg的返回结构体
type GetMsgData struct {
	Time        int64              `json:"time"`
	MessageType string             `json:"message_type"`
	MessageID   int64              `json:"message_id"`
	RealID      int64              `json:"real_id"`
	Sender      adapter.OB11Sender `json:"sender"`
	Message     interface{}        `json:"message"`
}

// GetMsg 获取消息
func (a *Adapter) GetMsg(msgID int64) (GetMsgData, error) {
	params := map[string]interface{}{
		"message_id": msgID,
	}
	res, err := a.CallAction("get_msg", params)
	if err != nil {
		return GetMsgData{}, err
	}
	var Data GetMsgData
	err = json.Unmarshal(res.Data, &Data)
	if err != nil {
		return GetMsgData{}, err
	}
	return Data, nil
}

// GetForwardMsgData get_forward_msg 接口data统一结构体
type GetForwardMsgData struct {
	Messages []struct {
		Message []struct {
			Type string                 `json:"type"`
			Data map[string]interface{} `json:"data"`
		} `json:"message"`
	} `json:"messages"`
}

// GetForwardMsg get_forward_msg
func (a *Adapter) GetForwardMsg(iD string) (GetForwardMsgData, error) {
	params := map[string]interface{}{
		"message_id": iD,
	}
	res, err := a.CallAction("get_forward_msg", params)
	if err != nil {
		return GetForwardMsgData{}, err
	}
	var Data GetForwardMsgData
	err = json.Unmarshal(res.Data, &Data)
	return Data, nil

}

// SendLike send_like 发送好友赞
func (a *Adapter) SendLike(userID int64, times int) error {
	params := map[string]interface{}{
		"user_id": userID,
		"times":   times,
	}
	_, err := a.CallAction("send_like", params)
	return err
}

// SetGroupKick set_group_kick 群组踢人
func (a *Adapter) SetGroupKick(groupID int64, userID int64, rejectAddRequest bool) error {
	params := map[string]interface{}{
		"user_id":            userID,
		"reject_add_request": rejectAddRequest,
		"group_id":           groupID,
	}
	_, err := a.CallAction("set_group_kick", params)
	return err
}

// SetGroupBan set_group_ban 群组单人禁言
func (a *Adapter) SetGroupBan(groupID int64, userID int64, duration int64) error {
	params := map[string]interface{}{
		"user_id":  userID,
		"duration": duration,
		"group_id": groupID,
	}
	_, err := a.CallAction("set_group_ban", params)
	return err
}

// SetGroupAnonymousBan set_group_anonymous_ban 群组匿名用户禁言
func (a *Adapter) SetGroupAnonymousBan(groupID int64, anonymous any, anonymousFlag string, duration int64) error {
	params := map[string]interface{}{}
	if anonymous != nil {
		params = map[string]interface{}{
			"anonymous": anonymous,
			"duration":  duration,
			"group_id":  groupID,
		}
	} else {
		params = map[string]interface{}{
			"duration":      duration,
			"group_id":      groupID,
			"anonymousFlag": anonymousFlag,
		}
	}

	_, err := a.CallAction("set_group_anonymous_ban", params)
	return err
}

// SetGroupWholeBan set_group_whole_ban 群组全员禁言
func (a *Adapter) SetGroupWholeBan(groupID int64, enable bool) error {
	params := map[string]interface{}{
		"enable":   enable,
		"group_id": groupID,
	}
	_, err := a.CallAction("set_group_whole_ban", params)
	return err
}

// SetGroupAdmin set_group_admin 群组设置管理员
func (a *Adapter) SetGroupAdmin(groupID int64, userId int64, enable bool) error {
	params := map[string]interface{}{
		"enable":   enable,
		"user_id":  userId,
		"group_id": groupID,
	}
	_, err := a.CallAction("set_group_admin", params)
	return err
}

// SetGroupAnonymous set_group_anonymous 群组匿名
func (a *Adapter) SetGroupAnonymous(groupID int64, enable bool) error {
	params := map[string]interface{}{
		"enable":   enable,
		"group_id": groupID,
	}
	_, err := a.CallAction("set_group_anonymous", params)
	return err
}

// SetGroupCard set_group_card 设置群名片（群备注）
// group_id	number	-	群号
// user_id	number	-	要设置的 QQ 号
// card	string	空	群名片内容，不填或空字符串表示删除群名片
func (a *Adapter) SetGroupCard(groupID int64, userID int64, card string) error {
	params := map[string]interface{}{
		"user_id":  userID,
		"group_id": groupID,
		"card":     card,
	}
	_, err := a.CallAction("set_group_card", params)
	return err
}

// SetGroupName set_group_name 设置群名
// group_id	number (int64)	群号
// group_name	string	新群名
func (a *Adapter) SetGroupName(groupID int64, groupName string) error {
	params := map[string]interface{}{
		"group_id":   groupID,
		"group_name": groupName,
	}
	_, err := a.CallAction("set_group_card", params)
	return err
}

// SetGroupLeave 退出群组
// group_id	number	-	群号
// is_dismiss	boolean	false	是否解散，如果登录号是群主，则仅在此项为 true 时能够解散
func (a *Adapter) SetGroupLeave(groupID int64, isDismiss bool) error {
	params := map[string]interface{}{
		"is_dismiss": isDismiss,
		"group_id":   groupID,
	}
	_, err := a.CallAction("set_group_leave", params)
	return err
}

// SetGroupSpecialTitle 设置群组专属头衔
// group_id	number	-	群号
// user_id	number	-	要设置的 QQ 号
// special_title	string	空	专属头衔，不填或空字符串表示删除专属头衔
// duration	number	-1	专属头衔有效期，单位秒，-1 表示永久，不过此项似乎没有效果，可能是只有某些特殊的时间长度有效，有待测试
func (a *Adapter) SetGroupSpecialTitle(groupID int64, userID int64, specialTitle string, duration int64) error {
	params := map[string]interface{}{
		"user_id":       userID,
		"group_id":      groupID,
		"special_title": specialTitle,
		"duration":      duration,
	}
	_, err := a.CallAction("set_group_member", params)
	return err
}

// SetFriendAddRequest 处理加好友请求
// flag	string	-	加好友请求的 flag（需从上报的数据中获得）
// approve	boolean	true	是否同意请求
// remark	string	空	添加后的好友备注（仅在同意时有效）
func (a *Adapter) SetFriendAddRequest(flag string, approve bool, remark string) error {
	params := map[string]interface{}{
		"approve": approve,
		"flag":    flag,
		"remark":  remark,
	}
	_, err := a.CallAction("set_friend_add_request", params)
	return err
}

// SetGroupAddRequest 处理加群请求／邀请
// flag      string  -     加群请求的 flag（需从上报的数据中获得）
// subType   string  -     add 或 invite，请求类型（需要和上报消息中的 sub_type 字段相符）
// approve   bool    true  是否同意请求／邀请
// reason    string  空    拒绝理由（仅在拒绝时有效）
func (a *Adapter) SetGroupAddRequest(flag string, subType string, approve bool, reason string) error {
	params := map[string]interface{}{
		"flag":     flag,
		"sub_type": subType,
		"approve":  approve,
		"reason":   reason,
	}
	_, err := a.CallAction("set_group_add_request", params)
	return err
}

// GetLoginInfoData 登录号信息结构体
type GetLoginInfoData struct {
	UserID   int64  `json:"user_id"`
	NickName string `json:"nickname"`
}

// GetLoginInfo 获取登录号结构体
func (a *Adapter) GetLoginInfo() (GetLoginInfoData, error) {
	res, err := a.CallAction("get_login_info", nil)
	var data GetLoginInfoData
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(res.Data, &data)
	return data, err
}

// GetStrangerInfoData 陌生人信息结构体
type GetStrangerInfoData struct {
	UserID   int64  `json:"user_id"`
	NickName string `json:"nickname"`
	Sex      string `json:"sex"`
	Age      int32  `json:"age"`
}

// GetStrangerInfo 获取陌生人信息
// user_id	number	-	QQ 号
// no_cache	boolean	false	是否不使用缓存（使用缓存可能更新不及时，但响应更快）
func (a *Adapter) GetStrangerInfo(userID int64, noCache bool) (GetStrangerInfoData, error) {
	params := map[string]interface{}{
		"user_id":  userID,
		"no_cache": noCache,
	}
	res, err := a.CallAction("get_stranger_info", params)
	var data GetStrangerInfoData
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(res.Data, &data)
	return data, err
}

// GetFriendListData 好友列表结构体
type GetFriendListData struct {
	UserID   int64  `json:"user_id"`
	NickName string `json:"nickname"`
	Remark   string `json:"remark"`
}

// GetFriendList 获取好友列表
func (a *Adapter) GetFriendList() ([]GetFriendListData, error) {
	res, err := a.CallAction("get_friend_list", nil)
	var data []GetFriendListData
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(res.Data, &data)
	return data, err
}

// GetGroupInfoData 获取群信息结构体
// group_id	number (int64)	群号
// group_name	string	群名称
// member_count	number (int32)	成员数
// max_member_count	number (int32)	最大成员数（群容量）
type GetGroupInfoData struct {
	GroupID        int64  `json:"group_id"`
	GroupName      string `json:"group_name"`
	MemberCount    int32  `json:"member_count"`
	MaxMemberCount int32  `json:"max_member_count"`
}

// GetGroupInfo 获取群信息
// group_id	number	-	群号
// no_cache	boolean	false	是否不使用缓存（使用缓存可能更新不及时，但响应更快）
func (a *Adapter) GetGroupInfo(groupID int64, noCache bool) (GetGroupInfoData, error) {
	params := map[string]interface{}{
		"group_id": groupID,
		"no_cache": noCache,
	}
	res, err := a.CallAction("get_group_info", params)
	var data GetGroupInfoData
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(res.Data, &data)
	return data, err
}

// GetGroupList 获取群列表
func (a *Adapter) GetGroupList() ([]GetGroupInfoData, error) {
	res, err := a.CallAction("get_group_list", nil)
	var data []GetGroupInfoData
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(res.Data, &data)
	return data, err
}

// GetGroupMemberInfoData 获取群成员信息结构体
// group_id	number (int64)	群号
// user_id	number (int64)	QQ 号
// nickname	string	昵称
// card	string	群名片／备注
// sex	string	性别，male 或 female 或 unknown
// age	number (int32)	年龄
// area	string	地区
// join_time	number (int32)	加群时间戳
// last_sent_time	number (int32)	最后发言时间戳
// level	string	成员等级
// role	string	角色，owner 或 admin 或 member
// unfriendly	boolean	是否不良记录成员
// title	string	专属头衔
// title_expire_time	number (int32)	专属头衔过期时间戳
// card_changeable	boolean	是否允许修改群名片
type GetGroupMemberInfoData struct {
	GroupID         int64  `json:"group_id"`
	NickName        string `json:"nickname"`
	Card            string `json:"card"`
	Sex             string `json:"sex"`
	Age             int32  `json:"age"`
	Area            string `json:"area"`
	JoinTime        int32  `json:"join_time"`
	LastSentTime    int32  `json:"last_sent_time"`
	Level           string `json:"level"`
	Role            string `json:"role"`
	Unfriendly      bool   `json:"unfriendly"`
	Title           string `json:"title"`
	Tile            string `json:"tile"`
	TitleExpireTime int32  `json:"title_expire_time"`
	CardChangeable  bool   `json:"card_changeable"`
}

// GetGroupMemberInfo 获取群成员信息
// group_id	number	-	群号
// user_id	number	-	QQ 号
// no_cache	boolean	false	是否不使用缓存（使用缓存可能更新不及时，但响应更快）
func (a *Adapter) GetGroupMemberInfo(groupID int64, noCache bool, userID int64) (GetGroupMemberInfoData, error) {
	params := map[string]interface{}{
		"group_id": groupID,
		"no_cache": noCache,
		"user_id":  userID,
	}
	res, err := a.CallAction("get_group_member_info", params)
	var data GetGroupMemberInfoData
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(res.Data, &data)
	return data, err
}

// GetGroupMemberList 获取群成员列表
// group_id	number (int64)	-	群号
// type	string	-	要获取的群荣誉类型，可传入 talkative performer legend strong_newbie emotion 以分别获取单个类型的群荣誉数据，或传入 all 获取所有数据
func (a *Adapter) GetGroupMemberList() ([]GetGroupMemberInfoData, error) {
	res, err := a.CallAction("get_group_member_list", nil)
	var data []GetGroupMemberInfoData
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(res.Data, &data)
	return data, err
}

// GetGroupHonorInfoData 获取群荣誉信息结构体
// group_id	number (int64)	群号
// current_talkative	object	当前龙王，仅 type 为 talkative 或 all 时有数据
// talkative_list	array	历史龙王，仅 type 为 talkative 或 all 时有数据
// performer_list	array	群聊之火，仅 type 为 performer 或 all 时有数据
// legend_list	array	群聊炽焰，仅 type 为 legend 或 all 时有数据
// strong_newbie_list	array	冒尖小春笋，仅 type 为 strong_newbie 或 all 时有数据
// emotion_list	array	快乐之源，仅 type 为 emotion 或 all 时有数据
// 其中 current_talkative 字段的内容如下：
// 字段名	数据类型	说明
// user_id	number (int64)	QQ 号
// nickname	string	昵称
// avatar	string	头像 URL
// day_count	number (int32)	持续天数
// 其它各 *_list 的每个元素是一个 JSON 对象，内容如下：
// 字段名	数据类型	说明
// user_id	number (int64)	QQ 号
// nickname	string	昵称
// avatar	string	头像 URL
// description	string	荣誉描述
type GetGroupHonorInfoData struct {
	GroupID          int64                   `json:"group_id"`
	CurrentTalkative CurrentTalkative        `json:"current_talkative"`
	TtalkativeList   []GetGroupHonorInfolist `json:"talkative_list"`
	PerformerList    []GetGroupHonorInfolist `json:"performer_list"`
	LegendList       []GetGroupHonorInfolist `json:"legend_list"`
	StrongNewbieList []GetGroupHonorInfolist `json:"strong_newbie_list"`
	EmotionList      []GetGroupHonorInfolist `json:"emotion_list"`
}

// CurrentTalkative 当前龙王
// user_id	number (int64)	QQ 号
// nickname	string	昵称
// avatar	string	头像 URL
// day_count	number (int32)	持续天数
type CurrentTalkative struct {
	user_id  int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Avatar   string
	DayCount int32 `json:"day_count"`
}

// GetGroupHonorInfolist  获取群荣誉信息:历史龙王等数据的类型
// user_id	number (int64)	QQ 号
// nickname	string	昵称
// avatar	string	头像 URL
// description	string	荣誉描述
type GetGroupHonorInfolist struct {
	UserID      int64  `json:"user_id"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
}

// GetCookies 获取 Cookies
// domain	string	空	需要获取 cookies 的域名
func (a *Adapter) GetCookies(domain string) (string, error) {
	res, err := a.CallAction("get_cookies", nil)
	var cookies string
	if err != nil {
		return cookies, err
	}
	err = json.Unmarshal(res.Data, &cookies)
	return cookies, err
}

// GetCsrfToken 获取 CSRF Token
func (a *Adapter) GetCsrfToken() int32 {
	res, err := a.CallAction("get_csrf_token", nil)
	var token int32
	if err != nil {
		return token
	}
	err = json.Unmarshal(res.Data, &token)
	return token
}

// GetCredentialsData 获取 QQ 相关接口凭证 结构体
// cookies	string	Cookies
// csrf_token	number (int32)	CSRF Token
type GetCredentialsData struct {
	Cookies   string `json:"cookies"`
	CsrfToken int32  `json:"csrf_token"`
}

// GetCredentials 获取 QQ 相关接口凭证 即上面两个接口的合并。
// domain	string	空	需要获取 cookies 的域名
func (a *Adapter) GetCredentials(domain string) GetCredentialsData {
	params := map[string]interface{}{
		"domain": domain,
	}
	res, err := a.CallAction("get_credentials", params)
	var data GetCredentialsData
	if err != nil {
		return data
	}
	err = json.Unmarshal(res.Data, &data)
	return data
}

// GetRecord 获取语音
// 提示：要使用此接口，通常需要安装 ffmpeg，请参考 OneBot 实现的相关说明。
// file	string	-	收到的语音文件名（消息段的 file 参数），如 0B38145AA44505000B38145AA4450500.silk
// out_format	string	-	要转换到的格式，目前支持 mp3、amr、wma、m4a、spx、ogg、wav、flac
func (a *Adapter) GetRecord(file string, outFormat string) string {
	params := map[string]interface{}{
		"file":       file,
		"out_format": outFormat,
	}
	res, err := a.CallAction("get_record", params)
	var record string
	if err != nil {
		return record
	}
	err = json.Unmarshal(res.Data, &record)
	return record
}

// GetImage 获取图片
// file	string	-	收到的图片文件名（消息段的 file 参数），如 6B4DE3DFD1BD271E3297859D41C530F5.jpg
func (a *Adapter) GetImage(file string) string {
	params := map[string]interface{}{
		"file": file,
	}
	res, err := a.CallAction("get_image", params)
	var image string
	if err != nil {
		return image
	}
	err = json.Unmarshal(res.Data, &image)
	return image

}

// Cansendimage 检查是否可以发送图片
func (a *Adapter) Cansendimage() bool {
	res, err := a.CallAction("can_send_image", nil)
	var canSendImage bool
	if err != nil {
		return canSendImage
	}
	err = json.Unmarshal(res.Data, &canSendImage)
	return canSendImage
}

// GetStatusData 获取运行状态结构体
// online	boolean	当前 QQ 在线，null 表示无法查询到在线状态
// good	boolean	状态符合预期，意味着各模块正常运行、功能正常，且 QQ 在线
type GetStatusData struct {
	Online bool `json:"online"`
	Good   bool `json:"good"`
}

// GetStatus 获取运行状态
func (a *Adapter) GetStatus() (GetStatusData, error) {
	res, err := a.CallAction("get_status", nil)
	var status GetStatusData
	if err != nil {
		return status, err
	}
	err = json.Unmarshal(res.Data, &status)
	return status, err
}

// GetVersionInfoData 获取版本信息结构体
// app_name	string	应用标识，如 mirai-native
// app_version	string	应用版本，如 1.2.3
// protocol_version	string	OneBot 标准版本，如 v11
// ……	-	OneBot 实现自行添加的其它内容
type GetVersionInfoData struct {
	AppName         string `json:"app_name"`
	AppVersion      string `json:"app_version"`
	ProtocolVersion string `json:"protocol_version"`
}

// SetRestart 重启 OneBot 实现
// 由于重启 OneBot 实现同时需要重启 API 服务，这意味着当前的 API 请求会被中断，因此需要异步地重启，接口返回的 status 是 async。
func (a *Adapter) SetRestart() error {
	_, err := a.CallAction("set_restart", nil)
	if err != nil {
		return err
	}
	return nil
}

// CleanCache 清理缓存
// 用于清理积攒了太多的缓存文件。
func (a *Adapter) CleanCache() error {
	_, err := a.CallAction("clea_cahe", nil)
	if err != nil {
		return err
	}
	return nil
}
