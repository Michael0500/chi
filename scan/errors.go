package scan

import "fmt"

// LexerError представляет ошибку лексического анализа
type LexerError struct {
    Message string
    Line    int
    Column  int
}

// Error возвращает строковое представление ошибки
func (e *LexerError) Error() string {
    return fmt.Sprintf("lexer error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

// NewLexerError создает новую ошибку лексического анализа
func NewLexerError(message string, line, column int) *LexerError {
    return &LexerError{
        Message: message,
        Line:    line,
        Column:  column,
    }
}