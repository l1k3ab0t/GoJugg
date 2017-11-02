package main

import (
	"github.com/l1k3ab0t/GoJugg/lib/FormatHTML"
	"github.com/l1k3ab0t/GoJugg/lib/GameEngine"
	"github.com/l1k3ab0t/GoJugg/lib/ReadConfig"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
	"os"
	"io"
)

type settings struct {
	webcfg     bool
	name       string
	list       string
	fields	int
	gameMode   int
	groupCont  int
	consoleLog bool
	port       int
	aTCount		int
}
type data struct {
	Name        template.HTML
	Table       template.HTML
	Bracket     template.HTML
	TextPhrase1 template.HTML
	TextPhrase2 template.HTML
	I           int
	I2          int
}

var s settings
var gg [][][]GameEngine.Game
var t []GameEngine.Team
var defaultData data
var setupDone bool
var templates = template.Must(template.ParseFiles("resources/control.html", "resources/tournament.html", "resources/games.html", "resources/setup.html", "resources/submitResult.html", "resources/rank.html", "resources/team.html"))

func loadCFG() (settings, []GameEngine.Team) {
	var s settings
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
		if ReadConfig.SplitConfig(v.Content)[0] == "AdvancingfromGroups" {
			i, err := strconv.Atoi(ReadConfig.SplitConfig(v.Content)[1])
			if err != nil {
				log.Fatalf("error setting Advancing Teams Count: %v", err)
			} else {
				s.aTCount = i
			}
		}
		if ReadConfig.SplitConfig(v.Content)[0] == "ConsoleLog" {
			if ReadConfig.SplitConfig(v.Content)[1] == "TRUE" {
				s.consoleLog = true
			} else {
				s.consoleLog = false
			}
		}
		if ReadConfig.SplitConfig(v.Content)[0] == "Port" {
			i, err := strconv.Atoi(ReadConfig.SplitConfig(v.Content)[1])
			if err != nil {
				log.Fatalf("error setting Port: %v", err)
			} else {
				s.port = i
			}

		}

	}

	log.Println(s)
	return s, ReadConfig.ReadTeamList(s.list)
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
	log.Printf("IP: " + r.RemoteAddr + " connected")
	if setupDone{
		log.Println("Setup Done Redirecting to /")
		http.Redirect(w, r, "/", http.StatusFound)
	}else {
		renderTemplate(w, "setup", defaultData)
	}
}
func cSettings(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	if r.Form["ChangeS"] != nil{
		if r.FormValue("GameMode")!=""{
			s.gameMode, _ = strconv.Atoi(r.FormValue("GameMode"))
			log.Println("Change GameModed to: ",s.gameMode)
		}
		if r.FormValue("NoF")!=""{
			s.fields, _ = strconv.Atoi(r.FormValue("NoF"))
			log.Println("Change Number of Fields to: ",s.fields)
		}
		if r.FormValue("NoG")!=""{
			s.groupCont, _ = strconv.Atoi(r.FormValue("NoG"))
			log.Println("Change Number of Groups to: ",s.groupCont)
		}
	}

	http.Redirect(w, r, "/setup/", http.StatusFound)
}

func uploadTList(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //get request method
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("FName")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	log.Println(w, "%v", handler.Header)
	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	s.list=handler.Filename
	log.Println("Loading Team List: ",s.list)
	exists:=false
	for _,v:= range ReadConfig.ReadTeamList(s.list){
		for _,v2:= range t{
			if v.Name==v2.Name{
				exists=true
			}
		}
		if !exists{
			t=append(t,v)
			http.HandleFunc("/"+v.Name, teamPage)
			log.Println("Adding Team: ",v)
		}else {
			log.Println("Team ",v," allready exsits, rejected")
			exists=false
		}
	}
	defaultData.Table = FormatHTML.FormatTeamLIst(t)
	http.Redirect(w, r, "/setup/", http.StatusFound)
}

func addTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	if r.Form["Add"] != nil && r.FormValue("Team")!="" {
		n:=FormatHTML.FormatTeamName(r.FormValue("Team"))
		t = append(t, GameEngine.Team{n, len(t)+1, len(t)+1, 1, nil, 0})
		http.HandleFunc("/"+n, teamPage)
		defaultData.Table = FormatHTML.FormatTeamLIst(t)
	}
	log.Println("Add Team: ", r.Form["Team"])

	http.Redirect(w, r, "/setup/", http.StatusFound)
}

func cTName(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	if r.Form["Change"] != nil && r.FormValue("TName")!="" {
		s.name=r.FormValue("TName")
		defaultData.Name=template.HTML(s.name)
	}
	log.Println("Change Turney Name: ", r.Form["TName"])

	http.Redirect(w, r, "/setup/", http.StatusFound)
}

func endSetup(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	if r.Form["end"] != nil{
		setupDone = true
		gg=nil
		var gq [][]GameEngine.Game
		if s.gameMode == 0 {
			t=GameEngine.ChangeTGroup(s.groupCont,t)
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
	}
	log.Println("Setup Done")
	http.Redirect(w, r, "/", http.StatusFound)
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
		if iP < s.groupCont {
			tdata.I = iP + 1
		} else {
			tdata.I = 0
		}
		tdata.TextPhrase1 = template.HTML("from Group " + strconv.Itoa(tdata.I))
		tdata.Table = FormatHTML.FormatRanking(GameEngine.SortByRank(gg[tdata.I], GameEngine.BuildGroups(s.groupCont, t)[tdata.I]))
	} else if r.Form["Custom"] != nil {

		log.Println("Custom Destination")
		iP, _ := strconv.Atoi(r.FormValue("Group"))
		if iP <= s.groupCont {
			tdata.I = iP
		} else {
			tdata.I = 0
		}
		tdata.TextPhrase1 = template.HTML("from Group " + strconv.Itoa(tdata.I))
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
						tdata.TextPhrase1 = template.HTML(v3.Opponent1.Name)
						tdata.TextPhrase2 = template.HTML(v3.Opponent2.Name)
						break
					} else if v3.Opponent2.ID == iD1 && v3.Opponent1.ID == iD2 {
						tdata.TextPhrase1 = template.HTML(v3.Opponent2.Name)
						tdata.TextPhrase2 = template.HTML(v3.Opponent1.Name)
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

func teamPage(w http.ResponseWriter, r *http.Request) {
	d := defaultData
	team:=GameEngine.TeamByName(string(FormatHTML.FormatURI(r.RequestURI)),t)
	d.Name = FormatHTML.FormatURI(r.RequestURI)
	d.I=GameEngine.TeamRank(team.Name,GameEngine.SortByRankInTourney(gg, GameEngine.BuildGroups(s.groupCont, t))).Rank
	d.I2=GameEngine.TeamRank(team.Name,GameEngine.SortByRank(gg[team.Group],GameEngine.BuildGroups(s.groupCont,t)[team.Group])).Rank
	d.TextPhrase1=FormatHTML.FormatBracket(GameEngine.TeamGames(team.Name,gg[team.Group]))
	renderTemplate(w, "team", d)
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func buildEliminationGames(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/setup/", http.StatusFound)
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
	/*
	lfName := time.Now().Format(time.RFC3339) + ".log"
	f, err := os.OpenFile(lfName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	if s.consoleLog {
		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
	} else {
		log.SetOutput(f)
	}
	log.Println("This is a test log entry")
	*/

	setupDone = !s.webcfg
	var gq [][]GameEngine.Game
	if s.gameMode == 0 {
		t=GameEngine.ChangeTGroup(s.groupCont,t)
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
	log.Println(time.Now().Format(time.RFC3339))
	defaultData.Name = template.HTML(s.name)
	defaultData.Table = FormatHTML.FormatTeamLIst(t)
	defaultData.Bracket = FormatHTML.FormatBracket(gg[0][0])
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.HandleFunc("/", defaultConnection)
	http.HandleFunc("/control/", controlPage)
	http.HandleFunc("/setup/", setupPage)
	http.HandleFunc("/rank", rank)
	http.HandleFunc("/games", games)
	http.HandleFunc("/cSettings", cSettings)
	http.HandleFunc("/addTeam", addTeam)
	http.HandleFunc("/cTName", cTName)
	http.HandleFunc("/uploadTList", uploadTList)
	http.HandleFunc("/endSetup", endSetup)
	http.HandleFunc("/submitResult", submitResult)
	http.HandleFunc("/receiveResult", receiveResult)
	http.HandleFunc("/tournament/", mainPage)
	http.HandleFunc("/bEGames/", buildEliminationGames)
	for _, v := range t {
		http.HandleFunc("/"+v.Name, teamPage)
	}
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(s.port), nil))
}
