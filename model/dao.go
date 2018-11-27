package model

import (
	"database/sql"
	"fmt"
	"strings"
)

// DAO .
type DAO struct{ DB *sql.DB }

// Make .
func Make(connStr string) (*DAO, error) {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	return &DAO{db}, nil
}

// Get .
func (d *DAO) Get(entity string, filter []string, limit, offset int) ([]interface{}, error) {

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT ? OFFSET ?", entity, whereValue(filter))
	rows, err := d.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []interface{}
	switch entity {
	case "books":
		if err := handleBooks(&result, rows); err != nil {
			return nil, err
		}
		break
	default:
		if err := handleItems(&result, rows); err != nil {
			return nil, err
		}
		break
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// GetDistinctItems .
func (d *DAO) GetDistinctItems(entity string, limit, offset int) ([]interface{}, error) {

	query := fmt.Sprintf("SELECT * FROM %s GROUP BY id LIMIT ? OFFSET ?", entity)
	rows, err := d.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []interface{}
	if err := handleItems(&items, rows); err != nil {
		return nil, err
	}

	return items, nil
}

// GetBookByID .
func (d *DAO) GetBookByID(id string) (*Book, error) {

	books, err := d.Get("books", []string{fmt.Sprintf("id = %s", id)}, 1, 0)
	if err != nil {
		return nil, err
	}

	authors, err := d.Get("authors", []string{fmt.Sprintf("book_id = %s", id)}, 99, 0)
	if err != nil {
		return nil, err
	}

	categories, err := d.Get("categories", []string{fmt.Sprintf("book_id = %s", id)}, 99, 0)
	if err != nil {
		return nil, err
	}

	tags, err := d.Get("tags", []string{fmt.Sprintf("book_id = %s", id)}, 99, 0)
	if err != nil {
		return nil, err
	}

	book := books[0].(Book)
	book.Authors = ToItems(authors)
	book.Categories = ToItems(categories)
	book.Tags = ToItems(tags)

	return &book, nil
}

// GetItemByID .
func (d *DAO) GetItemByID(entity string, id string, limit, offset int) (*Item, error) {

	result, err := d.Get(entity, []string{fmt.Sprintf("id = %s", id)}, 99, 0)
	if err != nil || len(result) == 0 {
		return nil, err
	}

	item, bookIDs := ItemAndIDs(result)

	books, err := d.Get("books", []string{fmt.Sprintf("id IN(%s)", strings.Join(bookIDs, ", "))}, limit, offset)
	if err != nil {
		return nil, err
	}

	item.Books = ToBooks(books)

	return &item, nil
}
