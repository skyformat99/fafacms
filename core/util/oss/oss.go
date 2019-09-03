package oss

import (
	"bytes"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type Key struct {
	Endpoint        string `json:"Endpoint"`
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	BucketName      string `json:"BucketName"`
}

func SaveFile(K Key, ObjectName string, raw []byte) error {
	// 创建OSSClient实例。
	client, err := oss.New(K.Endpoint, K.AccessKeyId, K.AccessKeySecret)
	if err != nil {
		return err
	}

	// 获取存储空间。
	bucket, err := client.Bucket(K.BucketName)
	if err != nil {
		return err
	}

	// 上传Byte数组。
	err = bucket.PutObject(ObjectName, bytes.NewReader(raw))
	if err != nil {
		return err
	}

	return nil
}
