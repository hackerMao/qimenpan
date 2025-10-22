package sizhu

import (
	"fmt"
	"log"
	"qimenpan/pkg/dizhi"
	"qimenpan/pkg/tiangan"
)

type SiZhu struct {
	YearGan  tiangan.TianGan
	YearZhi  dizhi.DiZhi
	MonthGan tiangan.TianGan
	MonthZhi dizhi.DiZhi
	DayGan   tiangan.TianGan
	DayZhi   dizhi.DiZhi
	HourGan  tiangan.TianGan
	HourZhi  dizhi.DiZhi
}

func (s *SiZhu) calcShiZhu(hour int) error {
	// 确定时辰的地支
	shiChenIndex := (hour + 1) / 2 % 12
	shiZhi := dizhi.ALL[shiChenIndex]

	// 根据日干确定起始天干
	startGan := map[string]int{
		"甲": 0, "乙": 2, "丙": 4, "丁": 6, "戊": 8,
		"己": 0, "庚": 2, "辛": 4, "壬": 6, "癸": 8,
	}[s.DayGan.Name]

	// 计算时干
	hourGanIndex := (startGan + shiChenIndex) % 10
	s.HourGan = tiangan.ALL[hourGanIndex]
	s.HourZhi = shiZhi

	return nil
}

func New(ng, nz, mg, mz, dg, dz string, hour int) (SiZhu, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("四柱转换失败")
		}
	}()

	siZhu := SiZhu{
		YearGan:  tiangan.Match(ng),
		YearZhi:  dizhi.Match(nz),
		MonthGan: tiangan.Match(mg),
		MonthZhi: dizhi.Match(mz),
		DayGan:   tiangan.Match(dg),
		DayZhi:   dizhi.Match(dz),
	}

	err := siZhu.calcShiZhu(hour)
	if err != nil {
		return SiZhu{}, err
	}

	return siZhu, nil
}

func (s *SiZhu) Print() {
	fmt.Printf("干支历：%s%s年 %s%s月 %s%s日 %s%s时",
		s.YearGan.Name, s.YearZhi.Name,
		s.MonthGan.Name, s.MonthZhi.Name,
		s.DayGan.Name, s.DayZhi.Name,
		s.HourGan.Name, s.HourZhi.Name)
}

func (s *SiZhu) String() string {
	return fmt.Sprintf("%s%s年 %s%s月 %s%s日 %s%s时",
		s.YearGan.Name, s.YearZhi.Name,
		s.MonthGan.Name, s.MonthZhi.Name,
		s.DayGan.Name, s.DayZhi.Name,
		s.HourGan.Name, s.HourZhi.Name)
}
