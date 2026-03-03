# go-policy

Deterministic, zero-dependency policy and rule evaluation engine for Go.

## Why

Declarative rule evaluation without embedding a scripting language.

Designed for:
- Feature flags
- Access control
- Input validation
- AI agent guardrails
- Workflow branching

## Install

go get github.com/njchilds90/go-policy

## Example

```go
package main

import (
	"context"
	"fmt"

	"github.com/njchilds90/go-policy"
)

func main() {
	rule := policy.Rule{
		Operator: policy.OpAnd,
		Rules: []policy.Rule{
			{
				Operator: policy.OpGreaterThan,
				Field:    "age",
				Value:    18,
			},
			{
				Operator: policy.OpEqual,
				Field:    "role",
				Value:    "admin",
			},
		},
	}

	input := map[string]any{
		"age":  25,
		"role": "admin",
	}

	result, err := policy.Evaluate(context.Background(), rule, input)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.Allowed) // true
}
