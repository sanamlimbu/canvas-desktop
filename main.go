package main

import (
	"canvas-desktop/canvas"
	"embed"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"golang.org/x/time/rate"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	baseURL := getenv("CANVAS_BASE_URL", "https://skillsaustralia.instructure.com")
	accessToken := getenv("CANVAS_ACCESS_TOKEN", "")
	pageSizeStr := getenv("CANVAS_PAGE_SIZE", "50")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		panic(err.Error)
	}

	rl := rate.NewLimiter(rate.Every(10*time.Second), 100) // 100 requests every 10 seconds
	client := canvas.NewAPIClient(baseURL, accessToken, pageSize, http.DefaultClient, rl)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Canvas",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			client,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func getenv(key string, other string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return other
	}

	return value
}
