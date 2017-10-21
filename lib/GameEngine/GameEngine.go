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
	Result    GameResult
}

type GameResult struct {
	Team1Juggs int
	Team2Juggs int
}

type Rank struct {
	Rank   int
	TName  string
	Result GameResult
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

func BuildGroupGames(t []Team) []Game {
	var games []Game
	var teams []Team
	teams = append(teams, t...)
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
	log.Println("Teams ", teams)
	return games
}

func buildGame(op1 Team, op2 Team) Game {
	return Game{op1, op2, GameResult{0, 0}}
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

func GroupPlayed(teams []Team, games []Game) []Team {
	for _, v := range games {
		for i := range teams {
			if v.Opponent1.ID == teams[i].ID {
				teams[i] = played(teams[i], v.Opponent2)
			}
			if v.Opponent2.ID == teams[i].ID {
				teams[i] = played(teams[i], v.Opponent1)
			}
		}
	}
	return teams
}

func played(t Team, opponent Team) Team {
	t.PlayedVS = append(t.PlayedVS, opponent.ID)
	return t
}

// Ranking ------------------------------------------------------------------------------------------------------

func SortByRankInTourney(games [][][]Game, teams [][]Team) []Rank {
	var r []Rank
	log.Println(len(games))
	for i, v := range games {
		for i2, v2 := range v {
			log.Println("Group:", i, " Round ", i2, " ", v2)
		}
	}
	log.Println(len(teams))
	for _, v := range teams {
		for _, v2 := range v {
			r = append(r, Rank{1, v2.Name, getResults(v2, games[v2.Group])})
		}
	}
	rs := sortByPoints(r)
	log.Println(len(rs))
	for i, v := range rs {
		log.Println("Team:", i, " Stats ", v)
	}
	return rs
}

func SortByRank(games [][]Game, teams []Team) []Rank {
	var r []Rank
	log.Println(len(games))
	for i, v := range games {
		for i2, v2 := range v {
			log.Println("Group:", i, " Round ", i2, " ", v2)
		}
	}
	log.Println(len(teams))
	for _, v := range teams {
		r = append(r, Rank{1, v.Name, getResults(v, games)})
	}
	rs := sortByPoints(r)
	log.Println(len(rs))
	for i, v := range rs {
		log.Println("Team:", i, " Stats ", v)
	}
	return rs
}

func sortByPoints(list []Rank) []Rank {
	var r []Rank
	saveu := 0
	i := 0
	var hRank []Rank
	for range list {
		if len(list) != 0 {
			hRank = highestByPoints(list)
			log.Println("highestRank: ", hRank)
			multiR := i + 1 + saveu
			for u, v := range hRank {
				saveu = u
				v.Rank = multiR
				i = multiR
				r = append(r, v)
				log.Println("TeamID:", v, " Result ", v.Result)
				list = findAndRemoveRanking(v.TName, list)
			}
		}
	}
	return r
}

func highestByPoints(list []Rank) []Rank {
	var hRank []Rank
	u := -100
	for _, v := range list {
		//log.Println("TeamID:", k, " Result ",v)
		i := v.Result.Team1Juggs - v.Result.Team2Juggs
		if i > u {
			u = i
			hRank = nil
			hRank = append(hRank, v)

		} else if i == u {
			hRank = append(hRank, v)
		}
	}
	log.Println(hRank)
	return hRank
}

func getResults(t Team, games [][]Game) GameResult {
	var result GameResult
	for _, v := range games {
		for _, v2 := range v {
			if t.ID == v2.Opponent1.ID {
				result.Team1Juggs = result.Team1Juggs + v2.Result.Team1Juggs
				result.Team2Juggs = result.Team2Juggs + v2.Result.Team2Juggs

			} else if t.ID == v2.Opponent2.ID {
				result.Team1Juggs = result.Team1Juggs + v2.Result.Team2Juggs
				result.Team2Juggs = result.Team2Juggs + v2.Result.Team1Juggs
			}
		}
	}
	return result
}

func findAndRemoveRanking(teamName string, ranking []Rank) []Rank {
	for i, v := range ranking {
		if v.TName == teamName {
			ranking = append(ranking[:i], ranking[i+1:]...)
		}
	}
	log.Println("Removed ", teamName, " ", ranking)
	return ranking
}
