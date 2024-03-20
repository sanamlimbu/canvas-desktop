package main

import (
	"canvas-desktop/canvas"
	"embed"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	baseURL := getenv("CANVAS_BASE_URL", "https://skillsaustralia.instructure.com/api/v1")
	accessToken := getenv("CANVAS_ACCESS_TOKEN", "")
	pageSizeStr := getenv("CANVAS_PAGE_SIZE", "100")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		println("Error:", err.Error())
	}

	if accessToken == "" {
		println("Error:", fmt.Errorf("missing access token"))
	}

	//rl := rate.NewLimiter(rate.Every(10*time.Second), 100)                           // 100 requests every 10 seconds
	client := canvas.NewAPIClient(baseURL, accessToken, pageSize, http.DefaultClient)
	controller := canvas.NewController(client)

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
			controller,
		},
		EnumBind: []interface {
		}{
			canvas.AllAssignmentBucket,
			canvas.AllCourseEnrollmentType,
			canvas.AllEnrollmentType,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func getenv(key string, other string) string {
	if other != "" {
		return other
	}

	value := os.Getenv(key)
	return value
}
