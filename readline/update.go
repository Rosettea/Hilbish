package readline

// updateHelpers is a key part of the whole refresh process:
// it should coordinate reprinting the input line, any Infos and completions
// and manage to get back to the current (computed) cursor coordinates
func (rl *Instance) updateHelpers() {

	// Load all Infos & completions before anything.
	// Thus overwrites anything having been dirtily added/forced/modified, like rl.SetInfoText()
	rl.getInfoText()
	rl.getHintText()
	if rl.modeTabCompletion {
		rl.getTabCompletion()
	}

	// We clear everything
	rl.clearHelpers()

	// We are at the prompt line (with the latter
	// not printed yet), then reprint everything
	rl.renderHelpers()
}

// Update reference should be called only once in a "loop" (not Readline(), but key control loop)
func (rl *Instance) updateReferences() {

	// We always need to work with clean data,
	// since we will have incrementers all around
	rl.posX = 0
	rl.fullX = 0
	rl.posY = 0
	rl.fullY = 0

	var fullLine, cPosLine int
	if len(rl.currentComp) > 0 {
		fullLine = len(rl.lineComp)
		cPosLine = len(rl.lineComp[:rl.pos])
	} else {
		fullLine = len(rl.line)
		cPosLine = len(rl.line[:rl.pos])
	}

	// We need the X offset of the whole line
	toEndLine := rl.promptLen + fullLine
	fullOffset := toEndLine / GetTermWidth()
	rl.fullY = fullOffset
	fullRest := toEndLine % GetTermWidth()
	rl.fullX = fullRest

	// Use rl.pos value to get the offset to go TO/FROM the CURRENT POSITION
	lineToCursorPos := rl.promptLen + cPosLine
	offsetToCursor := lineToCursorPos / GetTermWidth()
	cPosRest := lineToCursorPos % GetTermWidth()

	// If we are at the end of line
	if fullLine == rl.pos {
		rl.posY = fullOffset

		if fullRest == 0 {
			rl.posX = 0
		} else if fullRest > 0 {
			rl.posX = fullRest
		}
	} else if rl.pos < fullLine {
		// If we are somewhere in the middle of the line
		rl.posY = offsetToCursor

		if cPosRest == 0 {
		} else if cPosRest > 0 {
			rl.posX = cPosRest
		}
	}
}

func (rl *Instance) resetHelpers() {
	rl.modeAutoFind = false

	// Now reset all below-input helpers
	rl.resetInfoText()
	rl.resetTabCompletion()
}

// clearHelpers - Clears everything: prompt, input, Infos & comps,
// and comes back at the prompt.
func (rl *Instance) clearHelpers() {

	// Now go down to the last line of input
	moveCursorDown(rl.fullY - rl.posY)
	moveCursorBackwards(rl.posX)
	moveCursorForwards(rl.fullX)

	// Clear everything below
	print(seqClearScreenBelow)

	// Go back to current cursor position
	moveCursorBackwards(GetTermWidth())
	moveCursorUp(rl.fullY - rl.posY)
	moveCursorForwards(rl.posX)
}

// renderHelpers - pritns all components (prompt, line, Infos & comps)
// and replaces the cursor to its current position. This function never
// computes or refreshes any value, except from inside the echo function.
func (rl *Instance) renderHelpers() {
	
	// when the instance is in this state we want it to be "below" the user's
	// input for it to be aligned properly
	if !rl.compConfirmWait {
		rl.writeHintText()
	}
	rl.echo()
	if rl.modeTabCompletion {
		// in tab complete mode we want it to update
		// when something has been selected
		// (dynamic!!)
		rl.getHintText()
		rl.writeHintText()
	} else if !rl.compConfirmWait {
		// for the same reason above, do nothing here
	} else {
		rl.writeHintText()
	}

	// Go at beginning of first line after input remainder
	moveCursorDown(rl.fullY - rl.posY)
	moveCursorBackwards(GetTermWidth())

	// Print Infos, check for any confirmation Info current.
	// (do not overwrite the confirmation question Info)
	if !rl.compConfirmWait {
		if len(rl.infoText) > 0 {
			print("\n")
		}
		rl.writeInfoText()
		moveCursorBackwards(GetTermWidth())

		// Print completions and go back to beginning of this line
		print("\n")
		rl.writeTabCompletion()
		moveCursorBackwards(GetTermWidth())
		moveCursorUp(rl.tcUsedY)
	}

	// If we are still waiting for the user to confirm too long completions
	// Immediately refresh the Infos
	if rl.compConfirmWait {
		print("\n")
		rl.writeInfoText()
		rl.getInfoText()
		moveCursorBackwards(GetTermWidth())
	}

	// Anyway, compensate for Info printout
	if len(rl.infoText) > 0 {
		moveCursorUp(rl.infoY)
	} else if !rl.compConfirmWait {
		moveCursorUp(1)
	} else if rl.compConfirmWait {
		moveCursorUp(1)
	}

	// Go back to current cursor position
	moveCursorUp(rl.fullY - rl.posY)
	moveCursorForwards(rl.posX)
}
