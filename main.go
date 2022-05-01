package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

const (
	envName_AK = "QINIU_ACCESS_KEY"
	envName_SK = "QINIU_SECRET_KEY"

	keyPrefix = "typora/"
)

var (
	bucket string
	host   string
)

func init() {
	flag.StringVar(&bucket, "bucket", "", "存储桶")
	flag.StringVar(&host, "host", "", "文件访问 Host")
}

// typora-uploader -bucket=shine-doc -host=http://shine-doc.lixinshow.top --
func main() {
	flag.Parse()

	if bucket == "" {
		errorAndExit("未设置存储桶，需要通过参数 -bucket 指定")
	}
	if host == "" {
		errorAndExit("未设置文件访问 Host，需要通过参数 -host 指定")
	}
	validateHost(host)
	host = strings.TrimRight(host, "/")

	filenames := parseFilenames()
	if len(filenames) == 0 {
		errorAndExit(`未解析到待上传文件，请确认命令行 "--" 之后的参数存在且为本地文件`)
	}

	policy := storage.PutPolicy{
		Scope: bucket,
	}
	upToken := policy.UploadToken(getCredentials())

	uploader := storage.NewFormUploader(&storage.Config{})

	year := time.Now().Local().Format("2006/")

	for _, file := range filenames {
		data, err := os.ReadFile(file)
		if err != nil {
			errorAndExit(fmt.Sprintf("读取文件失败 %s %s\n", file, err))
		}
		_md5 := md5.Sum(data)
		key := keyPrefix + year + hex.EncodeToString(_md5[:])

		ret := storage.PutRet{}

		err = uploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(data), int64(len(data)), nil)
		if err != nil {
			errorAndExit(fmt.Sprintf("上传失败 %s %s", file, err))
		}

		fmt.Printf("%s/%s\n", host, ret.Key)
	}
}

func getCredentials() *auth.Credentials {
	ak := os.Getenv(envName_AK)
	sk := os.Getenv(envName_SK)
	if ak == "" || sk == "" {
		errorAndExit(fmt.Sprintf("需要通过环境变量配置AK,SK (环境变量名称：%s %s)", envName_AK, envName_SK))
	}

	return auth.New(ak, sk)
}

func parseFilenames() []string {
	var i int
	for idx, v := range os.Args {
		if v == "--" {
			i = idx
			break
		}
	}
	return os.Args[i+1:]
}

func validateHost(host string) {
	_, err := url.Parse(host)
	if err != nil {
		errorAndExit("文件访问 Host 格式有误，请检查")
	}
}

func errorAndExit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
