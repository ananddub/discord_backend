package mypg

import (
	"fmt"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []string{
		"UPDATE users SET name='John', age=25 WHERE id = 10 AND status = 'active'",
		"DELETE FROM orders WHERE price > 100",
		"SELECT id,name FROM product WHERE cat = 'food'",
		"INSERT INTO users(id,name) VALUES(1,'AA')",
		"Update user set email = 'test@gmail.com' where id =5 and is_accessed=true",
	}

	for _, q := range tests {
		fmt.Println("\nQUERY:", q)
		r, _ := Parser(q)

		fmt.Println("TYPE:", r.Type)
		fmt.Println("TABLE:", r.Table)
		fmt.Println("COLUMNS:", r.Columns)
		fmt.Println("WHERE PAIRS:", r.WherePairs)
		fmt.Println("SELECT COLS:", r.SelectCols)
	}
}

func TestInterpolate(t *testing.T) {
	query := "UPDATE users SET name=$1, age=$2 WHERE id = $3 AND status = $4"
	args := []interface{}{"John", 25, 10, "active"}
	result := Interpolate(query, args...)
	expected := "UPDATE users SET name='John', age=25 WHERE id = 10 AND status = 'active'"
	if strings.TrimSpace(result) != strings.TrimSpace(expected) {
		t.Errorf("Expected: %s, Got: %s", expected, result)
	}
	fmt.Println("Interpolated Query:", result)
}
