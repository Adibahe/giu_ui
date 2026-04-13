package rad_api

/*
#cgo CFLAGS: -O2
#include "rad_api.c"
*/
import "C"

import (
	"fmt"
	"log"
	"unsafe"
)

type RadCmd string

const (

	// Starts debugging a new instance of a target, then runs.
	CMD_LAUNCH_AND_RUN RadCmd = "launch_and_run"

	// Starts debugging a new instance of a target, then stops at the program's entry point.
	CMD_LAUNCH_AND_STEP_INTO RadCmd = "launch_and_step_into"

	// Kills the specified existing attached process(es).
	CMD_KILL RadCmd = "kill"

	// Kills all attached processes.
	CMD_KILL_ALL RadCmd = "kill_all"

	// Detaches the specified attached process(es).
	CMD_DETACH RadCmd = "detach"

	// Continues executing all attached processes.
	CMD_CONTINUE RadCmd = "continue"

	// Performs a step that goes into calls, at the instruction level.
	CMD_STEP_INTO_INST RadCmd = "step_into_inst"

	// Performs a step that skips calls, at the instruction level.
	CMD_STEP_OVER_INST RadCmd = "step_over_inst"

	// Performs a step that goes into calls, at the source code line level.
	CMD_STEP_INTO_LINE RadCmd = "step_into_line"

	// Performs a step that skips calls, at the source code line level.
	CMD_STEP_OVER_LINE RadCmd = "step_over_line"

	// Runs to the end of the current function and exits it.
	CMD_STEP_OUT RadCmd = "step_out"

	// Halts all attached processes.
	CMD_HALT RadCmd = "halt"

	// Sets the specified thread's instruction pointer at the specified address.
	CMD_SET_THREAD_IP RadCmd = "set_thread_ip"

	// Runs until a particular source line is hit. `Run` (`run`) Runs all targets after starting them if they have not been started yet.
	CMD_RUN_TO_LINE RadCmd = "run_to_line"

	// Runs all targets after starting them if they have not been
	CMD_RUN RadCmd = "run"

	// Kills all attached processes, then launches all active targets.
	CMD_RESTART RadCmd = "restart"

	// Steps once, possibly into function calls, for either source lines or instructions (whichever is selected).
	CMD_STEP_INTO RadCmd = "step_into"

	// Steps once, always over function calls, for either source lines or instructions.
	CMD_STEP_OVER RadCmd = "step_over"

	// Freezes the passed thread.
	CMD_FREEZE_THREAD RadCmd = "freeze_thread"

	// Thaws the passed thread.
	CMD_THAW_THREAD RadCmd = "thaw_thread"

	// Freezes the passed process.
	CMD_FREEZE_PROCESS RadCmd = "freeze_process"

	// Thaws the passed process.
	CMD_THAW_PROCESS RadCmd = "thaw_process"

	// Freezes the passed machine.
	CMD_FREEZE_MACHINE RadCmd = "freeze_machine"

	// Thaws the passed machine.
	CMD_THAW_MACHINE RadCmd = "thaw_machine"

	// Freezes the local machine.
	CMD_FREEZE_LOCAL_MACHINE RadCmd = "freeze_local_machine"

	// Thaws the local machine.
	CMD_THAW_LOCAL_MACHINE RadCmd = "thaw_local_machine"

	// Attaches to a process that is already running on the local machine.
	CMD_ATTACH RadCmd = "attach"

	// Exits the debugger.
	CMD_EXIT RadCmd = "exit"

	// Opens the palette.
	CMD_OPEN_PALETTE RadCmd = "open_palette"

	// Runs a command from the command palette.
	CMD_RUN_COMMAND RadCmd = "run_command"

	// Selects a thread.
	CMD_SELECT_THREAD RadCmd = "select_thread"

	// Selects an unwind frame number for the selected thread.
	CMD_SELECT_UNWIND RadCmd = "select_unwind"

	// Selects the call stack frame above the currently selected.
	CMD_UP_ONE_FRAME RadCmd = "up_one_frame"

	// Selects the call stack frame below the currently selected.
	CMD_DOWN_ONE_FRAME RadCmd = "down_one_frame"

	// Increases the window's font size by one point.
	CMD_INC_WINDOW_FONT_SIZE RadCmd = "inc_window_font_size"

	// Decreases the window's font size by one point.
	CMD_DEC_WINDOW_FONT_SIZE RadCmd = "dec_window_font_size"

	// Increases the view's font size by one point.
	CMD_INC_VIEW_FONT_SIZE RadCmd = "inc_view_font_size"

	// Decreases the view's font size by one point.
	CMD_DEC_VIEW_FONT_SIZE RadCmd = "dec_view_font_size"

	// Opens a new window.
	CMD_OPEN_WINDOW RadCmd = "open_window"

	// Opens settings for a window.
	CMD_WINDOW_SETTINGS RadCmd = "window_settings"

	// Closes an opened window.
	CMD_CLOSE_WINDOW RadCmd = "close_window"

	// Toggles fullscreen view on the active window.
	CMD_TOGGLE_FULLSCREEN RadCmd = "toggle_fullscreen"

	// Brings all windows to the front, and focuses the most recently focused window.
	CMD_BRING_TO_FRONT RadCmd = "bring_to_front"

	// Accepts the active popup prompt.
	CMD_POPUP_ACCEPT RadCmd = "popup_accept"

	// Cancels the active popup prompt.
	CMD_POPUP_CANCEL RadCmd = "popup_cancel"

	// Resets all keybindings to their defaults.
	CMD_RESET_TO_DEFAULT_BINDINGS RadCmd = "reset_to_default_bindings"

	// Resets the window to the default panel layout.
	CMD_RESET_TO_DEFAULT_PANELS RadCmd = "reset_to_default_panels"

	// Resets the window to the compact panel layout.
	CMD_RESET_TO_COMPACT_PANELS RadCmd = "reset_to_compact_panels"

	// Resets the window to the simple panel layout.
	CMD_RESET_TO_SIMPLE_PANELS RadCmd = "reset_to_simple_panels"

	// Creates a new panel to the left of the active panel.
	CMD_NEW_PANEL_LEFT RadCmd = "new_panel_left"

	// Creates a new panel at the top of the active panel.
	CMD_NEW_PANEL_UP RadCmd = "new_panel_up"

	// Creates a new panel to the right of the active panel.
	CMD_NEW_PANEL_RIGHT RadCmd = "new_panel_right"

	// Creates a new panel at the bottom of the active panel.
	CMD_NEW_PANEL_DOWN RadCmd = "new_panel_down"

	// Rotates all panels at the closest column level of the panel hierarchy.
	CMD_ROTATE_PANEL_COLUMNS RadCmd = "rotate_panel_columns"

	// Cycles the active panel forward.
	CMD_NEXT_PANEL RadCmd = "next_panel"

	// Cycles the active panel backwards.
	CMD_PREV_PANEL RadCmd = "prev_panel"

	// Focuses a panel rightward of the currently focused panel.
	CMD_FOCUS_PANEL_RIGHT RadCmd = "focus_panel_right"

	// Focuses a panel leftward of the currently focused panel.
	CMD_FOCUS_PANEL_LEFT RadCmd = "focus_panel_left"

	// Focuses a panel upward of the currently focused panel.
	CMD_FOCUS_PANEL_UP RadCmd = "focus_panel_up"

	// Focuses a panel downward of the currently focused panel.
	CMD_FOCUS_PANEL_DOWN RadCmd = "focus_panel_down"

	// Closes the currently active panel.
	CMD_CLOSE_PANEL RadCmd = "close_panel"

	// Focuses the next tab on the active panel.
	CMD_NEXT_TAB RadCmd = "next_tab"

	// Focuses the previous tab on the active panel.
	CMD_PREV_TAB RadCmd = "prev_tab"

	// Moves the selected tab right one slot.
	CMD_MOVE_TAB_RIGHT RadCmd = "move_tab_right"

	// Moves the selected tab left one slot.
	CMD_MOVE_TAB_LEFT RadCmd = "move_tab_left"

	// Opens a new tab.
	CMD_OPEN_TAB RadCmd = "open_tab"

	// Duplicates a tab.
	CMD_DUPLICATE_TAB RadCmd = "duplicate_tab"

	// Closes the currently opened tab.
	CMD_CLOSE_TAB RadCmd = "close_tab"

	// Anchors a panel's tab bar to the top of the panel.
	CMD_TAB_BAR_TOP RadCmd = "tab_bar_top"

	// Anchors a panel's tab bar to the bottom of the panel.
	CMD_TAB_BAR_BOTTOM RadCmd = "tab_bar_bottom"

	// Opens settings for a tab.
	CMD_TAB_SETTINGS RadCmd = "tab_settings"

	// Sets the debugger's current path, which is used as a starting point when browsing for files.
	CMD_SET_CURRENT_PATH RadCmd = "set_current_path"

	// Opens a file.
	CMD_OPEN RadCmd = "open"

	// Switches to the focused file's partner; or from header to implementation or vice versa.
	CMD_SWITCH_TO_PARTNER_FILE RadCmd = "switch_to_partner_file"

	// Opens the operating system's file explorer and shows the selected file.
	CMD_SHOW_FILE_IN_EXPLORER RadCmd = "show_file_in_explorer"

	// Goes to the disassembly, if any, for a given source code line.
	CMD_GO_TO_DISASSEMBLY RadCmd = "go_to_disassembly"

	// Goes to the source code, if any, for a given disassembly line.
	CMD_GO_TO_SOURCE RadCmd = "go_to_source"

	// Creates a new user file, and sets the current user path as that file's path.
	CMD_NEW_USER RadCmd = "new_user"

	// Creates a new project file, and sets the current project path as that file's path.
	CMD_NEW_PROJECT RadCmd = "new_project"

	// Opens a user file path, immediately loading it, and begins autosaving to it.
	CMD_OPEN_USER RadCmd = "open_user"

	// Opens a project file path, immediately loading it, and begins autosaving to it.
	CMD_OPEN_PROJECT RadCmd = "open_project"

	// Opens a recently used project file.
	CMD_OPEN_RECENT_PROJECT RadCmd = "open_recent_project"

	// Saves user data to a file, and sets the current user path as that path.
	CMD_SAVE_USER RadCmd = "save_user"

	// Saves project data to a file, and sets the current project path as that path.
	CMD_SAVE_PROJECT RadCmd = "save_project"

	// Writes user data to the active user file.
	CMD_WRITE_USER_DATA RadCmd = "write_user_data"

	// Writes project data to the active project file.
	CMD_WRITE_PROJECT_DATA RadCmd = "write_project_data"

	// Opens user settings.
	CMD_USER_SETTINGS RadCmd = "user_settings"

	// Opens project settings.
	CMD_PROJECT_SETTINGS RadCmd = "project_settings"

	// Edits the current selection.
	CMD_EDIT RadCmd = "edit"

	// Accepts current changes, or answers prompts in the affirmative.
	CMD_ACCEPT RadCmd = "accept"

	// Rejects current changes, exits temporary menus, or answers prompts in the negative.
	CMD_CANCEL RadCmd = "cancel"

	// Moves the cursor or selection left.
	CMD_MOVE_LEFT RadCmd = "move_left"

	// Moves the cursor or selection right.
	CMD_MOVE_RIGHT RadCmd = "move_right"

	// Moves the cursor or selection up.
	CMD_MOVE_UP RadCmd = "move_up"

	// Moves the cursor or selection down.
	CMD_MOVE_DOWN RadCmd = "move_down"

	// Moves the cursor or selection left, while selecting.
	CMD_MOVE_LEFT_SELECT RadCmd = "move_left_select"

	// Moves the cursor or selection right, while selecting.
	CMD_MOVE_RIGHT_SELECT RadCmd = "move_right_select"

	// Moves the cursor or selection up, while selecting.
	CMD_MOVE_UP_SELECT RadCmd = "move_up_select"

	// Moves the cursor or selection down, while selecting.
	CMD_MOVE_DOWN_SELECT RadCmd = "move_down_select"

	// Moves the cursor or selection left one chunk.
	CMD_MOVE_LEFT_CHUNK RadCmd = "move_left_chunk"

	// Moves the cursor or selection right one chunk.
	CMD_MOVE_RIGHT_CHUNK RadCmd = "move_right_chunk"

	// Moves the cursor or selection up one chunk.
	CMD_MOVE_UP_CHUNK RadCmd = "move_up_chunk"

	// Moves the cursor or selection down one chunk.
	CMD_MOVE_DOWN_CHUNK RadCmd = "move_down_chunk"

	// Moves the cursor or selection up one page.
	CMD_MOVE_UP_PAGE RadCmd = "move_up_page"

	// Moves the cursor or selection down one page.
	CMD_MOVE_DOWN_PAGE RadCmd = "move_down_page"

	// Moves the cursor or selection to the beginning of the relevant content.
	CMD_MOVE_UP_WHOLE RadCmd = "move_up_whole"

	// Moves the cursor or selection to the end of the relevant content.
	CMD_MOVE_DOWN_WHOLE RadCmd = "move_down_whole"

	// Moves the cursor or selection left one chunk.
	CMD_MOVE_LEFT_CHUNK_SELECT RadCmd = "move_left_chunk_select"

	// Moves the cursor or selection right one chunk.
	CMD_MOVE_RIGHT_CHUNK_SELECT RadCmd = "move_right_chunk_select"

	// Moves the cursor or selection up one chunk.
	CMD_MOVE_UP_CHUNK_SELECT RadCmd = "move_up_chunk_select"

	// Moves the cursor or selection down one chunk.
	CMD_MOVE_DOWN_CHUNK_SELECT RadCmd = "move_down_chunk_select"

	// Moves the cursor or selection up one page, while selecting.
	CMD_MOVE_UP_PAGE_SELECT RadCmd = "move_up_page_select"

	// Moves the cursor or selection down one page, while selecting.
	CMD_MOVE_DOWN_PAGE_SELECT RadCmd = "move_down_page_select"

	// Moves the cursor or selection to the beginning of the relevant content, while selecting.
	CMD_MOVE_UP_WHOLE_SELECT RadCmd = "move_up_whole_select"

	// Moves the cursor or selection to the end of the relevant content, while selecting.
	CMD_MOVE_DOWN_WHOLE_SELECT RadCmd = "move_down_whole_select"

	// Moves the cursor or selection up, while swapping the currently selected element with that upward.
	CMD_MOVE_UP_REORDER RadCmd = "move_up_reorder"

	// Moves the cursor or selection down, while swapping the currently selected element with that downward.
	CMD_MOVE_DOWN_REORDER RadCmd = "move_down_reorder"

	// Moves the cursor to the beginning of the line.
	CMD_MOVE_HOME RadCmd = "move_home"

	// Moves the cursor to the end of the line.
	CMD_MOVE_END RadCmd = "move_end"

	// Moves the cursor to the beginning of the line, while selecting.
	CMD_MOVE_HOME_SELECT RadCmd = "move_home_select"

	// Moves the cursor to the end of the line, while selecting.
	CMD_MOVE_END_SELECT RadCmd = "move_end_select"

	// Selects everything possible.
	CMD_SELECT_ALL RadCmd = "select_all"

	// Deletes a single element to the right of the cursor, or the active selection.
	CMD_DELETE_SINGLE RadCmd = "delete_single"

	// Deletes a chunk to the right of the cursor, or the active selection.
	CMD_DELETE_CHUNK RadCmd = "delete_chunk"

	// Deletes a single element to the left of the cursor, or the active selection.
	CMD_BACKSPACE_SINGLE RadCmd = "backspace_single"

	// Deletes a chunk to the left of the cursor, or the active selection.
	CMD_BACKSPACE_CHUNK RadCmd = "backspace_chunk"

	// Copies the active selection to the clipboard.
	CMD_COPY RadCmd = "copy"

	// Copies the active selection to the clipboard, then deletes it.
	CMD_CUT RadCmd = "cut"

	// Pastes the current contents of the clipboard.
	CMD_PASTE RadCmd = "paste"

	// Inserts the text that was used to cause this command.
	CMD_INSERT_TEXT RadCmd = "insert_text"

	// Moves the cursor or selection to the next element.
	CMD_MOVE_NEXT RadCmd = "move_next"

	// Moves the cursor or selection to the previous element.
	CMD_MOVE_PREV RadCmd = "move_prev"

	// Jumps to a line number in the current code file.
	CMD_GOTO_LINE RadCmd = "goto_line"

	// Jumps to an address in the current memory or disassembly view.
	CMD_GOTO_ADDRESS RadCmd = "goto_address"

	// Snaps the current code view to center the cursor.
	CMD_CENTER_CURSOR RadCmd = "center_cursor"

	// Snaps the current code view to contain the cursor.
	CMD_CONTAIN_CURSOR RadCmd = "contain_cursor"

	// Searches the current code file forward (from the cursor) for the last searched string.
	CMD_FIND_NEXT RadCmd = "find_next"

	// Searches the current code file backwards (from the cursor) for the last searched string.
	CMD_FIND_PREV RadCmd = "find_prev"

	// Jumps to the passed thread in either source code, disassembly, or both if they're already open.
	CMD_FIND_THREAD RadCmd = "find_thread"

	// Jumps to the selected thread in either source code, disassembly, or both if they're already open.
	CMD_FIND_SELECTED_THREAD RadCmd = "find_selected_thread"

	// Searches for the passed string as a file, a symbol in debug info, and more, then jumps to it if possible.
	CMD_GOTO_NAME RadCmd = "goto_name"

	// Searches for the text at the cursor as a file, a symbol in debug info, and more, then jumps to it if possible.
	CMD_GOTO_NAME_AT_CURSOR RadCmd = "goto_name_at_cursor"

	// Adds or removes an expression to an opened watch view.
	CMD_TOGGLE_WATCH_EXPR RadCmd = "toggle_watch_expr"

	// Adds or removes the expression that the cursor or selection is currently over to an opened watch view.
	CMD_TOGGLE_WATCH_EXPR_AT_CURSOR RadCmd = "toggle_watch_expr_at_cursor"

	// Adds or removes the expression that the mouse is currently over to an opened watch view.
	CMD_TOGGLE_WATCH_EXPR_AT_MOUSE RadCmd = "toggle_watch_expr_at_mouse"

	// Places a breakpoint at a given location (file path and line number, address, or symbol name).
	CMD_ADD_BREAKPOINT RadCmd = "add_breakpoint"

	// Places or removes a breakpoint at a given location (file path and line number, address, or symbol name).
	CMD_TOGGLE_BREAKPOINT RadCmd = "toggle_breakpoint"

	// Enables a breakpoint.
	CMD_ENABLE_BREAKPOINT RadCmd = "enable_breakpoint"

	// Disables a breakpoint.
	CMD_DISABLE_BREAKPOINT RadCmd = "disable_breakpoint"

	// Removes all breakpoints.
	CMD_CLEAR_BREAKPOINTS RadCmd = "clear_breakpoints"

	// Lists all breakpoints.
	CMD_LIST_BREAKPOINTS RadCmd = "list_breakpoints"

	// Clears all output.
	CMD_CLEAR_OUTPUT RadCmd = "clear_output"

	// Places a watch pin at a given location (file path and line number or address).
	CMD_ADD_WATCH_PIN RadCmd = "add_watch_pin"

	// Loads a debug info file.
	CMD_LOAD_DEBUG_INFO RadCmd = "load_debug_info"

	// Unloads a debug info file.
	CMD_UNLOAD_DEBUG_INFO RadCmd = "unload_debug_info"

	// Sets the selected thread's instruction pointer to the cursor's position.
	CMD_SET_NEXT_STATEMENT RadCmd = "set_next_statement"

	// Adds a new target.
	CMD_ADD_TARGET RadCmd = "add_target"

	// Selects a target.
	CMD_SELECT_TARGET RadCmd = "select_target"

	// Enables a target, in addition to all targets currently enabled.
	CMD_ENABLE_TARGET RadCmd = "enable_target"

	// Disables a target.
	CMD_DISABLE_TARGET RadCmd = "disable_target"

	// Removes a target.
	CMD_REMOVE_TARGET RadCmd = "remove_target"

	// Registers the RAD debugger as the just-in-time (JIT) debugger used by the operating system.
	CMD_REGISTER_AS_JIT_DEBUGGER RadCmd = "register_as_jit_debugger"

	// Finds a specific source code location given file, line, and column coordinates. Opens the file if necessary.
	CMD_FIND_CODE_LOCATION RadCmd = "find_code_location"

	// Begins searching within the active interface.
	CMD_SEARCH RadCmd = "search"

	// Begins searching backwards within the active interface.
	CMD_SEARCH_BACKWARDS RadCmd = "search_backwards"

	// Opens a new event buffer, to which debugger events will be written, for external processing.
	CMD_OPEN_EVENT_BUFFER RadCmd = "open_event_buffer"

	// Closes an existing event buffer.
	CMD_CLOSE_EVENT_BUFFER RadCmd = "close_event_buffer"

	// Opens and closes the developer menu.
	CMD_TOGGLE_DEV_MENU RadCmd = "toggle_dev_menu"

	// Logs a marker in the application log, to denote specific points in time within the log.
	CMD_LOG_MARKER RadCmd = "log_marker"
)

type RadError uint8

const (
	RD_SUCCESS        RadError = 0
	RD_PID_ERROR      RadError = 1
	RD_IPC_INIT_ERROR RadError = 2

	RD_OUT_OF_MEMMORY     RadError = 3
	RD_FAILED_S2M_LOCK    RadError = 4
	RD_TIMEOUT_REPLY_WAIT RadError = 5
	RD_TIMEOUT_M2S_WAIT   RadError = 6
)

type RadIpcState struct {
	is_connected bool
}

func radErrorToString(e RadError) string {
	res := ""

	switch e {
	case RD_SUCCESS:
		res = "RD_SUCCESS"

	case RD_PID_ERROR:
		res = "RD_PID_ERROR"

	case RD_IPC_INIT_ERROR:
		res = "RD_IPC_INIT_ERROR"

	case RD_OUT_OF_MEMMORY:
		res = "RD_OUT_OF_MEMMORY"

	case RD_FAILED_S2M_LOCK:
		res = "RD_FAILED_S2M_LOCK"

	case RD_TIMEOUT_REPLY_WAIT:
		res = "RD_TIMEOUT_REPLY_WAIT"

	case RD_TIMEOUT_M2S_WAIT:
		res = "RD_TIMEOUT_M2S_WAIT"
	}

	return res
}

// Initialize the connection to raddbg instance
func (r *RadIpcState) Init() error {
	res := RadError(C.RadInit())

	if res == RD_SUCCESS {
		r.is_connected = true
	} else {
		return fmt.Errorf(radErrorToString(res))
	}

	return nil
}

// Release the connection to raddbg instance
func (r *RadIpcState) Release() {
	if r.is_connected {
		C.RadRelease()
	}
}

// Send Commands to running raddbg instance,
// Pass Commands like CMD_RUN without any arguments(args => "")
// Commands like CMD_ADD_BREAKPOINT requires argument like
// "E:\tmp\proj-mini\p7\hookdll\hook.cpp:48"
func (r *RadIpcState) SendCommand(cmd_str RadCmd, args string) error {
	log.Println("Inside SendCommand")
	cmd_str = RadCmd(string(cmd_str) + " " + args)
	if r.is_connected {
		ln := len(cmd_str)
		cmd := C.CBytes(append([]byte(cmd_str), 0))

		res := RadError(C.RadSendCommand((*C.char)(cmd), C.uint64_t(ln)))

		C.free(unsafe.Pointer(cmd))
		if res == 0 {
			return nil
		} else {
			return fmt.Errorf(radErrorToString(res))
		}
	} else {
		return fmt.Errorf("Raddbg is not connected")
	}
}
