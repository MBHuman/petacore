// Code generated from sql.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // sql

import "github.com/antlr4-go/antlr/v4"

// sqlListener is a complete listener for a parse tree produced by sqlParser.
type sqlListener interface {
	antlr.ParseTreeListener

	// EnterQuery is called when entering the query production.
	EnterQuery(c *QueryContext)

	// EnterStatement is called when entering the statement production.
	EnterStatement(c *StatementContext)

	// EnterCreateTableStatement is called when entering the createTableStatement production.
	EnterCreateTableStatement(c *CreateTableStatementContext)

	// EnterColumnDefinition is called when entering the columnDefinition production.
	EnterColumnDefinition(c *ColumnDefinitionContext)

	// EnterColumnConstraints is called when entering the columnConstraints production.
	EnterColumnConstraints(c *ColumnConstraintsContext)

	// EnterDataType is called when entering the dataType production.
	EnterDataType(c *DataTypeContext)

	// EnterInsertStatement is called when entering the insertStatement production.
	EnterInsertStatement(c *InsertStatementContext)

	// EnterColumnList is called when entering the columnList production.
	EnterColumnList(c *ColumnListContext)

	// EnterValueList is called when entering the valueList production.
	EnterValueList(c *ValueListContext)

	// EnterDropTableStatement is called when entering the dropTableStatement production.
	EnterDropTableStatement(c *DropTableStatementContext)

	// EnterTruncateTableStatement is called when entering the truncateTableStatement production.
	EnterTruncateTableStatement(c *TruncateTableStatementContext)

	// EnterSetStatement is called when entering the setStatement production.
	EnterSetStatement(c *SetStatementContext)

	// EnterDescribeStatement is called when entering the describeStatement production.
	EnterDescribeStatement(c *DescribeStatementContext)

	// EnterShowStatement is called when entering the showStatement production.
	EnterShowStatement(c *ShowStatementContext)

	// EnterSelectStatement is called when entering the selectStatement production.
	EnterSelectStatement(c *SelectStatementContext)

	// EnterSelectList is called when entering the selectList production.
	EnterSelectList(c *SelectListContext)

	// EnterSelectItem is called when entering the selectItem production.
	EnterSelectItem(c *SelectItemContext)

	// EnterFromClause is called when entering the fromClause production.
	EnterFromClause(c *FromClauseContext)

	// EnterJoinClause is called when entering the joinClause production.
	EnterJoinClause(c *JoinClauseContext)

	// EnterWhereClause is called when entering the whereClause production.
	EnterWhereClause(c *WhereClauseContext)

	// EnterOrderByClause is called when entering the orderByClause production.
	EnterOrderByClause(c *OrderByClauseContext)

	// EnterOrderByItem is called when entering the orderByItem production.
	EnterOrderByItem(c *OrderByItemContext)

	// EnterLimitValue is called when entering the limitValue production.
	EnterLimitValue(c *LimitValueContext)

	// EnterOffsetValue is called when entering the offsetValue production.
	EnterOffsetValue(c *OffsetValueContext)

	// EnterAlias is called when entering the alias production.
	EnterAlias(c *AliasContext)

	// EnterExpression is called when entering the expression production.
	EnterExpression(c *ExpressionContext)

	// EnterOrExpression is called when entering the orExpression production.
	EnterOrExpression(c *OrExpressionContext)

	// EnterAndExpression is called when entering the andExpression production.
	EnterAndExpression(c *AndExpressionContext)

	// EnterNotExpression is called when entering the notExpression production.
	EnterNotExpression(c *NotExpressionContext)

	// EnterComparisonExpression is called when entering the comparisonExpression production.
	EnterComparisonExpression(c *ComparisonExpressionContext)

	// EnterConcatExpression is called when entering the concatExpression production.
	EnterConcatExpression(c *ConcatExpressionContext)

	// EnterAdditiveExpression is called when entering the additiveExpression production.
	EnterAdditiveExpression(c *AdditiveExpressionContext)

	// EnterMultiplicativeExpression is called when entering the multiplicativeExpression production.
	EnterMultiplicativeExpression(c *MultiplicativeExpressionContext)

	// EnterCastExpression is called when entering the castExpression production.
	EnterCastExpression(c *CastExpressionContext)

	// EnterPostfix is called when entering the postfix production.
	EnterPostfix(c *PostfixContext)

	// EnterTypeName is called when entering the typeName production.
	EnterTypeName(c *TypeNameContext)

	// EnterPrimaryExpression is called when entering the primaryExpression production.
	EnterPrimaryExpression(c *PrimaryExpressionContext)

	// EnterCaseExpression is called when entering the caseExpression production.
	EnterCaseExpression(c *CaseExpressionContext)

	// EnterFunctionCall is called when entering the functionCall production.
	EnterFunctionCall(c *FunctionCallContext)

	// EnterExtractFunction is called when entering the extractFunction production.
	EnterExtractFunction(c *ExtractFunctionContext)

	// EnterNamePart is called when entering the namePart production.
	EnterNamePart(c *NamePartContext)

	// EnterQualifiedName is called when entering the qualifiedName production.
	EnterQualifiedName(c *QualifiedNameContext)

	// EnterColumnName is called when entering the columnName production.
	EnterColumnName(c *ColumnNameContext)

	// EnterTableName is called when entering the tableName production.
	EnterTableName(c *TableNameContext)

	// EnterOperator is called when entering the operator production.
	EnterOperator(c *OperatorContext)

	// EnterOperatorExpr is called when entering the operatorExpr production.
	EnterOperatorExpr(c *OperatorExprContext)

	// EnterValue is called when entering the value production.
	EnterValue(c *ValueContext)

	// ExitQuery is called when exiting the query production.
	ExitQuery(c *QueryContext)

	// ExitStatement is called when exiting the statement production.
	ExitStatement(c *StatementContext)

	// ExitCreateTableStatement is called when exiting the createTableStatement production.
	ExitCreateTableStatement(c *CreateTableStatementContext)

	// ExitColumnDefinition is called when exiting the columnDefinition production.
	ExitColumnDefinition(c *ColumnDefinitionContext)

	// ExitColumnConstraints is called when exiting the columnConstraints production.
	ExitColumnConstraints(c *ColumnConstraintsContext)

	// ExitDataType is called when exiting the dataType production.
	ExitDataType(c *DataTypeContext)

	// ExitInsertStatement is called when exiting the insertStatement production.
	ExitInsertStatement(c *InsertStatementContext)

	// ExitColumnList is called when exiting the columnList production.
	ExitColumnList(c *ColumnListContext)

	// ExitValueList is called when exiting the valueList production.
	ExitValueList(c *ValueListContext)

	// ExitDropTableStatement is called when exiting the dropTableStatement production.
	ExitDropTableStatement(c *DropTableStatementContext)

	// ExitTruncateTableStatement is called when exiting the truncateTableStatement production.
	ExitTruncateTableStatement(c *TruncateTableStatementContext)

	// ExitSetStatement is called when exiting the setStatement production.
	ExitSetStatement(c *SetStatementContext)

	// ExitDescribeStatement is called when exiting the describeStatement production.
	ExitDescribeStatement(c *DescribeStatementContext)

	// ExitShowStatement is called when exiting the showStatement production.
	ExitShowStatement(c *ShowStatementContext)

	// ExitSelectStatement is called when exiting the selectStatement production.
	ExitSelectStatement(c *SelectStatementContext)

	// ExitSelectList is called when exiting the selectList production.
	ExitSelectList(c *SelectListContext)

	// ExitSelectItem is called when exiting the selectItem production.
	ExitSelectItem(c *SelectItemContext)

	// ExitFromClause is called when exiting the fromClause production.
	ExitFromClause(c *FromClauseContext)

	// ExitJoinClause is called when exiting the joinClause production.
	ExitJoinClause(c *JoinClauseContext)

	// ExitWhereClause is called when exiting the whereClause production.
	ExitWhereClause(c *WhereClauseContext)

	// ExitOrderByClause is called when exiting the orderByClause production.
	ExitOrderByClause(c *OrderByClauseContext)

	// ExitOrderByItem is called when exiting the orderByItem production.
	ExitOrderByItem(c *OrderByItemContext)

	// ExitLimitValue is called when exiting the limitValue production.
	ExitLimitValue(c *LimitValueContext)

	// ExitOffsetValue is called when exiting the offsetValue production.
	ExitOffsetValue(c *OffsetValueContext)

	// ExitAlias is called when exiting the alias production.
	ExitAlias(c *AliasContext)

	// ExitExpression is called when exiting the expression production.
	ExitExpression(c *ExpressionContext)

	// ExitOrExpression is called when exiting the orExpression production.
	ExitOrExpression(c *OrExpressionContext)

	// ExitAndExpression is called when exiting the andExpression production.
	ExitAndExpression(c *AndExpressionContext)

	// ExitNotExpression is called when exiting the notExpression production.
	ExitNotExpression(c *NotExpressionContext)

	// ExitComparisonExpression is called when exiting the comparisonExpression production.
	ExitComparisonExpression(c *ComparisonExpressionContext)

	// ExitConcatExpression is called when exiting the concatExpression production.
	ExitConcatExpression(c *ConcatExpressionContext)

	// ExitAdditiveExpression is called when exiting the additiveExpression production.
	ExitAdditiveExpression(c *AdditiveExpressionContext)

	// ExitMultiplicativeExpression is called when exiting the multiplicativeExpression production.
	ExitMultiplicativeExpression(c *MultiplicativeExpressionContext)

	// ExitCastExpression is called when exiting the castExpression production.
	ExitCastExpression(c *CastExpressionContext)

	// ExitPostfix is called when exiting the postfix production.
	ExitPostfix(c *PostfixContext)

	// ExitTypeName is called when exiting the typeName production.
	ExitTypeName(c *TypeNameContext)

	// ExitPrimaryExpression is called when exiting the primaryExpression production.
	ExitPrimaryExpression(c *PrimaryExpressionContext)

	// ExitCaseExpression is called when exiting the caseExpression production.
	ExitCaseExpression(c *CaseExpressionContext)

	// ExitFunctionCall is called when exiting the functionCall production.
	ExitFunctionCall(c *FunctionCallContext)

	// ExitExtractFunction is called when exiting the extractFunction production.
	ExitExtractFunction(c *ExtractFunctionContext)

	// ExitNamePart is called when exiting the namePart production.
	ExitNamePart(c *NamePartContext)

	// ExitQualifiedName is called when exiting the qualifiedName production.
	ExitQualifiedName(c *QualifiedNameContext)

	// ExitColumnName is called when exiting the columnName production.
	ExitColumnName(c *ColumnNameContext)

	// ExitTableName is called when exiting the tableName production.
	ExitTableName(c *TableNameContext)

	// ExitOperator is called when exiting the operator production.
	ExitOperator(c *OperatorContext)

	// ExitOperatorExpr is called when exiting the operatorExpr production.
	ExitOperatorExpr(c *OperatorExprContext)

	// ExitValue is called when exiting the value production.
	ExitValue(c *ValueContext)
}
