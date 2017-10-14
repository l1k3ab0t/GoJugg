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
	teams = append(teams,t...)
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
	log.Println("Teams ",teams)
	return games
}

func buildGame(op1 Team, op2 Team) Game {
	return Game{op1, op2,GameResult{0,0}}
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
	for _,v:=range games{
		for i:=range teams{
			if v.Opponent1.ID==teams[i].ID {
				teams[i]=played(teams[i],v.Opponent2)
			}
			if v.Opponent2.ID==teams[i].ID {
				teams[i]=played(teams[i],v.Opponent1)
			}
		}
	}
	return teams
}

func played(t Team, opponent Team) Team {
	t.PlayedVS =append(t.PlayedVS,opponent.ID)
	return t
}

func SortByRankInTourney(games [][][]Game, teams []Team ) map[int]GameResult{
	m := make(map[int]GameResult)
	log.Println(len(games))
	for i,v:=range games {
		for i2, v2 := range v {
			log.Println("Group:", i, " Round ",i2," ", v2)
		}
	}

	for _,v:=range teams {
		log.Println("Team:", v.Name, " Stats ",m[v.ID])
	}

	return nil
}

func SortByRankInGroup (games [][]Game, teams []Team ) map[int]GameResult{
	m := make(map[int]GameResult)
	log.Println(len(games))
	for i,v:=range games {
		for i2, v2 := range v {
			log.Println("Group:", i, " Round ",i2," ", v2)
		}
	}
	for _,v:=range teams {
		m[v.ID]=getResults(v,games)
	}
	u:=-100
	for k,v:=range m{
		//log.Println("TeamID:", k, " Result ",v)
		i:=v.Team1Juggs-v.Team2Juggs
		if  i>u{
			u=i

		}
	}


	for _,v:=range teams {
		log.Println("Team:", v.Name, " Stats ",m[v.ID])
	}

	return nil
}

func sortByPoints(m map[int]GameResult) map[int]GameResult{
	mSorted := make(map[int]GameResult)
	var highesKey []int
	u:=-100
	m2:=m
	for k,v:=range m{
		//log.Println("TeamID:", k, " Result ",v)
		{
			highesKey=hightesByPoints(m)
			if len(highesKey)>1 {
				mSorted[]
			}
		}
	}
}

func hightesByPoints(m map[int]GameResult) []int{
	var kSave []int
	u:=-100
	for k,v:=range m{
		//log.Println("TeamID:", k, " Result ",v)
		i:=v.Team1Juggs-v.Team2Juggs
		if  i>u{
			u=i
			kSave=nil
			kSave=append(kSave,k)

		}else if  i==u{
			kSave=append(kSave,k)
		}
	}
	return kSave
}

func getResults (t Team, games [][]Game ) GameResult{
	var result GameResult
	for _,v:=range games {
		for _,v2:=range v {
			if t.ID==v2.Opponent1.ID {
				result.Team1Juggs=result.Team1Juggs+v2.Result.Team1Juggs
				result.Team2Juggs=result.Team1Juggs+v2.Result.Team2Juggs

			}else if t.ID==v2.Opponent2.ID {
				result.Team1Juggs=result.Team1Juggs+v2.Result.Team2Juggs
				result.Team2Juggs=result.Team1Juggs+v2.Result.Team1Juggs
			}
		}
	}
	return result
	}