package Objects

import (
	"log"
	"math"
	"strconv"
)

type Tournament struct {
	S        Settings
	teamList [][]Team
	SetupTeamList []Team
	games    [][][]Game
	ranked   []Rank
	gRanked  [][]Rank
}

type Rank struct {
	Rank   int
	TName  string
	Result Result
}

func changeTGroup(gCount int, teams []Team) []Team {
	x := 0
	for i := range teams {

		if x <= gCount {
			teams[i].Group = x
			x++
		}
		if x > gCount {
			x = 0
		}
	}
	return teams
}

func buildGroups(gCount int, teams []Team) [][]Team {
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

func buildGroupGames(t []Team, fields []int,rounds int) [][]Game {
	var games [][]Game
	var teams []Team
	teams = append(teams, t...)
	//log.Println(teams)
	for r := 0; r <= rounds; r++ {
		maxgames := 0
		var rgames []Game

		teams = sortByLeastPlayed(teams)
		for i := 0; i < len(teams); i++ {

			for i2 := 0; i2 < len(teams); i2++ {
				if len(teams) >= 1 {
					if !playedAgainst(teams[i].ID, teams[i2]) && teams[i].ID != teams[i2].ID && maxgames < fields[r] {
						log.Println(i, " ", i2)
						rgames = append(rgames, buildGame(teams[i], teams[i2]))
						teams = findAndRemove(teams[i], teams)
						teams = findAndRemove(teams[i2-1], teams)
						i = 0
						i2 = 0
						maxgames++

						log.Println("found")
					}
				} else {
					break
				}
			}
		}
		games = append(games, rgames)
		rgames = nil
		maxgames = 0
		teams = nil
		teams = append(teams, t...)
		teams = groupPlayed(teams, games)

	}
	log.Println("Teams ", teams)
	return games
}

func buildGame(op1 Team, op2 Team) Game {
	return Game{op1, op2, Result{0, 0},"0"}
}

func playedAgainst(ID int, t Team) bool {
	for _, v := range t.PlayedVsID {
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

func groupPlayed(teams []Team, games [][]Game) []Team {
	for _, v := range games {
		for _, v2 := range v {
			for i := range teams {
				if v2.Opponent1.ID == teams[i].ID {
					teams[i] = played(teams[i], v2.Opponent2)
					teams[i].PlayedGames++
				}
				if v2.Opponent2.ID == teams[i].ID {
					teams[i] = played(teams[i], v2.Opponent1)
					teams[i].PlayedGames++
				}
			}
		}
	}
	return teams
}

func played(t Team, opponent Team) Team {
	t.PlayedVsID = append(t.PlayedVsID, opponent.ID)
	return t
}

func StartSetup() Tournament{
	tour := Tournament{}
	tour.S, tour.SetupTeamList = Settings{}.LoadCFG("config.cfg")
	return tour
}

func (t *Tournament) FinishSetup(){

	var fieldsagroup [][]int				// round group
	if t.S.fields>=t.S.groupCount{
		f :=int(math.Floor(float64(t.S.fields+1)/float64(t.S.groupCount+1)))
		for i:=0; i<= t.S.roundCount;i++{
			var fields []int
			fields=nil
			for u:=0; u<= t.S.groupCount;u++{
				fields=append(fields,f)
			}
			fieldsagroup=append(fieldsagroup,fields)
		}
	}
	if (t.S.fields+1)%(t.S.groupCount+1)!=0{
		fieldsagroup=shuffleRemainingFields(fieldsagroup,int(math.Floor(float64(t.S.fields)/float64(t.S.groupCount))-1))
	}
	t.SetupTeamList = changeTGroup(t.S.groupCount, t.SetupTeamList)
	gpGroups := buildGroups(t.S.groupCount, t.SetupTeamList)

	for i := 0; i <= t.S.groupCount; i++ {
		t.games = append(t.games, buildGroupGames(gpGroups[i],convertFields(fieldsagroup, i) ,t.S.roundCount))

	}
	if t.S.customFieldNames==false {
		for i, v := range fieldsagroup {
			u := 0
			for i2, v2 := range v {
				for i3 := 0; i3 < v2; i3++ {
					log.Println("Field Nr:",u," Game: ",i2," ",i," ",i3)
					t.games[i2][i][i3].Field = strconv.Itoa(u)
					u++
				}
			}
		}
	}

	for i, v := range t.games {
		for i2, v2 := range v {
			log.Println("Group:", i, " Round ", i2, " ", v2)
		}
	}

	t.teamList = gpGroups
	t.initializeGRanked()
	t.UpdateRanked()
}

func shuffleRemainingFields(fieldMap [][]int, rFields int) [][]int{
	var fairfields []int
	for range fieldMap[0]{
		fairfields=append(fairfields,0)
	}
	for i:=range fieldMap {
		for u := 0; u < rFields; u++ {
			fieldMap[i][leastFields(fairfields)]++
			fairfields[leastFields(fairfields)]++
		}
	}
	return fieldMap
}

func leastFields(fairFields []int) int {
	lp := fairFields[0]
	lpI:=0
	for i, v := range fairFields {
		if v < lp {
			lp = v
			lpI=i
		}
	}
	return lpI
}

func convertFields(a [][]int,i int) []int{
	var r []int
	for _,v:=range a{
		r=append(r,v[i])
	}
	return r
}
func NewTournament() Tournament {
	var t []Team
	tour := Tournament{}
	tour.S, t = Settings{}.LoadCFG("config.cfg")
	var fieldsagroup [][]int

	if tour.S.fields>=tour.S.groupCount{
		f :=int(math.Floor(float64(tour.S.fields+1)/float64(tour.S.groupCount+1)))
		for i:=0; i<= tour.S.roundCount;i++{
			var fields []int
			fields=nil
			for u:=0; u<= tour.S.groupCount;u++{
				fields=append(fields,f)
			}
			fieldsagroup=append(fieldsagroup,fields)
		}
	}
	if (tour.S.fields+1)%(tour.S.groupCount+1)!=0{
			fieldsagroup=shuffleRemainingFields(fieldsagroup,int(math.Floor(float64(tour.S.fields)/float64(tour.S.groupCount))-1))
	}

	t = changeTGroup(tour.S.groupCount, t)
	gpGroups := buildGroups(tour.S.groupCount, t)
	for i := 0; i <= tour.S.groupCount; i++ {
		tour.games = append(tour.games, buildGroupGames(gpGroups[i],convertFields(fieldsagroup, i) ,tour.S.roundCount))
	}
	for i, v := range tour.games {
		for i2, v2 := range v {
			log.Println("Group:", i, " Round ", i2, " ", v2)
		}
	}

	tour.teamList = gpGroups
	tour.initializeGRanked()
	tour.UpdateRanked()
	return tour
}

//-------------------------------------------------------------------------------------------------

func (t Tournament) TeamList() []Team {
	var teams []Team
	for _, v := range t.teamList {
		teams = append(teams, v...)
	}
	return teams
}

func (t Tournament) Games() [][][]Game {
	return t.games
}

func (t *Tournament) SetGameResult(i1, i2, i3 int, res Result) {
	t.games[i1][i2][i3].Result = res
}

func (t Tournament) Ranked() []Rank {
	return t.ranked
}

func (t Tournament) GroupRanked(GroupID int) []Rank {
	if len(t.gRanked) < GroupID {
		return nil
	} else {
		return t.gRanked[GroupID]
	}

}

func (t *Tournament) UpdateRanked() {
	var r []Rank
	log.Println(len(t.games))
	for i, v := range t.games {
		for i2, v2 := range v {
			log.Println("Group:", i, " Round ", i2, " ", v2)
		}
	}
	log.Println(len(t.teamList))
	for _, v := range t.teamList {
		for _, v2 := range v {
			r = append(r, Rank{1, v2.Name, getResults(v2, t.games[v2.Group])})
		}
	}
	t.ranked = sortByPoints(r)
	log.Println(len(t.ranked))
	for i, v := range t.ranked {
		log.Println("Team:", i, " Stats ", v)
	}
}

func (t *Tournament) initializeGRanked() {
	var r []Rank
	log.Println(len(t.games))
	for i, v := range t.games {
		for i2, v2 := range v {
			log.Println("Group:", i, " Round ", i2, " ", v2)
		}
	}

	for i, v := range t.teamList {
		for _,v2 :=range v{
			r = append(r, Rank{1, v2.Name, getResults(v2, t.games[i])})
		}
		t.gRanked =append(t.gRanked, sortByPoints(r))
		r=nil
		log.Println(len(t.gRanked[i]))
		for i, v := range t.gRanked[i] {
			log.Println("Team:", i, " Stats ", v)
		}
	}


}

func (t *Tournament) UpdateGRanked() {
	for i := range t.teamList {
		t.UpdateSingleGRanked(i)
	}
}

func (t *Tournament) UpdateSingleGRanked(Group int) {
	var r []Rank
	log.Println(len(t.games))
	for i, v := range t.games[Group] {
		for i2, v2 := range v {
			log.Println("Group:", i, " Round ", i2, " ", v2)
		}
	}
	log.Println(len(t.teamList[Group]))
	for _, v := range t.teamList[Group] {
		r = append(r, Rank{1, v.Name, getResults(v, t.games[Group])})
	}
	t.gRanked[Group] = sortByPoints(r)
	log.Println(len(t.gRanked[Group]))
	for i, v := range t.gRanked[Group] {
		log.Println("Team:", i, " Stats ", v)
	}
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

func getResults(t Team, games [][]Game) Result {
	var result Result
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

func (t Tournament) TeamRank(team Team) Rank {
	var r Rank
	for _, v := range t.ranked {
		if v.TName == team.Name {
			return v
		}
	}
	return r
}

func (t Tournament) TeamGRank(team Team) Rank {
	var r Rank
	for _, v := range t.gRanked {
		for _, v2 := range v {
			if v2.TName == team.Name {
				return v2
			}
		}
	}
	return r
}

func (t Tournament) TeamGames(team Team) []Game {
	var g []Game
	for _, v := range t.games {
		for _, v2 := range v {
			for _, v3 := range v2 {
				if v3.Opponent1.Name == team.Name || v3.Opponent2.Name == team.Name {
					g = append(g, v3)
				}
			}
		}
	}
	return g
}

func (t Tournament) TeamByName(name string) Team {
	var team Team
	for _, v := range t.teamList {
		for _, v2 := range v {
			if v2.Name == name {
				return v2
			}
		}
	}
	return team
}
