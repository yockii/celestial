package provider

import "github.com/yockii/celestial/internal/module/asset/model"

type Azure struct {
	model.OssConfig
}

func (o *Azure) Auth() error {
	return nil
}
