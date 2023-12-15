package handling

import (
	"context"
	"fmt"
	"strings"
)

type Service struct {
}

func (s *Service) Signup(ctx context.Context, email, password string) error {
	if !strings.Contains(email, "@") {
		return fmt.Errorf("invalid email: %w", ErrBadRequest)
	}
	return nil
}
