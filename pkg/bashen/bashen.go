package bashen

import (
	"qimenpan/pkg/wuxing"
)

// Shen 天盘八神
type Shen struct {
	Name          string `json:"name"` // 神名
	wuxing.WuXing        // 五行
}

var (
	ZHI_FU   = Shen{Name: "值符", WuXing: wuxing.SHUI}
	TEN_SHE  = Shen{Name: "滕蛇", WuXing: wuxing.TU}
	TAI_YIN  = Shen{Name: "太阴", WuXing: wuxing.MU}
	LIU_HE   = Shen{Name: "六合", WuXing: wuxing.MU}
	BAIHU    = Shen{Name: "白虎", WuXing: wuxing.HUO}
	XUAN_WU  = Shen{Name: "玄武", WuXing: wuxing.TU}
	JIU_DI   = Shen{Name: "九地", WuXing: wuxing.JIN}
	JIU_TIAN = Shen{Name: "九天", WuXing: wuxing.JIN}

	ALL = []Shen{ZHI_FU, TEN_SHE, TAI_YIN, LIU_HE, BAIHU, XUAN_WU, JIU_DI, JIU_TIAN}
)
