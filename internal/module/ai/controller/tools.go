package controller

import (
	"context"
	"fmt"
	"github.com/yockii/celestial/internal/constant"
	"github.com/yockii/ruomu-core/database"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`\s*Observation:?\s*$`)

type NeededTableInfo struct{}

func (t NeededTableInfo) Name() string {
	return "NeededTableInfo"
}

func (t NeededTableInfo) Description() string {
	return `
	数据库表及其具体字段信息，使用了建表语句。
	在生成查询数据的SQL前，应该将所有可能相关的数据库表结构获取到，并根据表结构进行sql生成，以防止出现错误的sql
	输入参数应该是以,分割的表名，这些表名应该是可能用到的所有表名
	示例: t_user,t_article
`
}

func (t NeededTableInfo) Call(ctx context.Context, input string) (string, error) {
	// 如果以\nObservation或\n\tObservation结尾，则删除
	input = re.ReplaceAllString(input, "")

	tables := strings.Split(input, ",")
	return controller.db.TableInfo(ctx, tables)
}

type SystemTableNames struct{}

func (t SystemTableNames) Name() string {
	return "SystemTableNames"
}

func (t SystemTableNames) Description() string {
	return `
	系统所用的数据库的所有表名。
	你应该在查询数据之前通过本工具获取所有表，并从中筛选出查询可能会用到的表。
	本工具只需要调用一次获取到所有表即可。
	输入参数因为固定值: tables
`
}

func (t SystemTableNames) Call(ctx context.Context, _ string) (string, error) {
	result := ""
	for _, m := range constant.Models {
		stmt := &gorm.Statement{DB: database.DB}
		_ = stmt.Parse(m)
		tableName := stmt.Schema.Table
		tableComment := m.TableComment()
		result += fmt.Sprintf("%s[%s]\n", tableName, tableComment)
	}
	return result, nil
}

type SqlExecutor struct{}

func (e SqlExecutor) Name() string {
	return "SqlExecutor"
}

func (e SqlExecutor) Description() string {
	return `
	Sql执行器。
	当需要执行SQL时调用.
	输入参数应该是纯粹的可执行的sql语句.
`
}

func (e SqlExecutor) Call(ctx context.Context, input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("no sql query generated")
	}
	// 如果input以\nObservation或\n\tObservation结尾，则删除。用正则表达式的方式删除
	input = re.ReplaceAllString(input, "")

	// 如果以```开头，则表示使用了Markdown代码块，需要删除前后的代码块标记，而代码块标记可能是```也可能是```sql
	if strings.HasPrefix(input, "```") {
		input = strings.TrimPrefix(input, "```sql")
		input = strings.TrimPrefix(input, "```")
		input = strings.TrimSuffix(input, "```")
	}

	return controller.db.Query(ctx, input)
}

type DBDialect struct{}

func (t DBDialect) Name() string {
	return "DataBaseDialect"
}

func (t DBDialect) Description() string {
	return `
	系统所用数据库的方言.
	在生成特定数据库方言的SQL时非常有用.
	输入参数应为固定值: dialect
`
}

func (t DBDialect) Call(ctx context.Context, _ string) (string, error) {
	return controller.db.Dialect(), nil
}
