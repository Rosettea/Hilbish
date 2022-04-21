package readline

import (
	"strings"
)

// insertCandidateVirtual - When a completion candidate is selected, we insert it virtually in the input line:
// this will not trigger further firltering against the other candidates. Each time this function
// is called, any previous candidate is dropped, after being used for moving the cursor around.
func (rl *Instance) insertCandidateVirtual(candidate []rune) {
	for {
		// I don't really understand why `0` is creaping in at the end of the
		// array but it only happens with unicode characters.
		if len(candidate) > 1 && candidate[len(candidate)-1] == 0 {
			candidate = candidate[:len(candidate)-1]
			continue
		}
		break
	}

	// Keep reference of previous candidate to compare its prefix
	oldComp := rl.currentComp
	// We place the cursor back at the beginning of the previous virtual candidate
	rl.pos -= len(rl.oldCOmp)

	// We delete the previous virtual completion, just
	// like we would delete a word in vim editing mode.
	if len(oldComp) == 1 {
		rl.deleteVirtual() // Delete a single character
	} else if len(oldComp) > 0 {
		rl.viDeleteByAdjustVirtual(rl.viJumpEVirtual(tokeniseSplitSpaces) + 1)
	}
	prefix := len(rl.tcPrefix)
	comp := candidate[prefix:]

	// We then keep a reference to the new candidate
	rl.currentComp = comp

	// We should not have a remaining virtual completion
	// line, so it is now identical to the real line.
	rl.lineComp = rl.line

	compPrefix := candidate[:prefix]
	// first clause checks if the previous input wouldve changed the prefix
	if (len(oldComp) > 0 && string(oldComp[:prefix]) != rl.tcPrefix) || string(compPrefix) != rl.tcPrefix {
		rl.viDeleteByAdjustVirtual(-prefix)
		rl.insertVirtual(compPrefix)
	}

	rl.insertVirtual(comp)
}

func (rl *Instance) insertVirtual(r []rune) {
	// we dont need that scuffed unicode fix here, since this function
	// is only called in the 2 instances above, which already handles it
	switch {
	case len(rl.lineComp) == 0:
		rl.lineComp = r
	case rl.pos == 0:
		rl.lineComp = append(r, rl.lineComp...)
	case rl.pos < len(rl.lineComp):
		r := append(r, rl.lineComp[rl.pos:]...)
		rl.lineComp = append(rl.lineComp[:rl.pos], r...)
	default:
		rl.lineComp = append(rl.lineComp, r...)
	}

	// We place the cursor at the end of our new virtually completed item
	rl.pos += len(r)
}

// Insert the current completion candidate into the input line.
// This candidate might either be the currently selected one (white frame),
// or the only candidate available, if the total number of candidates is 1.
func (rl *Instance) insertCandidate() {

	cur := rl.getCurrentGroup()

	if cur != nil {
		completion := cur.getCurrentCell(rl)
		prefix := len(rl.tcPrefix)

		// Special case for the only special escape, which
		// if not handled, will make us insert the first
		// character of our actual rl.tcPrefix in the candidate.
		if strings.HasPrefix(string(rl.tcPrefix), "%") {
			prefix++
		}

		// Ensure no indexing error happens with prefix
		if len(completion) >= prefix {
			comp := completion[prefix:]
			if completion[:prefix] != rl.tcPrefix {
				rl.viDeleteByAdjust(-prefix)
				comp = completion
			}
			rl.insert([]rune(comp))
			if !cur.TrimSlash && !cur.NoSpace {
				rl.insert([]rune(" "))
			}
		}
	}
}

// updateVirtualComp - Either insert the current completion candidate virtually, or on the real line.
func (rl *Instance) updateVirtualComp() {
	cur := rl.getCurrentGroup()
	if cur != nil {

		completion := cur.getCurrentCell(rl)
		prefix := len(rl.tcPrefix)

		// If the total number of completions is one, automatically insert it.
		if rl.hasOneCandidate() {
			rl.insertCandidate()
			// Quit the tab completion mode to avoid asking to the user to press
			// Enter twice to actually run the command
			// Refresh first, and then quit the completion mode
			rl.viUndoSkipAppend = true
			rl.resetTabCompletion()
		} else {

			// Special case for the only special escape, which
			// if not handled, will make us insert the first
			// character of our actual rl.tcPrefix in the candidate.
			if strings.HasPrefix(string(rl.tcPrefix), "%") {
				prefix++
			}

			// Or insert it virtually.
			if len(completion) >= prefix {
				rl.insertCandidateVirtual([]rune(completion))
			}
		}
	}
}

// resetVirtualComp - This function is called before most of our readline key handlers,
// and makes sure that the current completion (virtually inserted) is either inserted or dropped,
// and that all related parameters are reinitialized.
func (rl *Instance) resetVirtualComp(drop bool) {

	// If we don't have a current virtual completion, there's nothing to do.
	// IMPORTANT: this MUST be first, to avoid nil problems with empty comps.
	if len(rl.currentComp) == 0 {
		return
	}

	// Get the current candidate and its group.
	//It contains info on how we must process it
	cur := rl.getCurrentGroup()
	if cur == nil {
		return
	}
	completion := cur.getCurrentCell(rl)
	// Avoid problems with empty completions
	if completion == "" {
		rl.clearVirtualComp()
		return
	}

	// We will only insert the net difference between prefix and completion.
	prefix := len(rl.tcPrefix)
	// Special case for the only special escape, which
	// if not handled, will make us insert the first
	// character of our actual rl.tcPrefix in the candidate.
	if strings.HasPrefix(string(rl.tcPrefix), "%") {
		prefix++
	}

	// If we are asked to drop the completion, move it away from the line and return.
	if drop {
		rl.pos -= len([]rune(completion[prefix:]))
		rl.lineComp = rl.line
		rl.clearVirtualComp()
		return
	}

	// Insert the current candidate. A bit of processing happens:
	// - We trim the trailing slash if needed
	// - We add a space only if the group has not explicitely specified not to add one.
	if cur.TrimSlash {
		trimmed, hadSlash := trimTrailing(completion)
		if !hadSlash && rl.compAddSpace {
			if !cur.NoSpace {
				trimmed = trimmed + " "
			}
		}
		rl.insertCandidateVirtual([]rune(trimmed))
	} else {
		if rl.compAddSpace {
			if !cur.NoSpace {
				completion = completion + " "
			}
		}
		rl.insertCandidateVirtual([]rune(completion))
	}

	// Reset virtual
	rl.clearVirtualComp()
}

// trimTrailing - When the group to which the current candidate
// belongs has TrimSlash enabled, we process the candidate.
// This is used when the completions are directory/file paths.
func trimTrailing(comp string) (trimmed string, hadSlash bool) {
	// Unix paths
	if strings.HasSuffix(comp, "/") {
		return strings.TrimSuffix(comp, "/"), true
	}
	// Windows paths
	if strings.HasSuffix(comp, "\\") {
		return strings.TrimSuffix(comp, "\\"), true
	}
	return comp, false
}

// viDeleteByAdjustVirtual - Same as viDeleteByAdjust, but for our virtually completed input line.
func (rl *Instance) viDeleteByAdjustVirtual(adjust int) {
	var (
		newLine []rune
		backOne bool
	)

	// Avoid doing anything if input line is empty.
	if len(rl.lineComp) == 0 {
		return
	}

	switch {
	case adjust == 0:
		rl.viUndoSkipAppend = true
		return
	case rl.pos+adjust == len(rl.lineComp)-1:
		newLine = rl.lineComp[:rl.pos]
		// backOne = true // Deleted, otherwise the completion moves back when we don't want to.
	case rl.pos+adjust == 0:
		newLine = rl.lineComp[rl.pos:]
	case adjust < 0:
		newLine = append(rl.lineComp[:rl.pos+adjust], rl.lineComp[rl.pos:]...)
	default:
		newLine = append(rl.lineComp[:rl.pos], rl.lineComp[rl.pos+adjust:]...)
	}

	// We have our new line completed
	rl.lineComp = newLine

	if adjust < 0 {
		rl.moveCursorByAdjust(adjust)
	}

	if backOne {
		rl.pos--
	}
}

// viJumpEVirtual - Same as viJumpE, but for our virtually completed input line.
func (rl *Instance) viJumpEVirtual(tokeniser func([]rune, int) ([]string, int, int)) (adjust int) {
	split, index, pos := tokeniser(rl.lineComp, rl.pos)
	if len(split) == 0 {
		return
	}

	word := rTrimWhiteSpace(split[index])

	switch {
	case len(split) == 0:
		return
	case index == len(split)-1 && pos >= len(word)-1:
		return
	case pos >= len(word)-1:
		word = rTrimWhiteSpace(split[index+1])
		adjust = len(split[index]) - pos
		adjust += len(word) - 1
	default:
		adjust = len(word) - pos - 1
	}
	return
}

func (rl *Instance) deleteVirtual() {
	switch {
	case len(rl.lineComp) == 0:
		return
	case rl.pos == 0:
		rl.lineComp = rl.lineComp[1:]
	case rl.pos > len(rl.lineComp):
	case rl.pos == len(rl.lineComp):
		rl.lineComp = rl.lineComp[:rl.pos]
	default:
		rl.lineComp = append(rl.lineComp[:rl.pos], rl.lineComp[rl.pos+1:]...)
	}
	
	rl.pos--
}

// We are done with the current virtual completion candidate.
// Get ready for the next one
func (rl *Instance) clearVirtualComp() {
	rl.line = rl.lineComp
	rl.currentComp = []rune{}
	rl.compAddSpace = false
}
