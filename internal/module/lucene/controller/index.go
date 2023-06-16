package controller

import (
	"github.com/yockii/celestial/internal/core/middleware"
	"github.com/yockii/ruomu-core/server"
)

func InitRouter() {
	server.Get("/api/v1/search", middleware.NeedAuthorization("user"), LuceneController.Search)
}
