package graphql

import "github.com/graphql-go/graphql"

// Article holds article information with graphql object
var Article = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Article",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"content": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.DateTime,
			},
			"created_at": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)

// ArticleEdge holds article edge information with graphql object
var ArticleEdge = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ArticleEdge",
		Fields: graphql.Fields{
			"node": &graphql.Field{
				Type: Article,
			},
			"cursor": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// ArticleResult holds article result information with graphql object
var ArticleResult = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ArticleResult",
		Fields: graphql.Fields{
			"edges": &graphql.Field{
				Type: graphql.NewList(ArticleEdge),
			},
			"pageInfo": &graphql.Field{
				Type: pageInfo,
			},
		},
	},
)

var pageInfo = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "PageInfo",
		Fields: graphql.Fields{
			"endCursor": &graphql.Field{
				Type: graphql.String,
			},
			"hasNextPage": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)

// Schema is struct which has method for Query and Mutation. Please init this struct using constructor function.
type Schema struct {
	articleResolver Resolver
}

// Query initializes config schema query for graphql server.
func (s Schema) Query() *graphql.Object {
	objectConfig := graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"fetchArticle": &graphql.Field{
				Type:        ArticleResult,
				Description: "Fetch Article",
				Args: graphql.FieldConfigArgument{
					"first": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"after": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: s.articleResolver.FetchArticle,
			},
			"getArticleByID": &graphql.Field{
				Type:        Article,
				Description: "Get Article By ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: s.articleResolver.GetArticleByID,
			},
			"getArticleByTitle": &graphql.Field{
				Type:        Article,
				Description: "Get Article By Title",
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: s.articleResolver.GetArticleByTitle,
			},
		},
	}

	return graphql.NewObject(objectConfig)
}

// Mutation initializes config schema mutation for graphql server.
func (s Schema) Mutation() *graphql.Object {
	objectConfig := graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"updateArticle": &graphql.Field{
				Type:        graphql.String,
				Description: "Update article by certain ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: s.articleResolver.UpdateArticle,
			},
			"storeArticle": &graphql.Field{
				Type:        graphql.String,
				Description: "Store a new article",
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: s.articleResolver.StoreArticle,
			},
			"deleteArticle": &graphql.Field{
				Type:        graphql.String,
				Description: "Delete an article by its ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: s.articleResolver.DeleteArticle,
			},
		},
	}

	return graphql.NewObject(objectConfig)
}

// NewSchema initializes Schema struct which takes resolver as the argument.
func NewSchema(articleResolver Resolver) Schema {
	return Schema{
		articleResolver: articleResolver,
	}
}
