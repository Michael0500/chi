package main

import (
	"fmt"

	"chi-lang/scan"
)

func main() {
	// Тестовый код PHP
	input := `<?php
class MyClass {
    private $property = "Hello World";
    
    public function myMethod($param1, $param2 = null) {
        if ($param1 > 0) {
            return $param1 + $param2;
        } elseif ($param1 == 0) {
            echo "Zero value";
        } else {
            $array = [1, 2, 3];
            foreach ($array as $item) {
                echo $item . "\n";
            }
        }
        
        try {
            throw new Exception("Test exception");
        } catch (Exception $e) {
            echo $e->getMessage();
        } finally {
            echo "Finally block";
        }
    }
}

$object = new MyClass();
$result = $object->myMethod(5, 10);
echo $result;
`

	l := scan.New(input)

	fmt.Println("Лексический анализ PHP кода:")
	fmt.Println("============================")

	tokenCount := 0
	for {
		tok := l.NextToken()
		tokenCount++

		// Форматируем вывод
		lineInfo := fmt.Sprintf("[%d:%d]", tok.Line, tok.Column)
		typeInfo := fmt.Sprintf("%-15s", tok.Type.String())
		literalInfo := fmt.Sprintf("'%s'", tok.Literal)

		fmt.Printf("%-10s %-15s %s\n", lineInfo, typeInfo, literalInfo)

		if tok.Type == scan.EOF || tokenCount > 1000 { // Защита от бесконечного цикла
			break
		}
	}

	fmt.Printf("\nВсего токенов: %d\n", tokenCount-1)
}
