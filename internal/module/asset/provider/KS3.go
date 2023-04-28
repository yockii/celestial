package provider

import "github.com/yockii/celestial/internal/module/asset/model"

type KS3 struct {
	model.OssConfig
}

func (o *KS3) Auth() error {
	return nil
}
