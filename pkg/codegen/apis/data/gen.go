package data

//go:generate oapi-codegen -generate types -package data -templates ../../templates -o types.go -import-mapping=../common/common.yaml:github.com/Jacobbrewer1/f1-data/pkg/codegen/apis/common ./routes.yaml
//go:generate oapi-codegen -generate gorilla -package data -templates ../../templates -o server.go -import-mapping=../common/common.yaml:github.com/Jacobbrewer1/f1-data/pkg/codegen/apis/common ./routes.yaml
//go:generate oapi-codegen -generate client -package data -templates ../../templates -o client.go -import-mapping=../common/common.yaml:github.com/Jacobbrewer1/f1-data/pkg/codegen/apis/common ./routes.yaml
