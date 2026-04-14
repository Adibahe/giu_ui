package main

import (
	"fmt"
	"log"
	"os/exec"

	"giu_ui/winproc"

	"github.com/YUSHACOD/rad_api"
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

var state struct {
	targetPath     string
	dllPath        string
	targetProcess  *winproc.Process
	radProcess     *winproc.Process
	target_spawned bool
	rad            rad_api.RadIpcState
	auto_running   bool
	manual_running bool
}

func startSession() {

	state.auto_running = true

	tp, err := winproc.Start(state.targetPath, true)
	state.targetProcess = tp
	if err != nil {
		log.Fatalf("The target process execution failed %s", err)
	}

	r, err := winproc.Start("raddbg.exe", false)
	if err != nil {
		log.Fatalf("RAD process execution failed %s", err)
	}
	state.radProcess = r

	state.rad.Init()
	state.rad.SendCommand(rad_api.CMD_LAUNCH_AND_RUN, "")
}

func exitSession() {

	state.targetProcess.Close()
	state.rad.Release()
	state.radProcess.Close()
}

func DebugCommand(action string) {
	var err error

	switch action {

	case "manual_run":
		fmt.Println("manual_run")
		if !(state.manual_running || state.auto_running) {
			startSession()
		}

	case "auto_run":
		fmt.Println("auto_run")
		if state.manual_running {
			if state.auto_running {
			} else {
				state.auto_running = true
				go func() {
					state.targetProcess.Wait()
				}()
			}
		} else {
			if state.auto_running {
			} else {
				state.auto_running = true
				go func() {
					state.targetProcess.Wait()
					exitSession()
				}()
			}
		}

	case "step":
		err = state.rad.SendCommand(rad_api.CMD_RUN, "")

	case "halt":
		fmt.Println("halt")
		state.auto_running = false
		err = state.rad.SendCommand(rad_api.CMD_HALT, "")

	case "exit":
		exitSession()

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
