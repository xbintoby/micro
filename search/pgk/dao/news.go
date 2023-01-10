package dao

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"reflect"
	"strconv"
)

type Article struct {
	Id      int32    `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Url     string   `json:"url"`
	Tags    []string `json:"tags"`
}

type newsSearch struct {
	client *elastic.Client
}

func NewNewsSearch() *newsSearch {
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
	return &newsSearch{
		client: client,
	}
}
func (news *newsSearch) QueryNews(ctx context.Context, term string) []Article {
	var searchResult *elastic.SearchResult
	var err error
	suggesterName := "news_tags_suggest"
	cs := elastic.NewCompletionSuggester(suggesterName)
	cs = cs.Size(10)
	cs = cs.SkipDuplicates(true)
	cs = cs.Text(term)
	cs = cs.Field("tags")

	searchResult, err = news.client.Search().
		Index("news").
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
	var arr []Article
	for _, arts := range mySuggest.Options {
		var a Article
		id, _ := strconv.Atoi(arts.Id)
		a = Article{int32(id), arts.Text, "", "", []string{}}
		arr = append(arr, a)
	}
	return arr
}
func (news *newsSearch) SearchNews(ctx context.Context, term string) []Article {
	var searchResult *elastic.SearchResult
	var err error
	hl := elastic.NewHighlight()
	hl = hl.Fields(elastic.NewHighlighterField("title"))
	hl = hl.PreTags("<span class='highLight'>").PostTags("</span>")
	hl2 := elastic.NewHighlight()
	hl2 = hl.Fields(elastic.NewHighlighterField("content"))
	hl2 = hl.PreTags("<span class='highLight'>").PostTags("</span>")
	q := elastic.NewMultiMatchQuery(term, "title", "content")
	fsc := elastic.NewFetchSourceContext(true).Include("url", "title", "content")
	builder := elastic.NewSearchSource().Query(q).FetchSourceContext(fsc)
	_, err = builder.Source()

	searchResult, err = news.client.Search("news").Highlight(hl).
		Highlight(hl2).
		Query(q).
		Size(10).Pretty(true).
		Do(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	//for i := 0; i < len(searchResult.Hits.Hits); i++ {
	//	fmt.Println("Title:", searchResult.Hits.Hits[i].Highlight["title"])
	//}
	//hit := searchResult.Hits.Hits[0]
	//
	//fmt.Println(hit.Highlight["title"])
	var arr []Article
	var typ Article

	for i, item := range searchResult.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
		t := item.(Article)
		hit := searchResult.Hits.Hits[i]
		if hl, found := hit.Highlight["title"]; found {
			t.Title = hl[0]

		}
		if hl, found := hit.Highlight["content"]; found {
			t.Content = hl[0]

		}
		id := t.Id + '0'
		bb := Article{Id: id, Title: t.Title, Content: t.Content, Url: t.Url, Tags: t.Tags}

		arr = append(arr, bb)
		//fmt.Printf("%#v\n", t)

	}

	return arr
}
