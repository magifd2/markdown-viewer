# Development Plan & Project Documentation

This document outlines the development plan, architecture, and technical decisions for the Markdown Viewer project. It serves as a guide for current and future development.

## 1. Project Vision & Goals

The primary goal is to create a lightweight, portable, and easy-to-use tool for browsing and reading Markdown files locally. It should be distributed as a single binary to eliminate the need for complex setup or dependencies for the end-user.

## 2. Core Principles

- **Simplicity:** Prioritize simple, clear code over complex abstractions.
- **Standard Library First:** Rely on Go's standard library (`net/http`) for core web functionalities to maintain a small footprint.
- **Portability:** The final product must be a single, cross-platform binary.
- **Security:** While serving local files, be mindful of security best practices (e.g., preventing directory traversal attacks).

## 3. Architecture & Technology Stack

- **Backend Language:** Go
- **Web Server:** Standard `net/http` package.
- **Markdown Parsing:** [blackfriday v2](https://github.com/russross/blackfriday)
    - **Reasoning:** Initially, `goldmark` was chosen, but persistent issues with fetching its dependencies (`go get`) in the development environment necessitated a switch. `blackfriday` is a robust, widely-used, pure Go alternative that does not have these dependency issues.
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

### **Phase 5: UI/UX Improvements & Finalization**

- [ ] Apply CSS to improve the visual appearance of the file list and rendered Markdown.
- [ ] Implement breadcrumb navigation to easily move between directories.
- [ ] Refine error handling and provide user-friendly error pages.
- [ ] Create a build script or Makefile for easy compilation.
- [ ] Write final usage instructions in the `README.md`.