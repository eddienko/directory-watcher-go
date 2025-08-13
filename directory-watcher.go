package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func main() {
	includeDirs := flag.Bool("include-dirs", false, "Trigger command for directories too")
	ignoreHidden := flag.Bool("ignore-hidden", false, "Ignore hidden files and directories")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [--include-dirs] [--ignore-hidden] <directory_to_watch> <command> [args...]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	// Need at least dir + command
	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	rootDir := flag.Arg(0)
	cmd := flag.Arg(1)
	cmdArgs := flag.Args()[2:]

	// Validate rootDir
	info, err := os.Stat(rootDir)
	if err != nil {
		log.Fatalf("Error accessing directory: %v", err)
	}
	if !info.IsDir() {
		log.Fatalf("%s is not a directory", rootDir)
	}

	// Create watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %v", err)
	}
	defer watcher.Close()

	// Check if a path is hidden
	isHidden := func(path string) bool {
		base := filepath.Base(path)
		return strings.HasPrefix(base, ".")
	}

	// Add directories recursively
	addDirRecursive := func(dir string) error {
		return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if *ignoreHidden && isHidden(path) {
				if info.IsDir() {
					return filepath.SkipDir // Skip whole hidden directory
				}
				return nil
			}
			if info.IsDir() {
				err = watcher.Add(path)
				if err != nil {
					return fmt.Errorf("failed to watch %s: %w", path, err)
				}
				log.Printf("Watching: %s", path)
			}
			return nil
		})
	}

	if err := addDirRecursive(rootDir); err != nil {
		log.Fatalf("Error adding directories: %v", err)
	}

	// Event loop
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					if *ignoreHidden && isHidden(event.Name) {
						log.Printf("Ignoring hidden: %s", event.Name)
						continue
					}

					info, err := os.Stat(event.Name)
					if err == nil && info.IsDir() {
						log.Printf("New directory detected: %s", event.Name)
						if err := addDirRecursive(event.Name); err != nil {
							log.Printf("Error watching new directory: %v", err)
						}
						if *includeDirs {
							runCommand(cmd, cmdArgs, event.Name)
						}
					} else {
						log.Printf("New file detected: %s", event.Name)
						runCommand(cmd, cmdArgs, event.Name)
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	log.Printf("Recursive watching started at: %s (include dirs: %v, ignore hidden: %v)", rootDir, *includeDirs, *ignoreHidden)
	<-make(chan struct{}) // Block forever
}

func runCommand(cmd string, args []string, newPath string) {
	execCmd := exec.Command(cmd, append(args, newPath)...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	if err := execCmd.Run(); err != nil {
		log.Printf("Error running command: %v", err)
	}
}
