// Code generated from sql.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // sql

import "github.com/antlr4-go/antlr/v4"

// BasesqlListener is a complete listener for a parse tree produced by sqlParser.
type BasesqlListener struct{}

var _ sqlListener = &BasesqlListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasesqlListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasesqlListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasesqlListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasesqlListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterQuery is called when production query is entered.
func (s *BasesqlListener) EnterQuery(ctx *QueryContext) {}

// ExitQuery is called when production query is exited.
func (s *BasesqlListener) ExitQuery(ctx *QueryContext) {}

// EnterStatement is called when production statement is entered.
func (s *BasesqlListener) EnterStatement(ctx *StatementContext) {}

// ExitStatement is called when production statement is exited.
func (s *BasesqlListener) ExitStatement(ctx *StatementContext) {}

// EnterCreateTableStatement is called when production createTableStatement is entered.
func (s *BasesqlListener) EnterCreateTableStatement(ctx *CreateTableStatementContext) {}

// ExitCreateTableStatement is called when production createTableStatement is exited.
func (s *BasesqlListener) ExitCreateTableStatement(ctx *CreateTableStatementContext) {}

// EnterColumnDefinition is called when production columnDefinition is entered.
func (s *BasesqlListener) EnterColumnDefinition(ctx *ColumnDefinitionContext) {}

// ExitColumnDefinition is called when production columnDefinition is exited.
func (s *BasesqlListener) ExitColumnDefinition(ctx *ColumnDefinitionContext) {}

// EnterColumnConstraints is called when production columnConstraints is entered.
func (s *BasesqlListener) EnterColumnConstraints(ctx *ColumnConstraintsContext) {}

// ExitColumnConstraints is called when production columnConstraints is exited.
func (s *BasesqlListener) ExitColumnConstraints(ctx *ColumnConstraintsContext) {}

// EnterDataType is called when production dataType is entered.
func (s *BasesqlListener) EnterDataType(ctx *DataTypeContext) {}

// ExitDataType is called when production dataType is exited.
func (s *BasesqlListener) ExitDataType(ctx *DataTypeContext) {}

// EnterInsertStatement is called when production insertStatement is entered.
func (s *BasesqlListener) EnterInsertStatement(ctx *InsertStatementContext) {}

// ExitInsertStatement is called when production insertStatement is exited.
func (s *BasesqlListener) ExitInsertStatement(ctx *InsertStatementContext) {}

// EnterColumnList is called when production columnList is entered.
func (s *BasesqlListener) EnterColumnList(ctx *ColumnListContext) {}

// ExitColumnList is called when production columnList is exited.
func (s *BasesqlListener) ExitColumnList(ctx *ColumnListContext) {}

// EnterValueList is called when production valueList is entered.
func (s *BasesqlListener) EnterValueList(ctx *ValueListContext) {}

// ExitValueList is called when production valueList is exited.
func (s *BasesqlListener) ExitValueList(ctx *ValueListContext) {}

// EnterDropTableStatement is called when production dropTableStatement is entered.
func (s *BasesqlListener) EnterDropTableStatement(ctx *DropTableStatementContext) {}

// ExitDropTableStatement is called when production dropTableStatement is exited.
func (s *BasesqlListener) ExitDropTableStatement(ctx *DropTableStatementContext) {}

// EnterTruncateTableStatement is called when production truncateTableStatement is entered.
func (s *BasesqlListener) EnterTruncateTableStatement(ctx *TruncateTableStatementContext) {}

// ExitTruncateTableStatement is called when production truncateTableStatement is exited.
func (s *BasesqlListener) ExitTruncateTableStatement(ctx *TruncateTableStatementContext) {}

// EnterSetStatement is called when production setStatement is entered.
func (s *BasesqlListener) EnterSetStatement(ctx *SetStatementContext) {}

// ExitSetStatement is called when production setStatement is exited.
func (s *BasesqlListener) ExitSetStatement(ctx *SetStatementContext) {}

// EnterDescribeStatement is called when production describeStatement is entered.
func (s *BasesqlListener) EnterDescribeStatement(ctx *DescribeStatementContext) {}

// ExitDescribeStatement is called when production describeStatement is exited.
func (s *BasesqlListener) ExitDescribeStatement(ctx *DescribeStatementContext) {}

// EnterShowStatement is called when production showStatement is entered.
func (s *BasesqlListener) EnterShowStatement(ctx *ShowStatementContext) {}

// ExitShowStatement is called when production showStatement is exited.
func (s *BasesqlListener) ExitShowStatement(ctx *ShowStatementContext) {}

// EnterSelectStatement is called when production selectStatement is entered.
func (s *BasesqlListener) EnterSelectStatement(ctx *SelectStatementContext) {}

// ExitSelectStatement is called when production selectStatement is exited.
func (s *BasesqlListener) ExitSelectStatement(ctx *SelectStatementContext) {}

// EnterSelectList is called when production selectList is entered.
func (s *BasesqlListener) EnterSelectList(ctx *SelectListContext) {}

// ExitSelectList is called when production selectList is exited.
func (s *BasesqlListener) ExitSelectList(ctx *SelectListContext) {}

// EnterSelectItem is called when production selectItem is entered.
func (s *BasesqlListener) EnterSelectItem(ctx *SelectItemContext) {}

// ExitSelectItem is called when production selectItem is exited.
func (s *BasesqlListener) ExitSelectItem(ctx *SelectItemContext) {}

// EnterFromClause is called when production fromClause is entered.
func (s *BasesqlListener) EnterFromClause(ctx *FromClauseContext) {}

// ExitFromClause is called when production fromClause is exited.
func (s *BasesqlListener) ExitFromClause(ctx *FromClauseContext) {}

// EnterJoinClause is called when production joinClause is entered.
func (s *BasesqlListener) EnterJoinClause(ctx *JoinClauseContext) {}

// ExitJoinClause is called when production joinClause is exited.
func (s *BasesqlListener) ExitJoinClause(ctx *JoinClauseContext) {}

// EnterWhereClause is called when production whereClause is entered.
func (s *BasesqlListener) EnterWhereClause(ctx *WhereClauseContext) {}

// ExitWhereClause is called when production whereClause is exited.
func (s *BasesqlListener) ExitWhereClause(ctx *WhereClauseContext) {}

// EnterOrderByClause is called when production orderByClause is entered.
func (s *BasesqlListener) EnterOrderByClause(ctx *OrderByClauseContext) {}

// ExitOrderByClause is called when production orderByClause is exited.
func (s *BasesqlListener) ExitOrderByClause(ctx *OrderByClauseContext) {}

// EnterOrderByItem is called when production orderByItem is entered.
func (s *BasesqlListener) EnterOrderByItem(ctx *OrderByItemContext) {}

// ExitOrderByItem is called when production orderByItem is exited.
func (s *BasesqlListener) ExitOrderByItem(ctx *OrderByItemContext) {}

// EnterLimitValue is called when production limitValue is entered.
func (s *BasesqlListener) EnterLimitValue(ctx *LimitValueContext) {}

// ExitLimitValue is called when production limitValue is exited.
func (s *BasesqlListener) ExitLimitValue(ctx *LimitValueContext) {}

// EnterOffsetValue is called when production offsetValue is entered.
func (s *BasesqlListener) EnterOffsetValue(ctx *OffsetValueContext) {}

// ExitOffsetValue is called when production offsetValue is exited.
func (s *BasesqlListener) ExitOffsetValue(ctx *OffsetValueContext) {}

// EnterAlias is called when production alias is entered.
func (s *BasesqlListener) EnterAlias(ctx *AliasContext) {}

// ExitAlias is called when production alias is exited.
func (s *BasesqlListener) ExitAlias(ctx *AliasContext) {}

// EnterExpression is called when production expression is entered.
func (s *BasesqlListener) EnterExpression(ctx *ExpressionContext) {}

// ExitExpression is called when production expression is exited.
func (s *BasesqlListener) ExitExpression(ctx *ExpressionContext) {}

// EnterOrExpression is called when production orExpression is entered.
func (s *BasesqlListener) EnterOrExpression(ctx *OrExpressionContext) {}

// ExitOrExpression is called when production orExpression is exited.
func (s *BasesqlListener) ExitOrExpression(ctx *OrExpressionContext) {}

// EnterAndExpression is called when production andExpression is entered.
func (s *BasesqlListener) EnterAndExpression(ctx *AndExpressionContext) {}

// ExitAndExpression is called when production andExpression is exited.
func (s *BasesqlListener) ExitAndExpression(ctx *AndExpressionContext) {}

// EnterNotExpression is called when production notExpression is entered.
func (s *BasesqlListener) EnterNotExpression(ctx *NotExpressionContext) {}

// ExitNotExpression is called when production notExpression is exited.
func (s *BasesqlListener) ExitNotExpression(ctx *NotExpressionContext) {}

// EnterComparisonExpression is called when production comparisonExpression is entered.
func (s *BasesqlListener) EnterComparisonExpression(ctx *ComparisonExpressionContext) {}

// ExitComparisonExpression is called when production comparisonExpression is exited.
func (s *BasesqlListener) ExitComparisonExpression(ctx *ComparisonExpressionContext) {}

// EnterConcatExpression is called when production concatExpression is entered.
func (s *BasesqlListener) EnterConcatExpression(ctx *ConcatExpressionContext) {}

// ExitConcatExpression is called when production concatExpression is exited.
func (s *BasesqlListener) ExitConcatExpression(ctx *ConcatExpressionContext) {}

// EnterAdditiveExpression is called when production additiveExpression is entered.
func (s *BasesqlListener) EnterAdditiveExpression(ctx *AdditiveExpressionContext) {}

// ExitAdditiveExpression is called when production additiveExpression is exited.
func (s *BasesqlListener) ExitAdditiveExpression(ctx *AdditiveExpressionContext) {}

// EnterMultiplicativeExpression is called when production multiplicativeExpression is entered.
func (s *BasesqlListener) EnterMultiplicativeExpression(ctx *MultiplicativeExpressionContext) {}

// ExitMultiplicativeExpression is called when production multiplicativeExpression is exited.
func (s *BasesqlListener) ExitMultiplicativeExpression(ctx *MultiplicativeExpressionContext) {}

// EnterCastExpression is called when production castExpression is entered.
func (s *BasesqlListener) EnterCastExpression(ctx *CastExpressionContext) {}

// ExitCastExpression is called when production castExpression is exited.
func (s *BasesqlListener) ExitCastExpression(ctx *CastExpressionContext) {}

// EnterAtTimeZoneExpression is called when production atTimeZoneExpression is entered.
func (s *BasesqlListener) EnterAtTimeZoneExpression(ctx *AtTimeZoneExpressionContext) {}

// ExitAtTimeZoneExpression is called when production atTimeZoneExpression is exited.
func (s *BasesqlListener) ExitAtTimeZoneExpression(ctx *AtTimeZoneExpressionContext) {}

// EnterPrimaryExpression is called when production primaryExpression is entered.
func (s *BasesqlListener) EnterPrimaryExpression(ctx *PrimaryExpressionContext) {}

// ExitPrimaryExpression is called when production primaryExpression is exited.
func (s *BasesqlListener) ExitPrimaryExpression(ctx *PrimaryExpressionContext) {}

// EnterCaseExpression is called when production caseExpression is entered.
func (s *BasesqlListener) EnterCaseExpression(ctx *CaseExpressionContext) {}

// ExitCaseExpression is called when production caseExpression is exited.
func (s *BasesqlListener) ExitCaseExpression(ctx *CaseExpressionContext) {}

// EnterFunctionCall is called when production functionCall is entered.
func (s *BasesqlListener) EnterFunctionCall(ctx *FunctionCallContext) {}

// ExitFunctionCall is called when production functionCall is exited.
func (s *BasesqlListener) ExitFunctionCall(ctx *FunctionCallContext) {}

// EnterExtractFunction is called when production extractFunction is entered.
func (s *BasesqlListener) EnterExtractFunction(ctx *ExtractFunctionContext) {}

// ExitExtractFunction is called when production extractFunction is exited.
func (s *BasesqlListener) ExitExtractFunction(ctx *ExtractFunctionContext) {}

// EnterQualifiedName is called when production qualifiedName is entered.
func (s *BasesqlListener) EnterQualifiedName(ctx *QualifiedNameContext) {}

// ExitQualifiedName is called when production qualifiedName is exited.
func (s *BasesqlListener) ExitQualifiedName(ctx *QualifiedNameContext) {}

// EnterColumnName is called when production columnName is entered.
func (s *BasesqlListener) EnterColumnName(ctx *ColumnNameContext) {}

// ExitColumnName is called when production columnName is exited.
func (s *BasesqlListener) ExitColumnName(ctx *ColumnNameContext) {}

// EnterTableName is called when production tableName is entered.
func (s *BasesqlListener) EnterTableName(ctx *TableNameContext) {}

// ExitTableName is called when production tableName is exited.
func (s *BasesqlListener) ExitTableName(ctx *TableNameContext) {}

// EnterOperator is called when production operator is entered.
func (s *BasesqlListener) EnterOperator(ctx *OperatorContext) {}

// ExitOperator is called when production operator is exited.
func (s *BasesqlListener) ExitOperator(ctx *OperatorContext) {}

// EnterValue is called when production value is entered.
func (s *BasesqlListener) EnterValue(ctx *ValueContext) {}

// ExitValue is called when production value is exited.
func (s *BasesqlListener) ExitValue(ctx *ValueContext) {}
