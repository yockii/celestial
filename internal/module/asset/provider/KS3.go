package provider

import (
	"github.com/yockii/celestial/internal/module/asset/model"
	"io"
)

type KS3 struct {
	model.OssConfig
}

func (o *KS3) Auth() error {
	return nil
}

func (o *KS3) Close() error {
	return nil
}

func (o *KS3) PutObject(objName string, reader io.Reader) error {
	return nil
}

func (o *KS3) GetObject(objName string) (io.ReadCloser, error) {
	return nil, nil
}
