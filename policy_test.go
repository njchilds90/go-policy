package policy

import (
	"context"
	"testing"
)

func TestEvaluate_BasicComparisons(t *testing.T) {
	tests := []struct {
		name   string
		rule   Rule
		input  map[string]any
		expect bool
	}{
		{
			name: "equal true",
			rule: Rule{
				Operator: OpEqual,
				Field:    "role",
				Value:    "admin",
			},
			input:  map[string]any{"role": "admin"},
			expect: true,
		},
		{
			name: "greater than true",
			rule: Rule{
				Operator: OpGreaterThan,
				Field:    "age",
				Value:    18,
			},
			input:  map[string]any{"age": 21},
			expect: true,
		},
		{
			name: "and composite",
			rule: Rule{
				Operator: OpAnd,
				Rules: []Rule{
					{
						Operator: OpGreaterThan,
						Field:    "age",
						Value:    18,
					},
					{
						Operator: OpEqual,
						Field:    "role",
						Value:    "admin",
					},
				},
			},
			input:  map[string]any{"age": 25, "role": "admin"},
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Evaluate(context.Background(), tt.rule, tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.Allowed != tt.expect {
				t.Fatalf("expected %v, got %v", tt.expect, res.Allowed)
			}
		})
	}
}
