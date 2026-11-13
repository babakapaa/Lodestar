package study

type MyDefense struct {
	Format   DefenseFormat `json:"format"`
	Required bool          `json:"required"`
	Recorded bool          `json:"recorded"`
	Status   DefenseStatus `json:"status"`
	Slot     *DefenseSlot  `json:"slot"`
	Protocol *Defense      `json:"protocol"`
}
