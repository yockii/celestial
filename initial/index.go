package initial

import (
	"time"

	"github.com/yockii/ruomu-core/util"
)

func init() {
	time.Local = time.FixedZone("CST", 8*3600)
	_ = util.InitNode(1)
}
