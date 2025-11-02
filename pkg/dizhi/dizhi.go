package dizhi

import (
	"qimenpan/pkg/wuxing"
)

// DiZhi 地支
type DiZhi struct {
	Id            int           `json:"id"`   // 序数
	Name          string        `json:"name"` // 名称
	wuxing.WuXing               // 五行
	ShiLinWuXing  wuxing.WuXing // 时令五行
	YueJiang      string        `json:"yueJiang"`
	JianXing      string        `json:"jianXing"`
}

var (
	ZI   = DiZhi{Id: 1, Name: "子", WuXing: wuxing.SHUI, ShiLinWuXing: wuxing.SHUI}
	CHOU = DiZhi{Id: 2, Name: "丑", WuXing: wuxing.TU, ShiLinWuXing: wuxing.TU}
	YIN  = DiZhi{Id: 3, Name: "寅", WuXing: wuxing.MU, ShiLinWuXing: wuxing.MU}
	MAO  = DiZhi{Id: 4, Name: "卯", WuXing: wuxing.MU, ShiLinWuXing: wuxing.MU}
	CHEN = DiZhi{Id: 5, Name: "辰", WuXing: wuxing.TU, ShiLinWuXing: wuxing.TU}
	SI   = DiZhi{Id: 6, Name: "巳", WuXing: wuxing.HUO, ShiLinWuXing: wuxing.HUO}
	WU   = DiZhi{Id: 7, Name: "午", WuXing: wuxing.MU, ShiLinWuXing: wuxing.HUO}
	WEI  = DiZhi{Id: 8, Name: "未", WuXing: wuxing.TU, ShiLinWuXing: wuxing.TU}
	SHEN = DiZhi{Id: 9, Name: "申", WuXing: wuxing.JIN, ShiLinWuXing: wuxing.JIN}
	YOU  = DiZhi{Id: 10, Name: "酉", WuXing: wuxing.JIN, ShiLinWuXing: wuxing.JIN}
	XU   = DiZhi{Id: 11, Name: "戌", WuXing: wuxing.TU, ShiLinWuXing: wuxing.TU}
	HAI  = DiZhi{Id: 12, Name: "亥", WuXing: wuxing.SHUI, ShiLinWuXing: wuxing.SHUI}
)
var ALL = []DiZhi{ZI, CHOU, YIN, MAO, CHEN, SI, WU, WEI, SHEN, YOU, XU, HAI}
var YUE_JIANG = []string{"登明", "河魁", "从魁", "传送", "小吉", "胜光", "太乙", "天罡", "太冲", "功曹", "大吉", "神后"}
var Jianxing = []string{"建", "除", "满", "平", "定", "执", "破", "危", "成", "收", "开", "闭"}

var data = map[string]DiZhi{
	"子": ZI,
	"丑": CHOU,
	"寅": YIN,
	"卯": MAO,
	"辰": CHEN,
	"巳": SI,
	"午": WU,
	"未": WEI,
	"申": SHEN,
	"酉": YOU,
	"戌": XU,
	"亥": HAI,
}
var yuejiang = map[int]string{
	ZI.Id:   "神后",
	CHOU.Id: "大吉",
	YIN.Id:  "功曹",
	MAO.Id:  "太冲",
	CHEN.Id: "天罡",
	SI.Id:   "太乙",
	WU.Id:   "胜光",
	WEI.Id:  "小吉",
	SHEN.Id: "传送",
	YOU.Id:  "从魁",
	XU.Id:   "河魁",
	HAI.Id:  "登明",
}

func Match(z string) DiZhi {
	data, ok := data[z]
	if !ok {
		panic("地支匹配失败")
	}
	return data
}

func GetYueJiang(id int) string {
	return yuejiang[id]
}
func Index(index int) DiZhi {
	i := index - 1
	if i < 0 {
		i = 12 + i
	}
	return ALL[i]
}

type He struct {
	wuxing.WuXing
	He []DiZhi
}

func (h *He) Equal(he He) bool {
	return h.WuXing.GetName() == he.WuXing.GetName()
}

var (
	Shui = He{wuxing.SHUI, []DiZhi{SHEN, ZI, CHEN}}
	Mu   = He{wuxing.MU, []DiZhi{HAI, MAO, WEI}}
	Huo  = He{wuxing.HUO, []DiZhi{YIN, WU, XU}}
	Jin  = He{wuxing.JIN, []DiZhi{SI, YOU, CHOU}}
)

func SanHe(dz DiZhi) He {
	if dz == ZI || dz == CHEN || dz == SHEN {
		return Shui
	}
	if dz == HAI || dz == MAO || dz == WEI {
		return Mu
	}
	if dz == YIN || dz == WU || dz == XU {
		return Huo
	}
	return Jin
}
