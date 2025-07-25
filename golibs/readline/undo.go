package readline

type undoItem struct {
	line string
	pos  int
}

func (rl *Readline) undoAppendHistory() {
	defer func() { rl.viUndoSkipAppend = false }()

	if rl.viUndoSkipAppend {
		return
	}

	rl.viUndoHistory = append(rl.viUndoHistory, undoItem{
		line: string(rl.line),
		pos:  rl.pos,
	})
}

func (rl *Readline) undoLast() {
	var undo undoItem
	for {
		if len(rl.viUndoHistory) == 0 {
			return
		}
		undo = rl.viUndoHistory[len(rl.viUndoHistory)-1]
		rl.viUndoHistory = rl.viUndoHistory[:len(rl.viUndoHistory)-1]
		if string(undo.line) != string(rl.line) {
			break
		}
	}

	rl.line = []rune(undo.line)
	rl.pos = undo.pos

	if rl.modeViMode != VimInsert && len(rl.line) > 0 && rl.pos == len(rl.line) {
		rl.pos--
	}

	rl.updateHelpers()
}
