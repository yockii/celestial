package provider

import (
	"github.com/yockii/celestial/internal/module/asset/model"
	"strings"
)

type OsManager interface {
	Auth() error
}

func GetProvider(c *model.OssConfig) (m OsManager) {
	switch strings.ToLower(c.Type) {
	case "minio":
		m = &Minio{
			OssConfig: *c,
		}
	case "ks3":
		m = &KS3{
			OssConfig: *c,
		}
	case "obs":
		m = &OBS{
			OssConfig: *c,
		}
	case "oss":
		m = &OSS{
			OssConfig: *c,
		}
	case "azure":
		m = &Azure{
			OssConfig: *c,
		}
	default:
		m = &OSS{
			OssConfig: *c,
		}
	}
	return
}
