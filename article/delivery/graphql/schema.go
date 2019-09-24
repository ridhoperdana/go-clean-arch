package graphql

import "github.com/graphql-go/graphql"

// ArticleGraphQL holds article information with graphql object
var ArticleGraphQL = graphql.NewObject(
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

// ArticleEdgeGraphQL holds article edge information with graphql object
var ArticleEdgeGraphQL = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ArticleEdge",
		Fields: graphql.Fields{
			"node": &graphql.Field{
				Type: ArticleGraphQL,
			},
			"cursor": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

// ArticleResultGraphQL holds article result information with graphql object
var ArticleResultGraphQL = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ArticleResult",
		Fields: graphql.Fields{
			"edges": &graphql.Field{
				Type: graphql.NewList(ArticleEdgeGraphQL),
			},
			"pageInfo": &graphql.Field{
				Type: pageInfoGraphQL,
			},
		},
	},
)

var pageInfoGraphQL = graphql.NewObject(
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
			"FetchArticle": &graphql.Field{
				Type:        ArticleResultGraphQL,
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
			"GetArticleByID": &graphql.Field{
				Type:        ArticleGraphQL,
				Description: "Get Article By ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: s.articleResolver.GetArticleByID,
			},
			"GetArticleByTitle": &graphql.Field{
				Type:        ArticleGraphQL,
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
			"UpdateArticle": &graphql.Field{
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
			"StoreArticle": &graphql.Field{
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
			"DeleteArticle": &graphql.Field{
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
