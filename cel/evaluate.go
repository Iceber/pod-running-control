package cel

import (
	"context"
	"fmt"
	"time"

	"github.com/google/cel-go/interpreter"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/admission/plugin/cel"
	pkgcel "k8s.io/apiserver/pkg/cel"
)

func Evaluate(ctx context.Context, object runtime.Object, compilationResult cel.CompilationResult) (cel.EvaluationResult, error) {
	o, err := runtime.DefaultUnstructuredConverter.ToUnstructured(object)
	if err != nil {
		return cel.EvaluationResult{}, err
	}

	activation := &evaluationActivation{object: o}
	return activation.Evaluate(ctx, compilationResult)
}

type evaluationActivation struct{ object, params any }

// ResolveName returns a value from the activation by qualified name, or false if the name
// could not be found.
func (a *evaluationActivation) ResolveName(name string) (interface{}, bool) {
	switch name {
	case cel.ObjectVarName:
		return a.object, true
	default:
		return nil, false
	}
}

// Parent returns the parent of the current activation, may be nil.
// If non-nil, the parent will be searched during resolve calls.
func (a *evaluationActivation) Parent() interpreter.Activation {
	return nil
}

func (a *evaluationActivation) Evaluate(ctx context.Context, compilationResult cel.CompilationResult) (cel.EvaluationResult, error) {
	var evaluation = cel.EvaluationResult{}
	if compilationResult.ExpressionAccessor == nil { // in case of placeholder
		return evaluation, nil
	}

	evaluation.ExpressionAccessor = compilationResult.ExpressionAccessor
	if compilationResult.Error != nil {
		evaluation.Error = &pkgcel.Error{
			Type:   pkgcel.ErrorTypeInvalid,
			Detail: fmt.Sprintf("compilation error: %v", compilationResult.Error),
			Cause:  compilationResult.Error,
		}
		return evaluation, nil
	}
	if compilationResult.Program == nil {
		evaluation.Error = &pkgcel.Error{
			Type:   pkgcel.ErrorTypeInternal,
			Detail: "unexpected internal error compiling expression",
		}
		return evaluation, nil
	}

	t1 := time.Now()
	evalResult, evalDetails, err := compilationResult.Program.ContextEval(ctx, a)
	elapsed := time.Since(t1)
	evaluation.Elapsed = elapsed
	if evalDetails == nil {
		return evaluation, &pkgcel.Error{
			Type:   pkgcel.ErrorTypeInternal,
			Detail: fmt.Sprintf("runtime cost could not be calculated for expression: %v, no further expression will be run", compilationResult.ExpressionAccessor.GetExpression()),
		}
	}
	if err != nil {
		evaluation.Error = &pkgcel.Error{
			Type:   pkgcel.ErrorTypeInvalid,
			Detail: fmt.Sprintf("expression '%v' resulted in error: %v", compilationResult.ExpressionAccessor.GetExpression(), err),
		}
	} else {
		evaluation.EvalResult = evalResult
	}
	return evaluation, nil
}
