package eep

// Eval eval expression
func Eval(expression string) (interface{}, error) {
	return eval(expression, make(environ))
}

// EvalWithEnv eval expression with env
func EvalWithEnv(expression string, env map[string]interface{}) (interface{}, error) {
	return eval(expression, env)
}

func eval(expression string, env map[string]interface{}) (interface{}, error) {
	tokens, err := newScanner(expression).scanTokens()
	if err != nil {
		return nil, err
	}

	exp, parseErr := newParser(tokens).parse()
	if parseErr != nil {
		return nil, parseErr
	}

	result, evalErr := newEvaluator(env).Eval(exp)
	if evalErr != nil {
		return nil, evalErr
	}

	return result, nil
}
