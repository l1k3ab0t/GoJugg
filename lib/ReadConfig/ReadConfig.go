package ReadConfig

import (
	"bufio"
	"log"
	"os"
	"strings"
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
