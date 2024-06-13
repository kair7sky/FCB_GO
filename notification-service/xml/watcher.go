package xml

import (
    "io/ioutil"
    "log"

    "github.com/fsnotify/fsnotify"
)

func WatchXMLFile(filePath string, onChange func(changes string)) error {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return err
    }
    defer watcher.Close()

    done := make(chan bool)

    go func() {
        var lastContent []byte
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                if event.Op&fsnotify.Write == fsnotify.Write {
                    content, err := ioutil.ReadFile(filePath)
                    if err != nil {
                        log.Printf("Error reading file: %v", err)
                        continue
                    }
                    if string(content) != string(lastContent) {
                        lastContent = content
                        onChange(string(content))
                    }
                }
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Printf("Error: %v", err)
            }
        }
    }()

    err = watcher.Add(filePath)
    if err != nil {
        return err
    }

    <-done
    return nil
}
