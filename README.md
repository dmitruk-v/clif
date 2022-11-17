# Clif
Clif is a small and simple framework to build CLI-based applications

## How to use
Basic example:

```go
import (
  "github.com/dmitruk-v/clif/v0"
)

func main() {
  cfg := clif.AppConfig{
    Commands: clif.Commands{
      clif.NewCommand(`command:\+ amount:\d+ currency:\w{3}`, &depositController{}, clif.CommandHelp{
        Info:  "Deposit some amount of currency to balance:",
        Usage: []string{
          "+ AMOUNT CURRENCY",
        },
        Examples: []string{
          "+ 100 usd   Adds 100 USD to balance",
        },
      }),
    },
  }
  app := clif.NewApp(cfg)
  if err := app.Run(); err != nil {
    log.Fatal(err)
  }
}

type depositController struct{}

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
