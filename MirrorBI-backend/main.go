package main

import (
	"mrbi/cmd"
	_ "mrbi/docs"
	_ "mrbi/pkg/rds"
)

// @title           MirrorBI
// @version         1.0
// @description		明镜数据分析接口文档
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cmd.Main()
}
