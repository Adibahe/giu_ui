// Package main presents how to implement large, procedurally-generated table.
package main

import (
	"fmt"
	"image/color"

	"github.com/AllenDang/giu"
	g "github.com/AllenDang/giu"
)

// giu State variables
var sashPos float32 = 500
var selectedFunction int = -1
var details string = "Hello World!!"

func buildRows(messages []message) []*g.TableRowWidget {
	rows := make([]*g.TableRowWidget, len(messages))
	for i := range rows {
		rows[i] = g.TableRow(
			g.Selectable(messages[i].Id).
				Selected(selectedFunction == i).
				OnClick(func() {
					selectedFunction = i
					details = "Selected:" + messages[i].Id
					fmt.Println("Selected:", messages[i].Id)
				}),
			g.Selectable(messages[i].Name).
				Selected(selectedFunction == i).
				OnClick(func() {
					selectedFunction = i
					details = "Selected:" + messages[i].Id
					fmt.Println("Selected:", messages[i].Id)
				}),
		)
		if selectedFunction == i {
			rows[i].BgColor(&color.RGBA{80, 120, 200, 100})
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
				g.Table().Freeze(0, 1).FastMode(true).Columns(g.TableColumn("Id"), g.TableColumn("Function Name")).Rows(buildRows(messages)...),
			},
			g.Layout{
				g.Style().
					SetFontSize(15).
					To(
						g.Label("Description"),
					),
				g.Separator(),
				g.InputTextMultiline(&details).
					Size(-1, 200).
					Flags(giu.InputTextFlagsReadOnly).ID("Desc"),
			},
		),
	)

	return messages
}
