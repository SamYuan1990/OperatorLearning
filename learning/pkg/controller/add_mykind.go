package controller

import (
	"github.com/SamYuan1990/OperatorLearning/learning/pkg/controller/mykind"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, mykind.Add)
}
