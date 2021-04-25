package _package

import (
	"fmt"
	"github.com/wang1137095129/go-git-k8s/config"
	"github.com/wang1137095129/go-git-k8s/utils"
	"os"
	"path/filepath"
)


type Handler interface {
	Package(c *config.Config) (string, error)
}

func ExposeStr(c *config.Config) string {
	var result = ""
	for _, port := range c.Build.Expose {
		result += fmt.Sprintf("\nEXPOSE:%d", port)
	}
	result += "\n"
	return result
}

func DockerfileHandle() (string,error) {
	dockerfilePath := filepath.Join(utils.GetConfigDir(),"Dockerfile")
	_, err := os.Stat(dockerfilePath)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(dockerfilePath)
			if err != nil {
				return "", err
			}
			file.Close()
		}else {
			return "", err
		}
	}
	return dockerfilePath, nil
}
