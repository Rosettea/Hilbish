//go:build plan9
// +build plan9

package readline

import "errors"

func (rl *Readline) launchEditor(multiline []rune) ([]rune, error) {
	return rl.line, errors.New("Not currently supported on Plan 9")
}
