# logrus-hooks

* Rotate Hook, rotate the log file by daily(default) or hour
* Source Hook, append the caller to log's message 

## How to use

```go
package main

import (
    "github.com/Sirupsen/logrus"
    hooks "github.com/git-hulk/logrus-hooks"
)

func main() {
    logger := logrus.New()
    // Create rotate hook
    rotateHook, err := hooks.NewRotateHook(logger, "/www/mydir", "test")
    if err != nil {
        // do something and exit
    }
    logger.Hooks.Add(rotateHook)
    // Create source hook
    logger.Hooks.Add(hooks.NewSourceHook(logrus.InfoLevel))
    logger.Info("foo")
}
```
