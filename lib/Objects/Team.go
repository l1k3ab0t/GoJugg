package Objects

type Team struct {
	Name        string	`xml:"TeamName"`
	Rank        int		`xml:"TeamRank"`
	ID          int		`xml:"TeamID"`
	Group       int		`xml:"TeamGroup"`
	PlayedVsID  []int	`xml:"PlayedOPs"`
	PlayedGames int		`xml:"PlayedGames"`
}
