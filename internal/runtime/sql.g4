grammar sql;

// Parser Rules
query
    : statement SEMICOLON?
    ;

statement
    : selectStatement 
    | createTableStatement 
    | insertStatement 
    | dropTableStatement 
    | truncateTableStatement 
    | setStatement 
    | showStatement 
    ;

// ============ CREATE TABLE ======================

createTableStatement
    : CREATE TABLE (IF NOT EXISTS)? tableName LPAREN 
        columnDefinition (COMMA columnDefinition)* COMMA? 
        (PRIMARY KEY LPAREN columnName (COMMA columnName)* RPAREN)?
    RPAREN
    ;

columnDefinition
    : columnName dataType columnConstraints PRIMARY KEY
    | columnName dataType columnConstraints
    ;

columnConstraints
    : (NOT NULL | UNIQUE | DEFAULT value)*
    ;

dataType
    : STRING_TYPE 
    | INT_TYPE 
    | FLOAT_TYPE 
    | BOOL_TYPE 
    | TEXT_TYPE 
    | VARCHAR_TYPE (LPAREN NUMBER RPAREN)? 
    | CHAR_TYPE (LPAREN NUMBER RPAREN)?
    | SERIAL_TYPE 
    | TIMESTAMP_TYPE
    ;

// ============ INSERT ======================

insertStatement
    : INSERT INTO tableName LPAREN 
        columnList 
    RPAREN VALUES valueList (COMMA valueList)*
    ;

columnList
    : IDENTIFIER (COMMA IDENTIFIER)*
    ;

valueList
    : LPAREN expression (COMMA expression)* RPAREN
    ;

// ============ DROP TABLE ======================

dropTableStatement
    : DROP TABLE tableName
    ;

// ============ TRUNCATE TABLE ======================

truncateTableStatement
    : TRUNCATE TABLE? tableName
    ;

// ============ SET ======================

setStatement
    : SET IDENTIFIER (TO | EQ) value
    ;


// ============ SHOW ======================

showStatement
    : SHOW IDENTIFIER (IDENTIFIER)*
    ;

// ============ SELECT ======================

selectStatement
    : SELECT selectList 
        (FROM fromClause 
            (WHERE whereClause)? 
            (ORDER BY orderByClause)? 
            (LIMIT limitValue)?
            (OFFSET offsetValue)?
        )?
    ;

selectList
    : selectItem (COMMA selectItem)*
    ;

selectItem
    : expression alias?
    | selectAll
    ;

selectAll
    : STAR
    ;

fromClause
    : tableFactor (joinClause)*
    ;

tableFactor
    : tableName alias?
    | LPAREN selectStatement RPAREN alias        // derived table требует alias
    ;

joinClause
    : (INNER | LEFT | RIGHT | FULL | CROSS)? JOIN qualifiedName alias? 
        ON expression
    ;

whereClause
    : expression
    ;

orderByClause
    : orderByItem (COMMA orderByItem)*
    ;

orderByItem
    : expression (ASC | DESC)?
    ;

limitValue
    : NUMBER
    ;

offsetValue
    : NUMBER
    ;

alias
    : AS? IDENTIFIER
    ;

// ============ EXPRESSIONS (no left recursion) ======================

expression
    : orExpression
    ;

orExpression
    : andExpression (OR andExpression)*
    ;

andExpression
    : notExpression (AND notExpression)*
    ;

notExpression
    : NOT? comparisonExpression
    ;

comparisonExpression
    : concatExpression (
        ((operator | operatorExpr) concatExpression)
        | (NOT? IN LPAREN expression (COMMA expression)* RPAREN)
        | (NOT? LIKE concatExpression)
        | (IS NOT? NULL)
    )?
    ;

concatExpression
    : additiveExpression (CONCAT additiveExpression)*
    ;

additiveExpression
    : multiplicativeExpression ((PLUS | MINUS) multiplicativeExpression)*
    ;

multiplicativeExpression
    : unaryExpression ((STAR | SLASH) unaryExpression)*
    ;

unaryExpression
    : (PLUS | MINUS)? castExpression
    ;

castExpression
    : primaryExpression postfix*
    ;

postfix
    : AT TIME ZONE STRING_LITERAL
    | COLLATE qualifiedName
    | (COLON_COLON typeName)+
    ;

typeName
    : qualifiedName
    ;

primaryExpression
    : LPAREN expression RPAREN
    | subqueryExpression
    | caseExpression
    | functionCall
    | extractFunction
    | columnName
    | value
    | PARAMETER
    ;

subqueryExpression
    : LPAREN selectStatement RPAREN
    ;

caseExpression
    : CASE 
        (WHEN expression THEN expression)+ 
        (ELSE expression)? 
    END
    ;

functionCall
    : qualifiedName LPAREN (expression (COMMA expression)*)? RPAREN
    ;

extractFunction
    : EXTRACT LPAREN IDENTIFIER FROM expression RPAREN
    ;

namePart
    : IDENTIFIER
    | DEFAULT        // чтобы работало pg_catalog.default
    | TEXT_TYPE      // for pg_catalog.text
    | INT_TYPE       // for casting ::integer, ::int
    | FLOAT_TYPE     // for casting ::float, ::decimal
    | BOOL_TYPE      // for casting ::boolean, ::bool
    | VARCHAR_TYPE   // for casting ::varchar
    | CHAR_TYPE      // for casting ::char
    | SERIAL_TYPE    // for casting ::serial
    | TIMESTAMP_TYPE // for casting ::timestamp
    ;

qualifiedName
    : namePart (DOT (namePart | TILDE))*
    ;
    
columnName
    : qualifiedName
    ;

tableName
    : qualifiedName
    ;

operator
    : EQ | GT | LT | GE | LE | NE
    | TILDE | NREGEX | IREGEX | NIREGEX
    ;

operatorExpr
    : OPERATOR_KW LPAREN qualifiedName RPAREN
    ;

value
    : STRING_LITERAL 
    | NUMBER 
    | CURRENT_TIMESTAMP 
    | TRUE 
    | FALSE
    | NULL
    ;

// ======= LEXER RULES ====================

CREATE : 'CREATE' | 'create';
TABLE : 'TABLE' | 'table';
INSERT : 'INSERT' | 'insert';
INTO : 'INTO' | 'into';
VALUES : 'VALUES' | 'values';
PRIMARY : 'PRIMARY' | 'primary';
KEY : 'KEY' | 'key';
DROP : 'DROP' | 'drop';
TRUNCATE : 'TRUNCATE' | 'truncate';
SET : 'SET' | 'set';
TO : 'TO' | 'to';
IF : 'IF' | 'if';
NOT : 'NOT' | 'not';
EXISTS : 'EXISTS' | 'exists';
NULL : 'NULL' | 'null';
UNIQUE : 'UNIQUE' | 'unique';
DEFAULT : 'DEFAULT' | 'default';
SHOW : 'SHOW' | 'show';
SELECT : 'SELECT' | 'select';
FROM : 'FROM' | 'from';
WHERE : 'WHERE' | 'where';
LIMIT : 'LIMIT' | 'limit';
OFFSET : 'OFFSET' | 'offset';
ORDER : 'ORDER' | 'order';
BY : 'BY' | 'by';
ASC : 'ASC' | 'asc';
DESC : 'DESC' | 'desc';

JOIN : 'JOIN' | 'join';
INNER : 'INNER' | 'inner';
LEFT : 'LEFT' | 'left';
RIGHT : 'RIGHT' | 'right';
FULL : 'FULL' | 'full';
CROSS : 'CROSS' | 'cross';
ON : 'ON' | 'on';

IN : 'IN' | 'in';
AND : 'AND' | 'and';
OR : 'OR' | 'or';
IS : 'IS' | 'is';
LIKE : 'LIKE' | 'like';
ILIKE : 'ILIKE' | 'ilike';

CASE : 'CASE' | 'case';
WHEN : 'WHEN' | 'when';
THEN : 'THEN' | 'then';
ELSE : 'ELSE' | 'else';
END : 'END' | 'end';

AS : 'AS' | 'as';
AT : 'AT' | 'at';
TIME : 'TIME' | 'time';
ZONE : 'ZONE' | 'zone';
EXTRACT : 'EXTRACT' | 'extract';

STRING_TYPE : 'STRING' | 'string';
INT_TYPE : 'INT' | 'int' | 'INTEGER' | 'integer' | 'BIGINT' | 'bigint';
FLOAT_TYPE : 'FLOAT' | 'float' | 'DECIMAL' | 'decimal';
BOOL_TYPE : 'BOOL' | 'bool' | 'BOOLEAN' | 'boolean';
TEXT_TYPE : 'TEXT' | 'text';
VARCHAR_TYPE : 'VARCHAR' | 'varchar';
CHAR_TYPE : 'CHAR' | 'char';
SERIAL_TYPE : 'SERIAL' | 'serial';
TIMESTAMP_TYPE : 'TIMESTAMP' | 'timestamp';

CURRENT_TIMESTAMP : 'CURRENT_TIMESTAMP' | 'current_timestamp';
TRUE : 'TRUE' | 'true';
FALSE : 'FALSE' | 'false';

LPAREN : '(';
RPAREN : ')';
SEMICOLON : ';';
COMMA : ',';
DOT : '.';
STAR : '*';
PLUS : '+';
MINUS : '-';
SLASH : '/';
CONCAT : '||';
COLON_COLON : '::';
TILDE       : '~';
NREGEX  : '!~';
IREGEX  : '~*';
NIREGEX : '!~*';

EQ : '=';
GT : '>';
LT : '<';
GE : '>=';
LE : '<=';
NE : '!=' | '<>';

OPERATOR_KW : 'OPERATOR' | 'operator';
COLLATE     : 'COLLATE'  | 'collate';

IDENTIFIER : [a-zA-Z_][a-zA-Z0-9_]*;
PARAMETER : '$' [0-9]+;
STRING_LITERAL : '\'' (~['\\] | '\\' .)* '\'';
NUMBER : [0-9]+ ('.' [0-9]+)?;

BLOCK_COMMENT : '/*' .*? '*/' -> skip;
LINE_COMMENT : '--' ~[\r\n]* -> skip;
WS : [ \t\r\n]+ -> skip;
