package tiangan

import (
	"qimenpan/pkg/wuxing"
)

type TianGan struct {
	Id   int
	Name string
	wuxing.WuXing
}

var (
	JIA  = TianGan{Id: 1, Name: "甲", WuXing: wuxing.MU}
	YI   = TianGan{Id: 2, Name: "乙", WuXing: wuxing.MU}
	BING = TianGan{Id: 3, Name: "丙", WuXing: wuxing.HUO}
	DING = TianGan{Id: 4, Name: "丁", WuXing: wuxing.HUO}
	WU   = TianGan{Id: 5, Name: "戊", WuXing: wuxing.TU}
	JI   = TianGan{Id: 6, Name: "己", WuXing: wuxing.TU}
	GENG = TianGan{Id: 7, Name: "庚", WuXing: wuxing.JIN}
	XIN  = TianGan{Id: 8, Name: "辛", WuXing: wuxing.JIN}
	REN  = TianGan{Id: 9, Name: "壬", WuXing: wuxing.SHUI}
	GUI  = TianGan{Id: 10, Name: "癸", WuXing: wuxing.SHUI}
)

var (
	ALL = []TianGan{
		JIA,
		YI,
		BING,
		DING,
		WU,
		JI,
		GENG,
		XIN,
		REN,
		GUI,
	}
	SORT = []TianGan{
		WU,
		JI,
		GENG,
		XIN,
		REN,
		GUI,
		DING,
		BING,
		YI,
	}
	data = map[string]TianGan{
		"甲": JIA,
		"乙": YI,
		"丙": BING,
		"丁": DING,
		"戊": WU,
		"己": JI,
		"庚": GENG,
		"辛": XIN,
		"壬": REN,
		"癸": GUI,
	}
)

func Match(g string) TianGan {
	data, ok := data[g]
	if !ok {
		panic("天干匹配失败")
	}
	return data
}
