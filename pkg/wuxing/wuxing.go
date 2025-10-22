package wuxing

// WuXing 五行
// 金木水火土
type WuXing struct {
	Name string `json:"name"`
}

func (f *WuXing) GetName() string {
	return f.Name
}

var (
	JIN  = WuXing{Name: "金"}
	MU   = WuXing{Name: "木"}
	SHUI = WuXing{Name: "水"}
	HUO  = WuXing{Name: "火"}
	TU   = WuXing{Name: "土"}

	genMap = map[WuXing]WuXing{
		JIN:  SHUI,
		MU:   HUO,
		SHUI: MU,
		HUO:  TU,
		TU:   JIN,
	}
	keMap = map[WuXing]WuXing{
		JIN:  MU,
		MU:   TU,
		SHUI: HUO,
		HUO:  JIN,
		TU:   SHUI,
	}
)

func IsKe(x1, x2 WuXing) bool {
	ke := keMap[x1]
	return ke.Name == x2.Name
}

func (f *WuXing) Gen() WuXing {
	return genMap[*f]
}

func (f *WuXing) Ke() WuXing {
	return keMap[*f]
}

func (f *WuXing) Tong(other *WuXing) bool {
	return f.Name == other.Name
}
