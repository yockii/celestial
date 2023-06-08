package provider

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/asset/model"
	"io"
)

type Minio struct {
	model.OssConfig
	Client *minio.Client
}

func (o *Minio) GetOssConfigID() uint64 {
	return o.OssConfig.ID
}

func (o *Minio) Auth() error {
	client, err := minio.New(o.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(o.AccessKeyID, o.SecretAccessKey, ""),
		Region: o.Region,
		Secure: o.Secure == 1,
	})
	if err != nil {
		logger.Error(err)
		return err
	}

	found, err := client.BucketExists(context.Background(), o.Bucket)
	if err != nil {
		logger.Error(err)
		return err
	}
	if !found {
		err = client.MakeBucket(context.Background(), o.Bucket, minio.MakeBucketOptions{Region: o.Region})
		if err != nil {
			logger.Error(err)
			return err
		}
	}
	o.Client = client
	return nil
}

func (o *Minio) Close() error {
	return nil
}

func (o *Minio) PutObject(objName string, reader io.Reader) error {
	_, err := o.Client.PutObject(context.Background(), o.Bucket, objName, reader, -1, minio.PutObjectOptions{})
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (o *Minio) GetObject(objName string) (io.ReadCloser, error) {
	object, err := o.Client.GetObject(context.Background(), o.Bucket, objName, minio.GetObjectOptions{})
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return object, nil
}
