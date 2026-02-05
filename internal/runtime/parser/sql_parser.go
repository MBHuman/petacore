// Code generated from sql.g4 by ANTLR 4.13.2. DO NOT EDIT.

package parser // sql

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr4-go/antlr/v4"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type sqlParser struct {
	*antlr.BaseParser
}

var SqlParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	LiteralNames           []string
	SymbolicNames          []string
	RuleNames              []string
	PredictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func sqlParserInit() {
	staticData := &SqlParserStaticData
	staticData.LiteralNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "'('", "')'", "';'", "','",
		"'.'", "'*'", "'+'", "'-'", "'/'", "'||'", "'::'", "'~'", "'!~'", "'~*'",
		"'!~*'", "'='", "'>'", "'<'", "'>='", "'<='",
	}
	staticData.SymbolicNames = []string{
		"", "CREATE", "TABLE", "INSERT", "INTO", "VALUES", "PRIMARY", "KEY",
		"DROP", "TRUNCATE", "SET", "TO", "IF", "NOT", "EXISTS", "NULL", "UNIQUE",
		"DEFAULT", "SHOW", "SELECT", "FROM", "WHERE", "GROUP", "LIMIT", "OFFSET",
		"ORDER", "BY", "ASC", "DESC", "UNION", "INTERSECT", "EXCEPT", "ALL",
		"JOIN", "INNER", "LEFT", "RIGHT", "FULL", "CROSS", "ON", "IN", "AND",
		"OR", "IS", "LIKE", "ILIKE", "CASE", "WHEN", "THEN", "ELSE", "END",
		"AS", "AT", "TIME", "ZONE", "EXTRACT", "STRING_TYPE", "INT_TYPE", "FLOAT_TYPE",
		"BOOL_TYPE", "TEXT_TYPE", "VARCHAR_TYPE", "CHAR_TYPE", "SERIAL_TYPE",
		"TIMESTAMP_TYPE", "DATE_TYPE", "INTERVAL_TYPE", "YEAR", "MONTH", "DAY",
		"HOUR", "MINUTE", "SECOND", "WITH", "WITHOUT", "CURRENT_TIMESTAMP",
		"CURRENT_USER", "TRUE", "FALSE", "LPAREN", "RPAREN", "SEMICOLON", "COMMA",
		"DOT", "STAR", "PLUS", "MINUS", "SLASH", "CONCAT", "COLON_COLON", "TILDE",
		"NREGEX", "IREGEX", "NIREGEX", "EQ", "GT", "LT", "GE", "LE", "NE", "OPERATOR_KW",
		"COLLATE", "IDENTIFIER", "PARAMETER", "STRING_LITERAL", "NUMBER", "BLOCK_COMMENT",
		"LINE_COMMENT", "WS",
	}
	staticData.RuleNames = []string{
		"query", "statement", "createTableStatement", "columnDefinition", "columnConstraints",
		"dataType", "intervalFields", "insertStatement", "columnList", "valueList",
		"dropTableStatement", "truncateTableStatement", "setStatement", "showStatement",
		"selectStatement", "unionExceptStatement", "intersectStatement", "primarySelectStatement",
		"selectList", "selectItem", "selectAll", "fromClause", "tableFactor",
		"joinClause", "whereClause", "orderByClause", "orderByItem", "limitValue",
		"offsetValue", "groupByClause", "alias", "expression", "orExpression",
		"andExpression", "notExpression", "comparisonExpression", "concatExpression",
		"additiveExpression", "multiplicativeExpression", "unaryExpression",
		"castExpression", "postfix", "typeName", "primaryExpression", "subqueryExpression",
		"caseExpression", "functionCall", "functionArg", "extractFunction",
		"namePart", "qualifiedName", "columnName", "tableName", "operator",
		"operatorExpr", "value", "typedLiteral",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 108, 702, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2,
		42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47,
		7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 2, 50, 7, 50, 2, 51, 7, 51, 2, 52, 7,
		52, 2, 53, 7, 53, 2, 54, 7, 54, 2, 55, 7, 55, 2, 56, 7, 56, 1, 0, 1, 0,
		3, 0, 117, 8, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 126, 8,
		1, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 3, 2, 133, 8, 2, 1, 2, 1, 2, 1, 2, 1,
		2, 1, 2, 5, 2, 140, 8, 2, 10, 2, 12, 2, 143, 9, 2, 1, 2, 3, 2, 146, 8,
		2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 5, 2, 154, 8, 2, 10, 2, 12, 2, 157,
		9, 2, 1, 2, 1, 2, 3, 2, 161, 8, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 3,
		1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 3, 3, 175, 8, 3, 1, 4, 1, 4, 1, 4,
		1, 4, 1, 4, 5, 4, 182, 8, 4, 10, 4, 12, 4, 185, 9, 4, 1, 5, 1, 5, 1, 5,
		1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 196, 8, 5, 1, 5, 1, 5, 1, 5,
		1, 5, 3, 5, 202, 8, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 209, 8, 5, 1,
		5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 217, 8, 5, 1, 5, 1, 5, 1, 5, 1,
		5, 1, 5, 3, 5, 224, 8, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 232,
		8, 5, 1, 5, 1, 5, 3, 5, 236, 8, 5, 1, 5, 1, 5, 1, 5, 3, 5, 241, 8, 5, 3,
		5, 243, 8, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1,
		6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1,
		6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 3, 6, 272, 8, 6, 1, 7, 1, 7, 1, 7, 1,
		7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 5, 7, 284, 8, 7, 10, 7, 12, 7, 287,
		9, 7, 1, 8, 1, 8, 1, 8, 5, 8, 292, 8, 8, 10, 8, 12, 8, 295, 9, 8, 1, 9,
		1, 9, 1, 9, 1, 9, 5, 9, 301, 8, 9, 10, 9, 12, 9, 304, 9, 9, 1, 9, 1, 9,
		1, 10, 1, 10, 1, 10, 1, 10, 3, 10, 312, 8, 10, 1, 10, 1, 10, 1, 11, 1,
		11, 3, 11, 318, 8, 11, 1, 11, 1, 11, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12,
		1, 13, 1, 13, 1, 13, 5, 13, 330, 8, 13, 10, 13, 12, 13, 333, 9, 13, 1,
		14, 1, 14, 1, 15, 1, 15, 1, 15, 3, 15, 340, 8, 15, 1, 15, 5, 15, 343, 8,
		15, 10, 15, 12, 15, 346, 9, 15, 1, 16, 1, 16, 1, 16, 3, 16, 351, 8, 16,
		1, 16, 5, 16, 354, 8, 16, 10, 16, 12, 16, 357, 9, 16, 1, 17, 1, 17, 1,
		17, 1, 17, 3, 17, 363, 8, 17, 1, 17, 1, 17, 3, 17, 367, 8, 17, 1, 17, 1,
		17, 1, 17, 3, 17, 372, 8, 17, 1, 17, 1, 17, 1, 17, 3, 17, 377, 8, 17, 1,
		17, 1, 17, 3, 17, 381, 8, 17, 1, 17, 1, 17, 3, 17, 385, 8, 17, 1, 17, 1,
		17, 1, 17, 1, 17, 3, 17, 391, 8, 17, 1, 18, 1, 18, 1, 18, 5, 18, 396, 8,
		18, 10, 18, 12, 18, 399, 9, 18, 1, 19, 1, 19, 3, 19, 403, 8, 19, 1, 19,
		3, 19, 406, 8, 19, 1, 20, 1, 20, 1, 21, 1, 21, 1, 21, 5, 21, 413, 8, 21,
		10, 21, 12, 21, 416, 9, 21, 1, 21, 5, 21, 419, 8, 21, 10, 21, 12, 21, 422,
		9, 21, 1, 22, 1, 22, 3, 22, 426, 8, 22, 1, 22, 1, 22, 1, 22, 1, 22, 1,
		22, 3, 22, 433, 8, 22, 1, 23, 3, 23, 436, 8, 23, 1, 23, 1, 23, 1, 23, 3,
		23, 441, 8, 23, 1, 23, 1, 23, 1, 23, 1, 24, 1, 24, 1, 25, 1, 25, 1, 25,
		5, 25, 451, 8, 25, 10, 25, 12, 25, 454, 9, 25, 1, 26, 1, 26, 3, 26, 458,
		8, 26, 1, 27, 1, 27, 1, 28, 1, 28, 1, 29, 1, 29, 1, 29, 5, 29, 467, 8,
		29, 10, 29, 12, 29, 470, 9, 29, 1, 30, 3, 30, 473, 8, 30, 1, 30, 1, 30,
		1, 31, 1, 31, 1, 32, 1, 32, 1, 32, 5, 32, 482, 8, 32, 10, 32, 12, 32, 485,
		9, 32, 1, 33, 1, 33, 1, 33, 5, 33, 490, 8, 33, 10, 33, 12, 33, 493, 9,
		33, 1, 34, 3, 34, 496, 8, 34, 1, 34, 1, 34, 3, 34, 500, 8, 34, 1, 34, 1,
		34, 3, 34, 504, 8, 34, 1, 35, 1, 35, 1, 35, 3, 35, 509, 8, 35, 1, 35, 1,
		35, 1, 35, 3, 35, 514, 8, 35, 1, 35, 1, 35, 1, 35, 1, 35, 1, 35, 5, 35,
		521, 8, 35, 10, 35, 12, 35, 524, 9, 35, 1, 35, 1, 35, 1, 35, 3, 35, 529,
		8, 35, 1, 35, 1, 35, 1, 35, 3, 35, 534, 8, 35, 1, 35, 1, 35, 1, 35, 1,
		35, 3, 35, 540, 8, 35, 1, 35, 3, 35, 543, 8, 35, 1, 36, 1, 36, 1, 36, 5,
		36, 548, 8, 36, 10, 36, 12, 36, 551, 9, 36, 1, 37, 1, 37, 1, 37, 5, 37,
		556, 8, 37, 10, 37, 12, 37, 559, 9, 37, 1, 38, 1, 38, 1, 38, 5, 38, 564,
		8, 38, 10, 38, 12, 38, 567, 9, 38, 1, 39, 3, 39, 570, 8, 39, 1, 39, 1,
		39, 1, 40, 1, 40, 5, 40, 576, 8, 40, 10, 40, 12, 40, 579, 9, 40, 1, 41,
		1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 4, 41, 589, 8, 41, 11,
		41, 12, 41, 590, 3, 41, 593, 8, 41, 1, 42, 1, 42, 1, 43, 1, 43, 1, 43,
		1, 43, 1, 43, 1, 43, 1, 43, 1, 43, 1, 43, 1, 43, 1, 43, 3, 43, 608, 8,
		43, 1, 44, 1, 44, 1, 44, 1, 44, 1, 45, 1, 45, 1, 45, 1, 45, 1, 45, 1, 45,
		4, 45, 620, 8, 45, 11, 45, 12, 45, 621, 1, 45, 1, 45, 3, 45, 626, 8, 45,
		1, 45, 1, 45, 1, 46, 1, 46, 1, 46, 1, 46, 1, 46, 5, 46, 635, 8, 46, 10,
		46, 12, 46, 638, 9, 46, 3, 46, 640, 8, 46, 1, 46, 1, 46, 1, 47, 1, 47,
		3, 47, 646, 8, 47, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1, 48, 1,
		49, 1, 49, 1, 50, 1, 50, 1, 50, 1, 50, 3, 50, 661, 8, 50, 5, 50, 663, 8,
		50, 10, 50, 12, 50, 666, 9, 50, 1, 51, 1, 51, 1, 52, 1, 52, 1, 53, 1, 53,
		1, 54, 1, 54, 1, 54, 1, 54, 1, 54, 1, 55, 1, 55, 1, 55, 1, 55, 1, 55, 1,
		55, 1, 55, 1, 55, 3, 55, 687, 8, 55, 1, 56, 1, 56, 1, 56, 1, 56, 1, 56,
		1, 56, 1, 56, 1, 56, 1, 56, 3, 56, 698, 8, 56, 3, 56, 700, 8, 56, 1, 56,
		0, 0, 57, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32,
		34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68,
		70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98, 100, 102, 104,
		106, 108, 110, 112, 0, 8, 2, 0, 11, 11, 94, 94, 2, 0, 29, 29, 31, 31, 1,
		0, 34, 38, 1, 0, 27, 28, 1, 0, 85, 86, 2, 0, 84, 84, 87, 87, 4, 0, 17,
		17, 53, 53, 57, 66, 102, 102, 1, 0, 90, 99, 772, 0, 114, 1, 0, 0, 0, 2,
		125, 1, 0, 0, 0, 4, 127, 1, 0, 0, 0, 6, 174, 1, 0, 0, 0, 8, 183, 1, 0,
		0, 0, 10, 242, 1, 0, 0, 0, 12, 271, 1, 0, 0, 0, 14, 273, 1, 0, 0, 0, 16,
		288, 1, 0, 0, 0, 18, 296, 1, 0, 0, 0, 20, 307, 1, 0, 0, 0, 22, 315, 1,
		0, 0, 0, 24, 321, 1, 0, 0, 0, 26, 326, 1, 0, 0, 0, 28, 334, 1, 0, 0, 0,
		30, 336, 1, 0, 0, 0, 32, 347, 1, 0, 0, 0, 34, 390, 1, 0, 0, 0, 36, 392,
		1, 0, 0, 0, 38, 405, 1, 0, 0, 0, 40, 407, 1, 0, 0, 0, 42, 409, 1, 0, 0,
		0, 44, 432, 1, 0, 0, 0, 46, 435, 1, 0, 0, 0, 48, 445, 1, 0, 0, 0, 50, 447,
		1, 0, 0, 0, 52, 455, 1, 0, 0, 0, 54, 459, 1, 0, 0, 0, 56, 461, 1, 0, 0,
		0, 58, 463, 1, 0, 0, 0, 60, 472, 1, 0, 0, 0, 62, 476, 1, 0, 0, 0, 64, 478,
		1, 0, 0, 0, 66, 486, 1, 0, 0, 0, 68, 503, 1, 0, 0, 0, 70, 505, 1, 0, 0,
		0, 72, 544, 1, 0, 0, 0, 74, 552, 1, 0, 0, 0, 76, 560, 1, 0, 0, 0, 78, 569,
		1, 0, 0, 0, 80, 573, 1, 0, 0, 0, 82, 592, 1, 0, 0, 0, 84, 594, 1, 0, 0,
		0, 86, 607, 1, 0, 0, 0, 88, 609, 1, 0, 0, 0, 90, 613, 1, 0, 0, 0, 92, 629,
		1, 0, 0, 0, 94, 645, 1, 0, 0, 0, 96, 647, 1, 0, 0, 0, 98, 654, 1, 0, 0,
		0, 100, 656, 1, 0, 0, 0, 102, 667, 1, 0, 0, 0, 104, 669, 1, 0, 0, 0, 106,
		671, 1, 0, 0, 0, 108, 673, 1, 0, 0, 0, 110, 686, 1, 0, 0, 0, 112, 699,
		1, 0, 0, 0, 114, 116, 3, 2, 1, 0, 115, 117, 5, 81, 0, 0, 116, 115, 1, 0,
		0, 0, 116, 117, 1, 0, 0, 0, 117, 1, 1, 0, 0, 0, 118, 126, 3, 28, 14, 0,
		119, 126, 3, 4, 2, 0, 120, 126, 3, 14, 7, 0, 121, 126, 3, 20, 10, 0, 122,
		126, 3, 22, 11, 0, 123, 126, 3, 24, 12, 0, 124, 126, 3, 26, 13, 0, 125,
		118, 1, 0, 0, 0, 125, 119, 1, 0, 0, 0, 125, 120, 1, 0, 0, 0, 125, 121,
		1, 0, 0, 0, 125, 122, 1, 0, 0, 0, 125, 123, 1, 0, 0, 0, 125, 124, 1, 0,
		0, 0, 126, 3, 1, 0, 0, 0, 127, 128, 5, 1, 0, 0, 128, 132, 5, 2, 0, 0, 129,
		130, 5, 12, 0, 0, 130, 131, 5, 13, 0, 0, 131, 133, 5, 14, 0, 0, 132, 129,
		1, 0, 0, 0, 132, 133, 1, 0, 0, 0, 133, 134, 1, 0, 0, 0, 134, 135, 3, 104,
		52, 0, 135, 136, 5, 79, 0, 0, 136, 141, 3, 6, 3, 0, 137, 138, 5, 82, 0,
		0, 138, 140, 3, 6, 3, 0, 139, 137, 1, 0, 0, 0, 140, 143, 1, 0, 0, 0, 141,
		139, 1, 0, 0, 0, 141, 142, 1, 0, 0, 0, 142, 145, 1, 0, 0, 0, 143, 141,
		1, 0, 0, 0, 144, 146, 5, 82, 0, 0, 145, 144, 1, 0, 0, 0, 145, 146, 1, 0,
		0, 0, 146, 160, 1, 0, 0, 0, 147, 148, 5, 6, 0, 0, 148, 149, 5, 7, 0, 0,
		149, 150, 5, 79, 0, 0, 150, 155, 3, 102, 51, 0, 151, 152, 5, 82, 0, 0,
		152, 154, 3, 102, 51, 0, 153, 151, 1, 0, 0, 0, 154, 157, 1, 0, 0, 0, 155,
		153, 1, 0, 0, 0, 155, 156, 1, 0, 0, 0, 156, 158, 1, 0, 0, 0, 157, 155,
		1, 0, 0, 0, 158, 159, 5, 80, 0, 0, 159, 161, 1, 0, 0, 0, 160, 147, 1, 0,
		0, 0, 160, 161, 1, 0, 0, 0, 161, 162, 1, 0, 0, 0, 162, 163, 5, 80, 0, 0,
		163, 5, 1, 0, 0, 0, 164, 165, 3, 102, 51, 0, 165, 166, 3, 10, 5, 0, 166,
		167, 3, 8, 4, 0, 167, 168, 5, 6, 0, 0, 168, 169, 5, 7, 0, 0, 169, 175,
		1, 0, 0, 0, 170, 171, 3, 102, 51, 0, 171, 172, 3, 10, 5, 0, 172, 173, 3,
		8, 4, 0, 173, 175, 1, 0, 0, 0, 174, 164, 1, 0, 0, 0, 174, 170, 1, 0, 0,
		0, 175, 7, 1, 0, 0, 0, 176, 177, 5, 13, 0, 0, 177, 182, 5, 15, 0, 0, 178,
		182, 5, 16, 0, 0, 179, 180, 5, 17, 0, 0, 180, 182, 3, 110, 55, 0, 181,
		176, 1, 0, 0, 0, 181, 178, 1, 0, 0, 0, 181, 179, 1, 0, 0, 0, 182, 185,
		1, 0, 0, 0, 183, 181, 1, 0, 0, 0, 183, 184, 1, 0, 0, 0, 184, 9, 1, 0, 0,
		0, 185, 183, 1, 0, 0, 0, 186, 243, 5, 56, 0, 0, 187, 243, 5, 57, 0, 0,
		188, 243, 5, 58, 0, 0, 189, 243, 5, 59, 0, 0, 190, 243, 5, 60, 0, 0, 191,
		195, 5, 61, 0, 0, 192, 193, 5, 79, 0, 0, 193, 194, 5, 105, 0, 0, 194, 196,
		5, 80, 0, 0, 195, 192, 1, 0, 0, 0, 195, 196, 1, 0, 0, 0, 196, 243, 1, 0,
		0, 0, 197, 201, 5, 62, 0, 0, 198, 199, 5, 79, 0, 0, 199, 200, 5, 105, 0,
		0, 200, 202, 5, 80, 0, 0, 201, 198, 1, 0, 0, 0, 201, 202, 1, 0, 0, 0, 202,
		243, 1, 0, 0, 0, 203, 243, 5, 63, 0, 0, 204, 208, 5, 64, 0, 0, 205, 206,
		5, 79, 0, 0, 206, 207, 5, 105, 0, 0, 207, 209, 5, 80, 0, 0, 208, 205, 1,
		0, 0, 0, 208, 209, 1, 0, 0, 0, 209, 216, 1, 0, 0, 0, 210, 211, 5, 73, 0,
		0, 211, 212, 5, 53, 0, 0, 212, 217, 5, 54, 0, 0, 213, 214, 5, 74, 0, 0,
		214, 215, 5, 53, 0, 0, 215, 217, 5, 54, 0, 0, 216, 210, 1, 0, 0, 0, 216,
		213, 1, 0, 0, 0, 216, 217, 1, 0, 0, 0, 217, 243, 1, 0, 0, 0, 218, 243,
		5, 65, 0, 0, 219, 223, 5, 53, 0, 0, 220, 221, 5, 79, 0, 0, 221, 222, 5,
		105, 0, 0, 222, 224, 5, 80, 0, 0, 223, 220, 1, 0, 0, 0, 223, 224, 1, 0,
		0, 0, 224, 231, 1, 0, 0, 0, 225, 226, 5, 73, 0, 0, 226, 227, 5, 53, 0,
		0, 227, 232, 5, 54, 0, 0, 228, 229, 5, 74, 0, 0, 229, 230, 5, 53, 0, 0,
		230, 232, 5, 54, 0, 0, 231, 225, 1, 0, 0, 0, 231, 228, 1, 0, 0, 0, 231,
		232, 1, 0, 0, 0, 232, 243, 1, 0, 0, 0, 233, 235, 5, 66, 0, 0, 234, 236,
		3, 12, 6, 0, 235, 234, 1, 0, 0, 0, 235, 236, 1, 0, 0, 0, 236, 240, 1, 0,
		0, 0, 237, 238, 5, 79, 0, 0, 238, 239, 5, 105, 0, 0, 239, 241, 5, 80, 0,
		0, 240, 237, 1, 0, 0, 0, 240, 241, 1, 0, 0, 0, 241, 243, 1, 0, 0, 0, 242,
		186, 1, 0, 0, 0, 242, 187, 1, 0, 0, 0, 242, 188, 1, 0, 0, 0, 242, 189,
		1, 0, 0, 0, 242, 190, 1, 0, 0, 0, 242, 191, 1, 0, 0, 0, 242, 197, 1, 0,
		0, 0, 242, 203, 1, 0, 0, 0, 242, 204, 1, 0, 0, 0, 242, 218, 1, 0, 0, 0,
		242, 219, 1, 0, 0, 0, 242, 233, 1, 0, 0, 0, 243, 11, 1, 0, 0, 0, 244, 272,
		5, 67, 0, 0, 245, 272, 5, 68, 0, 0, 246, 272, 5, 69, 0, 0, 247, 272, 5,
		70, 0, 0, 248, 272, 5, 71, 0, 0, 249, 272, 5, 72, 0, 0, 250, 251, 5, 67,
		0, 0, 251, 252, 5, 11, 0, 0, 252, 272, 5, 68, 0, 0, 253, 254, 5, 69, 0,
		0, 254, 255, 5, 11, 0, 0, 255, 272, 5, 70, 0, 0, 256, 257, 5, 69, 0, 0,
		257, 258, 5, 11, 0, 0, 258, 272, 5, 71, 0, 0, 259, 260, 5, 69, 0, 0, 260,
		261, 5, 11, 0, 0, 261, 272, 5, 72, 0, 0, 262, 263, 5, 70, 0, 0, 263, 264,
		5, 11, 0, 0, 264, 272, 5, 71, 0, 0, 265, 266, 5, 70, 0, 0, 266, 267, 5,
		11, 0, 0, 267, 272, 5, 72, 0, 0, 268, 269, 5, 71, 0, 0, 269, 270, 5, 11,
		0, 0, 270, 272, 5, 72, 0, 0, 271, 244, 1, 0, 0, 0, 271, 245, 1, 0, 0, 0,
		271, 246, 1, 0, 0, 0, 271, 247, 1, 0, 0, 0, 271, 248, 1, 0, 0, 0, 271,
		249, 1, 0, 0, 0, 271, 250, 1, 0, 0, 0, 271, 253, 1, 0, 0, 0, 271, 256,
		1, 0, 0, 0, 271, 259, 1, 0, 0, 0, 271, 262, 1, 0, 0, 0, 271, 265, 1, 0,
		0, 0, 271, 268, 1, 0, 0, 0, 272, 13, 1, 0, 0, 0, 273, 274, 5, 3, 0, 0,
		274, 275, 5, 4, 0, 0, 275, 276, 3, 104, 52, 0, 276, 277, 5, 79, 0, 0, 277,
		278, 3, 16, 8, 0, 278, 279, 5, 80, 0, 0, 279, 280, 5, 5, 0, 0, 280, 285,
		3, 18, 9, 0, 281, 282, 5, 82, 0, 0, 282, 284, 3, 18, 9, 0, 283, 281, 1,
		0, 0, 0, 284, 287, 1, 0, 0, 0, 285, 283, 1, 0, 0, 0, 285, 286, 1, 0, 0,
		0, 286, 15, 1, 0, 0, 0, 287, 285, 1, 0, 0, 0, 288, 293, 5, 102, 0, 0, 289,
		290, 5, 82, 0, 0, 290, 292, 5, 102, 0, 0, 291, 289, 1, 0, 0, 0, 292, 295,
		1, 0, 0, 0, 293, 291, 1, 0, 0, 0, 293, 294, 1, 0, 0, 0, 294, 17, 1, 0,
		0, 0, 295, 293, 1, 0, 0, 0, 296, 297, 5, 79, 0, 0, 297, 302, 3, 62, 31,
		0, 298, 299, 5, 82, 0, 0, 299, 301, 3, 62, 31, 0, 300, 298, 1, 0, 0, 0,
		301, 304, 1, 0, 0, 0, 302, 300, 1, 0, 0, 0, 302, 303, 1, 0, 0, 0, 303,
		305, 1, 0, 0, 0, 304, 302, 1, 0, 0, 0, 305, 306, 5, 80, 0, 0, 306, 19,
		1, 0, 0, 0, 307, 308, 5, 8, 0, 0, 308, 311, 5, 2, 0, 0, 309, 310, 5, 12,
		0, 0, 310, 312, 5, 14, 0, 0, 311, 309, 1, 0, 0, 0, 311, 312, 1, 0, 0, 0,
		312, 313, 1, 0, 0, 0, 313, 314, 3, 104, 52, 0, 314, 21, 1, 0, 0, 0, 315,
		317, 5, 9, 0, 0, 316, 318, 5, 2, 0, 0, 317, 316, 1, 0, 0, 0, 317, 318,
		1, 0, 0, 0, 318, 319, 1, 0, 0, 0, 319, 320, 3, 104, 52, 0, 320, 23, 1,
		0, 0, 0, 321, 322, 5, 10, 0, 0, 322, 323, 5, 102, 0, 0, 323, 324, 7, 0,
		0, 0, 324, 325, 3, 110, 55, 0, 325, 25, 1, 0, 0, 0, 326, 327, 5, 18, 0,
		0, 327, 331, 5, 102, 0, 0, 328, 330, 5, 102, 0, 0, 329, 328, 1, 0, 0, 0,
		330, 333, 1, 0, 0, 0, 331, 329, 1, 0, 0, 0, 331, 332, 1, 0, 0, 0, 332,
		27, 1, 0, 0, 0, 333, 331, 1, 0, 0, 0, 334, 335, 3, 30, 15, 0, 335, 29,
		1, 0, 0, 0, 336, 344, 3, 32, 16, 0, 337, 339, 7, 1, 0, 0, 338, 340, 5,
		32, 0, 0, 339, 338, 1, 0, 0, 0, 339, 340, 1, 0, 0, 0, 340, 341, 1, 0, 0,
		0, 341, 343, 3, 32, 16, 0, 342, 337, 1, 0, 0, 0, 343, 346, 1, 0, 0, 0,
		344, 342, 1, 0, 0, 0, 344, 345, 1, 0, 0, 0, 345, 31, 1, 0, 0, 0, 346, 344,
		1, 0, 0, 0, 347, 355, 3, 34, 17, 0, 348, 350, 5, 30, 0, 0, 349, 351, 5,
		32, 0, 0, 350, 349, 1, 0, 0, 0, 350, 351, 1, 0, 0, 0, 351, 352, 1, 0, 0,
		0, 352, 354, 3, 34, 17, 0, 353, 348, 1, 0, 0, 0, 354, 357, 1, 0, 0, 0,
		355, 353, 1, 0, 0, 0, 355, 356, 1, 0, 0, 0, 356, 33, 1, 0, 0, 0, 357, 355,
		1, 0, 0, 0, 358, 359, 5, 19, 0, 0, 359, 362, 3, 36, 18, 0, 360, 361, 5,
		20, 0, 0, 361, 363, 3, 42, 21, 0, 362, 360, 1, 0, 0, 0, 362, 363, 1, 0,
		0, 0, 363, 366, 1, 0, 0, 0, 364, 365, 5, 21, 0, 0, 365, 367, 3, 48, 24,
		0, 366, 364, 1, 0, 0, 0, 366, 367, 1, 0, 0, 0, 367, 371, 1, 0, 0, 0, 368,
		369, 5, 22, 0, 0, 369, 370, 5, 26, 0, 0, 370, 372, 3, 58, 29, 0, 371, 368,
		1, 0, 0, 0, 371, 372, 1, 0, 0, 0, 372, 376, 1, 0, 0, 0, 373, 374, 5, 25,
		0, 0, 374, 375, 5, 26, 0, 0, 375, 377, 3, 50, 25, 0, 376, 373, 1, 0, 0,
		0, 376, 377, 1, 0, 0, 0, 377, 380, 1, 0, 0, 0, 378, 379, 5, 23, 0, 0, 379,
		381, 3, 54, 27, 0, 380, 378, 1, 0, 0, 0, 380, 381, 1, 0, 0, 0, 381, 384,
		1, 0, 0, 0, 382, 383, 5, 24, 0, 0, 383, 385, 3, 56, 28, 0, 384, 382, 1,
		0, 0, 0, 384, 385, 1, 0, 0, 0, 385, 391, 1, 0, 0, 0, 386, 387, 5, 79, 0,
		0, 387, 388, 3, 28, 14, 0, 388, 389, 5, 80, 0, 0, 389, 391, 1, 0, 0, 0,
		390, 358, 1, 0, 0, 0, 390, 386, 1, 0, 0, 0, 391, 35, 1, 0, 0, 0, 392, 397,
		3, 38, 19, 0, 393, 394, 5, 82, 0, 0, 394, 396, 3, 38, 19, 0, 395, 393,
		1, 0, 0, 0, 396, 399, 1, 0, 0, 0, 397, 395, 1, 0, 0, 0, 397, 398, 1, 0,
		0, 0, 398, 37, 1, 0, 0, 0, 399, 397, 1, 0, 0, 0, 400, 402, 3, 62, 31, 0,
		401, 403, 3, 60, 30, 0, 402, 401, 1, 0, 0, 0, 402, 403, 1, 0, 0, 0, 403,
		406, 1, 0, 0, 0, 404, 406, 3, 40, 20, 0, 405, 400, 1, 0, 0, 0, 405, 404,
		1, 0, 0, 0, 406, 39, 1, 0, 0, 0, 407, 408, 5, 84, 0, 0, 408, 41, 1, 0,
		0, 0, 409, 414, 3, 44, 22, 0, 410, 411, 5, 82, 0, 0, 411, 413, 3, 44, 22,
		0, 412, 410, 1, 0, 0, 0, 413, 416, 1, 0, 0, 0, 414, 412, 1, 0, 0, 0, 414,
		415, 1, 0, 0, 0, 415, 420, 1, 0, 0, 0, 416, 414, 1, 0, 0, 0, 417, 419,
		3, 46, 23, 0, 418, 417, 1, 0, 0, 0, 419, 422, 1, 0, 0, 0, 420, 418, 1,
		0, 0, 0, 420, 421, 1, 0, 0, 0, 421, 43, 1, 0, 0, 0, 422, 420, 1, 0, 0,
		0, 423, 425, 3, 104, 52, 0, 424, 426, 3, 60, 30, 0, 425, 424, 1, 0, 0,
		0, 425, 426, 1, 0, 0, 0, 426, 433, 1, 0, 0, 0, 427, 428, 5, 79, 0, 0, 428,
		429, 3, 28, 14, 0, 429, 430, 5, 80, 0, 0, 430, 431, 3, 60, 30, 0, 431,
		433, 1, 0, 0, 0, 432, 423, 1, 0, 0, 0, 432, 427, 1, 0, 0, 0, 433, 45, 1,
		0, 0, 0, 434, 436, 7, 2, 0, 0, 435, 434, 1, 0, 0, 0, 435, 436, 1, 0, 0,
		0, 436, 437, 1, 0, 0, 0, 437, 438, 5, 33, 0, 0, 438, 440, 3, 100, 50, 0,
		439, 441, 3, 60, 30, 0, 440, 439, 1, 0, 0, 0, 440, 441, 1, 0, 0, 0, 441,
		442, 1, 0, 0, 0, 442, 443, 5, 39, 0, 0, 443, 444, 3, 62, 31, 0, 444, 47,
		1, 0, 0, 0, 445, 446, 3, 62, 31, 0, 446, 49, 1, 0, 0, 0, 447, 452, 3, 52,
		26, 0, 448, 449, 5, 82, 0, 0, 449, 451, 3, 52, 26, 0, 450, 448, 1, 0, 0,
		0, 451, 454, 1, 0, 0, 0, 452, 450, 1, 0, 0, 0, 452, 453, 1, 0, 0, 0, 453,
		51, 1, 0, 0, 0, 454, 452, 1, 0, 0, 0, 455, 457, 3, 62, 31, 0, 456, 458,
		7, 3, 0, 0, 457, 456, 1, 0, 0, 0, 457, 458, 1, 0, 0, 0, 458, 53, 1, 0,
		0, 0, 459, 460, 5, 105, 0, 0, 460, 55, 1, 0, 0, 0, 461, 462, 5, 105, 0,
		0, 462, 57, 1, 0, 0, 0, 463, 468, 3, 62, 31, 0, 464, 465, 5, 82, 0, 0,
		465, 467, 3, 62, 31, 0, 466, 464, 1, 0, 0, 0, 467, 470, 1, 0, 0, 0, 468,
		466, 1, 0, 0, 0, 468, 469, 1, 0, 0, 0, 469, 59, 1, 0, 0, 0, 470, 468, 1,
		0, 0, 0, 471, 473, 5, 51, 0, 0, 472, 471, 1, 0, 0, 0, 472, 473, 1, 0, 0,
		0, 473, 474, 1, 0, 0, 0, 474, 475, 5, 102, 0, 0, 475, 61, 1, 0, 0, 0, 476,
		477, 3, 64, 32, 0, 477, 63, 1, 0, 0, 0, 478, 483, 3, 66, 33, 0, 479, 480,
		5, 42, 0, 0, 480, 482, 3, 66, 33, 0, 481, 479, 1, 0, 0, 0, 482, 485, 1,
		0, 0, 0, 483, 481, 1, 0, 0, 0, 483, 484, 1, 0, 0, 0, 484, 65, 1, 0, 0,
		0, 485, 483, 1, 0, 0, 0, 486, 491, 3, 68, 34, 0, 487, 488, 5, 41, 0, 0,
		488, 490, 3, 68, 34, 0, 489, 487, 1, 0, 0, 0, 490, 493, 1, 0, 0, 0, 491,
		489, 1, 0, 0, 0, 491, 492, 1, 0, 0, 0, 492, 67, 1, 0, 0, 0, 493, 491, 1,
		0, 0, 0, 494, 496, 5, 13, 0, 0, 495, 494, 1, 0, 0, 0, 495, 496, 1, 0, 0,
		0, 496, 497, 1, 0, 0, 0, 497, 504, 3, 70, 35, 0, 498, 500, 5, 13, 0, 0,
		499, 498, 1, 0, 0, 0, 499, 500, 1, 0, 0, 0, 500, 501, 1, 0, 0, 0, 501,
		502, 5, 14, 0, 0, 502, 504, 3, 88, 44, 0, 503, 495, 1, 0, 0, 0, 503, 499,
		1, 0, 0, 0, 504, 69, 1, 0, 0, 0, 505, 542, 3, 72, 36, 0, 506, 509, 3, 106,
		53, 0, 507, 509, 3, 108, 54, 0, 508, 506, 1, 0, 0, 0, 508, 507, 1, 0, 0,
		0, 509, 510, 1, 0, 0, 0, 510, 511, 3, 72, 36, 0, 511, 543, 1, 0, 0, 0,
		512, 514, 5, 13, 0, 0, 513, 512, 1, 0, 0, 0, 513, 514, 1, 0, 0, 0, 514,
		515, 1, 0, 0, 0, 515, 516, 5, 40, 0, 0, 516, 517, 5, 79, 0, 0, 517, 522,
		3, 62, 31, 0, 518, 519, 5, 82, 0, 0, 519, 521, 3, 62, 31, 0, 520, 518,
		1, 0, 0, 0, 521, 524, 1, 0, 0, 0, 522, 520, 1, 0, 0, 0, 522, 523, 1, 0,
		0, 0, 523, 525, 1, 0, 0, 0, 524, 522, 1, 0, 0, 0, 525, 526, 5, 80, 0, 0,
		526, 543, 1, 0, 0, 0, 527, 529, 5, 13, 0, 0, 528, 527, 1, 0, 0, 0, 528,
		529, 1, 0, 0, 0, 529, 530, 1, 0, 0, 0, 530, 531, 5, 40, 0, 0, 531, 543,
		3, 88, 44, 0, 532, 534, 5, 13, 0, 0, 533, 532, 1, 0, 0, 0, 533, 534, 1,
		0, 0, 0, 534, 535, 1, 0, 0, 0, 535, 536, 5, 44, 0, 0, 536, 543, 3, 72,
		36, 0, 537, 539, 5, 43, 0, 0, 538, 540, 5, 13, 0, 0, 539, 538, 1, 0, 0,
		0, 539, 540, 1, 0, 0, 0, 540, 541, 1, 0, 0, 0, 541, 543, 5, 15, 0, 0, 542,
		508, 1, 0, 0, 0, 542, 513, 1, 0, 0, 0, 542, 528, 1, 0, 0, 0, 542, 533,
		1, 0, 0, 0, 542, 537, 1, 0, 0, 0, 542, 543, 1, 0, 0, 0, 543, 71, 1, 0,
		0, 0, 544, 549, 3, 74, 37, 0, 545, 546, 5, 88, 0, 0, 546, 548, 3, 74, 37,
		0, 547, 545, 1, 0, 0, 0, 548, 551, 1, 0, 0, 0, 549, 547, 1, 0, 0, 0, 549,
		550, 1, 0, 0, 0, 550, 73, 1, 0, 0, 0, 551, 549, 1, 0, 0, 0, 552, 557, 3,
		76, 38, 0, 553, 554, 7, 4, 0, 0, 554, 556, 3, 76, 38, 0, 555, 553, 1, 0,
		0, 0, 556, 559, 1, 0, 0, 0, 557, 555, 1, 0, 0, 0, 557, 558, 1, 0, 0, 0,
		558, 75, 1, 0, 0, 0, 559, 557, 1, 0, 0, 0, 560, 565, 3, 78, 39, 0, 561,
		562, 7, 5, 0, 0, 562, 564, 3, 78, 39, 0, 563, 561, 1, 0, 0, 0, 564, 567,
		1, 0, 0, 0, 565, 563, 1, 0, 0, 0, 565, 566, 1, 0, 0, 0, 566, 77, 1, 0,
		0, 0, 567, 565, 1, 0, 0, 0, 568, 570, 7, 4, 0, 0, 569, 568, 1, 0, 0, 0,
		569, 570, 1, 0, 0, 0, 570, 571, 1, 0, 0, 0, 571, 572, 3, 80, 40, 0, 572,
		79, 1, 0, 0, 0, 573, 577, 3, 86, 43, 0, 574, 576, 3, 82, 41, 0, 575, 574,
		1, 0, 0, 0, 576, 579, 1, 0, 0, 0, 577, 575, 1, 0, 0, 0, 577, 578, 1, 0,
		0, 0, 578, 81, 1, 0, 0, 0, 579, 577, 1, 0, 0, 0, 580, 581, 5, 52, 0, 0,
		581, 582, 5, 53, 0, 0, 582, 583, 5, 54, 0, 0, 583, 593, 5, 104, 0, 0, 584,
		585, 5, 101, 0, 0, 585, 593, 3, 100, 50, 0, 586, 587, 5, 89, 0, 0, 587,
		589, 3, 84, 42, 0, 588, 586, 1, 0, 0, 0, 589, 590, 1, 0, 0, 0, 590, 588,
		1, 0, 0, 0, 590, 591, 1, 0, 0, 0, 591, 593, 1, 0, 0, 0, 592, 580, 1, 0,
		0, 0, 592, 584, 1, 0, 0, 0, 592, 588, 1, 0, 0, 0, 593, 83, 1, 0, 0, 0,
		594, 595, 3, 100, 50, 0, 595, 85, 1, 0, 0, 0, 596, 597, 5, 79, 0, 0, 597,
		598, 3, 62, 31, 0, 598, 599, 5, 80, 0, 0, 599, 608, 1, 0, 0, 0, 600, 608,
		3, 88, 44, 0, 601, 608, 3, 90, 45, 0, 602, 608, 3, 92, 46, 0, 603, 608,
		3, 96, 48, 0, 604, 608, 3, 102, 51, 0, 605, 608, 3, 110, 55, 0, 606, 608,
		5, 103, 0, 0, 607, 596, 1, 0, 0, 0, 607, 600, 1, 0, 0, 0, 607, 601, 1,
		0, 0, 0, 607, 602, 1, 0, 0, 0, 607, 603, 1, 0, 0, 0, 607, 604, 1, 0, 0,
		0, 607, 605, 1, 0, 0, 0, 607, 606, 1, 0, 0, 0, 608, 87, 1, 0, 0, 0, 609,
		610, 5, 79, 0, 0, 610, 611, 3, 28, 14, 0, 611, 612, 5, 80, 0, 0, 612, 89,
		1, 0, 0, 0, 613, 619, 5, 46, 0, 0, 614, 615, 5, 47, 0, 0, 615, 616, 3,
		62, 31, 0, 616, 617, 5, 48, 0, 0, 617, 618, 3, 62, 31, 0, 618, 620, 1,
		0, 0, 0, 619, 614, 1, 0, 0, 0, 620, 621, 1, 0, 0, 0, 621, 619, 1, 0, 0,
		0, 621, 622, 1, 0, 0, 0, 622, 625, 1, 0, 0, 0, 623, 624, 5, 49, 0, 0, 624,
		626, 3, 62, 31, 0, 625, 623, 1, 0, 0, 0, 625, 626, 1, 0, 0, 0, 626, 627,
		1, 0, 0, 0, 627, 628, 5, 50, 0, 0, 628, 91, 1, 0, 0, 0, 629, 630, 3, 100,
		50, 0, 630, 639, 5, 79, 0, 0, 631, 636, 3, 94, 47, 0, 632, 633, 5, 82,
		0, 0, 633, 635, 3, 94, 47, 0, 634, 632, 1, 0, 0, 0, 635, 638, 1, 0, 0,
		0, 636, 634, 1, 0, 0, 0, 636, 637, 1, 0, 0, 0, 637, 640, 1, 0, 0, 0, 638,
		636, 1, 0, 0, 0, 639, 631, 1, 0, 0, 0, 639, 640, 1, 0, 0, 0, 640, 641,
		1, 0, 0, 0, 641, 642, 5, 80, 0, 0, 642, 93, 1, 0, 0, 0, 643, 646, 3, 62,
		31, 0, 644, 646, 5, 84, 0, 0, 645, 643, 1, 0, 0, 0, 645, 644, 1, 0, 0,
		0, 646, 95, 1, 0, 0, 0, 647, 648, 5, 55, 0, 0, 648, 649, 5, 79, 0, 0, 649,
		650, 5, 102, 0, 0, 650, 651, 5, 20, 0, 0, 651, 652, 3, 62, 31, 0, 652,
		653, 5, 80, 0, 0, 653, 97, 1, 0, 0, 0, 654, 655, 7, 6, 0, 0, 655, 99, 1,
		0, 0, 0, 656, 664, 3, 98, 49, 0, 657, 660, 5, 83, 0, 0, 658, 661, 3, 98,
		49, 0, 659, 661, 5, 90, 0, 0, 660, 658, 1, 0, 0, 0, 660, 659, 1, 0, 0,
		0, 661, 663, 1, 0, 0, 0, 662, 657, 1, 0, 0, 0, 663, 666, 1, 0, 0, 0, 664,
		662, 1, 0, 0, 0, 664, 665, 1, 0, 0, 0, 665, 101, 1, 0, 0, 0, 666, 664,
		1, 0, 0, 0, 667, 668, 3, 100, 50, 0, 668, 103, 1, 0, 0, 0, 669, 670, 3,
		100, 50, 0, 670, 105, 1, 0, 0, 0, 671, 672, 7, 7, 0, 0, 672, 107, 1, 0,
		0, 0, 673, 674, 5, 100, 0, 0, 674, 675, 5, 79, 0, 0, 675, 676, 3, 100,
		50, 0, 676, 677, 5, 80, 0, 0, 677, 109, 1, 0, 0, 0, 678, 687, 5, 104, 0,
		0, 679, 687, 5, 105, 0, 0, 680, 687, 5, 75, 0, 0, 681, 687, 5, 76, 0, 0,
		682, 687, 5, 77, 0, 0, 683, 687, 5, 78, 0, 0, 684, 687, 5, 15, 0, 0, 685,
		687, 3, 112, 56, 0, 686, 678, 1, 0, 0, 0, 686, 679, 1, 0, 0, 0, 686, 680,
		1, 0, 0, 0, 686, 681, 1, 0, 0, 0, 686, 682, 1, 0, 0, 0, 686, 683, 1, 0,
		0, 0, 686, 684, 1, 0, 0, 0, 686, 685, 1, 0, 0, 0, 687, 111, 1, 0, 0, 0,
		688, 689, 5, 65, 0, 0, 689, 700, 5, 104, 0, 0, 690, 691, 5, 64, 0, 0, 691,
		700, 5, 104, 0, 0, 692, 693, 5, 53, 0, 0, 693, 700, 5, 104, 0, 0, 694,
		695, 5, 66, 0, 0, 695, 697, 5, 104, 0, 0, 696, 698, 3, 12, 6, 0, 697, 696,
		1, 0, 0, 0, 697, 698, 1, 0, 0, 0, 698, 700, 1, 0, 0, 0, 699, 688, 1, 0,
		0, 0, 699, 690, 1, 0, 0, 0, 699, 692, 1, 0, 0, 0, 699, 694, 1, 0, 0, 0,
		700, 113, 1, 0, 0, 0, 80, 116, 125, 132, 141, 145, 155, 160, 174, 181,
		183, 195, 201, 208, 216, 223, 231, 235, 240, 242, 271, 285, 293, 302, 311,
		317, 331, 339, 344, 350, 355, 362, 366, 371, 376, 380, 384, 390, 397, 402,
		405, 414, 420, 425, 432, 435, 440, 452, 457, 468, 472, 483, 491, 495, 499,
		503, 508, 513, 522, 528, 533, 539, 542, 549, 557, 565, 569, 577, 590, 592,
		607, 621, 625, 636, 639, 645, 660, 664, 686, 697, 699,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// sqlParserInit initializes any static state used to implement sqlParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewsqlParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func SqlParserInit() {
	staticData := &SqlParserStaticData
	staticData.once.Do(sqlParserInit)
}

// NewsqlParser produces a new parser instance for the optional input antlr.TokenStream.
func NewsqlParser(input antlr.TokenStream) *sqlParser {
	SqlParserInit()
	this := new(sqlParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &SqlParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.PredictionContextCache)
	this.RuleNames = staticData.RuleNames
	this.LiteralNames = staticData.LiteralNames
	this.SymbolicNames = staticData.SymbolicNames
	this.GrammarFileName = "sql.g4"

	return this
}

// sqlParser tokens.
const (
	sqlParserEOF               = antlr.TokenEOF
	sqlParserCREATE            = 1
	sqlParserTABLE             = 2
	sqlParserINSERT            = 3
	sqlParserINTO              = 4
	sqlParserVALUES            = 5
	sqlParserPRIMARY           = 6
	sqlParserKEY               = 7
	sqlParserDROP              = 8
	sqlParserTRUNCATE          = 9
	sqlParserSET               = 10
	sqlParserTO                = 11
	sqlParserIF                = 12
	sqlParserNOT               = 13
	sqlParserEXISTS            = 14
	sqlParserNULL              = 15
	sqlParserUNIQUE            = 16
	sqlParserDEFAULT           = 17
	sqlParserSHOW              = 18
	sqlParserSELECT            = 19
	sqlParserFROM              = 20
	sqlParserWHERE             = 21
	sqlParserGROUP             = 22
	sqlParserLIMIT             = 23
	sqlParserOFFSET            = 24
	sqlParserORDER             = 25
	sqlParserBY                = 26
	sqlParserASC               = 27
	sqlParserDESC              = 28
	sqlParserUNION             = 29
	sqlParserINTERSECT         = 30
	sqlParserEXCEPT            = 31
	sqlParserALL               = 32
	sqlParserJOIN              = 33
	sqlParserINNER             = 34
	sqlParserLEFT              = 35
	sqlParserRIGHT             = 36
	sqlParserFULL              = 37
	sqlParserCROSS             = 38
	sqlParserON                = 39
	sqlParserIN                = 40
	sqlParserAND               = 41
	sqlParserOR                = 42
	sqlParserIS                = 43
	sqlParserLIKE              = 44
	sqlParserILIKE             = 45
	sqlParserCASE              = 46
	sqlParserWHEN              = 47
	sqlParserTHEN              = 48
	sqlParserELSE              = 49
	sqlParserEND               = 50
	sqlParserAS                = 51
	sqlParserAT                = 52
	sqlParserTIME              = 53
	sqlParserZONE              = 54
	sqlParserEXTRACT           = 55
	sqlParserSTRING_TYPE       = 56
	sqlParserINT_TYPE          = 57
	sqlParserFLOAT_TYPE        = 58
	sqlParserBOOL_TYPE         = 59
	sqlParserTEXT_TYPE         = 60
	sqlParserVARCHAR_TYPE      = 61
	sqlParserCHAR_TYPE         = 62
	sqlParserSERIAL_TYPE       = 63
	sqlParserTIMESTAMP_TYPE    = 64
	sqlParserDATE_TYPE         = 65
	sqlParserINTERVAL_TYPE     = 66
	sqlParserYEAR              = 67
	sqlParserMONTH             = 68
	sqlParserDAY               = 69
	sqlParserHOUR              = 70
	sqlParserMINUTE            = 71
	sqlParserSECOND            = 72
	sqlParserWITH              = 73
	sqlParserWITHOUT           = 74
	sqlParserCURRENT_TIMESTAMP = 75
	sqlParserCURRENT_USER      = 76
	sqlParserTRUE              = 77
	sqlParserFALSE             = 78
	sqlParserLPAREN            = 79
	sqlParserRPAREN            = 80
	sqlParserSEMICOLON         = 81
	sqlParserCOMMA             = 82
	sqlParserDOT               = 83
	sqlParserSTAR              = 84
	sqlParserPLUS              = 85
	sqlParserMINUS             = 86
	sqlParserSLASH             = 87
	sqlParserCONCAT            = 88
	sqlParserCOLON_COLON       = 89
	sqlParserTILDE             = 90
	sqlParserNREGEX            = 91
	sqlParserIREGEX            = 92
	sqlParserNIREGEX           = 93
	sqlParserEQ                = 94
	sqlParserGT                = 95
	sqlParserLT                = 96
	sqlParserGE                = 97
	sqlParserLE                = 98
	sqlParserNE                = 99
	sqlParserOPERATOR_KW       = 100
	sqlParserCOLLATE           = 101
	sqlParserIDENTIFIER        = 102
	sqlParserPARAMETER         = 103
	sqlParserSTRING_LITERAL    = 104
	sqlParserNUMBER            = 105
	sqlParserBLOCK_COMMENT     = 106
	sqlParserLINE_COMMENT      = 107
	sqlParserWS                = 108
)

// sqlParser rules.
const (
	sqlParserRULE_query                    = 0
	sqlParserRULE_statement                = 1
	sqlParserRULE_createTableStatement     = 2
	sqlParserRULE_columnDefinition         = 3
	sqlParserRULE_columnConstraints        = 4
	sqlParserRULE_dataType                 = 5
	sqlParserRULE_intervalFields           = 6
	sqlParserRULE_insertStatement          = 7
	sqlParserRULE_columnList               = 8
	sqlParserRULE_valueList                = 9
	sqlParserRULE_dropTableStatement       = 10
	sqlParserRULE_truncateTableStatement   = 11
	sqlParserRULE_setStatement             = 12
	sqlParserRULE_showStatement            = 13
	sqlParserRULE_selectStatement          = 14
	sqlParserRULE_unionExceptStatement     = 15
	sqlParserRULE_intersectStatement       = 16
	sqlParserRULE_primarySelectStatement   = 17
	sqlParserRULE_selectList               = 18
	sqlParserRULE_selectItem               = 19
	sqlParserRULE_selectAll                = 20
	sqlParserRULE_fromClause               = 21
	sqlParserRULE_tableFactor              = 22
	sqlParserRULE_joinClause               = 23
	sqlParserRULE_whereClause              = 24
	sqlParserRULE_orderByClause            = 25
	sqlParserRULE_orderByItem              = 26
	sqlParserRULE_limitValue               = 27
	sqlParserRULE_offsetValue              = 28
	sqlParserRULE_groupByClause            = 29
	sqlParserRULE_alias                    = 30
	sqlParserRULE_expression               = 31
	sqlParserRULE_orExpression             = 32
	sqlParserRULE_andExpression            = 33
	sqlParserRULE_notExpression            = 34
	sqlParserRULE_comparisonExpression     = 35
	sqlParserRULE_concatExpression         = 36
	sqlParserRULE_additiveExpression       = 37
	sqlParserRULE_multiplicativeExpression = 38
	sqlParserRULE_unaryExpression          = 39
	sqlParserRULE_castExpression           = 40
	sqlParserRULE_postfix                  = 41
	sqlParserRULE_typeName                 = 42
	sqlParserRULE_primaryExpression        = 43
	sqlParserRULE_subqueryExpression       = 44
	sqlParserRULE_caseExpression           = 45
	sqlParserRULE_functionCall             = 46
	sqlParserRULE_functionArg              = 47
	sqlParserRULE_extractFunction          = 48
	sqlParserRULE_namePart                 = 49
	sqlParserRULE_qualifiedName            = 50
	sqlParserRULE_columnName               = 51
	sqlParserRULE_tableName                = 52
	sqlParserRULE_operator                 = 53
	sqlParserRULE_operatorExpr             = 54
	sqlParserRULE_value                    = 55
	sqlParserRULE_typedLiteral             = 56
)

// IQueryContext is an interface to support dynamic dispatch.
type IQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Statement() IStatementContext
	SEMICOLON() antlr.TerminalNode

	// IsQueryContext differentiates from other interfaces.
	IsQueryContext()
}

type QueryContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQueryContext() *QueryContext {
	var p = new(QueryContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_query
	return p
}

func InitEmptyQueryContext(p *QueryContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_query
}

func (*QueryContext) IsQueryContext() {}

func NewQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QueryContext {
	var p = new(QueryContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_query

	return p
}

func (s *QueryContext) GetParser() antlr.Parser { return s.parser }

func (s *QueryContext) Statement() IStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IStatementContext)
}

func (s *QueryContext) SEMICOLON() antlr.TerminalNode {
	return s.GetToken(sqlParserSEMICOLON, 0)
}

func (s *QueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterQuery(s)
	}
}

func (s *QueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitQuery(s)
	}
}

func (p *sqlParser) Query() (localctx IQueryContext) {
	localctx = NewQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, sqlParserRULE_query)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(114)
		p.Statement()
	}
	p.SetState(116)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserSEMICOLON {
		{
			p.SetState(115)
			p.Match(sqlParserSEMICOLON)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IStatementContext is an interface to support dynamic dispatch.
type IStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SelectStatement() ISelectStatementContext
	CreateTableStatement() ICreateTableStatementContext
	InsertStatement() IInsertStatementContext
	DropTableStatement() IDropTableStatementContext
	TruncateTableStatement() ITruncateTableStatementContext
	SetStatement() ISetStatementContext
	ShowStatement() IShowStatementContext

	// IsStatementContext differentiates from other interfaces.
	IsStatementContext()
}

type StatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyStatementContext() *StatementContext {
	var p = new(StatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_statement
	return p
}

func InitEmptyStatementContext(p *StatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_statement
}

func (*StatementContext) IsStatementContext() {}

func NewStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StatementContext {
	var p = new(StatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_statement

	return p
}

func (s *StatementContext) GetParser() antlr.Parser { return s.parser }

func (s *StatementContext) SelectStatement() ISelectStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectStatementContext)
}

func (s *StatementContext) CreateTableStatement() ICreateTableStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICreateTableStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICreateTableStatementContext)
}

func (s *StatementContext) InsertStatement() IInsertStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInsertStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInsertStatementContext)
}

func (s *StatementContext) DropTableStatement() IDropTableStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDropTableStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDropTableStatementContext)
}

func (s *StatementContext) TruncateTableStatement() ITruncateTableStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITruncateTableStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITruncateTableStatementContext)
}

func (s *StatementContext) SetStatement() ISetStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISetStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISetStatementContext)
}

func (s *StatementContext) ShowStatement() IShowStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IShowStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IShowStatementContext)
}

func (s *StatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *StatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *StatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterStatement(s)
	}
}

func (s *StatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitStatement(s)
	}
}

func (p *sqlParser) Statement() (localctx IStatementContext) {
	localctx = NewStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, sqlParserRULE_statement)
	p.SetState(125)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case sqlParserSELECT, sqlParserLPAREN:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(118)
			p.SelectStatement()
		}

	case sqlParserCREATE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(119)
			p.CreateTableStatement()
		}

	case sqlParserINSERT:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(120)
			p.InsertStatement()
		}

	case sqlParserDROP:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(121)
			p.DropTableStatement()
		}

	case sqlParserTRUNCATE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(122)
			p.TruncateTableStatement()
		}

	case sqlParserSET:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(123)
			p.SetStatement()
		}

	case sqlParserSHOW:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(124)
			p.ShowStatement()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICreateTableStatementContext is an interface to support dynamic dispatch.
type ICreateTableStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CREATE() antlr.TerminalNode
	TABLE() antlr.TerminalNode
	TableName() ITableNameContext
	AllLPAREN() []antlr.TerminalNode
	LPAREN(i int) antlr.TerminalNode
	AllColumnDefinition() []IColumnDefinitionContext
	ColumnDefinition(i int) IColumnDefinitionContext
	AllRPAREN() []antlr.TerminalNode
	RPAREN(i int) antlr.TerminalNode
	IF() antlr.TerminalNode
	NOT() antlr.TerminalNode
	EXISTS() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	PRIMARY() antlr.TerminalNode
	KEY() antlr.TerminalNode
	AllColumnName() []IColumnNameContext
	ColumnName(i int) IColumnNameContext

	// IsCreateTableStatementContext differentiates from other interfaces.
	IsCreateTableStatementContext()
}

type CreateTableStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCreateTableStatementContext() *CreateTableStatementContext {
	var p = new(CreateTableStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_createTableStatement
	return p
}

func InitEmptyCreateTableStatementContext(p *CreateTableStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_createTableStatement
}

func (*CreateTableStatementContext) IsCreateTableStatementContext() {}

func NewCreateTableStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CreateTableStatementContext {
	var p = new(CreateTableStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_createTableStatement

	return p
}

func (s *CreateTableStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *CreateTableStatementContext) CREATE() antlr.TerminalNode {
	return s.GetToken(sqlParserCREATE, 0)
}

func (s *CreateTableStatementContext) TABLE() antlr.TerminalNode {
	return s.GetToken(sqlParserTABLE, 0)
}

func (s *CreateTableStatementContext) TableName() ITableNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableNameContext)
}

func (s *CreateTableStatementContext) AllLPAREN() []antlr.TerminalNode {
	return s.GetTokens(sqlParserLPAREN)
}

func (s *CreateTableStatementContext) LPAREN(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserLPAREN, i)
}

func (s *CreateTableStatementContext) AllColumnDefinition() []IColumnDefinitionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IColumnDefinitionContext); ok {
			len++
		}
	}

	tst := make([]IColumnDefinitionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IColumnDefinitionContext); ok {
			tst[i] = t.(IColumnDefinitionContext)
			i++
		}
	}

	return tst
}

func (s *CreateTableStatementContext) ColumnDefinition(i int) IColumnDefinitionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnDefinitionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnDefinitionContext)
}

func (s *CreateTableStatementContext) AllRPAREN() []antlr.TerminalNode {
	return s.GetTokens(sqlParserRPAREN)
}

func (s *CreateTableStatementContext) RPAREN(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserRPAREN, i)
}

func (s *CreateTableStatementContext) IF() antlr.TerminalNode {
	return s.GetToken(sqlParserIF, 0)
}

func (s *CreateTableStatementContext) NOT() antlr.TerminalNode {
	return s.GetToken(sqlParserNOT, 0)
}

func (s *CreateTableStatementContext) EXISTS() antlr.TerminalNode {
	return s.GetToken(sqlParserEXISTS, 0)
}

func (s *CreateTableStatementContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(sqlParserCOMMA)
}

func (s *CreateTableStatementContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserCOMMA, i)
}

func (s *CreateTableStatementContext) PRIMARY() antlr.TerminalNode {
	return s.GetToken(sqlParserPRIMARY, 0)
}

func (s *CreateTableStatementContext) KEY() antlr.TerminalNode {
	return s.GetToken(sqlParserKEY, 0)
}

func (s *CreateTableStatementContext) AllColumnName() []IColumnNameContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IColumnNameContext); ok {
			len++
		}
	}

	tst := make([]IColumnNameContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IColumnNameContext); ok {
			tst[i] = t.(IColumnNameContext)
			i++
		}
	}

	return tst
}

func (s *CreateTableStatementContext) ColumnName(i int) IColumnNameContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnNameContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnNameContext)
}

func (s *CreateTableStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CreateTableStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CreateTableStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterCreateTableStatement(s)
	}
}

func (s *CreateTableStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitCreateTableStatement(s)
	}
}

func (p *sqlParser) CreateTableStatement() (localctx ICreateTableStatementContext) {
	localctx = NewCreateTableStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, sqlParserRULE_createTableStatement)
	var _la int

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(127)
		p.Match(sqlParserCREATE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(128)
		p.Match(sqlParserTABLE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(132)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserIF {
		{
			p.SetState(129)
			p.Match(sqlParserIF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(130)
			p.Match(sqlParserNOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(131)
			p.Match(sqlParserEXISTS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(134)
		p.TableName()
	}
	{
		p.SetState(135)
		p.Match(sqlParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(136)
		p.ColumnDefinition()
	}
	p.SetState(141)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext())
	if p.HasError() {
		goto errorExit
	}
	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			{
				p.SetState(137)
				p.Match(sqlParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(138)
				p.ColumnDefinition()
			}

		}
		p.SetState(143)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(145)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserCOMMA {
		{
			p.SetState(144)
			p.Match(sqlParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	p.SetState(160)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserPRIMARY {
		{
			p.SetState(147)
			p.Match(sqlParserPRIMARY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(148)
			p.Match(sqlParserKEY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(149)
			p.Match(sqlParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(150)
			p.ColumnName()
		}
		p.SetState(155)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == sqlParserCOMMA {
			{
				p.SetState(151)
				p.Match(sqlParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(152)
				p.ColumnName()
			}

			p.SetState(157)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(158)
			p.Match(sqlParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(162)
		p.Match(sqlParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IColumnDefinitionContext is an interface to support dynamic dispatch.
type IColumnDefinitionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ColumnName() IColumnNameContext
	DataType() IDataTypeContext
	ColumnConstraints() IColumnConstraintsContext
	PRIMARY() antlr.TerminalNode
	KEY() antlr.TerminalNode

	// IsColumnDefinitionContext differentiates from other interfaces.
	IsColumnDefinitionContext()
}

type ColumnDefinitionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnDefinitionContext() *ColumnDefinitionContext {
	var p = new(ColumnDefinitionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_columnDefinition
	return p
}

func InitEmptyColumnDefinitionContext(p *ColumnDefinitionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_columnDefinition
}

func (*ColumnDefinitionContext) IsColumnDefinitionContext() {}

func NewColumnDefinitionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnDefinitionContext {
	var p = new(ColumnDefinitionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_columnDefinition

	return p
}

func (s *ColumnDefinitionContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnDefinitionContext) ColumnName() IColumnNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnNameContext)
}

func (s *ColumnDefinitionContext) DataType() IDataTypeContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IDataTypeContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IDataTypeContext)
}

func (s *ColumnDefinitionContext) ColumnConstraints() IColumnConstraintsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnConstraintsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnConstraintsContext)
}

func (s *ColumnDefinitionContext) PRIMARY() antlr.TerminalNode {
	return s.GetToken(sqlParserPRIMARY, 0)
}

func (s *ColumnDefinitionContext) KEY() antlr.TerminalNode {
	return s.GetToken(sqlParserKEY, 0)
}

func (s *ColumnDefinitionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnDefinitionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ColumnDefinitionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterColumnDefinition(s)
	}
}

func (s *ColumnDefinitionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitColumnDefinition(s)
	}
}

func (p *sqlParser) ColumnDefinition() (localctx IColumnDefinitionContext) {
	localctx = NewColumnDefinitionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, sqlParserRULE_columnDefinition)
	p.SetState(174)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 7, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(164)
			p.ColumnName()
		}
		{
			p.SetState(165)
			p.DataType()
		}
		{
			p.SetState(166)
			p.ColumnConstraints()
		}
		{
			p.SetState(167)
			p.Match(sqlParserPRIMARY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(168)
			p.Match(sqlParserKEY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(170)
			p.ColumnName()
		}
		{
			p.SetState(171)
			p.DataType()
		}
		{
			p.SetState(172)
			p.ColumnConstraints()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IColumnConstraintsContext is an interface to support dynamic dispatch.
type IColumnConstraintsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNOT() []antlr.TerminalNode
	NOT(i int) antlr.TerminalNode
	AllNULL() []antlr.TerminalNode
	NULL(i int) antlr.TerminalNode
	AllUNIQUE() []antlr.TerminalNode
	UNIQUE(i int) antlr.TerminalNode
	AllDEFAULT() []antlr.TerminalNode
	DEFAULT(i int) antlr.TerminalNode
	AllValue() []IValueContext
	Value(i int) IValueContext

	// IsColumnConstraintsContext differentiates from other interfaces.
	IsColumnConstraintsContext()
}

type ColumnConstraintsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnConstraintsContext() *ColumnConstraintsContext {
	var p = new(ColumnConstraintsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_columnConstraints
	return p
}

func InitEmptyColumnConstraintsContext(p *ColumnConstraintsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_columnConstraints
}

func (*ColumnConstraintsContext) IsColumnConstraintsContext() {}

func NewColumnConstraintsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnConstraintsContext {
	var p = new(ColumnConstraintsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_columnConstraints

	return p
}

func (s *ColumnConstraintsContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnConstraintsContext) AllNOT() []antlr.TerminalNode {
	return s.GetTokens(sqlParserNOT)
}

func (s *ColumnConstraintsContext) NOT(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserNOT, i)
}

func (s *ColumnConstraintsContext) AllNULL() []antlr.TerminalNode {
	return s.GetTokens(sqlParserNULL)
}

func (s *ColumnConstraintsContext) NULL(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserNULL, i)
}

func (s *ColumnConstraintsContext) AllUNIQUE() []antlr.TerminalNode {
	return s.GetTokens(sqlParserUNIQUE)
}

func (s *ColumnConstraintsContext) UNIQUE(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserUNIQUE, i)
}

func (s *ColumnConstraintsContext) AllDEFAULT() []antlr.TerminalNode {
	return s.GetTokens(sqlParserDEFAULT)
}

func (s *ColumnConstraintsContext) DEFAULT(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserDEFAULT, i)
}

func (s *ColumnConstraintsContext) AllValue() []IValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IValueContext); ok {
			len++
		}
	}

	tst := make([]IValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IValueContext); ok {
			tst[i] = t.(IValueContext)
			i++
		}
	}

	return tst
}

func (s *ColumnConstraintsContext) Value(i int) IValueContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *ColumnConstraintsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnConstraintsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ColumnConstraintsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterColumnConstraints(s)
	}
}

func (s *ColumnConstraintsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitColumnConstraints(s)
	}
}

func (p *sqlParser) ColumnConstraints() (localctx IColumnConstraintsContext) {
	localctx = NewColumnConstraintsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, sqlParserRULE_columnConstraints)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(183)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&204800) != 0 {
		p.SetState(181)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case sqlParserNOT:
			{
				p.SetState(176)
				p.Match(sqlParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(177)
				p.Match(sqlParserNULL)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case sqlParserUNIQUE:
			{
				p.SetState(178)
				p.Match(sqlParserUNIQUE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case sqlParserDEFAULT:
			{
				p.SetState(179)
				p.Match(sqlParserDEFAULT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(180)
				p.Value()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(185)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDataTypeContext is an interface to support dynamic dispatch.
type IDataTypeContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRING_TYPE() antlr.TerminalNode
	INT_TYPE() antlr.TerminalNode
	FLOAT_TYPE() antlr.TerminalNode
	BOOL_TYPE() antlr.TerminalNode
	TEXT_TYPE() antlr.TerminalNode
	VARCHAR_TYPE() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	NUMBER() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	CHAR_TYPE() antlr.TerminalNode
	SERIAL_TYPE() antlr.TerminalNode
	TIMESTAMP_TYPE() antlr.TerminalNode
	WITH() antlr.TerminalNode
	AllTIME() []antlr.TerminalNode
	TIME(i int) antlr.TerminalNode
	ZONE() antlr.TerminalNode
	WITHOUT() antlr.TerminalNode
	DATE_TYPE() antlr.TerminalNode
	INTERVAL_TYPE() antlr.TerminalNode
	IntervalFields() IIntervalFieldsContext

	// IsDataTypeContext differentiates from other interfaces.
	IsDataTypeContext()
}

type DataTypeContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDataTypeContext() *DataTypeContext {
	var p = new(DataTypeContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_dataType
	return p
}

func InitEmptyDataTypeContext(p *DataTypeContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_dataType
}

func (*DataTypeContext) IsDataTypeContext() {}

func NewDataTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DataTypeContext {
	var p = new(DataTypeContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_dataType

	return p
}

func (s *DataTypeContext) GetParser() antlr.Parser { return s.parser }

func (s *DataTypeContext) STRING_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserSTRING_TYPE, 0)
}

func (s *DataTypeContext) INT_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserINT_TYPE, 0)
}

func (s *DataTypeContext) FLOAT_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserFLOAT_TYPE, 0)
}

func (s *DataTypeContext) BOOL_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserBOOL_TYPE, 0)
}

func (s *DataTypeContext) TEXT_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserTEXT_TYPE, 0)
}

func (s *DataTypeContext) VARCHAR_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserVARCHAR_TYPE, 0)
}

func (s *DataTypeContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserLPAREN, 0)
}

func (s *DataTypeContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(sqlParserNUMBER, 0)
}

func (s *DataTypeContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserRPAREN, 0)
}

func (s *DataTypeContext) CHAR_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserCHAR_TYPE, 0)
}

func (s *DataTypeContext) SERIAL_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserSERIAL_TYPE, 0)
}

func (s *DataTypeContext) TIMESTAMP_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserTIMESTAMP_TYPE, 0)
}

func (s *DataTypeContext) WITH() antlr.TerminalNode {
	return s.GetToken(sqlParserWITH, 0)
}

func (s *DataTypeContext) AllTIME() []antlr.TerminalNode {
	return s.GetTokens(sqlParserTIME)
}

func (s *DataTypeContext) TIME(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserTIME, i)
}

func (s *DataTypeContext) ZONE() antlr.TerminalNode {
	return s.GetToken(sqlParserZONE, 0)
}

func (s *DataTypeContext) WITHOUT() antlr.TerminalNode {
	return s.GetToken(sqlParserWITHOUT, 0)
}

func (s *DataTypeContext) DATE_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserDATE_TYPE, 0)
}

func (s *DataTypeContext) INTERVAL_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserINTERVAL_TYPE, 0)
}

func (s *DataTypeContext) IntervalFields() IIntervalFieldsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntervalFieldsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntervalFieldsContext)
}

func (s *DataTypeContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DataTypeContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DataTypeContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterDataType(s)
	}
}

func (s *DataTypeContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitDataType(s)
	}
}

func (p *sqlParser) DataType() (localctx IDataTypeContext) {
	localctx = NewDataTypeContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, sqlParserRULE_dataType)
	var _la int

	p.SetState(242)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case sqlParserSTRING_TYPE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(186)
			p.Match(sqlParserSTRING_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserINT_TYPE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(187)
			p.Match(sqlParserINT_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserFLOAT_TYPE:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(188)
			p.Match(sqlParserFLOAT_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserBOOL_TYPE:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(189)
			p.Match(sqlParserBOOL_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserTEXT_TYPE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(190)
			p.Match(sqlParserTEXT_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserVARCHAR_TYPE:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(191)
			p.Match(sqlParserVARCHAR_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(195)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserLPAREN {
			{
				p.SetState(192)
				p.Match(sqlParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(193)
				p.Match(sqlParserNUMBER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(194)
				p.Match(sqlParserRPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case sqlParserCHAR_TYPE:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(197)
			p.Match(sqlParserCHAR_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(201)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserLPAREN {
			{
				p.SetState(198)
				p.Match(sqlParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(199)
				p.Match(sqlParserNUMBER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(200)
				p.Match(sqlParserRPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	case sqlParserSERIAL_TYPE:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(203)
			p.Match(sqlParserSERIAL_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserTIMESTAMP_TYPE:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(204)
			p.Match(sqlParserTIMESTAMP_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(208)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserLPAREN {
			{
				p.SetState(205)
				p.Match(sqlParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(206)
				p.Match(sqlParserNUMBER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(207)
				p.Match(sqlParserRPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		p.SetState(216)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		switch p.GetTokenStream().LA(1) {
		case sqlParserWITH:
			{
				p.SetState(210)
				p.Match(sqlParserWITH)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(211)
				p.Match(sqlParserTIME)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(212)
				p.Match(sqlParserZONE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case sqlParserWITHOUT:
			{
				p.SetState(213)
				p.Match(sqlParserWITHOUT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(214)
				p.Match(sqlParserTIME)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(215)
				p.Match(sqlParserZONE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case sqlParserPRIMARY, sqlParserNOT, sqlParserUNIQUE, sqlParserDEFAULT, sqlParserRPAREN, sqlParserCOMMA:

		default:
		}

	case sqlParserDATE_TYPE:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(218)
			p.Match(sqlParserDATE_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserTIME:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(219)
			p.Match(sqlParserTIME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(223)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserLPAREN {
			{
				p.SetState(220)
				p.Match(sqlParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(221)
				p.Match(sqlParserNUMBER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(222)
				p.Match(sqlParserRPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		p.SetState(231)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		switch p.GetTokenStream().LA(1) {
		case sqlParserWITH:
			{
				p.SetState(225)
				p.Match(sqlParserWITH)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(226)
				p.Match(sqlParserTIME)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(227)
				p.Match(sqlParserZONE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case sqlParserWITHOUT:
			{
				p.SetState(228)
				p.Match(sqlParserWITHOUT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(229)
				p.Match(sqlParserTIME)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(230)
				p.Match(sqlParserZONE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case sqlParserPRIMARY, sqlParserNOT, sqlParserUNIQUE, sqlParserDEFAULT, sqlParserRPAREN, sqlParserCOMMA:

		default:
		}

	case sqlParserINTERVAL_TYPE:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(233)
			p.Match(sqlParserINTERVAL_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(235)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64((_la-67)) & ^0x3f) == 0 && ((int64(1)<<(_la-67))&63) != 0 {
			{
				p.SetState(234)
				p.IntervalFields()
			}

		}
		p.SetState(240)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserLPAREN {
			{
				p.SetState(237)
				p.Match(sqlParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(238)
				p.Match(sqlParserNUMBER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(239)
				p.Match(sqlParserRPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIntervalFieldsContext is an interface to support dynamic dispatch.
type IIntervalFieldsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	YEAR() antlr.TerminalNode
	MONTH() antlr.TerminalNode
	DAY() antlr.TerminalNode
	HOUR() antlr.TerminalNode
	MINUTE() antlr.TerminalNode
	SECOND() antlr.TerminalNode
	TO() antlr.TerminalNode

	// IsIntervalFieldsContext differentiates from other interfaces.
	IsIntervalFieldsContext()
}

type IntervalFieldsContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntervalFieldsContext() *IntervalFieldsContext {
	var p = new(IntervalFieldsContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_intervalFields
	return p
}

func InitEmptyIntervalFieldsContext(p *IntervalFieldsContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_intervalFields
}

func (*IntervalFieldsContext) IsIntervalFieldsContext() {}

func NewIntervalFieldsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntervalFieldsContext {
	var p = new(IntervalFieldsContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_intervalFields

	return p
}

func (s *IntervalFieldsContext) GetParser() antlr.Parser { return s.parser }

func (s *IntervalFieldsContext) YEAR() antlr.TerminalNode {
	return s.GetToken(sqlParserYEAR, 0)
}

func (s *IntervalFieldsContext) MONTH() antlr.TerminalNode {
	return s.GetToken(sqlParserMONTH, 0)
}

func (s *IntervalFieldsContext) DAY() antlr.TerminalNode {
	return s.GetToken(sqlParserDAY, 0)
}

func (s *IntervalFieldsContext) HOUR() antlr.TerminalNode {
	return s.GetToken(sqlParserHOUR, 0)
}

func (s *IntervalFieldsContext) MINUTE() antlr.TerminalNode {
	return s.GetToken(sqlParserMINUTE, 0)
}

func (s *IntervalFieldsContext) SECOND() antlr.TerminalNode {
	return s.GetToken(sqlParserSECOND, 0)
}

func (s *IntervalFieldsContext) TO() antlr.TerminalNode {
	return s.GetToken(sqlParserTO, 0)
}

func (s *IntervalFieldsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntervalFieldsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntervalFieldsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterIntervalFields(s)
	}
}

func (s *IntervalFieldsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitIntervalFields(s)
	}
}

func (p *sqlParser) IntervalFields() (localctx IIntervalFieldsContext) {
	localctx = NewIntervalFieldsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, sqlParserRULE_intervalFields)
	p.SetState(271)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 19, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(244)
			p.Match(sqlParserYEAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(245)
			p.Match(sqlParserMONTH)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(246)
			p.Match(sqlParserDAY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(247)
			p.Match(sqlParserHOUR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(248)
			p.Match(sqlParserMINUTE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(249)
			p.Match(sqlParserSECOND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(250)
			p.Match(sqlParserYEAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(251)
			p.Match(sqlParserTO)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(252)
			p.Match(sqlParserMONTH)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(253)
			p.Match(sqlParserDAY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(254)
			p.Match(sqlParserTO)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(255)
			p.Match(sqlParserHOUR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 9:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(256)
			p.Match(sqlParserDAY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(257)
			p.Match(sqlParserTO)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(258)
			p.Match(sqlParserMINUTE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 10:
		p.EnterOuterAlt(localctx, 10)
		{
			p.SetState(259)
			p.Match(sqlParserDAY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(260)
			p.Match(sqlParserTO)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(261)
			p.Match(sqlParserSECOND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 11:
		p.EnterOuterAlt(localctx, 11)
		{
			p.SetState(262)
			p.Match(sqlParserHOUR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(263)
			p.Match(sqlParserTO)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(264)
			p.Match(sqlParserMINUTE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 12:
		p.EnterOuterAlt(localctx, 12)
		{
			p.SetState(265)
			p.Match(sqlParserHOUR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(266)
			p.Match(sqlParserTO)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(267)
			p.Match(sqlParserSECOND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 13:
		p.EnterOuterAlt(localctx, 13)
		{
			p.SetState(268)
			p.Match(sqlParserMINUTE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(269)
			p.Match(sqlParserTO)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(270)
			p.Match(sqlParserSECOND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IInsertStatementContext is an interface to support dynamic dispatch.
type IInsertStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	INSERT() antlr.TerminalNode
	INTO() antlr.TerminalNode
	TableName() ITableNameContext
	LPAREN() antlr.TerminalNode
	ColumnList() IColumnListContext
	RPAREN() antlr.TerminalNode
	VALUES() antlr.TerminalNode
	AllValueList() []IValueListContext
	ValueList(i int) IValueListContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsInsertStatementContext differentiates from other interfaces.
	IsInsertStatementContext()
}

type InsertStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInsertStatementContext() *InsertStatementContext {
	var p = new(InsertStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_insertStatement
	return p
}

func InitEmptyInsertStatementContext(p *InsertStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_insertStatement
}

func (*InsertStatementContext) IsInsertStatementContext() {}

func NewInsertStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InsertStatementContext {
	var p = new(InsertStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_insertStatement

	return p
}

func (s *InsertStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *InsertStatementContext) INSERT() antlr.TerminalNode {
	return s.GetToken(sqlParserINSERT, 0)
}

func (s *InsertStatementContext) INTO() antlr.TerminalNode {
	return s.GetToken(sqlParserINTO, 0)
}

func (s *InsertStatementContext) TableName() ITableNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableNameContext)
}

func (s *InsertStatementContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserLPAREN, 0)
}

func (s *InsertStatementContext) ColumnList() IColumnListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnListContext)
}

func (s *InsertStatementContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserRPAREN, 0)
}

func (s *InsertStatementContext) VALUES() antlr.TerminalNode {
	return s.GetToken(sqlParserVALUES, 0)
}

func (s *InsertStatementContext) AllValueList() []IValueListContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IValueListContext); ok {
			len++
		}
	}

	tst := make([]IValueListContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IValueListContext); ok {
			tst[i] = t.(IValueListContext)
			i++
		}
	}

	return tst
}

func (s *InsertStatementContext) ValueList(i int) IValueListContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueListContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueListContext)
}

func (s *InsertStatementContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(sqlParserCOMMA)
}

func (s *InsertStatementContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserCOMMA, i)
}

func (s *InsertStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InsertStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *InsertStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterInsertStatement(s)
	}
}

func (s *InsertStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitInsertStatement(s)
	}
}

func (p *sqlParser) InsertStatement() (localctx IInsertStatementContext) {
	localctx = NewInsertStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, sqlParserRULE_insertStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(273)
		p.Match(sqlParserINSERT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(274)
		p.Match(sqlParserINTO)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(275)
		p.TableName()
	}
	{
		p.SetState(276)
		p.Match(sqlParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(277)
		p.ColumnList()
	}
	{
		p.SetState(278)
		p.Match(sqlParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(279)
		p.Match(sqlParserVALUES)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(280)
		p.ValueList()
	}
	p.SetState(285)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserCOMMA {
		{
			p.SetState(281)
			p.Match(sqlParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(282)
			p.ValueList()
		}

		p.SetState(287)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IColumnListContext is an interface to support dynamic dispatch.
type IColumnListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsColumnListContext differentiates from other interfaces.
	IsColumnListContext()
}

type ColumnListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnListContext() *ColumnListContext {
	var p = new(ColumnListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_columnList
	return p
}

func InitEmptyColumnListContext(p *ColumnListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_columnList
}

func (*ColumnListContext) IsColumnListContext() {}

func NewColumnListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnListContext {
	var p = new(ColumnListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_columnList

	return p
}

func (s *ColumnListContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnListContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(sqlParserIDENTIFIER)
}

func (s *ColumnListContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserIDENTIFIER, i)
}

func (s *ColumnListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(sqlParserCOMMA)
}

func (s *ColumnListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserCOMMA, i)
}

func (s *ColumnListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ColumnListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterColumnList(s)
	}
}

func (s *ColumnListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitColumnList(s)
	}
}

func (p *sqlParser) ColumnList() (localctx IColumnListContext) {
	localctx = NewColumnListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, sqlParserRULE_columnList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(288)
		p.Match(sqlParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(293)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserCOMMA {
		{
			p.SetState(289)
			p.Match(sqlParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(290)
			p.Match(sqlParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(295)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IValueListContext is an interface to support dynamic dispatch.
type IValueListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	RPAREN() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsValueListContext differentiates from other interfaces.
	IsValueListContext()
}

type ValueListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueListContext() *ValueListContext {
	var p = new(ValueListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_valueList
	return p
}

func InitEmptyValueListContext(p *ValueListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_valueList
}

func (*ValueListContext) IsValueListContext() {}

func NewValueListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueListContext {
	var p = new(ValueListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_valueList

	return p
}

func (s *ValueListContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueListContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserLPAREN, 0)
}

func (s *ValueListContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ValueListContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ValueListContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserRPAREN, 0)
}

func (s *ValueListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(sqlParserCOMMA)
}

func (s *ValueListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserCOMMA, i)
}

func (s *ValueListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ValueListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterValueList(s)
	}
}

func (s *ValueListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitValueList(s)
	}
}

func (p *sqlParser) ValueList() (localctx IValueListContext) {
	localctx = NewValueListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, sqlParserRULE_valueList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(296)
		p.Match(sqlParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(297)
		p.Expression()
	}
	p.SetState(302)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserCOMMA {
		{
			p.SetState(298)
			p.Match(sqlParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(299)
			p.Expression()
		}

		p.SetState(304)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(305)
		p.Match(sqlParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IDropTableStatementContext is an interface to support dynamic dispatch.
type IDropTableStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DROP() antlr.TerminalNode
	TABLE() antlr.TerminalNode
	TableName() ITableNameContext
	IF() antlr.TerminalNode
	EXISTS() antlr.TerminalNode

	// IsDropTableStatementContext differentiates from other interfaces.
	IsDropTableStatementContext()
}

type DropTableStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDropTableStatementContext() *DropTableStatementContext {
	var p = new(DropTableStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_dropTableStatement
	return p
}

func InitEmptyDropTableStatementContext(p *DropTableStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_dropTableStatement
}

func (*DropTableStatementContext) IsDropTableStatementContext() {}

func NewDropTableStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DropTableStatementContext {
	var p = new(DropTableStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_dropTableStatement

	return p
}

func (s *DropTableStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *DropTableStatementContext) DROP() antlr.TerminalNode {
	return s.GetToken(sqlParserDROP, 0)
}

func (s *DropTableStatementContext) TABLE() antlr.TerminalNode {
	return s.GetToken(sqlParserTABLE, 0)
}

func (s *DropTableStatementContext) TableName() ITableNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableNameContext)
}

func (s *DropTableStatementContext) IF() antlr.TerminalNode {
	return s.GetToken(sqlParserIF, 0)
}

func (s *DropTableStatementContext) EXISTS() antlr.TerminalNode {
	return s.GetToken(sqlParserEXISTS, 0)
}

func (s *DropTableStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DropTableStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DropTableStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterDropTableStatement(s)
	}
}

func (s *DropTableStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitDropTableStatement(s)
	}
}

func (p *sqlParser) DropTableStatement() (localctx IDropTableStatementContext) {
	localctx = NewDropTableStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, sqlParserRULE_dropTableStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(307)
		p.Match(sqlParserDROP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(308)
		p.Match(sqlParserTABLE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(311)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserIF {
		{
			p.SetState(309)
			p.Match(sqlParserIF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(310)
			p.Match(sqlParserEXISTS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(313)
		p.TableName()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITruncateTableStatementContext is an interface to support dynamic dispatch.
type ITruncateTableStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TRUNCATE() antlr.TerminalNode
	TableName() ITableNameContext
	TABLE() antlr.TerminalNode

	// IsTruncateTableStatementContext differentiates from other interfaces.
	IsTruncateTableStatementContext()
}

type TruncateTableStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTruncateTableStatementContext() *TruncateTableStatementContext {
	var p = new(TruncateTableStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_truncateTableStatement
	return p
}

func InitEmptyTruncateTableStatementContext(p *TruncateTableStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_truncateTableStatement
}

func (*TruncateTableStatementContext) IsTruncateTableStatementContext() {}

func NewTruncateTableStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TruncateTableStatementContext {
	var p = new(TruncateTableStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_truncateTableStatement

	return p
}

func (s *TruncateTableStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *TruncateTableStatementContext) TRUNCATE() antlr.TerminalNode {
	return s.GetToken(sqlParserTRUNCATE, 0)
}

func (s *TruncateTableStatementContext) TableName() ITableNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableNameContext)
}

func (s *TruncateTableStatementContext) TABLE() antlr.TerminalNode {
	return s.GetToken(sqlParserTABLE, 0)
}

func (s *TruncateTableStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TruncateTableStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TruncateTableStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterTruncateTableStatement(s)
	}
}

func (s *TruncateTableStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitTruncateTableStatement(s)
	}
}

func (p *sqlParser) TruncateTableStatement() (localctx ITruncateTableStatementContext) {
	localctx = NewTruncateTableStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, sqlParserRULE_truncateTableStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(315)
		p.Match(sqlParserTRUNCATE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(317)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserTABLE {
		{
			p.SetState(316)
			p.Match(sqlParserTABLE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(319)
		p.TableName()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISetStatementContext is an interface to support dynamic dispatch.
type ISetStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SET() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	Value() IValueContext
	TO() antlr.TerminalNode
	EQ() antlr.TerminalNode

	// IsSetStatementContext differentiates from other interfaces.
	IsSetStatementContext()
}

type SetStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySetStatementContext() *SetStatementContext {
	var p = new(SetStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_setStatement
	return p
}

func InitEmptySetStatementContext(p *SetStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_setStatement
}

func (*SetStatementContext) IsSetStatementContext() {}

func NewSetStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SetStatementContext {
	var p = new(SetStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_setStatement

	return p
}

func (s *SetStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *SetStatementContext) SET() antlr.TerminalNode {
	return s.GetToken(sqlParserSET, 0)
}

func (s *SetStatementContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(sqlParserIDENTIFIER, 0)
}

func (s *SetStatementContext) Value() IValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *SetStatementContext) TO() antlr.TerminalNode {
	return s.GetToken(sqlParserTO, 0)
}

func (s *SetStatementContext) EQ() antlr.TerminalNode {
	return s.GetToken(sqlParserEQ, 0)
}

func (s *SetStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SetStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SetStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterSetStatement(s)
	}
}

func (s *SetStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitSetStatement(s)
	}
}

func (p *sqlParser) SetStatement() (localctx ISetStatementContext) {
	localctx = NewSetStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, sqlParserRULE_setStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(321)
		p.Match(sqlParserSET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(322)
		p.Match(sqlParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(323)
		_la = p.GetTokenStream().LA(1)

		if !(_la == sqlParserTO || _la == sqlParserEQ) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(324)
		p.Value()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IShowStatementContext is an interface to support dynamic dispatch.
type IShowStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SHOW() antlr.TerminalNode
	AllIDENTIFIER() []antlr.TerminalNode
	IDENTIFIER(i int) antlr.TerminalNode

	// IsShowStatementContext differentiates from other interfaces.
	IsShowStatementContext()
}

type ShowStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyShowStatementContext() *ShowStatementContext {
	var p = new(ShowStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_showStatement
	return p
}

func InitEmptyShowStatementContext(p *ShowStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_showStatement
}

func (*ShowStatementContext) IsShowStatementContext() {}

func NewShowStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ShowStatementContext {
	var p = new(ShowStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_showStatement

	return p
}

func (s *ShowStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *ShowStatementContext) SHOW() antlr.TerminalNode {
	return s.GetToken(sqlParserSHOW, 0)
}

func (s *ShowStatementContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(sqlParserIDENTIFIER)
}

func (s *ShowStatementContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserIDENTIFIER, i)
}

func (s *ShowStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ShowStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ShowStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterShowStatement(s)
	}
}

func (s *ShowStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitShowStatement(s)
	}
}

func (p *sqlParser) ShowStatement() (localctx IShowStatementContext) {
	localctx = NewShowStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, sqlParserRULE_showStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(326)
		p.Match(sqlParserSHOW)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(327)
		p.Match(sqlParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(331)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserIDENTIFIER {
		{
			p.SetState(328)
			p.Match(sqlParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(333)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISelectStatementContext is an interface to support dynamic dispatch.
type ISelectStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	UnionExceptStatement() IUnionExceptStatementContext

	// IsSelectStatementContext differentiates from other interfaces.
	IsSelectStatementContext()
}

type SelectStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectStatementContext() *SelectStatementContext {
	var p = new(SelectStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_selectStatement
	return p
}

func InitEmptySelectStatementContext(p *SelectStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_selectStatement
}

func (*SelectStatementContext) IsSelectStatementContext() {}

func NewSelectStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectStatementContext {
	var p = new(SelectStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_selectStatement

	return p
}

func (s *SelectStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectStatementContext) UnionExceptStatement() IUnionExceptStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnionExceptStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnionExceptStatementContext)
}

func (s *SelectStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SelectStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterSelectStatement(s)
	}
}

func (s *SelectStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitSelectStatement(s)
	}
}

func (p *sqlParser) SelectStatement() (localctx ISelectStatementContext) {
	localctx = NewSelectStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, sqlParserRULE_selectStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(334)
		p.UnionExceptStatement()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUnionExceptStatementContext is an interface to support dynamic dispatch.
type IUnionExceptStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllIntersectStatement() []IIntersectStatementContext
	IntersectStatement(i int) IIntersectStatementContext
	AllUNION() []antlr.TerminalNode
	UNION(i int) antlr.TerminalNode
	AllEXCEPT() []antlr.TerminalNode
	EXCEPT(i int) antlr.TerminalNode
	AllALL() []antlr.TerminalNode
	ALL(i int) antlr.TerminalNode

	// IsUnionExceptStatementContext differentiates from other interfaces.
	IsUnionExceptStatementContext()
}

type UnionExceptStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnionExceptStatementContext() *UnionExceptStatementContext {
	var p = new(UnionExceptStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_unionExceptStatement
	return p
}

func InitEmptyUnionExceptStatementContext(p *UnionExceptStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_unionExceptStatement
}

func (*UnionExceptStatementContext) IsUnionExceptStatementContext() {}

func NewUnionExceptStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnionExceptStatementContext {
	var p = new(UnionExceptStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_unionExceptStatement

	return p
}

func (s *UnionExceptStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *UnionExceptStatementContext) AllIntersectStatement() []IIntersectStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IIntersectStatementContext); ok {
			len++
		}
	}

	tst := make([]IIntersectStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IIntersectStatementContext); ok {
			tst[i] = t.(IIntersectStatementContext)
			i++
		}
	}

	return tst
}

func (s *UnionExceptStatementContext) IntersectStatement(i int) IIntersectStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntersectStatementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntersectStatementContext)
}

func (s *UnionExceptStatementContext) AllUNION() []antlr.TerminalNode {
	return s.GetTokens(sqlParserUNION)
}

func (s *UnionExceptStatementContext) UNION(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserUNION, i)
}

func (s *UnionExceptStatementContext) AllEXCEPT() []antlr.TerminalNode {
	return s.GetTokens(sqlParserEXCEPT)
}

func (s *UnionExceptStatementContext) EXCEPT(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserEXCEPT, i)
}

func (s *UnionExceptStatementContext) AllALL() []antlr.TerminalNode {
	return s.GetTokens(sqlParserALL)
}

func (s *UnionExceptStatementContext) ALL(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserALL, i)
}

func (s *UnionExceptStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnionExceptStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnionExceptStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterUnionExceptStatement(s)
	}
}

func (s *UnionExceptStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitUnionExceptStatement(s)
	}
}

func (p *sqlParser) UnionExceptStatement() (localctx IUnionExceptStatementContext) {
	localctx = NewUnionExceptStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, sqlParserRULE_unionExceptStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(336)
		p.IntersectStatement()
	}
	p.SetState(344)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserUNION || _la == sqlParserEXCEPT {
		{
			p.SetState(337)
			_la = p.GetTokenStream().LA(1)

			if !(_la == sqlParserUNION || _la == sqlParserEXCEPT) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		p.SetState(339)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserALL {
			{
				p.SetState(338)
				p.Match(sqlParserALL)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(341)
			p.IntersectStatement()
		}

		p.SetState(346)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IIntersectStatementContext is an interface to support dynamic dispatch.
type IIntersectStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllPrimarySelectStatement() []IPrimarySelectStatementContext
	PrimarySelectStatement(i int) IPrimarySelectStatementContext
	AllINTERSECT() []antlr.TerminalNode
	INTERSECT(i int) antlr.TerminalNode
	AllALL() []antlr.TerminalNode
	ALL(i int) antlr.TerminalNode

	// IsIntersectStatementContext differentiates from other interfaces.
	IsIntersectStatementContext()
}

type IntersectStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIntersectStatementContext() *IntersectStatementContext {
	var p = new(IntersectStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_intersectStatement
	return p
}

func InitEmptyIntersectStatementContext(p *IntersectStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_intersectStatement
}

func (*IntersectStatementContext) IsIntersectStatementContext() {}

func NewIntersectStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *IntersectStatementContext {
	var p = new(IntersectStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_intersectStatement

	return p
}

func (s *IntersectStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *IntersectStatementContext) AllPrimarySelectStatement() []IPrimarySelectStatementContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPrimarySelectStatementContext); ok {
			len++
		}
	}

	tst := make([]IPrimarySelectStatementContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPrimarySelectStatementContext); ok {
			tst[i] = t.(IPrimarySelectStatementContext)
			i++
		}
	}

	return tst
}

func (s *IntersectStatementContext) PrimarySelectStatement(i int) IPrimarySelectStatementContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimarySelectStatementContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimarySelectStatementContext)
}

func (s *IntersectStatementContext) AllINTERSECT() []antlr.TerminalNode {
	return s.GetTokens(sqlParserINTERSECT)
}

func (s *IntersectStatementContext) INTERSECT(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserINTERSECT, i)
}

func (s *IntersectStatementContext) AllALL() []antlr.TerminalNode {
	return s.GetTokens(sqlParserALL)
}

func (s *IntersectStatementContext) ALL(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserALL, i)
}

func (s *IntersectStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IntersectStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *IntersectStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterIntersectStatement(s)
	}
}

func (s *IntersectStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitIntersectStatement(s)
	}
}

func (p *sqlParser) IntersectStatement() (localctx IIntersectStatementContext) {
	localctx = NewIntersectStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, sqlParserRULE_intersectStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(347)
		p.PrimarySelectStatement()
	}
	p.SetState(355)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserINTERSECT {
		{
			p.SetState(348)
			p.Match(sqlParserINTERSECT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(350)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserALL {
			{
				p.SetState(349)
				p.Match(sqlParserALL)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(352)
			p.PrimarySelectStatement()
		}

		p.SetState(357)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPrimarySelectStatementContext is an interface to support dynamic dispatch.
type IPrimarySelectStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	SELECT() antlr.TerminalNode
	SelectList() ISelectListContext
	FROM() antlr.TerminalNode
	FromClause() IFromClauseContext
	WHERE() antlr.TerminalNode
	WhereClause() IWhereClauseContext
	GROUP() antlr.TerminalNode
	AllBY() []antlr.TerminalNode
	BY(i int) antlr.TerminalNode
	GroupByClause() IGroupByClauseContext
	ORDER() antlr.TerminalNode
	OrderByClause() IOrderByClauseContext
	LIMIT() antlr.TerminalNode
	LimitValue() ILimitValueContext
	OFFSET() antlr.TerminalNode
	OffsetValue() IOffsetValueContext
	LPAREN() antlr.TerminalNode
	SelectStatement() ISelectStatementContext
	RPAREN() antlr.TerminalNode

	// IsPrimarySelectStatementContext differentiates from other interfaces.
	IsPrimarySelectStatementContext()
}

type PrimarySelectStatementContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimarySelectStatementContext() *PrimarySelectStatementContext {
	var p = new(PrimarySelectStatementContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_primarySelectStatement
	return p
}

func InitEmptyPrimarySelectStatementContext(p *PrimarySelectStatementContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_primarySelectStatement
}

func (*PrimarySelectStatementContext) IsPrimarySelectStatementContext() {}

func NewPrimarySelectStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimarySelectStatementContext {
	var p = new(PrimarySelectStatementContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_primarySelectStatement

	return p
}

func (s *PrimarySelectStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimarySelectStatementContext) SELECT() antlr.TerminalNode {
	return s.GetToken(sqlParserSELECT, 0)
}

func (s *PrimarySelectStatementContext) SelectList() ISelectListContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectListContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectListContext)
}

func (s *PrimarySelectStatementContext) FROM() antlr.TerminalNode {
	return s.GetToken(sqlParserFROM, 0)
}

func (s *PrimarySelectStatementContext) FromClause() IFromClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFromClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFromClauseContext)
}

func (s *PrimarySelectStatementContext) WHERE() antlr.TerminalNode {
	return s.GetToken(sqlParserWHERE, 0)
}

func (s *PrimarySelectStatementContext) WhereClause() IWhereClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWhereClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWhereClauseContext)
}

func (s *PrimarySelectStatementContext) GROUP() antlr.TerminalNode {
	return s.GetToken(sqlParserGROUP, 0)
}

func (s *PrimarySelectStatementContext) AllBY() []antlr.TerminalNode {
	return s.GetTokens(sqlParserBY)
}

func (s *PrimarySelectStatementContext) BY(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserBY, i)
}

func (s *PrimarySelectStatementContext) GroupByClause() IGroupByClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupByClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupByClauseContext)
}

func (s *PrimarySelectStatementContext) ORDER() antlr.TerminalNode {
	return s.GetToken(sqlParserORDER, 0)
}

func (s *PrimarySelectStatementContext) OrderByClause() IOrderByClauseContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOrderByClauseContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOrderByClauseContext)
}

func (s *PrimarySelectStatementContext) LIMIT() antlr.TerminalNode {
	return s.GetToken(sqlParserLIMIT, 0)
}

func (s *PrimarySelectStatementContext) LimitValue() ILimitValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILimitValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILimitValueContext)
}

func (s *PrimarySelectStatementContext) OFFSET() antlr.TerminalNode {
	return s.GetToken(sqlParserOFFSET, 0)
}

func (s *PrimarySelectStatementContext) OffsetValue() IOffsetValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOffsetValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOffsetValueContext)
}

func (s *PrimarySelectStatementContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserLPAREN, 0)
}

func (s *PrimarySelectStatementContext) SelectStatement() ISelectStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectStatementContext)
}

func (s *PrimarySelectStatementContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserRPAREN, 0)
}

func (s *PrimarySelectStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimarySelectStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PrimarySelectStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterPrimarySelectStatement(s)
	}
}

func (s *PrimarySelectStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitPrimarySelectStatement(s)
	}
}

func (p *sqlParser) PrimarySelectStatement() (localctx IPrimarySelectStatementContext) {
	localctx = NewPrimarySelectStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 34, sqlParserRULE_primarySelectStatement)
	var _la int

	p.SetState(390)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case sqlParserSELECT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(358)
			p.Match(sqlParserSELECT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(359)
			p.SelectList()
		}
		p.SetState(362)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserFROM {
			{
				p.SetState(360)
				p.Match(sqlParserFROM)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(361)
				p.FromClause()
			}

		}
		p.SetState(366)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserWHERE {
			{
				p.SetState(364)
				p.Match(sqlParserWHERE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(365)
				p.WhereClause()
			}

		}
		p.SetState(371)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserGROUP {
			{
				p.SetState(368)
				p.Match(sqlParserGROUP)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(369)
				p.Match(sqlParserBY)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(370)
				p.GroupByClause()
			}

		}
		p.SetState(376)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserORDER {
			{
				p.SetState(373)
				p.Match(sqlParserORDER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(374)
				p.Match(sqlParserBY)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(375)
				p.OrderByClause()
			}

		}
		p.SetState(380)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserLIMIT {
			{
				p.SetState(378)
				p.Match(sqlParserLIMIT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(379)
				p.LimitValue()
			}

		}
		p.SetState(384)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserOFFSET {
			{
				p.SetState(382)
				p.Match(sqlParserOFFSET)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(383)
				p.OffsetValue()
			}

		}

	case sqlParserLPAREN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(386)
			p.Match(sqlParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(387)
			p.SelectStatement()
		}
		{
			p.SetState(388)
			p.Match(sqlParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISelectListContext is an interface to support dynamic dispatch.
type ISelectListContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllSelectItem() []ISelectItemContext
	SelectItem(i int) ISelectItemContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsSelectListContext differentiates from other interfaces.
	IsSelectListContext()
}

type SelectListContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectListContext() *SelectListContext {
	var p = new(SelectListContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_selectList
	return p
}

func InitEmptySelectListContext(p *SelectListContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_selectList
}

func (*SelectListContext) IsSelectListContext() {}

func NewSelectListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectListContext {
	var p = new(SelectListContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_selectList

	return p
}

func (s *SelectListContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectListContext) AllSelectItem() []ISelectItemContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISelectItemContext); ok {
			len++
		}
	}

	tst := make([]ISelectItemContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISelectItemContext); ok {
			tst[i] = t.(ISelectItemContext)
			i++
		}
	}

	return tst
}

func (s *SelectListContext) SelectItem(i int) ISelectItemContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectItemContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectItemContext)
}

func (s *SelectListContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(sqlParserCOMMA)
}

func (s *SelectListContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserCOMMA, i)
}

func (s *SelectListContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectListContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SelectListContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterSelectList(s)
	}
}

func (s *SelectListContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitSelectList(s)
	}
}

func (p *sqlParser) SelectList() (localctx ISelectListContext) {
	localctx = NewSelectListContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 36, sqlParserRULE_selectList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(392)
		p.SelectItem()
	}
	p.SetState(397)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserCOMMA {
		{
			p.SetState(393)
			p.Match(sqlParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(394)
			p.SelectItem()
		}

		p.SetState(399)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISelectItemContext is an interface to support dynamic dispatch.
type ISelectItemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	Alias() IAliasContext
	SelectAll() ISelectAllContext

	// IsSelectItemContext differentiates from other interfaces.
	IsSelectItemContext()
}

type SelectItemContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectItemContext() *SelectItemContext {
	var p = new(SelectItemContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_selectItem
	return p
}

func InitEmptySelectItemContext(p *SelectItemContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_selectItem
}

func (*SelectItemContext) IsSelectItemContext() {}

func NewSelectItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectItemContext {
	var p = new(SelectItemContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_selectItem

	return p
}

func (s *SelectItemContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectItemContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *SelectItemContext) Alias() IAliasContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAliasContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAliasContext)
}

func (s *SelectItemContext) SelectAll() ISelectAllContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectAllContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectAllContext)
}

func (s *SelectItemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectItemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SelectItemContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterSelectItem(s)
	}
}

func (s *SelectItemContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitSelectItem(s)
	}
}

func (p *sqlParser) SelectItem() (localctx ISelectItemContext) {
	localctx = NewSelectItemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 38, sqlParserRULE_selectItem)
	var _la int

	p.SetState(405)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case sqlParserNOT, sqlParserEXISTS, sqlParserNULL, sqlParserDEFAULT, sqlParserCASE, sqlParserTIME, sqlParserEXTRACT, sqlParserINT_TYPE, sqlParserFLOAT_TYPE, sqlParserBOOL_TYPE, sqlParserTEXT_TYPE, sqlParserVARCHAR_TYPE, sqlParserCHAR_TYPE, sqlParserSERIAL_TYPE, sqlParserTIMESTAMP_TYPE, sqlParserDATE_TYPE, sqlParserINTERVAL_TYPE, sqlParserCURRENT_TIMESTAMP, sqlParserCURRENT_USER, sqlParserTRUE, sqlParserFALSE, sqlParserLPAREN, sqlParserPLUS, sqlParserMINUS, sqlParserIDENTIFIER, sqlParserPARAMETER, sqlParserSTRING_LITERAL, sqlParserNUMBER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(400)
			p.Expression()
		}
		p.SetState(402)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserAS || _la == sqlParserIDENTIFIER {
			{
				p.SetState(401)
				p.Alias()
			}

		}

	case sqlParserSTAR:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(404)
			p.SelectAll()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISelectAllContext is an interface to support dynamic dispatch.
type ISelectAllContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STAR() antlr.TerminalNode

	// IsSelectAllContext differentiates from other interfaces.
	IsSelectAllContext()
}

type SelectAllContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectAllContext() *SelectAllContext {
	var p = new(SelectAllContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_selectAll
	return p
}

func InitEmptySelectAllContext(p *SelectAllContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_selectAll
}

func (*SelectAllContext) IsSelectAllContext() {}

func NewSelectAllContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectAllContext {
	var p = new(SelectAllContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_selectAll

	return p
}

func (s *SelectAllContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectAllContext) STAR() antlr.TerminalNode {
	return s.GetToken(sqlParserSTAR, 0)
}

func (s *SelectAllContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectAllContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SelectAllContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterSelectAll(s)
	}
}

func (s *SelectAllContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitSelectAll(s)
	}
}

func (p *sqlParser) SelectAll() (localctx ISelectAllContext) {
	localctx = NewSelectAllContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 40, sqlParserRULE_selectAll)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(407)
		p.Match(sqlParserSTAR)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFromClauseContext is an interface to support dynamic dispatch.
type IFromClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllTableFactor() []ITableFactorContext
	TableFactor(i int) ITableFactorContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode
	AllJoinClause() []IJoinClauseContext
	JoinClause(i int) IJoinClauseContext

	// IsFromClauseContext differentiates from other interfaces.
	IsFromClauseContext()
}

type FromClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFromClauseContext() *FromClauseContext {
	var p = new(FromClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_fromClause
	return p
}

func InitEmptyFromClauseContext(p *FromClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_fromClause
}

func (*FromClauseContext) IsFromClauseContext() {}

func NewFromClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FromClauseContext {
	var p = new(FromClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_fromClause

	return p
}

func (s *FromClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *FromClauseContext) AllTableFactor() []ITableFactorContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITableFactorContext); ok {
			len++
		}
	}

	tst := make([]ITableFactorContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITableFactorContext); ok {
			tst[i] = t.(ITableFactorContext)
			i++
		}
	}

	return tst
}

func (s *FromClauseContext) TableFactor(i int) ITableFactorContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableFactorContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableFactorContext)
}

func (s *FromClauseContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(sqlParserCOMMA)
}

func (s *FromClauseContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserCOMMA, i)
}

func (s *FromClauseContext) AllJoinClause() []IJoinClauseContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IJoinClauseContext); ok {
			len++
		}
	}

	tst := make([]IJoinClauseContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IJoinClauseContext); ok {
			tst[i] = t.(IJoinClauseContext)
			i++
		}
	}

	return tst
}

func (s *FromClauseContext) JoinClause(i int) IJoinClauseContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IJoinClauseContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IJoinClauseContext)
}

func (s *FromClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FromClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FromClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterFromClause(s)
	}
}

func (s *FromClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitFromClause(s)
	}
}

func (p *sqlParser) FromClause() (localctx IFromClauseContext) {
	localctx = NewFromClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 42, sqlParserRULE_fromClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(409)
		p.TableFactor()
	}
	p.SetState(414)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserCOMMA {
		{
			p.SetState(410)
			p.Match(sqlParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(411)
			p.TableFactor()
		}

		p.SetState(416)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(420)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&541165879296) != 0 {
		{
			p.SetState(417)
			p.JoinClause()
		}

		p.SetState(422)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITableFactorContext is an interface to support dynamic dispatch.
type ITableFactorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	TableName() ITableNameContext
	Alias() IAliasContext
	LPAREN() antlr.TerminalNode
	SelectStatement() ISelectStatementContext
	RPAREN() antlr.TerminalNode

	// IsTableFactorContext differentiates from other interfaces.
	IsTableFactorContext()
}

type TableFactorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTableFactorContext() *TableFactorContext {
	var p = new(TableFactorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_tableFactor
	return p
}

func InitEmptyTableFactorContext(p *TableFactorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_tableFactor
}

func (*TableFactorContext) IsTableFactorContext() {}

func NewTableFactorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TableFactorContext {
	var p = new(TableFactorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_tableFactor

	return p
}

func (s *TableFactorContext) GetParser() antlr.Parser { return s.parser }

func (s *TableFactorContext) TableName() ITableNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableNameContext)
}

func (s *TableFactorContext) Alias() IAliasContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAliasContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAliasContext)
}

func (s *TableFactorContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserLPAREN, 0)
}

func (s *TableFactorContext) SelectStatement() ISelectStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectStatementContext)
}

func (s *TableFactorContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserRPAREN, 0)
}

func (s *TableFactorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TableFactorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TableFactorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterTableFactor(s)
	}
}

func (s *TableFactorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitTableFactor(s)
	}
}

func (p *sqlParser) TableFactor() (localctx ITableFactorContext) {
	localctx = NewTableFactorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 44, sqlParserRULE_tableFactor)
	var _la int

	p.SetState(432)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case sqlParserDEFAULT, sqlParserTIME, sqlParserINT_TYPE, sqlParserFLOAT_TYPE, sqlParserBOOL_TYPE, sqlParserTEXT_TYPE, sqlParserVARCHAR_TYPE, sqlParserCHAR_TYPE, sqlParserSERIAL_TYPE, sqlParserTIMESTAMP_TYPE, sqlParserDATE_TYPE, sqlParserINTERVAL_TYPE, sqlParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(423)
			p.TableName()
		}
		p.SetState(425)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserAS || _la == sqlParserIDENTIFIER {
			{
				p.SetState(424)
				p.Alias()
			}

		}

	case sqlParserLPAREN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(427)
			p.Match(sqlParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(428)
			p.SelectStatement()
		}
		{
			p.SetState(429)
			p.Match(sqlParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(430)
			p.Alias()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IJoinClauseContext is an interface to support dynamic dispatch.
type IJoinClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	JOIN() antlr.TerminalNode
	QualifiedName() IQualifiedNameContext
	ON() antlr.TerminalNode
	Expression() IExpressionContext
	Alias() IAliasContext
	INNER() antlr.TerminalNode
	LEFT() antlr.TerminalNode
	RIGHT() antlr.TerminalNode
	FULL() antlr.TerminalNode
	CROSS() antlr.TerminalNode

	// IsJoinClauseContext differentiates from other interfaces.
	IsJoinClauseContext()
}

type JoinClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyJoinClauseContext() *JoinClauseContext {
	var p = new(JoinClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_joinClause
	return p
}

func InitEmptyJoinClauseContext(p *JoinClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_joinClause
}

func (*JoinClauseContext) IsJoinClauseContext() {}

func NewJoinClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *JoinClauseContext {
	var p = new(JoinClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_joinClause

	return p
}

func (s *JoinClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *JoinClauseContext) JOIN() antlr.TerminalNode {
	return s.GetToken(sqlParserJOIN, 0)
}

func (s *JoinClauseContext) QualifiedName() IQualifiedNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedNameContext)
}

func (s *JoinClauseContext) ON() antlr.TerminalNode {
	return s.GetToken(sqlParserON, 0)
}

func (s *JoinClauseContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *JoinClauseContext) Alias() IAliasContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAliasContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAliasContext)
}

func (s *JoinClauseContext) INNER() antlr.TerminalNode {
	return s.GetToken(sqlParserINNER, 0)
}

func (s *JoinClauseContext) LEFT() antlr.TerminalNode {
	return s.GetToken(sqlParserLEFT, 0)
}

func (s *JoinClauseContext) RIGHT() antlr.TerminalNode {
	return s.GetToken(sqlParserRIGHT, 0)
}

func (s *JoinClauseContext) FULL() antlr.TerminalNode {
	return s.GetToken(sqlParserFULL, 0)
}

func (s *JoinClauseContext) CROSS() antlr.TerminalNode {
	return s.GetToken(sqlParserCROSS, 0)
}

func (s *JoinClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *JoinClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *JoinClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterJoinClause(s)
	}
}

func (s *JoinClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitJoinClause(s)
	}
}

func (p *sqlParser) JoinClause() (localctx IJoinClauseContext) {
	localctx = NewJoinClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 46, sqlParserRULE_joinClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(435)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&532575944704) != 0 {
		{
			p.SetState(434)
			_la = p.GetTokenStream().LA(1)

			if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&532575944704) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}
	{
		p.SetState(437)
		p.Match(sqlParserJOIN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(438)
		p.QualifiedName()
	}
	p.SetState(440)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserAS || _la == sqlParserIDENTIFIER {
		{
			p.SetState(439)
			p.Alias()
		}

	}
	{
		p.SetState(442)
		p.Match(sqlParserON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(443)
		p.Expression()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IWhereClauseContext is an interface to support dynamic dispatch.
type IWhereClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext

	// IsWhereClauseContext differentiates from other interfaces.
	IsWhereClauseContext()
}

type WhereClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhereClauseContext() *WhereClauseContext {
	var p = new(WhereClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_whereClause
	return p
}

func InitEmptyWhereClauseContext(p *WhereClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_whereClause
}

func (*WhereClauseContext) IsWhereClauseContext() {}

func NewWhereClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhereClauseContext {
	var p = new(WhereClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_whereClause

	return p
}

func (s *WhereClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *WhereClauseContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *WhereClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhereClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WhereClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterWhereClause(s)
	}
}

func (s *WhereClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitWhereClause(s)
	}
}

func (p *sqlParser) WhereClause() (localctx IWhereClauseContext) {
	localctx = NewWhereClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 48, sqlParserRULE_whereClause)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(445)
		p.Expression()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOrderByClauseContext is an interface to support dynamic dispatch.
type IOrderByClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllOrderByItem() []IOrderByItemContext
	OrderByItem(i int) IOrderByItemContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsOrderByClauseContext differentiates from other interfaces.
	IsOrderByClauseContext()
}

type OrderByClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOrderByClauseContext() *OrderByClauseContext {
	var p = new(OrderByClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_orderByClause
	return p
}

func InitEmptyOrderByClauseContext(p *OrderByClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_orderByClause
}

func (*OrderByClauseContext) IsOrderByClauseContext() {}

func NewOrderByClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OrderByClauseContext {
	var p = new(OrderByClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_orderByClause

	return p
}

func (s *OrderByClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *OrderByClauseContext) AllOrderByItem() []IOrderByItemContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IOrderByItemContext); ok {
			len++
		}
	}

	tst := make([]IOrderByItemContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IOrderByItemContext); ok {
			tst[i] = t.(IOrderByItemContext)
			i++
		}
	}

	return tst
}

func (s *OrderByClauseContext) OrderByItem(i int) IOrderByItemContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOrderByItemContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOrderByItemContext)
}

func (s *OrderByClauseContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(sqlParserCOMMA)
}

func (s *OrderByClauseContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserCOMMA, i)
}

func (s *OrderByClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrderByClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OrderByClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterOrderByClause(s)
	}
}

func (s *OrderByClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitOrderByClause(s)
	}
}

func (p *sqlParser) OrderByClause() (localctx IOrderByClauseContext) {
	localctx = NewOrderByClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 50, sqlParserRULE_orderByClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(447)
		p.OrderByItem()
	}
	p.SetState(452)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserCOMMA {
		{
			p.SetState(448)
			p.Match(sqlParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(449)
			p.OrderByItem()
		}

		p.SetState(454)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOrderByItemContext is an interface to support dynamic dispatch.
type IOrderByItemContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	ASC() antlr.TerminalNode
	DESC() antlr.TerminalNode

	// IsOrderByItemContext differentiates from other interfaces.
	IsOrderByItemContext()
}

type OrderByItemContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOrderByItemContext() *OrderByItemContext {
	var p = new(OrderByItemContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_orderByItem
	return p
}

func InitEmptyOrderByItemContext(p *OrderByItemContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_orderByItem
}

func (*OrderByItemContext) IsOrderByItemContext() {}

func NewOrderByItemContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OrderByItemContext {
	var p = new(OrderByItemContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_orderByItem

	return p
}

func (s *OrderByItemContext) GetParser() antlr.Parser { return s.parser }

func (s *OrderByItemContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *OrderByItemContext) ASC() antlr.TerminalNode {
	return s.GetToken(sqlParserASC, 0)
}

func (s *OrderByItemContext) DESC() antlr.TerminalNode {
	return s.GetToken(sqlParserDESC, 0)
}

func (s *OrderByItemContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrderByItemContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OrderByItemContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterOrderByItem(s)
	}
}

func (s *OrderByItemContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitOrderByItem(s)
	}
}

func (p *sqlParser) OrderByItem() (localctx IOrderByItemContext) {
	localctx = NewOrderByItemContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 52, sqlParserRULE_orderByItem)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(455)
		p.Expression()
	}
	p.SetState(457)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserASC || _la == sqlParserDESC {
		{
			p.SetState(456)
			_la = p.GetTokenStream().LA(1)

			if !(_la == sqlParserASC || _la == sqlParserDESC) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ILimitValueContext is an interface to support dynamic dispatch.
type ILimitValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NUMBER() antlr.TerminalNode

	// IsLimitValueContext differentiates from other interfaces.
	IsLimitValueContext()
}

type LimitValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLimitValueContext() *LimitValueContext {
	var p = new(LimitValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_limitValue
	return p
}

func InitEmptyLimitValueContext(p *LimitValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_limitValue
}

func (*LimitValueContext) IsLimitValueContext() {}

func NewLimitValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LimitValueContext {
	var p = new(LimitValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_limitValue

	return p
}

func (s *LimitValueContext) GetParser() antlr.Parser { return s.parser }

func (s *LimitValueContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(sqlParserNUMBER, 0)
}

func (s *LimitValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LimitValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LimitValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterLimitValue(s)
	}
}

func (s *LimitValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitLimitValue(s)
	}
}

func (p *sqlParser) LimitValue() (localctx ILimitValueContext) {
	localctx = NewLimitValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 54, sqlParserRULE_limitValue)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(459)
		p.Match(sqlParserNUMBER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOffsetValueContext is an interface to support dynamic dispatch.
type IOffsetValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	NUMBER() antlr.TerminalNode

	// IsOffsetValueContext differentiates from other interfaces.
	IsOffsetValueContext()
}

type OffsetValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOffsetValueContext() *OffsetValueContext {
	var p = new(OffsetValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_offsetValue
	return p
}

func InitEmptyOffsetValueContext(p *OffsetValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_offsetValue
}

func (*OffsetValueContext) IsOffsetValueContext() {}

func NewOffsetValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OffsetValueContext {
	var p = new(OffsetValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_offsetValue

	return p
}

func (s *OffsetValueContext) GetParser() antlr.Parser { return s.parser }

func (s *OffsetValueContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(sqlParserNUMBER, 0)
}

func (s *OffsetValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OffsetValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OffsetValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterOffsetValue(s)
	}
}

func (s *OffsetValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitOffsetValue(s)
	}
}

func (p *sqlParser) OffsetValue() (localctx IOffsetValueContext) {
	localctx = NewOffsetValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 56, sqlParserRULE_offsetValue)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(461)
		p.Match(sqlParserNUMBER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IGroupByClauseContext is an interface to support dynamic dispatch.
type IGroupByClauseContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsGroupByClauseContext differentiates from other interfaces.
	IsGroupByClauseContext()
}

type GroupByClauseContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupByClauseContext() *GroupByClauseContext {
	var p = new(GroupByClauseContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_groupByClause
	return p
}

func InitEmptyGroupByClauseContext(p *GroupByClauseContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_groupByClause
}

func (*GroupByClauseContext) IsGroupByClauseContext() {}

func NewGroupByClauseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupByClauseContext {
	var p = new(GroupByClauseContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_groupByClause

	return p
}

func (s *GroupByClauseContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupByClauseContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *GroupByClauseContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *GroupByClauseContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(sqlParserCOMMA)
}

func (s *GroupByClauseContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserCOMMA, i)
}

func (s *GroupByClauseContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupByClauseContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GroupByClauseContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterGroupByClause(s)
	}
}

func (s *GroupByClauseContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitGroupByClause(s)
	}
}

func (p *sqlParser) GroupByClause() (localctx IGroupByClauseContext) {
	localctx = NewGroupByClauseContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 58, sqlParserRULE_groupByClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(463)
		p.Expression()
	}
	p.SetState(468)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserCOMMA {
		{
			p.SetState(464)
			p.Match(sqlParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(465)
			p.Expression()
		}

		p.SetState(470)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAliasContext is an interface to support dynamic dispatch.
type IAliasContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	AS() antlr.TerminalNode

	// IsAliasContext differentiates from other interfaces.
	IsAliasContext()
}

type AliasContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAliasContext() *AliasContext {
	var p = new(AliasContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_alias
	return p
}

func InitEmptyAliasContext(p *AliasContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_alias
}

func (*AliasContext) IsAliasContext() {}

func NewAliasContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AliasContext {
	var p = new(AliasContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_alias

	return p
}

func (s *AliasContext) GetParser() antlr.Parser { return s.parser }

func (s *AliasContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(sqlParserIDENTIFIER, 0)
}

func (s *AliasContext) AS() antlr.TerminalNode {
	return s.GetToken(sqlParserAS, 0)
}

func (s *AliasContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AliasContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AliasContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterAlias(s)
	}
}

func (s *AliasContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitAlias(s)
	}
}

func (p *sqlParser) Alias() (localctx IAliasContext) {
	localctx = NewAliasContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 60, sqlParserRULE_alias)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(472)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserAS {
		{
			p.SetState(471)
			p.Match(sqlParserAS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(474)
		p.Match(sqlParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExpressionContext is an interface to support dynamic dispatch.
type IExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OrExpression() IOrExpressionContext

	// IsExpressionContext differentiates from other interfaces.
	IsExpressionContext()
}

type ExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExpressionContext() *ExpressionContext {
	var p = new(ExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_expression
	return p
}

func InitEmptyExpressionContext(p *ExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_expression
}

func (*ExpressionContext) IsExpressionContext() {}

func NewExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExpressionContext {
	var p = new(ExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_expression

	return p
}

func (s *ExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExpressionContext) OrExpression() IOrExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOrExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOrExpressionContext)
}

func (s *ExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterExpression(s)
	}
}

func (s *ExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitExpression(s)
	}
}

func (p *sqlParser) Expression() (localctx IExpressionContext) {
	localctx = NewExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 62, sqlParserRULE_expression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(476)
		p.OrExpression()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOrExpressionContext is an interface to support dynamic dispatch.
type IOrExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAndExpression() []IAndExpressionContext
	AndExpression(i int) IAndExpressionContext
	AllOR() []antlr.TerminalNode
	OR(i int) antlr.TerminalNode

	// IsOrExpressionContext differentiates from other interfaces.
	IsOrExpressionContext()
}

type OrExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOrExpressionContext() *OrExpressionContext {
	var p = new(OrExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_orExpression
	return p
}

func InitEmptyOrExpressionContext(p *OrExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_orExpression
}

func (*OrExpressionContext) IsOrExpressionContext() {}

func NewOrExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OrExpressionContext {
	var p = new(OrExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_orExpression

	return p
}

func (s *OrExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *OrExpressionContext) AllAndExpression() []IAndExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAndExpressionContext); ok {
			len++
		}
	}

	tst := make([]IAndExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAndExpressionContext); ok {
			tst[i] = t.(IAndExpressionContext)
			i++
		}
	}

	return tst
}

func (s *OrExpressionContext) AndExpression(i int) IAndExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAndExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAndExpressionContext)
}

func (s *OrExpressionContext) AllOR() []antlr.TerminalNode {
	return s.GetTokens(sqlParserOR)
}

func (s *OrExpressionContext) OR(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserOR, i)
}

func (s *OrExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OrExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OrExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterOrExpression(s)
	}
}

func (s *OrExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitOrExpression(s)
	}
}

func (p *sqlParser) OrExpression() (localctx IOrExpressionContext) {
	localctx = NewOrExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 64, sqlParserRULE_orExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(478)
		p.AndExpression()
	}
	p.SetState(483)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserOR {
		{
			p.SetState(479)
			p.Match(sqlParserOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(480)
			p.AndExpression()
		}

		p.SetState(485)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAndExpressionContext is an interface to support dynamic dispatch.
type IAndExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNotExpression() []INotExpressionContext
	NotExpression(i int) INotExpressionContext
	AllAND() []antlr.TerminalNode
	AND(i int) antlr.TerminalNode

	// IsAndExpressionContext differentiates from other interfaces.
	IsAndExpressionContext()
}

type AndExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAndExpressionContext() *AndExpressionContext {
	var p = new(AndExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_andExpression
	return p
}

func InitEmptyAndExpressionContext(p *AndExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_andExpression
}

func (*AndExpressionContext) IsAndExpressionContext() {}

func NewAndExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AndExpressionContext {
	var p = new(AndExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_andExpression

	return p
}

func (s *AndExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *AndExpressionContext) AllNotExpression() []INotExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(INotExpressionContext); ok {
			len++
		}
	}

	tst := make([]INotExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(INotExpressionContext); ok {
			tst[i] = t.(INotExpressionContext)
			i++
		}
	}

	return tst
}

func (s *AndExpressionContext) NotExpression(i int) INotExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INotExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(INotExpressionContext)
}

func (s *AndExpressionContext) AllAND() []antlr.TerminalNode {
	return s.GetTokens(sqlParserAND)
}

func (s *AndExpressionContext) AND(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserAND, i)
}

func (s *AndExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AndExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AndExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterAndExpression(s)
	}
}

func (s *AndExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitAndExpression(s)
	}
}

func (p *sqlParser) AndExpression() (localctx IAndExpressionContext) {
	localctx = NewAndExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 66, sqlParserRULE_andExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(486)
		p.NotExpression()
	}
	p.SetState(491)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserAND {
		{
			p.SetState(487)
			p.Match(sqlParserAND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(488)
			p.NotExpression()
		}

		p.SetState(493)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INotExpressionContext is an interface to support dynamic dispatch.
type INotExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	ComparisonExpression() IComparisonExpressionContext
	NOT() antlr.TerminalNode
	EXISTS() antlr.TerminalNode
	SubqueryExpression() ISubqueryExpressionContext

	// IsNotExpressionContext differentiates from other interfaces.
	IsNotExpressionContext()
}

type NotExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNotExpressionContext() *NotExpressionContext {
	var p = new(NotExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_notExpression
	return p
}

func InitEmptyNotExpressionContext(p *NotExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_notExpression
}

func (*NotExpressionContext) IsNotExpressionContext() {}

func NewNotExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NotExpressionContext {
	var p = new(NotExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_notExpression

	return p
}

func (s *NotExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *NotExpressionContext) ComparisonExpression() IComparisonExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComparisonExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComparisonExpressionContext)
}

func (s *NotExpressionContext) NOT() antlr.TerminalNode {
	return s.GetToken(sqlParserNOT, 0)
}

func (s *NotExpressionContext) EXISTS() antlr.TerminalNode {
	return s.GetToken(sqlParserEXISTS, 0)
}

func (s *NotExpressionContext) SubqueryExpression() ISubqueryExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISubqueryExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISubqueryExpressionContext)
}

func (s *NotExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NotExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NotExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterNotExpression(s)
	}
}

func (s *NotExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitNotExpression(s)
	}
}

func (p *sqlParser) NotExpression() (localctx INotExpressionContext) {
	localctx = NewNotExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 68, sqlParserRULE_notExpression)
	var _la int

	p.SetState(503)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 54, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(495)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserNOT {
			{
				p.SetState(494)
				p.Match(sqlParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(497)
			p.ComparisonExpression()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		p.SetState(499)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserNOT {
			{
				p.SetState(498)
				p.Match(sqlParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(501)
			p.Match(sqlParserEXISTS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(502)
			p.SubqueryExpression()
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IComparisonExpressionContext is an interface to support dynamic dispatch.
type IComparisonExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllConcatExpression() []IConcatExpressionContext
	ConcatExpression(i int) IConcatExpressionContext
	IN() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	RPAREN() antlr.TerminalNode
	SubqueryExpression() ISubqueryExpressionContext
	LIKE() antlr.TerminalNode
	IS() antlr.TerminalNode
	NULL() antlr.TerminalNode
	Operator() IOperatorContext
	OperatorExpr() IOperatorExprContext
	NOT() antlr.TerminalNode
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsComparisonExpressionContext differentiates from other interfaces.
	IsComparisonExpressionContext()
}

type ComparisonExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComparisonExpressionContext() *ComparisonExpressionContext {
	var p = new(ComparisonExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_comparisonExpression
	return p
}

func InitEmptyComparisonExpressionContext(p *ComparisonExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_comparisonExpression
}

func (*ComparisonExpressionContext) IsComparisonExpressionContext() {}

func NewComparisonExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComparisonExpressionContext {
	var p = new(ComparisonExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_comparisonExpression

	return p
}

func (s *ComparisonExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ComparisonExpressionContext) AllConcatExpression() []IConcatExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IConcatExpressionContext); ok {
			len++
		}
	}

	tst := make([]IConcatExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IConcatExpressionContext); ok {
			tst[i] = t.(IConcatExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ComparisonExpressionContext) ConcatExpression(i int) IConcatExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IConcatExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IConcatExpressionContext)
}

func (s *ComparisonExpressionContext) IN() antlr.TerminalNode {
	return s.GetToken(sqlParserIN, 0)
}

func (s *ComparisonExpressionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserLPAREN, 0)
}

func (s *ComparisonExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ComparisonExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ComparisonExpressionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserRPAREN, 0)
}

func (s *ComparisonExpressionContext) SubqueryExpression() ISubqueryExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISubqueryExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISubqueryExpressionContext)
}

func (s *ComparisonExpressionContext) LIKE() antlr.TerminalNode {
	return s.GetToken(sqlParserLIKE, 0)
}

func (s *ComparisonExpressionContext) IS() antlr.TerminalNode {
	return s.GetToken(sqlParserIS, 0)
}

func (s *ComparisonExpressionContext) NULL() antlr.TerminalNode {
	return s.GetToken(sqlParserNULL, 0)
}

func (s *ComparisonExpressionContext) Operator() IOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOperatorContext)
}

func (s *ComparisonExpressionContext) OperatorExpr() IOperatorExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IOperatorExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IOperatorExprContext)
}

func (s *ComparisonExpressionContext) NOT() antlr.TerminalNode {
	return s.GetToken(sqlParserNOT, 0)
}

func (s *ComparisonExpressionContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(sqlParserCOMMA)
}

func (s *ComparisonExpressionContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserCOMMA, i)
}

func (s *ComparisonExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComparisonExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterComparisonExpression(s)
	}
}

func (s *ComparisonExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitComparisonExpression(s)
	}
}

func (p *sqlParser) ComparisonExpression() (localctx IComparisonExpressionContext) {
	localctx = NewComparisonExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 70, sqlParserRULE_comparisonExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(505)
		p.ConcatExpression()
	}
	p.SetState(542)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 61, p.GetParserRuleContext()) == 1 {
		p.SetState(508)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case sqlParserTILDE, sqlParserNREGEX, sqlParserIREGEX, sqlParserNIREGEX, sqlParserEQ, sqlParserGT, sqlParserLT, sqlParserGE, sqlParserLE, sqlParserNE:
			{
				p.SetState(506)
				p.Operator()
			}

		case sqlParserOPERATOR_KW:
			{
				p.SetState(507)
				p.OperatorExpr()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}
		{
			p.SetState(510)
			p.ConcatExpression()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	} else if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 61, p.GetParserRuleContext()) == 2 {
		p.SetState(513)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserNOT {
			{
				p.SetState(512)
				p.Match(sqlParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(515)
			p.Match(sqlParserIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(516)
			p.Match(sqlParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(517)
			p.Expression()
		}
		p.SetState(522)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == sqlParserCOMMA {
			{
				p.SetState(518)
				p.Match(sqlParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(519)
				p.Expression()
			}

			p.SetState(524)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(525)
			p.Match(sqlParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	} else if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 61, p.GetParserRuleContext()) == 3 {
		p.SetState(528)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserNOT {
			{
				p.SetState(527)
				p.Match(sqlParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(530)
			p.Match(sqlParserIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(531)
			p.SubqueryExpression()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	} else if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 61, p.GetParserRuleContext()) == 4 {
		p.SetState(533)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserNOT {
			{
				p.SetState(532)
				p.Match(sqlParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(535)
			p.Match(sqlParserLIKE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(536)
			p.ConcatExpression()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	} else if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 61, p.GetParserRuleContext()) == 5 {
		{
			p.SetState(537)
			p.Match(sqlParserIS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(539)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserNOT {
			{
				p.SetState(538)
				p.Match(sqlParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(541)
			p.Match(sqlParserNULL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IConcatExpressionContext is an interface to support dynamic dispatch.
type IConcatExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllAdditiveExpression() []IAdditiveExpressionContext
	AdditiveExpression(i int) IAdditiveExpressionContext
	AllCONCAT() []antlr.TerminalNode
	CONCAT(i int) antlr.TerminalNode

	// IsConcatExpressionContext differentiates from other interfaces.
	IsConcatExpressionContext()
}

type ConcatExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyConcatExpressionContext() *ConcatExpressionContext {
	var p = new(ConcatExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_concatExpression
	return p
}

func InitEmptyConcatExpressionContext(p *ConcatExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_concatExpression
}

func (*ConcatExpressionContext) IsConcatExpressionContext() {}

func NewConcatExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConcatExpressionContext {
	var p = new(ConcatExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_concatExpression

	return p
}

func (s *ConcatExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *ConcatExpressionContext) AllAdditiveExpression() []IAdditiveExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAdditiveExpressionContext); ok {
			len++
		}
	}

	tst := make([]IAdditiveExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAdditiveExpressionContext); ok {
			tst[i] = t.(IAdditiveExpressionContext)
			i++
		}
	}

	return tst
}

func (s *ConcatExpressionContext) AdditiveExpression(i int) IAdditiveExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAdditiveExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAdditiveExpressionContext)
}

func (s *ConcatExpressionContext) AllCONCAT() []antlr.TerminalNode {
	return s.GetTokens(sqlParserCONCAT)
}

func (s *ConcatExpressionContext) CONCAT(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserCONCAT, i)
}

func (s *ConcatExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ConcatExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ConcatExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterConcatExpression(s)
	}
}

func (s *ConcatExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitConcatExpression(s)
	}
}

func (p *sqlParser) ConcatExpression() (localctx IConcatExpressionContext) {
	localctx = NewConcatExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 72, sqlParserRULE_concatExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(544)
		p.AdditiveExpression()
	}
	p.SetState(549)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserCONCAT {
		{
			p.SetState(545)
			p.Match(sqlParserCONCAT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(546)
			p.AdditiveExpression()
		}

		p.SetState(551)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IAdditiveExpressionContext is an interface to support dynamic dispatch.
type IAdditiveExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllMultiplicativeExpression() []IMultiplicativeExpressionContext
	MultiplicativeExpression(i int) IMultiplicativeExpressionContext
	AllPLUS() []antlr.TerminalNode
	PLUS(i int) antlr.TerminalNode
	AllMINUS() []antlr.TerminalNode
	MINUS(i int) antlr.TerminalNode

	// IsAdditiveExpressionContext differentiates from other interfaces.
	IsAdditiveExpressionContext()
}

type AdditiveExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAdditiveExpressionContext() *AdditiveExpressionContext {
	var p = new(AdditiveExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_additiveExpression
	return p
}

func InitEmptyAdditiveExpressionContext(p *AdditiveExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_additiveExpression
}

func (*AdditiveExpressionContext) IsAdditiveExpressionContext() {}

func NewAdditiveExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AdditiveExpressionContext {
	var p = new(AdditiveExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_additiveExpression

	return p
}

func (s *AdditiveExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *AdditiveExpressionContext) AllMultiplicativeExpression() []IMultiplicativeExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IMultiplicativeExpressionContext); ok {
			len++
		}
	}

	tst := make([]IMultiplicativeExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IMultiplicativeExpressionContext); ok {
			tst[i] = t.(IMultiplicativeExpressionContext)
			i++
		}
	}

	return tst
}

func (s *AdditiveExpressionContext) MultiplicativeExpression(i int) IMultiplicativeExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IMultiplicativeExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IMultiplicativeExpressionContext)
}

func (s *AdditiveExpressionContext) AllPLUS() []antlr.TerminalNode {
	return s.GetTokens(sqlParserPLUS)
}

func (s *AdditiveExpressionContext) PLUS(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserPLUS, i)
}

func (s *AdditiveExpressionContext) AllMINUS() []antlr.TerminalNode {
	return s.GetTokens(sqlParserMINUS)
}

func (s *AdditiveExpressionContext) MINUS(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserMINUS, i)
}

func (s *AdditiveExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AdditiveExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AdditiveExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterAdditiveExpression(s)
	}
}

func (s *AdditiveExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitAdditiveExpression(s)
	}
}

func (p *sqlParser) AdditiveExpression() (localctx IAdditiveExpressionContext) {
	localctx = NewAdditiveExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 74, sqlParserRULE_additiveExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(552)
		p.MultiplicativeExpression()
	}
	p.SetState(557)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserPLUS || _la == sqlParserMINUS {
		{
			p.SetState(553)
			_la = p.GetTokenStream().LA(1)

			if !(_la == sqlParserPLUS || _la == sqlParserMINUS) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(554)
			p.MultiplicativeExpression()
		}

		p.SetState(559)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IMultiplicativeExpressionContext is an interface to support dynamic dispatch.
type IMultiplicativeExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllUnaryExpression() []IUnaryExpressionContext
	UnaryExpression(i int) IUnaryExpressionContext
	AllSTAR() []antlr.TerminalNode
	STAR(i int) antlr.TerminalNode
	AllSLASH() []antlr.TerminalNode
	SLASH(i int) antlr.TerminalNode

	// IsMultiplicativeExpressionContext differentiates from other interfaces.
	IsMultiplicativeExpressionContext()
}

type MultiplicativeExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyMultiplicativeExpressionContext() *MultiplicativeExpressionContext {
	var p = new(MultiplicativeExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_multiplicativeExpression
	return p
}

func InitEmptyMultiplicativeExpressionContext(p *MultiplicativeExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_multiplicativeExpression
}

func (*MultiplicativeExpressionContext) IsMultiplicativeExpressionContext() {}

func NewMultiplicativeExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MultiplicativeExpressionContext {
	var p = new(MultiplicativeExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_multiplicativeExpression

	return p
}

func (s *MultiplicativeExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *MultiplicativeExpressionContext) AllUnaryExpression() []IUnaryExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IUnaryExpressionContext); ok {
			len++
		}
	}

	tst := make([]IUnaryExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IUnaryExpressionContext); ok {
			tst[i] = t.(IUnaryExpressionContext)
			i++
		}
	}

	return tst
}

func (s *MultiplicativeExpressionContext) UnaryExpression(i int) IUnaryExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IUnaryExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IUnaryExpressionContext)
}

func (s *MultiplicativeExpressionContext) AllSTAR() []antlr.TerminalNode {
	return s.GetTokens(sqlParserSTAR)
}

func (s *MultiplicativeExpressionContext) STAR(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserSTAR, i)
}

func (s *MultiplicativeExpressionContext) AllSLASH() []antlr.TerminalNode {
	return s.GetTokens(sqlParserSLASH)
}

func (s *MultiplicativeExpressionContext) SLASH(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserSLASH, i)
}

func (s *MultiplicativeExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *MultiplicativeExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *MultiplicativeExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterMultiplicativeExpression(s)
	}
}

func (s *MultiplicativeExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitMultiplicativeExpression(s)
	}
}

func (p *sqlParser) MultiplicativeExpression() (localctx IMultiplicativeExpressionContext) {
	localctx = NewMultiplicativeExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 76, sqlParserRULE_multiplicativeExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(560)
		p.UnaryExpression()
	}
	p.SetState(565)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserSTAR || _la == sqlParserSLASH {
		{
			p.SetState(561)
			_la = p.GetTokenStream().LA(1)

			if !(_la == sqlParserSTAR || _la == sqlParserSLASH) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(562)
			p.UnaryExpression()
		}

		p.SetState(567)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IUnaryExpressionContext is an interface to support dynamic dispatch.
type IUnaryExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CastExpression() ICastExpressionContext
	PLUS() antlr.TerminalNode
	MINUS() antlr.TerminalNode

	// IsUnaryExpressionContext differentiates from other interfaces.
	IsUnaryExpressionContext()
}

type UnaryExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyUnaryExpressionContext() *UnaryExpressionContext {
	var p = new(UnaryExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_unaryExpression
	return p
}

func InitEmptyUnaryExpressionContext(p *UnaryExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_unaryExpression
}

func (*UnaryExpressionContext) IsUnaryExpressionContext() {}

func NewUnaryExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnaryExpressionContext {
	var p = new(UnaryExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_unaryExpression

	return p
}

func (s *UnaryExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *UnaryExpressionContext) CastExpression() ICastExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICastExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICastExpressionContext)
}

func (s *UnaryExpressionContext) PLUS() antlr.TerminalNode {
	return s.GetToken(sqlParserPLUS, 0)
}

func (s *UnaryExpressionContext) MINUS() antlr.TerminalNode {
	return s.GetToken(sqlParserMINUS, 0)
}

func (s *UnaryExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *UnaryExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *UnaryExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterUnaryExpression(s)
	}
}

func (s *UnaryExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitUnaryExpression(s)
	}
}

func (p *sqlParser) UnaryExpression() (localctx IUnaryExpressionContext) {
	localctx = NewUnaryExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 78, sqlParserRULE_unaryExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(569)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserPLUS || _la == sqlParserMINUS {
		{
			p.SetState(568)
			_la = p.GetTokenStream().LA(1)

			if !(_la == sqlParserPLUS || _la == sqlParserMINUS) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}
	{
		p.SetState(571)
		p.CastExpression()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICastExpressionContext is an interface to support dynamic dispatch.
type ICastExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	PrimaryExpression() IPrimaryExpressionContext
	AllPostfix() []IPostfixContext
	Postfix(i int) IPostfixContext

	// IsCastExpressionContext differentiates from other interfaces.
	IsCastExpressionContext()
}

type CastExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCastExpressionContext() *CastExpressionContext {
	var p = new(CastExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_castExpression
	return p
}

func InitEmptyCastExpressionContext(p *CastExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_castExpression
}

func (*CastExpressionContext) IsCastExpressionContext() {}

func NewCastExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CastExpressionContext {
	var p = new(CastExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_castExpression

	return p
}

func (s *CastExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *CastExpressionContext) PrimaryExpression() IPrimaryExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPrimaryExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPrimaryExpressionContext)
}

func (s *CastExpressionContext) AllPostfix() []IPostfixContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IPostfixContext); ok {
			len++
		}
	}

	tst := make([]IPostfixContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IPostfixContext); ok {
			tst[i] = t.(IPostfixContext)
			i++
		}
	}

	return tst
}

func (s *CastExpressionContext) Postfix(i int) IPostfixContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IPostfixContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IPostfixContext)
}

func (s *CastExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CastExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CastExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterCastExpression(s)
	}
}

func (s *CastExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitCastExpression(s)
	}
}

func (p *sqlParser) CastExpression() (localctx ICastExpressionContext) {
	localctx = NewCastExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 80, sqlParserRULE_castExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(573)
		p.PrimaryExpression()
	}
	p.SetState(577)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64((_la-52)) & ^0x3f) == 0 && ((int64(1)<<(_la-52))&563087392374785) != 0 {
		{
			p.SetState(574)
			p.Postfix()
		}

		p.SetState(579)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPostfixContext is an interface to support dynamic dispatch.
type IPostfixContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AT() antlr.TerminalNode
	TIME() antlr.TerminalNode
	ZONE() antlr.TerminalNode
	STRING_LITERAL() antlr.TerminalNode
	COLLATE() antlr.TerminalNode
	QualifiedName() IQualifiedNameContext
	AllCOLON_COLON() []antlr.TerminalNode
	COLON_COLON(i int) antlr.TerminalNode
	AllTypeName() []ITypeNameContext
	TypeName(i int) ITypeNameContext

	// IsPostfixContext differentiates from other interfaces.
	IsPostfixContext()
}

type PostfixContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPostfixContext() *PostfixContext {
	var p = new(PostfixContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_postfix
	return p
}

func InitEmptyPostfixContext(p *PostfixContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_postfix
}

func (*PostfixContext) IsPostfixContext() {}

func NewPostfixContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PostfixContext {
	var p = new(PostfixContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_postfix

	return p
}

func (s *PostfixContext) GetParser() antlr.Parser { return s.parser }

func (s *PostfixContext) AT() antlr.TerminalNode {
	return s.GetToken(sqlParserAT, 0)
}

func (s *PostfixContext) TIME() antlr.TerminalNode {
	return s.GetToken(sqlParserTIME, 0)
}

func (s *PostfixContext) ZONE() antlr.TerminalNode {
	return s.GetToken(sqlParserZONE, 0)
}

func (s *PostfixContext) STRING_LITERAL() antlr.TerminalNode {
	return s.GetToken(sqlParserSTRING_LITERAL, 0)
}

func (s *PostfixContext) COLLATE() antlr.TerminalNode {
	return s.GetToken(sqlParserCOLLATE, 0)
}

func (s *PostfixContext) QualifiedName() IQualifiedNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedNameContext)
}

func (s *PostfixContext) AllCOLON_COLON() []antlr.TerminalNode {
	return s.GetTokens(sqlParserCOLON_COLON)
}

func (s *PostfixContext) COLON_COLON(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserCOLON_COLON, i)
}

func (s *PostfixContext) AllTypeName() []ITypeNameContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ITypeNameContext); ok {
			len++
		}
	}

	tst := make([]ITypeNameContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ITypeNameContext); ok {
			tst[i] = t.(ITypeNameContext)
			i++
		}
	}

	return tst
}

func (s *PostfixContext) TypeName(i int) ITypeNameContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypeNameContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypeNameContext)
}

func (s *PostfixContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PostfixContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PostfixContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterPostfix(s)
	}
}

func (s *PostfixContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitPostfix(s)
	}
}

func (p *sqlParser) Postfix() (localctx IPostfixContext) {
	localctx = NewPostfixContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 82, sqlParserRULE_postfix)
	var _alt int

	p.SetState(592)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case sqlParserAT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(580)
			p.Match(sqlParserAT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(581)
			p.Match(sqlParserTIME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(582)
			p.Match(sqlParserZONE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(583)
			p.Match(sqlParserSTRING_LITERAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserCOLLATE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(584)
			p.Match(sqlParserCOLLATE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(585)
			p.QualifiedName()
		}

	case sqlParserCOLON_COLON:
		p.EnterOuterAlt(localctx, 3)
		p.SetState(588)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(586)
					p.Match(sqlParserCOLON_COLON)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(587)
					p.TypeName()
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

			p.SetState(590)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 67, p.GetParserRuleContext())
			if p.HasError() {
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypeNameContext is an interface to support dynamic dispatch.
type ITypeNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	QualifiedName() IQualifiedNameContext

	// IsTypeNameContext differentiates from other interfaces.
	IsTypeNameContext()
}

type TypeNameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypeNameContext() *TypeNameContext {
	var p = new(TypeNameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_typeName
	return p
}

func InitEmptyTypeNameContext(p *TypeNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_typeName
}

func (*TypeNameContext) IsTypeNameContext() {}

func NewTypeNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeNameContext {
	var p = new(TypeNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_typeName

	return p
}

func (s *TypeNameContext) GetParser() antlr.Parser { return s.parser }

func (s *TypeNameContext) QualifiedName() IQualifiedNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedNameContext)
}

func (s *TypeNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypeNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypeNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterTypeName(s)
	}
}

func (s *TypeNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitTypeName(s)
	}
}

func (p *sqlParser) TypeName() (localctx ITypeNameContext) {
	localctx = NewTypeNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 84, sqlParserRULE_typeName)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(594)
		p.QualifiedName()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IPrimaryExpressionContext is an interface to support dynamic dispatch.
type IPrimaryExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	Expression() IExpressionContext
	RPAREN() antlr.TerminalNode
	SubqueryExpression() ISubqueryExpressionContext
	CaseExpression() ICaseExpressionContext
	FunctionCall() IFunctionCallContext
	ExtractFunction() IExtractFunctionContext
	ColumnName() IColumnNameContext
	Value() IValueContext
	PARAMETER() antlr.TerminalNode

	// IsPrimaryExpressionContext differentiates from other interfaces.
	IsPrimaryExpressionContext()
}

type PrimaryExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyPrimaryExpressionContext() *PrimaryExpressionContext {
	var p = new(PrimaryExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_primaryExpression
	return p
}

func InitEmptyPrimaryExpressionContext(p *PrimaryExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_primaryExpression
}

func (*PrimaryExpressionContext) IsPrimaryExpressionContext() {}

func NewPrimaryExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryExpressionContext {
	var p = new(PrimaryExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_primaryExpression

	return p
}

func (s *PrimaryExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *PrimaryExpressionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserLPAREN, 0)
}

func (s *PrimaryExpressionContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *PrimaryExpressionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserRPAREN, 0)
}

func (s *PrimaryExpressionContext) SubqueryExpression() ISubqueryExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISubqueryExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISubqueryExpressionContext)
}

func (s *PrimaryExpressionContext) CaseExpression() ICaseExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICaseExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICaseExpressionContext)
}

func (s *PrimaryExpressionContext) FunctionCall() IFunctionCallContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionCallContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionCallContext)
}

func (s *PrimaryExpressionContext) ExtractFunction() IExtractFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExtractFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExtractFunctionContext)
}

func (s *PrimaryExpressionContext) ColumnName() IColumnNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnNameContext)
}

func (s *PrimaryExpressionContext) Value() IValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IValueContext)
}

func (s *PrimaryExpressionContext) PARAMETER() antlr.TerminalNode {
	return s.GetToken(sqlParserPARAMETER, 0)
}

func (s *PrimaryExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *PrimaryExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *PrimaryExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterPrimaryExpression(s)
	}
}

func (s *PrimaryExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitPrimaryExpression(s)
	}
}

func (p *sqlParser) PrimaryExpression() (localctx IPrimaryExpressionContext) {
	localctx = NewPrimaryExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 86, sqlParserRULE_primaryExpression)
	p.SetState(607)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 69, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(596)
			p.Match(sqlParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(597)
			p.Expression()
		}
		{
			p.SetState(598)
			p.Match(sqlParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(600)
			p.SubqueryExpression()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(601)
			p.CaseExpression()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(602)
			p.FunctionCall()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(603)
			p.ExtractFunction()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(604)
			p.ColumnName()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(605)
			p.Value()
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(606)
			p.Match(sqlParserPARAMETER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case antlr.ATNInvalidAltNumber:
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ISubqueryExpressionContext is an interface to support dynamic dispatch.
type ISubqueryExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	LPAREN() antlr.TerminalNode
	SelectStatement() ISelectStatementContext
	RPAREN() antlr.TerminalNode

	// IsSubqueryExpressionContext differentiates from other interfaces.
	IsSubqueryExpressionContext()
}

type SubqueryExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySubqueryExpressionContext() *SubqueryExpressionContext {
	var p = new(SubqueryExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_subqueryExpression
	return p
}

func InitEmptySubqueryExpressionContext(p *SubqueryExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_subqueryExpression
}

func (*SubqueryExpressionContext) IsSubqueryExpressionContext() {}

func NewSubqueryExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SubqueryExpressionContext {
	var p = new(SubqueryExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_subqueryExpression

	return p
}

func (s *SubqueryExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *SubqueryExpressionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserLPAREN, 0)
}

func (s *SubqueryExpressionContext) SelectStatement() ISelectStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectStatementContext)
}

func (s *SubqueryExpressionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserRPAREN, 0)
}

func (s *SubqueryExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SubqueryExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SubqueryExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterSubqueryExpression(s)
	}
}

func (s *SubqueryExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitSubqueryExpression(s)
	}
}

func (p *sqlParser) SubqueryExpression() (localctx ISubqueryExpressionContext) {
	localctx = NewSubqueryExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 88, sqlParserRULE_subqueryExpression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(609)
		p.Match(sqlParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(610)
		p.SelectStatement()
	}
	{
		p.SetState(611)
		p.Match(sqlParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ICaseExpressionContext is an interface to support dynamic dispatch.
type ICaseExpressionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	CASE() antlr.TerminalNode
	END() antlr.TerminalNode
	AllWHEN() []antlr.TerminalNode
	WHEN(i int) antlr.TerminalNode
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
	AllTHEN() []antlr.TerminalNode
	THEN(i int) antlr.TerminalNode
	ELSE() antlr.TerminalNode

	// IsCaseExpressionContext differentiates from other interfaces.
	IsCaseExpressionContext()
}

type CaseExpressionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCaseExpressionContext() *CaseExpressionContext {
	var p = new(CaseExpressionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_caseExpression
	return p
}

func InitEmptyCaseExpressionContext(p *CaseExpressionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_caseExpression
}

func (*CaseExpressionContext) IsCaseExpressionContext() {}

func NewCaseExpressionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CaseExpressionContext {
	var p = new(CaseExpressionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_caseExpression

	return p
}

func (s *CaseExpressionContext) GetParser() antlr.Parser { return s.parser }

func (s *CaseExpressionContext) CASE() antlr.TerminalNode {
	return s.GetToken(sqlParserCASE, 0)
}

func (s *CaseExpressionContext) END() antlr.TerminalNode {
	return s.GetToken(sqlParserEND, 0)
}

func (s *CaseExpressionContext) AllWHEN() []antlr.TerminalNode {
	return s.GetTokens(sqlParserWHEN)
}

func (s *CaseExpressionContext) WHEN(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserWHEN, i)
}

func (s *CaseExpressionContext) AllExpression() []IExpressionContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExpressionContext); ok {
			len++
		}
	}

	tst := make([]IExpressionContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExpressionContext); ok {
			tst[i] = t.(IExpressionContext)
			i++
		}
	}

	return tst
}

func (s *CaseExpressionContext) Expression(i int) IExpressionContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *CaseExpressionContext) AllTHEN() []antlr.TerminalNode {
	return s.GetTokens(sqlParserTHEN)
}

func (s *CaseExpressionContext) THEN(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserTHEN, i)
}

func (s *CaseExpressionContext) ELSE() antlr.TerminalNode {
	return s.GetToken(sqlParserELSE, 0)
}

func (s *CaseExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CaseExpressionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CaseExpressionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterCaseExpression(s)
	}
}

func (s *CaseExpressionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitCaseExpression(s)
	}
}

func (p *sqlParser) CaseExpression() (localctx ICaseExpressionContext) {
	localctx = NewCaseExpressionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 90, sqlParserRULE_caseExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(613)
		p.Match(sqlParserCASE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(619)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == sqlParserWHEN {
		{
			p.SetState(614)
			p.Match(sqlParserWHEN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(615)
			p.Expression()
		}
		{
			p.SetState(616)
			p.Match(sqlParserTHEN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(617)
			p.Expression()
		}

		p.SetState(621)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(625)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserELSE {
		{
			p.SetState(623)
			p.Match(sqlParserELSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(624)
			p.Expression()
		}

	}
	{
		p.SetState(627)
		p.Match(sqlParserEND)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunctionCallContext is an interface to support dynamic dispatch.
type IFunctionCallContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	QualifiedName() IQualifiedNameContext
	LPAREN() antlr.TerminalNode
	RPAREN() antlr.TerminalNode
	AllFunctionArg() []IFunctionArgContext
	FunctionArg(i int) IFunctionArgContext
	AllCOMMA() []antlr.TerminalNode
	COMMA(i int) antlr.TerminalNode

	// IsFunctionCallContext differentiates from other interfaces.
	IsFunctionCallContext()
}

type FunctionCallContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionCallContext() *FunctionCallContext {
	var p = new(FunctionCallContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_functionCall
	return p
}

func InitEmptyFunctionCallContext(p *FunctionCallContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_functionCall
}

func (*FunctionCallContext) IsFunctionCallContext() {}

func NewFunctionCallContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionCallContext {
	var p = new(FunctionCallContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_functionCall

	return p
}

func (s *FunctionCallContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionCallContext) QualifiedName() IQualifiedNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedNameContext)
}

func (s *FunctionCallContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserLPAREN, 0)
}

func (s *FunctionCallContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserRPAREN, 0)
}

func (s *FunctionCallContext) AllFunctionArg() []IFunctionArgContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IFunctionArgContext); ok {
			len++
		}
	}

	tst := make([]IFunctionArgContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IFunctionArgContext); ok {
			tst[i] = t.(IFunctionArgContext)
			i++
		}
	}

	return tst
}

func (s *FunctionCallContext) FunctionArg(i int) IFunctionArgContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionArgContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionArgContext)
}

func (s *FunctionCallContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(sqlParserCOMMA)
}

func (s *FunctionCallContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserCOMMA, i)
}

func (s *FunctionCallContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionCallContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionCallContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterFunctionCall(s)
	}
}

func (s *FunctionCallContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitFunctionCall(s)
	}
}

func (p *sqlParser) FunctionCall() (localctx IFunctionCallContext) {
	localctx = NewFunctionCallContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 92, sqlParserRULE_functionCall)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(629)
		p.QualifiedName()
	}
	{
		p.SetState(630)
		p.Match(sqlParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(639)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&-99008823057784832) != 0) || ((int64((_la-64)) & ^0x3f) == 0 && ((int64(1)<<(_la-64))&4123176007687) != 0) {
		{
			p.SetState(631)
			p.FunctionArg()
		}
		p.SetState(636)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == sqlParserCOMMA {
			{
				p.SetState(632)
				p.Match(sqlParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(633)
				p.FunctionArg()
			}

			p.SetState(638)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	}
	{
		p.SetState(641)
		p.Match(sqlParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IFunctionArgContext is an interface to support dynamic dispatch.
type IFunctionArgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	Expression() IExpressionContext
	STAR() antlr.TerminalNode

	// IsFunctionArgContext differentiates from other interfaces.
	IsFunctionArgContext()
}

type FunctionArgContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionArgContext() *FunctionArgContext {
	var p = new(FunctionArgContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_functionArg
	return p
}

func InitEmptyFunctionArgContext(p *FunctionArgContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_functionArg
}

func (*FunctionArgContext) IsFunctionArgContext() {}

func NewFunctionArgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionArgContext {
	var p = new(FunctionArgContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_functionArg

	return p
}

func (s *FunctionArgContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionArgContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *FunctionArgContext) STAR() antlr.TerminalNode {
	return s.GetToken(sqlParserSTAR, 0)
}

func (s *FunctionArgContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionArgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionArgContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterFunctionArg(s)
	}
}

func (s *FunctionArgContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitFunctionArg(s)
	}
}

func (p *sqlParser) FunctionArg() (localctx IFunctionArgContext) {
	localctx = NewFunctionArgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 94, sqlParserRULE_functionArg)
	p.SetState(645)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case sqlParserNOT, sqlParserEXISTS, sqlParserNULL, sqlParserDEFAULT, sqlParserCASE, sqlParserTIME, sqlParserEXTRACT, sqlParserINT_TYPE, sqlParserFLOAT_TYPE, sqlParserBOOL_TYPE, sqlParserTEXT_TYPE, sqlParserVARCHAR_TYPE, sqlParserCHAR_TYPE, sqlParserSERIAL_TYPE, sqlParserTIMESTAMP_TYPE, sqlParserDATE_TYPE, sqlParserINTERVAL_TYPE, sqlParserCURRENT_TIMESTAMP, sqlParserCURRENT_USER, sqlParserTRUE, sqlParserFALSE, sqlParserLPAREN, sqlParserPLUS, sqlParserMINUS, sqlParserIDENTIFIER, sqlParserPARAMETER, sqlParserSTRING_LITERAL, sqlParserNUMBER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(643)
			p.Expression()
		}

	case sqlParserSTAR:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(644)
			p.Match(sqlParserSTAR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IExtractFunctionContext is an interface to support dynamic dispatch.
type IExtractFunctionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EXTRACT() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	IDENTIFIER() antlr.TerminalNode
	FROM() antlr.TerminalNode
	Expression() IExpressionContext
	RPAREN() antlr.TerminalNode

	// IsExtractFunctionContext differentiates from other interfaces.
	IsExtractFunctionContext()
}

type ExtractFunctionContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExtractFunctionContext() *ExtractFunctionContext {
	var p = new(ExtractFunctionContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_extractFunction
	return p
}

func InitEmptyExtractFunctionContext(p *ExtractFunctionContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_extractFunction
}

func (*ExtractFunctionContext) IsExtractFunctionContext() {}

func NewExtractFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExtractFunctionContext {
	var p = new(ExtractFunctionContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_extractFunction

	return p
}

func (s *ExtractFunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *ExtractFunctionContext) EXTRACT() antlr.TerminalNode {
	return s.GetToken(sqlParserEXTRACT, 0)
}

func (s *ExtractFunctionContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserLPAREN, 0)
}

func (s *ExtractFunctionContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(sqlParserIDENTIFIER, 0)
}

func (s *ExtractFunctionContext) FROM() antlr.TerminalNode {
	return s.GetToken(sqlParserFROM, 0)
}

func (s *ExtractFunctionContext) Expression() IExpressionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExpressionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExpressionContext)
}

func (s *ExtractFunctionContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserRPAREN, 0)
}

func (s *ExtractFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExtractFunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExtractFunctionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterExtractFunction(s)
	}
}

func (s *ExtractFunctionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitExtractFunction(s)
	}
}

func (p *sqlParser) ExtractFunction() (localctx IExtractFunctionContext) {
	localctx = NewExtractFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 96, sqlParserRULE_extractFunction)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(647)
		p.Match(sqlParserEXTRACT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(648)
		p.Match(sqlParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(649)
		p.Match(sqlParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(650)
		p.Match(sqlParserFROM)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(651)
		p.Expression()
	}
	{
		p.SetState(652)
		p.Match(sqlParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// INamePartContext is an interface to support dynamic dispatch.
type INamePartContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	IDENTIFIER() antlr.TerminalNode
	DEFAULT() antlr.TerminalNode
	TEXT_TYPE() antlr.TerminalNode
	INT_TYPE() antlr.TerminalNode
	FLOAT_TYPE() antlr.TerminalNode
	BOOL_TYPE() antlr.TerminalNode
	VARCHAR_TYPE() antlr.TerminalNode
	CHAR_TYPE() antlr.TerminalNode
	SERIAL_TYPE() antlr.TerminalNode
	TIMESTAMP_TYPE() antlr.TerminalNode
	DATE_TYPE() antlr.TerminalNode
	TIME() antlr.TerminalNode
	INTERVAL_TYPE() antlr.TerminalNode

	// IsNamePartContext differentiates from other interfaces.
	IsNamePartContext()
}

type NamePartContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyNamePartContext() *NamePartContext {
	var p = new(NamePartContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_namePart
	return p
}

func InitEmptyNamePartContext(p *NamePartContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_namePart
}

func (*NamePartContext) IsNamePartContext() {}

func NewNamePartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *NamePartContext {
	var p = new(NamePartContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_namePart

	return p
}

func (s *NamePartContext) GetParser() antlr.Parser { return s.parser }

func (s *NamePartContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(sqlParserIDENTIFIER, 0)
}

func (s *NamePartContext) DEFAULT() antlr.TerminalNode {
	return s.GetToken(sqlParserDEFAULT, 0)
}

func (s *NamePartContext) TEXT_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserTEXT_TYPE, 0)
}

func (s *NamePartContext) INT_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserINT_TYPE, 0)
}

func (s *NamePartContext) FLOAT_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserFLOAT_TYPE, 0)
}

func (s *NamePartContext) BOOL_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserBOOL_TYPE, 0)
}

func (s *NamePartContext) VARCHAR_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserVARCHAR_TYPE, 0)
}

func (s *NamePartContext) CHAR_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserCHAR_TYPE, 0)
}

func (s *NamePartContext) SERIAL_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserSERIAL_TYPE, 0)
}

func (s *NamePartContext) TIMESTAMP_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserTIMESTAMP_TYPE, 0)
}

func (s *NamePartContext) DATE_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserDATE_TYPE, 0)
}

func (s *NamePartContext) TIME() antlr.TerminalNode {
	return s.GetToken(sqlParserTIME, 0)
}

func (s *NamePartContext) INTERVAL_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserINTERVAL_TYPE, 0)
}

func (s *NamePartContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NamePartContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *NamePartContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterNamePart(s)
	}
}

func (s *NamePartContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitNamePart(s)
	}
}

func (p *sqlParser) NamePart() (localctx INamePartContext) {
	localctx = NewNamePartContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 98, sqlParserRULE_namePart)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(654)
		_la = p.GetTokenStream().LA(1)

		if !(((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&-135107988820983808) != 0) || ((int64((_la-64)) & ^0x3f) == 0 && ((int64(1)<<(_la-64))&274877906951) != 0)) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IQualifiedNameContext is an interface to support dynamic dispatch.
type IQualifiedNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	AllNamePart() []INamePartContext
	NamePart(i int) INamePartContext
	AllDOT() []antlr.TerminalNode
	DOT(i int) antlr.TerminalNode
	AllTILDE() []antlr.TerminalNode
	TILDE(i int) antlr.TerminalNode

	// IsQualifiedNameContext differentiates from other interfaces.
	IsQualifiedNameContext()
}

type QualifiedNameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyQualifiedNameContext() *QualifiedNameContext {
	var p = new(QualifiedNameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_qualifiedName
	return p
}

func InitEmptyQualifiedNameContext(p *QualifiedNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_qualifiedName
}

func (*QualifiedNameContext) IsQualifiedNameContext() {}

func NewQualifiedNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QualifiedNameContext {
	var p = new(QualifiedNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_qualifiedName

	return p
}

func (s *QualifiedNameContext) GetParser() antlr.Parser { return s.parser }

func (s *QualifiedNameContext) AllNamePart() []INamePartContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(INamePartContext); ok {
			len++
		}
	}

	tst := make([]INamePartContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(INamePartContext); ok {
			tst[i] = t.(INamePartContext)
			i++
		}
	}

	return tst
}

func (s *QualifiedNameContext) NamePart(i int) INamePartContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(INamePartContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(INamePartContext)
}

func (s *QualifiedNameContext) AllDOT() []antlr.TerminalNode {
	return s.GetTokens(sqlParserDOT)
}

func (s *QualifiedNameContext) DOT(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserDOT, i)
}

func (s *QualifiedNameContext) AllTILDE() []antlr.TerminalNode {
	return s.GetTokens(sqlParserTILDE)
}

func (s *QualifiedNameContext) TILDE(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserTILDE, i)
}

func (s *QualifiedNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *QualifiedNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *QualifiedNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterQualifiedName(s)
	}
}

func (s *QualifiedNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitQualifiedName(s)
	}
}

func (p *sqlParser) QualifiedName() (localctx IQualifiedNameContext) {
	localctx = NewQualifiedNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 100, sqlParserRULE_qualifiedName)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(656)
		p.NamePart()
	}
	p.SetState(664)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserDOT {
		{
			p.SetState(657)
			p.Match(sqlParserDOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(660)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case sqlParserDEFAULT, sqlParserTIME, sqlParserINT_TYPE, sqlParserFLOAT_TYPE, sqlParserBOOL_TYPE, sqlParserTEXT_TYPE, sqlParserVARCHAR_TYPE, sqlParserCHAR_TYPE, sqlParserSERIAL_TYPE, sqlParserTIMESTAMP_TYPE, sqlParserDATE_TYPE, sqlParserINTERVAL_TYPE, sqlParserIDENTIFIER:
			{
				p.SetState(658)
				p.NamePart()
			}

		case sqlParserTILDE:
			{
				p.SetState(659)
				p.Match(sqlParserTILDE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(666)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IColumnNameContext is an interface to support dynamic dispatch.
type IColumnNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	QualifiedName() IQualifiedNameContext

	// IsColumnNameContext differentiates from other interfaces.
	IsColumnNameContext()
}

type ColumnNameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnNameContext() *ColumnNameContext {
	var p = new(ColumnNameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_columnName
	return p
}

func InitEmptyColumnNameContext(p *ColumnNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_columnName
}

func (*ColumnNameContext) IsColumnNameContext() {}

func NewColumnNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnNameContext {
	var p = new(ColumnNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_columnName

	return p
}

func (s *ColumnNameContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnNameContext) QualifiedName() IQualifiedNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedNameContext)
}

func (s *ColumnNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ColumnNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterColumnName(s)
	}
}

func (s *ColumnNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitColumnName(s)
	}
}

func (p *sqlParser) ColumnName() (localctx IColumnNameContext) {
	localctx = NewColumnNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 102, sqlParserRULE_columnName)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(667)
		p.QualifiedName()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITableNameContext is an interface to support dynamic dispatch.
type ITableNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	QualifiedName() IQualifiedNameContext

	// IsTableNameContext differentiates from other interfaces.
	IsTableNameContext()
}

type TableNameContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTableNameContext() *TableNameContext {
	var p = new(TableNameContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_tableName
	return p
}

func InitEmptyTableNameContext(p *TableNameContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_tableName
}

func (*TableNameContext) IsTableNameContext() {}

func NewTableNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TableNameContext {
	var p = new(TableNameContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_tableName

	return p
}

func (s *TableNameContext) GetParser() antlr.Parser { return s.parser }

func (s *TableNameContext) QualifiedName() IQualifiedNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedNameContext)
}

func (s *TableNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TableNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TableNameContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterTableName(s)
	}
}

func (s *TableNameContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitTableName(s)
	}
}

func (p *sqlParser) TableName() (localctx ITableNameContext) {
	localctx = NewTableNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 104, sqlParserRULE_tableName)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(669)
		p.QualifiedName()
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOperatorContext is an interface to support dynamic dispatch.
type IOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	EQ() antlr.TerminalNode
	GT() antlr.TerminalNode
	LT() antlr.TerminalNode
	GE() antlr.TerminalNode
	LE() antlr.TerminalNode
	NE() antlr.TerminalNode
	TILDE() antlr.TerminalNode
	NREGEX() antlr.TerminalNode
	IREGEX() antlr.TerminalNode
	NIREGEX() antlr.TerminalNode

	// IsOperatorContext differentiates from other interfaces.
	IsOperatorContext()
}

type OperatorContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOperatorContext() *OperatorContext {
	var p = new(OperatorContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_operator
	return p
}

func InitEmptyOperatorContext(p *OperatorContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_operator
}

func (*OperatorContext) IsOperatorContext() {}

func NewOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OperatorContext {
	var p = new(OperatorContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_operator

	return p
}

func (s *OperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *OperatorContext) EQ() antlr.TerminalNode {
	return s.GetToken(sqlParserEQ, 0)
}

func (s *OperatorContext) GT() antlr.TerminalNode {
	return s.GetToken(sqlParserGT, 0)
}

func (s *OperatorContext) LT() antlr.TerminalNode {
	return s.GetToken(sqlParserLT, 0)
}

func (s *OperatorContext) GE() antlr.TerminalNode {
	return s.GetToken(sqlParserGE, 0)
}

func (s *OperatorContext) LE() antlr.TerminalNode {
	return s.GetToken(sqlParserLE, 0)
}

func (s *OperatorContext) NE() antlr.TerminalNode {
	return s.GetToken(sqlParserNE, 0)
}

func (s *OperatorContext) TILDE() antlr.TerminalNode {
	return s.GetToken(sqlParserTILDE, 0)
}

func (s *OperatorContext) NREGEX() antlr.TerminalNode {
	return s.GetToken(sqlParserNREGEX, 0)
}

func (s *OperatorContext) IREGEX() antlr.TerminalNode {
	return s.GetToken(sqlParserIREGEX, 0)
}

func (s *OperatorContext) NIREGEX() antlr.TerminalNode {
	return s.GetToken(sqlParserNIREGEX, 0)
}

func (s *OperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterOperator(s)
	}
}

func (s *OperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitOperator(s)
	}
}

func (p *sqlParser) Operator() (localctx IOperatorContext) {
	localctx = NewOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 106, sqlParserRULE_operator)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(671)
		_la = p.GetTokenStream().LA(1)

		if !((int64((_la-90)) & ^0x3f) == 0 && ((int64(1)<<(_la-90))&1023) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IOperatorExprContext is an interface to support dynamic dispatch.
type IOperatorExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	OPERATOR_KW() antlr.TerminalNode
	LPAREN() antlr.TerminalNode
	QualifiedName() IQualifiedNameContext
	RPAREN() antlr.TerminalNode

	// IsOperatorExprContext differentiates from other interfaces.
	IsOperatorExprContext()
}

type OperatorExprContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyOperatorExprContext() *OperatorExprContext {
	var p = new(OperatorExprContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_operatorExpr
	return p
}

func InitEmptyOperatorExprContext(p *OperatorExprContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_operatorExpr
}

func (*OperatorExprContext) IsOperatorExprContext() {}

func NewOperatorExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OperatorExprContext {
	var p = new(OperatorExprContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_operatorExpr

	return p
}

func (s *OperatorExprContext) GetParser() antlr.Parser { return s.parser }

func (s *OperatorExprContext) OPERATOR_KW() antlr.TerminalNode {
	return s.GetToken(sqlParserOPERATOR_KW, 0)
}

func (s *OperatorExprContext) LPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserLPAREN, 0)
}

func (s *OperatorExprContext) QualifiedName() IQualifiedNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IQualifiedNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IQualifiedNameContext)
}

func (s *OperatorExprContext) RPAREN() antlr.TerminalNode {
	return s.GetToken(sqlParserRPAREN, 0)
}

func (s *OperatorExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *OperatorExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *OperatorExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterOperatorExpr(s)
	}
}

func (s *OperatorExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitOperatorExpr(s)
	}
}

func (p *sqlParser) OperatorExpr() (localctx IOperatorExprContext) {
	localctx = NewOperatorExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 108, sqlParserRULE_operatorExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(673)
		p.Match(sqlParserOPERATOR_KW)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(674)
		p.Match(sqlParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(675)
		p.QualifiedName()
	}
	{
		p.SetState(676)
		p.Match(sqlParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// IValueContext is an interface to support dynamic dispatch.
type IValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	STRING_LITERAL() antlr.TerminalNode
	NUMBER() antlr.TerminalNode
	CURRENT_TIMESTAMP() antlr.TerminalNode
	CURRENT_USER() antlr.TerminalNode
	TRUE() antlr.TerminalNode
	FALSE() antlr.TerminalNode
	NULL() antlr.TerminalNode
	TypedLiteral() ITypedLiteralContext

	// IsValueContext differentiates from other interfaces.
	IsValueContext()
}

type ValueContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyValueContext() *ValueContext {
	var p = new(ValueContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_value
	return p
}

func InitEmptyValueContext(p *ValueContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_value
}

func (*ValueContext) IsValueContext() {}

func NewValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ValueContext {
	var p = new(ValueContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_value

	return p
}

func (s *ValueContext) GetParser() antlr.Parser { return s.parser }

func (s *ValueContext) STRING_LITERAL() antlr.TerminalNode {
	return s.GetToken(sqlParserSTRING_LITERAL, 0)
}

func (s *ValueContext) NUMBER() antlr.TerminalNode {
	return s.GetToken(sqlParserNUMBER, 0)
}

func (s *ValueContext) CURRENT_TIMESTAMP() antlr.TerminalNode {
	return s.GetToken(sqlParserCURRENT_TIMESTAMP, 0)
}

func (s *ValueContext) CURRENT_USER() antlr.TerminalNode {
	return s.GetToken(sqlParserCURRENT_USER, 0)
}

func (s *ValueContext) TRUE() antlr.TerminalNode {
	return s.GetToken(sqlParserTRUE, 0)
}

func (s *ValueContext) FALSE() antlr.TerminalNode {
	return s.GetToken(sqlParserFALSE, 0)
}

func (s *ValueContext) NULL() antlr.TerminalNode {
	return s.GetToken(sqlParserNULL, 0)
}

func (s *ValueContext) TypedLiteral() ITypedLiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITypedLiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITypedLiteralContext)
}

func (s *ValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterValue(s)
	}
}

func (s *ValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitValue(s)
	}
}

func (p *sqlParser) Value() (localctx IValueContext) {
	localctx = NewValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 110, sqlParserRULE_value)
	p.SetState(686)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case sqlParserSTRING_LITERAL:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(678)
			p.Match(sqlParserSTRING_LITERAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserNUMBER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(679)
			p.Match(sqlParserNUMBER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserCURRENT_TIMESTAMP:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(680)
			p.Match(sqlParserCURRENT_TIMESTAMP)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserCURRENT_USER:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(681)
			p.Match(sqlParserCURRENT_USER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserTRUE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(682)
			p.Match(sqlParserTRUE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserFALSE:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(683)
			p.Match(sqlParserFALSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserNULL:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(684)
			p.Match(sqlParserNULL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserTIME, sqlParserTIMESTAMP_TYPE, sqlParserDATE_TYPE, sqlParserINTERVAL_TYPE:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(685)
			p.TypedLiteral()
		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}

// ITypedLiteralContext is an interface to support dynamic dispatch.
type ITypedLiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// Getter signatures
	DATE_TYPE() antlr.TerminalNode
	STRING_LITERAL() antlr.TerminalNode
	TIMESTAMP_TYPE() antlr.TerminalNode
	TIME() antlr.TerminalNode
	INTERVAL_TYPE() antlr.TerminalNode
	IntervalFields() IIntervalFieldsContext

	// IsTypedLiteralContext differentiates from other interfaces.
	IsTypedLiteralContext()
}

type TypedLiteralContext struct {
	antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTypedLiteralContext() *TypedLiteralContext {
	var p = new(TypedLiteralContext)
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_typedLiteral
	return p
}

func InitEmptyTypedLiteralContext(p *TypedLiteralContext) {
	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, nil, -1)
	p.RuleIndex = sqlParserRULE_typedLiteral
}

func (*TypedLiteralContext) IsTypedLiteralContext() {}

func NewTypedLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypedLiteralContext {
	var p = new(TypedLiteralContext)

	antlr.InitBaseParserRuleContext(&p.BaseParserRuleContext, parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_typedLiteral

	return p
}

func (s *TypedLiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *TypedLiteralContext) DATE_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserDATE_TYPE, 0)
}

func (s *TypedLiteralContext) STRING_LITERAL() antlr.TerminalNode {
	return s.GetToken(sqlParserSTRING_LITERAL, 0)
}

func (s *TypedLiteralContext) TIMESTAMP_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserTIMESTAMP_TYPE, 0)
}

func (s *TypedLiteralContext) TIME() antlr.TerminalNode {
	return s.GetToken(sqlParserTIME, 0)
}

func (s *TypedLiteralContext) INTERVAL_TYPE() antlr.TerminalNode {
	return s.GetToken(sqlParserINTERVAL_TYPE, 0)
}

func (s *TypedLiteralContext) IntervalFields() IIntervalFieldsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIntervalFieldsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIntervalFieldsContext)
}

func (s *TypedLiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TypedLiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TypedLiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterTypedLiteral(s)
	}
}

func (s *TypedLiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitTypedLiteral(s)
	}
}

func (p *sqlParser) TypedLiteral() (localctx ITypedLiteralContext) {
	localctx = NewTypedLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 112, sqlParserRULE_typedLiteral)
	var _la int

	p.SetState(699)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case sqlParserDATE_TYPE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(688)
			p.Match(sqlParserDATE_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(689)
			p.Match(sqlParserSTRING_LITERAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserTIMESTAMP_TYPE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(690)
			p.Match(sqlParserTIMESTAMP_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(691)
			p.Match(sqlParserSTRING_LITERAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserTIME:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(692)
			p.Match(sqlParserTIME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(693)
			p.Match(sqlParserSTRING_LITERAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserINTERVAL_TYPE:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(694)
			p.Match(sqlParserINTERVAL_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(695)
			p.Match(sqlParserSTRING_LITERAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(697)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if (int64((_la-67)) & ^0x3f) == 0 && ((int64(1)<<(_la-67))&63) != 0 {
			{
				p.SetState(696)
				p.IntervalFields()
			}

		}

	default:
		p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		goto errorExit
	}

errorExit:
	if p.HasError() {
		v := p.GetError()
		localctx.SetException(v)
		p.GetErrorHandler().ReportError(p, v)
		p.GetErrorHandler().Recover(p, v)
		p.SetError(nil)
	}
	p.ExitRule()
	return localctx
	goto errorExit // Trick to prevent compiler error if the label is not used
}
