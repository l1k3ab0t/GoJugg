package main

import (
	"github.com/gorilla/sessions"
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

type data struct {
	Name        template.HTML
	Table       template.HTML
	Bracket     template.HTML
	TextPhrase1 template.HTML
	TextPhrase2 template.HTML
	I           int
	I2          int
	I3          int
	I4          int
	Options     []option
}

type option struct {
	TextPhrase1 template.HTML
	B           bool
	Value       template.HTML
}

var t = Objects.Tournament{}
var running bool
var defaultData data
var setupDone bool
var templates = template.Must(template.ParseGlob("resources/*"))
var store = sessions.NewCookieStore([]byte("change-me-pls"))
var sessionName string
var pwM Objects.PWHandler

func startAutoSaving(tme time.Duration) {
	running = true
	Ticker := time.NewTicker(tme)
	for tik := range Ticker.C {
		if running {
			log.Println("Tick at", tik)
			autoSave("save.xml")
		} else {
			break
		}
	}
	Ticker.Stop()
	log.Println("Ticker Stoped")
}
func stopAutoSaving() {
	running = false
	log.Println("Stoping Autosaving")
}

func autoSave(filename string) {
	log.Println("AutoSaving:", filename)
	Objects.CreateSave(filename, t)
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
	//stopAutoSaving()
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
			t.S.SetFields(noF - 1)
		}
		if r.FormValue("NoG") != "" {
			noG, _ := strconv.Atoi(r.FormValue("NoG"))
			t.S.SetGroupCount(noG - 1)
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
	fname := handler.Filename
	log.Println("Loading Team List: ", fname)
	tL := Objects.ReadTeamList(fname)
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

func uploadSave(w http.ResponseWriter, r *http.Request) {
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
	fname := handler.Filename
	t = Objects.LoadTournament(Objects.LoadFile(fname))
	setupDone = true
	http.Redirect(w, r, "/", http.StatusFound)
}

func addTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	if r.Form["Add"] != nil && r.FormValue("Team") != "" {
		n := FormatHTML.FormatTeamName(r.FormValue("Team"))
		t.SetupTeamList = append(t.SetupTeamList, Objects.Team{n, len(t.SetupTeamList) + 1, len(t.SetupTeamList) + 1, 1, nil, 0}) //Setter Requiert
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
		t.S.SetName(r.FormValue("TName"))
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
		t.FinishSetup() //load extra teamlist
		sessionName = FormatHTML.FormatTeamName(t.S.Name()) + "session"
		Objects.CreateSave("test.xml", t)
		var pwL []Objects.Password
		for _, v := range t.TeamList() {
			pwL = append(pwL, Objects.Password{v.Name, "1234"})
		}
		pwM = Objects.CreateDB("test.db", pwL)
		log.Println(pwM.Check(Objects.Password{"Rigor-Mortis", "1234"}))
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
	var id int
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	log.Println("Group: ", r.FormValue("Group"))
	log.Println("Round: ", r.FormValue("Round"))
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session.Values["ID"] != nil {
		id = session.Values["ID"].(int)
	} else {
		id = 24000
	}
	if r.Form["CGroup"] != nil {
		log.Println("Group Changed")
		tdata.I2 = 0

		iP, _ := strconv.Atoi(r.FormValue("Group"))
		if iP < t.S.GroupCount() {
			tdata.I = iP + 1
		} else {
			tdata.I = 0
		}
		tdata.Bracket = FormatHTML.FormatBracket(t.Games()[tdata.I][tdata.I2], id)
	} else if r.Form["CRound"] != nil {
		log.Println("Round Changed")
		iR, _ := strconv.Atoi(r.FormValue("Round"))
		if iR < t.S.RoundCount() {
			tdata.I2 = iR + 1
		} else {
			tdata.I2 = 0
		}
		iP, _ := strconv.Atoi(r.FormValue("Group"))
		tdata.I = iP
		tdata.Bracket = FormatHTML.FormatBracket(t.Games()[tdata.I][tdata.I2], id)
	} else if r.Form["Custom"] != nil {

		log.Println("Custom Destination")
		iP, _ := strconv.Atoi(r.FormValue("Group"))
		if iP <= t.S.GroupCount() {
			tdata.I = iP
		} else {
			tdata.I = 0
		}

		iR, _ := strconv.Atoi(r.FormValue("Round"))
		if iR < t.S.RoundCount() {
			tdata.I2 = iR
		} else {
			tdata.I2 = 0
		}
		tdata.Bracket = FormatHTML.FormatBracket(t.Games()[tdata.I][tdata.I2], id)
	} else {

		if session.Values["Group"] != nil {
			tdata.I = session.Values["Group"].(int)
			tdata.Bracket = FormatHTML.FormatBracket(t.Games()[session.Values["Group"].(int)][0], id)
		} else {
			tdata.Bracket = FormatHTML.FormatBracket(t.Games()[0][0], id)
		}
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
		session, err := store.Get(r, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if session.Values["ID"] == iD1 || session.Values["ID"] == iD2 || session.Values["Admin"] == true {
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
			tdata.TextPhrase1 = template.HTML("not logged in!")
			renderTemplate(w, "error", tdata)
		}
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
							t.SetGameResult(i1, i2, i3, Objects.Result{t1j, t2j})
							log.Println("1")
							log.Println(t.Games()[i1][i2][i3])
							t.UpdateSingleGRanked(i1)
							log.Println(t.GroupRanked(i1))
							break
						} else if v3.Opponent2.Name == r.FormValue("Team1") && v3.Opponent1.Name == r.FormValue("Team2") {
							t.SetGameResult(i1, i2, i3, Objects.Result{t2j, t1j})
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
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	log.Println("Password: ", r.FormValue("pw"))
	d := defaultData
	team := t.TeamByName(string(FormatHTML.FormatURI(r.RequestURI)))
	d.Name = FormatHTML.FormatURI(r.RequestURI)
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Form["Login"] != nil {
		if pwM.Check(Objects.Password{team.Name, r.FormValue("pw")}) {
			session.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   86400 * 5,
				HttpOnly: true,
			}
			// Set some session values.
			session.Values["Admin"] = false
			session.Values["Round"] = 0
			session.Values["Group"] = team.Group
			session.Values["ID"] = team.ID
			session.Values["loggedIn"] = true
			session.ID = team.Name
			// Save it before we write to the response/return from the handler.
			session.Save(r, w)
			log.Println(session.ID, " logged in")
			log.Println(r.RemoteAddr)
		} else {
			log.Println("Authentication failed, wrong Password. User: ", team.Name)
		}

	}
	d.I = t.TeamRank(team).Rank
	d.I2 = t.TeamGRank(team).Rank
	if session.Values["ID"] != nil {
		d.TextPhrase1 = FormatHTML.FormatBracket(t.TeamGames(team), session.Values["ID"].(int))
	} else {
		d.TextPhrase1 = FormatHTML.FormatBracket(t.TeamGames(team), 23000)
	}

	renderTemplate(w, "team", d)
	log.Printf("IP: " + r.RemoteAddr + " connected")
	log.Println(session.Values["ID"])
}

func buildEliminationGames(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/setup/", http.StatusFound)
}

func saveTournament(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	log.Println(" Save File Name: ", r.FormValue("FileName"))
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session.Values["Admin"] == true && r.Form["Save"] != nil {
		log.Println("Saving:", r.FormValue("FileName"))
		Objects.CreateSave(r.FormValue("FileName"), t)
		http.Redirect(w, r, "/admin", http.StatusFound)
	} else {
		errorData := defaultData
		errorData.TextPhrase1 = template.HTML("No Rights/ No fileName")
		renderTemplate(w, "error", errorData)
		log.Printf("IP: " + r.RemoteAddr + " connected")
	}
}

func admin(w http.ResponseWriter, r *http.Request) {

	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	log.Println("User: ", r.FormValue("User"))
	log.Println("Password: ", r.FormValue("Password")) //Remove Later
	log.Println(r.Form["User"])
	session, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session.Values["Admin"] == true {
		adminData := defaultData
		if r.Form["Select"] != nil {

			log.Println("Selected Team: ", r.FormValue("TName"))
			adminData.Bracket = FormatHTML.FormatBracket(t.TeamGames(t.TeamByName(r.FormValue("TName"))), 0)

		} else {
			adminData.Bracket = FormatHTML.FormatBracket(t.Games()[0][0], 0)
		}
		for _, v := range t.TeamList() {
			if r.FormValue("TName") == v.Name {
				adminData.Options = append(adminData.Options, option{template.HTML(v.Name), true, template.HTML("")})
			} else {
				adminData.Options = append(adminData.Options, option{template.HTML(v.Name), false, template.HTML("")})

			}
		}

		renderTemplate(w, "admin", adminData)
	} else if r.Form["Login"] != nil && r.FormValue("User") == "Admin" {
		session, err := store.Get(r, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 5,
			HttpOnly: true,
		}

		session.Values["ID"] = 0
		session.Values["Admin"] = true
		session.Values["loggedIn"] = true
		session.ID = "Admin"
		session.Save(r, w)
		log.Println(session.ID, " logged in")

		log.Println(r.RemoteAddr)
		http.Redirect(w, r, "/admin", http.StatusFound)
	} else {
		adminData := defaultData
		adminData.TextPhrase1 = template.HTML("admin")
		adminData.Options = append(adminData.Options, option{template.HTML("Admin"), true, template.HTML("")})

		renderTemplate(w, "login", adminData)
	}
	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func login(w http.ResponseWriter, r *http.Request) {
	log.Println("method:", r.Method) //get request method
	r.ParseForm()
	log.Println("User: ", r.FormValue("User"))
	log.Println("Password: ", r.FormValue("Password")) //Remove Later
	log.Println(r.Form["User"])
	if r.Form["User"] == nil {
		loginData := defaultData
		loginData.TextPhrase1 = template.HTML("login")
		for _, v := range t.TeamList() {
			loginData.Options = append(loginData.Options, option{template.HTML(v.Name), false, template.HTML("")})
		}
		renderTemplate(w, "login", loginData)
	} else {
		team := t.TeamByName(r.FormValue("User"))
		if pwM.Check(Objects.Password{team.Name, r.FormValue("pw")}) {
			session, err := store.Get(r, sessionName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			session.Options = &sessions.Options{
				Path:     "/",
				MaxAge:   86400 * 5,
				HttpOnly: true,
			}

			// Set some session values.
			session.Values["Admin"] = false
			session.Values["Round"] = 0
			session.Values["Group"] = team.Group
			session.Values["ID"] = team.ID
			session.Values["loggedIn"] = true
			session.ID = team.Name
			// Save it before we write to the response/return from the handler.
			session.Save(r, w)
			log.Println(session.ID, " logged in")
			log.Println(r.RemoteAddr)
			http.Redirect(w, r, "/"+r.FormValue("User"), http.StatusFound)
		} else {
			log.Println("Authentication failed, wrong Password. User: ", team.Name)
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}

	log.Printf("IP: " + r.RemoteAddr + " connected")
}

func renderTemplate(w http.ResponseWriter, tmpl string, t data) {
	err := templates.ExecuteTemplate(w, tmpl+".html", t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	//t = Objects.NewTournament()
	t = Objects.StartSetup()
	if t.S.Webcfg() == false {
		setupDone = true
		t.FinishSetup()
		sessionName = FormatHTML.FormatTeamName(t.S.Name()) + "session"
	} else {
		setupDone = false
	}
	//go startAutoSaving(time.Second * 250)

	/*
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
	*/
	log.Println("This is a test log entry")

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
	http.HandleFunc("/uploadSave", uploadSave)
	http.HandleFunc("/endSetup", endSetup)
	http.HandleFunc("/submitResult", submitResult)
	http.HandleFunc("/receiveResult", receiveResult)
	http.HandleFunc("/tournament/", mainPage)
	http.HandleFunc("/bEGames/", buildEliminationGames)
	http.HandleFunc("/save", saveTournament)
	http.HandleFunc("/admin", admin)
	http.HandleFunc("/login", login)
	for _, v := range t.SetupTeamList {
		http.HandleFunc("/"+v.Name, teamPage)
	}
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(t.S.Port()), nil))
}
