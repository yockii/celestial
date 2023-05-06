package provider

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/asset/model"
	"io"
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

func (o *OSS) Close() error {
	return nil
}

func (o *OSS) PutObject(objName string, reader io.Reader) error {
	err := o.Bucket.PutObject(objName, reader)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (o *OSS) GetObject(objName string) (io.ReadCloser, error) {
	body, err := o.Bucket.GetObject(objName)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return body, nil
}
