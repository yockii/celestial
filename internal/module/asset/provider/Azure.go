package provider

import (
	"github.com/yockii/celestial/internal/module/asset/model"
	"io"
)

type Azure struct {
	model.OssConfig
}

func (o *Azure) GetOssConfigID() uint64 {
	return o.OssConfig.ID
}

func (o *Azure) Auth() error {
	return nil
}

func (o *Azure) Close() error {
	return nil
}

func (o *Azure) PutObject(objName string, reader io.Reader) error {
	return nil
}

func (o *Azure) GetObject(objName string) (io.ReadCloser, error) {
	return nil, nil
}
