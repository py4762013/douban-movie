package parse

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type DoubanMovie struct {
	Title		string
	Subtitle	string
	Other		string
	Desc 		string
	Year		int
	Area		string
	Tag 		string
	Star 		string
	Comment 	string
	Quote 		string
}

type Page struct {
	Page 	int
	Url 	string
}

func Fetch(url string) *goquery.Document {
	fmt.Println("Fetch Url", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Http get err:", err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Http status code:", resp.StatusCode)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

// 获取分页
func GetPages(url string) []Page {
	/*doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}*/
	doc := Fetch(url)
	//fmt.Println(doc.Html())

	return ParsePages(doc)
}

// 分析分页
func ParsePages(doc *goquery.Document) (pages []Page) {
	pages = append(pages, Page{Page: 1, Url: ""})
	doc.Find("#content > div > div.article > div.paginator > a").Each(func(i int, selection *goquery.Selection) {
		page, _ := strconv.Atoi(selection.Text())
		url, _ := selection.Attr("href")

		pages = append(pages, Page{
			Page: page,
			Url: url,
		})
	})
	return pages
}

// 分析电影数据
func ParseMovies(doc *goquery.Document) (movies []DoubanMovie) {
	doc.Find("#content > div > div.article > ol > li").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".hd a span").Eq(0).Text()

		subtitle := s.Find(".hd a span").Eq(1).Text()
		subtitle = strings.TrimLeft(subtitle, " /")

		other := s.Find(".hd a span").Eq(2).Text()
		other = strings.TrimLeft(other, " /")

		desc := strings.TrimSpace(s.Find(".bd p").Eq(0).Text())
		DescInfo := strings.Split(desc, "\n")
		desc = DescInfo[0]

		movieDesc := strings.Split(DescInfo[1], "/")
		year, _ := strconv.Atoi(strings.TrimSpace(movieDesc[0]))
		area := strings.TrimSpace(movieDesc[1])
		tag := strings.TrimSpace(movieDesc[2])

		star := s.Find(".bd .star .rating_num").Text()

		comment := strings.TrimSpace(s.Find(".bd .star span").Eq(3).Text())
		compile := regexp.MustCompile("[0-9]")
		comment = strings.Join(compile.FindAllString(comment, -1), "")

		quote := s.Find(".quote .inq").Text()

		movie := DoubanMovie{
			Title:		title,
			Subtitle: 	subtitle,
			Other:		other,
			Desc:		desc,
			Year:		year,
			Area:		area,
			Tag:		tag,
			Star:		star,
			Comment:	comment,
			Quote:		quote,
		}

		log.Printf("i: %d, movie: %v", i, movie)

		movies = append(movies, movie)
	})
	return movies
}