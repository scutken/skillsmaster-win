//go:build wails
// +build wails

package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"

	"skillsmaster-win/models"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Set up logging to file so we can see errors
	logFile, err := os.Create("skillsmaster-debug.log")
	if err == nil {
		log.SetOutput(logFile)
		defer logFile.Close()
	}

	log.Println("SkillsMaster starting...")
	wd, _ := os.Getwd()
	log.Printf("Working directory: %s", wd)

	resolver := models.NewDefaultPathResolver()
	app := NewApp(resolver)
	agents := app.GetAgents()
	log.Printf("Loaded %d agents", len(agents))

	fmt.Printf("SkillsMaster starting with %d agents...\n", len(agents))

	err = wails.Run(&options.App{
		Title:     "SkillsMaster",
		Width:     1024,
		Height:    768,
		MinWidth:  800,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			log.Println("Window started successfully")
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		errMsg := fmt.Sprintf("SkillsMaster error: %v", err)
		log.Println(errMsg)
		fmt.Fprintln(os.Stderr, errMsg)
	}
	log.Println("SkillsMaster exiting")
}
