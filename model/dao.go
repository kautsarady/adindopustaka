package model

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/kautsarady/adindopustaka/model/item"
)

// DAO .
type DAO struct{ DB *sql.DB }

// Make .
func Make(cs string) (*DAO, error) {
	db, err := sql.Open("mysql", cs)
	if err != nil {
		return nil, err
	}
	return &DAO{db}, nil
}

// GetAllAuthors .
func (c *DAO) GetAllAuthors(offset, limit int) ([]item.Author, error) {

	rows, err := c.DB.Query("SELECT DISTINCT id, name FROM authors ORDER BY name LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	var authors []item.Author

	for rows.Next() {
		var a item.Author
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, err
		}
		a.Name = strings.Title(a.Name)
		authors = append(authors, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return authors, nil
}

// GetAllCategories .
func (c *DAO) GetAllCategories(offset, limit int) ([]item.Category, error) {

	rows, err := c.DB.Query("SELECT DISTINCT id, name FROM categories ORDER BY name LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	var categories []item.Category

	for rows.Next() {
		var a item.Category
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, err
		}
		a.Name = strings.Title(a.Name)
		categories = append(categories, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

// GetAllTags .
func (c *DAO) GetAllTags(offset, limit int) ([]item.Tag, error) {

	rows, err := c.DB.Query("SELECT DISTINCT id, name FROM tags ORDER BY name LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	var tags []item.Tag

	for rows.Next() {
		var a item.Tag
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, err
		}
		a.Name = strings.Title(a.Name)
		tags = append(tags, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

// GetAllBooks .
func (c *DAO) GetAllBooks(offset, limit int) ([]Book, error) {

	rows, err := c.DB.Query("SELECT * FROM books ORDER BY id LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}

	books, err := c.handleBooksRows(rows)
	if err != nil {
		return nil, err
	}

	return books, nil
}

// GetFilterBooks .
func (c *DAO) GetFilterBooks(filter, query string, offset, limit int) ([]Book, error) {

	bookIDs, err := c.getBookIDsBy(filter, query)
	if err != nil {
		return nil, err
	}

	rows, err := c.DB.Query(fmt.Sprintf("SELECT * FROM books WHERE id IN(%s) ORDER BY id LIMIT ? OFFSET ?", strings.Join(bookIDs, ", ")), limit, offset)
	if err != nil {
		return nil, err
	}

	books, err := c.handleBooksRows(rows)
	if err != nil {
		return nil, err
	}

	return books, nil
}

// GetBookByID .
func (c *DAO) GetBookByID(id int) (Book, error) {

	var b Book
	if err := c.DB.QueryRow("SELECT * FROM books WHERE id = ?", id).
		Scan(&b.ID, &b.Title, &b.ImageURL, &b.GramedURL, &b.Description); err != nil {
		return Book{}, err
	}

	if err := c.setBookItems(&b); err != nil {
		return Book{}, err
	}

	return b, nil
}

// AbsItem .
type AbsItem struct {
	ID    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Books []int  `json:"books,omitempty"`
}

// CompAbs .
type CompAbs struct {
	Item  AbsItem `json:"item,omitempty"`
	Books []Book  `json:"books,omitempty"`
}

// GetItemByID .
func (c *DAO) GetItemByID(filter string, id int) (i AbsItem, err error) {
	err = c.DB.QueryRow(fmt.Sprintf("SELECT id, name FROM %s WHERE id = ?", filter), id).Scan(&i.ID, &i.Name)
	i.Name = strings.Title(i.Name)
	return
}

func (c *DAO) handleBooksRows(rows *sql.Rows) ([]Book, error) {
	var books []Book

	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.ImageURL, &b.GramedURL, &b.Description); err != nil {
			return nil, err
		}

		// if err := c.setBookItems(&b); err != nil {
		// 	return nil, err
		// }

		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (c *DAO) setBookItems(b *Book) error {

	tt := []string{"authors", "categories", "tags"}

	for _, tn := range tt {
		rows, err := c.DB.Query(fmt.Sprintf("SELECT id, name FROM %s WHERE book_id = ?", tn), b.ID)
		if err != nil {
			return err
		}

		for rows.Next() {
			var i AbsItem
			if err := rows.Scan(&i.ID, &i.Name); err != nil {
				return err
			}
			i.Name = strings.Title(i.Name)
			switch tn {
			case "authors":
				b.Authors = append(b.Authors, item.Author(i))
				break
			case "categories":
				b.Categories = append(b.Categories, item.Category(i))
				break
			case "tags":
				b.Tags = append(b.Tags, item.Tag(i))
				break
			}
		}

		if err := rows.Err(); err != nil {
			return err
		}
	}

	return nil
}

func (c *DAO) getBookIDsBy(tname, query string) ([]string, error) {

	rows, err := c.DB.Query(fmt.Sprintf("SELECT book_id FROM %s WHERE name = ?", tname), query)
	if err != nil {
		return nil, err
	}

	var bookIDs []string

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		bookIDs = append(bookIDs, strconv.Itoa(id))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return bookIDs, nil
}
