package maven

import (
	"bufio"
	"fmt"
	"github.com/wang1137095129/go-git-k8s/config"
	_package "github.com/wang1137095129/go-git-k8s/pkg/handler/package"
	"github.com/wang1137095129/go-git-k8s/utils"
	"os"
	"path/filepath"
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
	path := filepath.Join(c.Git.Local, c.Git.Repository)
	pkgStr := fmt.Sprintf(packageStr, path, _package.ExposeStr(c), c.Git.Repository)
	dockerfilePath, err := _package.DockerfileHandle()
	utils.CheckError(err)
	file, err := os.OpenFile(dockerfilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	utils.CheckError(err)
	defer file.Close()
	writer := bufio.NewWriter(file)
	_, err = writer.Write(utils.StringToBytes(pkgStr))
	utils.CheckError(err)
	return "", nil
}
