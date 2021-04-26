package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
	"runtime"
	"time"
	"unsafe"
)

func CheckError(err error) {
	if err != nil {
		logrus.Fatal(err)
	}
}

//GetConfigDir 获取配置文件所在路径
func GetConfigDir() string {
	if dir := os.Getenv("GIT_K8S_CONFIG"); len(dir) > 0 {
		return dir
	}
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

func StringToBytes(s string) (byteArrayResult []byte) {
	stringHeader := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	byteArrayHeader := (*reflect.SliceHeader)(unsafe.Pointer(&byteArrayResult))
	byteArrayHeader.Data, byteArrayHeader.Len, byteArrayHeader.Cap = stringHeader.Data, stringHeader.Len, stringHeader.Len
	return
}

func BytesToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func DateFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}