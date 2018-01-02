package FormatHTML

import (
	"github.com/l1k3ab0t/GoJugg/lib/GameEngine"
	"html/template"
	"strconv"
)

func FormatTeamLIst(teams []GameEngine.Team) template.HTML {
	tmpl := template.HTML("<table>")
	tmpl = tmpl + "<tr> <th> Team ID </th> <th> Team Name </th> </tr>"
	for _, v := range teams {
		tmpl = tmpl + template.HTML("\n		<tr> <td>"+strconv.Itoa(v.ID)+"</td>"+"<td>"+string(FormatTeamLink(v.Name))+"</td> </tr> ")
	}
	tmpl = tmpl + "\n	 </table>"
	return tmpl
}

func FormatBracket(games []GameEngine.Game) template.HTML {
	tmpl := template.HTML("<table>")
	tmpl = tmpl + "<tr> <th> Team 1 ID </th> <th> Team 1 </th> <th></th> <th> Team 2 ID </th> <th> Team 2 </th> <th>Result </th> </tr>"
	for _, v := range games {
		tmpl = tmpl + template.HTML("\n		<tr> <td>"+strconv.Itoa(v.Opponent1.ID)+"</td>"+"<td>"+string(FormatTeamLink(v.Opponent1.Name))+"</td> <td>	vs		</td> ")
		tmpl = tmpl + template.HTML(" <td>"+strconv.Itoa(v.Opponent2.ID)+"</td>"+"<td>"+string(FormatTeamLink(v.Opponent2.Name))+"</td> <td>"+strconv.Itoa(v.Result.Team1Juggs)+" : "+strconv.Itoa(v.Result.Team2Juggs)+"</td>")
		tmpl = tmpl + "<td>" + formatSubmitButton(v) + "</td></tr>"
	}
	tmpl = tmpl + "\n	 </table>"
	return tmpl
}

func formatSubmitButton(game GameEngine.Game) template.HTML {
	tmpl := template.HTML("\n<form action=\"/submitResult\" method=\"post\">")
	tmpl = tmpl + template.HTML("\n <input hidden=\"hidden\" type=\"text\" value=\""+strconv.Itoa(game.Opponent1.ID)+"\" name=\"T1ID\">")
	tmpl = tmpl + template.HTML("\n <input hidden=\"hidden\" type=\"text\" value=\""+strconv.Itoa(game.Opponent2.ID)+"\" name=\"T2ID\">")
	tmpl = tmpl + template.HTML("\n  <input type=\"submit\"  name=\"Custom\" value=\"Submit Result\"> </form>")
	return tmpl
}

func FormatRanking(r []GameEngine.Rank) template.HTML {
	tmpl := template.HTML("<table>")
	tmpl = tmpl + "<tr> <th> Rank </th> <th> Team </th> <th> Score </th> </tr>"
	for _, v := range r {
		tmpl = tmpl + template.HTML("\n		<tr> <td>"+strconv.Itoa(v.Rank)+"</td>"+"<td>"+string(FormatTeamLink(v.TName))+"</td> <td>"+strconv.Itoa(v.Result.Team1Juggs)+"-"+strconv.Itoa(v.Result.Team2Juggs)+"</td></tr> ")
	}
	tmpl = tmpl + "\n	 </table>"
	return tmpl

}

func FormatTeamName(name string) string {
	var rname string
	for _, char := range name {
		if char == 32 { //32 == " "
			rname = rname + "-"
		}else if char == 228{
			rname = rname + "ae"
		}else if char == 196{
			rname = rname + "Ae"
		}else if char == 246{
			rname = rname + "oe"
		}else if char == 214{
			rname = rname + "Oe"
		}else if char == 252{
			rname = rname + "ue"
		}else if char == 220{
			rname = rname + "Ue"
		} else {
			rname = rname + string(char)
		}
	}
	return rname
}

func FormatURI(uri string) template.HTML {
	var str string
	for _, v := range uri {
		if v != 47 {
			str = str + string(v)
		}
	}
	return template.HTML(str)
}

func FormatTeamLink(uri string) template.HTML {
	return template.HTML("<a href=\"/" + uri + "\">" + uri + "</a>")
}

func FormatGSave(g [][][]GameEngine.Game) string{
	var str string
	for _,v:=range g{
		str=str+"<ggrp>"
		for _,v2:=range v{
			str=str+"<round>"
			for _,v3:=range v2{
				str=str+"<gfield>"
				str=str+"<op1>"
				str=str+strconv.Itoa(v3.Opponent1.ID)
				str=str+"</op1>"
				str=str+"<op2>"
				str=str+strconv.Itoa(v3.Opponent2.ID)
				str=str+"</op2>"
				str=str+"<res>"
				str=str+"<t1j>"
				str=str+strconv.Itoa(v3.Result.Team1Juggs)
				str=str+"</t1j>"
				str=str+"</t2j>"
				str=str+strconv.Itoa(v3.Result.Team2Juggs)
				str=str+"</t2j>"
				str=str+"</res>"
				str=str+"</gfield>"
			}
			str=str+"</round>"
		}
		str=str+"</ggrp>"
	}
	return str
}

