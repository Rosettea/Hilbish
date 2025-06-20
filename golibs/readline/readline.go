package readline

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"syscall"
)

var rxMultiline = regexp.MustCompile(`[\r\n]+`)

// Readline displays the readline prompt.
// It will return a string (user entered data) or an error.
func (rl *Readline) Readline() (string, error) {
	fd := int(os.Stdin.Fd())
	state, err := MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer Restore(fd, state)

	// In Vim mode, we always start in Input mode. The prompt needs this.
	rl.modeViMode = VimInsert

	// Prompt Init
	// Here we have to either print prompt
	// and return new line (multiline)
	if rl.Multiline {
		fmt.Println(rl.mainPrompt)
	}
	rl.stillOnRefresh = false
	rl.computePrompt() // initialise the prompt for first print

	// Line Init & Cursor
	rl.line = []rune{}
	rl.currentComp = []rune{} // No virtual completion yet
	rl.lineComp = []rune{}    // So no virtual line either
	rl.modeViMode = VimInsert
	rl.pos = 0
	rl.posY = 0
	rl.tcPrefix = ""

	// Completion && infos init
	rl.resetInfoText()
	rl.resetTabCompletion()
	rl.getInfoText()

	// History Init
	// We need this set to the last command, so that we can access it quickly
	rl.histOffset = 0
	rl.viUndoHistory = []undoItem{{line: "", pos: 0}}

	// Multisplit
	if len(rl.multisplit) > 0 {
		r := []rune(rl.multisplit[0])
		if len(r) >= 1 {
			rl.editorInput(r)
		}

		rl.carridgeReturn()
		if len(rl.multisplit) > 1 {
			rl.multisplit = rl.multisplit[1:]
		} else {
			rl.multisplit = []string{}
		}
		return string(rl.line), nil
	}

	// Finally, print any info or completions
	// if the TabCompletion engines so desires
	rl.renderHelpers()

	// Start handling keystrokes. Classified by subject for most.
	for {
		rl.viUndoSkipAppend = false
		b := make([]byte, 1024)
		var i int

		if !rl.skipStdinRead {
			var err error
			i, err = os.Stdin.Read(b)
			if err != nil {
				if errors.Is(err, syscall.EAGAIN) {
					err = syscall.SetNonblock(syscall.Stdin, false)
					if err == nil {
						continue
					}
				}
				return "", err
			}
		}

		rl.skipStdinRead = false
		r := []rune(string(b))
		if rl.RawInputCallback != nil {
			rl.RawInputCallback(r[:i])
		}

		if isMultiline(r[:i]) || len(rl.multiline) > 0 {
			rl.multiline = append(rl.multiline, b[:i]...)
			if i == len(b) {
				continue
			}

			if !rl.allowMultiline(rl.multiline) {
				rl.multiline = []byte{}
				continue
			}

			s := string(rl.multiline)
			rl.multisplit = rxMultiline.Split(s, -1)

			r = []rune(rl.multisplit[0])
			rl.modeViMode = VimInsert
			rl.editorInput(r)
			rl.carridgeReturn()
			rl.multiline = []byte{}
			if len(rl.multisplit) > 1 {
				rl.multisplit = rl.multisplit[1:]
			} else {
				rl.multisplit = []string{}
			}
			return string(rl.line), nil
		}

		s := string(r[:i])
		if rl.evtKeyPress[s] != nil {
			rl.clearHelpers()

			ret := rl.evtKeyPress[s](s, rl.line, rl.pos)

			rl.clearLine()
			rl.line = append(ret.NewLine, []rune{}...)
			rl.updateHelpers() // rl.echo
			rl.pos = ret.NewPos

			if ret.ClearHelpers {
				rl.resetHelpers()
			} else {
				rl.updateHelpers()
			}

			if len(ret.InfoText) > 0 {
				rl.infoText = ret.InfoText
				rl.clearHelpers()
				rl.renderHelpers()
			}
			if !ret.ForwardKey {
				continue
			}
			if ret.CloseReadline {
				rl.clearHelpers()
				return string(rl.line), nil
			}
		}

		// Before anything: we can never be both in modeTabCompletion and compConfirmWait,
		// because we need to confirm before entering completion. If both are true, there
		// is a problem (at least, the user has escaped the confirm hint some way).
		if (rl.modeTabCompletion && rl.searchMode != HistoryFind) && rl.compConfirmWait {
			rl.compConfirmWait = false
		}

		switch b[0] {
		// Errors & Returns --------------------------------------------------------------------------------
		case charCtrlC:
			if rl.modeTabCompletion {
				rl.resetVirtualComp(true)
				rl.resetHelpers()
				rl.renderHelpers()
				continue
			}
			rl.clearHelpers()
			return "", CtrlC

		case charEOF: // ctrl d
			if len(rl.line) == 0 {
				rl.clearHelpers()
				return "", EOF
			}
			if rl.modeTabFind {
				rl.backspaceTabFind()
			} else {
				if rl.pos < len(rl.line) {
					rl.deleteBackspace(true)
				}
			}

		// Clear screen
		case charCtrlL:
			print(seqClearScreen)
			print(seqCursorTopLeft)
			if rl.Multiline {
				fmt.Println(rl.mainPrompt)
			}
			print(seqClearScreenBelow)

			rl.resetInfoText()
			rl.getInfoText()
			rl.renderHelpers()

		// Line Editing ------------------------------------------------------------------------------------
		case charCtrlU:
			if rl.modeTabCompletion {
				rl.resetVirtualComp(true)
			}
			// Delete everything from the beginning of the line to the cursor position
			rl.saveBufToRegister(rl.line[:rl.pos])
			rl.deleteToBeginning()
			rl.resetHelpers()
			rl.updateHelpers()

		case charCtrlK:
			if rl.modeTabCompletion {
				rl.resetVirtualComp(true)
			}
			// Delete everything after the cursor position
			rl.saveBufToRegister(rl.line[rl.pos:])
			rl.deleteToEnd()
			rl.resetHelpers()
			rl.updateHelpers()

		case charBackspace, charBackspace2:
			// When currently in history completion, we refresh and automatically
			// insert the first (filtered) candidate, virtually
			if rl.modeAutoFind && rl.searchMode == HistoryFind {
				rl.resetVirtualComp(true)
				rl.backspaceTabFind()

				// Then update the printing, with the new candidate
				rl.updateVirtualComp()
				rl.renderHelpers()
				rl.viUndoSkipAppend = true
				continue
			}

			// Normal completion search does only refresh the search pattern and the comps
			if rl.modeTabFind || rl.modeAutoFind {
				rl.resetVirtualComp(false)
				rl.backspaceTabFind()
				rl.renderHelpers()
				rl.viUndoSkipAppend = true
			} else {
				// Always cancel any virtual completion
				rl.resetVirtualComp(false)

				// Vim mode has different behaviors
				if rl.InputMode == Vim {
					if rl.modeViMode == VimInsert {
						rl.backspace(false)
					} else if rl.pos != 0 {
						rl.pos--
					}
					rl.renderHelpers()
					continue
				}

				// Else emacs deletes a character
				rl.backspace(false)
				rl.renderHelpers()
			}

		// Emacs Bindings ----------------------------------------------------------------------------------
		case charCtrlW:
			if rl.modeTabCompletion {
				rl.resetVirtualComp(false)
			}
			// This is only available in Insert mode
			if rl.modeViMode != VimInsert {
				continue
			}
			rl.saveToRegister(rl.viJumpB(tokeniseLine))
			rl.viDeleteByAdjust(rl.viJumpB(tokeniseLine))
			rl.updateHelpers()

		case charCtrlY:
			if rl.modeTabCompletion {
				rl.resetVirtualComp(false)
			}
			// paste after the cursor position
			rl.viUndoSkipAppend = true
			buffer := rl.pasteFromRegister()
			rl.insert(buffer)
			rl.updateHelpers()

		case charCtrlE:
			if rl.modeTabCompletion {
				rl.resetVirtualComp(false)
			}
			// This is only available in Insert mode
			if rl.modeViMode != VimInsert {
				continue
			}
			if len(rl.line) > 0 {
				rl.pos = len(rl.line)
			}
			rl.viUndoSkipAppend = true
			rl.updateHelpers()

		case charCtrlA:
			if rl.modeTabCompletion {
				rl.resetVirtualComp(false)
			}
			// This is only available in Insert mode
			if rl.modeViMode != VimInsert {
				continue
			}
			rl.viUndoSkipAppend = true
			rl.pos = 0
			rl.updateHelpers()

		// Command History ---------------------------------------------------------------------------------

		// NOTE: The alternative history source is triggered by Alt+r,
		// but because this is a sequence, the alternative history code
		// trigger is in the below rl.escapeSeq(r) function.
		case charCtrlR:
			rl.resetVirtualComp(false)
			// For some modes only, if we are in vim Keys mode,
			// we toogle back to insert mode. For others, we return
			// without getting the completions.
			if rl.modeViMode != VimInsert {
				rl.modeViMode = VimInsert
				rl.computePrompt()
			}

			rl.mainHist = true // false before
			rl.searchMode = HistoryFind
			rl.modeAutoFind = true
			rl.modeTabCompletion = true

			rl.modeTabFind = true
			rl.updateTabFind([]rune{})
			rl.updateVirtualComp()
			rl.renderHelpers()
			rl.viUndoSkipAppend = true

		// Tab Completion & Completion Search ---------------------------------------------------------------
		case charTab:
			// The user cannot show completions if currently in Vim Normal mode
			if rl.InputMode == Vim && rl.modeViMode != VimInsert {
				continue
			}

			// If we have asked for completions, already printed, and we want to move selection.
			if rl.modeTabCompletion && !rl.compConfirmWait {
				rl.tabCompletionSelect = true
				rl.moveTabCompletionHighlight(1, 0)
				rl.updateVirtualComp()
				rl.renderHelpers()
				rl.viUndoSkipAppend = true
			} else {
				// Else we might be asked to confirm printing (if too many suggestions), or not.
				rl.getTabCompletion()

				// If too many completions and no yet confirmed, ask user for completion
				// comps, lines := rl.getCompletionCount()
				// if ((lines > GetTermLength()) || (lines > rl.MaxTabCompleterRows)) && !rl.compConfirmWait {
				//         sentence := fmt.Sprintf("%s show all %d completions (%d lines) ? tab to confirm",
				//                 FOREWHITE, comps, lines)
				//         rl.promptCompletionConfirm(sentence)
				//         continue
				// }

				rl.compConfirmWait = false
				rl.modeTabCompletion = true

				// Also here, if only one candidate is available, automatically
				// insert it and don't bother printing completions.
				// Quit the tab completion mode to avoid asking to the user
				// to press Enter twice to actually run the command.
				if rl.hasOneCandidate() {
					rl.insertCandidate()

					// Refresh first, and then quit the completion mode
					rl.updateHelpers() // REDUNDANT WITH getTabCompletion()
					rl.viUndoSkipAppend = true
					rl.resetTabCompletion()
					continue
				}

				rl.updateHelpers() // REDUNDANT WITH getTabCompletion()
				rl.viUndoSkipAppend = true
				continue
			}

		case charCtrlF:
			rl.resetVirtualComp(true)

			if !rl.modeTabCompletion {
				rl.modeTabCompletion = true
			}

			if rl.compConfirmWait {
				rl.resetHelpers()
			}

			// Both these settings apply to when we already
			// are in completion mode and when we are not.
			rl.searchMode = CompletionFind
			rl.modeAutoFind = true

			// Switch from history to completion search
			if rl.modeTabCompletion && rl.searchMode == HistoryFind {
				rl.searchMode = CompletionFind
			}

			rl.updateTabFind([]rune{})
			rl.viUndoSkipAppend = true

		case charCtrlG:
			if rl.modeAutoFind && rl.searchMode == HistoryFind {
				rl.resetVirtualComp(false)
				rl.resetTabFind()
				rl.resetHelpers()
				rl.renderHelpers()
				continue
			}

			if rl.modeAutoFind {
				rl.resetTabFind()
				rl.resetHelpers()
				rl.renderHelpers()
			}

		case charCtrlUnderscore:
			rl.undoLast()
			rl.viUndoSkipAppend = true

		case '\r':
			fallthrough
		case '\n':
			if rl.modeTabCompletion {
				cur := rl.getCurrentGroup()

				// Check that there is a group indeed, as we might have no completions.
				if cur == nil {
					rl.clearHelpers()
					rl.resetTabCompletion()
					rl.renderHelpers()
					continue
				}

				// IF we have a prefix and completions printed, but no candidate
				// (in which case the completion is ""), we immediately return.
				completion := cur.getCurrentCell(rl)
				prefix := len(rl.tcPrefix)
				if prefix > len(completion) {
					rl.carridgeReturn()
					return string(rl.line), nil
				}

				// Else, we insert the completion candidate in the real input line.
				// By default we add a space, unless completion group asks otherwise.
				rl.compAddSpace = true
				rl.resetVirtualComp(false)

				// If we were in history completion, immediately execute the line.
				if rl.modeAutoFind && rl.searchMode == HistoryFind {
					rl.carridgeReturn()
					return string(rl.line), nil
				}

				// Reset completions and update input line
				rl.clearHelpers()
				rl.resetTabCompletion()
				rl.renderHelpers()

				continue
			}
			rl.carridgeReturn()
			return string(rl.line), nil

		// Vim --------------------------------------------------------------------------------------
		case charEscape:

			// If we were waiting for completion confirm, abort
			if rl.compConfirmWait {
				rl.compConfirmWait = false
				rl.renderHelpers()
			}

			// We always refresh the completion candidates, except if we are currently
			// cycling through them, because then it would just append the candidate.
			if rl.modeTabCompletion {
				if string(r[:i]) != seqShiftTab &&
					string(r[:i]) != seqForwards && string(r[:i]) != seqBackwards &&
					string(r[:i]) != seqUp && string(r[:i]) != seqDown {
					// basically only applies except on 1st ctrl r open
					// so if we have not explicitly selected something
					// (tabCompletionSelect is false) drop virtual completion
					rl.resetVirtualComp(!rl.tabCompletionSelect)
				}
			}

			// Once helpers of all sorts are cleared, we can process
			// the change of input modes, etc.
			rl.escapeSeq(r[:i])

		// Dispatch --------------------------------------------------------------------------------------
		default:

			// If we were waiting for completion confirm, abort
			if rl.compConfirmWait {
				rl.resetVirtualComp(false)
				rl.compConfirmWait = false
				rl.renderHelpers()
			}

			// When currently in history completion, we refresh and automatically
			// insert the first (filtered) candidate, virtually
			if rl.modeAutoFind && rl.searchMode == HistoryFind {
				rl.resetVirtualComp(true)
				rl.updateTabFind(r[:i])
				rl.updateVirtualComp()
				rl.renderHelpers()
				rl.viUndoSkipAppend = true
				continue
			}

			// Not sure that CompletionFind is useful, nor one of the other two
			if rl.modeAutoFind || rl.modeTabFind {
				rl.resetVirtualComp(false)
				rl.updateTabFind(r[:i])
				rl.renderHelpers()
				rl.viUndoSkipAppend = true
				continue
			} else {
				rl.resetVirtualComp(false)
				rl.editorInput(r[:i])
				if len(rl.multiline) > 0 && rl.modeViMode == VimKeys {
					rl.skipStdinRead = true
				}
			}

			rl.clearHelpers()
		}

		rl.undoAppendHistory()
	}
}

// editorInput is an unexported function used to determine what mode of text
// entry readline is currently configured for and then update the line entries
// accordingly.
func (rl *Readline) editorInput(r []rune) {
	if len(r) == 0 {
		return
	}

	switch rl.modeViMode {
	case VimKeys:
		rl.vi(r[0])
		rl.refreshVimStatus()

	case VimDelete:
		rl.viDelete(r[0])
		rl.refreshVimStatus()

	case VimReplaceOnce:
		rl.modeViMode = VimKeys
		rl.deleteX()
		rl.insert([]rune{r[0]})
		rl.refreshVimStatus()

	case VimReplaceMany:
		for _, char := range r {
			if rl.pos != len(rl.line) {
				rl.deleteX()
			}
			rl.insert([]rune{char})
		}
		rl.refreshVimStatus()

	default:
		// Don't insert control keys
		if r[0] >= 1 && r[0] <= 31 {
			return
		}
		// We reset the history nav counter each time we come here:
		// We don't need it when inserting text.
		rl.histNavIdx = 0
		rl.insert(r)
		rl.writeHintText()
	}

	rl.echoRightPrompt()

	if len(rl.multisplit) == 0 {
		rl.syntaxCompletion()
	}
}

// viEscape - In case th user is using Vim input, and the escape sequence has not
// been handled by other cases, we dispatch it to Vim and handle a few cases here.
func (rl *Readline) viEscape(r []rune) {

	// Sometimes the escape sequence is interleaved with another one,
	// but key strokes might be in the wrong order, so we double check
	// and escape the mode only if needed.
	if rl.modeViMode == VimInsert && len(r) == 1 && r[0] == 27 {
		if len(rl.line) > 0 && rl.pos > 0 {
			rl.pos--
		}
		rl.modeViMode = VimKeys
		rl.viIteration = ""
		rl.refreshVimStatus()
		return
	}
}

func (rl *Readline) escapeSeq(r []rune) {
	switch string(r) {
	// Vim escape sequences & dispatching --------------------------------------------------------
	case string(charEscape):
		switch {
		case rl.modeAutoFind:
			rl.resetVirtualComp(true)
			rl.resetTabFind()
			rl.clearHelpers()
			rl.resetTabCompletion()
			rl.resetHelpers()
			rl.renderHelpers()

		case rl.modeTabFind:
			rl.resetVirtualComp(true)
			rl.resetTabFind()
			rl.resetTabCompletion()

		case rl.modeTabCompletion:
			rl.clearHelpers()
			rl.resetTabCompletion()
			rl.renderHelpers()

		default:
			// No matter the input mode, we exit
			// any completion confirm if there's one.
			if rl.compConfirmWait {
				rl.compConfirmWait = false
				rl.clearHelpers()
				rl.renderHelpers()
				return
			}

			// If we are in Vim mode, the escape key has its usage.
			// Otherwise in emacs mode the escape key does nothing.
			if rl.InputMode == Vim {
				rl.viEscape(r)
				return
			}

			// This refreshed and actually prints the new Vim status
			// if we have indeed change the Vim mode.
			rl.clearHelpers()
			rl.renderHelpers()

		}
		rl.viUndoSkipAppend = true

	// Tab completion movements ------------------------------------------------------------------
	case seqShiftTab:
		if rl.modeTabCompletion && !rl.compConfirmWait {
			rl.tabCompletionReverse = true
			rl.moveTabCompletionHighlight(-1, 0)
			rl.updateVirtualComp()
			rl.tabCompletionReverse = false
			rl.renderHelpers()
			rl.viUndoSkipAppend = true
			return
		}

	case seqUp:
		if rl.modeTabCompletion {
			rl.tabCompletionSelect = true
			rl.tabCompletionReverse = true
			rl.moveTabCompletionHighlight(-1, 0)
			rl.updateVirtualComp()
			rl.tabCompletionReverse = false
			rl.renderHelpers()
			return
		}
		rl.mainHist = true
		rl.walkHistory(1)
		moveCursorForwards(len(rl.line) - rl.pos)
		rl.pos = len(rl.line)

	case seqDown:
		if rl.modeTabCompletion {
			rl.tabCompletionSelect = true
			rl.moveTabCompletionHighlight(1, 0)
			rl.updateVirtualComp()
			rl.renderHelpers()
			return
		}
		rl.mainHist = true
		rl.walkHistory(-1)
		moveCursorForwards(len(rl.line) - rl.pos)
		rl.pos = len(rl.line)

	case seqForwards:
		if rl.modeTabCompletion {
			rl.tabCompletionSelect = true
			rl.moveTabCompletionHighlight(1, 0)
			rl.updateVirtualComp()
			rl.renderHelpers()
			return
		}

		rl.insertHintText()

		if (rl.modeViMode == VimInsert && rl.pos < len(rl.line)) ||
			(rl.modeViMode != VimInsert && rl.pos < len(rl.line)-1) {
			rl.moveCursorByAdjust(1)
		}
		rl.updateHelpers()
		rl.viUndoSkipAppend = true

	case seqBackwards:
		if rl.modeTabCompletion {
			rl.tabCompletionSelect = true
			rl.tabCompletionReverse = true
			rl.moveTabCompletionHighlight(-1, 0)
			rl.updateVirtualComp()
			rl.tabCompletionReverse = false
			rl.renderHelpers()
			return
		}
		rl.moveCursorByAdjust(-1)
		rl.viUndoSkipAppend = true
		rl.updateHelpers()

	// Registers -------------------------------------------------------------------------------
	case seqAltQuote:
		if rl.modeViMode != VimInsert {
			return
		}
		rl.modeTabCompletion = true
		rl.modeAutoFind = true
		rl.searchMode = RegisterFind
		// Else we might be asked to confirm printing (if too many suggestions), or not.
		rl.getTabCompletion()
		rl.viUndoSkipAppend = true
		rl.renderHelpers()

	// Movement -------------------------------------------------------------------------------
	case seqCtrlLeftArrow:
		rl.moveCursorByAdjust(rl.viJumpB(tokeniseLine))
		rl.updateHelpers()
		return
	case seqCtrlRightArrow:
		rl.insert(rl.hintText)
		rl.moveCursorByAdjust(rl.viJumpW(tokeniseLine))
		rl.updateHelpers()
		return

	case seqDelete, seqDelete2:
		if rl.modeTabFind {
			rl.backspaceTabFind()
		} else {
			if rl.pos < len(rl.line) {
				rl.deleteBackspace(true)
			}
		}

	case seqHome, seqHomeSc:
		if rl.modeTabCompletion {
			return
		}
		rl.moveCursorByAdjust(-rl.pos)
		rl.updateHelpers()
		rl.viUndoSkipAppend = true

	case seqEnd, seqEndSc:
		if rl.modeTabCompletion {
			return
		}
		rl.moveCursorByAdjust(len(rl.line) - rl.pos)
		rl.updateHelpers()
		rl.viUndoSkipAppend = true

	case seqAltB:
		if rl.modeTabCompletion {
			return
		}

		// This is only available in Insert mode
		if rl.modeViMode != VimInsert {
			return
		}

		move := rl.emacsBackwardWord(tokeniseLine)
		rl.moveCursorByAdjust(-move)
		rl.updateHelpers()

	case seqAltF:
		if rl.modeTabCompletion {
			return
		}

		// This is only available in Insert mode
		if rl.modeViMode != VimInsert {
			return
		}

		move := rl.emacsForwardWord(tokeniseLine)
		rl.moveCursorByAdjust(move)
		rl.updateHelpers()

	case seqAltR:
		rl.resetVirtualComp(false)
		// For some modes only, if we are in vim Keys mode,
		// we toogle back to insert mode. For others, we return
		// without getting the completions.
		if rl.modeViMode != VimInsert {
			rl.modeViMode = VimInsert
		}

		rl.mainHist = false // true before
		rl.searchMode = HistoryFind
		rl.modeAutoFind = true
		rl.modeTabCompletion = true

		rl.modeTabFind = true
		rl.updateTabFind([]rune{})
		rl.viUndoSkipAppend = true

	case seqAltBackspace:
		if rl.modeTabCompletion {
			rl.resetVirtualComp(false)
		}
		// This is only available in Insert mode
		if rl.modeViMode != VimInsert {
			return
		}

		rl.saveToRegister(rl.viJumpB(tokeniseLine))
		rl.viDeleteByAdjust(rl.viJumpB(tokeniseLine))
		rl.updateHelpers()

	case seqCtrlDelete, seqCtrlDelete2, seqAltD:
		if rl.modeTabCompletion {
			rl.resetVirtualComp(false)
		}
		rl.saveToRegister(rl.emacsForwardWord(tokeniseLine))
		// vi delete, emacs forward, funny huh
		rl.viDeleteByAdjust(rl.emacsForwardWord(tokeniseLine))
		rl.updateHelpers()

	case seqAltDelete:
		if rl.modeTabCompletion {
			rl.resetVirtualComp(false)
		}
		rl.saveToRegister(-rl.emacsBackwardWord(tokeniseLine))
		rl.viDeleteByAdjust(-rl.emacsBackwardWord(tokeniseLine))
		rl.updateHelpers()

	default:
		if rl.modeTabFind {
			return
		}
		// alt+numeric append / delete
		if len(r) == 2 && '1' <= r[1] && r[1] <= '9' {
			if rl.modeViMode == VimDelete {
				rl.viDelete(r[1])
				return
			}

			line, err := rl.mainHistory.GetLine(rl.mainHistory.Len() - 1)
			if err != nil {
				return
			}
			if !rl.mainHist && rl.altHistory != nil {
				line, err = rl.altHistory.GetLine(rl.altHistory.Len() - 1)
				if err != nil {
					return
				}
			}

			tokens, _, _ := tokeniseSplitSpaces([]rune(line), 0)
			pos := int(r[1]) - 48 // convert ASCII to integer
			if pos > len(tokens) {
				return
			}
			rl.insert([]rune(tokens[pos-1]))
		} else {
			rl.viUndoSkipAppend = true
		}
	}
}

func (rl *Readline) carridgeReturn() {
	rl.moveCursorByAdjust(len(rl.line))
	rl.updateHelpers()
	rl.clearHelpers()
	print("\r\n")
	if rl.HistoryAutoWrite {
		var err error

		// Main history
		if rl.mainHistory != nil {
			rl.histPos, err = rl.mainHistory.Write(string(rl.line))
			if err != nil {
				print(err.Error() + "\r\n")
			}
		}
		// Alternative history
		if rl.altHistory != nil {
			rl.histPos, err = rl.altHistory.Write(string(rl.line))
			if err != nil {
				print(err.Error() + "\r\n")
			}
		}
	}
}

func isMultiline(r []rune) bool {
	for i := range r {
		if (r[i] == '\r' || r[i] == '\n') && i != len(r)-1 {
			return true
		}
	}
	return false
}

func (rl *Readline) allowMultiline(data []byte) bool {
	rl.clearHelpers()
	printf("\r\nWARNING: %d bytes of multiline data was dumped into the shell!", len(data))
	for {
		print("\r\nDo you wish to proceed (yes|no|preview)? [y/n/p] ")

		b := make([]byte, 1024)

		i, err := os.Stdin.Read(b)
		if err != nil {
			return false
		}

		s := string(b[:i])
		print(s)

		switch s {
		case "y", "Y":
			print("\r\n" + rl.mainPrompt)
			return true

		case "n", "N":
			print("\r\n" + rl.mainPrompt)
			return false

		case "p", "P":
			preview := string(bytes.Replace(data, []byte{'\r'}, []byte{'\r', '\n'}, -1))
			if rl.SyntaxHighlighter != nil {
				preview = rl.SyntaxHighlighter([]rune(preview))
			}
			print("\r\n" + preview)

		default:
			print("\r\nInvalid response. Please answer `y` (yes), `n` (no) or `p` (preview)")
		}
	}
}
