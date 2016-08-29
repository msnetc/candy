package store

import (
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/dearcode/candy/meta"
)

// Store save user, message.
type Store struct {
	host    string
	dbPath  string
	user    *userDB
	message *messageDB
}

// NewStore new Store server.
func NewStore(host, dbPath string) *Store {
	return &Store{host: host, dbPath: dbPath, user: newUserDB(dbPath), message: newMessageDB(dbPath)}
}

// Start start service.
func (s *Store) Start() error {
	serv := grpc.NewServer()
	meta.RegisterStoreServer(serv, s)

	lis, err := net.Listen("tcp", s.host)
	if err != nil {
		return err
	}

	if err = s.user.start(); err != nil {
		return err
	}

	if err = s.message.start(); err != nil {
		return err
	}

	return serv.Serve(lis)
}

// Register add user.
func (s *Store) Register(_ context.Context, req *meta.RegisterRequest) (*meta.RegisterResponse, error) {
	if err := s.user.register(req.User, req.Password, req.ID); err != nil {
		return &meta.RegisterResponse{Header: &meta.ResponseHeader{Code: -1, Msg: err.Error()}}, nil
	}

	return &meta.RegisterResponse{}, nil
}

// Auth check password.
func (s *Store) Auth(_ context.Context, req *meta.AuthRequest) (*meta.AuthResponse, error) {
	id, err := s.user.auth(req.User, req.Password)
	if err != nil {
		return &meta.AuthResponse{Header: &meta.ResponseHeader{Code: -1, Msg: err.Error()}}, nil
	}
	return &meta.AuthResponse{ID: id}, nil
}

// FindUser 根据字符串的用户名查的用户信息.
func (s *Store) FindUser(_ context.Context, req *meta.FindUserRequest) (*meta.FindUserResponse, error) {
	id, err := s.user.findUser(req.User)
	if err != nil {
		return &meta.FindUserResponse{Header: &meta.ResponseHeader{Code: -1, Msg: err.Error()}}, nil
	}
	return &meta.FindUserResponse{ID: id}, nil
}

// AddFriend 添加好友，两人都添加过对方后才可以聊天.
func (s *Store) AddFriend(_ context.Context, req *meta.AddFriendRequest) (*meta.AddFriendResponse, error) {
	if err := s.user.friend.add(req.From, req.To, req.Confirm); err != nil {
		return &meta.AddFriendResponse{Header: &meta.ResponseHeader{Code: -1, Msg: err.Error()}}, nil
	}

	if req.Confirm {
		if err := s.user.friend.add(req.To, req.From, req.Confirm); err != nil {
			return &meta.AddFriendResponse{Header: &meta.ResponseHeader{Code: -1, Msg: err.Error()}}, nil
		}
		return &meta.AddFriendResponse{Confirm: true}, nil
	}

	return &meta.AddFriendResponse{}, nil
}

func (s *Store) CreateGroup(_ context.Context, req *meta.CreateGroupRequest) (*meta.CreateGroupResponse, error) {
	return nil, nil
}

func (s *Store) NewMessage(_ context.Context, req *meta.NewMessageRequest) (*meta.NewMessageResponse, error) {
	return nil, nil
}
