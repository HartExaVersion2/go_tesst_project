package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type News struct {
	ID         int64   `json:"Id"`
	Title      string  `json:"Title"`
	Content    string  `json:"Content"`
	Categories []int64 `json:"Categories"`
}

type Categories struct {
	ID      int64
	News_ID int64
}

// init is invoked before main()
var db *sql.DB

func init() {

	godotenv.Load()

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/reform", user, password))
	if err != nil {
		log.Fatal(err)
	}

}

func GetAllNews() ([]News, error) {
	var news_list []News

	rows, err := db.Query("SELECT News.Id, News.Title, News.Content, (SELECT GROUP_CONCAT(NC.CategoryId) FROM NewsCategories NC WHERE NC.NewsId = News.Id) as Categories FROM News;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var news News
		var categoriesStr string
		if err := rows.Scan(&news.ID, &news.Title, &news.Content, &categoriesStr); err != nil {
			return nil, err
		}
		// Разбить строку на список строк по запятым
		categoriesArr := strings.Split(categoriesStr, ",")
		// Преобразовать каждую строку в int64
		var categoriesInt []int64
		for _, catStr := range categoriesArr {
			catInt, err := strconv.ParseInt(catStr, 10, 64)
			if err != nil {
				return nil, err
			}
			categoriesInt = append(categoriesInt, catInt)
		}
		news.Categories = categoriesInt
		news_list = append(news_list, news)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return news_list, nil
}

func PushOne(news News) {
	result, err := db.Exec("INSERT INTO News (Title, Content) VALUES (?, ?)", news.Title, news.Content)
	if err != nil {
		fmt.Println("get_list_news: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("get_list_news: %v", err)
	}
	list_categories := news.Categories
	for _, value := range list_categories {
		_, err := db.Exec("INSERT INTO NewsCategories (NewsId, CategoryId) VALUES (?, ?)", id, value)
		if err != nil {
			fmt.Println("get_list_news: %v", err)
		}
	}
}

func UpdateOne(news News) {
	_, err := db.Exec("UPDATE News SET Title = ?, Content = ? WHERE Id = ?", news.Title, news.Content, news.ID)
	if err != nil {
		fmt.Println("get_list_news: %v", err)
	}

	_, err = db.Exec("DELETE FROM NewsCategories WHERE NewsId = ?", news.ID)
	if err != nil {
		fmt.Println("get_list_news: %v", err)
	}

	// Добавляем категории в таблицу NewsCategories
	for _, categoryID := range news.Categories {
		_, err = db.Exec("INSERT INTO NewsCategories (NewsId, CategoryId) VALUES (?, ?)", news.ID, categoryID)
		if err != nil {
			fmt.Println("get_list_news: %v", err)
		}
	}
}
