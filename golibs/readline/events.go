package readline

// EventReturn is a structure returned by the callback event function.
// This is used by readline to determine what state the API should
// return to after the readline event.
type EventReturn struct {
	ForwardKey    bool
	ClearHelpers  bool
	CloseReadline bool
	InfoText      []rune
	NewLine       []rune
	NewPos        int
}

// AddEvent registers a new keypress handler
func (rl *Readline) AddEvent(keyPress string, callback func(string, []rune, int) *EventReturn) {
	rl.evtKeyPress[keyPress] = callback
}

// DelEvent deregisters an existing keypress handler
func (rl *Readline) DelEvent(keyPress string) {
	delete(rl.evtKeyPress, keyPress)
}
