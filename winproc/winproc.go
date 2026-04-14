package winproc

import (
	"fmt"
	"unsafe" // Added for Sizeof
	"golang.org/x/sys/windows"
)

type Process struct {
	Pid      uint32
	hProcess windows.Handle
	hThread  windows.Handle
}

func Start(cmdLine string, suspended bool) (*Process, error) {
	argv, err := windows.UTF16PtrFromString(cmdLine)
	if err != nil {
		return nil, fmt.Errorf("invalid command line: %w", err)
	}

	var si windows.StartupInfo
	var pi windows.ProcessInformation
	
	// Correct way: Use unsafe.Sizeof to get the struct size
	si.Cb = uint32(unsafe.Sizeof(si))

	var flags uint32
	if suspended {
		flags = windows.CREATE_SUSPENDED
	}

	err = windows.CreateProcess(
		nil,
		argv,
		nil,
		nil,
		false,
		flags,
		nil,
		nil,
		&si,
		&pi,
	)

	if err != nil {
		return nil, err
	}

	return &Process{
		Pid:      pi.ProcessId,
		hProcess: pi.Process,
		hThread:  pi.Thread,
	}, nil
}

func (p *Process) Resume() error {
	_, err := windows.ResumeThread(p.hThread)
	return err
}

func (p *Process) Wait() (uint32, error) {
	_, err := windows.WaitForSingleObject(p.hProcess, windows.INFINITE)
	if err != nil {
		return 0, err
	}

	var exitCode uint32
	err = windows.GetExitCodeProcess(p.hProcess, &exitCode)
	return exitCode, err
}

func (p *Process) Close() {
	windows.CloseHandle(p.hThread)
	windows.CloseHandle(p.hProcess)
}
