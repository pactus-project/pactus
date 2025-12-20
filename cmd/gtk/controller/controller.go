//go:build gtk

package controller

// Controller represents a UI controller that can be executed.
// All controllers should have a Run method to execute their primary function.
// The Run method signature may vary based on the controller's needs:
// - Run() - for controllers that need no parameters
// - Run(addr string) - for controllers that need an address parameter
// - Run() (string, bool) - for controllers that return user input

type PasswordProvider func() (string, bool)
