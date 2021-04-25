package _ssh

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
)

// 连接状态
const (
	// StatusOff 关闭
	StatusOff = iota
	// StatusOnline 连接中
	StatusOnline
)

// 密码类型
const (
	// PasswordText 文本
	PasswordText = iota
	// PasswordPrivateKey 私钥
	PasswordPrivateKey
)

// Config ssh链接配置
type Config struct {
	// 连接名称
	Name string
	// 服务器
	Host string
	// 端口
	Port string
	// 用户名
	Username string
	// 密码或私钥密码
	Password string
	// 私钥路径
	PrivateKeyPath string
	// 密码类型
	PasswordType int
}

// Client ssh连接客户端
type Client struct {
	// 连接状态
	status int
	// 连接配置
	conf *Config
	// 连接session
	Session *ssh.Session
}

// BuildClient 建立ssh连接
func BuildClient(config *Config) *Client {
	session := buildSession(config)
	return &Client{
		status:  StatusOnline,
		conf:    config,
		Session: session,
	}
}

// 建立ssh连接session
func buildSession(config *Config) *ssh.Session {
	//ssh链接配置
	clientConfig := &ssh.ClientConfig{
		User:            config.Username,
		Auth:            getAuthMethod(config),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	//服务器地址
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	//建立连接客户端
	sshClient, err := ssh.Dial("tcp", addr, clientConfig)
	errLog(err, "ssh链接")
	session, err := sshClient.NewSession()
	errLog(err, "session")

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}

	return session
}

// 获取ssh验证方法
func getAuthMethod(config *Config) []ssh.AuthMethod {
	var authMethods []ssh.AuthMethod

	switch config.PasswordType {
	case PasswordText:
		if len(config.Password) != 0 {
			authMethods = []ssh.AuthMethod{ssh.Password(config.Password)}
		} else {
			errLog(errors.New("未输入"), "密码")
		}
	case PasswordPrivateKey:
		// 如果私钥地址为空则获取计算机默认路径文件
		if len(config.PrivateKeyPath) == 0 {
			//用户目录
			userDir, err := os.UserHomeDir()
			errLog(err, "用户目录")
			//默认私钥路径
			config.PrivateKeyPath = fmt.Sprintf("%s/.ssh/id_rsa", userDir)
		}
		// 私钥字节
		privateBytes := getFilePrivateKey(config.PrivateKeyPath)
		authMethods = getPublicKeyAuth(privateBytes, config.Password)
	}

	return authMethods
}

// 获取公钥Auth
func getPublicKeyAuth(privateKey []byte, password string) []ssh.AuthMethod {
	// 私钥签名
	var signer ssh.Signer
	var err error
	// 是否传递私钥密码
	if len(password) != 0 {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(password))
	} else {
		signer, err = ssh.ParsePrivateKey(privateKey)
	}
	errLog(err, "ssh签名错误")
	return []ssh.AuthMethod{ssh.PublicKeys(signer)}
}

// 获取私钥文件key
func getFilePrivateKey(path string) []byte {
	key, err := ioutil.ReadFile(path)
	errLog(err, "ssh私钥获取错误")
	return key
}

// Close 关闭ssh连接
func (client *Client) Close() error {
	if client.status != StatusOnline {
		return errors.New("ssh隧道已经关闭")
	}
	client.status = StatusOff
	return client.Session.Close()
}

// Restart 重新启动
func (client *Client) Restart() {
	if client.status == StatusOnline {
		_ = client.Close()
	}
	client.Session = buildSession(client.conf)
}

func errLog(err error, msg string) {
	if err != nil {
		log.Fatalf("%s error: %v", msg, err)
	}
}
