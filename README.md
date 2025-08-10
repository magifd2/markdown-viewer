# Markdown Viewer

A simple, single-binary web server that acts as a file browser and Markdown renderer.

## Features

- **2-Pane View:** A persistent, expandable file tree on the left and a content viewer on the right.
- **Single Binary:** All UI assets (CSS, JavaScript, HTML templates) are embedded directly into the executable; no installation or external dependencies required.
- **Local First:** Serves files and directories from the current location where you run it, with a self-contained UI.
- **Markdown Rendering:** Renders Markdown files with GitHub-like styling, built-in syntax highlighting, and bundled Mermaid diagram support.
- **Secure by Default:** Includes protection against directory traversal attacks and sanitizes HTML to prevent Cross-Site Scripting (XSS).
- **Graceful Shutdown:** Shuts down safely via a UI button or `Ctrl+C`.

## Usage

1.  Download the latest `mdv` binary for your operating system from the [Releases](https://github.com/magifd2/markdown-viewer/releases) page.
2.  Run the executable from your terminal:
    ```bash
    ./mdv
    ```
4.  Open your web browser and navigate to `http://127.0.0.1:8888` (or the port you specify).

### Configuration

`mdv` can be configured using a `config.json` file or command-line flags. The configuration is loaded in the following order of precedence (later ones override earlier ones):

1.  Default values
2.  `config.json` in `$HOME/.config/mdv/`
3.  `config.json` in the current working directory
4.  Environment variables (prefixed with `MDV_`, e.g., `MDV_PORT=9000`)
5.  Command-line flags

For an example configuration file, see `config.json.example`.

### Command-Line Options

- `-p, --port int`: Specifies the port to listen on (default 8888).
- `-o, --open`: Automatically opens the default web browser to the application URL upon startup.
- `-d, --dir string`: Specifies the root directory to serve (default ".").
- `-h, --help`: Displays the help message.
- `--version`: Displays the application version.

## Build from Source

To build the project from source, you need to have Go and Make installed.

1.  Clone the repository:
    ```bash
    git clone https://github.com/magifd2/markdown-viewer.git
    cd markdown-viewer
    ```
2.  Build the `mdv` binary:
    ```bash
    make build
    ```
    The `mdv` binary will be generated in the `bin/<OS>-<ARCH>/` directory (e.g., `bin/darwin-arm64/mdv`).
3.  Run the application:
    ```bash
    make run
    ```
    This will build and then run the `mdv` application.

For cross-platform builds, see the `Makefile` for more options (`make help`).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

### Third-Party Licenses

This project bundles third-party software components. Their licenses and attribution notices are available in the `NOTICE.md` file included in the distribution.
