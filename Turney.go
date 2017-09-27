package main

import (
	"github.com/l1k3ab0t/GoJugg/lib/GameEngine"
	"github.com/l1k3ab0t/GoJugg/lib/ReadConfig"
	"log"
	"strconv"
	"net/http"
	"html/template"
	"github.com/l1k3ab0t/GoJugg/lib/FormatHTML"
	"fmt"
)

type settings struct {
	webcfg bool
	name string
	list      string
	gameMode  int
	groupCont int

}
type data struct {
	Name string
	Table template.HTML

}

var tdata data
//var tName string
var setupDone   bool
var templates = template.Must(template.ParseFiles("resources/control.html","resources/tournament.html","resources/setup.html"))

func loadCFG() (settings, []GameEngine.Team) {
	var s settings
	var t []GameEngine.Team
	cfg := ReadConfig.ReadFile("config.cfg")
	for _, v := range cfg {
		if ReadConfig.SplitConfig(v.Content)[0] == "EnableWebConfig" {
			if ReadConfig.SplitConfig(v.Content)[1]=="TRUE" {
				s.webcfg = true
			}else {
				s.webcfg = false
			}
		}
		if ReadConfig.SplitConfig(v.Content)[0] == "TournamentName" {
			s.name=ReadConfig.SplitConfig(v.Content)[1]
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
	}else {
		http.Redirect(w, r, "/setup/", http.StatusFound)
	}
}

func controlPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "control")
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func setupPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "setup")
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func addTeam(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		r.ParseForm()
		// logic part of log in
		log.Println("Add Team: ", r.Form["Team"])
	} else {
		r.ParseForm()
		// logic part of log in
		log.Println("Add Team: ", r.Form["Team"])
	}
	http.Redirect(w, r, "/setup/", http.StatusFound)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "tournament")
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", tdata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	setupDone=false
	s, t := loadCFG()
	setupDone=!s.webcfg
	if s.gameMode == 0 {
		tg := GameEngine.BuildGroups(s.groupCont, t)
		g := GameEngine.BuildGroupGames(tg[3])
		log.Println("Test  ",tg[3])
		for _, v := range g {
			log.Println(v.Opponent1, " vs ", v.Opponent2)
		}
		tg[3] = GameEngine.GroupPlayed(tg[3], g)
		g = GameEngine.BuildGroupGames(tg[3])

		for _, v := range g {
			log.Println(v.Opponent1, " vs ", v.Opponent2)

		}
	}

	tdata.Name=s.name
	tdata.Table=FormatHTML.FormatTeamLIst(t)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/", defaultConnection)
	http.HandleFunc("/control/", controlPage)
	http.HandleFunc("/setup/", setupPage)
	http.HandleFunc("/addTeam", addTeam)
	http.HandleFunc("/tournament/", mainPage)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
