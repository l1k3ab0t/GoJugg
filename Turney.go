package main

import (
	"fmt"
	"github.com/l1k3ab0t/GoJugg/lib/GameEngine"
	"github.com/l1k3ab0t/GoJugg/lib/ReadConfig"
	"log"
	"strconv"
	"net/http"
	"html/template"
)

type settings struct {
	name string
	list      string
	gameMode  int
	groupCont int

}

//var tName string
var setupDone   bool
var templates = template.Must(template.ParseFiles("resources/control.html","resources/tournament.html","resources/setup.html"))

func loadCFG() (settings, []GameEngine.Team) {
	var s settings
	var t []GameEngine.Team
	cfg := ReadConfig.ReadFile("config.cfg")
	for _, v := range cfg {
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

func mainPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "tournament")
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	setupDone=false
	s, t := loadCFG()
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

	//tName=s.name

	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/", defaultConnection)
	http.HandleFunc("/control/", controlPage)
	http.HandleFunc("/setup/", setupPage)
	http.HandleFunc("/tournament/", mainPage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
