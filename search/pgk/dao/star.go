package dao

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"jam3.com/search/config"
	"log"
	"os"
)

type Star struct {
	Id   string `json:"name"`
	Name string `json:"name"`
}

var host = config.C.ES.Addr

type starSearch struct {
	client *elastic.Client
}

//var StartSeach = NewStarSearch()

func NewStarSearch() *starSearch {
	errorlog := log.New(os.Stdout, "[app-search]", log.LstdFlags)
	var err error
	var client *elastic.Client
	client, err = elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetURL(host), elastic.SetSniff(false))
	fmt.Println(client)
	if err != nil {
		panic(err)
	}
	info, code, err := client.Ping(host).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Es returned with code %d and version %s\n", code, info.Version.Number)
	esversionCode, err := client.ElasticsearchVersion(host)
	if err != nil {
		panic(err)
	}
	fmt.Printf("es version %s\n", esversionCode)
	return &starSearch{
		client: client,
	}
}

func (star *starSearch) Query(ctx context.Context, term string) []Star {
	var searchResult *elastic.SearchResult
	var err error
	suggesterName := "star_name_suggest"
	cs := elastic.NewCompletionSuggester(suggesterName)
	cs = cs.Size(10)
	cs = cs.SkipDuplicates(true)
	cs = cs.Text(term)
	cs = cs.Field("name")

	searchResult, err = star.client.Search().
		Index("stars").
		//Query(elastic.NewMatchAllQuery()).
		Suggester(cs).
		Pretty(true).
		Do(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	if searchResult.Suggest == nil {
		fmt.Errorf("expected SearchResult.Suggest != nil; got nil")
	}
	mySuggestions, found := searchResult.Suggest[suggesterName]
	if !found {
		fmt.Errorf("expected to find SearchResult.Suggest[%s]; got false", suggesterName)
	}
	if mySuggestions == nil {
		fmt.Errorf("expected SearchResult.Suggest[%s] != nil; got nil", suggesterName)
	}

	mySuggest := mySuggestions[0]
	var arr []Star
	for _, star := range mySuggest.Options {
		var a Star
		a = Star{star.Id, star.Text}
		arr = append(arr, a)
	}
	return arr
}
