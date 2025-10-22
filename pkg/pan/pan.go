package pan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"qimenpan/pkg/bamen"
	"qimenpan/pkg/bashen"
	"qimenpan/pkg/dizhi"
	"qimenpan/pkg/jiugong"
	"qimenpan/pkg/jiuxing"
	"qimenpan/pkg/sizhu"
	"qimenpan/pkg/tiangan"
	"qimenpan/pkg/wuxing"
	"qimenpan/pkg/xunshou"
	"qimenpan/pkg/zhangsheng"
	"strconv"
	"sync"
	"time"
)

// 宫位类型
type PalaceType int

const (
	PalaceKan PalaceType = iota + 1
	PalaceKun
	PalaceZhen
	PalaceXun
	PalaceZhong
	PalaceQian
	PalaceDui
	PalaceGen
	PalaceLi
)

// 宫位状态位
type GongWeiStatus uint32

const (
	StatusMaxing GongWeiStatus = 1 << iota
	StatusKongWang
	StatusMenPo
	StatusTGXing
	StatusDGXing
	StatusJTGXing
	StatusJDGXing
	StatusTGMu
	StatusDGMu
	StatusJTGMu
	StatusJDGMu
)

var sanHeMaMap = map[wuxing.WuXing]int{
	wuxing.SHUI: jiugong.GEN.Id,
	wuxing.JIN:  jiugong.QIAN.Id,
	wuxing.HUO:  jiugong.KUN.Id,
	wuxing.MU:   jiugong.XUN.Id,
}

// 击刑入墓规则
type XingMuRule struct {
	GongId   int
	XingGans []int
	MuGans   []int
}

// 缓存管理器
type CacheManager struct {
	palaceCache     sync.Map
	xunshouCache    sync.Map
	wuxingCache     sync.Map
	zhangshengCache sync.Map
	lunarCache      sync.Map
}

var cacheMgr = &CacheManager{}

// 预定义配置
var palaceConfigs = map[PalaceType]struct {
	OriginalStar jiuxing.Star
	OriginalMen  bamen.Men
	OriginalShen bashen.Shen
	DiZHiList    []dizhi.DiZhi
}{
	PalaceKan:   {jiuxing.TIAN_PEN, bamen.XIU, bashen.ZHI_FU, []dizhi.DiZhi{dizhi.ZI}},
	PalaceKun:   {jiuxing.TIAN_RUI, bamen.SI, bashen.XUAN_WU, []dizhi.DiZhi{dizhi.WEI, dizhi.SHEN}},
	PalaceZhen:  {jiuxing.TIAN_CHONG, bamen.SHANG, bashen.TAI_YIN, []dizhi.DiZhi{dizhi.MAO}},
	PalaceXun:   {jiuxing.TIAN_FU, bamen.DU, bashen.LIU_HE, []dizhi.DiZhi{dizhi.CHEN, dizhi.SI}},
	PalaceZhong: {jiuxing.TIAN_QIN, bamen.Men{}, bashen.Shen{}, []dizhi.DiZhi{}},
	PalaceQian:  {jiuxing.TIAN_XIN, bamen.KAI, bashen.JIU_TIAN, []dizhi.DiZhi{dizhi.XU, dizhi.HAI}},
	PalaceDui:   {jiuxing.TIAN_ZHU, bamen.JING2, bashen.JIU_DI, []dizhi.DiZhi{dizhi.YOU}},
	PalaceGen:   {jiuxing.TIAN_REN, bamen.SHENG, bashen.TEN_SHE, []dizhi.DiZhi{dizhi.CHOU, dizhi.YIN}},
	PalaceLi:    {jiuxing.TIAN_YING, bamen.JING, bashen.BAIHU, []dizhi.DiZhi{dizhi.WU}},
}

var xunshouEmptyMap = map[xunshou.Xun][]int{
	xunshou.JiaZiWu:     {6},
	xunshou.JiaYinGui:   {8, 1},
	xunshou.JiaChenRen:  {3, 8},
	xunshou.JiaWuXin:    {4},
	xunshou.JiaShenGeng: {2, 9},
	xunshou.JiaXuJi:     {7, 2},
}

var xingMuRules = []XingMuRule{
	{jiugong.GEN.Id, []int{tiangan.GENG.Id}, []int{tiangan.GENG.Id, tiangan.JI.Id, tiangan.DING.Id}},
	{jiugong.ZHEN.Id, []int{tiangan.WU.Id}, []int{}},
	{jiugong.XUN.Id, []int{tiangan.REN.Id, tiangan.GUI.Id}, []int{tiangan.XIN.Id, tiangan.REN.Id}},
	{jiugong.LI.Id, []int{tiangan.XIN.Id}, []int{}},
	{jiugong.KUN.Id, []int{tiangan.JI.Id}, []int{tiangan.GUI.Id}},
	{jiugong.QIAN.Id, []int{}, []int{tiangan.YI.Id, tiangan.BING.Id, tiangan.WU.Id}},
}

type GongWei struct {
	jiugong.Gong
	IsNeiPan         bool
	OriginalStar     jiuxing.Star
	OriginalMen      bamen.Men
	OriginalShen     bashen.Shen
	DiZHiList        []dizhi.DiZhi
	DiPanGan         tiangan.TianGan
	TianPanGan       tiangan.TianGan
	YinGan           []tiangan.TianGan
	JiDiGan          tiangan.TianGan
	JiTianGan        tiangan.TianGan
	TianPan          jiuxing.Star
	ShenPan          bashen.Shen
	RenPan           bamen.Men
	Status           GongWeiStatus
	GongStatus       string
	JiuXingStatus    string
	JiuXingYueStatus string
	BamenStatus      string
	BamenYueStatus   string
	TianPanGanStatus string
	DiPanGanStatus   string
	JiDiGanStatus    string
	JiTianGanStatus  string
}

// 状态操作方法
func (g *GongWei) SetStatus(status GongWeiStatus) {
	g.Status |= status
}

func (g *GongWei) ClearStatus(status GongWeiStatus) {
	g.Status &^= status
}

func (g *GongWei) HasStatus(status GongWeiStatus) bool {
	return g.Status&status != 0
}

// 初始化宫位
func initGongWeiS() []GongWei {
	gws := make([]GongWei, 9)

	types := []PalaceType{PalaceKan, PalaceKun, PalaceZhen, PalaceXun, PalaceZhong, PalaceQian, PalaceDui, PalaceGen, PalaceLi}
	gongRefs := []*jiugong.Gong{&jiugong.KAN, &jiugong.KUN, &jiugong.ZHEN, &jiugong.XUN, &jiugong.ZHONG, &jiugong.QIAN, &jiugong.DUI, &jiugong.GEN, &jiugong.LI}

	for i, palaceType := range types {
		config := palaceConfigs[palaceType]
		gws[i] = GongWei{
			Gong:         *gongRefs[i],
			OriginalStar: config.OriginalStar,
			OriginalMen:  config.OriginalMen,
			OriginalShen: config.OriginalShen,
			DiZHiList:    config.DiZHiList,
		}
	}

	return gws
}

var GWS = initGongWeiS()
var BaGong = []GongWei{GWS[0], GWS[7], GWS[2], GWS[3], GWS[8], GWS[1], GWS[6], GWS[5]}

type QTime struct {
	GongLi   time.Time
	YinLi    time.Time
	YinLiStr string
	JieQi    string
	sizhu.SiZhu
}

func (q *QTime) Print() {
	fmt.Printf("公历：%d年%d月%d日%d:%d\n",
		q.GongLi.Year(), q.GongLi.Month(), q.GongLi.Day(), q.GongLi.Hour(), q.GongLi.Minute())
	fmt.Printf("陰历：%d年%d月%d日\n",
		q.YinLi.Year(), q.YinLi.Month(), q.YinLi.Day())
	q.SiZhu.Print()
}

// 排盘策略接口
type ArrangeStrategy interface {
	Arrange(pan *Pan) error
	Name() string
}

// 具体策略实现
type DiPanGanStrategy struct{}

func (s *DiPanGanStrategy) Arrange(p *Pan) error {
	index := p.JuShu - 1
	gws := make([]*GongWei, 9)

	if !p.YinYangDun {
		for i := 0; i < p.JuShu; i++ {
			gws[i] = p.GongWeiS[index-i]
		}
		sub := 0
		for i := index + 1; i < 9; i++ {
			gws[i] = p.GongWeiS[8-sub]
			sub++
		}
	} else {
		gws = p.GongWeiS[index:]
		if p.JuShu != 1 {
			gws = append(gws, p.GongWeiS[:index]...)
		}
	}

	gongMap := make(map[int]*GongWei, 9)
	for i, v := range gws {
		v.DiPanGan = tiangan.SORT[i]
		gongMap[v.Id] = v
	}

	if p.diGanMap == nil {
		p.diGanMap = make(map[tiangan.TianGan]*GongWei, 9)
	}
	for _, v := range gws {
		p.diGanMap[v.DiPanGan] = v
	}

	// 中五宫寄坤二宫
	zg := gongMap[5]
	gongMap[2].JiDiGan = zg.DiPanGan

	return nil
}

func (s *DiPanGanStrategy) Name() string {
	return "DiPanGanStrategy"
}

type ZhiFuShiStrategy struct{}

func (s *ZhiFuShiStrategy) Arrange(p *Pan) error {
	dunJia := p.XunShou.TianGan
	gw := p.diGanMap[dunJia]
	p.ZhiFu = gw.OriginalStar
	if gw.Gong.Id == 5 {
		gw = p.gongMap[2]
	}
	p.ZhiShi = gw.OriginalMen
	return nil
}

func (s *ZhiFuShiStrategy) Name() string {
	return "ZhiFuShiStrategy"
}

type TianPanStrategy struct{}

func (s *TianPanStrategy) Arrange(p *Pan) error {
	gongMap := p.gongMap
	gws := make([]*GongWei, 8)
	sort := []int{1, 8, 3, 4, 9, 2, 7, 6}
	for i, index := range sort {
		gws[i] = gongMap[index]
	}

	xunShouGong := p.diGanMap[p.XunShou.TianGan]
	gongIndex := -1
	for i, v := range gws {
		if v.Id == xunShouGong.Gong.Id {
			gongIndex = i
		}
	}
	if gongIndex == -1 {
		gongIndex = 5
	}

	gwsCycle := make([]*GongWei, 8)
	marshal, _ := json.Marshal(gws)
	_ = json.Unmarshal(marshal, &gwsCycle)

	if gongIndex > 0 {
		for i := 0; i < 8-gongIndex; i++ {
			gwsCycle[i] = gws[i+gongIndex]
		}
		for i := 8 - gongIndex; i < 8; i++ {
			gwsCycle[i] = gws[i-8+gongIndex]
		}
	}

	hourGanIndex := -1
	hourGanId := p.HourGan.Id
	if hourGanId == 1 {
		hourGanId = p.XunShou.TianGan.Id
	}

	for i, v := range gws {
		if v.DiPanGan.Id == hourGanId {
			hourGanIndex = i
		}
	}
	if hourGanIndex == -1 {
		hourGanIndex = 5
	}

	cycle2 := make([]*GongWei, 8)
	copy(cycle2, gws)
	if hourGanIndex > 0 {
		for i := 0; i < 8-hourGanIndex; i++ {
			cycle2[i] = gws[i+hourGanIndex]
		}
		for i := 8 - hourGanIndex; i < 8; i++ {
			cycle2[i] = gws[i-8+hourGanIndex]
		}
	}

	for i, v := range cycle2 {
		v.TianPan = gwsCycle[i].OriginalStar
		v.TianPanGan = gwsCycle[i].DiPanGan
		v.JiTianGan = gwsCycle[i].JiDiGan
	}

	gongWei := p.GongWeiS[0]
	if gongWei.TianPanGan.Id == gongWei.DiPanGan.Id {
		p.FuYin = true
	}

	return nil
}

func (s *TianPanStrategy) Name() string {
	return "TianPanStrategy"
}

type ShenPanStrategy struct{}

func (s *ShenPanStrategy) Arrange(p *Pan) error {
	gws := make([]*GongWei, 8)
	sort := []int{1, 8, 3, 4, 9, 2, 7, 6}
	for i, index := range sort {
		gws[i] = p.gongMap[index]
	}

	shiGan := p.HourGan
	if shiGan.Id == 1 {
		shiGan = p.XunShou.TianGan
	}
	houGong := p.diGanMap[shiGan]

	gongIndex := -1
	for i, v := range gws {
		if v.Id == houGong.Gong.Id {
			gongIndex = i
		}
	}
	if gongIndex == -1 {
		gongIndex = 5
	}

	gwsCycle := make([]*GongWei, 8)
	marshal, _ := json.Marshal(gws)
	_ = json.Unmarshal(marshal, &gwsCycle)

	if gongIndex > 0 {
		if p.YinYangDun {
			for i := 0; i < 8-gongIndex; i++ {
				gwsCycle[i] = gws[i+gongIndex]
			}
			for i := 8 - gongIndex; i < 8; i++ {
				gwsCycle[i] = gws[i-8+gongIndex]
			}
		} else {
			for i := 0; i < gongIndex+1; i++ {
				gwsCycle[i] = gws[gongIndex-i]
			}

			sub := 0
			for i := gongIndex + 1; i < 8; i++ {
				gwsCycle[i] = gws[7-sub]
				sub++
			}
		}
	}

	shenMap := make(map[int]bashen.Shen, 8)
	for i, v := range gwsCycle {
		shenMap[v.Id] = bashen.ALL[i]
	}

	for _, gw := range p.GongWeiS {
		gw.ShenPan = shenMap[gw.Id]
	}

	return nil
}

func (s *ShenPanStrategy) Name() string {
	return "ShenPanStrategy"
}

type RenPanStrategy struct{}

func (s *RenPanStrategy) Arrange(p *Pan) error {
	touGan := p.XunShou.TianGan
	touGanGong := p.diGanMap[touGan]
	step := touGanGong.Id
	if step == 5 {
		step = 2
	}

	tgsMap := make(map[tiangan.TianGan]int, 10)
	if p.YinYangDun {
		for _, gan := range tiangan.ALL {
			tgsMap[gan] = step
			if step == 9 {
				step = 0
			}
			step++
		}
	} else {
		for _, gan := range tiangan.ALL {
			tgsMap[gan] = step
			step--
			if step == 0 {
				step = 9
			}
		}
	}

	gongWei := tgsMap[p.HourGan]
	gws := make([]*GongWei, 8)
	sort := []int{1, 8, 3, 4, 9, 2, 7, 6}
	for i, index := range sort {
		gws[i] = p.gongMap[index]
	}

	gongIndex := 0
	for i, v := range gws {
		if v.Id == gongWei {
			gongIndex = i
		}
	}

	gwsCycle := make([]*GongWei, 8)
	gwsCycle = gws[gongIndex:]
	gwsCycle = append(gwsCycle, gws[:gongIndex]...)

	menCycle := make([]bamen.Men, 8)
	menIndex := 0
	for i, v := range gws {
		if v.Gong.Id == p.diGanMap[p.XunShou.TianGan].Id {
			menIndex = i
		}
	}
	menCycle = bamen.ALL[menIndex:]
	menCycle = append(menCycle, bamen.ALL[:menIndex]...)

	for i, gw := range gwsCycle {
		gw.RenPan = menCycle[i]
	}

	menMap := make(map[int]bamen.Men, 8)
	for i, v := range gwsCycle {
		menMap[v.Id] = menCycle[i]
	}

	for _, gw := range p.GongWeiS {
		gw.RenPan = menMap[gw.Id]
	}

	return nil
}

func (s *RenPanStrategy) Name() string {
	return "RenPanStrategy"
}

type YinGanStrategy struct{}

func (s *YinGanStrategy) Arrange(p *Pan) error {
	if p.FuYin {
		gongSort := make([]jiugong.Gong, 9)
		gongMap := make(map[int][]tiangan.TianGan, 9)
		if p.YinYangDun {
			gongSort = append(gongSort, jiugong.ZHONG, jiugong.QIAN, jiugong.DUI, jiugong.GEN, jiugong.LI)
			gongSort = append(gongSort, jiugong.KAN, jiugong.KUN, jiugong.ZHEN, jiugong.XUN)
		} else {
			gongSort = append(gongSort, jiugong.ZHONG, jiugong.XUN, jiugong.ZHEN, jiugong.KUN, jiugong.KAN)
			gongSort = append(gongSort, jiugong.LI, jiugong.GEN, jiugong.DUI, jiugong.QIAN)
		}

		startIndex := 0
		hourGanId := p.HourGan.Id
		if hourGanId == 1 {
			hourGanId = p.XunShou.TianGan.Id
		}
		for i, v := range tiangan.SORT {
			if v.Id == hourGanId {
				startIndex = i
			}
		}

		for i, v := range gongSort {
			gongMap[v.Id] = []tiangan.TianGan{tiangan.SORT[(startIndex+i)%9]}
		}

		// 中5宫寄坤2宫
		yinGan := gongMap[5]
		gans := gongMap[2]
		gans = append(gans, yinGan...)
		gongMap[2] = gans

		for _, v := range p.GongWeiS {
			v.YinGan = gongMap[v.Id]
		}
	} else {
		gws := make([]*GongWei, 8)
		sort := []int{1, 8, 3, 4, 9, 2, 7, 6}
		i := 0
		for _, index := range sort {
			gws[i] = p.gongMap[index]
			i++
		}

		// 找到值使宫位置
		zhiShiIndex := 0
		tgs := make([][]tiangan.TianGan, 8)
		for i, v := range gws {
			if v.RenPan.Name == p.ZhiShi.Name {
				zhiShiIndex = i
			}
		}

		cycle := make([]*GongWei, 8)
		cycle = gws[zhiShiIndex:]
		cycle = append(cycle, gws[:zhiShiIndex]...)

		for i, index := range sort {
			g1 := p.gongMap[index].DiPanGan
			tgs[i] = []tiangan.TianGan{g1}
			g2 := p.gongMap[index].JiDiGan
			if g2.Id > 0 {
				tgs[i] = append(tgs[i], g2)
			}
		}

		tgIndex := 0
		hourGanId := p.HourGan.Id
		if hourGanId == 1 {
			hourGanId = p.XunShou.TianGan.Id
		}
		for i, v := range tgs {
			if len(v) == 1 {
				if v[0].Id == hourGanId {
					tgIndex = i
				}
			} else {
				if v[0].Id == hourGanId || v[1].Id == hourGanId {
					tgIndex = i
				}
			}
		}

		tgCycle := make([][]tiangan.TianGan, 8)
		tgCycle = tgs[tgIndex:]
		tgCycle = append(tgCycle, tgs[:tgIndex]...)

		for i, gw := range cycle {
			gw.YinGan = tgCycle[i]
		}

		yinGanMap := make(map[int][]tiangan.TianGan, 8)
		for _, v := range cycle {
			yinGanMap[v.Id] = v.YinGan
		}

		for _, gw := range p.GongWeiS {
			gw.YinGan = yinGanMap[gw.Id]
		}
	}

	return nil
}

func (s *YinGanStrategy) Name() string {
	return "YinGanStrategy"
}

type EmptyMaXingStrategy struct{}

func (s *EmptyMaXingStrategy) Arrange(p *Pan) error {
	// 排空亡
	if emptyGong, exists := xunshouEmptyMap[p.XunShou]; exists {
		for _, v := range p.GongWeiS {
			for _, emptyId := range emptyGong {
				if v.Id == emptyId {
					v.SetStatus(StatusKongWang)
					break
				}
			}
		}
	}

	// 排马星
	sanHe := dizhi.SanHe(p.HourZhi)
	if ma, exists := sanHeMaMap[sanHe.WuXing]; exists {
		for _, v := range p.GongWeiS {
			if v.Id == ma {
				v.SetStatus(StatusMaxing)
			}
		}
	}

	return nil
}

func (s *EmptyMaXingStrategy) Name() string {
	return "EmptyMaXingStrategy"
}

type SiHaiStrategy struct{}

func (s *SiHaiStrategy) Arrange(p *Pan) error {
	p.batchUpdatePalaces(func(gw *GongWei) {
		gw.ClearStatus(StatusMenPo | StatusTGXing | StatusDGXing | StatusJTGXing | StatusJDGXing |
			StatusTGMu | StatusDGMu | StatusJTGMu | StatusJDGMu)

		// 门迫判断
		if gw.RenPan.Name != "" && wuxing.IsKe(gw.RenPan.Gong.WuXing, gw.Gong.WuXing) {
			gw.SetStatus(StatusMenPo)
		}

		// 应用击刑和入墓规则
		p.applyXingMuRules(gw)
	})

	return nil
}

func (s *SiHaiStrategy) Name() string {
	return "SiHaiStrategy"
}

type StatusStrategy struct{}

func (s *StatusStrategy) Arrange(p *Pan) error {
	p.arrangeStatus()
	return nil
}

func (s *StatusStrategy) Name() string {
	return "StatusStrategy"
}

type InnerOutStrategy struct{}

func (s *InnerOutStrategy) Arrange(p *Pan) error {
	p.arrangeInnerOut()
	return nil
}

func (s *InnerOutStrategy) Name() string {
	return "InnerOutStrategy"
}

// 策略管理器
type ArrangeManager struct {
	strategies []ArrangeStrategy
}

func NewArrangeManager() *ArrangeManager {
	return &ArrangeManager{
		strategies: []ArrangeStrategy{
			&DiPanGanStrategy{},
			&ZhiFuShiStrategy{},
			&TianPanStrategy{},
			&ShenPanStrategy{},
			&RenPanStrategy{},
			&YinGanStrategy{},
			&EmptyMaXingStrategy{},
			&SiHaiStrategy{},
			&StatusStrategy{},
			&InnerOutStrategy{},
		},
	}
}

func (am *ArrangeManager) Execute(pan *Pan) error {
	for _, strategy := range am.strategies {
		if err := strategy.Arrange(pan); err != nil {
			return fmt.Errorf("strategy %s failed: %w", strategy.Name(), err)
		}
	}
	return nil
}

type Pan struct {
	QTime
	YinYangDun bool
	JuShu      int
	XunShou    xunshou.Xun
	ZhiFu      jiuxing.Star
	ZhiShi     bamen.Men
	GongWeiS   []*GongWei
	diGanMap   map[tiangan.TianGan]*GongWei
	gongMap    map[int]*GongWei
	FuYin      bool
	arrangeMgr *ArrangeManager
}

// 批量处理宫位
func (p *Pan) batchUpdatePalaces(updater func(*GongWei)) {
	for i := range p.GongWeiS {
		updater(p.GongWeiS[i])
	}
}

// 应用击刑入墓规则
func (p *Pan) applyXingMuRules(gw *GongWei) {
	for _, rule := range xingMuRules {
		if gw.Gong.Id == rule.GongId {
			// 击刑判断
			for _, ganId := range rule.XingGans {
				if gw.TianPanGan.Id == ganId {
					gw.SetStatus(StatusTGXing)
				}
				if gw.DiPanGan.Id == ganId {
					gw.SetStatus(StatusDGXing)
				}
				if gw.JiTianGan.Id == ganId {
					gw.SetStatus(StatusJTGXing)
				}
				if gw.JiDiGan.Id == ganId {
					gw.SetStatus(StatusJDGXing)
				}
			}

			// 入墓判断
			for _, ganId := range rule.MuGans {
				if gw.TianPanGan.Id == ganId {
					gw.SetStatus(StatusTGMu)
				}
				if gw.DiPanGan.Id == ganId {
					gw.SetStatus(StatusDGMu)
				}
				if gw.JiTianGan.Id == ganId {
					gw.SetStatus(StatusJTGMu)
				}
				if gw.JiDiGan.Id == ganId {
					gw.SetStatus(StatusJDGMu)
				}
			}
		}
	}
}

// 农历服务
type CalendarService struct {
	client *http.Client
}

func NewCalendarService() *CalendarService {
	return &CalendarService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (cs *CalendarService) GetLunarInfo(year, month, day int) ([]byte, error) {
	cacheKey := fmt.Sprintf("%d-%d-%d", year, month, day)

	if cached, found := cacheMgr.lunarCache.Load(cacheKey); found {
		return cached.([]byte), nil
	}

	data, err := cs.fetchLunarInfo(year, month, day)
	if err != nil {
		return nil, err
	}

	cacheMgr.lunarCache.Store(cacheKey, data)
	return data, nil
}

func (cs *CalendarService) fetchLunarInfo(year, month, day int) ([]byte, error) {
	id := 10006676
	key := "4da1df18ac273b15942f3815205d0fe2"
	url := fmt.Sprintf("%s?id=%d&key=%s&nian=%d&yue=%d&ri=%d",
		"https://cn.apihz.cn/api/time/getzdday.php", id, key, year, month, day)

	resp, err := cs.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func timeTransfer(gongLi time.Time) (QTime, error) {
	if gongLi == (time.Time{}) {
		gongLi = time.Now()
	}

	service := NewCalendarService()
	body, err := service.GetLunarInfo(gongLi.Year(), int(gongLi.Month()), gongLi.Day())
	if err != nil {
		return QTime{}, err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		return QTime{}, err
	}

	jieQi := data["jieqi"].(string)
	nianZhu := data["ganzhinian"].(string)
	yueZhu := data["ganzhiyue"].(string)
	riZhu := data["ganzhiri"].(string)
	yinYue := data["nyue"].(string)
	yinRi := data["nri"].(string)

	if nianZhu == "" || yueZhu == "" || riZhu == "" {
		return QTime{}, errors.New("四柱转换异常")
	}

	if len([]rune(nianZhu)) != 2 || len([]rune(yueZhu)) != 2 || len([]rune(riZhu)) != 2 {
		return QTime{}, errors.New("四柱转换异常")
	}

	ng := string([]rune(nianZhu)[0:1])
	nz := string([]rune(nianZhu)[1:])
	yg := string([]rune(yueZhu)[0:1])
	yz := string([]rune(yueZhu)[1:])
	rg := string([]rune(riZhu)[0:1])
	rz := string([]rune(riZhu)[1:])

	yinMonth, err := strconv.Atoi(data["YIMONTH"].(string))
	if err != nil {
		return QTime{}, errors.New("四柱转换异常")
	}
	yinDay, err := strconv.Atoi(data["YIDAY"].(string))
	if err != nil {
		return QTime{}, errors.New("四柱转换异常")
	}

	yinLi := time.Date(gongLi.Year(), time.Month(yinMonth), yinDay, 0, 0, 0, 0, time.Local)

	siZhu, err := sizhu.New(ng, nz, yg, yz, rg, rz, gongLi.Hour())
	if err != nil {
		return QTime{}, err
	}

	return QTime{
		GongLi:   gongLi,
		YinLi:    yinLi,
		YinLiStr: yinYue + yinRi,
		JieQi:    jieQi,
		SiZhu:    siZhu,
	}, nil
}

func getYinYangDun(lunarMonth int) bool {
	if lunarMonth >= 11 || lunarMonth <= 5 {
		return true
	}
	return false
}

func calcJuShu(q QTime) int {
	js := (q.YearZhi.Id + int(q.YinLi.Month()) + q.YinLi.Day() + q.HourZhi.Id) % 9
	if js == 0 {
		js = 9
	}
	return js
}

func calcXunShou(q QTime) xunshou.Xun {
	shiGanIndex := q.HourGan.Id - 1
	index := q.HourZhi.Id - shiGanIndex
	z := dizhi.Index(index)
	return xunshou.Match(z)
}

func New(gongLi time.Time) Pan {
	qTime, err := timeTransfer(gongLi)
	if err != nil {
		return Pan{}
	}

	gongWeiS := make([]*GongWei, 9)
	marshal, _ := json.Marshal(GWS)
	_ = json.Unmarshal(marshal, &gongWeiS)

	p := Pan{
		QTime:      qTime,
		YinYangDun: getYinYangDun(int(qTime.YinLi.Month())),
		JuShu:      calcJuShu(qTime),
		XunShou:    calcXunShou(qTime),
		GongWeiS:   gongWeiS,
		arrangeMgr: NewArrangeManager(),
	}

	gongMap := map[int]*GongWei{}
	for _, v := range p.GongWeiS {
		gongMap[v.Id] = v
	}
	p.gongMap = gongMap

	// 使用策略模式执行所有排盘步骤
	if err := p.arrangeMgr.Execute(&p); err != nil {
		fmt.Printf("排盘错误: %v\n", err)
	}

	return p
}

// 旺衰状态相关方法
func (p *Pan) arrangeStatus() {
	for _, gw := range p.GongWeiS {
		p.setGongStatus(gw)
		p.setJiuXingStatus(gw)
		p.setBaMenStatus(gw)
		p.setTianGanStatus(gw)
	}
}

func (p *Pan) setJiuXingStatus(gw *GongWei) {
	status := ""
	if gw.TianPan.WuXing.Gen().Name == gw.WuXing.Name {
		status = "旺"
	} else if gw.WuXing.Tong(&gw.TianPan.WuXing) {
		status = "相"
	} else if gw.TianPan.WuXing.Ke().Name == gw.WuXing.Name {
		status = "休"
	} else if gw.WuXing.Gen().Name == gw.TianPan.WuXing.Name {
		status = "废"
	} else if gw.WuXing.Ke().Name == gw.TianPan.WuXing.Name {
		status = "囚"
	}
	gw.JiuXingStatus = status

	yueStatus := ""
	if gw.TianPan.WuXing.Gen().Name == p.MonthZhi.WuXing.Name {
		yueStatus = "月旺"
	} else if p.MonthZhi.Tong(&gw.TianPan.WuXing) {
		yueStatus = "月相"
	} else if gw.TianPan.WuXing.Ke().Name == p.MonthZhi.WuXing.Name {
		yueStatus = "月休"
	} else if p.MonthZhi.WuXing.Gen().Name == gw.TianPan.WuXing.Name {
		yueStatus = "月废"
	} else if p.MonthZhi.WuXing.Ke().Name == gw.TianPan.WuXing.Name {
		yueStatus = "月囚"
	}
	gw.JiuXingYueStatus = yueStatus
}

func (p *Pan) setBaMenStatus(gw *GongWei) {
	status := ""
	if gw.WuXing.Tong(&gw.RenPan.WuXing) {
		status = "旺"
	} else if gw.RenPan.WuXing.Name == gw.Gen().Name {
		status = "相"
	} else if gw.RenPan.Gen().Name == gw.WuXing.Name {
		status = "休"
	} else if gw.RenPan.Ke().Name == gw.WuXing.Name {
		status = "囚"
	} else if gw.Ke().Name == gw.RenPan.WuXing.Name {
		status = "死"
	}
	gw.BamenStatus = status

	yueStatus := ""
	if p.MonthZhi.WuXing.Tong(&gw.RenPan.WuXing) {
		yueStatus = "月旺"
	} else if gw.RenPan.WuXing.Name == p.MonthZhi.Gen().Name {
		yueStatus = "月相"
	} else if gw.RenPan.Gen().Name == p.MonthZhi.WuXing.Name {
		yueStatus = "月休"
	} else if gw.RenPan.Ke().Name == p.MonthZhi.WuXing.Name {
		yueStatus = "月囚"
	} else if p.MonthZhi.Ke().Name == gw.RenPan.WuXing.Name {
		yueStatus = "月死"
	}
	gw.BamenYueStatus = yueStatus
}

func (p *Pan) setGongStatus(gw *GongWei) {
	status := ""
	if p.MonthZhi.ShiLinWuXing.Name == gw.WuXing.Name {
		status = "旺"
	} else if p.MonthZhi.ShiLinWuXing.Gen().Name == gw.WuXing.Name {
		status = "相"
	} else if gw.WuXing.Gen().Name == p.MonthZhi.ShiLinWuXing.Name {
		status = "休"
	} else if gw.WuXing.Ke().Name == p.MonthZhi.ShiLinWuXing.Name {
		status = "囚"
	} else if p.MonthZhi.ShiLinWuXing.Ke().Name == gw.WuXing.Name {
		status = "死"
	}
	gw.GongStatus = status
}

func (p *Pan) setTianGanStatus(gw *GongWei) {
	mapping := zhangsheng.ZhangShengMapping[gw.TianPanGan.Id]
	s1 := ""
	s2 := ""
	for _, dz := range gw.DiZHiList {
		s1 += mapping[dz.Id]
	}
	gw.TianPanGanStatus = s1

	mapping = zhangsheng.ZhangShengMapping[gw.DiPanGan.Id]
	for _, dz := range gw.DiZHiList {
		s2 += mapping[dz.Id]
	}
	gw.DiPanGanStatus = s2

	if gw.JiTianGan.Id > 0 {
		mapping = zhangsheng.ZhangShengMapping[gw.JiTianGan.Id]
		s3 := ""
		for _, dz := range gw.DiZHiList {
			s3 += mapping[dz.Id]
		}
		gw.JiTianGanStatus = s3
	}
	if gw.JiDiGan.Id > 0 {
		mapping = zhangsheng.ZhangShengMapping[gw.JiDiGan.Id]
		s4 := ""
		for _, dz := range gw.DiZHiList {
			s4 += mapping[dz.Id]
		}
		gw.JiDiGanStatus = s4
	}
}

func (p *Pan) arrangeInnerOut() {
	if p.YinYangDun {
		for _, gw := range p.GongWeiS {
			gw.IsNeiPan = (gw.Id == 1 || gw.Id == 3 || gw.Id == 4 || gw.Id == 8)
		}
	} else {
		for _, gw := range p.GongWeiS {
			gw.IsNeiPan = (gw.Id == 9 || gw.Id == 2 || gw.Id == 7 || gw.Id == 6)
		}
	}
}

// 输出数据结构
type GwInfo struct {
	PalaceId         int    `json:"palaceId"`
	PalaceName       string `json:"palaceName"`
	Yingan           string `json:"yingan"`
	Bashen           string `json:"bashen"`
	Jiuxing          string `json:"jiuxing"`
	Tianpangan       string `json:"tianpangan"`
	Dipangan         string `json:"dipangan"`
	Jigan            string `json:"jigan"`
	JiTiangan        string `json:"jiTianGan"`
	Bamen            string `json:"bamen"`
	Maxing           bool   `json:"maxing"`
	Kongwang         bool   `json:"kongwang"`
	Menpo            bool   `json:"menpo"`
	DipanganJixing   bool   `json:"dipanganJixing"`
	TianpanganJixing bool   `json:"tianpanganJixing"`
	JiganJixing      bool   `json:"jiganJixing"`
	JiTianGanJixing  bool   `json:"jiTianGanJixing"`
	TianpanganRumu   bool   `json:"tianpanganRumu"`
	DipanganRumu     bool   `json:"dipanganRumu"`
	JiganRumu        bool   `json:"jiganRumu"`
	JiTianGanRumu    bool   `json:"jiTianGanRumu"`
	GongStatus       string `json:"gongStatus"`
	IsNeiPan         bool   `json:"isNeiPan"`
	JiuxingStatus    string `json:"jiuxingStatus"`
	JiuxingYueStatus string `json:"jiuxingYueStatus"`
	BamenStatus      string `json:"bamenStatus"`
	BamenYueStatus   string `json:"bamenYueStatus"`
	TianPanGanStatus string `json:"tianPanGanStatus"`
	DiPanGanStatus   string `json:"diPanGanStatus"`
	JiTianGanStatus  string `json:"JiTianGanStatus"`
	JiDiGanStatus    string `json:"jiDiGanStatus"`
}

type PanData struct {
	GregorianTime string   `json:"gregorianTime"`
	LunarTime     string   `json:"lunarTime"`
	GanzhiTime    string   `json:"ganzhiTime"`
	JuNumber      string   `json:"juNumber"`
	XunShou       string   `json:"xunShou"`
	ZhiFu         string   `json:"zhiFu"`
	ZhiShi        string   `json:"zhiShi"`
	Palaces       []GwInfo `json:"palaces"`
	PalaceInfo    string   `json:"palaceInfo"`
}

func (p *Pan) Parse() PanData {
	data := PanData{
		GregorianTime: p.GongLi.Format("2006年01月02日15:04"),
		LunarTime:     p.YinLiStr,
		GanzhiTime:    p.SiZhu.String(),
		JuNumber:      getJuShu(p.JuShu, p.YinYangDun),
		XunShou:       p.XunShou.Name,
		ZhiFu:         p.ZhiFu.Name,
		ZhiShi:        p.ZhiShi.Name,
		Palaces:       getGwInfo(p.GongWeiS),
	}
	data.PalaceInfo = generateAIPrompt(data.Palaces)
	return data
}

func generateAIPrompt(gws []GwInfo) string {
	tpl := "%s宫：\n    天盘干：%s;\n    地盘干：%s;\n    九星：%s;\n    八神：%s;\n    八门：%s;\n    引干：%s;\n    寄宫天干：%s;\n"
	tip := ""
	for _, gw := range gws {
		tip += fmt.Sprintf(tpl,
			gw.PalaceName,
			gw.Tianpangan,
			gw.Dipangan,
			gw.Jiuxing,
			gw.Bashen,
			gw.Bamen,
			gw.Yingan,
			gw.JiTiangan,
		)
		if gw.Kongwang {
			tip += fmt.Sprintf("落空亡;\n")
		}
		if gw.Maxing {
			tip += fmt.Sprintf("落马星;\n")
		}
	}
	return tip
}

func getGwInfo(s []*GongWei) []GwInfo {
	if len(s) == 0 {
		return []GwInfo{}
	}
	gws := make([]GwInfo, 0, len(s))
	for i := range s {
		gw := s[i]
		yg := ""
		if len(gw.YinGan) > 0 {
			yg = gw.YinGan[0].Name
			if len(gw.YinGan) > 1 {
				yg = yg + gw.YinGan[1].Name
			}
		}

		g := GwInfo{
			PalaceId:         gw.Id,
			PalaceName:       gw.Name,
			Yingan:           yg,
			Bashen:           gw.ShenPan.Name,
			Jiuxing:          gw.TianPan.Name,
			Tianpangan:       gw.TianPanGan.Name,
			Dipangan:         gw.DiPanGan.Name,
			Jigan:            gw.JiDiGan.Name,
			JiTiangan:        gw.JiTianGan.Name,
			Bamen:            gw.RenPan.Name,
			Maxing:           gw.HasStatus(StatusMaxing),
			Kongwang:         gw.HasStatus(StatusKongWang),
			Menpo:            gw.HasStatus(StatusMenPo),
			DipanganJixing:   gw.HasStatus(StatusDGXing),
			TianpanganJixing: gw.HasStatus(StatusTGXing),
			JiganJixing:      gw.HasStatus(StatusJDGXing),
			JiTianGanJixing:  gw.HasStatus(StatusJTGXing),
			TianpanganRumu:   gw.HasStatus(StatusTGMu),
			DipanganRumu:     gw.HasStatus(StatusDGMu),
			JiganRumu:        gw.HasStatus(StatusJDGMu),
			JiTianGanRumu:    gw.HasStatus(StatusJTGMu),
			JiuxingStatus:    gw.JiuXingStatus,
			JiuxingYueStatus: gw.JiuXingYueStatus,
			BamenStatus:      gw.BamenStatus,
			BamenYueStatus:   gw.BamenYueStatus,
			TianPanGanStatus: gw.TianPanGanStatus,
			DiPanGanStatus:   gw.DiPanGanStatus,
			GongStatus:       gw.GongStatus,
			IsNeiPan:         gw.IsNeiPan,
			JiTianGanStatus:  gw.JiTianGanStatus,
			JiDiGanStatus:    gw.JiDiGanStatus,
		}
		gws = append(gws, g)
	}
	return gws
}

func getJuShu(shu int, yinYang bool) string {
	yinYangDun := "阳遁"
	if !yinYang {
		yinYangDun = "阴遁"
	}

	juMap := map[int]string{
		1: "一局", 2: "二局", 3: "三局", 4: "四局", 5: "五局",
		6: "六局", 7: "七局", 8: "八局", 9: "九局",
	}

	if js, exists := juMap[shu]; exists {
		return yinYangDun + js
	}
	return yinYangDun + "未知局"
}

// 辅助函数
func contains(arr []int, target int) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
}
