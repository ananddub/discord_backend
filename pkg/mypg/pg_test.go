package mypg

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/sys/windows"
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

func TestGoRoutine(t *testing.T) {

	// Check available resources
	fmt.Println("CPU cores:", runtime.NumCPU())
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))

	var wg sync.WaitGroup

	// Launch 10000 goroutines (light!)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// This goroutine might run on ANY M+P combination
			fmt.Printf("G%d on thread %d\n", id, windows.GetCurrentThreadId())
		}(i)
	}

	// Only 4-8 OS threads (M) handling 10000 goroutines (G)
	fmt.Println("OS threads:", runtime.NumGoroutine())

	wg.Wait()
}
func tmp(i int) {
	fmt.Println(i)
	go goRoutine()
	tmp(i + 1)
}
func goRoutine() {
	for {
	}
}

func TestRec(t *testing.T) {
	tmp(0)
	time.Sleep(10 * time.Minute)
}
