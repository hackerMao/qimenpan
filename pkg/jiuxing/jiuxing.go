package jiuxing

import (
	"qimenpan/pkg/dizhi"
	"qimenpan/pkg/jiugong"
	"qimenpan/pkg/wuxing"
)

// Star 天盘九星（北斗九星）
type Star struct {
	Id            int8                   `json:"id"`
	Name          string                 `json:"name"`  // 星名
	Color         string                 `json:"color"` // 颜色
	jiugong.Gong  `json:"jiugong.Gong"`  // 原始八卦宫位
	wuxing.WuXing `json:"wuxing.WuXing"` // 五行
}

var (
	TIAN_PEN   = Star{Id: 1, Name: "天蓬", Color: "白", Gong: jiugong.KAN, WuXing: wuxing.SHUI} // 天蓬/贪狼星
	TIAN_RUI   = Star{Id: 2, Name: "禽芮", Color: "黑", Gong: jiugong.KUN, WuXing: wuxing.TU}   // 天芮/巨门星
	TIAN_CHONG = Star{Id: 3, Name: "天冲", Color: "碧", Gong: jiugong.ZHEN, WuXing: wuxing.MU}  // 天冲/禄存星
	TIAN_FU    = Star{Id: 4, Name: "天辅", Color: "绿", Gong: jiugong.XUN, WuXing: wuxing.MU}   // 天辅/文曲星
	TIAN_QIN   = Star{Id: 5, Name: "天禽", Color: "黄", Gong: jiugong.ZHONG, WuXing: wuxing.TU} // 天禽/廉贞星
	TIAN_XIN   = Star{Id: 6, Name: "天心", Color: "白", Gong: jiugong.QIAN, WuXing: wuxing.JIN} // 天心/武曲星
	TIAN_ZHU   = Star{Id: 7, Name: "天柱", Color: "赤", Gong: jiugong.DUI, WuXing: wuxing.JIN}  // 天柱/破军星
	TIAN_REN   = Star{Id: 8, Name: "天任", Color: "白", Gong: jiugong.GEN, WuXing: wuxing.TU}   // 天任/左辅星
	TIAN_YING  = Star{Id: 9, Name: "天英", Color: "紫", Gong: jiugong.LI, WuXing: wuxing.HUO}   // 天英/右弼星
)

var YueStatus = map[int8]map[dizhi.DiZhi]string{
	1: {
		dizhi.YIN:  "旺",
		dizhi.MAO:  "旺",
		dizhi.CHEN: "囚",
		dizhi.SI:   "休",
		dizhi.WU:   "休",
		dizhi.WEI:  "囚",
		dizhi.SHEN: "废",
		dizhi.YOU:  "废",
		dizhi.XU:   "囚",
		dizhi.HAI:  "相",
		dizhi.ZI:   "相",
		dizhi.CHOU: "囚",
	},
	2: {
		dizhi.YIN:  "囚",
		dizhi.MAO:  "囚",
		dizhi.CHEN: "相",
		dizhi.SI:   "废",
		dizhi.WU:   "废",
		dizhi.WEI:  "相",
		dizhi.SHEN: "旺",
		dizhi.YOU:  "旺",
		dizhi.XU:   "相",
		dizhi.HAI:  "休",
		dizhi.ZI:   "休",
		dizhi.CHOU: "相",
	},
	3: {
		dizhi.YIN:  "相",
		dizhi.MAO:  "相",
		dizhi.CHEN: "休",
		dizhi.SI:   "旺",
		dizhi.WU:   "旺",
		dizhi.WEI:  "休",
		dizhi.SHEN: "囚",
		dizhi.YOU:  "囚",
		dizhi.XU:   "休",
		dizhi.HAI:  "废",
		dizhi.ZI:   "废",
		dizhi.CHOU: "休",
	},
	4: {
		dizhi.YIN:  "相",
		dizhi.MAO:  "相",
		dizhi.CHEN: "休",
		dizhi.SI:   "旺",
		dizhi.WU:   "旺",
		dizhi.WEI:  "休",
		dizhi.SHEN: "囚",
		dizhi.YOU:  "囚",
		dizhi.XU:   "休",
		dizhi.HAI:  "废",
		dizhi.ZI:   "废",
		dizhi.CHOU: "休",
	},
	6: {
		dizhi.YIN:  "休",
		dizhi.MAO:  "休",
		dizhi.CHEN: "废",
		dizhi.SI:   "囚",
		dizhi.WU:   "囚",
		dizhi.WEI:  "废",
		dizhi.SHEN: "相",
		dizhi.YOU:  "相",
		dizhi.XU:   "废",
		dizhi.HAI:  "旺",
		dizhi.ZI:   "旺",
		dizhi.CHOU: "废",
	},
	7: {
		dizhi.YIN:  "休",
		dizhi.MAO:  "休",
		dizhi.CHEN: "废",
		dizhi.SI:   "囚",
		dizhi.WU:   "囚",
		dizhi.WEI:  "废",
		dizhi.SHEN: "相",
		dizhi.YOU:  "相",
		dizhi.XU:   "废",
		dizhi.HAI:  "旺",
		dizhi.ZI:   "旺",
		dizhi.CHOU: "废",
	},
	8: {
		dizhi.YIN:  "囚",
		dizhi.MAO:  "囚",
		dizhi.CHEN: "相",
		dizhi.SI:   "废",
		dizhi.WU:   "废",
		dizhi.WEI:  "相",
		dizhi.SHEN: "旺",
		dizhi.YOU:  "旺",
		dizhi.XU:   "相",
		dizhi.HAI:  "休",
		dizhi.ZI:   "休",
		dizhi.CHOU: "相",
	},
	9: {
		dizhi.YIN:  "废",
		dizhi.MAO:  "废",
		dizhi.CHEN: "旺",
		dizhi.SI:   "相",
		dizhi.WU:   "相",
		dizhi.WEI:  "旺",
		dizhi.SHEN: "休",
		dizhi.YOU:  "休",
		dizhi.XU:   "旺",
		dizhi.HAI:  "囚",
		dizhi.ZI:   "囚",
		dizhi.CHOU: "旺",
	},
}

func (s *Star) Index() int {
	return int(s.Id - 1)
}

func (s *Star) YueStatus(yue dizhi.DiZhi) string {
	return YueStatus[s.Id][yue]
}
