# go-logger
simple logger for my golang based applications

## import and use:

```
import (
	"github.com/gtsteffaniak/go-logger/logger"
)

func main() {
	err := logger.SetupLogger("STDOUT", "INFO,DEBUG", "INFO,ERROR", false)
	if err != nil {
		// Handle error
	}

	logger.Info("This is an info message from the logger.")
	logger.Debug("This is a debug message with details: %s", "some detail")
	logger.Api("API call successful", 200)
}
```
## linting

```
go mod tidy
go tool golangci-lint run
```