package scan

import "fmt"

// TokenType представляет тип токена
type TokenType int

const (
	// Специальные токены
	EOF TokenType = iota
	ILLEGAL
	COMMENT

	// Идентификаторы и литералы
	IDENT    // main, x, y
	INT      // 123
	FLOAT    // 123.45
	STRING   // "hello"
	VARIABLE // $var

	// Операторы
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	BANG     // !
	ASTERISK // *
	SLASH    // /
	MOD      // %

	LT        // <
	GT        // >
	LTE       // <=
	GTE       // >=
	EQ        // ==
	NOT_EQ    // !=
	IDENTICAL // ===
	NOT_IDENT // !==

	LOGICAL_AND // &&
	LOGICAL_OR  // ||

	// Разделители
	COMMA        // ,
	SEMICOLON    // ;
	COLON        // :
	DOT          // .
	QUESTION     // ?
	ARROW        // ->
	DOUBLE_ARROW // =>
	COALESCE     // ??

	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]

	// Ключевые слова
	FUNCTION
	CLASS
	INTERFACE
	TRAIT
	EXTENDS
	IMPLEMENTS
	PUBLIC
	PRIVATE
	PROTECTED
	STATIC
	ABSTRACT
	FINAL
	CONST

	IF
	ELSE
	ELSEIF
	SWITCH
	CASE
	DEFAULT
	WHILE
	DO
	FOR
	FOREACH
	BREAK
	CONTINUE
	RETURN
	TRY
	CATCH
	FINALLY
	THROW
	NEW
	CLONE
	INSTANCEOF

	ECHO
	PRINT
	INCLUDE
	INCLUDE_ONCE
	REQUIRE
	REQUIRE_ONCE
	USE
	NAMESPACE
	AS

	NULL
	TRUE
	FALSE

	// Типы данных
	ARRAY
	CALLABLE
	BOOL
	INT_TYPE
	FLOAT_TYPE
	STRING_TYPE
	OBJECT

	// Специальные конструкции
	HEREDOC_START
	HEREDOC_END
	NOWDOC_START
	NOWDOC_END
)

// Token представляет токен
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// String возвращает строковое представление TokenType
func (tt TokenType) String() string {
	switch tt {
	case EOF:
		return "EOF"
	case ILLEGAL:
		return "ILLEGAL"
	case IDENT:
		return "IDENT"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case STRING:
		return "STRING"
	case VARIABLE:
		return "VARIABLE"
	// ... остальные токены
	default:
		return fmt.Sprintf("TokenType(%d)", int(tt))
	}
}

// keywords содержит ключевые слова PHP
var keywords = map[string]TokenType{
	"function":   FUNCTION,
	"class":      CLASS,
	"interface":  INTERFACE,
	"trait":      TRAIT,
	"extends":    EXTENDS,
	"implements": IMPLEMENTS,
	"public":     PUBLIC,
	"private":    PRIVATE,
	"protected":  PROTECTED,
	"static":     STATIC,
	"abstract":   ABSTRACT,
	"final":      FINAL,
	"const":      CONST,

	"if":         IF,
	"else":       ELSE,
	"elseif":     ELSEIF,
	"switch":     SWITCH,
	"case":       CASE,
	"default":    DEFAULT,
	"while":      WHILE,
	"do":         DO,
	"for":        FOR,
	"foreach":    FOREACH,
	"break":      BREAK,
	"continue":   CONTINUE,
	"return":     RETURN,
	"try":        TRY,
	"catch":      CATCH,
	"finally":    FINALLY,
	"throw":      THROW,
	"new":        NEW,
	"clone":      CLONE,
	"instanceof": INSTANCEOF,

	"echo":         ECHO,
	"print":        PRINT,
	"include":      INCLUDE,
	"include_once": INCLUDE_ONCE,
	"require":      REQUIRE,
	"require_once": REQUIRE_ONCE,
	"use":          USE,
	"namespace":    NAMESPACE,
	"as":           AS,

	"null":  NULL,
	"true":  TRUE,
	"false": FALSE,

	"array":    ARRAY,
	"callable": CALLABLE,
	"bool":     BOOL,
	"int":      INT_TYPE,
	"float":    FLOAT_TYPE,
	"string":   STRING_TYPE,
	"object":   OBJECT,
}
