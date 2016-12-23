package api

import (
	"fmt"
)

type Sample struct {
	Name string
}

func New(name string) (sample * Sample) {
	return &Sample{
		Name: name,
	}
}

func (sample * Sample) Print() {
	fmt.Println("Sample Name:", sample.Name)
}