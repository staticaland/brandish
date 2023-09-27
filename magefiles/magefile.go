//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Build the project.
func Build() error {
	return sh.Run("go", "build", "-o", "main", "main.go")
}

// Run the project.
func Run() error {
	return sh.Run("go", "run", "main.go")
}
