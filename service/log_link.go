package service

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"time"
	"we2log/config"
	"we2log/model/_ssh"
	"we2log/model/log"
)

// 消息map
var (
	MsgChanMap map[string]chan string
	groups     []string
	GroupNum   int
	sshClients []*_ssh.Client
)

// CreateSshClient 创建ssh链接
func CreateSshClient() {
	// 所有开启的分组
	groups = GetSshGroups()
	GroupNum = len(groups)
	MsgChanMap = make(map[string]chan string, GroupNum)
	// 创建链接并推送
	for _, group := range *config.Yml.Log.Group {
		if !group.OnOff {
			continue
		}
		for _, conf := range *group.Ssh {
			// 如果未开启并且分组未开启则跳过
			if !conf.OnOff {
				continue
			}
			// 创建shh客户端
			client := _ssh.BuildClient(&_ssh.Config{
				Name:           conf.Name,
				Host:           conf.Host,
				Port:           conf.Port,
				Username:       conf.Username,
				Password:       conf.Password,
				PrivateKeyPath: conf.PriKeyPath,
				PasswordType:   conf.PwType,
			})
			sshClients = append(sshClients, client)
			// 执行日志查看命令
			go func() {
				err := client.Session.Run(fmt.Sprintf("tail -f %s", conf.LogPath))
				if err != nil {
					log.Warn(fmt.Sprintf("ssh 链接发送命令错误或已停止: %s", err))
				}
			}()
			// 获取数据流
			out, err := client.Session.StdoutPipe()
			if err != nil {
				log.Fatal("ssh 链接错误")
			}
			// 推送数据到通道
			go pushMsg(conf.Name, &out, group.Name)
			// 延迟10毫秒防止高并发问题
			time.Sleep(time.Millisecond * 10)
		}
	}
}

// GetSshGroups 获取ssh所有分组
func GetSshGroups() []string {
	groupConf := *config.Yml.Log.Group
	grs := make([]string, 0, len(groupConf))
	for _, group := range groupConf {
		if group.OnOff {
			grs = append(grs, group.Name)
		}
	}
	return grs
}

// 推送消息
func pushMsg(name string, out *io.Reader, group string) {
	msgChan := MsgChanMap[group]
	if msgChan == nil {
		msgChan = make(chan string, 10)
		MsgChanMap[group] = msgChan
	}
	reader := bufio.NewReader(*out)
	for {
		// 最长1096
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Warn(fmt.Sprintf("ssh推送消息错误或已停止: %v", err))
			return
		}
		mc := regexp.MustCompile("[\n\t]")
		bytes := mc.ReplaceAll([]byte(s), []byte(""))
		s = string(bytes)
		if s == "" {
			continue
		}
		msgChan <- fmt.Sprintf("%s: %s", name, s)
	}
}

// CloseSshClient 关闭所有ssh链接
func CloseSshClient() {
	for _, client := range sshClients {
		_ = client.Close()
	}
}
