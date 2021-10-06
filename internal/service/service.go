package service

import (
	"context"

	"github.com/Asymmetriq/url_shortener/internal/encoding"
	"github.com/Asymmetriq/url_shortener/pkg/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const SHORT_URL_LEN = 10

type CustomStorage interface {
	GetURLByShortURL(ctx context.Context, shortURL string) (string, error)
	SaveShortURL(ctx context.Context, shortUrl, ogUrl string) error
}

type Service struct {
	api.UnimplementedServiceServer

	// public
	Storage CustomStorage
}

func (s *Service) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	ogUrl := req.GetUrl()

	shortUrl := encoding.GenerateRandomString(SHORT_URL_LEN)

	err := s.Storage.SaveShortURL(ctx, shortUrl, ogUrl)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.CreateResponse{
		ShortUrl: shortUrl,
	}, nil
}

func (s *Service) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	url, err := s.Storage.GetURLByShortURL(ctx, req.GetShortUrl())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.GetResponse{
		Url: url,
	}, nil
}
