module github.com/{{ .Scaffold.GitHubUser }}/{{ .Project }}

go 1.24.0

require (
    github.com/aws/aws-sdk-go v1.55.5
    github.com/aws/aws-sdk-go-v2/config v1.28.7
	github.com/aws/aws-sdk-go-v2/credentials v1.17.48
	github.com/gin-gonic/gin v1.10.0
	github.com/go-resty/resty/v2 v2.16.2
    github.com/pkg/errors v0.9.1
	github.com/xdatadev/superapp-packages/superapp-common v0.0.6
	github.com/xdatadev/superapp-packages/superappdb v0.0.10
	google.golang.org/api v0.186.0
    github.com/google/uuid v1.6.0
    gorm.io/gorm v1.25.12
)
