package users

import "context"

type Service struct {
	Repo *Repository
}

func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (string, error) {

	id, err := s.Repo.Create(ctx, req.Email, req.Name)
	if err != nil {
		return "", err
	}

	return id, nil
}
