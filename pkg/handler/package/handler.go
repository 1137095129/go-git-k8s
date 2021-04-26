package _package

import (
	"bufio"
	"fmt"
	"github.com/wang1137095129/go-git-k8s/config"
	"github.com/wang1137095129/go-git-k8s/utils"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

//DockerPackageStr eg:docker build -t [PATH] [IMAGE_NAME]
var DockerPackageStr = "docker build -t %s %s"

var once = &sync.Once{}

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

func DockerfileHandle() (string, error) {
	dockerfilePath := filepath.Join(utils.GetConfigDir(), "Dockerfile")
	_, err := os.Stat(dockerfilePath)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(dockerfilePath)
			if err != nil {
				return "", err
			}
			file.Close()
		} else {
			return "", err
		}
	}
	return dockerfilePath, nil
}

func DockerTagName(c *config.Config) string {
	return fmt.Sprintf("%s/%s:%s", c.Docker.Username, c.Docker.Repository, utils.DateFormat(time.Now()))
}

func DockerfileCreate(c *config.Config, dockerfileTemplate string) (string, error) {
	path := filepath.Join(c.Git.Local, c.Git.Repository)
	pkgStr := fmt.Sprintf(dockerfileTemplate, path, ExposeStr(c), c.Git.Repository)
	dockerfilePath, err := DockerfileHandle()
	if err != nil {
		return "", err
	}
	file, err := os.OpenFile(dockerfilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	_, err = writer.Write(utils.StringToBytes(pkgStr))
	if err != nil {
		return "", err
	}
	return dockerfilePath, nil
}

func DockerPackage(c *config.Config, dockerfilePath string) (string, error) {
	once.Do(func() {
		cmd := exec.Command("docker", "login", "-u", c.Docker.Username, "-p", c.Docker.Password)
		err := cmd.Run()
		utils.CheckError(err)
	})
	tagName := DockerTagName(c)
	cmd := exec.Command("docker", "build", "-t", tagName, "-f", dockerfilePath)
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return tagName, err
}
