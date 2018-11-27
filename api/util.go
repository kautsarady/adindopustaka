package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type dataContext struct {
	Metadata metadata    `json:"metadata"`
	Data     interface{} `json:"data"`
}

type metadata struct {
	Entity  string `json:"entity"`
	Page    int    `json:"page"`
	PerPage int    `json:"per_page"`
	Next    int    `json:"next"`
	Prev    int    `json:"prev"`
}

func paginate(ctx *gin.Context) (limit int, offset int, err error) {
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

func wrapData(entity string, limit, offset int, data interface{}) dataContext {
	page, perPage := (offset/limit)+1, limit
	next, prev := page+1, 0
	if page-1 > 0 {
		prev = page - 1
	}
	return dataContext{
		metadata{entity, page, perPage, next, prev},
		data,
	}
}
