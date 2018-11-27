package model

import (
	"database/sql"
	"strconv"
	"strings"
)

// ToBooks .
func ToBooks(result []interface{}) (books []Book) {
	for _, book := range result {
		books = append(books, book.(Book))
	}
	return
}

// ToItems .
func ToItems(result []interface{}) (items []Item) {
	for _, item := range result {
		items = append(items, item.(Item))
	}
	return
}

// ItemAndIDs .
func ItemAndIDs(result []interface{}) (Item, []string) {
	items := ToItems(result)
	var IDs []string
	for _, item := range items {
		IDs = append(IDs, strconv.Itoa(item.BookID))
	}
	return items[0], IDs
}

func whereValue(filter []string) string {
	if filter == nil {
		return "1"
	}
	return strings.Join(filter, " AND ")
}

func handleBooks(result *[]interface{}, rows *sql.Rows) error {
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.ImageURL, &book.GramedURL, &book.Description); err != nil {
			return err
		}
		*result = append(*result, book)
	}
	return nil
}

func handleItems(result *[]interface{}, rows *sql.Rows) error {
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.BookID, &item.Name); err != nil {
			return err
		}
		item.Name = strings.Title(item.Name)
		*result = append(*result, item)
	}
	return nil
}
