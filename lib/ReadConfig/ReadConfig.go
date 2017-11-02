package ReadConfig

import (
	"bufio"
	"log"
	"os"
	"strings"
	"github.com/l1k3ab0t/GoJugg/lib/GameEngine"
	"strconv"
	"github.com/l1k3ab0t/GoJugg/lib/FormatHTML"
)

type Line struct {
	Linenumber int
	Content    string
}

func ReadFile(path string) []Line {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	var str []Line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" && strings.Index(scanner.Text(), "#") != 0 {
			i++
			str = append(str, Line{i, scanner.Text()})
		}
	}
	file.Close()
	return str
}

func SplitConfig(config string) []string {
	splits := strings.Split(config, "=")
	return splits

}

func ReadTeamList(f string) []GameEngine.Team{
	teamList := ReadFile(f)
	var t []GameEngine.Team
	for i := range teamList {
		log.Println(strconv.Itoa(teamList[i].Linenumber) + " " + teamList[i].Content)
		t = append(t, GameEngine.Team{FormatHTML.FormatTeamName(teamList[i].Content), teamList[i].Linenumber, teamList[i].Linenumber, 1, nil, 0})
	}
	return t
}
