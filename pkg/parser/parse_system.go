package parser

import (
	"fmt"
	"strings"

	"github.com/golang-collections/collections/stack"
)

// ParseSystem is a structure for parsing, which consists `stack1`,
// `stack2` and `expression` components as like in top down parsing algorithm.
type ParseSystem struct {
	Expression string
	stack1     *stack.Stack
	stack2     *stack.Stack
	Grammar    []Rule
}

// Rule describes a method of grammar definition in my realization.
type Rule struct {
	Term         string
	Alternatives []string
}

type step struct {
	Term  string
	Index int
}

func (p *ParseSystem) printResult() {
	stepsLen := p.stack1.Len() - 1
	resultSteps := make([]step, stepsLen)
	for i := stepsLen - 1; i >= 0; i-- {
		resultSteps[i] = p.stack1.Pop().(step)
	}
	for _, step := range resultSteps {
		if step.Index != -1 {
			alternatives, _ := p.isTerm(step.Term)
			fmt.Printf("%s_%d) %s -> %s\n", step.Term, step.Index, step.Term, alternatives[step.Index])
		}
	}
}

func (p *ParseSystem) isTerm(a string) ([]string, bool) {
	for _, rule := range p.Grammar {
		if rule.Term == a {
			return rule.Alternatives, true
		}
	}
	return nil, false
}

func (p *ParseSystem) back() {
	previousStep, ok := p.stack1.Pop().(step)
	// we catch the empty word. and we cant to continue return.
	if !ok {
		panic("The grammar is not applied to this expression!")
	}
	previousExprStack2 := p.stack2.Pop().(string)
	alternatives, ok := p.isTerm(previousStep.Term)
	if ok {
		if len(alternatives) <= previousStep.Index+1 {
			p.back()
			return
		}
		symbolsExprStack2 := strings.Split(previousExprStack2, " ")
		oldAltLen := len(strings.Split(alternatives[previousStep.Index], " "))
		newExprStack2 := strings.Join(append([]string{alternatives[previousStep.Index+1]}, symbolsExprStack2[oldAltLen:]...), " ")
		p.stack2.Push(newExprStack2)
		newStep := step{
			Term:  previousStep.Term,
			Index: previousStep.Index + 1,
		}
		p.stack1.Push(newStep)
		return
	} else {
		p.Expression = previousStep.Term + " " + p.Expression
		p.back()
	}

}

// Parse parses expression
func (p *ParseSystem) Parse() {
	p.stack1 = stack.New()
	p.stack2 = stack.New()
	p.stack1.Push("")
	p.stack2.Push(p.Grammar[0].Term)
	p.parse()
}

func (p *ParseSystem) parse() {
	for len(p.Expression) > 0 && p.stack2.Peek() != "" {
		exprStack2 := p.stack2.Peek().(string)
		symbolsStack2 := strings.Split(exprStack2, " ")
		inputSymbol := symbolsStack2[0]
		tail := symbolsStack2[1:]
		alternatives, ok := p.isTerm(inputSymbol)
		if ok {
			//if input symbol is term: apply rule
			newExprStack2 := strings.Join(append([]string{alternatives[0]}, tail...), " ")
			p.stack2.Push(newExprStack2)
			p.stack1.Push(step{
				Term:  inputSymbol,
				Index: 0,
			})
			continue
		} else {
			symbolsExpr := strings.Split(p.Expression, " ")
			var i, j int
			// compare stack2 and expression
			for i = 0; i < len(symbolsStack2) && i < len(symbolsExpr); i++ {
				if symbolsStack2[i] != symbolsExpr[i] {
					break
				}
			}
			// push on the stack1 symbol, which is not a term, but equals with expression symbol
			for j = 0; j < i; j++ {
				p.stack1.Push(step{
					Term:  symbolsExpr[j],
					Index: -1,
				})
				cutExprStack2 := strings.Join(symbolsStack2[j+1:], " ")
				cutExpr := strings.Join(symbolsExpr[j+1:], " ")
				p.Expression = cutExpr
				p.stack2.Push(cutExprStack2)
			}
			// if (stack2 is not empty) or (stack2 and expression are not equal)
			if (len(p.Expression) == 0 && p.stack2.Peek().(string) != "") || (i == 0) {
				p.back()
				continue
			}
		}
	}
	p.printResult()
}
