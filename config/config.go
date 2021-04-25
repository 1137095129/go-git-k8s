package config

import (
	"github.com/wang1137095129/go-git-k8s/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

//ConfigFileName 配置文件名称
var ConfigFileName = "git_k8s_config.yaml"

type Config struct {
	Build Build `json:"build" yaml:"build"`
	Git   Git   `json:"git" yaml:"git"`
}

type Build struct {
	//Kind 编译方式，目前只支持maven和go
	Kind        string `json:"kind" yaml:"kind,omitempty"`
	PackageName string `json:"packageName" yaml:"packageName,omitempty"`
	Expose      []int  `json:"expose" yaml:"expose,omitempty"`
}

//Git 监控的git仓库配置
type Git struct {
	//Url 监控的git仓库地址
	Url string `json:"url" yaml:"url,omitempty"`
	//Remote git远端名称
	Remote string `json:"remote" yaml:"remote,omitempty"`
	//Branch 监控的git远端分支名称
	Branch string `json:"branch" yaml:"branch,omitempty"`
	//Certificate 私钥路径
	Certificate string `json:"certificate" yaml:"certificate,omitempty"`
	//Username
	Username string `json:"username" yaml:"username,omitempty"`
	//Password
	Password string `json:"password" yaml:"password,omitempty"`
	//Local 克隆到本地的路径
	Local string `json:"local" yaml:"local,omitempty"`
	//Repository 本地仓库名称
	Repository string
}

//Load 从配置文件中加载配置
func (c *Config) Load() error {
	err := createIfNotExists()
	if err != nil {
		return err
	}
	path := filepath.Join(utils.GetConfigDir(), ConfigFileName)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if len(b) > 0 {
		if err := yaml.Unmarshal(b, c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) Write() error {
	path := filepath.Join(utils.GetConfigDir(), ConfigFileName)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2)
	return encoder.Encode(c)
}

func New() (*Config, error) {
	c := &Config{}
	if err := c.Load(); err != nil {
		return c, err
	}
	return c, nil
}

//createIfNotExists 如果配置文件不存在，则创建它
func createIfNotExists() error {
	path := filepath.Join(utils.GetConfigDir(), ConfigFileName)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(path)
			if err != nil {
				return err
			}
			file.Close()
		} else {
			return err
		}
	}
	return nil
}

