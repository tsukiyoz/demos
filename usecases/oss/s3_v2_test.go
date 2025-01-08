package oss

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestS3_V2(t *testing.T) {
	ctx := context.Background()

	accessKeyId, ok := os.LookupEnv("ACCESS_KEY_ID")
	if !ok {
		t.Fatal("not access_key_id found")
	}
	accessKeySecret, ok := os.LookupEnv("ACCESS_KEY_SECRET")
	if !ok {
		t.Fatal("no access_key_secret found")
	}

	baseEndPoint := "https://oss-cn-hangzhou.aliyuncs.com"
	client := s3.New(s3.Options{
		Region:       "oss-cn-hangzhou",
		BaseEndpoint: &baseEndPoint,
		Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
	})

	t.Logf("[id]: %v\n[secret]: %v\n", accessKeyId, accessKeySecret)

	//_, err := client.PutObject(ctx, &s3.PutObjectInput{
	//	Bucket:      GetStringPtr("tsukiyo"),
	//	Key:         GetStringPtr("webook/hello"),
	//	Body:        bytes.NewReader([]byte("hello, world :)")),
	//	ContentType: GetStringPtr("text/plain;charset=utf-8"),
	//})
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	//_, err := client.PutBucketAcl(ctx, &s3.PutBucketAclInput{
	//	Bucket: GetStringPtr("tsukiyo"),
	//	ACL:    types.BucketCannedACLPublicRead,
	//})
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	resp, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: GetStringPtr("tsukiyo"),
		Key:    GetStringPtr("webook/hello.txt"),
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	t.Log(string(data))
}

func GetStringPtr(str string) *string {
	s := str
	return &s
}
