package bamen

import (
	"qimenpan/pkg/jiugong"
	"qimenpan/pkg/wuxing"
)

// Men 人盘八门
type Men struct {
	Name string
	jiugong.Gong
	wuxing.WuXing
}

var (
	XIU   = Men{Name: "休门", Gong: jiugong.KAN, WuXing: wuxing.SHUI}
	SHENG = Men{Name: "生门", Gong: jiugong.GEN, WuXing: wuxing.TU}
	SHANG = Men{Name: "伤门", Gong: jiugong.ZHEN, WuXing: wuxing.MU}
	DU    = Men{Name: "杜门", Gong: jiugong.XUN, WuXing: wuxing.MU}
	JING  = Men{Name: "景门", Gong: jiugong.LI, WuXing: wuxing.HUO}
	SI    = Men{Name: "死门", Gong: jiugong.KUN, WuXing: wuxing.TU}
	JING2 = Men{Name: "惊门", Gong: jiugong.DUI, WuXing: wuxing.JIN}
	KAI   = Men{Name: "开门", Gong: jiugong.QIAN, WuXing: wuxing.JIN}

	ALL = []Men{XIU, SHENG, SHANG, DU, JING, SI, JING2, KAI}
)
