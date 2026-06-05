//go:build !wails
// +build !wails

package main

import (
	"fmt"
	"skillsmaster-win/models"
)

func main() {
	resolver := models.NewDefaultPathResolver()
	app := NewApp(resolver)
	agents := app.GetAgents()
	fmt.Printf("SkillsMaster loaded with %d agents\n", len(agents))
}
