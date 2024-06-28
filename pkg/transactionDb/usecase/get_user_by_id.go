package usecase

import (
	"context"
	"time"
)

func (s *service) UpdateUserId(ctx context.Context, id int) error {

	tx := r.gormDB.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	user, err := s.userService.GetUserById(tx, id)
	if err != nil {
		return err
	}

	if user.Status != "active" {
		return nil
	}

	user.UpdatedAt = time.Now()

	err = s.userService.UpdateName(tx, name)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
