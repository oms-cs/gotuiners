# gotuiners
TUI for docker containers in go

A lightweight, terminal-based dashboard for managing Docker containers and images, built with Go and the [Charm Bubble Tea](https://github.com/charmbracelet/bubbletea) framework.

![gotuiners Demo](assets/tui.gif)

## 🚀 Features
* **Real-time Monitoring:** View running containers and local images in a side-by-side layout.
* **Interactive TUI:** Use keyboard shortcuts to navigate and inspect resources.
* **Logs & Inspect:** Instant access to container logs and image inspection via a scrollable viewport.
* **Keyboard Driven:** Optimized for speed—no mouse required.

## ⌨️ Keyboard Shortcuts
| Key | Action |
| :--- | :--- |
| `Tab` | Switch focus between Containers, Images, and Details |
| `Enter` | Select an item to view Logs/Details |
| `↑/↓` | Navigate through lists |
| `j/k` | Scroll logs in the Details panel |
| `q` / `Ctrl+C` | Quit Application |

## 🛠️ Installation

Ensure you have [Go](https://go.dev/) and [Docker](https://www.docker.com/) installed.

```bash
# Clone the repository
git clone [https://github.com/oms-cs/gotuiners](https://github.com/oms-cs/gotuiners)

# Navigate to the project
cd gotuiners

# Install dependencies
go mod tidy

# Run the application
go run main.go