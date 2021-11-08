package main

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/hugolgst/rich-go/client"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getIcon(command string) (string, string) {
	switch command {
	case "git":
		return "git", "git source control"
	case "curl":
		return "curl", "the downloady thing"
		//case "go":
		//	return "golang", "go lang"
		//case "rustup":
		//	fallthrough
		//case "cargo":
		//	return "rust", "crab crab crab"
	}
	return "fish", "the friendly, interactive shell"
}

var startTime = time.Now()

func updateStatus(command string) {
	parts := strings.Split(command, " ")
	icon, iconText := getIcon(parts[0])
	err := client.SetActivity(client.Activity{
		State:   fmt.Sprintf("Running `%s`", parts[0]),
		Details: "",
		Timestamps: &client.Timestamps{
			Start: &startTime,
		},
		LargeImage: icon,
		LargeText:  iconText,
		SmallImage: "bash",
		SmallText:  "The friendly interactive shell",
	})
	if err != nil {
		fmt.Println("couldn't publish discord rich status", err)
	}
}

func trimLeftChars(s string, n int) string {
	m := 0
	for i := range s {
		if m >= n {
			return s[i:]
		}
		m++
	}
	return s[:0]
}

func main() {
	err := client.Login("907099827547033610")
	if err != nil {
		panic(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {

		}
	}(watcher)

	done := make(chan bool)
	go func() {
		var line string
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					content, err := ioutil.ReadFile(event.Name)
					if err != nil {
						fmt.Println("error:", err)
					}
					hist := bytes.Split(content, []byte{'\n'})
					last := string(hist[len(hist)-3])
					last = trimLeftChars(last, 7)
					if last == line {
						continue
					}
					line = last
					updateStatus(line)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					continue
				}
				fmt.Println("error:", err)
			}
		}
	}()

	homedir, _ := os.UserHomeDir()
	err = watcher.Add(filepath.Join(homedir, ".local/share/fish/fish_history"))
	if err != nil {
		panic(err)
	}
	<-done
}
