package provider

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/asset/model"
)

type OSS struct {
	model.OssConfig
	Client *oss.Client
	Bucket *oss.Bucket
}

func (o *OSS) Auth() error {
	client, err := oss.New(o.Endpoint, o.AccessKeyID, o.SecretAccessKey, oss.UseCname(o.SelfDomain == 1))
	if err != nil {
		logger.Error(err)
		return err
	}

	var bucket *oss.Bucket
	bucket, err = client.Bucket(o.BucketName)
	if err != nil {
		logger.Error(err)
		return err
	}
	o.Client = client
	o.Bucket = bucket
	return nil
}
