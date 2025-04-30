package uuid

import (
	"github.com/google/uuid"
)

func NewGenerator() UUID {
	return &Generator{}
}

type Generator struct{}

func (Generator) New() string { return uuid.New().String() }
