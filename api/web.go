package api

import (
	"errors"
	"net/http"

	"github.com/kautsarady/adindopustaka/httputil"
	"github.com/kautsarady/adindopustaka/model"

	"github.com/gin-gonic/gin"
)

// PageLanding .
func (ctr *Controller) PageLanding(ctx *gin.Context) {
	limit, offset, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	books, err := ctr.DAO.Get("books", nil, limit, offset)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, errors.New("database query failure"))
		ctx.Error(err)
		return
	}

	if books == nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("no corresponding data found"))
		return
	}

	ctx.HTML(http.StatusOK, "index.html", wrapData("all", limit, offset, books))
}

// PageBook .
func (ctr *Controller) PageBook(ctx *gin.Context) {
	id := ctx.Param("id")

	book, err := ctr.DAO.GetBookByID(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, errors.New("database query failure"))
		ctx.Error(err)
		return
	}

	if book == nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("no corresponding data found"))
		return
	}

	ctx.HTML(http.StatusOK, "detail.html", book)
}

// PageAuthor .
func (ctr *Controller) PageAuthor(ctx *gin.Context) {
	id := ctx.Param("id")

	limit, offset, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	author, err := ctr.DAO.GetItemByID("authors", id, limit, offset)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, errors.New("database query failure"))
		ctx.Error(err)
		return
	}

	if author == nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("no corresponding data found"))
		return
	}

	ctx.HTML(http.StatusOK, "entity.html", wrapData("author/"+id, limit, offset, author))
}

// PageCategory .
func (ctr *Controller) PageCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	limit, offset, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	categories, err := ctr.DAO.GetItemByID("categories", id, limit, offset)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, errors.New("database query failure"))
		ctx.Error(err)
		return
	}

	if categories == nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("no corresponding data found"))
		return
	}

	ctx.HTML(http.StatusOK, "entity.html", wrapData("category/"+id, limit, offset, categories))
}

// PageTag .
func (ctr *Controller) PageTag(ctx *gin.Context) {
	id := ctx.Param("id")

	limit, offset, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	tags, err := ctr.DAO.GetItemByID("tags", id, limit, offset)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, errors.New("database query failure"))
		ctx.Error(err)
		return
	}

	if tags == nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("no corresponding data found"))
		return
	}

	ctx.HTML(http.StatusOK, "entity.html", wrapData("tag/"+id, limit, offset, tags))
}

// PageFilter .
func (ctr *Controller) PageFilter(ctx *gin.Context) {
	limit, offset, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	authors, err := ctr.DAO.GetDistinctItems("authors", limit, offset)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, errors.New("database query failure"))
		ctx.Error(err)
		return
	}

	if authors == nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("no corresponding data found"))
		return
	}

	categories, err := ctr.DAO.GetDistinctItems("categories", limit, offset)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, errors.New("database query failure"))
		ctx.Error(err)
		return
	}

	if categories == nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("no corresponding data found"))
		return
	}

	tags, err := ctr.DAO.GetDistinctItems("tags", limit, offset)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, errors.New("database query failure"))
		ctx.Error(err)
		return
	}

	if tags == nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("no corresponding data found"))
		return
	}

	data := struct {
		Authors    []model.Item
		Categories []model.Item
		Tags       []model.Item
	}{
		Authors:    model.ToItems(authors),
		Categories: model.ToItems(categories),
		Tags:       model.ToItems(tags),
	}

	ctx.HTML(http.StatusOK, "filter.html", wrapData("ignore", limit, offset, data))
}
