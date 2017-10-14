package FormatHTML

import (
	"github.com/l1k3ab0t/GoJugg/lib/GameEngine"
	"strconv"
	"html/template"
)

func FormatTeamLIst(teams []GameEngine.Team) template.HTML {
	tmpl :=template.HTML("<table>")
	tmpl = tmpl +"<tr> <th> Team ID </th> <th> Team Name </th> </tr>"
	for _,v:=range teams{
		tmpl = tmpl + template.HTML("\n		<tr> <td>" + strconv.Itoa(v.ID) + "</td>" + "<td>" + v.Name + "</td> </tr> " )
	}
	tmpl = tmpl +"\n	 </table>"
	return tmpl
}

func FormatBracket(games []GameEngine.Game) template.HTML {
	tmpl :=template.HTML("<table>")
	tmpl = tmpl +"<tr> <th> Team 1 ID </th> <th> Team 1 </th> <th></th> <th> Team 2 ID </th> <th> Team 2 </th> <th>Result </th> </tr>"
	for _,v:=range games{
		tmpl = tmpl + template.HTML("\n		<tr> <td>" + strconv.Itoa(v.Opponent1.ID) + "</td>" + "<td>" + v.Opponent1.Name + "</td> <td>	vs		</td> " )
		tmpl = tmpl + template.HTML(" <td>" + strconv.Itoa(v.Opponent2.ID) + "</td>" + "<td>" + v.Opponent2.Name + "</td> <td>"+strconv.Itoa(v.Result.Team1Juggs)+" : "+strconv.Itoa(v.Result.Team2Juggs) + "</td>" )
		tmpl = tmpl +"<td>" + formatSubmitButton(v) + "</td></tr>"
	}
	tmpl = tmpl +"\n	 </table>"
	return tmpl
}

func formatSubmitButton(game GameEngine.Game) template.HTML{
	tmpl:=template.HTML("\n<form action=\"/submitResult\" method=\"post\">")
	tmpl=tmpl+template.HTML("\n <input hidden=\"hidden\" type=\"text\" value=\""+ strconv.Itoa(game.Opponent1.ID)+"\" name=\"T1ID\">")
	tmpl=tmpl+template.HTML("\n <input hidden=\"hidden\" type=\"text\" value=\""+strconv.Itoa(game.Opponent2.ID)+"\" name=\"T2ID\">")
	tmpl=tmpl+template.HTML("\n  <input type=\"submit\"  name=\"Custom\" value=\"Submit Result\"> </form>")
	return tmpl
}