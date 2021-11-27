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

var newID int64

func (s *Service) Save(ctx context.Context, item *Banner, file multipart.File) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// check if id is equal to 0
	if item.ID == 0 {
		newID++
		item.ID = newID
		if item.Image != "" {
			item.Image = fmt.Sprint(item.ID) + "." + item.Image

			data, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, errors.New("not readible data")
			}

			err = ioutil.WriteFile(".web/banners/"+item.Image, data, 0666)
			if err != nil {
				return nil, err
			}
		}

		s.items = append(s.items, item)
		return item, nil
	} else {

		sBanner, err := s.ByID(ctx, item.ID)
		if err != nil {
			log.Print(err)
			return nil, errors.New("item not found")
		}

		if item.Image != "" {
			item.Image = fmt.Sprint(item.ID) + "." + item.Image

			data, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, errors.New("not readible data")
			}

			err = ioutil.WriteFile(".web/banners/"+item.Image, data, 0666)
			if err != nil {
				return nil, err
			}
		} else {
			sBanner = item
			return sBanner, nil
		}

	}

	return nil, errors.New("item not found")
}
