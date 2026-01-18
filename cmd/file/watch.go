package file

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

// watchCmd represents the watch command
var watchCmd = &cobra.Command{
	Use:   "watch [path]",
	Short: "Watch files for changes",
	Long: `Watch files and directories for changes and execute commands.

Examples:
  devkit file watch ./src
  devkit file watch ./src --on-change "go build"
  devkit file watch . --pattern "*.go" --on-change "go test ./..."`,
	Args: cobra.MaximumNArgs(1),
	RunE: runWatch,
}

func init() {
	fileCmd.AddCommand(watchCmd)

	watchCmd.Flags().StringP("pattern", "p", "*", "File pattern to watch")
	watchCmd.Flags().String("on-change", "", "Command to execute on file change")
	watchCmd.Flags().BoolP("recursive", "r", true, "Watch recursively")
}

func runWatch(cmd *cobra.Command, args []string) error {
	watchPath := "."
	if len(args) > 0 {
		watchPath = args[0]
	}

	pattern, _ := cmd.Flags().GetString("pattern")
	onChange, _ := cmd.Flags().GetString("on-change")
	recursive, _ := cmd.Flags().GetBool("recursive")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	// Add path to watcher
	if recursive {
		err = filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return watcher.Add(path)
			}
			return nil
		})
	} else {
		err = watcher.Add(watchPath)
	}

	if err != nil {
		return fmt.Errorf("failed to add path to watcher: %w", err)
	}

	fmt.Printf("Watching: %s (pattern: %s)\n", watchPath, pattern)
	if onChange != "" {
		fmt.Printf("On change: %s\n", onChange)
	}
	fmt.Println("Press Ctrl+C to stop...\n")

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				matched, _ := filepath.Match(pattern, filepath.Base(event.Name))
				if !matched {
					continue
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Printf("[%s] Modified: %s\n", time.Now().Format("15:04:05"), event.Name)

					if onChange != "" {
						cmd := exec.Command("sh", "-c", onChange)
						cmd.Stdout = os.Stdout
						cmd.Stderr = os.Stderr
						if err := cmd.Run(); err != nil {
							fmt.Printf("Error executing command: %v\n", err)
						}
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("Watcher error: %v\n", err)
			}
		}
	}()

	<-done
	return nil
}
