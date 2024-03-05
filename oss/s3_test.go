package oss

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"
)

// 你可以用这个来单独测试你的 OSS 配置对不对，有没有权限
func TestS3(t *testing.T) {
	// 腾讯云中对标 s3 和 OSS 的产品叫做 COS
	cosId, ok := os.LookupEnv("ACCESS_KEY_ID")
	if !ok {
		panic("没有找到环境变量 ACCESS_KEY_ID")
	}
	cosKey, ok := os.LookupEnv("ACCESS_KEY_SECRET")
	if !ok {
		panic("没有找到环境变量 ACCESS_KEY_SECRET")
	}
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(cosId, cosKey, ""),
		Region:      ekit.ToPtr[string]("oss-cn-hangzhou"),
		Endpoint:    ekit.ToPtr[string]("https://oss-cn-hangzhou.aliyuncs.com"),
		// 强制使用 /bucket/key 的形态
		S3ForcePathStyle: ekit.ToPtr[bool](true),
	})
	assert.NoError(t, err)
	client := s3.New(sess)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err = client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      ekit.ToPtr[string]("tsukiyo"),
		Key:         ekit.ToPtr[string]("webook/test-v1.txt"),
		Body:        bytes.NewReader([]byte("测试内容 abc")),
		ContentType: ekit.ToPtr[string]("text/plain;charset=utf-8"),
	})
	assert.NoError(t, err)
	res, err := client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: ekit.ToPtr[string]("tsukiyo"),
		Key:    ekit.ToPtr[string]("webook/test-v1.txt"),
	})
	assert.NoError(t, err)
	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	t.Log(string(data))
}
