package controller

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
	"github.com/tmc/langchaingo/tools/sqldatabase"
	_ "github.com/tmc/langchaingo/tools/sqldatabase/mysql"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/celestial/internal/core/middleware"
	"github.com/yockii/ruomu-core/config"
	"github.com/yockii/ruomu-core/server"
)

const (
	mysqlConnStringFmt = "%s:%s@tcp(%s:%d)/%s"
)

func InitRouter() {
	server.Get("/api/v1/ai/dataSearch", middleware.NeedAuthorization(constant.ResourceAI), controller.Generate)
}

func init() {

	client, err := openai.New(
		openai.WithToken(config.GetString("openai.openaiToken")),
		openai.WithBaseURL(config.GetString("openai.openaiBaseURL")),
		openai.WithModel(config.GetString("openai.modelName")),
	)

	if err != nil {
		logrus.Warnln(err)
		return
	}
	db, err := sqldatabase.NewSQLDatabaseWithDSN("mysql", fmt.Sprintf(mysqlConnStringFmt,
		config.GetString("database.readonlyUser"),
		config.GetString("database.readonlyPassword"),
		config.GetString("database.host"),
		config.GetInt("database.port"),
		config.GetString("database.db"),
	), nil)
	if err != nil {
		logrus.Fatal(err)
	}

	agentTools := []tools.Tool{
		SystemTableNames{},
		NeededTableInfo{},
		DBDialect{},
		SqlExecutor{},
	}

	agent := agents.NewOneShotAgent(client, agentTools)

	controller = &aiController{
		client: client,
		db:     db,
		agent:  agent,
	}
}

var controller *aiController

type aiController struct {
	client llms.Model
	db     *sqldatabase.SQLDatabase
	agent  *agents.OneShotZeroAgent
}

func (c *aiController) Generate(ctx *fiber.Ctx) error {
	question := ctx.Query("question")
	if question == "" {
		return ctx.JSON(&server.CommonResponse{
			Code: server.ResponseCodeParamNotEnough,
			Msg:  server.ResponseMsgParamNotEnough + " question",
		})
	}

	//sqlDatabaseChain := chains.NewSQLDatabaseChain(c.client, 100, c.db)
	//input := map[string]any{
	//	"query": question,
	//}
	//result, err := chains.Predict(context.Background(), sqlDatabaseChain, input)

	executor := agents.NewExecutor(
		c.agent,
		agents.WithCallbacksHandler(callbacks.LogHandler{}),
		agents.WithMaxIterations(10),
	)
	result, err := chains.Run(context.Background(), executor,
		question+"\n请遵循以下原则：1.若无数据则直接告诉我没有数据；2.若发现有字段是在其他表中关联的，则继续请求关联表信息; 3.不要猜测字段，请获取确切的表及字段信息",
	)

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
