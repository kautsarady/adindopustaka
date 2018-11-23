//go:generate gorunpkg github.com/99designs/gqlgen

package gql

import (
	context "context"

	model "github.com/kautsarady/adindopustaka/model"
	item "github.com/kautsarady/adindopustaka/model/item"
)

type Resolver struct{ DAO *model.DAO }

func (r *Resolver) Author() AuthorResolver {
	return &authorResolver{r}
}
func (r *Resolver) Category() CategoryResolver {
	return &categoryResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Tag() TagResolver {
	return &tagResolver{r}
}

type authorResolver struct{ *Resolver }

func (r *authorResolver) Books(ctx context.Context, obj *item.Author, page int, per_page int) ([]model.Book, error) {
	offset, limit := (page-1)*per_page, per_page
	books, err := r.DAO.GetFilterBooks("authors", obj.Name, offset, limit)
	if err != nil {
		return nil, err
	}
	return books, nil
}

type categoryResolver struct{ *Resolver }

func (r *categoryResolver) Books(ctx context.Context, obj *item.Category, page int, per_page int) ([]model.Book, error) {
	offset, limit := (page-1)*per_page, per_page
	books, err := r.DAO.GetFilterBooks("categories", obj.Name, offset, limit)
	if err != nil {
		return nil, err
	}
	return books, nil
}

type tagResolver struct{ *Resolver }

func (r *tagResolver) Books(ctx context.Context, obj *item.Tag, page int, per_page int) ([]model.Book, error) {
	offset, limit := (page-1)*per_page, per_page
	books, err := r.DAO.GetFilterBooks("tags", obj.Name, offset, limit)
	if err != nil {
		return nil, err
	}
	return books, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Book(ctx context.Context, id int) (model.Book, error) {
	return r.DAO.GetBookByID(id)
}

func (r *queryResolver) Author(ctx context.Context, id int) (item.Author, error) {
	abs, err := r.DAO.GetItemByID("authors", id)
	if err != nil {
		return item.Author{}, err
	}
	return item.Author(abs), nil
}

func (r *queryResolver) Category(ctx context.Context, id int) (item.Category, error) {
	abs, err := r.DAO.GetItemByID("categories", id)
	if err != nil {
		return item.Category{}, err
	}
	return item.Category(abs), nil
}

func (r *queryResolver) Tag(ctx context.Context, id int) (item.Tag, error) {
	abs, err := r.DAO.GetItemByID("tags", id)
	if err != nil {
		return item.Tag{}, err
	}
	return item.Tag(abs), nil
}

func (r *queryResolver) AllBooks(ctx context.Context, page int, per_page int) ([]model.Book, error) {
	offset, limit := (page-1)*per_page, per_page
	books, err := r.DAO.GetAllBooks(offset, limit)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (r *queryResolver) AllAuthors(ctx context.Context, page int, per_page int) ([]item.Author, error) {
	offset, limit := (page-1)*per_page, per_page
	authors, err := r.DAO.GetAllAuthors(offset, limit)
	if err != nil {
		return nil, err
	}
	return authors, nil
}

func (r *queryResolver) AllCategories(ctx context.Context, page int, per_page int) ([]item.Category, error) {
	offset, limit := (page-1)*per_page, per_page
	categories, err := r.DAO.GetAllCategories(offset, limit)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *queryResolver) AllTags(ctx context.Context, page int, per_page int) ([]item.Tag, error) {
	offset, limit := (page-1)*per_page, per_page
	tags, err := r.DAO.GetAllTags(offset, limit)
	if err != nil {
		return nil, err
	}
	return tags, nil
}
