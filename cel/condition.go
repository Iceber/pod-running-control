package cel

import (
	celgo "github.com/google/cel-go/cel"
	"k8s.io/apiserver/pkg/admission/plugin/cel"
)

type ValidationCondition struct {
	Expression string
}

var _ cel.ExpressionAccessor = &ValidationCondition{}

func (v *ValidationCondition) GetExpression() string {
	return v.Expression
}
func (v *ValidationCondition) ReturnTypes() []*celgo.Type {
	return []*celgo.Type{celgo.BoolType}
}
