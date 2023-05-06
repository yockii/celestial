package provider

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/asset/model"
	"io"
)

type OBS struct {
	model.OssConfig
	Client *obs.ObsClient
}

func (o *OBS) Auth() error {
	client, err := obs.New(o.AccessKeyID, o.SecretAccessKey, o.Endpoint, obs.WithCustomDomainName(o.SelfDomain == 1), obs.WithRegion(o.Region))
	if err != nil {
		logger.Error(err)
		return err
	}
	_, err = client.GetBucketLocation(o.BucketName)
	if err != nil {
		logger.Error(err)
		return err
	}

	o.Client = client
	return nil
}

func (o *OBS) Close() error {
	o.Client.Close()
	return nil
}

func (o *OBS) PutObject(objName string, reader io.Reader) error {
	input := &obs.PutObjectInput{}
	input.Bucket = o.BucketName
	input.Key = objName
	input.Body = reader

	_, err := o.Client.PutObject(input)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (o *OBS) GetObject(objName string) (io.ReadCloser, error) {
	input := &obs.GetObjectInput{}
	input.Bucket = o.BucketName
	input.Key = objName

	output, err := o.Client.GetObject(input)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return output.Body, nil
}
