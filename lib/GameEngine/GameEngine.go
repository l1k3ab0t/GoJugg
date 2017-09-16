package GameEngine

import (
	"log"
)

type Team struct {
	Name        string
	Rank        int
	ID          int
	Group       int
	PlayedVS    []int
	PlayedGames int
}

type Game struct {
	Opponent1 Team
	Opponent2 Team
	//Result    GameResult
}

type GameResult struct {
	winner Team
}

func BuildGroups(gCount int, teams []Team) [][]Team {
	x := 0
	tg := make([][]Team, gCount+1)
	for _, t := range teams {

		if x <= gCount {
			t.Group = x
			tg[x] = append(tg[x], t)
			x++
		}
		if x > gCount {
			x = 0
		}
	}
	return tg
}

func BuildGroupGames(teams []Team) []Game {
	var games []Game
	teams = sortByLeastPlayed(teams)
	//log.Println(teams)
	for i := 0; i < len(teams); i++ {
		for i2 := 0; i2 < len(teams); i2++ {
			if len(teams) >= 1 {

				if !playedAgainst(teams[i].ID, teams[i2]) && teams[i].ID != teams[i2].ID {
					log.Println(i, " ", i2)
					games = append(games, buildGame(teams[i], teams[i2]))
					teams = findAndRemove(teams[i], teams)
					teams = findAndRemove(teams[i2-1], teams)
					i = 0
					i2 = 0

					log.Println("found")
				}
			} else {
				break
			}
		}
	}
	return games
}

func buildGame(op1 Team, op2 Team) Game {
	return Game{op1, op2}
}

func playedAgainst(ID int, t Team) bool {
	for _, v := range t.PlayedVS {
		if ID == v {
			return true
		}
	}
	return false
}
func sortByLeastPlayed(teams []Team) []Team {
	var ts []Team
	for range teams {
		t := leastPlayed(teams)
		ts = append(ts, t)
		teams = findAndRemove(t, teams)
	}
	return ts
}

func leastPlayed(teams []Team) Team {
	lp := teams[0]
	for _, v := range teams {
		if v.PlayedGames < lp.PlayedGames {
			lp = v
		}
	}
	return lp
}
func findAndRemove(t Team, teams []Team) []Team {
	for i, v := range teams {
		if v.ID == t.ID {
			teams = append(teams[:i], teams[i+1:]...)
		}
	}
	log.Println("Removed ", t.Name, " ", teams)
	return teams
}
