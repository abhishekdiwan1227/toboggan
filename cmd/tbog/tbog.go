package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/abhishekdiwan1227/toboggan/internal/ast"
)

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	for {
		if stdin.Scan() {
			line := stdin.Text()

			parser := ast.InitParser(line)

			tree := parser.Parse()

			print(tree, "")
		}
	}
}

func print(root ast.SyntaxNode, indent string) {
	fmt.Print(indent)
	if root.TokenValue() != nil {
		fmt.Print(root.TokenValue())
	} else {
		fmt.Print(root.Kind().String())
	}
	fmt.Println()
	indent += "    "
	for _, child := range root.Children() {
		print(child, indent)
	}
}
