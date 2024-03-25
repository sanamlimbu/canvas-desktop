package main

import (
	"bufio"
	"canvas-desktop/canvas"
	"embed"
	"net/http"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed .env
var config embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	env, err := readEnv(config)
	if err != nil {
		panic(err)
	}

	baseURL := env["CANVAS_BASE_URL"]
	if baseURL == "" {
		panic("no canvas base url provided")
	}

	accessToken := env["CANVAS_ACCESS_TOKEN"]
	if accessToken == "" {
		panic("no canvas access token provided")
	}

	pageSizeStr := env["CANVAS_PAGE_SIZE"]
	if pageSizeStr == "" {
		panic("no canvas page size provided")
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		panic(err)
	}

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

func readEnv(config embed.FS) (map[string]string, error) {
	env := make(map[string]string)
	file, err := config.Open(".env")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		arr := strings.Split(scanner.Text(), "=")
		env[arr[0]] = arr[1]
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return env, nil
}
