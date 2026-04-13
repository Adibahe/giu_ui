package main

import (
	"fmt"
	"giu_ui/rad_api"
	"log"
	"os/exec"

	"github.com/sqweek/dialog"

	webview "github.com/webview/webview_go"
)

func addTableRow(w webview.WebView, id string, functionName string) {
	w.Dispatch(func() {
		js := fmt.Sprintf(`window.addRow(%q, %q)`, id, functionName)
		w.Eval(js)
		log.Printf("sent: %s, %s \n", id, functionName)
	})
}

func uiUpdater(w webview.WebView, msgChan <-chan message, messages *[]message) {
	go func() {
		for msg := range msgChan {

			msg.Name = getName(msg.Id)

			*messages = append(*messages, msg)
			addTableRow(w, msg.Id, msg.Name)
		}
	}()
}

func onPageReload(messages []message) {
	fmt.Println("Webpage loaded/reloaded")
}

func giveToJs(w webview.WebView) {
	err := w.Bind("getDescFromGo", func(id string) string {
		var details string
		getDesc(id, &details)
		return details
	})
	if err != nil {
		log.Fatal(err)
	}
}

func openExternalLink(url string) string {
	// log.Println("Opening ", url)
	err := exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	if err != nil {
		return err.Error()
	}
	return "ok"
}

var r_api rad_api.RadIpcState

func DebugCommand(action string) {
	var err error

	switch action {

	case "step_into":
		err = r_api.SendCommand(rad_api.CMD_STEP_INTO, "")

	case "step_over":
		err = r_api.SendCommand(rad_api.CMD_STEP_OVER, "")

	case "run":
		err = r_api.SendCommand(rad_api.CMD_RUN, "")

	case "stop":
		err = r_api.SendCommand(rad_api.CMD_HALT, "")

	default:
		fmt.Println("Unknown action:", action)
		return
	}

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func OpenFileDialog(fileType string) string {
	log.Println("In OpenFileDialog")
	var path string
	var err error

	if fileType == "dll" {
		path, err = dialog.File().Filter("DLL Files", "dll").Load()
	} else if fileType == "exe" {
		path, err = dialog.File().Filter("Executable Files", "exe").Load()
	} else {
		path, err = dialog.File().Load()
	}

	if err != nil {
		return ""
	}

	return path
}
