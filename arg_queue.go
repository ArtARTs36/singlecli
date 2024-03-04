package cli

import (
	"slices"
)

type argQueue struct {
	args []*ArgDefinition
}

func newArgQueue(args []*ArgDefinition) *argQueue {
	return &argQueue{args: slices.Clone(args)}
}

func (s *argQueue) isEmpty() bool {
	return len(s.args) == 0
}

func (s *argQueue) valid() bool {
	return len(s.args) > 0
}

func (s *argQueue) pop() *ArgDefinition {
	if s.isEmpty() {
		return nil
	}

	a := s.args[0]
	s.args = s.args[1:]

	return a
}

func (s *argQueue) firstRequired() *ArgDefinition {
	for _, def := range s.args {
		if def.Required {
			return def
		}
	}

	return nil
}

func (s *argQueue) clean() {
	s.args = nil
}
