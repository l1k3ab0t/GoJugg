package main

import (
	"github.com/l1k3ab0t/GoJugg/lib/FormatHTML"
	"github.com/l1k3ab0t/GoJugg/lib/GameEngine"
	"github.com/l1k3ab0t/GoJugg/lib/ReadConfig"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type settings struct {
	webcfg    bool
	name      string
	list      string
	gameMode  int
	groupCont int
}
type data struct {
	Name        string
	Table       template.HTML
	Bracket     template.HTML
	TextPhrase1 string
	TextPhrase2 string
	I           int
	I2          int
}

var s settings
var gg [][][]GameEngine.Game
var t []GameEngine.Team
var defaultData data
var setupDone bool
var templates = template.Must(template.ParseFiles("resources/control.html", "resources/tournament.html", "resources/games.html", "resources/setup.html", "resources/submitResult.html", "resources/rank.html"))

func loadCFG() (settings, []GameEngine.Team) {
	var s settings
	var t []GameEngine.Team
	cfg := ReadConfig.ReadFile("config.cfg")
	for _, v := range cfg {
		if ReadConfig.SplitConfig(v.Content)[0] == "EnableWebConfig" {
			if ReadConfig.SplitConfig(v.Content)[1] == "TRUE" {
				s.webcfg = true
			} else {
				s.webcfg = false
			}
		}
		if ReadConfig.SplitConfig(v.Content)[0] == "TournamentName" {
			s.name = ReadConfig.SplitConfig(v.Content)[1]
		}
		if ReadConfig.SplitConfig(v.Content)[0] == "List" {
			s.list = ReadConfig.SplitConfig(v.Content)[1]
		}
		if ReadConfig.SplitConfig(v.Content)[0] == "GameMode" {

			i, err := strconv.Atoi(ReadConfig.SplitConfig(v.Content)[1])
			if err == nil {
				s.gameMode = i
			} else {
				log.Println("Wrong Value set in Gamemode")
			}
		}
		if ReadConfig.SplitConfig(v.Content)[0] == "GroupCount" {

			i, err := strconv.Atoi(ReadConfig.SplitConfig(v.Content)[1])
			if err == nil {
				s.groupCont = i - 1
			} else {
				log.Println("Wrong Value set in GroupCount")
			}
		}

	}

	log.Println(s)
	teamList := ReadConfig.ReadFile(s.list)
	for i := range teamList {
		log.Println(strconv.Itoa(teamList[i].Linenumber) + " " + teamList[i].Content)
		t = append(t, GameEngine.Team{teamList[i].Content, teamList[i].Linenumber, teamList[i].Linenumber, 1, nil, 0})
	}

	return s, t
}

func defaultConnection(w http.ResponseWriter, r *http.Request) {
	if setupDone {
		http.Redirect(w, r, "/tournament/", http.StatusFound)
	} else {
		http.Redirect(w, r, "/setup/", http.StatusFound)
	}
}

func controlPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "control", defaultData)
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func setupPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "setup", defaultData)
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func addTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	log.Println("Add Team: ", r.Form["Team"])
	http.Redirect(w, r, "/setup/", http.StatusFound)
}

func rank(w http.ResponseWriter, r *http.Request) {
	tdata := defaultData
	tdata.TextPhrase1 = "Overall"
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	log.Println("Round: ", r.FormValue("Group"))

	if r.Form["CGroup"] != nil {
		log.Println("Group Changed")

		iP, _ := strconv.Atoi(r.FormValue("Group"))
		if iP < 6 {
			tdata.I = iP + 1
		} else {
			tdata.I = 0
		}
		tdata.TextPhrase1 = "from Group " + strconv.Itoa(tdata.I)
		tdata.Table = FormatHTML.FormatRanking(GameEngine.SortByRank(gg[tdata.I], GameEngine.BuildGroups(s.groupCont, t)[tdata.I]))
	} else if r.Form["Custom"] != nil {

		log.Println("Custom Destination")
		iP, _ := strconv.Atoi(r.FormValue("Group"))
		if iP <= s.groupCont {
			tdata.I = iP
		} else {
			tdata.I = 0
		}
		tdata.TextPhrase1 = "from Group " + strconv.Itoa(tdata.I)
		tdata.Table = FormatHTML.FormatRanking(GameEngine.SortByRank(gg[tdata.I], GameEngine.BuildGroups(s.groupCont, t)[tdata.I]))

	} else {
		tdata.I = 0
		tdata.Table = FormatHTML.FormatRanking(GameEngine.SortByRankInTourney(gg, GameEngine.BuildGroups(s.groupCont, t)))
	}
	renderTemplate(w, "rank", tdata)
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func games(w http.ResponseWriter, r *http.Request) {
	tdata := defaultData
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	log.Println("Group: ", r.FormValue("Page"))
	log.Println("Round: ", r.FormValue("Round"))

	if r.Form["CGroup"] != nil {
		log.Println("Group Changed")
		tdata.I2 = 0

		iP, _ := strconv.Atoi(r.FormValue("Page"))
		if iP < 6 {
			tdata.I = iP + 1
		} else {
			tdata.I = 0
		}
		tdata.Bracket = FormatHTML.FormatBracket(gg[tdata.I][tdata.I2])
	} else if r.Form["CRound"] != nil {
		log.Println("Round Changed")
		iR, _ := strconv.Atoi(r.FormValue("Round"))
		if iR < 6 {
			tdata.I2 = iR + 1
		} else {
			tdata.I2 = 0
		}
		iP, _ := strconv.Atoi(r.FormValue("Page"))
		tdata.I = iP
		tdata.Bracket = FormatHTML.FormatBracket(gg[tdata.I][tdata.I2])
	} else if r.Form["Custom"] != nil {

		log.Println("Custom Destination")
		iP, _ := strconv.Atoi(r.FormValue("Page"))
		if iP <= s.groupCont {
			tdata.I = iP
		} else {
			tdata.I = 0
		}

		iR, _ := strconv.Atoi(r.FormValue("Round"))
		if iR < 6 {
			tdata.I2 = iR
		} else {
			tdata.I2 = 0
		}
		tdata.Bracket = FormatHTML.FormatBracket(gg[tdata.I][tdata.I2])
	} else {
		tdata.Bracket = FormatHTML.FormatBracket(gg[0][0])
	}

	renderTemplate(w, "games", tdata)
	/*
		rank := GameEngine.SortByRank(gg[0], t)
		log.Println(rank)
		for i, v := range rank {
			log.Println("Team:", i, " Stats ", v)
		}
		log.Printf("IP: " + r.RemoteAddr + " connected to /games")
	*/
}

func submitResult(w http.ResponseWriter, r *http.Request) {
	tdata := defaultData
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	log.Println("Team 1: ", r.FormValue("T1ID"))
	log.Println("Team 2: ", r.FormValue("T2ID"))
	iD1, _ := strconv.Atoi(r.FormValue("T1ID"))
	iD2, _ := strconv.Atoi(r.FormValue("T2ID"))
	if r.Form["T1ID"] != nil && r.Form["T2ID"] != nil {
		for _, v := range gg {
			for _, v2 := range v {
				for _, v3 := range v2 {
					if v3.Opponent1.ID == iD1 && v3.Opponent2.ID == iD2 {
						tdata.TextPhrase1 = v3.Opponent1.Name
						tdata.TextPhrase2 = v3.Opponent2.Name
						break
					} else if v3.Opponent2.ID == iD1 && v3.Opponent1.ID == iD2 {
						tdata.TextPhrase1 = v3.Opponent2.Name
						tdata.TextPhrase2 = v3.Opponent1.Name
						break
					}
				}
			}
		}
		renderTemplate(w, "submitResult", tdata)
	} else {
		http.Redirect(w, r, "/games", http.StatusFound)
	}
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func receiveResult(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	log.Println("Team 1 (", r.FormValue("Team1"), ") Juggs: ", r.FormValue("T1Juggs"))
	log.Println("Team 2 (", r.FormValue("Team2"), ") Juggs: ", r.FormValue("T2Juggs"))
	if r.Form["Result"] != nil {
		t1j, _ := strconv.Atoi(r.FormValue("T1Juggs"))
		t2j, _ := strconv.Atoi(r.FormValue("T2Juggs"))
		if t1j == t2j {
			log.Println("can´t accept Result")
			http.Redirect(w, r, "/games", http.StatusFound)
		} else {
			for i1, v1 := range gg {
				for i2, v2 := range v1 {
					for i3, v3 := range v2 {
						if v3.Opponent1.Name == r.FormValue("Team1") && v3.Opponent2.Name == r.FormValue("Team2") {
							gg[i1][i2][i3].Result.Team1Juggs = t1j
							gg[i1][i2][i3].Result.Team2Juggs = t2j
							log.Println("1")
							log.Println(gg[i1][i2][i3])
							break
						} else if v3.Opponent2.Name == r.FormValue("Team1") && v3.Opponent1.Name == r.FormValue("Team2") {
							gg[i1][i2][i3].Result.Team1Juggs = t2j
							gg[i1][i2][i3].Result.Team2Juggs = t1j
							log.Println("2")
							log.Println(gg[i1][i2][i3])
							break
						}
					}
				}
			}
		}

	}
	http.Redirect(w, r, "/games", http.StatusFound)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "tournament", defaultData)
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func renderTemplate(w http.ResponseWriter, tmpl string, t data) {
	err := templates.ExecuteTemplate(w, tmpl+".html", t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
}

func main() {
	setupDone = false
	s, t = loadCFG()
	setupDone = !s.webcfg
	var gq [][]GameEngine.Game
	if s.gameMode == 0 {
		gpGroups := GameEngine.BuildGroups(s.groupCont, t)
		for i := 0; i <= s.groupCont; i++ {
			for u := 0; u <= 6; u++ {
				gq = append(gq, GameEngine.BuildGroupGames(gpGroups[i]))
				gpGroups[i] = GameEngine.GroupPlayed(gpGroups[i], gq[u])
			}
			gg = append(gg, gq)
			gq = nil
		}
		for i, v := range gg {
			for i2, v2 := range v {
				log.Println("Group:", i, " Round ", i2, " ", v2)
			}
		}
	}

	//r := GameEngine.SortByRank(gg[0], GameEngine.BuildGroups(s.groupCont, t)[0])
	//log.Println("Rank: ", r)

	/*
			gg := GameEngine.BuildGroups(s.groupCont, t)
			g := GameEngine.BuildGroupGames(gg[3])
			log.Println("Test  ",gg[3])
			for _, v := range g {
				log.Println(v.Opponent1, " vs ", v.Opponent2)
			}
			gg[3] = GameEngine.GroupPlayed(gg[3], g)
			g = GameEngine.BuildGroupGames(gg[3])

			for _, v := range g {
				log.Println(v.Opponent1, " vs ", v.Opponent2)

			}
			defaultData.Bracked=FormatHTML.FormatBracket(g)
		}
	*/
	defaultData.Name = s.name
	defaultData.Table = FormatHTML.FormatTeamLIst(t)
	defaultData.Bracket = FormatHTML.FormatBracket(gg[0][0])
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/", defaultConnection)
	http.HandleFunc("/control/", controlPage)
	http.HandleFunc("/setup/", setupPage)
	http.HandleFunc("/rank", rank)
	http.HandleFunc("/games", games)
	http.HandleFunc("/addTeam", addTeam)
	http.HandleFunc("/submitResult", submitResult)
	http.HandleFunc("/receiveResult", receiveResult)
	http.HandleFunc("/tournament/", mainPage)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
