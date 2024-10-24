package controller

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools/sqldatabase"
	"github.com/yockii/celestial/internal/core/middleware"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/server"
)

const (
	mysqlConnStringFmt = "%s:%s@tcp(%s:%d)/%s"
)

func InitRouter() {
	server.Get("/api/v1/aiSearch", middleware.NeedAuthorization("user"), controller.Generate)
}

func init() {
	client, err := openai.New(
		openai.WithToken(openaiToken),
		openai.WithBaseURL(openaiBaseURL),
		openai.WithModel(modelName),
	)
	if err != nil {
		logrus.Fatal(err)
		return
	}
	db, err := sqldatabase.NewSQLDatabaseWithDSN("mysql", fmt.Sprintf(mysqlConnStringFmt,
		config.GetString("database.readonlyUser"),
		config.GetString("database.readonlyPassword"),
		config.GetString("database.host"),
		config.GetInt("database.port"),
		config.GetString("database.name"),
	), nil)
	if err != nil {
		logrus.Fatal(err)
	}

	controller = &aiController{
		client: client,
		db:     db,
	}
}

var controller *aiController

type aiController struct {
	client *openai.LLM
	db     *sqldatabase.SQLDatabase
}

func (c *aiController) Generate(ctx *fiber.Ctx) error {
	query := ctx.Query("query")
	if query == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " query",
		})
	}

	sqlDatabaseChain := chains.NewSQLDatabaseChain(c.client, 100, c.db)
	input := map[string]any{
		"query": query,
	}
	result, err := chains.Predict(context.Background(), sqlDatabaseChain, input)
	if err != nil {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeGeneration,
			Msg:  err.Error(),
		})
	}
	return ctx.JSON(&server.CommonResponse{
		Data: result,
	})
}
