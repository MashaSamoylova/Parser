package parser

import (
    "fmt"
    "strings"

    "github.com/golang-collections/collections/stack"
)


type ParseSystem struct {
    Expression string
    stack1         *stack.Stack
    stack2         *stack.Stack
    Grammar    []Rule
}

type Rule struct {
    Term         string
    Alternatives []string
}

type step struct {
    Term  string
    Index int
}

func (p *ParseSystem) PrintResult() {
    for p.stack1.Len() > 0 {
        fmt.Println(p.stack1.Pop())
    }
}
func (p *ParseSystem) IsTerm(a string) ([]string, bool) {
    for _, rule := range p.Grammar {
        if rule.Term == a {
            return rule.Alternatives, true
        }
    }
    return nil, false
}

func (p *ParseSystem) Back() {
    fmt.Println("###############\nBack")
    previousStep, ok := p.stack1.Pop().(step)
    if !ok {
        panic("Грамматика не применима к данному выражению!")
    }
    previousExprStack2 := p.stack2.Pop().(string)
    fmt.Printf("previousStep: %v\n", previousStep)
    fmt.Printf("prev exprStack2: %v\n", previousExprStack2)
    alternatives, ok := p.IsTerm(previousStep.Term)
    if ok {
        fmt.Println("Это терм, значит ищем альтернативы или продолжаем возвращаться")
        if len(alternatives) <= previousStep.Index+1 {
            fmt.Printf("Закончились альтернативы у %v, продолжаем возвращение\n", previousStep.Term)
            p.Back()
            return
        }
        fmt.Println("Альтернативы не закончились")
        symbolsExprStack2 := strings.Split(previousExprStack2, " ")
        oldAltLen := len(strings.Split(alternatives[previousStep.Index], " "))
        newExprStack2 := strings.Join(append([]string{alternatives[previousStep.Index+1]}, symbolsExprStack2[oldAltLen:]...), " ")
        fmt.Printf("Заменили на другую альтернативу, сформировали новое выражение в Stack2: %v\n", newExprStack2)
        p.stack2.Push(newExprStack2)
        newStep := step{
            Term:  previousStep.Term,
            Index: previousStep.Index + 1,
        }
        fmt.Printf("Пушим шаг: %v\n", newStep)
        p.stack1.Push(newStep)
        fmt.Println("Вышел из рекурсии")
        return
    } else {
        fmt.Println("Это не терм, продолжаем возвращаться")
        p.Expression = previousStep.Term + " " + p.Expression
        fmt.Printf("Вернули терм в выражение: %v\n", p.Expression)
        p.Back()
    }

}

func (p *ParseSystem) Parse() {
    p.stack1 = stack.New()
    p.stack2 = stack.New()
    p.stack1.Push("")
    p.stack2.Push(p.Grammar[0].Term)
    p.parse()
}

func (p *ParseSystem) parse() {
    for len(p.Expression) > 0 && p.stack2.Peek() != "" {
        fmt.Println("**************************")
        fmt.Printf("Expression: %v\n", p.Expression)
        exprStack2 := p.stack2.Peek().(string)
        symbolsStack2 := strings.Split(exprStack2, " ")
        inputSymbol := symbolsStack2[0]
        tail := symbolsStack2[1:]
        alternatives, ok := p.IsTerm(inputSymbol)
        if ok {
            fmt.Printf("Input symbol is a term: %v\n", inputSymbol)
            newExprStack2 := strings.Join(append([]string{alternatives[0]}, tail...), " ")
            fmt.Printf("Create new expression Stack2: %v\n", newExprStack2)
            p.stack2.Push(newExprStack2)
            p.stack1.Push(step{
                Term:  inputSymbol,
                Index: 0,
            })
            continue
        } else {
            fmt.Printf("Input symbol is not a term: %v\n", inputSymbol)
            symbolsExpr := strings.Split(p.Expression, " ")
            var i, j int
            for i = 0; i < len(symbolsStack2) && i < len(symbolsExpr); i++ {
                if symbolsStack2[i] != symbolsExpr[i] {
                    break
                }
            }
            if i == 0 {
                fmt.Println("Выражения не совпали, откатываемся")
                p.Back()
                continue
            }
            for j = 0; j < i; j++ {
                fmt.Printf("Symbols are same: %v and %v\n", symbolsStack2[j], symbolsExpr[j])
                p.stack1.Push(step{
                    Term:  symbolsExpr[j],
                    Index: -1,
                })
                cutExprStack2 := strings.Join(symbolsStack2[j+1:], " ")
                cutExpr := strings.Join(symbolsExpr[j+1:], " ")
                fmt.Printf("Cut Stack2 expr: %v\n", cutExprStack2)
                fmt.Printf("Cut expr: %v\n", cutExpr)
                p.Expression = cutExpr
                p.stack2.Push(cutExprStack2)
            }
            if len(p.Expression) == 0 && p.stack2.Peek().(string) != "" {
                fmt.Println("Откат когда stack2 не закончился, а выражение пусто")
                p.Back()
                continue
            }
        }
    }
    fmt.Println("Вышли с успехом")
    p.PrintResult()
}

