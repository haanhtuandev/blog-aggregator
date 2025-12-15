package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// crafting a request and send
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	req.Header.Set("User-Agent", "gator")
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	rssFeed := &RSSFeed{}
	err = xml.Unmarshal(data, rssFeed)
	if err != nil {
		return nil, err
	}
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for _, feed := range rssFeed.Channel.Item {
		feed.Description = html.UnescapeString(feed.Description)
		feed.Title = html.UnescapeString(feed.Title)
		feed.Link = html.UnescapeString(feed.Link)
		feed.PubDate = html.UnescapeString(feed.PubDate)
	}

	return rssFeed, nil
}
