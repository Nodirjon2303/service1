package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"service1/internal/database"
	"service1/internal/models"
	"service1/proto"
)

type Collector struct {
	DB *mydatabase.DB
	proto.UnimplementedCollectorServer
}

func NewCollector(db *mydatabase.DB) *Collector {
	return &Collector{
		DB: db,
	}
}

func (c *Collector) CollectPosts(ctx context.Context, req *proto.CollectPostsRequest) (*proto.CollectPostsResponse, error) {
	// Perform data collection logic here
	var posts []models.Post

	for page := 1; page <= 50; page++ {
		url := fmt.Sprintf("https://gorest.co.in/public/v1/posts?page=%d", page)

		response, err := http.Get(url)
		if err != nil {
			log.Printf("Error retrieving posts from page %d: %s", page, err)
			continue
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			log.Printf("Failed to retrieve posts from page %d: %s", page, response.Status)
			continue
		}

		var result struct {
			Data []models.Post `json:"data"`
		}

		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			log.Printf("Error decoding response from page %d: %s", page, err)
			continue
		}

		posts = append(posts, result.Data...)
	}

	// Save posts to the database
	for _, post := range posts {
		err := c.DB.SavePost(&post)
		if err != nil {
			log.Printf("Error saving post with ID %d: %s", post.ID, err)
		}
	}

	return &proto.CollectPostsResponse{
		Success: true,
	}, nil
}
