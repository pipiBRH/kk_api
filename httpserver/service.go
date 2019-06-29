package httpserver

import (
	"net/http"
	"strconv"
	"time"

	"github.com/pipiBRH/kk_api/config"
	"github.com/pipiBRH/kk_api/dal"
	database "github.com/pipiBRH/kk_database"

	"github.com/gin-gonic/gin"
)

var defaultError = gin.H{
	"code":  http.StatusBadRequest,
	"error": "invalid parameter",
}

// For health check
func HC(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": "ok",
	})
}

// SearchStationByText search by user input text
func SearchStationByText(ctx *gin.Context) {
	text := ctx.Param("text")
	index := config.Config.Elasticsearch.Index["youbike"]
	pageSize, err := strconv.ParseInt(ctx.DefaultQuery("page_size", "10"), 10, 32)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			defaultError,
		)
		return
	}
	page, err := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 32)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			defaultError,
		)
		return
	}

	if page < 1 || pageSize < 1 {
		ctx.JSON(
			http.StatusBadRequest,
			defaultError,
		)
		return
	}

	esDAL := dal.NewElasticsearchDAL(database.EsClient)
	data, pagination, err := esDAL.SearchYouBikeInfoByText(index, text, int(page), int(pageSize))
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code":  http.StatusInternalServerError,
				"error": err.Error(),
			})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":       http.StatusOK,
		"data":       data,
		"pagination": pagination,
	})
}

// SearchSpaceRank search top three free space
func SearchSpaceRank(ctx *gin.Context) {
	esDAL := dal.NewElasticsearchDAL(database.EsClient)
	data, err := esDAL.SearchTopThreeFreeSpaceYouBikeInfo(config.Config.Elasticsearch.Index["youbike"])
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code":  http.StatusInternalServerError,
				"error": err.Error(),
			})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": data,
	})
}

// SearchMaxAndMinFreeSpaceByTime search max and min free space in specified time range and group by 10 minunt
func SearchMaxAndMinFreeSpaceByTime(ctx *gin.Context) {
	st := ctx.Query("st")
	et := ctx.Query("et")
	_, err := time.Parse("2006-01-02 15:04:05", st)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			defaultError,
		)
		return
	}

	_, err = time.Parse("2006-01-02 15:04:05", et)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			defaultError,
		)
		return
	}

	esDAL := dal.NewElasticsearchDAL(database.EsClient)
	data, err := esDAL.SearchMaxAndMinFreeSpaceYouBikeInfoByTimeRange(
		config.Config.Elasticsearch.Index["youbike_history"],
		st,
		et,
	)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"code":  http.StatusInternalServerError,
				"error": err.Error(),
			})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": data,
	})
}
