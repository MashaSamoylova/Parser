package main
import (
    "fmt"

    "github.com/golang-collections/collections/stack"
)

type Rule struct {
    Term string
    Alternatives []string
}


func Parse(expression string) {
    
}

func main() {
    grammar := []Rule{
        {
            Term: "B",
            Alternatives: []string{
                "T + B",
                "T",
            },
        },
        {
            Term: "T",
            Alternatives: []string{
                "M",
                "M * T",
            },
        },
        {
            Term: "M",
            Alternatives: []string{
              "a",
              "b",
            },
        },

    }
    for _, g := range grammar {
        fmt.Println(g)
    }
    var inputExpression string
    fmt.Println("Input your expression:")
    fmt.Scanf("%s", inputExpression)
}
