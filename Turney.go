package main

import (
	"github.com/l1k3ab0t/GoJugg/lib/FormatHTML"
	"github.com/l1k3ab0t/GoJugg/lib/Objects"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

)

type save struct {
	t Objects.Tournament
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
var t= Objects.Tournament{}
var running bool
var defaultData data
var setupDone bool
var templates = template.Must(template.ParseFiles("resources/control.html", "resources/tournament.html", "resources/games.html", "resources/setup.html", "resources/submitResult.html", "resources/rank.html", "resources/team.html"))


func startAutosaving(tme time.Duration) {
	running = true
	Ticker := time.NewTicker(tme)
	for tik := range Ticker.C {
		if running {
			log.Println("Tick at", tik)
			autosave("save.xml", save{t})
		} else {
			break
		}
	}
	Ticker.Stop()
	log.Println("Ticker Stoped")
}
func stopAutosaving() {
	running = false
	log.Println("Stoping Autosaving")
}

func autosave(filename string, content save) {
	log.Println("AutoSaving:", filename)
	f, err := os.OpenFile("saves/"+time.Now().Format(time.RFC3339)+".as", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	f.WriteString(FormatHTML.FormatGSave(t.Games()))
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
	stopAutosaving()
	log.Printf("IP: " + r.RemoteAddr + " connected")
	if setupDone {
		log.Println("Setup Done Redirecting to /")
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		renderTemplate(w, "setup", defaultData)
	}
}
func cSettings(w http.ResponseWriter, r *http.Request) { //Setter Requert
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	if r.Form["ChangeS"] != nil {
		if r.FormValue("GameMode") != "" {
			 gm, _ := strconv.Atoi(r.FormValue("GameMode"))
			 t.S.SetGameMode(gm)
			log.Println("Change GameModed to: ", t.S.GameMode())
		}
		if r.FormValue("NoF") != "" {
			noF, _ := strconv.Atoi(r.FormValue("NoF"))
			log.Println("Change Number of Fields to: ", t.S.Fields())
			t.S.SetFields(noF-1)
		}
		if r.FormValue("NoG") != "" {
			noG, _ := strconv.Atoi(r.FormValue("NoG"))
			t.S.SetGroupCount(noG-1)
			log.Println("Change Number of Groups to: ", t.S.GroupCount())

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
	fname:= handler.Filename
	log.Println("Loading Team List: ", fname)
	tL:=Objects.ReadTeamList(fname)
	exists := false
	for _, v := range tL {
		for _, v2 := range t.SetupTeamList {
			if v.Name == v2.Name {
				exists = true
			}
		}
		if !exists {
			t.SetupTeamList = append(t.SetupTeamList, v)
			http.HandleFunc("/"+v.Name, teamPage)
			log.Println("Adding Team: ", v)
		} else {
			log.Println("Team ", v, " allready exsits, rejected")
			exists = false
		}
	}
	defaultData.Table = FormatHTML.FormatTeamLIst(t.SetupTeamList)
	http.Redirect(w, r, "/setup/", http.StatusFound)
}

func addTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	if r.Form["Add"] != nil && r.FormValue("Team") != "" {
		n := FormatHTML.FormatTeamName(r.FormValue("Team"))
		t.SetupTeamList= append(t.SetupTeamList, Objects.Team{n, len(t.SetupTeamList) + 1, len(t.SetupTeamList) + 1, 1, nil, 0}) //Setter Requiert
		http.HandleFunc("/"+n, teamPage)
		defaultData.Table = FormatHTML.FormatTeamLIst(t.SetupTeamList)
	}
	log.Println("Add Team: ", r.Form["Team"])

	http.Redirect(w, r, "/setup/", http.StatusFound)
}

func cTName(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	if r.Form["Change"] != nil && r.FormValue("TName") != "" {
		t.S.SetName(r.FormValue("TName"))          //???????????????????
		defaultData.Name = template.HTML(t.S.Name())
	}
	log.Println("Change Turney Name: ", r.Form["TName"])

	http.Redirect(w, r, "/setup/", http.StatusFound)
}

func endSetup(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	if r.Form["end"] != nil {
		setupDone = true
		t.FinishSetup()//load extra teamlist missing
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
		if iP < t.S.GroupCount() {
			tdata.I = iP + 1
		} else {
			tdata.I = 0
		}
		tdata.TextPhrase1 = template.HTML("from Group " + strconv.Itoa(tdata.I))
		tdata.Table = FormatHTML.FormatRanking(t.GroupRanked(tdata.I))
	} else if r.Form["Custom"] != nil {

		log.Println("Custom Destination")
		iP, _ := strconv.Atoi(r.FormValue("Group"))
		if iP <= t.S.GroupCount() {
			tdata.I = iP
		} else {
			tdata.I = 0
		}
		tdata.TextPhrase1 = template.HTML("from Group " + strconv.Itoa(tdata.I))
		tdata.Table = FormatHTML.FormatRanking(t.GroupRanked(tdata.I))

	} else {
		tdata.I = 0
		tdata.Table = FormatHTML.FormatRanking(t.Ranked())
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
		if iP < t.S.GroupCount() {
			tdata.I = iP + 1
		} else {
			tdata.I = 0
		}
		tdata.Bracket = FormatHTML.FormatBracket(t.Games()[tdata.I][tdata.I2])
	} else if r.Form["CRound"] != nil {
		log.Println("Round Changed")
		iR, _ := strconv.Atoi(r.FormValue("Round"))
		if iR < t.S.RoundCount() {
			tdata.I2 = iR + 1
		} else {
			tdata.I2 = 0
		}
		iP, _ := strconv.Atoi(r.FormValue("Page"))
		tdata.I = iP
		tdata.Bracket = FormatHTML.FormatBracket(t.Games()[tdata.I][tdata.I2])
	} else if r.Form["Custom"] != nil {

		log.Println("Custom Destination")
		iP, _ := strconv.Atoi(r.FormValue("Page"))
		if iP <= t.S.GroupCount() {
			tdata.I = iP
		} else {
			tdata.I = 0
		}

		iR, _ := strconv.Atoi(r.FormValue("Round"))
		if iR < t.S.GroupCount() {
			tdata.I2 = iR
		} else {
			tdata.I2 = 0
		}
		tdata.Bracket = FormatHTML.FormatBracket(t.Games()[tdata.I][tdata.I2])
	} else {
		tdata.Bracket = FormatHTML.FormatBracket(t.Games()[0][0])
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
		for _, v := range t.Games() {
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
			for i1, v1 := range t.Games() {
				for i2, v2 := range v1 {
					for i3, v3 := range v2 {
						if v3.Opponent1.Name == r.FormValue("Team1") && v3.Opponent2.Name == r.FormValue("Team2") { //Setter Requiert
							t.SetGameResult(i1,i2,i3,Objects.Result{t1j,t2j})
							log.Println("1")
							log.Println(t.Games()[i1][i2][i3])
							t.UpdateSingleGRanked(i1)
							log.Println(t.GroupRanked(i1))
							break
						} else if v3.Opponent2.Name == r.FormValue("Team1") && v3.Opponent1.Name == r.FormValue("Team2") {
							t.SetGameResult(i1,i2,i3,Objects.Result{t2j,t1j})
							log.Println("2")
							log.Println(t.Games()[i1][i2][i3])
							t.UpdateSingleGRanked(i1)
							log.Println(t.GroupRanked(i1))
							break
						}
					}
				}
			}
		}

	}
	t.UpdateRanked()
	http.Redirect(w, r, "/games", http.StatusFound)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "tournament", defaultData)
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func teamPage(w http.ResponseWriter, r *http.Request) {
	d := defaultData
	team := t.TeamByName(string(FormatHTML.FormatURI(r.RequestURI)))
	d.Name = FormatHTML.FormatURI(r.RequestURI)
	d.I = t.TeamRank(team).Rank
	d.I2 = t.TeamGRank(team).Rank
	d.TextPhrase1 = FormatHTML.FormatBracket(t.TeamGames(team))
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
	//t = Objects.NewTournament()
	t=Objects.StartSetup()
	go startAutosaving(time.Second * 25)

	setupDone = false

	lfName := "log/" + time.Now().Format(time.RFC3339) + ".log"
	f, err := os.OpenFile(lfName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	if t.S.ConsoleLog() {
		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
	} else {
		log.SetOutput(f)
	}
	log.Println("This is a test log entry")

	setupDone = !t.S.Webcfg()



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
	defaultData.Name = template.HTML(t.S.Name())
	defaultData.Table = FormatHTML.FormatTeamLIst(t.SetupTeamList)
	//defaultData.Bracket = FormatHTML.FormatBracket(t.Games()[0][0])
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
	for _, v := range t.SetupTeamList {
		http.HandleFunc("/"+v.Name, teamPage)
	}
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(t.S.Port()), nil))
}
