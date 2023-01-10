package search_service_v1

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"jam3.com/search/pgk/dao"
	"jam3.com/search/pgk/repo"
	"time"
)

type SearchService struct {
	UnimplementedSearchServiceServer
	cache repo.Cache
	db    *gorm.DB
}

func New() *SearchService {
	return &SearchService{
		cache: dao.Rc,
	}
}
func (ls *SearchService) StarQuery(ctx context.Context, msg *StarMessage) (*StarResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	term := msg.Term
	s := dao.NewStarSearch()
	stars := s.Query(c, term)
	mp := []string{}
	for _, v := range stars {
		mp = append(mp, v.Name)
	}
	fmt.Println("starts:", mp)
	return &StarResponse{Stars: mp}, nil
}

func (ls *SearchService) NewsQuery(ctx context.Context, msg *NewsMessage) (*NewsResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	term := msg.Term
	s := dao.NewNewsSearch()
	articles := s.QueryNews(c, term)
	mp := []string{}
	for _, v := range articles {
		mp = append(mp, v.Title)
	}
	fmt.Println("news:", mp)
	return &NewsResponse{Title: mp}, nil
}

func (ls *SearchService) SearchNews(ctx context.Context, msg *ArticleMessage) (*ArticleResponse, error) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	term := msg.Text
	s := dao.NewNewsSearch()
	arts := s.SearchNews(c, term)
	var res []*ArticleType
	for _, v := range arts {
		a := ArticleType{
			Id:      v.Id,
			Title:   v.Title,
			Tags:    v.Tags,
			Url:     v.Url,
			Content: v.Content,
		}
		res = append(res, &a)
	}
	return &ArticleResponse{Arts: res}, nil
}
