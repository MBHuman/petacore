package rparser

import (
	"fmt"
	"petacore/internal/logger"
	"petacore/internal/runtime/parser"
	"petacore/internal/runtime/rsql/statements"
)

// ParseFromClause парсит FROM clause
func ParseFromClause(ctx parser.IFromClauseContext) (*statements.FromClause, error) {
	if ctx == nil {
		return nil, nil
	}

	from := &statements.FromClause{}

	// Handle multiple table factors separated by commas.
	// The first factor becomes the main table; subsequent factors are converted
	// into implicit CROSS joins so that downstream planner sees all tables.
	tableFactors := ctx.AllTableFactor()
	if len(tableFactors) > 0 {
		for idx, tf := range tableFactors {
			if idx == 0 {
				if tableName := tf.TableName(); tableName != nil {
					from.TableName = tableName.GetText()
					if tf.Alias() != nil {
						from.Alias = tf.Alias().GetText()
					}
				}

				if subSelect := tf.SelectStatement(); subSelect != nil {
					logger.Debugf("DEBUG: parsing FROM subquery")
					if selectStmt, ok := subSelect.(*parser.SelectStatementContext); ok {
						var err error
						from.SelectStatement, err = ParseSelectStatement(selectStmt)
						if err != nil {
							return nil, fmt.Errorf("Error parsing subquery in FROM clause: %v", err)
						}
						if tf.Alias() != nil {
							from.Alias = tf.Alias().GetText()
						}
					}
				}
				continue
			}

			// For subsequent table factors, create implicit CROSS joins
			join := statements.JoinClause{Type: "CROSS"}
			if tableName := tf.TableName(); tableName != nil {
				join.TableName = tableName.GetText()
				if tf.Alias() != nil {
					join.Alias = tf.Alias().GetText()
				}
			}
			// Note: we don't support comma-separated subselects here yet
			if subSelect := tf.SelectStatement(); subSelect != nil {
				logger.Debugf("DEBUG: encountered subquery as additional tableFactor — treating as unnamed subquery join not yet supported")
			}

			from.Joins = append(from.Joins, join)
		}
	}

	// Parse joins
	for _, jc := range ctx.AllJoinClause() {
		join := statements.JoinClause{}
		if jc.INNER() != nil {
			join.Type = "INNER"
		} else if jc.LEFT() != nil {
			join.Type = "LEFT"
		} else if jc.RIGHT() != nil {
			join.Type = "RIGHT"
		} else if jc.FULL() != nil {
			join.Type = "FULL"
		} else if jc.CROSS() != nil {
			join.Type = "CROSS"
		} else {
			join.Type = "INNER" // default
		}
		if jc.QualifiedName() != nil {
			join.TableName = jc.QualifiedName().GetText()
		}
		if alias := jc.Alias(); alias != nil {
			if id := alias.IDENTIFIER(); id != nil {
				join.Alias = id.GetText()
			}
		}
		if jc.ON() != nil && jc.Expression() != nil {
			join.OnCondition = jc.Expression()
		}
		from.Joins = append(from.Joins, join)
	}

	return from, nil
}
