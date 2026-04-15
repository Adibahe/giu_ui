package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"

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

func onPageReload(_ []message) {
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

	state.manual_running = true

	tp, err := winproc.Start(state.targetPath, true)
	state.targetProcess = tp
	if err != nil {
		log.Fatalf("The target process execution failed %s, %s", err, state.targetPath)
	}

	r, err := winproc.Start("raddbg.exe", false)
	if err != nil {
		log.Fatalf("RAD process execution failed %s", err)
	}
	state.radProcess = r

	time.Sleep(1000 * time.Millisecond)

	state.rad.Init()
	err = state.rad.SendCommand(rad_api.CMD_ATTACH, fmt.Sprintf("%d", state.targetProcess.Pid))
	if err != nil {
		log.Fatalf("Failure attaching to the target %v", err)
	}

	time.Sleep(1000 * time.Millisecond)

	err = state.targetProcess.Resume()
	if err != nil {
		log.Fatalf("Failure resuming the targetProcess %v", err)
	}

	time.Sleep(400 * time.Millisecond)

	err = state.rad.SendCommand(rad_api.CMD_RUN, "")
	if err != nil {
		log.Fatalf("Failure sending first run command %v", err)
	}

	go waitForTargetExit()
}

func exitSession() {

	state.manual_running = false
	state.auto_running = false
	state.rad.SendCommand(rad_api.CMD_KILL_ALL, "")
	state.targetProcess.Close()
	state.rad.Release()
	state.radProcess.Close()
}

func waitForTargetExit() {
	state.targetProcess.Wait()
	exitSession()
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
				go waitForTargetExit()
			}
		} else {
			if state.auto_running {
			} else {
				state.auto_running = true
				startSession()
				go waitForTargetExit()
			}
		}

	case "step":
		err = state.rad.SendCommand(rad_api.CMD_RUN, "")

	case "halt":
		fmt.Println("halt")
		state.auto_running = false

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

	switch fileType {
	case "dll":
		{
			path, err = dialog.File().Filter("DLL Files", "dll").Load()
		}
	case "exe":
		{
			path, err = dialog.File().Filter("Executable Files", "exe").Load()
		}
	default:
		{
			path, err = dialog.File().Load()
		}
	}

	if err != nil {
		return ""
	}

	return path
}
