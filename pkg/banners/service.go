package banners

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"sync"
)

// this is for managing banners
type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
	Image   string
}

// it creats new service
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

// All
func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.items != nil {
		return s.items, nil
	}

	return nil, errors.New("No banners")
}

// GetById
func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}

	return nil, errors.New("item not found")
}

var newID int64

//Save and update
func (s *Service) Save(ctx context.Context, item *Banner, file multipart.File) (*Banner, error) {

	if item.ID == 0 {
		newID++

		if item.Image != "" {
			item.Image = fmt.Sprint(item.ID) + "." + item.Image
			err := uploadFile(file, "./web/banners/"+item.Image)
			if err != nil {
				return nil, err
			}
		}

		newBnner := &Banner{
			ID:      newID,
			Title:   item.Title,
			Content: item.Content,
			Button:  item.Button,
			Link:    item.Link,
			Image:   item.Image,
		}
		s.items = append(s.items, newBnner)
		return newBnner, nil
	}

	ExistB, err := s.ByID(ctx, item.ID)

	if err != nil {
		log.Print(err)
		return nil, errors.New("item not found")
	}

	ExistB.Button = item.Button
	ExistB.Title = item.Title
	ExistB.Content = item.Content
	ExistB.Link = item.Link
	if item.Image != "" {
		item.Image = fmt.Sprint(item.ID) + "." + item.Image
		err := uploadFile(file, "./web/banners/"+item.Image)
		if err != nil {
			return nil, err
		}
		ExistB.Image = item.Image
	}

	return ExistB, nil

}

func uploadFile(file multipart.File, path string) error {

	var data, err = ioutil.ReadAll(file)

	if err != nil {
		return errors.New("not readble data")
	}

	err = ioutil.WriteFile(path, data, 0666)

	if err != nil {
		return errors.New("not saved from folder ")
	}

	return nil
}

// Remove by id
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	sBanner, err := s.ByID(ctx, id)
	if err != nil {
		log.Print(err)
		return nil, errors.New("item not found")
	}
	for i, banner := range s.items {
		if banner.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			break
		}
	}

	return sBanner, nil
}
