package graphql

import "github.com/graphql-go/graphql"

type Resolver interface {
	FetchArticle(params graphql.ResolveParams) (interface{}, error)
	GetArticleByID(params graphql.ResolveParams) (interface{}, error)
	GetArticleByTitle(params graphql.ResolveParams) (interface{}, error)

	UpdateArticle(params graphql.ResolveParams) (interface{}, error)
	StoreArticle(params graphql.ResolveParams) (interface{}, error)
	DeleteArticle(params graphql.ResolveParams) (interface{}, error)
}
