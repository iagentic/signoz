package parser

import (
	"context"
	"fmt"
	"strings"

	"github.com/SigNoz/signoz/pkg/types"
	"github.com/antlr4-go/antlr/v4"
)

// ClickHouseVisitor implements the FilterQueryVisitor interface
// to convert the parsed filter expressions into ClickHouse WHERE clauses
type ClickHouseVisitor struct {
	metadataStore types.Metadata
	columnMapper  types.KeyToColumnMapper
}

type partialQuery struct {
	paritalQuery string
	args         []any
	warnings     []string
	errors       []string
}

// NewClickHouseVisitor creates a new ClickHouseVisitor
func NewClickHouseVisitor(
	metadataStore types.Metadata,
	columnMapper types.KeyToColumnMapper,
) *ClickHouseVisitor {
	return &ClickHouseVisitor{
		columnMapper: columnMapper,
	}
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

func getFieldSelectorFromKey(key string) types.FieldKeySelector {

	keyTextParts := strings.Split(key, ".")

	var explicitFieldContextProvided, explicitFieldDataTypeProvided bool
	var explicitFieldContext types.FieldContext
	var explicitFieldDataType types.FieldDataType

	if len(keyTextParts) > 1 {
		explicitFieldContext = types.FieldContextFromString(keyTextParts[0])
		if explicitFieldContext != types.FieldContextAll {
			explicitFieldContextProvided = true
		}
	}

	if explicitFieldContextProvided {
		keyTextParts = keyTextParts[1:]
	}

	// check if there is a field data type provided
	if len(keyTextParts) > 1 {
		lastPart := keyTextParts[len(keyTextParts)-1]
		lastPartParts := strings.Split(lastPart, ":")
		if len(lastPartParts) > 1 {
			explicitFieldDataType = types.FieldDataTypeFromString(lastPartParts[1])
			if explicitFieldDataType != types.FieldDataTypeAll {
				explicitFieldDataTypeProvided = true
			}
		}

		if explicitFieldDataTypeProvided {
			keyTextParts[len(keyTextParts)-1] = lastPartParts[0]
		}
	}

	realKey := strings.Join(keyTextParts, ".")

	fieldKeySelector := types.FieldKeySelector{
		Name: realKey,
	}

	if explicitFieldContextProvided {
		fieldKeySelector.FieldContext = explicitFieldContext
	} else {
		fieldKeySelector.FieldContext = types.FieldContextAll
	}

	if explicitFieldDataTypeProvided {
		fieldKeySelector.FieldDataType = explicitFieldDataType
	} else {
		fieldKeySelector.FieldDataType = types.FieldDataTypeAll
	}

	return fieldKeySelector
}

// PrepareWhereClause generates a ClickHouse compatible WHERE clause from the filter query
func PrepareWhereClause(query string, visitor *ClickHouseVisitor) (string, error) {
	// Setup the ANTLR parsing pipeline
	input := antlr.NewInputStream(query)
	lexer := NewFilterQueryLexer(input)

	lexerForKeysEnrich := NewFilterQueryLexer(input)
	fieldKeySelectors := []types.FieldKeySelector{}
	for {
		tok := lexerForKeysEnrich.NextToken()
		if tok.GetTokenType() == antlr.TokenEOF {
			break
		}
		if tok.GetTokenType() == FilterQueryLexerKEY {
			fieldKeySelectors = append(fieldKeySelectors, getFieldSelectorFromKey(tok.GetText()))
		}
	}

	fieldKeys, err := visitor.metadataStore.GetKeysMulti(context.Background(), fieldKeySelectors)
	if err != nil {
		return "", fmt.Errorf("error getting field keys: %s", err)
	}

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
	whereClause := visitor.Visit(tree, fieldKeys)

	// Convert result to string, handling nil cases
	if whereClause == nil {
		return "", fmt.Errorf("failed to generate WHERE clause for query: %s", query)
	}

	return whereClause.(string), nil
}

// Visit dispatches to the specific visit method based on node type
func (v *ClickHouseVisitor) Visit(tree antlr.ParseTree, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
	// Handle nil nodes to prevent panic
	if tree == nil {
		return ""
	}

	switch t := tree.(type) {
	case *QueryContext:
		fmt.Println("QueryContext")
		return v.VisitQuery(t, fieldKeys)
	case *ExpressionContext:
		fmt.Println("ExpressionContext")
		return v.VisitExpression(t, fieldKeys)
	case *OrExpressionContext:
		fmt.Println("OrExpressionContext")
		return v.VisitOrExpression(t, fieldKeys)
	case *AndExpressionContext:
		fmt.Println("AndExpressionContext")
		return v.VisitAndExpression(t, fieldKeys)
	case *UnaryExpressionContext:
		fmt.Println("UnaryExpressionContext")
		return v.VisitUnaryExpression(t, fieldKeys)
	case *PrimaryContext:
		fmt.Println("PrimaryContext")
		return v.VisitPrimary(t, fieldKeys)
	case *ComparisonContext:
		fmt.Println("ComparisonContext")
		return v.VisitComparison(t, fieldKeys)
	case *InClauseContext:
		fmt.Println("InClauseContext")
		return v.VisitInClause(t, fieldKeys)
	case *NotInClauseContext:
		fmt.Println("NotInClauseContext")
		return v.VisitNotInClause(t, fieldKeys)
	case *ValueListContext:
		fmt.Println("ValueListContext")
		return v.VisitValueList(t, fieldKeys)
	case *FullTextContext:
		fmt.Println("FullTextContext")
		return v.VisitFullText(t, fieldKeys)
	case *FunctionCallContext:
		fmt.Println("FunctionCallContext")
		return v.VisitFunctionCall(t, fieldKeys)
	case *FunctionParamListContext:
		fmt.Println("FunctionParamListContext")
		return v.VisitFunctionParamList(t, fieldKeys)
	case *FunctionParamContext:
		fmt.Println("FunctionParamContext")
		return v.VisitFunctionParam(t, fieldKeys)
	case *ArrayContext:
		fmt.Println("ArrayContext")
		return v.VisitArray(t, fieldKeys)
	case *ValueContext:
		val := v.VisitValue(t, fieldKeys)
		return val
	case *KeyContext:
		fmt.Println("KeyContext")
		return v.VisitKey(t, fieldKeys)
	default:
		fmt.Println("default")
		return ""
	}
}

// VisitQuery handles the root query node
func (v *ClickHouseVisitor) VisitQuery(ctx *QueryContext, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
	expressions := ctx.AllExpression()
	if len(expressions) == 0 {
		return ""
	}

	// Visit the first expression
	result := v.Visit(expressions[0], fieldKeys).(string)

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

		exprResult := v.Visit(expressions[i], fieldKeys).(string)
		result = fmt.Sprintf("(%s)%s(%s)", result, op, exprResult)
	}

	return result
}

// VisitExpression passes through to the orExpression
func (v *ClickHouseVisitor) VisitExpression(ctx *ExpressionContext, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
	return v.Visit(ctx.OrExpression(), fieldKeys)
}

// VisitOrExpression handles OR expressions
func (v *ClickHouseVisitor) VisitOrExpression(ctx *OrExpressionContext, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
	andExpressions := ctx.AllAndExpression()
	if len(andExpressions) == 1 {
		return v.Visit(andExpressions[0], fieldKeys)
	}

	parts := make([]string, len(andExpressions))
	for i, expr := range andExpressions {
		parts[i] = v.Visit(expr, fieldKeys).(string)
	}

	return strings.Join(parts, " OR ")
}

// VisitAndExpression handles AND expressions
func (v *ClickHouseVisitor) VisitAndExpression(ctx *AndExpressionContext, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
	unaryExpressions := ctx.AllUnaryExpression()
	if len(unaryExpressions) == 1 {
		return v.Visit(unaryExpressions[0], fieldKeys)
	}

	parts := make([]string, len(unaryExpressions))
	for i, expr := range unaryExpressions {
		parts[i] = v.Visit(expr, fieldKeys).(string)
	}

	return strings.Join(parts, " AND ")
}

// VisitUnaryExpression handles NOT expressions
func (v *ClickHouseVisitor) VisitUnaryExpression(ctx *UnaryExpressionContext, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
	result := v.Visit(ctx.Primary(), fieldKeys).(string)

	// Check if this is a NOT expression
	if ctx.NOT() != nil {
		return fmt.Sprintf("NOT (%s)", result)
	}

	return result
}

// VisitPrimary handles grouped expressions, comparisons, function calls, and full-text search
func (v *ClickHouseVisitor) VisitPrimary(ctx *PrimaryContext, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
	if ctx.OrExpression() != nil {
		// This is a parenthesized expression
		return fmt.Sprintf("(%s)", v.Visit(ctx.OrExpression(), fieldKeys).(string))
	} else if ctx.Comparison() != nil {
		return v.Visit(ctx.Comparison(), fieldKeys)
	} else if ctx.FunctionCall() != nil {
		return v.Visit(ctx.FunctionCall(), fieldKeys)
	} else if ctx.FullText() != nil {
		return v.Visit(ctx.FullText(), fieldKeys)
	}

	fmt.Println("ctx.GetChildCount()", ctx.GetChildCount())
	// Handle standalone key as a full text search term
	if ctx.GetChildCount() == 1 {
		child := ctx.GetChild(0)
		if keyCtx, ok := child.(*KeyContext); ok {
			// create a full text search condition on the body field
			keyText := keyCtx.GetText()
			return fmt.Sprintf("lower(body) LIKE '%%%s%%' OR has(mapValues(attributes_string, '%%%s%%'))", strings.ToLower(keyText), keyText)
		}
	}

	return "" // Should not happen with valid input
}

// VisitComparison handles all comparison operators
func (v *ClickHouseVisitor) VisitComparison(ctx *ComparisonContext, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
	keyName := v.Visit(ctx.Key(), fieldKeys).(string)

	// Handle EXISTS specially
	if ctx.EXISTS() != nil {
		if ctx.NOT() != nil {
			return partialQuery{
				paritalQuery: fmt.Sprintf("not has(%s)", keyName),
				args:         []any{},
			}
		}
		return partialQuery{
			paritalQuery: fmt.Sprintf("has(%s)", keyName),
			args:         []any{},
		}
	}

	// Handle IN clause
	if ctx.InClause() != nil {
		parts := v.Visit(ctx.InClause(), fieldKeys).([]any)
		partsQuery := ""

		for idx := range parts {
			partsQuery += "?"
			if idx < len(parts)-1 {
				partsQuery += ","
			}
		}
		return fmt.Sprintf("%s IN (%s)", keyName, partsQuery)
	}

	// Handle NOT IN clause
	if ctx.NotInClause() != nil {
		notInClause := v.Visit(ctx.NotInClause(), fieldKeys).(string)
		return fmt.Sprintf("%s %s", keyName, notInClause)
	}

	// Get all values for operations that need them
	values := ctx.AllValue()
	if len(values) > 0 {
		value := v.Visit(values[0], fieldKeys).(string)

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
		value1 := v.Visit(values[0], fieldKeys).(string)
		value2 := v.Visit(values[1], fieldKeys).(string)

		if ctx.NOT() != nil {
			return fmt.Sprintf("(%s < %s OR %s > %s)", keyName, value1, keyName, value2)
		}
		return fmt.Sprintf("%s BETWEEN %s AND %s", keyName, value1, value2)
	}

	return "" // Should not happen with valid input
}

// VisitInClause handles IN expressions
func (v *ClickHouseVisitor) VisitInClause(ctx *InClauseContext, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
	values := v.Visit(ctx.ValueList(), fieldKeys).(string)
	return fmt.Sprintf("IN (%s)", values)
}

// VisitNotInClause handles NOT IN expressions
func (v *ClickHouseVisitor) VisitNotInClause(ctx *NotInClauseContext, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
	values := v.Visit(ctx.ValueList(), fieldKeys).(string)
	return fmt.Sprintf("NOT IN (%s)", values)
}

// VisitValueList handles comma-separated value lists
func (v *ClickHouseVisitor) VisitValueList(ctx *ValueListContext, fieldKeys map[string][]types.TelemetryFieldKey) any {
	values := ctx.AllValue()
	if len(values) == 0 {
		return ""
	}

	parts := []any{}
	for _, val := range values {
		parts = append(parts, v.Visit(val, fieldKeys))
	}

	return parts
}

// VisitFullText handles standalone quoted strings for full-text search
func (v *ClickHouseVisitor) VisitFullText(ctx *FullTextContext, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
	// remove quotes from the quotedText
	quotedText := strings.Trim(ctx.QUOTED_TEXT().GetText(), "\"'")
	return fmt.Sprintf("lower(body) LIKE '%%%s%%' OR has(mapValues(attributes_string, '%%%s%%'))", strings.ToLower(quotedText), quotedText)
}

// VisitFunctionCall handles function calls like has(), hasAny(), etc.
func (v *ClickHouseVisitor) VisitFunctionCall(ctx *FunctionCallContext, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
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
	params := v.Visit(ctx.FunctionParamList(), fieldKeys).(string)

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
func (v *ClickHouseVisitor) VisitFunctionParamList(ctx *FunctionParamListContext, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
	params := ctx.AllFunctionParam()
	if len(params) == 0 {
		return ""
	}

	parts := make([]string, len(params))
	for i, param := range params {
		parts[i] = v.Visit(param, fieldKeys).(string)
	}

	return strings.Join(parts, ", ")
}

// VisitFunctionParam handles individual parameters in function calls
func (v *ClickHouseVisitor) VisitFunctionParam(ctx *FunctionParamContext, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
	if ctx.Key() != nil {
		return v.Visit(ctx.Key(), fieldKeys)
	} else if ctx.Value() != nil {
		return v.Visit(ctx.Value(), fieldKeys)
	} else if ctx.Array() != nil {
		return v.Visit(ctx.Array(), fieldKeys)
	}

	return "" // Should not happen with valid input
}

// VisitArray handles array literals
func (v *ClickHouseVisitor) VisitArray(ctx *ArrayContext, fieldKeys map[string][]types.TelemetryFieldKey) interface{} {
	values := v.Visit(ctx.ValueList(), fieldKeys).(string)
	return fmt.Sprintf("[%s]", values)
}

// VisitValue handles literal values: strings, numbers, booleans
func (v *ClickHouseVisitor) VisitValue(ctx *ValueContext, fieldKeys map[string][]types.TelemetryFieldKey) any {
	if ctx.QUOTED_TEXT() != nil {
		txt := ctx.QUOTED_TEXT().GetText()
		return strings.Trim(txt, "\"'")
	} else if ctx.NUMBER() != nil {
		return ctx.NUMBER().GetText()
	} else if ctx.BOOL() != nil {
		// Convert to ClickHouse boolean literal
		boolText := strings.ToLower(ctx.BOOL().GetText())
		return boolText == "true"
	} else if ctx.KEY() != nil {
		return ctx.KEY().GetText()
	}

	return "" // Should not happen with valid input
}

// VisitKey handles field/column references
func (v *ClickHouseVisitor) VisitKey(ctx *KeyContext, fieldKeys map[string][]types.TelemetryFieldKey) any {

	fieldKeySelector := getFieldSelectorFromKey(ctx.KEY().GetText())

	fieldKeysForName := fieldKeys[fieldKeySelector.Name]

	if len(fieldKeysForName) == 0 {
		return fmt.Errorf("Key %s not found", fieldKeySelector.Name)
	}

	if len(fieldKeysForName) > 1 {
		// this is error state, we must have a unambiguous key
		return fmt.Errorf("Key %s is ambiguous, found %d different combinations of field context and data type", fieldKeySelector.Name, len(fieldKeysForName))
	}

	return fieldKeysForName[0]
}
