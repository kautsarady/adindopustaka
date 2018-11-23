package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/handler"
	"github.com/gin-gonic/gin"

	// doc.json
	_ "github.com/kautsarady/adindopustaka/docs"
	"github.com/kautsarady/adindopustaka/gql"
	"github.com/kautsarady/adindopustaka/httputil"
	"github.com/kautsarady/adindopustaka/model"
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
	api := ctr.Router.Group("/api")
	{
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		api.GET("/book", ctr.HandleBooks)
		api.GET("/author", ctr.HandleAuthors)
		api.GET("/category", ctr.HandleCategories)
		api.GET("/tag", ctr.HandleTags)
		api.GET("/book/:id", ctr.HandleBook)
		api.GET("/author/:id", ctr.HandleAuthor)
		api.GET("/category/:id", ctr.HandleCategory)
		api.GET("/tag/:id", ctr.HandleTag)
	}
	gqlAPI := ctr.Router.Group("/gql")
	gqlHandlerFunc := gin.WrapF(handler.GraphQL(gql.NewExecutableSchema(gql.Config{Resolvers: &gql.Resolver{DAO: dao}})))
	{
		gqlAPI.GET("/", gin.WrapF(handler.Playground("Adindopustaka GraphQL", "/gql/query")))
		gqlAPI.GET("/query", gqlHandlerFunc)
		gqlAPI.POST("/query", gqlHandlerFunc)
		gqlAPI.OPTIONS("/query", gqlHandlerFunc)
	}
	return ctr
}

// @title Adindopustaka API
// @version 1.0
// @description Adindopustaka API documentation
// @contact.name kautsarady
// @contact.email kautsarady@gmail.com

// HandleBooks godoc
// @Summary Get All Books
// @ID get-all-books
// @Accept json
// @Produce json
// @Param page query string false "page number (default=1)" Format(string)
// @Param per_page query string false "per_page product count (default=20)" Format(string)
// @Param filter query string false "search filter ('authors', 'categories', 'tags', default='all')"
// @Param query query string false "search query ('tere liye', 'agama', 'smp', default='')"
// @Success 200 {array} model.Book
// @Failure 400 {object} httputil.HTTPError
// @Router /api/book [get]
func (ctr *Controller) HandleBooks(ctx *gin.Context) {

	offset, limit, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("query string 'page' AND 'per_page' MUST be a VALID NUMBER"))
		return
	}

	query := ctx.DefaultQuery("query", "")
	filter := ctx.DefaultQuery("filter", "all")

	if filter == "all" {
		books, err := ctr.DAO.GetAllBooks(offset, limit)
		if err != nil {
			httputil.NewError(ctx, http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, books)
		return
	}

	if query == "" {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("query string 'query' NOT FOUND"))
		return
	}

	books, err := ctr.DAO.GetFilterBooks(filter, query, offset, limit)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	if books == nil {
		httputil.NewError(ctx, http.StatusSeeOther, errors.New("no corresponding result found"))
	}

	ctx.JSON(http.StatusOK, books)
}

// HandleAuthors godoc
// @Summary Get All Authors
// @ID get-all-authors
// @Accept json
// @Produce json
// @Param page query string false "page number (default=1)" Format(string)
// @Param per_page query string false "per_page product count (default=20)" Format(string)
// @Success 200 {array} item.Author
// @Failure 400 {object} httputil.HTTPError
// @Router /api/author [get]
func (ctr *Controller) HandleAuthors(ctx *gin.Context) {

	offset, limit, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("query string 'page' AND 'per_page' MUST be a VALID NUMBER"))
		return
	}

	authors, err := ctr.DAO.GetAllAuthors(offset, limit)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, authors)
}

// HandleCategories godoc
// @Summary Get All Categories
// @ID get-all-categories
// @Accept json
// @Produce json
// @Param page query string false "page number (default=1)" Format(string)
// @Param per_page query string false "per_page product count (default=20)" Format(string)
// @Success 200 {array} item.Category
// @Failure 400 {object} httputil.HTTPError
// @Router /api/category [get]
func (ctr *Controller) HandleCategories(ctx *gin.Context) {

	offset, limit, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("query string 'page' AND 'per_page' MUST be a VALID NUMBER"))
		return
	}

	categories, err := ctr.DAO.GetAllCategories(offset, limit)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

// HandleTags godoc
// @Summary Get All Tags
// @ID get-all-tags
// @Accept json
// @Produce json
// @Param page query string false "page number (default=1)" Format(string)
// @Param per_page query string false "per_page product count (default=20)" Format(string)
// @Success 200 {array} item.Tag
// @Failure 400 {object} httputil.HTTPError
// @Router /api/tag [get]
func (ctr *Controller) HandleTags(ctx *gin.Context) {

	offset, limit, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("query string 'page' AND 'per_page' MUST be a VALID NUMBER"))
		return
	}

	tags, err := ctr.DAO.GetAllTags(offset, limit)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, tags)
}

// HandleBook godoc
// @Summary Get Book By ID
// @ID get-all-getBookByID
// @Accept json
// @Produce json
// @Param id path string true "book id to search"
// @Success 200 {object} model.Book
// @Failure 400 {object} httputil.HTTPError
// @Router /api/book/{id} [get]
func (ctr *Controller) HandleBook(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("parameter 'id' must be a valid number"))
		return
	}

	book, err := ctr.DAO.GetBookByID(id)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, book)
}

// HandleAuthor godoc
// @Summary Get Author By ID
// @ID get-all-getAuthorByID
// @Accept json
// @Produce json
// @Param id path string true "author id to search"
// @Param page query string false "page number of the item books (default=1)" Format(string)
// @Param per_page query string false "per_page product count of the item books (default=20)" Format(string)
// @Success 200 {object} model.CompAbs
// @Failure 400 {object} httputil.HTTPError
// @Router /api/author/{id} [get]
func (ctr *Controller) HandleAuthor(ctx *gin.Context) {
	filter := "authors"
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("parameter 'id' must be a valid number"))
		return
	}

	offset, limit, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("query string 'page' AND 'per_page' MUST be a VALID NUMBER"))
		return
	}

	a, err := ctr.DAO.GetItemByID(filter, id)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	b, err := ctr.DAO.GetFilterBooks(filter, a.Name, offset, limit)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, model.CompAbs{Item: a, Books: b})
}

// HandleCategory godoc
// @Summary Get Category By ID
// @ID get-all-getCategoryByID
// @Accept json
// @Produce json
// @Param id path string true "category id to search"
// @Param page query string false "page number of the item books (default=1)" Format(string)
// @Param per_page query string false "per_page product count of the item books (default=20)" Format(string)
// @Success 200 {object} model.CompAbs
// @Failure 400 {object} httputil.HTTPError
// @Router /api/category/{id} [get]
func (ctr *Controller) HandleCategory(ctx *gin.Context) {
	filter := "categories"
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("parameter 'id' must be a valid number"))
		return
	}

	offset, limit, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("query string 'page' AND 'per_page' MUST be a VALID NUMBER"))
		return
	}

	a, err := ctr.DAO.GetItemByID(filter, id)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	b, err := ctr.DAO.GetFilterBooks(filter, a.Name, offset, limit)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, model.CompAbs{Item: a, Books: b})
}

// HandleTag godoc
// @Summary Get Tag By ID
// @ID get-all-getTagByID
// @Accept json
// @Produce json
// @Param id path string true "tag id to search"
// @Param page query string false "page number of the item books (default=1)" Format(string)
// @Param per_page query string false "per_page product count of the item books (default=20)" Format(string)
// @Success 200 {object} model.CompAbs
// @Failure 400 {object} httputil.HTTPError
// @Router /api/tag/{id} [get]
func (ctr *Controller) HandleTag(ctx *gin.Context) {
	filter := "tags"
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("parameter 'id' must be a valid number"))
		return
	}

	offset, limit, err := paginate(ctx)
	if err != nil {
		httputil.NewError(ctx, http.StatusBadRequest, errors.New("query string 'page' AND 'per_page' MUST be a VALID NUMBER"))
		return
	}

	a, err := ctr.DAO.GetItemByID(filter, id)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	b, err := ctr.DAO.GetFilterBooks(filter, a.Name, offset, limit)
	if err != nil {
		httputil.NewError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, model.CompAbs{Item: a, Books: b})
}

func paginate(ctx *gin.Context) (offset int, limit int, err error) {
	pageNumS := ctx.DefaultQuery("page", "1")
	perPageS := ctx.DefaultQuery("per_page", "20")
	pageNum, err := strconv.Atoi(pageNumS)
	if err != nil {
		return -1, -1, err
	}
	limit, err = strconv.Atoi(perPageS)
	if err != nil {
		return -1, -1, err
	}
	offset = (pageNum - 1) * limit
	return
}
