# 🗂️ Go Recursive Directory Watcher

A lightweight Go program that recursively watches a directory for **newly created files and/or directories** and executes a command for each match.

## ✨ Features

- **Recursive watching** — all subdirectories are watched from startup.
- **Automatic new directory watching** — newly created folders are added to the watch list automatically.
- **Optional directory triggering** — run the command for new directories as well as files.
- **Hidden file/directory ignoring** — skip hidden files and folders if desired.
- **Glob-based filtering** — include or exclude files/directories using `filepath.Match` patterns.
- **Cross-platform** — works on Linux, macOS, and Windows.

## 📦 Installation

Clone the repository and build:

```bash
git clone <your-repo-url>
cd <repo-name>
go build -o watcher
```

## ⚙️ Usage

```bash
watcher [--include-dirs] [--ignore-hidden] \
        [--include pattern] [--exclude pattern] \
        <directory_to_watch> <command> [args...]
```

## 🛠️ Options

| Flag              | Description                                                                                        |
| ----------------- | -------------------------------------------------------------------------------------------------- |
| `--include-dirs`  | Trigger the command for newly created directories as well as files.                                |
| `--ignore-hidden` | Ignore hidden files and directories (names starting with `.`).                                     |
| `--include`       | Glob pattern to include. Can be specified multiple times. Only matching paths trigger the command. |
| `--exclude`       | Glob pattern to exclude. Can be specified multiple times. Matching paths will be skipped.          |

Patterns are matched **relative to the watched root directory**, using `/` as the path separator.

## 💡 Examples

Watch `/tmp` and echo new files:

```bash
./watcher /tmp echo "Detected:"
```

Watch `/tmp`, including directories:

```bash
./watcher --include-dirs /tmp echo "Detected:"
```

gnore hidden files and only trigger for `.txt` files:

```bash
./watcher --ignore-hidden --include '*.txt' /tmp echo "Detected text file:"
```

Include `.jpg` files but exclude anything in `cache/`:

```bash
./watcher --ignore-hidden \
  --include '*.jpg' \
  --exclude 'cache/*' \
  /tmp echo "New image:"
```

## 🎯 Pattern Matching

* Uses Go's filepath.Match.
* Patterns are checked relative to the root directory:
  * For `/tmp/images/photo.jpg` with root `/tmp`, the relative path is `images/photo.jpg`.
* Examples:
  * `*.txt` — matches `file.txt` in the root.
  * `images/*.jpg` — matches `images/pic.jpg` but not `nested/images/pic.jpg`.
  * `**/*.log` — Note: Go's filepath.Match does not support `**` by default, so use `*/*.log` or more specific patterns.

## 📜 License

MIT

## 👤 Author

Eduardo Gonzalez Solares

📧 Email: [your-email@example.com]

🐙 GitLab: [your-gitlab-profile-url]

🌐 Website: [optional-website]