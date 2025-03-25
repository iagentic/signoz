package parser

import (
	"fmt"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

// ClickHouseVisitor implements the FilterQueryVisitor interface
// to convert the parsed filter expressions into ClickHouse WHERE clauses
type ClickHouseVisitor struct {
	BaseFilterQueryVisitor
}

// NewClickHouseVisitor creates a new ClickHouseVisitor
func NewClickHouseVisitor() *ClickHouseVisitor {
	return &ClickHouseVisitor{}
}

// ErrorListener is a custom error listener to capture syntax errors
type ErrorListener struct {
	*antlr.DefaultErrorListener
	Errors []string
}

// NewErrorListener creates a new error listener
func NewErrorListener() *ErrorListener {
	return &ErrorListener{
		DefaultErrorListener: antlr.NewDefaultErrorListener(),
		Errors:               []string{},
	}
}

// SyntaxError captures syntax errors during parsing
func (l *ErrorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	fmt.Println("SyntaxError", msg, "line", line, "column", column, "offendingSymbol", offendingSymbol, "e", e)
	symbol := offendingSymbol.(*antlr.CommonToken)
	fmt.Println("symbol", symbol.GetText(), "symbol.GetTokenType()", symbol.GetTokenType(), "start", symbol.GetStart(), "stop", symbol.GetStop())
	error := fmt.Sprintf("line %d:%d %s", line, column, msg)
	l.Errors = append(l.Errors, error)
}

// PrepareWhereClause generates a ClickHouse compatible WHERE clause from the filter query
func PrepareWhereClause(query string) (string, error) {
	// Setup the ANTLR parsing pipeline
	input := antlr.NewInputStream(query)
	lexer := NewFilterQueryLexer(input)

	// Set up error handling
	errorListener := NewErrorListener()
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)

	tokens := antlr.NewCommonTokenStream(lexer, 0)
	parser := NewFilterQueryParser(tokens)
	parser.RemoveErrorListeners()
	parser.AddErrorListener(errorListener)

	// Parse the query
	tree := parser.Query()

	// Handle syntax errors
	if len(errorListener.Errors) > 0 {
		return "", fmt.Errorf("syntax error in filter query: %s", strings.Join(errorListener.Errors, "; "))
	}

	// Visit the parse tree with our ClickHouse visitor
	visitor := NewClickHouseVisitor()
	whereClause := visitor.Visit(tree)

	// Convert result to string, handling nil cases
	if whereClause == nil {
		return "", fmt.Errorf("failed to generate WHERE clause for query: %s", query)
	}

	return whereClause.(string), nil
}

// Visit dispatches to the specific visit method based on node type
func (v *ClickHouseVisitor) Visit(tree antlr.ParseTree) interface{} {
	// Handle nil nodes to prevent panic
	if tree == nil {
		return ""
	}

	switch t := tree.(type) {
	case *QueryContext:
		return v.VisitQuery(t)
	case *ExpressionContext:
		return v.VisitExpression(t)
	case *OrExpressionContext:
		return v.VisitOrExpression(t)
	case *AndExpressionContext:
		return v.VisitAndExpression(t)
	case *UnaryExpressionContext:
		return v.VisitUnaryExpression(t)
	case *PrimaryContext:
		return v.VisitPrimary(t)
	case *ComparisonContext:
		return v.VisitComparison(t)
	case *InClauseContext:
		return v.VisitInClause(t)
	case *NotInClauseContext:
		return v.VisitNotInClause(t)
	case *ValueListContext:
		return v.VisitValueList(t)
	case *FullTextContext:
		return v.VisitFullText(t)
	case *FunctionCallContext:
		return v.VisitFunctionCall(t)
	case *FunctionParamListContext:
		return v.VisitFunctionParamList(t)
	case *FunctionParamContext:
		return v.VisitFunctionParam(t)
	case *ArrayContext:
		return v.VisitArray(t)
	case *ValueContext:
		return v.VisitValue(t)
	case *KeyContext:
		return v.VisitKey(t)
	default:
		return ""
	}
}

// VisitQuery handles the root query node
func (v *ClickHouseVisitor) VisitQuery(ctx *QueryContext) interface{} {
	expressions := ctx.AllExpression()
	if len(expressions) == 0 {
		return ""
	}

	// Visit the first expression
	result := v.Visit(expressions[0]).(string)

	// Process any additional expressions (with implicit/explicit AND/OR)
	for i := 1; i < len(expressions); i++ {
		// Check if there's an AND/OR token between this expression and the previous one
		var op string
		// Check the token type at the position preceding the current expression
		// This requires examining the tokens in the input stream
		if i < len(expressions) && i > 0 {
			// Check for explicit operators in the original query
			// This is a simplification - in a real implementation, we'd need to examine tokens
			// For now, we'll default to AND for simplicity
			op = " AND "
		} else {
			// Implicit AND
			op = " AND "
		}

		exprResult := v.Visit(expressions[i]).(string)
		result = fmt.Sprintf("(%s)%s(%s)", result, op, exprResult)
	}

	return result
}

// VisitExpression passes through to the orExpression
func (v *ClickHouseVisitor) VisitExpression(ctx *ExpressionContext) interface{} {
	return v.Visit(ctx.OrExpression())
}

// VisitOrExpression handles OR expressions
func (v *ClickHouseVisitor) VisitOrExpression(ctx *OrExpressionContext) interface{} {
	andExpressions := ctx.AllAndExpression()
	if len(andExpressions) == 1 {
		return v.Visit(andExpressions[0])
	}

	parts := make([]string, len(andExpressions))
	for i, expr := range andExpressions {
		parts[i] = v.Visit(expr).(string)
	}

	return strings.Join(parts, " OR ")
}

// VisitAndExpression handles AND expressions
func (v *ClickHouseVisitor) VisitAndExpression(ctx *AndExpressionContext) interface{} {
	unaryExpressions := ctx.AllUnaryExpression()
	if len(unaryExpressions) == 1 {
		return v.Visit(unaryExpressions[0])
	}

	parts := make([]string, len(unaryExpressions))
	for i, expr := range unaryExpressions {
		parts[i] = v.Visit(expr).(string)
	}

	return strings.Join(parts, " AND ")
}

// VisitUnaryExpression handles NOT expressions
func (v *ClickHouseVisitor) VisitUnaryExpression(ctx *UnaryExpressionContext) interface{} {
	result := v.Visit(ctx.Primary()).(string)

	// Check if this is a NOT expression
	if ctx.NOT() != nil {
		return fmt.Sprintf("NOT (%s)", result)
	}

	return result
}

// VisitPrimary handles grouped expressions, comparisons, function calls, and full-text search
func (v *ClickHouseVisitor) VisitPrimary(ctx *PrimaryContext) interface{} {
	if ctx.OrExpression() != nil {
		// This is a parenthesized expression
		return fmt.Sprintf("(%s)", v.Visit(ctx.OrExpression()).(string))
	} else if ctx.Comparison() != nil {
		return v.Visit(ctx.Comparison())
	} else if ctx.FunctionCall() != nil {
		return v.Visit(ctx.FunctionCall())
	} else if ctx.FullText() != nil {
		return v.Visit(ctx.FullText())
	}

	return "" // Should not happen with valid input
}

// VisitComparison handles all comparison operators
func (v *ClickHouseVisitor) VisitComparison(ctx *ComparisonContext) interface{} {
	keyName := v.Visit(ctx.Key()).(string)

	// Handle EXISTS specially
	if ctx.EXISTS() != nil {
		if ctx.NOT() != nil {
			return fmt.Sprintf("not has(%s)", keyName)
		}
		return fmt.Sprintf("has(%s)", keyName)
	}

	// Handle IN clause
	if ctx.InClause() != nil {
		inClause := v.Visit(ctx.InClause()).(string)
		return fmt.Sprintf("%s %s", keyName, inClause)
	}

	// Handle NOT IN clause
	if ctx.NotInClause() != nil {
		notInClause := v.Visit(ctx.NotInClause()).(string)
		return fmt.Sprintf("%s %s", keyName, notInClause)
	}

	// Get all values for operations that need them
	values := ctx.AllValue()
	if len(values) > 0 {
		value := v.Visit(values[0]).(string)

		// Handle each type of comparison
		if ctx.EQUALS() != nil {
			return fmt.Sprintf("%s = %s", keyName, value)
		} else if ctx.NOT_EQUALS() != nil || ctx.NEQ() != nil {
			return fmt.Sprintf("%s != %s", keyName, value)
		} else if ctx.LT() != nil {
			return fmt.Sprintf("%s < %s", keyName, value)
		} else if ctx.LE() != nil {
			return fmt.Sprintf("%s <= %s", keyName, value)
		} else if ctx.GT() != nil {
			return fmt.Sprintf("%s > %s", keyName, value)
		} else if ctx.GE() != nil {
			return fmt.Sprintf("%s >= %s", keyName, value)
		} else if ctx.LIKE() != nil {
			// Convert SQL LIKE to ClickHouse LIKE
			return fmt.Sprintf("%s LIKE %s", keyName, value)
		} else if ctx.ILIKE() != nil {
			// ClickHouse has ilike
			return fmt.Sprintf("%s ILIKE %s", keyName, value)
		} else if ctx.NOT_LIKE() != nil {
			return fmt.Sprintf("%s NOT LIKE %s", keyName, value)
		} else if ctx.NOT_ILIKE() != nil {
			return fmt.Sprintf("%s NOT ILIKE %s", keyName, value)
		} else if ctx.REGEXP() != nil {
			// ClickHouse uses match for regex
			return fmt.Sprintf("match(%s, %s)", keyName, value)
		} else if ctx.NOT() != nil && ctx.REGEXP() != nil {
			return fmt.Sprintf("not match(%s, %s)", keyName, value)
		} else if ctx.CONTAINS() != nil {
			// In ClickHouse, we can use position or like
			return fmt.Sprintf("position(%s, %s) > 0", keyName, value)
		} else if ctx.NOT() != nil && ctx.CONTAINS() != nil {
			return fmt.Sprintf("position(%s, %s) = 0", keyName, value)
		}
	}

	// Handle BETWEEN
	if ctx.BETWEEN() != nil && len(values) >= 2 {
		value1 := v.Visit(values[0]).(string)
		value2 := v.Visit(values[1]).(string)

		if ctx.NOT() != nil {
			return fmt.Sprintf("(%s < %s OR %s > %s)", keyName, value1, keyName, value2)
		}
		return fmt.Sprintf("%s BETWEEN %s AND %s", keyName, value1, value2)
	}

	return "" // Should not happen with valid input
}

// VisitInClause handles IN expressions
func (v *ClickHouseVisitor) VisitInClause(ctx *InClauseContext) interface{} {
	values := v.Visit(ctx.ValueList()).(string)
	return fmt.Sprintf("IN (%s)", values)
}

// VisitNotInClause handles NOT IN expressions
func (v *ClickHouseVisitor) VisitNotInClause(ctx *NotInClauseContext) interface{} {
	values := v.Visit(ctx.ValueList()).(string)
	return fmt.Sprintf("NOT IN (%s)", values)
}

// VisitValueList handles comma-separated value lists
func (v *ClickHouseVisitor) VisitValueList(ctx *ValueListContext) interface{} {
	values := ctx.AllValue()
	if len(values) == 0 {
		return ""
	}

	parts := make([]string, len(values))
	for i, val := range values {
		parts[i] = v.Visit(val).(string)
	}

	return strings.Join(parts, ", ")
}

// VisitFullText handles standalone quoted strings for full-text search
func (v *ClickHouseVisitor) VisitFullText(ctx *FullTextContext) interface{} {
	// remove quotes from the quotedText
	quotedText := strings.Trim(ctx.QUOTED_TEXT().GetText(), "\"'")
	return fmt.Sprintf("lower(body) LIKE '%%%s%%' OR has(mapValues(attributes_string, '%%%s%%'))", strings.ToLower(quotedText), quotedText)
}

// VisitFunctionCall handles function calls like has(), hasAny(), etc.
func (v *ClickHouseVisitor) VisitFunctionCall(ctx *FunctionCallContext) interface{} {
	// Get function name based on which token is present
	var functionName string
	if ctx.HAS() != nil {
		functionName = "has"
	} else if ctx.HASANY() != nil {
		functionName = "hasany"
	} else if ctx.HASALL() != nil {
		functionName = "hasall"
	} else if ctx.HASNONE() != nil {
		functionName = "hasnone"
	} else {
		// Default fallback
		functionName = "unknown_function"
	}
	params := v.Visit(ctx.FunctionParamList()).(string)

	// Map our functions to ClickHouse equivalents
	switch functionName {
	case "has":
		return fmt.Sprintf("has(%s)", params)
	case "hasany":
		return fmt.Sprintf("hasAny(%s)", params)
	case "hasall":
		return fmt.Sprintf("hasAll(%s)", params)
	case "hasnone":
		// ClickHouse doesn't have hasNone directly, so we negate hasAny
		return fmt.Sprintf("not hasAny(%s)", params)
	default:
		return fmt.Sprintf("%s(%s)", functionName, params)
	}
}

// VisitFunctionParamList handles the parameter list for function calls
func (v *ClickHouseVisitor) VisitFunctionParamList(ctx *FunctionParamListContext) interface{} {
	params := ctx.AllFunctionParam()
	if len(params) == 0 {
		return ""
	}

	parts := make([]string, len(params))
	for i, param := range params {
		parts[i] = v.Visit(param).(string)
	}

	return strings.Join(parts, ", ")
}

// VisitFunctionParam handles individual parameters in function calls
func (v *ClickHouseVisitor) VisitFunctionParam(ctx *FunctionParamContext) interface{} {
	if ctx.Key() != nil {
		return v.Visit(ctx.Key())
	} else if ctx.Value() != nil {
		return v.Visit(ctx.Value())
	} else if ctx.Array() != nil {
		return v.Visit(ctx.Array())
	}

	return "" // Should not happen with valid input
}

// VisitArray handles array literals
func (v *ClickHouseVisitor) VisitArray(ctx *ArrayContext) interface{} {
	values := v.Visit(ctx.ValueList()).(string)
	return fmt.Sprintf("[%s]", values)
}

// VisitValue handles literal values: strings, numbers, booleans
func (v *ClickHouseVisitor) VisitValue(ctx *ValueContext) interface{} {
	if ctx.QUOTED_TEXT() != nil {
		return ctx.QUOTED_TEXT().GetText()
	} else if ctx.NUMBER() != nil {
		return ctx.NUMBER().GetText()
	} else if ctx.BOOL() != nil {
		// Convert to ClickHouse boolean literal
		boolText := strings.ToLower(ctx.BOOL().GetText())
		if boolText == "true" {
			return "1"
		}
		return "0"
	}

	return "" // Should not happen with valid input
}

// VisitKey handles field/column references
func (v *ClickHouseVisitor) VisitKey(ctx *KeyContext) interface{} {

	keyText := ctx.KEY().GetText()

	return keyText
}
