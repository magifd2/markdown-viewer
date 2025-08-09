# Markdown Viewer

A simple, single-binary web server that acts as a file browser and Markdown renderer.

## Features

- **2-Pane View:** A persistent, expandable file tree on the left and a content viewer on the right.
- **Single Binary:** No installation or external dependencies required.
- **Local First:** Serves files and directories from the current location where you run it.
- **Markdown Rendering:** Renders Markdown files with GitHub-like styling, syntax highlighting, and Mermaid diagram support.
- **Graceful Shutdown:** Shuts down safely via a UI button or `Ctrl+C`.

## Usage

1.  Download the latest binary for your operating system from the [Releases](https://github.com/magifd2/markdown-viewer/releases) page.
2.  Place the binary in the directory you want to browse.
3.  Run the executable from your terminal:
    ```bash
    ./markdown-viewer
    ```
4.  Open your web browser and navigate to `http://127.0.0.1:8080`.

## Build from Source

To build the project from source, you need to have Go and Make installed.

1.  Clone the repository:
    ```bash
    git clone https://github.com/magifd2/markdown-viewer.git
    cd markdown-viewer
    ```
2.  Build the binary:
    ```bash
    make build
    ```
3.  Run the application:
    ```bash
    make run
    ```

For cross-platform builds, see the `Makefile` for more options (`make help`).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.