package main

import (
	"database/sql"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
	"log"
	"net/url"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"github.com/graphql-go/handler"

	_articleHttpDeliver "github.com/ridhoperdana/go-clean-arch/article/delivery/http"
	_graphQLArticleDelivery "github.com/ridhoperdana/go-clean-arch/article/delivery/graphql"
	_articleRepo "github.com/ridhoperdana/go-clean-arch/article/repository"
	_articleUcase "github.com/ridhoperdana/go-clean-arch/article/usecase"
	_authorRepo "github.com/ridhoperdana/go-clean-arch/author/repository"
	"github.com/ridhoperdana/go-clean-arch/middleware"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil && viper.GetBool("debug") {
		fmt.Println(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)
	authorRepo := _authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := _articleRepo.NewMysqlArticleRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	au := _articleUcase.NewArticleUsecase(ar, authorRepo, timeoutContext)
	_articleHttpDeliver.NewArticleHandler(e, au)

	schema := _graphQLArticleDelivery.NewSchema(_graphQLArticleDelivery.NewResolver(au))
	graphqlSchema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    schema.Query(),
		Mutation: schema.Mutation(),
	})
	if err != nil {
		logrus.Fatal(err)
	}

	graphQLHandler := handler.New(&handler.Config{
		Schema: &graphqlSchema,
		GraphiQL:true,
		Pretty:true,
	})

	e.GET("/graphql", echo.WrapHandler(graphQLHandler))
	e.POST("/graphql", echo.WrapHandler(graphQLHandler))

	log.Fatal(e.Start(viper.GetString("server.address")))
}
