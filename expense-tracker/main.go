package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Expense struct {
	ID          int     `json:"id"`
	Date        string  `json:"date"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category,omitempty"`
}

const dataFile = "expenses.json"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		cmdAdd(os.Args[2:])
	case "update":
		cmdUpdate(os.Args[2:])
	case "delete":
		cmdDelete(os.Args[2:])
	case "list":
		cmdList(os.Args[2:])
	case "summary":
		cmdSummary(os.Args[2:])
	case "export":
		cmdExport(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Perintah tidak dikenal: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Expense Tracker CLI

Penggunaan:
  expense-tracker add --description "Lunch" --amount 20 [--category Food]
  expense-tracker update --id 1 --description "New" --amount 15 [--category Food]
  expense-tracker delete --id 1
  expense-tracker list [--category Food]
  expense-tracker summary [--month 8]
  expense-tracker export [--file expenses.csv]`)
}

func dataFilePath() string {
	return dataFile
}

func loadExpenses() ([]Expense, error) {
	path := dataFilePath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []Expense{}, nil
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return []Expense{}, nil
	}
	var expenses []Expense
	if err := json.Unmarshal(bytes, &expenses); err != nil {
		return nil, err
	}
	return expenses, nil
}

func saveExpenses(expenses []Expense) error {
	bytes, err := json.MarshalIndent(expenses, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFilePath(), bytes, 0644)
}

func nextID(expenses []Expense) int {
	max := 0
	for _, e := range expenses {
		if e.ID > max {
			max = e.ID
		}
	}
	return max + 1
}

func cmdAdd(args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	description := fs.String("description", "", "Deskripsi pengeluaran (wajib)")
	amount := fs.Float64("amount", 0, "Jumlah pengeluaran (wajib, harus > 0)")
	category := fs.String("category", "", "Kategori pengeluaran (opsional)")
	fs.Parse(args)

	if *description == "" {
		fmt.Println("Error: --description tidak boleh kosong")
		os.Exit(1)
	}
	if *amount <= 0 {
		fmt.Println("Error: --amount harus lebih besar dari 0")
		os.Exit(1)
	}

	expenses, err := loadExpenses()
	if err != nil {
		fmt.Println("Error membaca data:", err)
		os.Exit(1)
	}

	newExpense := Expense{
		ID:          nextID(expenses),
		Date:        time.Now().Format("2006-01-02"),
		Description: *description,
		Amount:      *amount,
		Category:    *category,
	}
	expenses = append(expenses, newExpense)

	if err := saveExpenses(expenses); err != nil {
		fmt.Println("Error menyimpan data:", err)
		os.Exit(1)
	}

	fmt.Printf("Expense added successfully (ID: %d)\n", newExpense.ID)
}

func cmdUpdate(args []string) {
	fs := flag.NewFlagSet("update", flag.ExitOnError)
	id := fs.Int("id", 0, "ID pengeluaran yang ingin diupdate (wajib)")
	description := fs.String("description", "", "Deskripsi baru (opsional)")
	amount := fs.Float64("amount", -1, "Jumlah baru (opsional)")
	category := fs.String("category", "", "Kategori baru (opsional)")
	fs.Parse(args)

	if *id <= 0 {
		fmt.Println("Error: --id wajib diisi dan harus lebih besar dari 0")
		os.Exit(1)
	}

	expenses, err := loadExpenses()
	if err != nil {
		fmt.Println("Error membaca data:", err)
		os.Exit(1)
	}

	found := false
	for i, e := range expenses {
		if e.ID == *id {
			found = true
			if *description != "" {
				expenses[i].Description = *description
			}
			if *amount >= 0 {
				expenses[i].Amount = *amount
			}
			if *category != "" {
				expenses[i].Category = *category
			}
			break
		}
	}

	if !found {
		fmt.Printf("Error: expense dengan ID %d tidak ditemukan\n", *id)
		os.Exit(1)
	}

	if err := saveExpenses(expenses); err != nil {
		fmt.Println("Error menyimpan data:", err)
		os.Exit(1)
	}

	fmt.Println("Expense updated successfully")
}

func cmdDelete(args []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	id := fs.Int("id", 0, "ID pengeluaran yang ingin dihapus (wajib)")
	fs.Parse(args)

	if *id <= 0 {
		fmt.Println("Error: --id wajib diisi dan harus lebih besar dari 0")
		os.Exit(1)
	}

	expenses, err := loadExpenses()
	if err != nil {
		fmt.Println("Error membaca data:", err)
		os.Exit(1)
	}

	newExpenses := make([]Expense, 0, len(expenses))
	found := false
	for _, e := range expenses {
		if e.ID == *id {
			found = true
			continue
		}
		newExpenses = append(newExpenses, e)
	}

	if !found {
		fmt.Printf("Error: expense dengan ID %d tidak ditemukan\n", *id)
		os.Exit(1)
	}

	if err := saveExpenses(newExpenses); err != nil {
		fmt.Println("Error menyimpan data:", err)
		os.Exit(1)
	}

	fmt.Println("Expense deleted successfully")
}

func cmdList(args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	category := fs.String("category", "", "Filter berdasarkan kategori (opsional)")
	fs.Parse(args)

	expenses, err := loadExpenses()
	if err != nil {
		fmt.Println("Error membaca data:", err)
		os.Exit(1)
	}

	if len(expenses) == 0 {
		fmt.Println("Belum ada expense yang tercatat.")
		return
	}

	fmt.Printf("%-4s%-12s%-20s%-10s%-12s\n", "ID", "Date", "Description", "Amount", "Category")
	for _, e := range expenses {
		if *category != "" && e.Category != *category {
			continue
		}
		fmt.Printf("%-4d%-12s%-20s$%-9.2f%-12s\n", e.ID, e.Date, e.Description, e.Amount, e.Category)
	}
}

func cmdSummary(args []string) {
	fs := flag.NewFlagSet("summary", flag.ExitOnError)
	month := fs.Int("month", 0, "Filter bulan pada tahun berjalan, 1-12 (opsional)")
	category := fs.String("category", "", "Filter berdasarkan kategori (opsional)")
	fs.Parse(args)

	if *month < 0 || *month > 12 {
		fmt.Println("Error: --month harus antara 1 dan 12")
		os.Exit(1)
	}

	expenses, err := loadExpenses()
	if err != nil {
		fmt.Println("Error membaca data:", err)
		os.Exit(1)
	}

	currentYear := time.Now().Year()
	total := 0.0

	for _, e := range expenses {
		if *category != "" && e.Category != *category {
			continue
		}
		if *month != 0 {
			d, err := time.Parse("2006-01-02", e.Date)
			if err != nil {
				continue
			}
			if d.Year() != currentYear || int(d.Month()) != *month {
				continue
			}
		}
		total += e.Amount
	}

	if *month != 0 {
		fmt.Printf("Total expenses for %s: $%.2f\n", time.Month(*month).String(), total)
	} else {
		fmt.Printf("Total expenses: $%.2f\n", total)
	}
}

func cmdExport(args []string) {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	outFile := fs.String("file", "expenses.csv", "Nama file CSV tujuan")
	fs.Parse(args)

	expenses, err := loadExpenses()
	if err != nil {
		fmt.Println("Error membaca data:", err)
		os.Exit(1)
	}

	f, err := os.Create(filepath.Clean(*outFile))
	if err != nil {
		fmt.Println("Error membuat file CSV:", err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintln(f, "ID,Date,Description,Amount,Category")
	for _, e := range expenses {
		fmt.Fprintf(f, "%d,%s,%s,%.2f,%s\n", e.ID, e.Date, e.Description, e.Amount, e.Category)
	}

	fmt.Printf("Expenses exported successfully to %s\n", *outFile)
}
