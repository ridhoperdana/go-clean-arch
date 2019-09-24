package graphql

import (
	"context"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/ridhoperdana/go-clean-arch/article"
	"github.com/ridhoperdana/go-clean-arch/article/repository"
	"github.com/ridhoperdana/go-clean-arch/models"
	"time"
)

// ArticleEdge holds information of article edge.
type ArticleEdge struct {
	Node   models.Article
	Cursor string
}

// ArticleResult holds information of article result.
type ArticleResult struct {
	Edges    []ArticleEdge
	PageInfo PageInfo
}

// PageInfo holds information of page info.
type PageInfo struct {
	EndCursor   string
	HasNextPage bool
}

type Resolver interface {
	FetchArticle(params graphql.ResolveParams) (interface{}, error)
	GetArticleByID(params graphql.ResolveParams) (interface{}, error)
	GetArticleByTitle(params graphql.ResolveParams) (interface{}, error)

	UpdateArticle(params graphql.ResolveParams) (interface{}, error)
	StoreArticle(params graphql.ResolveParams) (interface{}, error)
	DeleteArticle(params graphql.ResolveParams) (interface{}, error)
}

type resolver struct {
	articleService article.Usecase
}

func (r resolver) FetchArticle(params graphql.ResolveParams) (interface{}, error) {
	ctx := context.Background()
	num := 0
	cursor := ""
	if cursorFromClient, ok := params.Args["after"].(string); ok {
		cursor = cursorFromClient
	}

	if numFromClient, ok := params.Args["first"].(int); ok {
		num = numFromClient
	}

	results, cursorFromService, err := r.articleService.Fetch(ctx, cursor, int64(num))
	if err != nil {
		return nil, err
	}

	edges := make([]ArticleEdge, len(results))
	for index, result := range results {
		if result != nil {
			edges[index] = ArticleEdge{
				Node: *result,
				Cursor: repository.EncodeCursor(result.CreatedAt),
			}
		}
	}

	isHasNextPage := false
	if len(results) > 0 {
		results, _, err := r.articleService.Fetch(ctx, cursorFromService, int64(1))
		if err != nil {
			return nil, err
		}

		if len(results) > 0 {
			isHasNextPage = true
		}
	}

	return ArticleResult{
		Edges: edges,
		PageInfo:PageInfo{
			EndCursor: cursorFromService,
			HasNextPage:isHasNextPage,
		},
	}, nil
}

func (r resolver) GetArticleByID(params graphql.ResolveParams) (interface{}, error) {
	var (
		id int
		ok bool
	)

	ctx := context.Background()
	if id, ok = params.Args["id"].(int); !ok || id == 0 {
		return nil, fmt.Errorf("id is not integer or zero")
	}

	result, err := r.articleService.GetByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func (r resolver) GetArticleByTitle(params graphql.ResolveParams) (interface{}, error) {
	var (
		title string
		ok bool
	)

	ctx := context.Background()

	if title, ok = params.Args["title"].(string); !ok || title == "" {
		return nil, fmt.Errorf("title is empty or not string")
	}

	result, err := r.articleService.GetByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	return *result, nil
}

func (r resolver) UpdateArticle(params graphql.ResolveParams) (interface{}, error) {
	var (
		id int
		title, content string
		ok bool
	)

	ctx := context.Background()
	if id, ok = params.Args["id"].(int); !ok || id == 0 {
		return nil, fmt.Errorf("id is not integer or zero")
	}

	if title, ok = params.Args["title"].(string); !ok || title == "" {
		return nil, fmt.Errorf("title is empty or not string")
	}

	if content, ok = params.Args["content"].(string); !ok {
		return nil, fmt.Errorf("content is not string")
	}

	updatedArticle := &models.Article{
		ID: int64(id),
		Title: title,
		Content: content,
		UpdatedAt: time.Now(),
	}

	if err := r.articleService.Update(ctx, updatedArticle); err != nil {
		return nil, err
	}

	return *updatedArticle, nil
}

func (r resolver) StoreArticle(params graphql.ResolveParams) (interface{}, error) {
	var (
		title, content string
		ok bool
	)

	ctx := context.Background()

	if title, ok = params.Args["title"].(string); !ok || title == "" {
		return nil, fmt.Errorf("title is empty or not string")
	}

	if content, ok = params.Args["content"].(string); !ok {
		return nil, fmt.Errorf("content is not string")
	}

	storedArticle := &models.Article{
		Content: content,
		Title:title,
	}

	if err := r.articleService.Store(ctx, storedArticle); err != nil {
		return nil, err
	}

	return *storedArticle, nil
}

func (r resolver) DeleteArticle(params graphql.ResolveParams) (interface{}, error) {
	var (
		id int
		ok bool
	)

	ctx := context.Background()
	if id, ok = params.Args["id"].(int); !ok || id == 0 {
		return nil, fmt.Errorf("id is not integer or zero")
	}

	if err := r.articleService.Delete(ctx, int64(id)); err != nil {
		return nil, err
	}

	return id, nil
}

func NewResolver(articleService article.Usecase) Resolver {
	return &resolver{
		articleService:articleService,
	}
}
