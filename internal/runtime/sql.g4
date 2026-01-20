grammar sql;

// Parser Rules
query
    : (selectStatement | createTableStatement | insertStatement | dropTableStatement | setStatement | describeStatement | showStatement) SEMICOLON?
    ;

createTableStatement
    : CREATE TABLE (IF NOT EXISTS)? tableName LPAREN columnDefinition (COMMA columnDefinition)* COMMA? RPAREN
    ;

insertStatement
    : INSERT INTO tableName LPAREN columnList RPAREN VALUES valueList (COMMA valueList)*
    ;

dropTableStatement
    : DROP TABLE tableName
    ;

setStatement
    : SET IDENTIFIER (TO | EQ) value
    ;

describeStatement
    : DESCRIBE TABLE tableName
    ;

showStatement
    : SHOW IDENTIFIER (IDENTIFIER)*
    ;

columnDefinition
    : PRIMARY KEY columnName dataType columnConstraints
    | columnName dataType columnConstraints PRIMARY KEY
    | columnName dataType columnConstraints
    ;

columnConstraints
    : (NOT NULL | UNIQUE | DEFAULT value)*
    ;

dataType
    : STRING_TYPE | INT_TYPE | FLOAT_TYPE | BOOL_TYPE | TEXT_TYPE | VARCHAR_TYPE (LPAREN NUMBER RPAREN)? | SERIAL_TYPE | TIMESTAMP_TYPE
    ;

selectStatement
    : SELECT columnList (FROM fromClause (WHERE whereClause)? (ORDER BY orderByClause)? (LIMIT limitValue)?)?
    ;

fromClause
    : tableName (AS IDENTIFIER)? (joinClause)*
    ;

joinClause
    : (INNER | LEFT | RIGHT | FULL | CROSS | MERGE)? JOIN tableName (AS IDENTIFIER)? ON expression
    ;

columnList
    : STAR | selectItem (COMMA selectItem)*
    ;

selectItem
    : expression (AS IDENTIFIER)?
    ;

qualifiedName
    : IDENTIFIER (DOT IDENTIFIER)*
    ;

functionCall
    : qualifiedName LPAREN (expression (COMMA expression)*)? RPAREN
    ;

extractFunction
    : EXTRACT LPAREN IDENTIFIER FROM expression RPAREN
    ;

atTimeZoneExpression
    : primaryExpression AT TIME ZONE STRING_LITERAL
    ;

castExpression
    : primaryExpression COLON_COLON IDENTIFIER
    ;

caseExpression
    : CASE (WHEN expression THEN expression)+ (ELSE expression)? END
    ;

valueList
    : LPAREN expression (COMMA expression)* RPAREN
    ;

tableName
    : IDENTIFIER (DOT IDENTIFIER)?
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

expression
    : castExpression
    | atTimeZoneExpression
    | caseExpression
    | additiveExpression (operator additiveExpression)*
    ;

additiveExpression
    : multiplicativeExpression ((PLUS | MINUS) multiplicativeExpression)*
    ;

multiplicativeExpression
    : primaryExpression ((STAR | SLASH) primaryExpression)*
    ;

primaryExpression
    : LPAREN expression RPAREN
    | functionCall
    | extractFunction
    | columnName
    | value
    ;

operator
    : EQ | GT | LT | GE | LE | NE | LIKE | NOT LIKE | ILIKE
    ;

value
    : STRING_LITERAL | NUMBER | CURRENT_TIMESTAMP | TRUE | FALSE
    ;

limitValue
    : NUMBER
    ;

columnName
    : IDENTIFIER
    ;

// Lexer Rules
CREATE : 'CREATE' | 'create';
TABLE : 'TABLE';
INSERT : 'INSERT' | 'insert';
INTO : 'INTO' | 'into';
VALUES : 'VALUES' | 'values';
PRIMARY : 'PRIMARY';
KEY : 'KEY';
DROP : 'DROP' | 'drop';
DESCRIBE : 'DESCRIBE' | 'describe';
SET : 'SET' | 'set';
TO : 'TO' | 'to';
IF : 'IF' | 'if';
NOT : 'NOT' | 'not';
EXISTS : 'EXISTS' | 'exists';
NULL : 'NULL' | 'null';
UNIQUE : 'UNIQUE' | 'unique';
DEFAULT : 'DEFAULT' | 'default';
SHOW : 'SHOW' | 'show';
LPAREN : '(';
RPAREN : ')';
SEMICOLON : ';';
STRING_TYPE : 'STRING' | 'string';
INT_TYPE : 'INT' | 'int';
FLOAT_TYPE : 'FLOAT' | 'float';
BOOL_TYPE : 'BOOL' | 'bool';
TEXT_TYPE : 'TEXT' | 'text';
VARCHAR_TYPE : 'VARCHAR' | 'varchar';
SERIAL_TYPE : 'SERIAL' | 'serial';
TIMESTAMP_TYPE : 'TIMESTAMP' | 'timestamp';
CURRENT_TIMESTAMP : 'CURRENT_TIMESTAMP' | 'current_timestamp';

TRUE : 'TRUE' | 'true';
FALSE : 'FALSE' | 'false';
EXTRACT : 'EXTRACT' | 'extract';

SELECT : 'SELECT' | 'select';
FROM : 'FROM' | 'from';
WHERE : 'WHERE' | 'where';
LIMIT : 'LIMIT' | 'limit';
STAR : '*';
COMMA : ',';
DOT : '.';
COLON_COLON : '::';
EQ : '=';
GT : '>';
LT : '<';
GE : '>=';
LE : '<=';
NE : '!=' | '<>';
LIKE : 'LIKE' | 'like';
ILIKE : 'ILIKE' | 'ilike';
AS : 'AS' | 'as';
AT : 'AT' | 'at';
TIME : 'TIME' | 'time';
ZONE : 'ZONE' | 'zone';

JOIN : 'JOIN' | 'join';
INNER : 'INNER' | 'inner';
LEFT : 'LEFT' | 'left';
RIGHT : 'RIGHT' | 'right';
FULL : 'FULL' | 'full';
CROSS : 'CROSS' | 'cross';
MERGE : 'MERGE' | 'merge';
ON : 'ON' | 'on';

ORDER : 'ORDER' | 'order';
BY : 'BY' | 'by';
WHEN : 'WHEN' | 'when';
THEN : 'THEN' | 'then';
ELSE : 'ELSE' | 'else';
CASE : 'CASE' | 'case';
END : 'END' | 'end';
ASC : 'ASC' | 'asc';
DESC : 'DESC' | 'desc';

PLUS : '+';
MINUS : '-';
SLASH : '/';

IDENTIFIER : [a-zA-Z_][a-zA-Z0-9_]*;
STRING_LITERAL : '\'' (~['\\] | '\\' .)* '\'';
NUMBER : [0-9]+ ('.' [0-9]+)?;

BLOCK_COMMENT : '/*' .*? '*/' -> skip;
LINE_COMMENT : '--' ~[\r\n]* -> skip;
WS : [ \t\r\n]+ -> skip;