package zhangsheng

import (
	"qimenpan/pkg/dizhi"
	"qimenpan/pkg/tiangan"
)

var (
	ZhangShengOrder = []string{"生", "沐", "冠", "官", "旺", "衰", "病", "死", "墓", "绝", "胎", "养"}
	// 天干对应的起始长生地支
	ZhangShengStart = map[tiangan.TianGan]dizhi.DiZhi{
		tiangan.JIA: dizhi.HAI, tiangan.BING: dizhi.YIN, tiangan.WU: dizhi.YIN, tiangan.GENG: dizhi.SI, tiangan.REN: dizhi.SHEN, // 阳干
		tiangan.YI: dizhi.WU, tiangan.DING: dizhi.YOU, tiangan.JI: dizhi.YOU, tiangan.XIN: dizhi.ZI, tiangan.GUI: dizhi.MAO, // 阴干
	}
	ZhangShengMapping = map[int]map[int]string{}
)

func init() {
	for _, stem := range tiangan.ALL {
		startZhi := ZhangShengStart[stem]

		// 找到起始地支的索引
		startIndex := 0
		for i, zhi := range dizhi.ALL {
			if zhi == startZhi {
				startIndex = i
				break
			}
		}
		// 根据阴阳干确定顺序（阳干顺行，阴干逆行）
		isYang := stem == tiangan.JIA || stem == tiangan.BING || stem == tiangan.WU ||
			stem == tiangan.GENG || stem == tiangan.REN

		dizhiMapping := map[int]string{}
		for i := 0; i < 12; i++ {
			var idx int
			if isYang {
				idx = (startIndex + i) % 12
			} else {
				idx = (startIndex - i + 12) % 12
			}
			dizhiMapping[dizhi.ALL[idx].Id] = ZhangShengOrder[i]
		}
		ZhangShengMapping[stem.Id] = dizhiMapping
	}
}
