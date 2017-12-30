package eep

import (
	"fmt"
	"log"
)

type environ map[string]interface{}

type evaluator struct {
	env environ
}

// NewEvaluator new evaluator
func newEvaluator(env environ) *evaluator {
	return &evaluator{env: env}
}

func (elt *evaluator) Eval(exp Expr) (val interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("parser error: %v\n", e)
			switch e.(type) {
			case string:
				err = fmt.Errorf(e.(string))
			case error:
				err = e.(error)
			case parserErr:
				err = fmt.Errorf(e.(parserErr).msg)
			default:
				err = fmt.Errorf(fmt.Sprintf("%v", e))
			}
		}
	}()

	val = elt.evaluate(exp)
	return val, err
}

func (elt *evaluator) evaluate(exp Expr) interface{} {
	return exp.Accept(elt)
}

func (elt *evaluator) VisitBinaryExpr(exp *Binary) interface{} {
	left, right := elt.evaluate(exp.Left), elt.evaluate(exp.Right)

	switch exp.Operator.TokenType {
	case Minus:
		v1, v2 := elt.checkNumberOperands(*exp.Operator, left, right)
		return v1 - v2
	case Slash:
		v1, v2 := elt.checkNumberOperands(*exp.Operator, left, right)
		return v1 / v2
	case Star:
		v1, v2 := elt.checkNumberOperands(*exp.Operator, left, right)
		return v1 * v2
	case Plus:
		v, ok := left.(float64)
		v1, ok1 := right.(float64)
		if ok && ok1 {
			return v + v1
		}
		v2, ok2 := left.(string)
		v3, ok3 := right.(string)
		if ok2 && ok3 {
			return v2 + v3
		}
		panic(fmt.Sprintf("%s Operands must be two numbers or two strings!", exp.Operator.Lexeme))
	case Greater:
		v1, v2 := elt.checkNumberOperands(*exp.Operator, left, right)
		return v1 > v2
	case GreaterEqual:
		v1, v2 := elt.checkNumberOperands(*exp.Operator, left, right)
		return v1 >= v2
	case Less:
		v1, v2 := elt.checkNumberOperands(*exp.Operator, left, right)
		return v1 < v2
	case LessEqual:
		v1, v2 := elt.checkNumberOperands(*exp.Operator, left, right)
		return v1 <= v2
	case EqualEqual:
		return elt.isEqual(left, right)
	case BangEqual:
		return !elt.isEqual(left, right)
	default:
		return nil
	}
}

func (elt *evaluator) VisitCallExpr(expr *Call) interface{} {
	callee := elt.evaluate(expr.Callee)

	var arguments []interface{}
	for _, arg := range expr.Arguments {
		arguments = append(arguments, elt.evaluate(arg))
	}

	function, ok := callee.(func(...interface{}) interface{})
	if !ok {
		panic(fmt.Sprintf("callee: %v is illegal function type", callee))
	}

	return function(arguments...)
}

func (elt *evaluator) VisitGroupingExpr(exp *Grouping) interface{} {
	return elt.evaluate(exp.Expression)
}

func (elt *evaluator) VisitLiteralExpr(exp *Literal) interface{} {
	return exp.Value
}

func (elt *evaluator) VisitLogicalExpr(expr *Logical) interface{} {
	left := elt.evaluate(expr.Left)

	if expr.Operator.TokenType == OR {
		if elt.isTruthy(left) {
			return left
		}
	} else {
		if !elt.isTruthy(left) {
			return left
		}
	}

	return elt.evaluate(expr.Right)
}

func (elt *evaluator) VisitUnaryExpr(exp *Unary) interface{} {
	right := elt.evaluate(exp.Right)

	switch exp.Operator.TokenType {
	case Bang:
		return !elt.isTruthy(right)
	case Minus:
		v := elt.checkNumberOperand(*exp.Operator, right)
		return 0 - v
	}
	return nil
}

func (elt *evaluator) VisitVariableExpr(exp *Variable) interface{} {
	return elt.env[exp.Name.Lexeme]
}

func (elt *evaluator) checkNumberOperand(operator Token, obj interface{}) float64 {
	v, ok := obj.(float64)
	if !ok {
		panic(fmt.Sprintf("%s Operand must be a number.", operator.Lexeme))
	}
	return v
}

func (elt *evaluator) checkNumberOperands(operator Token, left, right interface{}) (float64, float64) {
	v1, ok := left.(float64)
	v2, ok1 := right.(float64)
	if !ok || !ok1 {
		panic(fmt.Sprintf("%s Operands must be numbers!", operator.Lexeme))
	}
	return v1, v2
}

func (elt *evaluator) isEqual(left, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}

	if left == nil {
		return false
	}

	return left == right
}

func (elt *evaluator) isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}

	v, ok := obj.(bool)
	if ok {
		return v
	}
	return true
}
