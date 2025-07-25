package scan

import "fmt"

type State struct {
	Message string
}

func (s *State) Get() string {
	return fmt.Sprintf("lexer state: %s", s.Message)
}

func NewState(message string) *State {
	return &State{
		Message: message,
	}
}
