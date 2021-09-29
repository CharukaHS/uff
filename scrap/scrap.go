package scrap

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"uff/db"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

var client http.Client
var count int = 0

func init() {
	// cookie jar init
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Error occured while creating cookie jar %s", err.Error())
	}

	// load .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env files")
	}

	// define cookies
	var cookies []*http.Cookie

	cookie_names := []string{"c_user", "datr", "fr", "locale", "sb", "xs"}
	for _, name := range cookie_names {
		cookie := &http.Cookie{
			Name:   name,
			Value:  os.Getenv(name),
			Path:   "/",
			Domain: ".facebook.com",
		}
		cookies = append(cookies, cookie)
	}

	// set cookies to jar
	u, _ := url.Parse("https://mbasic.facebook.com/")
	jar.SetCookies(u, cookies)

	client = http.Client{
		Jar:     jar,
		Timeout: 10 * time.Second,
	}
}

func FetchFriends(link string) {
	req, err := http.NewRequest("GET", "https://mbasic.facebook.com"+link, nil)
	if err != nil {
		log.Fatalf("Error occured while creating friend list fetch URL. %s", err.Error())
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error occured while fetching friend list. %s", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// parse HTML
	doc, err := goquery.NewDocumentFromReader(io.Reader(res.Body))
	if err != nil {
		fmt.Println("parsing error")
		log.Fatal(err)
	}

	// TODO: why td.x.u not working?
	doc.Find("td").Each(func(i int, s *goquery.Selection) {
		name := s.Find("a").Text()
		link, isLink := s.Find("a").Attr("href")
		mutual := s.Find("div").Text()

		// since goquery doesnt filter out by classes
		if !isLink {
			return
		}

		// every friend profile link contain ?fref=fr_tab
		if !strings.Contains(link, "?fref=fr_tab") {
			return
		}

		// remove unwanted parts in values
		mutual = strings.Replace(mutual, " mutual friends", "", 1)
		link = strings.Replace(link, "?fref=fr_tab", "", 1)

		// parse mutual value to int
		mutual_int, err := strconv.ParseInt(mutual, 10, 32)
		if err != nil {
			fmt.Println("Error occured while parsing mutual int of", name, mutual)
			mutual_int = 0
		}

		fmt.Printf("%s %s %s \n", name, mutual, link)
		db.InsertToFriend(name, int(mutual_int), link)
		count++
	})

	fmt.Printf("Scrapped total of %d friends", count)

	next, canNext := doc.Find("#m_more_friends a").Attr("href")

	if !canNext {
		fmt.Printf("Could not find next cursor link for fetching friends, maybe script has reached all friends")
		return
	}

	time.Sleep(200 * time.Millisecond)
	FetchFriends(next)
}
