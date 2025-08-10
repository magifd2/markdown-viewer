# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2025-08-10

### Security
- **Resolved GO-2025-3595 Vulnerability:** Updated `golang.org/x/net` from v0.26.0 to v0.38.0 to address a vulnerability related to incorrect neutralization of input during web page generation.

### Added
- **Local Bundling of Client-Side Assets:** `highlight.js` and `Mermaid.js` are now bundled directly with the application, removing external CDN dependencies for offline and closed-network environments.
- **Third-Party License Notice:** Included `NOTICE.md` in release packages for compliance with third-party software licenses.

### Changed
- **Updated Documentation:** `DEVELOPMENT.md` and `README.md` have been updated to reflect local asset bundling and license information.


## [1.0.0] - 2025-08-10

This is the first stable release after a major security and functionality overhaul.

### Security
- **Prevented Directory Traversal:** Implemented a custom router to validate request paths, blocking access to files outside the intended directory.
- **Prevented Cross-Site Scripting (XSS):** Integrated the `bluemonday` library to sanitize all HTML rendered from Markdown, mitigating XSS risks.
- **Hardened Browser Auto-Open:** Re-implemented the browser auto-open feature with strict URL validation and command path verification to prevent command injection vulnerabilities.
- **Disabled Unsafe Links:** The Markdown renderer now disables links to local non-Markdown files and requires user confirmation before opening external links in a new tab.
- **Updated Dependencies:** Updated all dependencies to their latest versions to patch known vulnerabilities.
- **Addressed `gosec` Findings:** Fixed all issues reported by the `gosec` static analysis tool, including potential Slowloris attacks and unhandled errors.

### Fixed
- **Markdown Rendering Engine:** Replaced the `blackfriday` parser with `goldmark` and its GFM extension. This fixed numerous rendering bugs, including incorrect handling of code blocks within lists, unwanted line breaks, and improved standards compliance.
- **Shutdown Mechanism:** Fixed a bug where the shutdown button was not working due to an uninitialized channel.

### Changed
- **Improved Layout:** Increased the maximum width of the Markdown rendering area to 1140px for better readability on wider screens.
