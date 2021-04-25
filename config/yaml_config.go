package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"we2log/model/log"
)

var Yml *YmlConfig

type YmlConfig struct {
	Log *LogConfig `yaml:"log"`
}

type LogConfig struct {
	Lines    int      `yaml:"lines"`
	FontSize float32  `yaml:"font-size"`
	Group    *[]Group `yaml:"group"`
}

type Group struct {
	Name  string `yaml:"name"`
	OnOff bool   `yaml:"on-off"`
	Ssh   *[]Ssh `yaml:"ssh"`
}

type Ssh struct {
	Name       string `yaml:"name"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	PriKeyPath string `yaml:"pri-key-path"`
	PwType     int    `yaml:"pw-type"`
	LogPath    string `yaml:"log-path"`
	OnOff      bool   `yaml:"on-off"`
}

// InitYaml 初始化yaml
func InitYaml() {
	// 获取程序内部目录
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("项目路径获取错误")
	}
	appYmlPath := fmt.Sprintf("%s/.we2log/config.yml", dir)
	var yml = new(YmlConfig)
	readYml(appYmlPath, yml)
	Yml = yml
}

// 读取yml文件
func readYml(path string, conf *YmlConfig) {
	if !pathExists(path) {
		// 默认配置
		*conf = YmlConfig{}
		conf.Log = new(LogConfig)
		conf.Log.Lines = 100
		conf.Log.FontSize = 11
		conf.Log.Group = &[]Group{}
		return
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("yml文件读取错误: %s", err))
	}

	err = yaml.Unmarshal(file, conf)
	if err != nil {
		log.Fatal(fmt.Sprintf("yml文件配置读取错误: %s", err))
	}
}

// SaveLocalCache 保存配置
func SaveLocalCache() {
	out, err := yaml.Marshal(Yml)
	if err != nil {
		log.Error(fmt.Sprintf("yaml保存出错: %s", err))
		return
	}
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Error(fmt.Sprintf("用户目录获取错误: %s", err))
		return
	}
	// 应用配置路径
	appPath := fmt.Sprintf("%s/.we2log", dir)
	if !pathExists(appPath) {
		err = os.Mkdir(appPath, os.ModePerm)
		if err != nil {
			log.Error(fmt.Sprintf("配置文件夹创建错误: %s", err))
			return
		}
	}
	// 配置路径
	configPath := fmt.Sprintf("%s/config.yml", appPath)
	if pathExists(configPath) {
		err = os.Remove(configPath)
		if err != nil {
			log.Error(fmt.Sprintf("配置文件删除错误: %s", err))
			return
		}
	}
	// 写入配置文件
	err = ioutil.WriteFile(configPath, out, os.ModePerm)
	if err != nil {
		log.Error(fmt.Sprintf("配置文件写入错误: %s", err))
	}
}

// 路径是否存在
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
