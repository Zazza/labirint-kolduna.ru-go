//go:build ignore
// +build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 user= password= dbname=gamebook sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Получаем все переходы
	rows, err := db.Query(`
		SELECT s.number, s2.number as target_number 
		FROM transitions t 
		JOIN sections s ON t.section_id = s.id 
		JOIN sections s2 ON t.target_section_id = s2.id 
		WHERE s.number != 9 AND s2.number != 9
		ORDER BY s.number, t.text_order
	`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Граф: map[номер_секции][]номеров_переходов
	graph := make(map[uint][]uint)
	numbers := make(map[uint]bool)

	for rows.Next() {
		var from, to uint
		if err := rows.Scan(&from, &to); err != nil {
			log.Fatal(err)
		}
		graph[from] = append(graph[from], to)
		numbers[from] = true
		numbers[to] = true
	}

	// Сортируем номера секций
	var sortedNumbers []uint
	for n := range numbers {
		sortedNumbers = append(sortedNumbers, n)
	}
	sort.Slice(sortedNumbers, func(i, j int) bool { return sortedNumbers[i] < sortedNumbers[j] })

	// Генерируем DOT-файл
	var sb strings.Builder
	sb.WriteString("digraph gamebook {\n")
	sb.WriteString("  rankdir=LR;\n")
	sb.WriteString("  node [shape=box, style=filled, fillcolor=lightblue];\n\n")

	// Записываем все узлы
	for _, n := range sortedNumbers {
		sb.WriteString(fmt.Sprintf("  %d [label=\"Секция %d\"];\n", n, n))
	}
	sb.WriteString("\n")

	// Записываем рёбра
	for _, from := range sortedNumbers {
		targets := graph[from]
		for _, to := range targets {
			sb.WriteString(fmt.Sprintf("  %d -> %d;\n", from, to))
		}
	}

	sb.WriteString("}\n")

	// Записываем в файл
	err = os.WriteFile("graph.dot", []byte(sb.String()), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Граф записан в graph.dot")
	fmt.Printf("Всего узлов: %d\n", len(numbers))
	fmt.Printf("Всего переходов: %d\n", len(graph))
}
