package GameEngine

type Team struct {
	Name     string
	Rank     int
	Startpos int
	Group    int
}

func BuildGroups(gCount int, t []Team) []Team {
	x := 0
	for i := range t {

		if x <= gCount {
			t[i].Group = x
			x++
		}
		if x > gCount {
			x = 0
		}
	}
	return t
}
