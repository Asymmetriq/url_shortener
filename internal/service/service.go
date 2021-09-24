package service

import (
	"context"
	"database/sql"
	"log"

	"github.com/Asymmetriq/url_shortener/internal/encoding"
	"github.com/Asymmetriq/url_shortener/pkg/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const SHORT_URL_LEN = 10

type Service struct {
	api.UnimplementedServiceServer

	// public
	DB *sql.DB
	// private
	b int
}

func (s *Service) Ping(ctx context.Context, in *api.Empty) (*api.Empty, error) {
	err := s.DB.PingContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	///
	rows, err := s.DB.QueryContext(ctx, "SELECT id, value FROM test_table")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id    int
			value int
		)
		err := rows.Scan(&id, &value)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		log.Println(id, value)
	}

	_, err = s.DB.ExecContext(ctx, "UPDATE test_table SET value = 0 WHERE id = $1", 1)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	//
	return &api.Empty{}, nil
}

func (s *Service) Create(ctx context.Context, req *api.CreateRequest) (*api.CreateResponse, error) {
	ogUrl := req.GetUrl()

	shortUrl := encoding.GenerateRandomString(SHORT_URL_LEN)

	_, err := s.DB.ExecContext(ctx, "INSERT INTO Links (short_url, url) VALUES ($1, $2)", shortUrl, ogUrl)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.CreateResponse{
		ShortUrl: shortUrl,
	}, nil
}

func (s *Service) Get(ctx context.Context, req *api.GetRequest) (*api.GetResponse, error) {
	row := s.DB.QueryRowContext(ctx, "SELECT url FROM Links WHERE short_url = $1", req.GetShortUrl())

	var url string
	err := row.Scan(&url)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.GetResponse{
		Url: url,
	}, nil
}
