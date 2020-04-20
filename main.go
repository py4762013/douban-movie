package main

import (
	"douban-movie/model"
	"douban-movie/parse"
	"log"
	"strings"
)

var (
	BaseUrl = "https://movie.douban.com/top250"
	//BaseUrl = "http://studygolang.com/topic"
)

// 新增数据
func Add(movies []parse.DoubanMovie)  {
	for index, movie := range movies {
		if err := model.DB.Create(&movie).Error; err != nil {
			log.Printf("db.Create index: %s, err: %v, data: %v", index, err, movie)
		}
	}
}

// 开始爬取
func Start()  {
	var movies []parse.DoubanMovie

	pages := parse.GetPages(BaseUrl)
	for _, page := range pages {
		doc := parse.Fetch(strings.Join([]string{BaseUrl, page.Url}, ""))
		/*doc, err := goquery.NewDocument(strings.Join([]string{BaseUrl, page.Url}, ""))
		if err != nil {
			log.Println(err)
		}*/

		movies = append(movies, parse.ParseMovies(doc)...)
	}

	Add(movies)
}

func main()  {
	Start()

	defer model.DB.Close()
}