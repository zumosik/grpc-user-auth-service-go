package server

import (
	"context"
	"errors"
	pb2 "github.com/zumosik/grpc-user-auth-service-go/pb"
	"github.com/zumosik/grpc-user-auth-service-go/storage"
	"log"
	"time"
)

type Storage interface {
	GetUserByID(ctx context.Context, id uint) (storage.User, error)
	GetUsers(ctx context.Context, limit uint) ([]storage.User, error)
	CreateUser(ctx context.Context, u storage.User) (uint, error)
	UpdateUser(ctx context.Context, u storage.User) error
}

//type UserServiceServer interface {
//	GetUserDetails(context.Context, *GetUserDetailsRequest) (*GetUserDetailsResponse, error)
//	GetUsers(context.Context, *GetUsersRequest) (*GetUsersResponse, error)
//	CreateUpdateUser(context.Context, *CreateUpdateUserRequest) (*CreateUpdateUserResponse, error)
//	mustEmbedUnimplementedUserServiceServer()
//}

type Server struct {
	pb2.UnimplementedUserServiceServer
	st      Storage
	timeout time.Duration
}

func New(st Storage, timeout time.Duration) *Server {
	return &Server{
		st:      st,
		timeout: timeout,
	}
}

func (s *Server) GetUsers(ctx context.Context, _ *pb2.GetUsersRequest) (*pb2.GetUsersResponse, error) {
	log.Println("GetUsers")
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	users, err := s.st.GetUsers(ctx, 100)
	if err != nil {

		return nil, err
	}

	var userResponses []*pb2.GetUserResponse

	for _, u := range users {
		resp := &pb2.GetUserResponse{
			Id:       uint64(u.ID),
			Email:    u.Email,
			Username: u.Username,
			Password: u.Password,
		}

		userResponses = append(userResponses, resp)
	}

	return &pb2.GetUsersResponse{Users: userResponses}, err
}

func (s *Server) GetUserDetails(ctx context.Context, req *pb2.GetUserDetailsRequest) (*pb2.GetUserDetailsResponse, error) {
	log.Println("GetUserDetails")
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	u, err := s.st.GetUserByID(ctx, uint(req.Id))
	if err != nil {
		if errors.Is(err, storage.ErrNothingReturned) {
			log.Println("nothing is returned")
			return &pb2.GetUserDetailsResponse{}, nil
		}
		log.Printf("some error: %v", err)
		return nil, err
	}

	return &pb2.GetUserDetailsResponse{
		Email:    u.Email,
		Username: u.Username,
		Password: u.Password,
	}, err
}

func (s *Server) CreateUpdateUser(ctx context.Context, req *pb2.CreateUpdateUserRequest) (*pb2.CreateUpdateUserResponse, error) {
	log.Println("CreateUpdateUser")
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if req.Operation == pb2.Operation_CREATE {
		log.Println("create")
		u := storage.User{
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		err := u.HashPassword()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		id, err := s.st.CreateUser(ctx, u)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return &pb2.CreateUpdateUserResponse{Id: uint64(id)}, err
	} else {
		log.Println("update")

		u := storage.User{
			ID:       uint(req.Id),
			Username: req.Username,
			Email:    req.Email,
			Password: req.Password,
		}

		err := u.HashPassword()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		err = s.st.UpdateUser(ctx, u)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return &pb2.CreateUpdateUserResponse{Id: req.Id}, err
	}
}
