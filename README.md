# Canvas LMS Desktop Client

## About

This is a desktop client for Canvas LMS built using the Wails React-TS template.
To run the app, you need to set the `CANVAS_ACCESS_TOKEN` environment variable. Obtain your access token from Canvas and set it as an environment variable before running the application.

## Development

Development dependencies

- Go 1.18+
- NPM (Node 15+)
- CANVAS_ACCESS_TOKEN (get access token from Canvas and set as environment variable)

You can configure the project by editing `wails.json`. For more detailed information about the project settings, please refer to the [Wails Project Configuration Documentation](https://wails.io/docs/reference/project-config).
To run in live development mode, execute `wails dev` in the project directory. This will start a Vite development server, enabling fast hot reload of your frontend changes. If you prefer to develop in a browser and access your Go methods, there is a dev server running on http://localhost:34115. Connect to this address in your browser to call your Go code from devtools.

## Building

To create a redistributable, production-ready package, use `wails build`. This will generate the necessary files for distribution.
