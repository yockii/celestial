package provider

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/asset/model"
)

type Minio struct {
	model.OssConfig
	Client *minio.Client
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

	found, err := client.BucketExists(context.Background(), o.BucketName)
	if err != nil {
		logger.Error(err)
		return err
	}
	if !found {
		err = client.MakeBucket(context.Background(), o.BucketName, minio.MakeBucketOptions{Region: o.Region})
		if err != nil {
			logger.Error(err)
			return err
		}
	}
	o.Client = client
	return nil
}
