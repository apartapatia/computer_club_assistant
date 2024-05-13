package club

type Club struct {
	WorkingTime *WorkingTime
	Price       int
	MaxTables   int
}

func NewClub(workingTime *WorkingTime, price, tablesCount int) *Club {
	return &Club{
		WorkingTime: workingTime,
		Price:       price,
		MaxTables:   tablesCount,
	}
}
