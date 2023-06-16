package controller

import (
	"github.com/gofiber/fiber/v2"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/celestial/pkg/search"
	"github.com/yockii/ruomu-core/server"
)

var LuceneController = new(luceneController)

type luceneController struct{}

func (c luceneController) Search(ctx *fiber.Ctx) error {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " keyword",
		})
	}

	paginate := new(server.Paginate)
	if err := ctx.QueryParser(paginate); err != nil {
		logger.Errorln(err)
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamParseError,
			Msg:  server.ResponseMsgParamParseError,
		})
	}
	if paginate.Limit == 0 {
		paginate.Limit = 10
	}
	total, documents, err := search.Search(keyword, int64(paginate.Limit), int64(paginate.Offset))
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeDatabase,
			Msg:  server.ResponseMsgDatabase + err.Error(),
		})
	}
	paginate.Total = total
	paginate.Items = documents
	return ctx.JSON(&server.CommonResponse{
		Data: paginate,
	})
}
