package main

import (
	"github.com/l1k3ab0t/GoJugg/lib/GameEngine"
	"github.com/l1k3ab0t/GoJugg/lib/ReadConfig"
	"log"
	"strconv"
)

type settings struct {
	list      string
	gameMode  int
	groupCont int
}

func loadCFG() (settings, []GameEngine.Team) {
	var s settings
	var t []GameEngine.Team
	cfg := ReadConfig.ReadFile("config.cfg")
	for _, v := range cfg {
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
		t = append(t, GameEngine.Team{teamList[i].Content, teamList[i].Linenumber, teamList[i].Linenumber, 1})
	}

	return s, t
}

func main() {
	s, t := loadCFG()
	if s.gameMode == 0 {
		t = GameEngine.BuildGroups(s.groupCont, t)
	}
	log.Println("_________________________")
	for _, v := range t {
		log.Println(v.Name + " " + strconv.Itoa(v.Group))
	}
}
