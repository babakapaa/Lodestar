package sandbox

type Estimate struct {
	MaxHours  float64          `json:"maxHours"`
	Starter   bool             `json:"starter"`
	Available float64          `json:"available"`
	Night     *NightWindow     `json:"night"`
	Options   []EstimateOption `json:"options"`
	Suggested *string          `json:"suggested"`
}

type NightWindow struct {
	From     int    `json:"from"`
	To       int    `json:"to"`
	Discount int    `json:"discount"`
	TimeZone string `json:"timeZone"`
}

type EstimateOption struct {
	ClassID   string  `json:"classId"`
	ClassName string  `json:"className"`
	MemoryGB  int     `json:"memoryGb"`
	Rate      float64 `json:"rate"`
	Cost      float64 `json:"cost"`
	NightCost float64 `json:"nightCost"`
	FreeGPUs  int     `json:"freeGpus"`
	Queued    int     `json:"queued"`
	Enough    bool    `json:"enough"`
}

func (o EstimateOption) better(than EstimateOption) bool {
	free, other := o.FreeGPUs > o.Queued, than.FreeGPUs > than.Queued
	switch {
	case free != other:
		return free
	case !free && o.Queued != than.Queued:
		return o.Queued < than.Queued
	default:
		return o.Cost < than.Cost
	}
}
