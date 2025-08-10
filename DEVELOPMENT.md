# Development Plan & Project Documentation

This document outlines the development plan, architecture, and technical decisions for the Markdown Viewer project. It serves as a guide for current and future development.

## 1. Project Vision & Goals

The primary goal is to create a lightweight, portable, and easy-to-use tool for browsing and reading Markdown files locally. It should be distributed as a single binary to eliminate the need for complex setup or dependencies for the end-user.

## 2. Core Principles

- **Simplicity:** Prioritize simple, clear code over complex abstractions.
- **Standard Library First:** Rely on Go's standard library (`net/http`) for core web functionalities to maintain a small footprint.
- **Portability:** The final product must be a single, cross-platform binary.
- **Security First:** Actively prevent vulnerabilities. All features and dependencies must be reviewed for security implications.

## 3. Architecture & Technology Stack

- **Backend Language:** Go
- **Web Server:** Standard `net/http` package.
- **Markdown Parsing:** [blackfriday v2](https://github.com/russross/blackfriday)
    - **Reasoning:** Initially, `goldmark` was chosen, but persistent issues with fetching its dependencies (`go get`) in the development environment necessitated a switch. `blackfriday` is a robust, widely-used, pure Go alternative that does not have these dependency issues.
- **HTML Sanitization:** [bluemonday](https://github.com/microcosm-cc/bluemonday)
    - **Reasoning:** Introduced to prevent XSS attacks from malicious Markdown content. It sanitizes the HTML output from `blackfriday` before rendering.
- **Syntax Highlighting:** [highlight.js](https://highlightjs.org/)
    - **Reasoning:** Switched from a Go-based highlighter to a client-side library to decouple it from the backend Markdown parser. `highlight.js` is powerful and supports a vast number of languages.
- **Diagram Rendering:** [Mermaid.js](https://mermaid-js.github.io/mermaid/)
    - **Reasoning:** It's a widely used standard for creating diagrams from text and can be rendered entirely on the client-side, fitting our architecture perfectly.

## 4. Development Roadmap

This roadmap is broken down into phases to ensure iterative and manageable development.

### **Phase 1: Project Initialization (Completed)**

- [x] Initialize Go module (`go mod init markdown-viewer`)
- [x] Create initial `main.go` structure.
- [x] Create `.gitignore` file.
- [x] Create `README.md` and `DEVELOPMENT.md`.

### **Phase 2: Web Server and File Listing (Completed)**

- [x] Implement an HTTP handler to serve the current working directory.
- [x] The handler will list all files and subdirectories.
- [x] Differentiate between files and directories in the listing.
- [x] Create clickable links for all entries.

### **Phase 3: Basic Markdown Rendering (Completed)**

- [x] Add a new handler that is triggered when a `.md` or `.markdown` file is clicked.
- [x] Read the content of the selected Markdown file.
- [x] Integrate a Markdown library to parse the content into HTML.
- [x] Render the resulting HTML in a clean, readable template.

### **Phase 4: Advanced Rendering Features (Completed)**

- [x] Switched Markdown library to `blackfriday` to resolve dependency issues.
- [x] Implemented GFM-like features (tables, etc.).
- [x] Added client-side syntax highlighting with `highlight.js`.
- [x] Included Mermaid.js library and ensured ````mermaid` blocks are rendered.
- [x] Created a `tests` directory with a test file for verifying functionality.

### **Phase 5: Security Hardening (v0.1.1) (Completed)**

- [x] **Directory Traversal:** Replaced `http.ServeMux` with a custom router to validate request paths and prevent access to files outside the target directory.
- [x] **Cross-Site Scripting (XSS):** Introduced `bluemonday` to sanitize HTML generated from Markdown, preventing malicious script execution.
- **Dependency Vulnerabilities:** Updated dependencies to patch known vulnerabilities.
- **Code Hardening:** Addressed multiple issues identified by `gosec` (unhandled errors, missing timeouts, command injection risks).
- **Browser Auto-Open:** Re-enabled the feature with enhanced security validation.

### **Phase 6: UI/UX Improvements & Finalization**

- [ ] Apply CSS to improve the visual appearance of the file list and rendered Markdown.
- [ ] Implement breadcrumb navigation to easily move between directories.
- [ ] Refine error handling and provide user-friendly error pages.
- [ ] Create a build script or Makefile for easy compilation.
- [ ] Write final usage instructions in the `README.md`.

## 5. Configuration & CLI

This section details the application's configuration system and command-line interface (CLI) options.

### 5.1. Configuration Loading

The application uses `spf13/viper` for flexible configuration management. Settings are loaded from multiple sources in the following order of precedence (later ones override earlier ones):

1.  **Default values:** Hardcoded defaults within the application (e.g., `port: 8080`, `open: false`, `target_dir: .`).
2.  **Configuration files:**
    *   `config.json` in `$HOME/.config/mdv/`
    *   `config.json` in the current working directory
3.  **Environment variables:** Prefixed with `MDV_` (e.g., `MDV_PORT=9000`). Dot notation in config keys is replaced with underscores (e.g., `MDV_TARGET_DIR` for `target_dir`).
4.  **Command-line flags:** Overrides all other settings.

### 5.2. Command-Line Options

The application uses `spf13/cobra` for its command-line interface.

*   `-p <port>` or `--port <port>`: Specifies the port to listen on (e.g., `-p 9000`).
*   `-o` or `--open`: Automatically opens the default web browser to the application URL upon startup.
*   `-d <directory>` or `--dir <directory>`: Specifies the root directory to serve. If not set, the current working directory is used.
*   `-h` or `--help`: Displays the help message with all available options.

## 6. Build & Deployment

This section outlines the build process and deployment considerations for the Markdown Viewer.

### 6.1. Binary Naming and Output Paths

The application binary is named `mdv`. When built using `make build`, the executable is placed in a platform-specific subdirectory within the `bin/` directory (e.g., `bin/darwin-arm64/mdv` for macOS ARM64).

### 6.2. Build Commands

*   `make build`: Compiles the application for the current operating system and architecture.
*   `make run`: Builds and then runs the application.
*   `make clean`: Removes all build artifacts and the `bin/` directory.
*   `make cross-compile`: Builds binaries for all supported platforms (macOS, Linux, Windows).
*   `make package-all`: Packages all cross-compiled binaries into archives.

## Known Issues

There are several known issues related to how Markdown code blocks are rendered:

1.  **Incorrect Line Breaks:** Code within fenced code blocks (` ``` `) may not display with correct line breaks, appearing as a single continuous line.
2.  **Incorrect Indentation:** The indentation of elements immediately following a code block may be incorrect, appearing at the same level as the code block itself.
3.  **Language Specifier Displayed:** The language specifier (e.g., ` ```bash `) may be rendered as plain text within the code block, rather than being correctly interpreted for syntax highlighting.

These issues are believed to stem from the interaction between the `blackfriday` Markdown parser's HTML output, and the client-side `highlight.js` library. Further investigation into `blackfriday`'s rendering behavior and its compatibility with client-side libraries is required.
