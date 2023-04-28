package provider

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/internal/module/asset/model"
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
