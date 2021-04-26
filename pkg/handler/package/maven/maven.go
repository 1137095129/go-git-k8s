package maven

import (
	"github.com/wang1137095129/go-git-k8s/config"
	_package "github.com/wang1137095129/go-git-k8s/pkg/handler/package"
	"os/exec"
)

var packageStr = `
FROM maven:3.8.1-jdk-8 AS builder
ADD %s /default
WORKDIR /default
RUN mvn package

FROM java:8-jdk-alpine%s
--form=builder /default/target/%s /app.jar
ENTRYPOINT ["java","-jar","app.jar"]
`

type Maven struct {
}

func (m *Maven) Package(c *config.Config) (string, error) {
	dockerfilePath, err := _package.DockerfileCreate(c, packageStr)
	if err != nil {
		return "", err
	}
	tagName, err := _package.DockerPackage(c, dockerfilePath)
	if err != nil {
		return "", err
	}
	cmd := exec.Command("docker", "push", tagName)
	err = cmd.Run()
	if err != nil {
		return "", err
	}
	return tagName, nil
}
