# Clim
Clim is a small and simple framework to build CLI-based applications

## How to use
Basic example:

```go
import (
  "github.com/dmitruk-v/clim/v0"
)

func main() {
  cfg := clim.AppConfig{
    Commands: clim.Commands{
      clim.NewCommand(`command:\+ amount:\d+ currency:\w{3}`, NewDepositController()),
      clim.NewQuitCommand(`command:quit|exit`),
    },
  }
  app := clim.NewApp(cfg)
  if err := app.Run(); err != nil {
    log.Fatal(err)
  }
}

type depositController struct{}

func NewDepositController() *depositController {
  return &depositController{}
}

func (ctrl *depositController) Handle(req map[string]string) error {
  fmt.Println("deposit controller got request:", req)
  return nil
}
```

For example, input command ```+ 100 usd``` will output:
```code
deposit controller got request: map[amount:100 command:+ currency:usd]
```
## License

MIT
