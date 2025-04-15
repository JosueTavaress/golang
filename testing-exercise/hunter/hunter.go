package hunter

import (
	"testdoubles/prey"
)

// Hunter is an interface that represents a hunter
type Hunter interface {
	// Hunt hunts the prey
	Hunt(prey prey.Prey) (err error)
}