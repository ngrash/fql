grammar FQL;

/*
 * Parser Rules
 */

query : expression+ EOF;

expression : expression OR expression
	   | PARENS_OPEN expression PARENS_CLOSE
           | negation
           | condition
           | value
           ;

negation : NOT condition ;
condition : key op value ;

key : WORD ;
op : OP ;
value : WORD | STRING ;

/*
 * Lexer Rules
 */

STRING : '"' .*? '"' ;

fragment EQ : (':' | '=') ;
fragment GE : ('>=' | '=>') ;
fragment LE : ('<=' | '=<') ;
fragment LT : '<' ;
fragment GT : '>' ;
fragment NE : ('!=' | '=!') ;

OP : (EQ | GE | LE | LT | GT | NE);

WHITESPACE : ' ' -> skip ;
NEWLINE : '\n' -> skip ;

OR : ('or' | 'OR') ;

PARENS_OPEN : ('(' | '[' | '{') ;
PARENS_CLOSE : (')' | ']' | '}') ;

WORD : [a-zA-Z0-9] [a-zA-Z0-9_.-]* ;

NOT : ('-' | '!') ;
