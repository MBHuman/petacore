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
		"", "", "", "", "", "", "", "", "", "", "", "", "'('", "')'", "';'",
		"','", "'.'", "'*'", "'+'", "'-'", "'/'", "'||'", "'::'", "'~'", "'!~'",
		"'~*'", "'!~*'", "'='", "'>'", "'<'", "'>='", "'<='",
	}
	staticData.SymbolicNames = []string{
		"", "CREATE", "TABLE", "INSERT", "INTO", "VALUES", "PRIMARY", "KEY",
		"DROP", "TRUNCATE", "SET", "TO", "IF", "NOT", "EXISTS", "NULL", "UNIQUE",
		"DEFAULT", "SHOW", "SELECT", "FROM", "WHERE", "LIMIT", "OFFSET", "ORDER",
		"BY", "ASC", "DESC", "JOIN", "INNER", "LEFT", "RIGHT", "FULL", "CROSS",
		"ON", "IN", "AND", "OR", "IS", "LIKE", "ILIKE", "CASE", "WHEN", "THEN",
		"ELSE", "END", "AS", "AT", "TIME", "ZONE", "EXTRACT", "STRING_TYPE",
		"INT_TYPE", "FLOAT_TYPE", "BOOL_TYPE", "TEXT_TYPE", "VARCHAR_TYPE",
		"CHAR_TYPE", "SERIAL_TYPE", "TIMESTAMP_TYPE", "CURRENT_TIMESTAMP", "TRUE",
		"FALSE", "LPAREN", "RPAREN", "SEMICOLON", "COMMA", "DOT", "STAR", "PLUS",
		"MINUS", "SLASH", "CONCAT", "COLON_COLON", "TILDE", "NREGEX", "IREGEX",
		"NIREGEX", "EQ", "GT", "LT", "GE", "LE", "NE", "OPERATOR_KW", "COLLATE",
		"IDENTIFIER", "PARAMETER", "STRING_LITERAL", "NUMBER", "BLOCK_COMMENT",
		"LINE_COMMENT", "WS",
	}
	staticData.RuleNames = []string{
		"query", "statement", "createTableStatement", "columnDefinition", "columnConstraints",
		"dataType", "insertStatement", "columnList", "valueList", "dropTableStatement",
		"truncateTableStatement", "setStatement", "showStatement", "selectStatement",
		"selectList", "selectItem", "selectAll", "fromClause", "tableFactor",
		"joinClause", "whereClause", "orderByClause", "orderByItem", "limitValue",
		"offsetValue", "alias", "expression", "orExpression", "andExpression",
		"notExpression", "comparisonExpression", "concatExpression", "additiveExpression",
		"multiplicativeExpression", "unaryExpression", "castExpression", "postfix",
		"typeName", "primaryExpression", "subqueryExpression", "caseExpression",
		"functionCall", "extractFunction", "namePart", "qualifiedName", "columnName",
		"tableName", "operator", "operatorExpr", "value",
	}
	staticData.PredictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 92, 532, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 2, 18, 7, 18, 2, 19, 7, 19, 2, 20, 7, 20, 2,
		21, 7, 21, 2, 22, 7, 22, 2, 23, 7, 23, 2, 24, 7, 24, 2, 25, 7, 25, 2, 26,
		7, 26, 2, 27, 7, 27, 2, 28, 7, 28, 2, 29, 7, 29, 2, 30, 7, 30, 2, 31, 7,
		31, 2, 32, 7, 32, 2, 33, 7, 33, 2, 34, 7, 34, 2, 35, 7, 35, 2, 36, 7, 36,
		2, 37, 7, 37, 2, 38, 7, 38, 2, 39, 7, 39, 2, 40, 7, 40, 2, 41, 7, 41, 2,
		42, 7, 42, 2, 43, 7, 43, 2, 44, 7, 44, 2, 45, 7, 45, 2, 46, 7, 46, 2, 47,
		7, 47, 2, 48, 7, 48, 2, 49, 7, 49, 1, 0, 1, 0, 3, 0, 103, 8, 0, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 112, 8, 1, 1, 2, 1, 2, 1, 2, 1,
		2, 1, 2, 3, 2, 119, 8, 2, 1, 2, 1, 2, 1, 2, 1, 2, 1, 2, 5, 2, 126, 8, 2,
		10, 2, 12, 2, 129, 9, 2, 1, 2, 3, 2, 132, 8, 2, 1, 2, 1, 2, 1, 2, 1, 2,
		1, 2, 1, 2, 5, 2, 140, 8, 2, 10, 2, 12, 2, 143, 9, 2, 1, 2, 1, 2, 3, 2,
		147, 8, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3,
		1, 3, 1, 3, 3, 3, 161, 8, 3, 1, 4, 1, 4, 1, 4, 1, 4, 1, 4, 5, 4, 168, 8,
		4, 10, 4, 12, 4, 171, 9, 4, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1, 5, 1,
		5, 1, 5, 3, 5, 182, 8, 5, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 188, 8, 5, 1, 5,
		1, 5, 3, 5, 192, 8, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6,
		1, 6, 1, 6, 5, 6, 204, 8, 6, 10, 6, 12, 6, 207, 9, 6, 1, 7, 1, 7, 1, 7,
		5, 7, 212, 8, 7, 10, 7, 12, 7, 215, 9, 7, 1, 8, 1, 8, 1, 8, 1, 8, 5, 8,
		221, 8, 8, 10, 8, 12, 8, 224, 9, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9,
		1, 10, 1, 10, 3, 10, 234, 8, 10, 1, 10, 1, 10, 1, 11, 1, 11, 1, 11, 1,
		11, 1, 11, 1, 12, 1, 12, 1, 12, 5, 12, 246, 8, 12, 10, 12, 12, 12, 249,
		9, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 3, 13, 257, 8, 13, 1,
		13, 1, 13, 1, 13, 3, 13, 262, 8, 13, 1, 13, 1, 13, 3, 13, 266, 8, 13, 1,
		13, 1, 13, 3, 13, 270, 8, 13, 3, 13, 272, 8, 13, 1, 14, 1, 14, 1, 14, 5,
		14, 277, 8, 14, 10, 14, 12, 14, 280, 9, 14, 1, 15, 1, 15, 3, 15, 284, 8,
		15, 1, 15, 3, 15, 287, 8, 15, 1, 16, 1, 16, 1, 17, 1, 17, 5, 17, 293, 8,
		17, 10, 17, 12, 17, 296, 9, 17, 1, 18, 1, 18, 3, 18, 300, 8, 18, 1, 18,
		1, 18, 1, 18, 1, 18, 1, 18, 3, 18, 307, 8, 18, 1, 19, 3, 19, 310, 8, 19,
		1, 19, 1, 19, 1, 19, 3, 19, 315, 8, 19, 1, 19, 1, 19, 1, 19, 1, 20, 1,
		20, 1, 21, 1, 21, 1, 21, 5, 21, 325, 8, 21, 10, 21, 12, 21, 328, 9, 21,
		1, 22, 1, 22, 3, 22, 332, 8, 22, 1, 23, 1, 23, 1, 24, 1, 24, 1, 25, 3,
		25, 339, 8, 25, 1, 25, 1, 25, 1, 26, 1, 26, 1, 27, 1, 27, 1, 27, 5, 27,
		348, 8, 27, 10, 27, 12, 27, 351, 9, 27, 1, 28, 1, 28, 1, 28, 5, 28, 356,
		8, 28, 10, 28, 12, 28, 359, 9, 28, 1, 29, 3, 29, 362, 8, 29, 1, 29, 1,
		29, 1, 30, 1, 30, 1, 30, 3, 30, 369, 8, 30, 1, 30, 1, 30, 1, 30, 3, 30,
		374, 8, 30, 1, 30, 1, 30, 1, 30, 1, 30, 1, 30, 5, 30, 381, 8, 30, 10, 30,
		12, 30, 384, 9, 30, 1, 30, 1, 30, 1, 30, 3, 30, 389, 8, 30, 1, 30, 1, 30,
		1, 30, 1, 30, 3, 30, 395, 8, 30, 1, 30, 3, 30, 398, 8, 30, 1, 31, 1, 31,
		1, 31, 5, 31, 403, 8, 31, 10, 31, 12, 31, 406, 9, 31, 1, 32, 1, 32, 1,
		32, 5, 32, 411, 8, 32, 10, 32, 12, 32, 414, 9, 32, 1, 33, 1, 33, 1, 33,
		5, 33, 419, 8, 33, 10, 33, 12, 33, 422, 9, 33, 1, 34, 3, 34, 425, 8, 34,
		1, 34, 1, 34, 1, 35, 1, 35, 5, 35, 431, 8, 35, 10, 35, 12, 35, 434, 9,
		35, 1, 36, 1, 36, 1, 36, 1, 36, 1, 36, 1, 36, 1, 36, 1, 36, 4, 36, 444,
		8, 36, 11, 36, 12, 36, 445, 3, 36, 448, 8, 36, 1, 37, 1, 37, 1, 38, 1,
		38, 1, 38, 1, 38, 1, 38, 1, 38, 1, 38, 1, 38, 1, 38, 1, 38, 1, 38, 3, 38,
		463, 8, 38, 1, 39, 1, 39, 1, 39, 1, 39, 1, 40, 1, 40, 1, 40, 1, 40, 1,
		40, 1, 40, 4, 40, 475, 8, 40, 11, 40, 12, 40, 476, 1, 40, 1, 40, 3, 40,
		481, 8, 40, 1, 40, 1, 40, 1, 41, 1, 41, 1, 41, 1, 41, 1, 41, 5, 41, 490,
		8, 41, 10, 41, 12, 41, 493, 9, 41, 3, 41, 495, 8, 41, 1, 41, 1, 41, 1,
		42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 42, 1, 43, 1, 43, 1, 44, 1, 44,
		1, 44, 1, 44, 3, 44, 512, 8, 44, 5, 44, 514, 8, 44, 10, 44, 12, 44, 517,
		9, 44, 1, 45, 1, 45, 1, 46, 1, 46, 1, 47, 1, 47, 1, 48, 1, 48, 1, 48, 1,
		48, 1, 48, 1, 49, 1, 49, 1, 49, 0, 0, 50, 0, 2, 4, 6, 8, 10, 12, 14, 16,
		18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52,
		54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88,
		90, 92, 94, 96, 98, 0, 8, 2, 0, 11, 11, 78, 78, 1, 0, 29, 33, 1, 0, 26,
		27, 1, 0, 69, 70, 2, 0, 68, 68, 71, 71, 3, 0, 17, 17, 52, 59, 86, 86, 1,
		0, 74, 83, 3, 0, 15, 15, 60, 62, 88, 89, 561, 0, 100, 1, 0, 0, 0, 2, 111,
		1, 0, 0, 0, 4, 113, 1, 0, 0, 0, 6, 160, 1, 0, 0, 0, 8, 169, 1, 0, 0, 0,
		10, 191, 1, 0, 0, 0, 12, 193, 1, 0, 0, 0, 14, 208, 1, 0, 0, 0, 16, 216,
		1, 0, 0, 0, 18, 227, 1, 0, 0, 0, 20, 231, 1, 0, 0, 0, 22, 237, 1, 0, 0,
		0, 24, 242, 1, 0, 0, 0, 26, 250, 1, 0, 0, 0, 28, 273, 1, 0, 0, 0, 30, 286,
		1, 0, 0, 0, 32, 288, 1, 0, 0, 0, 34, 290, 1, 0, 0, 0, 36, 306, 1, 0, 0,
		0, 38, 309, 1, 0, 0, 0, 40, 319, 1, 0, 0, 0, 42, 321, 1, 0, 0, 0, 44, 329,
		1, 0, 0, 0, 46, 333, 1, 0, 0, 0, 48, 335, 1, 0, 0, 0, 50, 338, 1, 0, 0,
		0, 52, 342, 1, 0, 0, 0, 54, 344, 1, 0, 0, 0, 56, 352, 1, 0, 0, 0, 58, 361,
		1, 0, 0, 0, 60, 365, 1, 0, 0, 0, 62, 399, 1, 0, 0, 0, 64, 407, 1, 0, 0,
		0, 66, 415, 1, 0, 0, 0, 68, 424, 1, 0, 0, 0, 70, 428, 1, 0, 0, 0, 72, 447,
		1, 0, 0, 0, 74, 449, 1, 0, 0, 0, 76, 462, 1, 0, 0, 0, 78, 464, 1, 0, 0,
		0, 80, 468, 1, 0, 0, 0, 82, 484, 1, 0, 0, 0, 84, 498, 1, 0, 0, 0, 86, 505,
		1, 0, 0, 0, 88, 507, 1, 0, 0, 0, 90, 518, 1, 0, 0, 0, 92, 520, 1, 0, 0,
		0, 94, 522, 1, 0, 0, 0, 96, 524, 1, 0, 0, 0, 98, 529, 1, 0, 0, 0, 100,
		102, 3, 2, 1, 0, 101, 103, 5, 65, 0, 0, 102, 101, 1, 0, 0, 0, 102, 103,
		1, 0, 0, 0, 103, 1, 1, 0, 0, 0, 104, 112, 3, 26, 13, 0, 105, 112, 3, 4,
		2, 0, 106, 112, 3, 12, 6, 0, 107, 112, 3, 18, 9, 0, 108, 112, 3, 20, 10,
		0, 109, 112, 3, 22, 11, 0, 110, 112, 3, 24, 12, 0, 111, 104, 1, 0, 0, 0,
		111, 105, 1, 0, 0, 0, 111, 106, 1, 0, 0, 0, 111, 107, 1, 0, 0, 0, 111,
		108, 1, 0, 0, 0, 111, 109, 1, 0, 0, 0, 111, 110, 1, 0, 0, 0, 112, 3, 1,
		0, 0, 0, 113, 114, 5, 1, 0, 0, 114, 118, 5, 2, 0, 0, 115, 116, 5, 12, 0,
		0, 116, 117, 5, 13, 0, 0, 117, 119, 5, 14, 0, 0, 118, 115, 1, 0, 0, 0,
		118, 119, 1, 0, 0, 0, 119, 120, 1, 0, 0, 0, 120, 121, 3, 92, 46, 0, 121,
		122, 5, 63, 0, 0, 122, 127, 3, 6, 3, 0, 123, 124, 5, 66, 0, 0, 124, 126,
		3, 6, 3, 0, 125, 123, 1, 0, 0, 0, 126, 129, 1, 0, 0, 0, 127, 125, 1, 0,
		0, 0, 127, 128, 1, 0, 0, 0, 128, 131, 1, 0, 0, 0, 129, 127, 1, 0, 0, 0,
		130, 132, 5, 66, 0, 0, 131, 130, 1, 0, 0, 0, 131, 132, 1, 0, 0, 0, 132,
		146, 1, 0, 0, 0, 133, 134, 5, 6, 0, 0, 134, 135, 5, 7, 0, 0, 135, 136,
		5, 63, 0, 0, 136, 141, 3, 90, 45, 0, 137, 138, 5, 66, 0, 0, 138, 140, 3,
		90, 45, 0, 139, 137, 1, 0, 0, 0, 140, 143, 1, 0, 0, 0, 141, 139, 1, 0,
		0, 0, 141, 142, 1, 0, 0, 0, 142, 144, 1, 0, 0, 0, 143, 141, 1, 0, 0, 0,
		144, 145, 5, 64, 0, 0, 145, 147, 1, 0, 0, 0, 146, 133, 1, 0, 0, 0, 146,
		147, 1, 0, 0, 0, 147, 148, 1, 0, 0, 0, 148, 149, 5, 64, 0, 0, 149, 5, 1,
		0, 0, 0, 150, 151, 3, 90, 45, 0, 151, 152, 3, 10, 5, 0, 152, 153, 3, 8,
		4, 0, 153, 154, 5, 6, 0, 0, 154, 155, 5, 7, 0, 0, 155, 161, 1, 0, 0, 0,
		156, 157, 3, 90, 45, 0, 157, 158, 3, 10, 5, 0, 158, 159, 3, 8, 4, 0, 159,
		161, 1, 0, 0, 0, 160, 150, 1, 0, 0, 0, 160, 156, 1, 0, 0, 0, 161, 7, 1,
		0, 0, 0, 162, 163, 5, 13, 0, 0, 163, 168, 5, 15, 0, 0, 164, 168, 5, 16,
		0, 0, 165, 166, 5, 17, 0, 0, 166, 168, 3, 98, 49, 0, 167, 162, 1, 0, 0,
		0, 167, 164, 1, 0, 0, 0, 167, 165, 1, 0, 0, 0, 168, 171, 1, 0, 0, 0, 169,
		167, 1, 0, 0, 0, 169, 170, 1, 0, 0, 0, 170, 9, 1, 0, 0, 0, 171, 169, 1,
		0, 0, 0, 172, 192, 5, 51, 0, 0, 173, 192, 5, 52, 0, 0, 174, 192, 5, 53,
		0, 0, 175, 192, 5, 54, 0, 0, 176, 192, 5, 55, 0, 0, 177, 181, 5, 56, 0,
		0, 178, 179, 5, 63, 0, 0, 179, 180, 5, 89, 0, 0, 180, 182, 5, 64, 0, 0,
		181, 178, 1, 0, 0, 0, 181, 182, 1, 0, 0, 0, 182, 192, 1, 0, 0, 0, 183,
		187, 5, 57, 0, 0, 184, 185, 5, 63, 0, 0, 185, 186, 5, 89, 0, 0, 186, 188,
		5, 64, 0, 0, 187, 184, 1, 0, 0, 0, 187, 188, 1, 0, 0, 0, 188, 192, 1, 0,
		0, 0, 189, 192, 5, 58, 0, 0, 190, 192, 5, 59, 0, 0, 191, 172, 1, 0, 0,
		0, 191, 173, 1, 0, 0, 0, 191, 174, 1, 0, 0, 0, 191, 175, 1, 0, 0, 0, 191,
		176, 1, 0, 0, 0, 191, 177, 1, 0, 0, 0, 191, 183, 1, 0, 0, 0, 191, 189,
		1, 0, 0, 0, 191, 190, 1, 0, 0, 0, 192, 11, 1, 0, 0, 0, 193, 194, 5, 3,
		0, 0, 194, 195, 5, 4, 0, 0, 195, 196, 3, 92, 46, 0, 196, 197, 5, 63, 0,
		0, 197, 198, 3, 14, 7, 0, 198, 199, 5, 64, 0, 0, 199, 200, 5, 5, 0, 0,
		200, 205, 3, 16, 8, 0, 201, 202, 5, 66, 0, 0, 202, 204, 3, 16, 8, 0, 203,
		201, 1, 0, 0, 0, 204, 207, 1, 0, 0, 0, 205, 203, 1, 0, 0, 0, 205, 206,
		1, 0, 0, 0, 206, 13, 1, 0, 0, 0, 207, 205, 1, 0, 0, 0, 208, 213, 5, 86,
		0, 0, 209, 210, 5, 66, 0, 0, 210, 212, 5, 86, 0, 0, 211, 209, 1, 0, 0,
		0, 212, 215, 1, 0, 0, 0, 213, 211, 1, 0, 0, 0, 213, 214, 1, 0, 0, 0, 214,
		15, 1, 0, 0, 0, 215, 213, 1, 0, 0, 0, 216, 217, 5, 63, 0, 0, 217, 222,
		3, 52, 26, 0, 218, 219, 5, 66, 0, 0, 219, 221, 3, 52, 26, 0, 220, 218,
		1, 0, 0, 0, 221, 224, 1, 0, 0, 0, 222, 220, 1, 0, 0, 0, 222, 223, 1, 0,
		0, 0, 223, 225, 1, 0, 0, 0, 224, 222, 1, 0, 0, 0, 225, 226, 5, 64, 0, 0,
		226, 17, 1, 0, 0, 0, 227, 228, 5, 8, 0, 0, 228, 229, 5, 2, 0, 0, 229, 230,
		3, 92, 46, 0, 230, 19, 1, 0, 0, 0, 231, 233, 5, 9, 0, 0, 232, 234, 5, 2,
		0, 0, 233, 232, 1, 0, 0, 0, 233, 234, 1, 0, 0, 0, 234, 235, 1, 0, 0, 0,
		235, 236, 3, 92, 46, 0, 236, 21, 1, 0, 0, 0, 237, 238, 5, 10, 0, 0, 238,
		239, 5, 86, 0, 0, 239, 240, 7, 0, 0, 0, 240, 241, 3, 98, 49, 0, 241, 23,
		1, 0, 0, 0, 242, 243, 5, 18, 0, 0, 243, 247, 5, 86, 0, 0, 244, 246, 5,
		86, 0, 0, 245, 244, 1, 0, 0, 0, 246, 249, 1, 0, 0, 0, 247, 245, 1, 0, 0,
		0, 247, 248, 1, 0, 0, 0, 248, 25, 1, 0, 0, 0, 249, 247, 1, 0, 0, 0, 250,
		251, 5, 19, 0, 0, 251, 271, 3, 28, 14, 0, 252, 253, 5, 20, 0, 0, 253, 256,
		3, 34, 17, 0, 254, 255, 5, 21, 0, 0, 255, 257, 3, 40, 20, 0, 256, 254,
		1, 0, 0, 0, 256, 257, 1, 0, 0, 0, 257, 261, 1, 0, 0, 0, 258, 259, 5, 24,
		0, 0, 259, 260, 5, 25, 0, 0, 260, 262, 3, 42, 21, 0, 261, 258, 1, 0, 0,
		0, 261, 262, 1, 0, 0, 0, 262, 265, 1, 0, 0, 0, 263, 264, 5, 22, 0, 0, 264,
		266, 3, 46, 23, 0, 265, 263, 1, 0, 0, 0, 265, 266, 1, 0, 0, 0, 266, 269,
		1, 0, 0, 0, 267, 268, 5, 23, 0, 0, 268, 270, 3, 48, 24, 0, 269, 267, 1,
		0, 0, 0, 269, 270, 1, 0, 0, 0, 270, 272, 1, 0, 0, 0, 271, 252, 1, 0, 0,
		0, 271, 272, 1, 0, 0, 0, 272, 27, 1, 0, 0, 0, 273, 278, 3, 30, 15, 0, 274,
		275, 5, 66, 0, 0, 275, 277, 3, 30, 15, 0, 276, 274, 1, 0, 0, 0, 277, 280,
		1, 0, 0, 0, 278, 276, 1, 0, 0, 0, 278, 279, 1, 0, 0, 0, 279, 29, 1, 0,
		0, 0, 280, 278, 1, 0, 0, 0, 281, 283, 3, 52, 26, 0, 282, 284, 3, 50, 25,
		0, 283, 282, 1, 0, 0, 0, 283, 284, 1, 0, 0, 0, 284, 287, 1, 0, 0, 0, 285,
		287, 3, 32, 16, 0, 286, 281, 1, 0, 0, 0, 286, 285, 1, 0, 0, 0, 287, 31,
		1, 0, 0, 0, 288, 289, 5, 68, 0, 0, 289, 33, 1, 0, 0, 0, 290, 294, 3, 36,
		18, 0, 291, 293, 3, 38, 19, 0, 292, 291, 1, 0, 0, 0, 293, 296, 1, 0, 0,
		0, 294, 292, 1, 0, 0, 0, 294, 295, 1, 0, 0, 0, 295, 35, 1, 0, 0, 0, 296,
		294, 1, 0, 0, 0, 297, 299, 3, 92, 46, 0, 298, 300, 3, 50, 25, 0, 299, 298,
		1, 0, 0, 0, 299, 300, 1, 0, 0, 0, 300, 307, 1, 0, 0, 0, 301, 302, 5, 63,
		0, 0, 302, 303, 3, 26, 13, 0, 303, 304, 5, 64, 0, 0, 304, 305, 3, 50, 25,
		0, 305, 307, 1, 0, 0, 0, 306, 297, 1, 0, 0, 0, 306, 301, 1, 0, 0, 0, 307,
		37, 1, 0, 0, 0, 308, 310, 7, 1, 0, 0, 309, 308, 1, 0, 0, 0, 309, 310, 1,
		0, 0, 0, 310, 311, 1, 0, 0, 0, 311, 312, 5, 28, 0, 0, 312, 314, 3, 88,
		44, 0, 313, 315, 3, 50, 25, 0, 314, 313, 1, 0, 0, 0, 314, 315, 1, 0, 0,
		0, 315, 316, 1, 0, 0, 0, 316, 317, 5, 34, 0, 0, 317, 318, 3, 52, 26, 0,
		318, 39, 1, 0, 0, 0, 319, 320, 3, 52, 26, 0, 320, 41, 1, 0, 0, 0, 321,
		326, 3, 44, 22, 0, 322, 323, 5, 66, 0, 0, 323, 325, 3, 44, 22, 0, 324,
		322, 1, 0, 0, 0, 325, 328, 1, 0, 0, 0, 326, 324, 1, 0, 0, 0, 326, 327,
		1, 0, 0, 0, 327, 43, 1, 0, 0, 0, 328, 326, 1, 0, 0, 0, 329, 331, 3, 52,
		26, 0, 330, 332, 7, 2, 0, 0, 331, 330, 1, 0, 0, 0, 331, 332, 1, 0, 0, 0,
		332, 45, 1, 0, 0, 0, 333, 334, 5, 89, 0, 0, 334, 47, 1, 0, 0, 0, 335, 336,
		5, 89, 0, 0, 336, 49, 1, 0, 0, 0, 337, 339, 5, 46, 0, 0, 338, 337, 1, 0,
		0, 0, 338, 339, 1, 0, 0, 0, 339, 340, 1, 0, 0, 0, 340, 341, 5, 86, 0, 0,
		341, 51, 1, 0, 0, 0, 342, 343, 3, 54, 27, 0, 343, 53, 1, 0, 0, 0, 344,
		349, 3, 56, 28, 0, 345, 346, 5, 37, 0, 0, 346, 348, 3, 56, 28, 0, 347,
		345, 1, 0, 0, 0, 348, 351, 1, 0, 0, 0, 349, 347, 1, 0, 0, 0, 349, 350,
		1, 0, 0, 0, 350, 55, 1, 0, 0, 0, 351, 349, 1, 0, 0, 0, 352, 357, 3, 58,
		29, 0, 353, 354, 5, 36, 0, 0, 354, 356, 3, 58, 29, 0, 355, 353, 1, 0, 0,
		0, 356, 359, 1, 0, 0, 0, 357, 355, 1, 0, 0, 0, 357, 358, 1, 0, 0, 0, 358,
		57, 1, 0, 0, 0, 359, 357, 1, 0, 0, 0, 360, 362, 5, 13, 0, 0, 361, 360,
		1, 0, 0, 0, 361, 362, 1, 0, 0, 0, 362, 363, 1, 0, 0, 0, 363, 364, 3, 60,
		30, 0, 364, 59, 1, 0, 0, 0, 365, 397, 3, 62, 31, 0, 366, 369, 3, 94, 47,
		0, 367, 369, 3, 96, 48, 0, 368, 366, 1, 0, 0, 0, 368, 367, 1, 0, 0, 0,
		369, 370, 1, 0, 0, 0, 370, 371, 3, 62, 31, 0, 371, 398, 1, 0, 0, 0, 372,
		374, 5, 13, 0, 0, 373, 372, 1, 0, 0, 0, 373, 374, 1, 0, 0, 0, 374, 375,
		1, 0, 0, 0, 375, 376, 5, 35, 0, 0, 376, 377, 5, 63, 0, 0, 377, 382, 3,
		52, 26, 0, 378, 379, 5, 66, 0, 0, 379, 381, 3, 52, 26, 0, 380, 378, 1,
		0, 0, 0, 381, 384, 1, 0, 0, 0, 382, 380, 1, 0, 0, 0, 382, 383, 1, 0, 0,
		0, 383, 385, 1, 0, 0, 0, 384, 382, 1, 0, 0, 0, 385, 386, 5, 64, 0, 0, 386,
		398, 1, 0, 0, 0, 387, 389, 5, 13, 0, 0, 388, 387, 1, 0, 0, 0, 388, 389,
		1, 0, 0, 0, 389, 390, 1, 0, 0, 0, 390, 391, 5, 39, 0, 0, 391, 398, 3, 62,
		31, 0, 392, 394, 5, 38, 0, 0, 393, 395, 5, 13, 0, 0, 394, 393, 1, 0, 0,
		0, 394, 395, 1, 0, 0, 0, 395, 396, 1, 0, 0, 0, 396, 398, 5, 15, 0, 0, 397,
		368, 1, 0, 0, 0, 397, 373, 1, 0, 0, 0, 397, 388, 1, 0, 0, 0, 397, 392,
		1, 0, 0, 0, 397, 398, 1, 0, 0, 0, 398, 61, 1, 0, 0, 0, 399, 404, 3, 64,
		32, 0, 400, 401, 5, 72, 0, 0, 401, 403, 3, 64, 32, 0, 402, 400, 1, 0, 0,
		0, 403, 406, 1, 0, 0, 0, 404, 402, 1, 0, 0, 0, 404, 405, 1, 0, 0, 0, 405,
		63, 1, 0, 0, 0, 406, 404, 1, 0, 0, 0, 407, 412, 3, 66, 33, 0, 408, 409,
		7, 3, 0, 0, 409, 411, 3, 66, 33, 0, 410, 408, 1, 0, 0, 0, 411, 414, 1,
		0, 0, 0, 412, 410, 1, 0, 0, 0, 412, 413, 1, 0, 0, 0, 413, 65, 1, 0, 0,
		0, 414, 412, 1, 0, 0, 0, 415, 420, 3, 68, 34, 0, 416, 417, 7, 4, 0, 0,
		417, 419, 3, 68, 34, 0, 418, 416, 1, 0, 0, 0, 419, 422, 1, 0, 0, 0, 420,
		418, 1, 0, 0, 0, 420, 421, 1, 0, 0, 0, 421, 67, 1, 0, 0, 0, 422, 420, 1,
		0, 0, 0, 423, 425, 7, 3, 0, 0, 424, 423, 1, 0, 0, 0, 424, 425, 1, 0, 0,
		0, 425, 426, 1, 0, 0, 0, 426, 427, 3, 70, 35, 0, 427, 69, 1, 0, 0, 0, 428,
		432, 3, 76, 38, 0, 429, 431, 3, 72, 36, 0, 430, 429, 1, 0, 0, 0, 431, 434,
		1, 0, 0, 0, 432, 430, 1, 0, 0, 0, 432, 433, 1, 0, 0, 0, 433, 71, 1, 0,
		0, 0, 434, 432, 1, 0, 0, 0, 435, 436, 5, 47, 0, 0, 436, 437, 5, 48, 0,
		0, 437, 438, 5, 49, 0, 0, 438, 448, 5, 88, 0, 0, 439, 440, 5, 85, 0, 0,
		440, 448, 3, 88, 44, 0, 441, 442, 5, 73, 0, 0, 442, 444, 3, 74, 37, 0,
		443, 441, 1, 0, 0, 0, 444, 445, 1, 0, 0, 0, 445, 443, 1, 0, 0, 0, 445,
		446, 1, 0, 0, 0, 446, 448, 1, 0, 0, 0, 447, 435, 1, 0, 0, 0, 447, 439,
		1, 0, 0, 0, 447, 443, 1, 0, 0, 0, 448, 73, 1, 0, 0, 0, 449, 450, 3, 88,
		44, 0, 450, 75, 1, 0, 0, 0, 451, 452, 5, 63, 0, 0, 452, 453, 3, 52, 26,
		0, 453, 454, 5, 64, 0, 0, 454, 463, 1, 0, 0, 0, 455, 463, 3, 78, 39, 0,
		456, 463, 3, 80, 40, 0, 457, 463, 3, 82, 41, 0, 458, 463, 3, 84, 42, 0,
		459, 463, 3, 90, 45, 0, 460, 463, 3, 98, 49, 0, 461, 463, 5, 87, 0, 0,
		462, 451, 1, 0, 0, 0, 462, 455, 1, 0, 0, 0, 462, 456, 1, 0, 0, 0, 462,
		457, 1, 0, 0, 0, 462, 458, 1, 0, 0, 0, 462, 459, 1, 0, 0, 0, 462, 460,
		1, 0, 0, 0, 462, 461, 1, 0, 0, 0, 463, 77, 1, 0, 0, 0, 464, 465, 5, 63,
		0, 0, 465, 466, 3, 26, 13, 0, 466, 467, 5, 64, 0, 0, 467, 79, 1, 0, 0,
		0, 468, 474, 5, 41, 0, 0, 469, 470, 5, 42, 0, 0, 470, 471, 3, 52, 26, 0,
		471, 472, 5, 43, 0, 0, 472, 473, 3, 52, 26, 0, 473, 475, 1, 0, 0, 0, 474,
		469, 1, 0, 0, 0, 475, 476, 1, 0, 0, 0, 476, 474, 1, 0, 0, 0, 476, 477,
		1, 0, 0, 0, 477, 480, 1, 0, 0, 0, 478, 479, 5, 44, 0, 0, 479, 481, 3, 52,
		26, 0, 480, 478, 1, 0, 0, 0, 480, 481, 1, 0, 0, 0, 481, 482, 1, 0, 0, 0,
		482, 483, 5, 45, 0, 0, 483, 81, 1, 0, 0, 0, 484, 485, 3, 88, 44, 0, 485,
		494, 5, 63, 0, 0, 486, 491, 3, 52, 26, 0, 487, 488, 5, 66, 0, 0, 488, 490,
		3, 52, 26, 0, 489, 487, 1, 0, 0, 0, 490, 493, 1, 0, 0, 0, 491, 489, 1,
		0, 0, 0, 491, 492, 1, 0, 0, 0, 492, 495, 1, 0, 0, 0, 493, 491, 1, 0, 0,
		0, 494, 486, 1, 0, 0, 0, 494, 495, 1, 0, 0, 0, 495, 496, 1, 0, 0, 0, 496,
		497, 5, 64, 0, 0, 497, 83, 1, 0, 0, 0, 498, 499, 5, 50, 0, 0, 499, 500,
		5, 63, 0, 0, 500, 501, 5, 86, 0, 0, 501, 502, 5, 20, 0, 0, 502, 503, 3,
		52, 26, 0, 503, 504, 5, 64, 0, 0, 504, 85, 1, 0, 0, 0, 505, 506, 7, 5,
		0, 0, 506, 87, 1, 0, 0, 0, 507, 515, 3, 86, 43, 0, 508, 511, 5, 67, 0,
		0, 509, 512, 3, 86, 43, 0, 510, 512, 5, 74, 0, 0, 511, 509, 1, 0, 0, 0,
		511, 510, 1, 0, 0, 0, 512, 514, 1, 0, 0, 0, 513, 508, 1, 0, 0, 0, 514,
		517, 1, 0, 0, 0, 515, 513, 1, 0, 0, 0, 515, 516, 1, 0, 0, 0, 516, 89, 1,
		0, 0, 0, 517, 515, 1, 0, 0, 0, 518, 519, 3, 88, 44, 0, 519, 91, 1, 0, 0,
		0, 520, 521, 3, 88, 44, 0, 521, 93, 1, 0, 0, 0, 522, 523, 7, 6, 0, 0, 523,
		95, 1, 0, 0, 0, 524, 525, 5, 84, 0, 0, 525, 526, 5, 63, 0, 0, 526, 527,
		3, 88, 44, 0, 527, 528, 5, 64, 0, 0, 528, 97, 1, 0, 0, 0, 529, 530, 7,
		7, 0, 0, 530, 99, 1, 0, 0, 0, 57, 102, 111, 118, 127, 131, 141, 146, 160,
		167, 169, 181, 187, 191, 205, 213, 222, 233, 247, 256, 261, 265, 269, 271,
		278, 283, 286, 294, 299, 306, 309, 314, 326, 331, 338, 349, 357, 361, 368,
		373, 382, 388, 394, 397, 404, 412, 420, 424, 432, 445, 447, 462, 476, 480,
		491, 494, 511, 515,
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
	sqlParserLIMIT             = 22
	sqlParserOFFSET            = 23
	sqlParserORDER             = 24
	sqlParserBY                = 25
	sqlParserASC               = 26
	sqlParserDESC              = 27
	sqlParserJOIN              = 28
	sqlParserINNER             = 29
	sqlParserLEFT              = 30
	sqlParserRIGHT             = 31
	sqlParserFULL              = 32
	sqlParserCROSS             = 33
	sqlParserON                = 34
	sqlParserIN                = 35
	sqlParserAND               = 36
	sqlParserOR                = 37
	sqlParserIS                = 38
	sqlParserLIKE              = 39
	sqlParserILIKE             = 40
	sqlParserCASE              = 41
	sqlParserWHEN              = 42
	sqlParserTHEN              = 43
	sqlParserELSE              = 44
	sqlParserEND               = 45
	sqlParserAS                = 46
	sqlParserAT                = 47
	sqlParserTIME              = 48
	sqlParserZONE              = 49
	sqlParserEXTRACT           = 50
	sqlParserSTRING_TYPE       = 51
	sqlParserINT_TYPE          = 52
	sqlParserFLOAT_TYPE        = 53
	sqlParserBOOL_TYPE         = 54
	sqlParserTEXT_TYPE         = 55
	sqlParserVARCHAR_TYPE      = 56
	sqlParserCHAR_TYPE         = 57
	sqlParserSERIAL_TYPE       = 58
	sqlParserTIMESTAMP_TYPE    = 59
	sqlParserCURRENT_TIMESTAMP = 60
	sqlParserTRUE              = 61
	sqlParserFALSE             = 62
	sqlParserLPAREN            = 63
	sqlParserRPAREN            = 64
	sqlParserSEMICOLON         = 65
	sqlParserCOMMA             = 66
	sqlParserDOT               = 67
	sqlParserSTAR              = 68
	sqlParserPLUS              = 69
	sqlParserMINUS             = 70
	sqlParserSLASH             = 71
	sqlParserCONCAT            = 72
	sqlParserCOLON_COLON       = 73
	sqlParserTILDE             = 74
	sqlParserNREGEX            = 75
	sqlParserIREGEX            = 76
	sqlParserNIREGEX           = 77
	sqlParserEQ                = 78
	sqlParserGT                = 79
	sqlParserLT                = 80
	sqlParserGE                = 81
	sqlParserLE                = 82
	sqlParserNE                = 83
	sqlParserOPERATOR_KW       = 84
	sqlParserCOLLATE           = 85
	sqlParserIDENTIFIER        = 86
	sqlParserPARAMETER         = 87
	sqlParserSTRING_LITERAL    = 88
	sqlParserNUMBER            = 89
	sqlParserBLOCK_COMMENT     = 90
	sqlParserLINE_COMMENT      = 91
	sqlParserWS                = 92
)

// sqlParser rules.
const (
	sqlParserRULE_query                    = 0
	sqlParserRULE_statement                = 1
	sqlParserRULE_createTableStatement     = 2
	sqlParserRULE_columnDefinition         = 3
	sqlParserRULE_columnConstraints        = 4
	sqlParserRULE_dataType                 = 5
	sqlParserRULE_insertStatement          = 6
	sqlParserRULE_columnList               = 7
	sqlParserRULE_valueList                = 8
	sqlParserRULE_dropTableStatement       = 9
	sqlParserRULE_truncateTableStatement   = 10
	sqlParserRULE_setStatement             = 11
	sqlParserRULE_showStatement            = 12
	sqlParserRULE_selectStatement          = 13
	sqlParserRULE_selectList               = 14
	sqlParserRULE_selectItem               = 15
	sqlParserRULE_selectAll                = 16
	sqlParserRULE_fromClause               = 17
	sqlParserRULE_tableFactor              = 18
	sqlParserRULE_joinClause               = 19
	sqlParserRULE_whereClause              = 20
	sqlParserRULE_orderByClause            = 21
	sqlParserRULE_orderByItem              = 22
	sqlParserRULE_limitValue               = 23
	sqlParserRULE_offsetValue              = 24
	sqlParserRULE_alias                    = 25
	sqlParserRULE_expression               = 26
	sqlParserRULE_orExpression             = 27
	sqlParserRULE_andExpression            = 28
	sqlParserRULE_notExpression            = 29
	sqlParserRULE_comparisonExpression     = 30
	sqlParserRULE_concatExpression         = 31
	sqlParserRULE_additiveExpression       = 32
	sqlParserRULE_multiplicativeExpression = 33
	sqlParserRULE_unaryExpression          = 34
	sqlParserRULE_castExpression           = 35
	sqlParserRULE_postfix                  = 36
	sqlParserRULE_typeName                 = 37
	sqlParserRULE_primaryExpression        = 38
	sqlParserRULE_subqueryExpression       = 39
	sqlParserRULE_caseExpression           = 40
	sqlParserRULE_functionCall             = 41
	sqlParserRULE_extractFunction          = 42
	sqlParserRULE_namePart                 = 43
	sqlParserRULE_qualifiedName            = 44
	sqlParserRULE_columnName               = 45
	sqlParserRULE_tableName                = 46
	sqlParserRULE_operator                 = 47
	sqlParserRULE_operatorExpr             = 48
	sqlParserRULE_value                    = 49
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
		p.SetState(100)
		p.Statement()
	}
	p.SetState(102)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserSEMICOLON {
		{
			p.SetState(101)
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
	p.SetState(111)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case sqlParserSELECT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(104)
			p.SelectStatement()
		}

	case sqlParserCREATE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(105)
			p.CreateTableStatement()
		}

	case sqlParserINSERT:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(106)
			p.InsertStatement()
		}

	case sqlParserDROP:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(107)
			p.DropTableStatement()
		}

	case sqlParserTRUNCATE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(108)
			p.TruncateTableStatement()
		}

	case sqlParserSET:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(109)
			p.SetStatement()
		}

	case sqlParserSHOW:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(110)
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
		p.SetState(113)
		p.Match(sqlParserCREATE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(114)
		p.Match(sqlParserTABLE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(118)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserIF {
		{
			p.SetState(115)
			p.Match(sqlParserIF)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(116)
			p.Match(sqlParserNOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(117)
			p.Match(sqlParserEXISTS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(120)
		p.TableName()
	}
	{
		p.SetState(121)
		p.Match(sqlParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(122)
		p.ColumnDefinition()
	}
	p.SetState(127)
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
				p.SetState(123)
				p.Match(sqlParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(124)
				p.ColumnDefinition()
			}

		}
		p.SetState(129)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 3, p.GetParserRuleContext())
		if p.HasError() {
			goto errorExit
		}
	}
	p.SetState(131)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserCOMMA {
		{
			p.SetState(130)
			p.Match(sqlParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	p.SetState(146)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserPRIMARY {
		{
			p.SetState(133)
			p.Match(sqlParserPRIMARY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(134)
			p.Match(sqlParserKEY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
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
			p.ColumnName()
		}
		p.SetState(141)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == sqlParserCOMMA {
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
				p.ColumnName()
			}

			p.SetState(143)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(144)
			p.Match(sqlParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(148)
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
	p.SetState(160)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 7, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(150)
			p.ColumnName()
		}
		{
			p.SetState(151)
			p.DataType()
		}
		{
			p.SetState(152)
			p.ColumnConstraints()
		}
		{
			p.SetState(153)
			p.Match(sqlParserPRIMARY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(154)
			p.Match(sqlParserKEY)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(156)
			p.ColumnName()
		}
		{
			p.SetState(157)
			p.DataType()
		}
		{
			p.SetState(158)
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
	p.SetState(169)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&204800) != 0 {
		p.SetState(167)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case sqlParserNOT:
			{
				p.SetState(162)
				p.Match(sqlParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(163)
				p.Match(sqlParserNULL)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case sqlParserUNIQUE:
			{
				p.SetState(164)
				p.Match(sqlParserUNIQUE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		case sqlParserDEFAULT:
			{
				p.SetState(165)
				p.Match(sqlParserDEFAULT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(166)
				p.Value()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}

		p.SetState(171)
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

	p.SetState(191)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case sqlParserSTRING_TYPE:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(172)
			p.Match(sqlParserSTRING_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserINT_TYPE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(173)
			p.Match(sqlParserINT_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserFLOAT_TYPE:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(174)
			p.Match(sqlParserFLOAT_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserBOOL_TYPE:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(175)
			p.Match(sqlParserBOOL_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserTEXT_TYPE:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(176)
			p.Match(sqlParserTEXT_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserVARCHAR_TYPE:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(177)
			p.Match(sqlParserVARCHAR_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(181)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserLPAREN {
			{
				p.SetState(178)
				p.Match(sqlParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(179)
				p.Match(sqlParserNUMBER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(180)
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
			p.SetState(183)
			p.Match(sqlParserCHAR_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(187)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserLPAREN {
			{
				p.SetState(184)
				p.Match(sqlParserLPAREN)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(185)
				p.Match(sqlParserNUMBER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(186)
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
			p.SetState(189)
			p.Match(sqlParserSERIAL_TYPE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserTIMESTAMP_TYPE:
		p.EnterOuterAlt(localctx, 9)
		{
			p.SetState(190)
			p.Match(sqlParserTIMESTAMP_TYPE)
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
	p.EnterRule(localctx, 12, sqlParserRULE_insertStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(193)
		p.Match(sqlParserINSERT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(194)
		p.Match(sqlParserINTO)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(195)
		p.TableName()
	}
	{
		p.SetState(196)
		p.Match(sqlParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(197)
		p.ColumnList()
	}
	{
		p.SetState(198)
		p.Match(sqlParserRPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(199)
		p.Match(sqlParserVALUES)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(200)
		p.ValueList()
	}
	p.SetState(205)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserCOMMA {
		{
			p.SetState(201)
			p.Match(sqlParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(202)
			p.ValueList()
		}

		p.SetState(207)
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
	p.EnterRule(localctx, 14, sqlParserRULE_columnList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(208)
		p.Match(sqlParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(213)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserCOMMA {
		{
			p.SetState(209)
			p.Match(sqlParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(210)
			p.Match(sqlParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(215)
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
	p.EnterRule(localctx, 16, sqlParserRULE_valueList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(216)
		p.Match(sqlParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(217)
		p.Expression()
	}
	p.SetState(222)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserCOMMA {
		{
			p.SetState(218)
			p.Match(sqlParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(219)
			p.Expression()
		}

		p.SetState(224)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(225)
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
	p.EnterRule(localctx, 18, sqlParserRULE_dropTableStatement)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(227)
		p.Match(sqlParserDROP)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(228)
		p.Match(sqlParserTABLE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(229)
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
	p.EnterRule(localctx, 20, sqlParserRULE_truncateTableStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(231)
		p.Match(sqlParserTRUNCATE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(233)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserTABLE {
		{
			p.SetState(232)
			p.Match(sqlParserTABLE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(235)
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
	p.EnterRule(localctx, 22, sqlParserRULE_setStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(237)
		p.Match(sqlParserSET)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(238)
		p.Match(sqlParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(239)
		_la = p.GetTokenStream().LA(1)

		if !(_la == sqlParserTO || _la == sqlParserEQ) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(240)
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
	p.EnterRule(localctx, 24, sqlParserRULE_showStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(242)
		p.Match(sqlParserSHOW)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(243)
		p.Match(sqlParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(247)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserIDENTIFIER {
		{
			p.SetState(244)
			p.Match(sqlParserIDENTIFIER)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

		p.SetState(249)
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
	SELECT() antlr.TerminalNode
	SelectList() ISelectListContext
	FROM() antlr.TerminalNode
	FromClause() IFromClauseContext
	WHERE() antlr.TerminalNode
	WhereClause() IWhereClauseContext
	ORDER() antlr.TerminalNode
	BY() antlr.TerminalNode
	OrderByClause() IOrderByClauseContext
	LIMIT() antlr.TerminalNode
	LimitValue() ILimitValueContext
	OFFSET() antlr.TerminalNode
	OffsetValue() IOffsetValueContext

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

func (s *SelectStatementContext) SELECT() antlr.TerminalNode {
	return s.GetToken(sqlParserSELECT, 0)
}

func (s *SelectStatementContext) SelectList() ISelectListContext {
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

func (s *SelectStatementContext) FROM() antlr.TerminalNode {
	return s.GetToken(sqlParserFROM, 0)
}

func (s *SelectStatementContext) FromClause() IFromClauseContext {
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

func (s *SelectStatementContext) WHERE() antlr.TerminalNode {
	return s.GetToken(sqlParserWHERE, 0)
}

func (s *SelectStatementContext) WhereClause() IWhereClauseContext {
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

func (s *SelectStatementContext) ORDER() antlr.TerminalNode {
	return s.GetToken(sqlParserORDER, 0)
}

func (s *SelectStatementContext) BY() antlr.TerminalNode {
	return s.GetToken(sqlParserBY, 0)
}

func (s *SelectStatementContext) OrderByClause() IOrderByClauseContext {
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

func (s *SelectStatementContext) LIMIT() antlr.TerminalNode {
	return s.GetToken(sqlParserLIMIT, 0)
}

func (s *SelectStatementContext) LimitValue() ILimitValueContext {
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

func (s *SelectStatementContext) OFFSET() antlr.TerminalNode {
	return s.GetToken(sqlParserOFFSET, 0)
}

func (s *SelectStatementContext) OffsetValue() IOffsetValueContext {
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
	p.EnterRule(localctx, 26, sqlParserRULE_selectStatement)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(250)
		p.Match(sqlParserSELECT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(251)
		p.SelectList()
	}
	p.SetState(271)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserFROM {
		{
			p.SetState(252)
			p.Match(sqlParserFROM)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(253)
			p.FromClause()
		}
		p.SetState(256)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserWHERE {
			{
				p.SetState(254)
				p.Match(sqlParserWHERE)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(255)
				p.WhereClause()
			}

		}
		p.SetState(261)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserORDER {
			{
				p.SetState(258)
				p.Match(sqlParserORDER)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(259)
				p.Match(sqlParserBY)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(260)
				p.OrderByClause()
			}

		}
		p.SetState(265)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserLIMIT {
			{
				p.SetState(263)
				p.Match(sqlParserLIMIT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(264)
				p.LimitValue()
			}

		}
		p.SetState(269)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserOFFSET {
			{
				p.SetState(267)
				p.Match(sqlParserOFFSET)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(268)
				p.OffsetValue()
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
	p.EnterRule(localctx, 28, sqlParserRULE_selectList)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(273)
		p.SelectItem()
	}
	p.SetState(278)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserCOMMA {
		{
			p.SetState(274)
			p.Match(sqlParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(275)
			p.SelectItem()
		}

		p.SetState(280)
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
	p.EnterRule(localctx, 30, sqlParserRULE_selectItem)
	var _la int

	p.SetState(286)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case sqlParserNOT, sqlParserNULL, sqlParserDEFAULT, sqlParserCASE, sqlParserEXTRACT, sqlParserINT_TYPE, sqlParserFLOAT_TYPE, sqlParserBOOL_TYPE, sqlParserTEXT_TYPE, sqlParserVARCHAR_TYPE, sqlParserCHAR_TYPE, sqlParserSERIAL_TYPE, sqlParserTIMESTAMP_TYPE, sqlParserCURRENT_TIMESTAMP, sqlParserTRUE, sqlParserFALSE, sqlParserLPAREN, sqlParserPLUS, sqlParserMINUS, sqlParserIDENTIFIER, sqlParserPARAMETER, sqlParserSTRING_LITERAL, sqlParserNUMBER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(281)
			p.Expression()
		}
		p.SetState(283)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserAS || _la == sqlParserIDENTIFIER {
			{
				p.SetState(282)
				p.Alias()
			}

		}

	case sqlParserSTAR:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(285)
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
	p.EnterRule(localctx, 32, sqlParserRULE_selectAll)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(288)
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
	TableFactor() ITableFactorContext
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

func (s *FromClauseContext) TableFactor() ITableFactorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITableFactorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITableFactorContext)
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
	p.EnterRule(localctx, 34, sqlParserRULE_fromClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(290)
		p.TableFactor()
	}
	p.SetState(294)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&16911433728) != 0 {
		{
			p.SetState(291)
			p.JoinClause()
		}

		p.SetState(296)
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
	p.EnterRule(localctx, 36, sqlParserRULE_tableFactor)
	var _la int

	p.SetState(306)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case sqlParserDEFAULT, sqlParserINT_TYPE, sqlParserFLOAT_TYPE, sqlParserBOOL_TYPE, sqlParserTEXT_TYPE, sqlParserVARCHAR_TYPE, sqlParserCHAR_TYPE, sqlParserSERIAL_TYPE, sqlParserTIMESTAMP_TYPE, sqlParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(297)
			p.TableName()
		}
		p.SetState(299)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserAS || _la == sqlParserIDENTIFIER {
			{
				p.SetState(298)
				p.Alias()
			}

		}

	case sqlParserLPAREN:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(301)
			p.Match(sqlParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(302)
			p.SelectStatement()
		}
		{
			p.SetState(303)
			p.Match(sqlParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(304)
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
	p.EnterRule(localctx, 38, sqlParserRULE_joinClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(309)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if (int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&16642998272) != 0 {
		{
			p.SetState(308)
			_la = p.GetTokenStream().LA(1)

			if !((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&16642998272) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}
	{
		p.SetState(311)
		p.Match(sqlParserJOIN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(312)
		p.QualifiedName()
	}
	p.SetState(314)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserAS || _la == sqlParserIDENTIFIER {
		{
			p.SetState(313)
			p.Alias()
		}

	}
	{
		p.SetState(316)
		p.Match(sqlParserON)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(317)
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
	p.EnterRule(localctx, 40, sqlParserRULE_whereClause)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(319)
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
	p.EnterRule(localctx, 42, sqlParserRULE_orderByClause)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(321)
		p.OrderByItem()
	}
	p.SetState(326)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserCOMMA {
		{
			p.SetState(322)
			p.Match(sqlParserCOMMA)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(323)
			p.OrderByItem()
		}

		p.SetState(328)
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
	p.EnterRule(localctx, 44, sqlParserRULE_orderByItem)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(329)
		p.Expression()
	}
	p.SetState(331)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserASC || _la == sqlParserDESC {
		{
			p.SetState(330)
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
	p.EnterRule(localctx, 46, sqlParserRULE_limitValue)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(333)
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
	p.EnterRule(localctx, 48, sqlParserRULE_offsetValue)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(335)
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
	p.EnterRule(localctx, 50, sqlParserRULE_alias)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(338)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserAS {
		{
			p.SetState(337)
			p.Match(sqlParserAS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(340)
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
	p.EnterRule(localctx, 52, sqlParserRULE_expression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(342)
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
	p.EnterRule(localctx, 54, sqlParserRULE_orExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(344)
		p.AndExpression()
	}
	p.SetState(349)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserOR {
		{
			p.SetState(345)
			p.Match(sqlParserOR)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(346)
			p.AndExpression()
		}

		p.SetState(351)
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
	p.EnterRule(localctx, 56, sqlParserRULE_andExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(352)
		p.NotExpression()
	}
	p.SetState(357)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserAND {
		{
			p.SetState(353)
			p.Match(sqlParserAND)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(354)
			p.NotExpression()
		}

		p.SetState(359)
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
	p.EnterRule(localctx, 58, sqlParserRULE_notExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(361)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserNOT {
		{
			p.SetState(360)
			p.Match(sqlParserNOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	}
	{
		p.SetState(363)
		p.ComparisonExpression()
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
	p.EnterRule(localctx, 60, sqlParserRULE_comparisonExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(365)
		p.ConcatExpression()
	}
	p.SetState(397)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 42, p.GetParserRuleContext()) == 1 {
		p.SetState(368)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case sqlParserTILDE, sqlParserNREGEX, sqlParserIREGEX, sqlParserNIREGEX, sqlParserEQ, sqlParserGT, sqlParserLT, sqlParserGE, sqlParserLE, sqlParserNE:
			{
				p.SetState(366)
				p.Operator()
			}

		case sqlParserOPERATOR_KW:
			{
				p.SetState(367)
				p.OperatorExpr()
			}

		default:
			p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
			goto errorExit
		}
		{
			p.SetState(370)
			p.ConcatExpression()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	} else if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 42, p.GetParserRuleContext()) == 2 {
		p.SetState(373)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserNOT {
			{
				p.SetState(372)
				p.Match(sqlParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(375)
			p.Match(sqlParserIN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(376)
			p.Match(sqlParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(377)
			p.Expression()
		}
		p.SetState(382)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == sqlParserCOMMA {
			{
				p.SetState(378)
				p.Match(sqlParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(379)
				p.Expression()
			}

			p.SetState(384)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(385)
			p.Match(sqlParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	} else if p.HasError() { // JIM
		goto errorExit
	} else if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 42, p.GetParserRuleContext()) == 3 {
		p.SetState(388)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserNOT {
			{
				p.SetState(387)
				p.Match(sqlParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(390)
			p.Match(sqlParserLIKE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(391)
			p.ConcatExpression()
		}

	} else if p.HasError() { // JIM
		goto errorExit
	} else if p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 42, p.GetParserRuleContext()) == 4 {
		{
			p.SetState(392)
			p.Match(sqlParserIS)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(394)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserNOT {
			{
				p.SetState(393)
				p.Match(sqlParserNOT)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}

		}
		{
			p.SetState(396)
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
	p.EnterRule(localctx, 62, sqlParserRULE_concatExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(399)
		p.AdditiveExpression()
	}
	p.SetState(404)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserCONCAT {
		{
			p.SetState(400)
			p.Match(sqlParserCONCAT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(401)
			p.AdditiveExpression()
		}

		p.SetState(406)
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
	p.EnterRule(localctx, 64, sqlParserRULE_additiveExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(407)
		p.MultiplicativeExpression()
	}
	p.SetState(412)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserPLUS || _la == sqlParserMINUS {
		{
			p.SetState(408)
			_la = p.GetTokenStream().LA(1)

			if !(_la == sqlParserPLUS || _la == sqlParserMINUS) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(409)
			p.MultiplicativeExpression()
		}

		p.SetState(414)
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
	p.EnterRule(localctx, 66, sqlParserRULE_multiplicativeExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(415)
		p.UnaryExpression()
	}
	p.SetState(420)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserSTAR || _la == sqlParserSLASH {
		{
			p.SetState(416)
			_la = p.GetTokenStream().LA(1)

			if !(_la == sqlParserSTAR || _la == sqlParserSLASH) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(417)
			p.UnaryExpression()
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
	p.EnterRule(localctx, 68, sqlParserRULE_unaryExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(424)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserPLUS || _la == sqlParserMINUS {
		{
			p.SetState(423)
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
		p.SetState(426)
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
	p.EnterRule(localctx, 70, sqlParserRULE_castExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(428)
		p.PrimaryExpression()
	}
	p.SetState(432)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for (int64((_la-47)) & ^0x3f) == 0 && ((int64(1)<<(_la-47))&274945015809) != 0 {
		{
			p.SetState(429)
			p.Postfix()
		}

		p.SetState(434)
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
	p.EnterRule(localctx, 72, sqlParserRULE_postfix)
	var _alt int

	p.SetState(447)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetTokenStream().LA(1) {
	case sqlParserAT:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(435)
			p.Match(sqlParserAT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(436)
			p.Match(sqlParserTIME)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(437)
			p.Match(sqlParserZONE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(438)
			p.Match(sqlParserSTRING_LITERAL)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case sqlParserCOLLATE:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(439)
			p.Match(sqlParserCOLLATE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(440)
			p.QualifiedName()
		}

	case sqlParserCOLON_COLON:
		p.EnterOuterAlt(localctx, 3)
		p.SetState(443)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_alt = 1
		for ok := true; ok; ok = _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
			switch _alt {
			case 1:
				{
					p.SetState(441)
					p.Match(sqlParserCOLON_COLON)
					if p.HasError() {
						// Recognition error - abort rule
						goto errorExit
					}
				}
				{
					p.SetState(442)
					p.TypeName()
				}

			default:
				p.SetError(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
				goto errorExit
			}

			p.SetState(445)
			p.GetErrorHandler().Sync(p)
			_alt = p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 48, p.GetParserRuleContext())
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
	p.EnterRule(localctx, 74, sqlParserRULE_typeName)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(449)
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
	p.EnterRule(localctx, 76, sqlParserRULE_primaryExpression)
	p.SetState(462)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}

	switch p.GetInterpreter().AdaptivePredict(p.BaseParser, p.GetTokenStream(), 50, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(451)
			p.Match(sqlParserLPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(452)
			p.Expression()
		}
		{
			p.SetState(453)
			p.Match(sqlParserRPAREN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(455)
			p.SubqueryExpression()
		}

	case 3:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(456)
			p.CaseExpression()
		}

	case 4:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(457)
			p.FunctionCall()
		}

	case 5:
		p.EnterOuterAlt(localctx, 5)
		{
			p.SetState(458)
			p.ExtractFunction()
		}

	case 6:
		p.EnterOuterAlt(localctx, 6)
		{
			p.SetState(459)
			p.ColumnName()
		}

	case 7:
		p.EnterOuterAlt(localctx, 7)
		{
			p.SetState(460)
			p.Value()
		}

	case 8:
		p.EnterOuterAlt(localctx, 8)
		{
			p.SetState(461)
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
	p.EnterRule(localctx, 78, sqlParserRULE_subqueryExpression)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(464)
		p.Match(sqlParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(465)
		p.SelectStatement()
	}
	{
		p.SetState(466)
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
	p.EnterRule(localctx, 80, sqlParserRULE_caseExpression)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(468)
		p.Match(sqlParserCASE)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(474)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == sqlParserWHEN {
		{
			p.SetState(469)
			p.Match(sqlParserWHEN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(470)
			p.Expression()
		}
		{
			p.SetState(471)
			p.Match(sqlParserTHEN)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(472)
			p.Expression()
		}

		p.SetState(476)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)
	}
	p.SetState(480)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if _la == sqlParserELSE {
		{
			p.SetState(478)
			p.Match(sqlParserELSE)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		{
			p.SetState(479)
			p.Expression()
		}

	}
	{
		p.SetState(482)
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
	AllExpression() []IExpressionContext
	Expression(i int) IExpressionContext
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

func (s *FunctionCallContext) AllExpression() []IExpressionContext {
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

func (s *FunctionCallContext) Expression(i int) IExpressionContext {
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
	p.EnterRule(localctx, 82, sqlParserRULE_functionCall)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(484)
		p.QualifiedName()
	}
	{
		p.SetState(485)
		p.Match(sqlParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	p.SetState(494)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	if ((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&-3375500697100288) != 0) || ((int64((_la-69)) & ^0x3f) == 0 && ((int64(1)<<(_la-69))&1966083) != 0) {
		{
			p.SetState(486)
			p.Expression()
		}
		p.SetState(491)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}
		_la = p.GetTokenStream().LA(1)

		for _la == sqlParserCOMMA {
			{
				p.SetState(487)
				p.Match(sqlParserCOMMA)
				if p.HasError() {
					// Recognition error - abort rule
					goto errorExit
				}
			}
			{
				p.SetState(488)
				p.Expression()
			}

			p.SetState(493)
			p.GetErrorHandler().Sync(p)
			if p.HasError() {
				goto errorExit
			}
			_la = p.GetTokenStream().LA(1)
		}

	}
	{
		p.SetState(496)
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
	p.EnterRule(localctx, 84, sqlParserRULE_extractFunction)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(498)
		p.Match(sqlParserEXTRACT)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(499)
		p.Match(sqlParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(500)
		p.Match(sqlParserIDENTIFIER)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(501)
		p.Match(sqlParserFROM)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(502)
		p.Expression()
	}
	{
		p.SetState(503)
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
	p.EnterRule(localctx, 86, sqlParserRULE_namePart)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(505)
		_la = p.GetTokenStream().LA(1)

		if !(((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&1148417904979607552) != 0) || _la == sqlParserIDENTIFIER) {
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
	p.EnterRule(localctx, 88, sqlParserRULE_qualifiedName)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(507)
		p.NamePart()
	}
	p.SetState(515)
	p.GetErrorHandler().Sync(p)
	if p.HasError() {
		goto errorExit
	}
	_la = p.GetTokenStream().LA(1)

	for _la == sqlParserDOT {
		{
			p.SetState(508)
			p.Match(sqlParserDOT)
			if p.HasError() {
				// Recognition error - abort rule
				goto errorExit
			}
		}
		p.SetState(511)
		p.GetErrorHandler().Sync(p)
		if p.HasError() {
			goto errorExit
		}

		switch p.GetTokenStream().LA(1) {
		case sqlParserDEFAULT, sqlParserINT_TYPE, sqlParserFLOAT_TYPE, sqlParserBOOL_TYPE, sqlParserTEXT_TYPE, sqlParserVARCHAR_TYPE, sqlParserCHAR_TYPE, sqlParserSERIAL_TYPE, sqlParserTIMESTAMP_TYPE, sqlParserIDENTIFIER:
			{
				p.SetState(509)
				p.NamePart()
			}

		case sqlParserTILDE:
			{
				p.SetState(510)
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

		p.SetState(517)
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
	p.EnterRule(localctx, 90, sqlParserRULE_columnName)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(518)
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
	p.EnterRule(localctx, 92, sqlParserRULE_tableName)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(520)
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
	p.EnterRule(localctx, 94, sqlParserRULE_operator)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(522)
		_la = p.GetTokenStream().LA(1)

		if !((int64((_la-74)) & ^0x3f) == 0 && ((int64(1)<<(_la-74))&1023) != 0) {
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
	p.EnterRule(localctx, 96, sqlParserRULE_operatorExpr)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(524)
		p.Match(sqlParserOPERATOR_KW)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(525)
		p.Match(sqlParserLPAREN)
		if p.HasError() {
			// Recognition error - abort rule
			goto errorExit
		}
	}
	{
		p.SetState(526)
		p.QualifiedName()
	}
	{
		p.SetState(527)
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
	TRUE() antlr.TerminalNode
	FALSE() antlr.TerminalNode
	NULL() antlr.TerminalNode

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

func (s *ValueContext) TRUE() antlr.TerminalNode {
	return s.GetToken(sqlParserTRUE, 0)
}

func (s *ValueContext) FALSE() antlr.TerminalNode {
	return s.GetToken(sqlParserFALSE, 0)
}

func (s *ValueContext) NULL() antlr.TerminalNode {
	return s.GetToken(sqlParserNULL, 0)
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
	p.EnterRule(localctx, 98, sqlParserRULE_value)
	var _la int

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(529)
		_la = p.GetTokenStream().LA(1)

		if !(((int64(_la) & ^0x3f) == 0 && ((int64(1)<<_la)&8070450532247961600) != 0) || _la == sqlParserSTRING_LITERAL || _la == sqlParserNUMBER) {
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
