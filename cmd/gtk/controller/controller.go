//go:build gtk

package controller

// Controller represents a UI controller that can be executed.
// All controllers should have a Run method to execute their primary function.

type PasswordProvider func() (string, bool)
