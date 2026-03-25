// Package main presents how to implement large, procedurally-generated table.
package main

import (
	"fmt"
	"image/color"
	"strconv"

	"giu_ui/giudoc"

	// "github.com/AllenDang/giu"
	g "github.com/AllenDang/giu"
)

// giu State variables
var sashPos float32 = 500
var selectedFunction int = -1
var details string = "(⌐■_■) click on function calls"

func buildRows(messages []message) []*g.TableRowWidget {
	rows := make([]*g.TableRowWidget, len(messages))
	if len(messages) == 0 {
		return []*g.TableRowWidget{
			g.TableRow(g.Label("	( ╥﹏╥) 	Nothing yet!!")),
		}
	} else {
		for i := range rows {
			if messages[i].Name == "" {
				messages[i].Name = getName(messages[i].Id)
			}
			rows[i] = g.TableRow(
				g.Label(strconv.Itoa(i+1)),
				g.Selectable(messages[i].Id).
					Selected(selectedFunction == i).
					OnClick(func() {
						selectedFunction = i
						go getDesc(messages[i].Id, &details)
						fmt.Println("Selected:", messages[i].Name)
					}),
				g.Selectable(messages[i].Name).
					Selected(selectedFunction == i).
					OnClick(func() {
						selectedFunction = i
						go getDesc(messages[i].Id, &details)
						fmt.Println("Selected:", messages[i].Name)
					}),
			)
			if selectedFunction == i {
				rows[i].BgColor(&color.RGBA{80, 120, 200, 100})
			}

		}
	}

	return rows
}

func loop(msgchan chan message, messages []message) []message {

	for {
		select {
		case msg := <-msgchan:
			messages = append(messages, msg)
		default:
			goto render
		}
	}
render:

	g.SingleWindow().Layout(
		g.SplitLayout(
			g.DirectionVertical,
			&sashPos,
			g.Layout{
				g.Style().
					SetFontSize(15).
					To(
						g.Label("Function List"),
					),
				g.Separator(),
				g.Table().Freeze(0, 1).FastMode(true).Columns(g.TableColumn("S/N."), g.TableColumn("Id"), g.TableColumn("Function Name")).Rows(buildRows(messages)...),
			},
			g.Layout{
				g.Style().
					SetFontSize(15).
					To(
						g.Label("Description"),
					),
				g.Separator(),
				giudoc.FromHTML(details),
			},
		),
	)

	return messages
}
