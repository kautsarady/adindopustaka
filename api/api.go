package api

import (
	"errors"
	"net/http"

	"github.com/kautsarady/adindopustaka/httputil"
	"github.com/kautsarady/adindopustaka/model"

	// doc.json
	_ "github.com/kautsarady/adindopustaka/docs"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Controller .
type Controller struct {
	DAO    *model.DAO
	Router *gin.Engine
}

// Make .
func Make(dao *model.DAO) *Controller {
	ctr := &Controller{dao, gin.Default()}
	ctr.Router.LoadHTMLGlob("public/*")
	ctr.Router.GET("/", ctr.PageLanding)
	ctr.Router.GET("/filter", ctr.PageFilter)
	ctr.Router.GET("/book/:id", ctr.PageBook)
	ctr.Router.GET("/author/:id", ctr.PageAuthor)
	ctr.Router.GET("/category/:id", ctr.PageCategory)
	ctr.Router.GET("/tag/:id", ctr.PageTag)
	api := ctr.Router.Group("/api")
	{
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		api.GET("/book", ctr.GetAllBook)
		api.GET("/author", ctr.GetAllAuthor)
		api.GET("/category", ctr.GetAllCategory)
		api.GET("/tag", ctr.GetAllTag)
		api.GET("/book/:id", ctr.GetBook)
		api.GET("/author/:id", ctr.GetAuthor)
		api.GET("/category/:id", ctr.GetCategory)
		api.GET("/tag/:id", ctr.GetTag)
	}
	return ctr
}

// @title github.com/kautsarady/Adindopustaka API
// @version 1.0
// @description github.com/kautsarady/Adindopustaka API documentation
// @contact.name kautsarady
// @contact.email kautsarady@gmail.com

// GetAllBook godoc
// @Summary Get All Book
// @ID get-all-book
// @Accept json
// @Produce json
// @Param page query string false "page number (default=1)" Format(string)
// @Param per_page query string false "per_page product count (default=20)" Format(string)
// @Success 200 {array} model.Book
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Router /api/book [get]
func (ctr *Controller) GetAllBook(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, books)
}

// GetAllAuthor godoc
// @Summary Get All Author
// @ID get-all-author
// @Accept json
// @Produce json
// @Param page query string false "page number (default=1)" Format(string)
// @Param per_page query string false "per_page product count (default=20)" Format(string)
// @Success 200 {array} model.Item
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Router /api/author [get]
func (ctr *Controller) GetAllAuthor(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, authors)
}

// GetAllCategory godoc
// @Summary Get All Category
// @ID get-all-category
// @Accept json
// @Produce json
// @Param page query string false "page number (default=1)" Format(string)
// @Param per_page query string false "per_page product count (default=20)" Format(string)
// @Success 200 {array} model.Item
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Router /api/category [get]
func (ctr *Controller) GetAllCategory(ctx *gin.Context) {
	limit, offset, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
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

	ctx.JSON(http.StatusOK, categories)
}

// GetAllTag godoc
// @Summary Get All Tag
// @ID get-all-tag
// @Accept json
// @Produce json
// @Param page query string false "page number (default=1)" Format(string)
// @Param per_page query string false "per_page product count (default=20)" Format(string)
// @Success 200 {array} model.Item
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Router /api/tag [get]
func (ctr *Controller) GetAllTag(ctx *gin.Context) {
	limit, offset, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
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

	ctx.JSON(http.StatusOK, tags)
}

// GetBook godoc
// @Summary Get Book By ID
// @ID get-all-getBookByID
// @Accept json
// @Produce json
// @Param id path string true "book id to search"
// @Success 200 {object} model.Book
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Router /api/book/{id} [get]
func (ctr *Controller) GetBook(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, book)
}

// GetAuthor godoc
// @Summary Get Author By ID
// @ID get-all-getAuthorByID
// @Accept json
// @Produce json
// @Param id path string true "author id to search"
// @Param page query string false "page number of the item books (default=1)" Format(string)
// @Param per_page query string false "per_page product count of the item books (default=20)" Format(string)
// @Success 200 {object} model.Item
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Router /api/author/{id} [get]
func (ctr *Controller) GetAuthor(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, author)
}

// GetCategory godoc
// @Summary Get Category By ID
// @ID get-all-getCategoryByID
// @Accept json
// @Produce json
// @Param id path string true "category id to search"
// @Param page query string false "page number of the item books (default=1)" Format(string)
// @Param per_page query string false "per_page product count of the item books (default=20)" Format(string)
// @Success 200 {object} model.Item
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Router /api/category/{id} [get]
func (ctr *Controller) GetCategory(ctx *gin.Context) {
	id := ctx.Param("id")

	limit, offset, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	category, err := ctr.DAO.GetItemByID("categories", id, limit, offset)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, errors.New("database query failure"))
		ctx.Error(err)
		return
	}

	if category == nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("no corresponding data found"))
		return
	}

	ctx.JSON(http.StatusOK, category)
}

// GetTag godoc
// @Summary Get Tag By ID
// @ID get-all-getTagByID
// @Accept json
// @Produce json
// @Param id path string true "tag id to search"
// @Param page query string false "page number of the item books (default=1)" Format(string)
// @Param per_page query string false "per_page product count of the item books (default=20)" Format(string)
// @Success 200 {object} model.Item
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Router /api/tag/{id} [get]
func (ctr *Controller) GetTag(ctx *gin.Context) {
	id := ctx.Param("id")

	limit, offset, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, err)
		return
	}

	tag, err := ctr.DAO.GetItemByID("tags", id, limit, offset)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, errors.New("database query failure"))
		ctx.Error(err)
		return
	}

	if tag == nil {
		httputil.NewError(ctx, http.StatusNotFound, errors.New("no corresponding data found"))
		return
	}

	ctx.JSON(http.StatusOK, tag)
}
