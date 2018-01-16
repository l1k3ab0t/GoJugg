package Objects

type Game struct {
	Opponent1 Team		`xml:"OP1"`
	Opponent2 Team		`xml:"OP2"`
	Result    Result	`xml:"Result"`
	Field   string		`xml:"Field"`
}
