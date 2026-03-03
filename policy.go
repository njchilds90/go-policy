package policy

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type Operator string

const (
	OpEqual              Operator = "equal"
	OpNotEqual           Operator = "not_equal"
	OpGreaterThan        Operator = "greater_than"
	OpLessThan           Operator = "less_than"
	OpGreaterOrEqual     Operator = "greater_or_equal"
	OpLessOrEqual        Operator = "less_or_equal"
	OpIn                 Operator = "in"
	OpAnd                Operator = "and"
	OpOr                 Operator = "or"
)

type Rule struct {
	Operator Operator `json:"operator"`
	Field    string   `json:"field,omitempty"`
	Value    any      `json:"value,omitempty"`
	Rules    []Rule   `json:"rules,omitempty"`
}

type Result struct {
	Allowed bool   `json:"allowed"`
	Error   string `json:"error,omitempty"`
}

type EvaluationError struct {
	Reason string
}

func (e EvaluationError) Error() string {
	return e.Reason
}

func Evaluate(ctx context.Context, rule Rule, input map[string]any) (Result, error) {
	select {
	case <-ctx.Done():
		return Result{}, ctx.Err()
	default:
	}

	ok, err := evaluateRule(rule, input)
	if err != nil {
		return Result{Allowed: false, Error: err.Error()}, err
	}

	return Result{Allowed: ok}, nil
}

func evaluateRule(rule Rule, input map[string]any) (bool, error) {
	switch rule.Operator {

	case OpAnd:
		for _, r := range rule.Rules {
			ok, err := evaluateRule(r, input)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
		}
		return true, nil

	case OpOr:
		for _, r := range rule.Rules {
			ok, err := evaluateRule(r, input)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
		}
		return false, nil

	default:
		val, exists := input[rule.Field]
		if !exists {
			return false, EvaluationError{
				Reason: fmt.Sprintf("field not found: %s", rule.Field),
			}
		}
		return compare(rule.Operator, val, rule.Value)
	}
}

func compare(op Operator, a, b any) (bool, error) {
	switch op {

	case OpEqual:
		return reflect.DeepEqual(a, b), nil

	case OpNotEqual:
		return !reflect.DeepEqual(a, b), nil

	case OpGreaterThan, OpLessThan, OpGreaterOrEqual, OpLessOrEqual:
		af, aok := toFloat(a)
		bf, bok := toFloat(b)
		if !aok || !bok {
			return false, EvaluationError{
				Reason: "numeric comparison requires numbers",
			}
		}

		switch op {
		case OpGreaterThan:
			return af > bf, nil
		case OpLessThan:
			return af < bf, nil
		case OpGreaterOrEqual:
			return af >= bf, nil
		case OpLessOrEqual:
			return af <= bf, nil
		}

	case OpIn:
		rv := reflect.ValueOf(b)
		if rv.Kind() != reflect.Slice {
			return false, EvaluationError{
				Reason: "in operator requires slice value",
			}
		}
		for i := 0; i < rv.Len(); i++ {
			if reflect.DeepEqual(a, rv.Index(i).Interface()) {
				return true, nil
			}
		}
		return false, nil
	}

	return false, errors.New("unsupported operator")
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case float64:
		return n, true
	case float32:
		return float64(n), true
	default:
		return 0, false
	}
}
