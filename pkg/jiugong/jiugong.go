package jiugong

import (
	"qimenpan/pkg/wuxing"
)

type Gong struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	wuxing.WuXing `json:"wuxing.WuXing"`
}

var (
	KAN = Gong{
		Id:     1,
		Name:   "坎",
		WuXing: wuxing.SHUI,
	}
	KUN = Gong{
		Id:     2,
		Name:   "坤",
		WuXing: wuxing.TU,
	}
	ZHEN = Gong{
		Id:     3,
		Name:   "震",
		WuXing: wuxing.MU,
	}
	XUN = Gong{
		Id:     4,
		Name:   "巽",
		WuXing: wuxing.MU,
	}
	ZHONG = Gong{
		Id:     5,
		Name:   "中",
		WuXing: wuxing.TU,
	}
	QIAN = Gong{
		Id:     6,
		Name:   "乾",
		WuXing: wuxing.JIN,
	}
	DUI = Gong{
		Id:     7,
		Name:   "兑",
		WuXing: wuxing.JIN,
	}
	GEN = Gong{
		Id:     8,
		Name:   "艮",
		WuXing: wuxing.TU,
	}
	LI = Gong{
		Id:     9,
		Name:   "离",
		WuXing: wuxing.HUO,
	}

	UP   = []Gong{KAN, KUN, ZHEN, XUN, ZHONG, QIAN, DUI, GEN, LI}
	DOWN = []Gong{LI, GEN, DUI, QIAN, ZHONG, XUN, ZHEN, KUN, KAN}
)
