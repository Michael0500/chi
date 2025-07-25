package scan

import (
	"unicode"
	"unicode/utf8"
)

// Lexer представляет лексический анализатор PHP
type Lexer struct {
	input        string
	position     int  // текущая позиция в input (указывает на текущий символ)
	readPosition int  // текущая позиция чтения в input (после текущего символа)
	ch           rune // текущий символ
	line         int  // текущая строка
	column       int  // текущая колонка
}

// New создает новый лексер
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 1,
	}
	l.readChar()
	return l
}

// readChar читает следующий символ из input
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		var size int
		l.ch, size = utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.position = l.readPosition
		l.readPosition += size

		// Обновляем позицию
		if l.ch == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
	}
}

// peekChar возвращает следующий символ без продвижения позиции
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

// peekCharN возвращает символ через N позиций
func (l *Lexer) peekCharN(n int) rune {
	pos := l.readPosition
	var r rune
	for i := 0; i <= n; i++ {
		if pos >= len(l.input) {
			return 0
		}
		var size int
		r, size = utf8.DecodeRuneInString(l.input[pos:])
		pos += size
	}
	return r
}

// NextToken возвращает следующий токен
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			if l.peekCharN(1) == '=' {
				ch := l.ch
				l.readChar()
				l.readChar()
				tok = Token{
					Type:    IDENTICAL,
					Literal: string(ch) + "==",
					Line:    l.line,
					Column:  l.column - 3,
				}
			} else {
				ch := l.ch
				l.readChar()
				tok = Token{
					Type:    EQ,
					Literal: string(ch) + "=",
					Line:    l.line,
					Column:  l.column - 2,
				}
			}
		} else {
			tok = newToken(ASSIGN, l.ch, l.line, l.column)
		}
	case '+':
		tok = newToken(PLUS, l.ch, l.line, l.column)
	case '-':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    ARROW,
				Literal: string(ch) + ">",
				Line:    l.line,
				Column:  l.column - 2,
			}
		} else {
			tok = newToken(MINUS, l.ch, l.line, l.column)
		}
	case '!':
		if l.peekChar() == '=' {
			if l.peekCharN(1) == '=' {
				ch := l.ch
				l.readChar()
				l.readChar()
				tok = Token{
					Type:    NOT_IDENT,
					Literal: string(ch) + "!=",
					Line:    l.line,
					Column:  l.column - 3,
				}
			} else {
				ch := l.ch
				l.readChar()
				tok = Token{
					Type:    NOT_EQ,
					Literal: string(ch) + "=",
					Line:    l.line,
					Column:  l.column - 2,
				}
			}
		} else {
			tok = newToken(BANG, l.ch, l.line, l.column)
		}
	case '*':
		tok = newToken(ASTERISK, l.ch, l.line, l.column)
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			tok = Token{
				Type:    COMMENT,
				Literal: l.readSingleLineComment(),
				Line:    l.line,
				Column:  l.column - 2,
			}
			return tok
		} else if l.peekChar() == '*' {
			l.readChar()
			tok = Token{
				Type:    COMMENT,
				Literal: l.readMultiLineComment(),
				Line:    l.line,
				Column:  l.column - 2,
			}
			return tok
		} else {
			tok = newToken(SLASH, l.ch, l.line, l.column)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    LTE,
				Literal: string(ch) + "=",
				Line:    l.line,
				Column:  l.column - 2,
			}
		} else if l.peekChar() == '<' {
			// Проверяем на heredoc/nowdoc
			if l.peekCharN(1) == '<' {
				return l.readHeredoc()
			}
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    LT,
				Literal: string(ch),
				Line:    l.line,
				Column:  l.column - 1,
			}
		} else {
			tok = newToken(LT, l.ch, l.line, l.column)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    GTE,
				Literal: string(ch) + "=",
				Line:    l.line,
				Column:  l.column - 2,
			}
		} else {
			tok = newToken(GT, l.ch, l.line, l.column)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    LOGICAL_AND,
				Literal: string(ch) + "&",
				Line:    l.line,
				Column:  l.column - 2,
			}
		} else {
			tok = newToken(ILLEGAL, l.ch, l.line, l.column)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    LOGICAL_OR,
				Literal: string(ch) + "|",
				Line:    l.line,
				Column:  l.column - 2,
			}
		} else {
			tok = newToken(ILLEGAL, l.ch, l.line, l.column)
		}
	case '?':
		if l.peekChar() == '?' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    COALESCE,
				Literal: string(ch) + "?",
				Line:    l.line,
				Column:  l.column - 2,
			}
		} else {
			tok = newToken(QUESTION, l.ch, l.line, l.column)
		}
	case ':':
		if l.peekChar() == ':' {
			ch := l.ch
			l.readChar()
			tok = Token{
				Type:    DOUBLE_ARROW, // Это не совсем правильно, но для демонстрации
				Literal: string(ch) + ":",
				Line:    l.line,
				Column:  l.column - 2,
			}
		} else {
			tok = newToken(COLON, l.ch, l.line, l.column)
		}
	case ';':
		tok = newToken(SEMICOLON, l.ch, l.line, l.column)
	case ',':
		tok = newToken(COMMA, l.ch, l.line, l.column)
	case '{':
		tok = newToken(LBRACE, l.ch, l.line, l.column)
	case '}':
		tok = newToken(RBRACE, l.ch, l.line, l.column)
	case '(':
		tok = newToken(LPAREN, l.ch, l.line, l.column)
	case ')':
		tok = newToken(RPAREN, l.ch, l.line, l.column)
	case '[':
		tok = newToken(LBRACKET, l.ch, l.line, l.column)
	case ']':
		tok = newToken(RBRACKET, l.ch, l.line, l.column)
	case '.':
		tok = newToken(DOT, l.ch, l.line, l.column)
	case '"':
		tok = l.readString()
	case '\'':
		tok = l.readSingleQuoteString()
	case '`':
		tok = l.readBacktickString()
	case '$':
		if unicode.IsLetter(l.peekChar()) || l.peekChar() == '_' {
			tok = l.readVariable()
		} else {
			tok = newToken(ILLEGAL, l.ch, l.line, l.column)
		}
	case 0:
		tok = Token{Type: EOF, Literal: "", Line: l.line, Column: l.column}
	default:
		if unicode.IsLetter(l.ch) || l.ch == '_' {
			literal := l.readIdentifier()
			tokenType := lookupIdent(literal)
			tok = Token{Type: tokenType, Literal: literal, Line: l.line, Column: l.column - len(literal)}
			return tok
		} else if unicode.IsDigit(l.ch) {
			return l.readNumber()
		} else {
			tok = newToken(ILLEGAL, l.ch, l.line, l.column)
		}
	}

	l.readChar()
	return tok
}

// newToken создает новый токен
func newToken(tokenType TokenType, ch rune, line, column int) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
		Line:    line,
		Column:  column,
	}
}

// skipWhitespace пропускает пробельные символы
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) && l.ch != 0 {
		l.readChar()
	}
}

// readIdentifier читает идентификатор
func (l *Lexer) readIdentifier() string {
	position := l.position
	for unicode.IsLetter(l.ch) || unicode.IsDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber читает число (целое или с плавающей точкой)
func (l *Lexer) readNumber() Token {
	position := l.position

	// Читаем целую часть
	for unicode.IsDigit(l.ch) {
		l.readChar()
	}

	// Проверяем наличие десятичной точки
	if l.ch == '.' && unicode.IsDigit(l.peekChar()) {
		l.readChar() // точка
		// Читаем дробную часть
		for unicode.IsDigit(l.ch) {
			l.readChar()
		}

		// Проверяем экспоненту
		if l.ch == 'e' || l.ch == 'E' {
			l.readChar()
			if l.ch == '+' || l.ch == '-' {
				l.readChar()
			}
			for unicode.IsDigit(l.ch) {
				l.readChar()
			}
		}

		return Token{
			Type:    FLOAT,
			Literal: l.input[position:l.position],
			Line:    l.line,
			Column:  l.column - (l.position - position),
		}
	}

	return Token{
		Type:    INT,
		Literal: l.input[position:l.position],
		Line:    l.line,
		Column:  l.column - (l.position - position),
	}
}

// readString читает строку в двойных кавычках
func (l *Lexer) readString() Token {
	position := l.position + 1 // пропускаем открывающую кавычку
	l.readChar()               // пропускаем открывающую кавычку

	for l.ch != '"' && l.ch != 0 {
		// Обрабатываем экранированные символы
		if l.ch == '\\' {
			l.readChar() // пропускаем обратный слеш
		}
		l.readChar()
	}

	if l.ch == '"' {
		l.readChar() // пропускаем закрывающую кавычку
	}

	return Token{
		Type:    STRING,
		Literal: l.input[position : l.position-1],
		Line:    l.line,
		Column:  l.column - (l.position - position),
	}
}

// readSingleQuoteString читает строку в одинарных кавычках
func (l *Lexer) readSingleQuoteString() Token {
	position := l.position + 1 // пропускаем открывающую кавычку
	l.readChar()               // пропускаем открывающую кавычку

	for l.ch != '\'' && l.ch != 0 {
		// Обрабатываем экранированные символы
		if l.ch == '\\' {
			l.readChar() // пропускаем обратный слеш
		}
		l.readChar()
	}

	if l.ch == '\'' {
		l.readChar() // пропускаем закрывающую кавычку
	}

	return Token{
		Type:    STRING,
		Literal: l.input[position : l.position-1],
		Line:    l.line,
		Column:  l.column - (l.position - position),
	}
}

// readBacktickString читает строку в обратных кавычках
func (l *Lexer) readBacktickString() Token {
	position := l.position + 1 // пропускаем открывающую кавычку
	l.readChar()               // пропускаем открывающую кавычку

	for l.ch != '`' && l.ch != 0 {
		l.readChar()
	}

	if l.ch == '`' {
		l.readChar() // пропускаем закрывающую кавычку
	}

	return Token{
		Type:    STRING,
		Literal: l.input[position : l.position-1],
		Line:    l.line,
		Column:  l.column - (l.position - position),
	}
}

// readVariable читает переменную PHP ($var)
func (l *Lexer) readVariable() Token {
	position := l.position // позиция $
	l.readChar()           // пропускаем $

	// Читаем имя переменной
	for unicode.IsLetter(l.ch) || unicode.IsDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}

	return Token{
		Type:    VARIABLE,
		Literal: l.input[position:l.position],
		Line:    l.line,
		Column:  l.column - (l.position - position),
	}
}

// readSingleLineComment читает однострочный комментарий
func (l *Lexer) readSingleLineComment() string {
	position := l.position + 1 // пропускаем //
	l.readChar()               // пропускаем второй /

	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}

	return l.input[position:l.position]
}

// readMultiLineComment читает многострочный комментарий
func (l *Lexer) readMultiLineComment() string {
	position := l.position + 1 // пропускаем *
	l.readChar()               // пропускаем *

	for {
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // *
			l.readChar() // /
			break
		}
		if l.ch == 0 {
			break
		}
		l.readChar()
	}

	return l.input[position : l.position-2]
}

// readHeredoc читает heredoc/nowdoc конструкции
func (l *Lexer) readHeredoc() Token {
	// Это упрощенная реализация heredoc
	// В реальной реализации нужно более точно обрабатывать heredoc
	position := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return Token{
		Type:    HEREDOC_START,
		Literal: l.input[position:l.position],
		Line:    l.line,
		Column:  l.column - (l.position - position),
	}
}

// lookupIdent проверяет, является ли идентификатор ключевым словом
func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
