package main

import (
	"github.com/l1k3ab0t/GoJugg/lib/GameEngine"
	"github.com/l1k3ab0t/GoJugg/lib/ReadConfig"
	"log"
	"strconv"
	"net/http"
	"html/template"
	"github.com/l1k3ab0t/GoJugg/lib/FormatHTML"
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
	Bracked template.HTML
	I int
	I2 int

}
var s settings
var tg [][][]GameEngine.Game
var t []GameEngine.Team
var defaultData data
//var tName string
var setupDone   bool
var templates = template.Must(template.ParseFiles("resources/control.html","resources/tournament.html","resources/games.html","resources/setup.html"))

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
	renderTemplate(w, "control",defaultData)
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func setupPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "setup",defaultData)
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func addTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	log.Println("Add Team: ", r.Form["Team"])
	http.Redirect(w, r, "/setup/", http.StatusFound)
}

func games(w http.ResponseWriter, r *http.Request) {
	tdata:=defaultData
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	log.Println("Group: ", r.FormValue("Page"))
	log.Println("Round: ", r.FormValue("Round"))

	if r.Form["CGroup"]!=nil {

		log.Println("Group Changed")
		tdata.I2=0

		iP,_:=strconv.Atoi(r.FormValue("Page"))
		if iP<6 {
			tdata.I = iP + 1
		}else {
			tdata.I=0
		}
		tdata.Bracked=FormatHTML.FormatBracket(tg[tdata.I][tdata.I2])
	}

	if r.Form["CRound"]!=nil {

		log.Println("Round Changed")
		iR,_:=strconv.Atoi(r.FormValue("Round"))
		if iR<6 {
			tdata.I2 = iR + 1
		}else {
			tdata.I2=0
		}
		iP,_:=strconv.Atoi(r.FormValue("Page"))
		tdata.I = iP
		tdata.Bracked=FormatHTML.FormatBracket(tg[tdata.I][tdata.I2])
	}
	if r.Form["Custom"]!=nil {

		log.Println("Custom Destination")
		iP,_:=strconv.Atoi(r.FormValue("Page"))
		if iP<6 {
			tdata.I = iP
		}else {
			tdata.I=0
		}

		iR,_:=strconv.Atoi(r.FormValue("Round"))
		if iR<6 {
			tdata.I2 = iR
		}else {
			tdata.I2=0
		}
		tdata.Bracked=FormatHTML.FormatBracket(tg[tdata.I][tdata.I2])
	}

	renderTemplate(w, "games",tdata)
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "tournament",defaultData)
	log.Printf("IP: " + r.RemoteAddr + " connected")
}


func renderTemplate(w http.ResponseWriter, tmpl string,t data) {
	err := templates.ExecuteTemplate(w, tmpl+".html", t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
}

func main() {
	setupDone=false
	s, t = loadCFG()
	setupDone=!s.webcfg
	var gq [][]GameEngine.Game
	if s.gameMode == 0 {
		gpGroups := GameEngine.BuildGroups(s.groupCont, t)
		for i:=0;i<=s.groupCont;i++{
			for u:=0;u<=6;u++{
				gq=append(gq,GameEngine.BuildGroupGames(gpGroups[i]))
				//tg[i][u]=append(tg[i][u],GameEngine.BuildGroupGames(gpGroups[i]))
				gpGroups[i]=GameEngine.GroupPlayed(gpGroups[i],gq[u])
			}
			tg=append(tg,gq)
			gq=nil
		}
		for i,v:=range tg {
			for i2, v2 := range v {
				log.Println("Group:", i, " Round ",i2," ", v2)
			}
		}
	}


		/*
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
		defaultData.Bracked=FormatHTML.FormatBracket(g)
	}
	*/
	defaultData.Name=s.name
	defaultData.Table=FormatHTML.FormatTeamLIst(t)
	defaultData.Bracked=FormatHTML.FormatBracket(tg[0][0])
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/", defaultConnection)
	http.HandleFunc("/control/", controlPage)
	http.HandleFunc("/setup/", setupPage)
	http.HandleFunc("/games", games)
	http.HandleFunc("/addTeam", addTeam)
	http.HandleFunc("/tournament/", mainPage)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
