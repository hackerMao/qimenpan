package xunshou

import (
	"fmt"
	"qimenpan/pkg/dizhi"
	"qimenpan/pkg/tiangan"
)

type Xun struct {
	Id    int8
	Name  string
	DiZhi dizhi.DiZhi
	tiangan.TianGan
}

func (x *Xun) Print() {
	fmt.Printf("旬首：%s\n", x.Name)
}

var (
	JiaZiWu     = Xun{Id: 5, Name: "甲子旬", TianGan: tiangan.WU, DiZhi: dizhi.ZI}
	JiaXuJi     = Xun{Id: 6, Name: "甲戌旬", TianGan: tiangan.JI, DiZhi: dizhi.XU}
	JiaShenGeng = Xun{Id: 7, Name: "甲申旬", TianGan: tiangan.GENG, DiZhi: dizhi.SHEN}
	JiaWuXin    = Xun{Id: 8, Name: "甲午旬", TianGan: tiangan.XIN, DiZhi: dizhi.WU}
	JiaChenRen  = Xun{Id: 9, Name: "甲辰旬", TianGan: tiangan.REN, DiZhi: dizhi.CHEN}
	JiaYinGui   = Xun{Id: 10, Name: "甲寅旬", TianGan: tiangan.GUI, DiZhi: dizhi.YIN}
)

var XunMap = map[dizhi.DiZhi]Xun{
	dizhi.ZI:   JiaZiWu,
	dizhi.XU:   JiaXuJi,
	dizhi.SHEN: JiaShenGeng,
	dizhi.WU:   JiaWuXin,
	dizhi.CHEN: JiaChenRen,
	dizhi.YIN:  JiaYinGui,
}

func Match(z dizhi.DiZhi) Xun {
	return XunMap[z]
}
