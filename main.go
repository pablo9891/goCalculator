package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
)

const (
	IntegerNumberTk int = iota
	DoubleNumberTk
	LeftParenthesisTk
	RightParenthesisTk
	AddTk
	SubstractionTk
	UnarySubstractionTk
	UnaryAddTk
	DivisionTk
	MultiplicationTk
	ExponentiationTk
	RootTk
	NotValidTk
)

const (
	UnarySubstractionPrec = 11
	UnaryAddPrec          = 11
	ExponentiationPrec    = 10
	RootPrec              = 10
	MultiplicationPrec    = 9
	DivisionPrec          = 9
	AddPrec               = 7
	SubstractionPrec      = 7
	UndefinedPrec         = -1
	ParenthesisPrec       = -5
)

const (
	Operand = iota
	Operator
	Parenthesis
	Undefined
)

type Token struct {
	token         string
	tokenType     int
	position      int
	precedence    int
	exprTokenType int
}

type Element struct {
	val *Token
}

type CompareFunc func(a *Element, b *Element) bool

type PriorityStack struct {
	size  int
	idx   int
	stack []*Element
	cmp   CompareFunc
}

func Create(params ...CompareFunc) *PriorityStack {
	if len(params) == 1 {
		return &PriorityStack{0, -1, make([]*Element, 0), params[0]}
	}
	return &PriorityStack{0, -1, make([]*Element, 0), nil}
}

func removeIndex(s []*Element, index int) []*Element { return append(s[:index], s[index+1:]...) }

func Push(val *Token, stack *PriorityStack) {
	stack.idx++
	stack.size++
	stack.stack = append(stack.stack, &Element{val})
	if stack.cmp != nil {
		i := stack.idx
		for i > 0 && stack.cmp(stack.stack[i], stack.stack[i-1]) {
			tmp := stack.stack[i-1]
			stack.stack[i-1] = stack.stack[i]
			stack.stack[i] = tmp
			i--
		}
	}
}

func comparePrec(a *Element, b *Element) bool { return a.val.precedence < b.val.precedence }

func Pop(stack *PriorityStack) (*Token, bool) {
	if stack.idx == -1 {
		return nil, false
	}
	e := stack.stack[stack.idx]
	stack.stack = removeIndex(stack.stack, stack.idx)
	stack.size--
	stack.idx--
	return e.val, true
}

func Peek(stack *PriorityStack) (*Token, bool) {
	if stack.idx == -1 {
		return nil, false
	}
	return stack.stack[stack.idx].val, true
}

func IsEmpty(stack *PriorityStack) bool { return stack.idx == -1 }

func isEmptySpaces(c string) bool { return (c == " " || c == "\n" || c == "\t") }

func skipEmptySpaces(index *int, expr string) {
	for *index < len(expr) && isEmptySpaces(string(expr[*index])) {
		*index++
	}
}

func isParenthesis(c string) bool { return c == "(" || c == ")" }

func isOperator(c string) bool { return c == "^" || c == "+" || c == "-" || c == "*" || c == "/" }

func isNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 10, 8)
	return err == nil
}

func createToken(position int, c string, tokenType int, precedence int, exprTokenType int) *Token {
	return &Token{c, tokenType, position, precedence, exprTokenType}
}

func processParenthesis(index *int, c string) *Token {
	tokType := NotValidTk

	if c == "(" {
		tokType = LeftParenthesisTk
	} else {
		tokType = RightParenthesisTk
	}
	return createToken(*index, c, tokType, ParenthesisPrec, Parenthesis)
}

func isBinary(token *Token) bool {
	if token == nil {
		return false
	}
	tkType := token.tokenType
	return tkType == RightParenthesisTk || tkType == IntegerNumberTk || tkType == DoubleNumberTk
}

func processNumber(index *int, expr string) (*Token, error) {
	number := ""
	t := createToken(*index, "", NotValidTk, UndefinedPrec, Operand)

	t.position = *index
	for *index < len(expr) && isNumeric(string(expr[*index])) {
		number = number + string(expr[*index])
		*index++
	}
	if *index < len(expr) && string(expr[*index]) == "." {
		number = number + string(expr[*index])
		*index++
		for *index < len(expr) && isNumeric(string(expr[*index])) {
			number = number + string(expr[*index])
			*index++
		}
		t.tokenType = DoubleNumberTk
	} else {
		t.tokenType = IntegerNumberTk
	}

	if *index < len(expr) && !isParenthesis(string(expr[*index])) && !isOperator(string(expr[*index])) && !isEmptySpaces(string(expr[*index])) {
		return nil, errors.New(fmt.Sprint("undefined token at position ", 1+(*index-(*index-t.position))))
	}
	t.token = number
	*index--
	return t, nil
}

func getNextToken(index *int, expr string, lastToken *Token) (*Token, error) {
	var t *Token
	var err error

	t = nil
	err = nil
	skipEmptySpaces(index, expr)

	if isParenthesis(string(expr[*index])) {
		t = processParenthesis(index, string(expr[*index]))
	} else if string(expr[*index]) == "+" {
		t = createToken(*index, string(expr[*index]), UnaryAddTk, AddPrec, Operator)
		if isBinary(lastToken) {
			t.tokenType = AddTk
		}
	} else if string(expr[*index]) == "-" {
		t = createToken(*index, string(expr[*index]), UnarySubstractionTk, SubstractionPrec, Operator)
		if isBinary(lastToken) {
			t.tokenType = SubstractionTk
		}
	} else if string(expr[*index]) == "*" {
		t = createToken(*index, string(expr[*index]), MultiplicationTk, MultiplicationPrec, Operator)
	} else if string(expr[*index]) == "/" {
		t = createToken(*index, string(expr[*index]), DivisionTk, DivisionPrec, Operator)
	} else if isNumeric(string(expr[*index])) {
		t, err = processNumber(index, expr)
	} else if string(expr[*index]) == "^" {
		t = createToken(*index, string(expr[*index]), ExponentiationTk, ExponentiationPrec, Operator)
	} else if string(expr[*index]) == "$" {
		t = createToken(*index, string(expr[*index]), RootTk, RootPrec, Operator)
	} else {
		t = createToken(*index, string(expr[*index]), NotValidTk, UndefinedPrec, Undefined)
		err = errors.New(fmt.Sprint("undefined token at position ", 1+(*index-(*index-t.position))))
	}

	*index++
	return t, err
}

func applyBinary(result []float64, idxResults *int, tk *Token) ([]float64, error) {
	var res float64
	if len(result) < 2 {
		return nil, errors.New("not valid operation")
	}
	n1 := result[*idxResults]
	result = append(result[:*idxResults], result[*idxResults+1:]...)
	*idxResults--
	n2 := result[*idxResults]
	result = append(result[:*idxResults], result[*idxResults+1:]...)
	*idxResults--
	switch tk.tokenType {
	case AddTk:
		res = n1 + n2
	case SubstractionTk:
		res = n2 - n1
	case MultiplicationTk:
		res = n1 * n2
	case RootTk:
		if n2 < 0 {
			return nil, errors.New("negative root base undefined")
		}
		if n1 == 0 {
			return nil, errors.New("exponent zero root undefined")
		}
		res = math.Pow(n2, 1/n1)
	case DivisionTk:
		if n1 == 0 {
			return nil, errors.New("division by zero undefined")
		}
		res = n2 / n1
	case ExponentiationTk:
		res = math.Pow(n2, n1)
	}
	result = append(result, res)
	*idxResults++
	return result, nil
}

func evaluate(postfixExpr []*Token) (float64, error) {
	result := make([]float64, 0)
	idx := 0
	idxResults := -1
	var err error

	for idx < len(postfixExpr) {
		tk := postfixExpr[idx]
		if tk.exprTokenType == Operator {
			switch tk.tokenType {
			case UnarySubstractionTk:
				result[idxResults] = result[idxResults] * (-1)
			case UnaryAddTk:
				result[idxResults] = result[idxResults] * 1
			case AddTk:
				result, err = applyBinary(result, &idxResults, tk)
			case SubstractionTk:
				result, err = applyBinary(result, &idxResults, tk)
			case MultiplicationTk:
				result, err = applyBinary(result, &idxResults, tk)
			case DivisionTk:
				result, err = applyBinary(result, &idxResults, tk)
			case ExponentiationTk:
				result, err = applyBinary(result, &idxResults, tk)
			case RootTk:
				result, err = applyBinary(result, &idxResults, tk)
			default:
				break
			}
			if err != nil {
				return 0.0, err
			}
		} else {
			n, _ := strconv.ParseFloat(tk.token, 64)
			result = append(result, n)
			idxResults++
		}
		idx++
	}

	return result[0], nil
}

func processExpression(expr string) (float64, error) {
	index := 0
	postfixIdx := -1
	postfixList := make([]*Token, 0)
	operatorStack := Create(comparePrec)
	var lastToken *Token
	lastToken = nil

	for index < len(expr) {
		t, errTk := getNextToken(&index, expr, lastToken)
		lastToken = t
		if errTk != nil {
			return 0.0, errTk
		}

		if t.exprTokenType == Operand {
			postfixIdx++
			postfixList = append(postfixList, t)
		}
		if t.exprTokenType == Operator || t.exprTokenType == Parenthesis {
			if t.tokenType == RightParenthesisTk {
				var err bool
				var tk *Token
				for {
					tk, err = Peek(operatorStack)
					if !err || tk.tokenType == LeftParenthesisTk {
						break
					}
					postfixList = append(postfixList, tk)
					Pop(operatorStack)
				}
				if !err {
					return 0.0, errors.New(fmt.Sprint("not balanced right parenthesis at position ", t.position))
				}
				Pop(operatorStack) // Removes left parenthesis from stack
			} else if t.tokenType == LeftParenthesisTk {
				operatorStack.idx++
				operatorStack.size++
				operatorStack.stack = append(operatorStack.stack, &Element{t})
			} else {
				var err bool
				var tk *Token
				for {
					tk, err = Peek(operatorStack)
					if !err || tk.precedence < t.precedence || (tk.precedence == t.precedence && tk.tokenType == ExponentiationTk) {
						break
					}
					postfixList = append(postfixList, tk)
					Pop(operatorStack)
				}
				Push(t, operatorStack)
			}
		}
	}
	for !IsEmpty(operatorStack) {
		tk, _ := Pop(operatorStack)
		if tk.tokenType == LeftParenthesisTk {
			return 0.0, errors.New(fmt.Sprint("not balanced left parenthesis at position ", tk.position))
		}
		postfixList = append(postfixList, tk)
	}
	return evaluate(postfixList)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("error in parameters")
		return
	}
	res, err := processExpression(os.Args[1])
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%v\n", res)
	}
}
