package repository

//type userRepository struct {
//}
//
//func (r *userRepository) GetUserById(tx *gorm.DB, id int) (entity.User, error) {
//	var (
//		user entity.User
//		err  error
//	)
//	err = tx.Table(user.TableName()).Where("id = ?", id).First(&user).Error
//	if err != nil {
//		return user, err
//	}
//	return user, nil
//}

//type PaginationQuery struct {
//	CountSQL    string
//	QuerySQL    string
//	Params      []interface{}
//	Page        int
//	Take        int
//	ItemMappers func(*sql.Rows) interface{}
//}
//
//func (r *Repository[T]) PaginationQuery(ctx context.Context, option *PaginationQuery) (int, *[]interface{}, *AppError) {
//	var countValue int64
//
//	err := r.QueryBuilder(ctx).Raw(option.CountSQL, option.Params...).First(&countValue).Error
//
//	if err != nil {
//		if !errors.Is(err, gorm.ErrRecordNotFound) {
//			return 0, nil, QueryInvalid(err.Error())
//		}
//	}
//
//	totalItems := int(countValue)
//	take := option.Take
//	offset := (option.Page - 1) * option.Take
//	items := make([]interface{}, 0)
//
//	if totalItems == 0 || totalItems < offset {
//		return 0, &items, nil
//	}
//
//	rows, err := r.QueryBuilder(ctx).Raw(option.QuerySQL+fmt.Sprintf(" LIMIT %v OFFSET %v", take, offset), option.Params...).Rows()
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			return 0, &items, nil
//		}
//		return 0, &items, QueryInvalid(err.Error())
//	}
//
//	defer rows.Close()
//	for rows.Next() {
//		it := option.ItemMappers(rows)
//		items = append(items, it)
//	}
//
//	return totalItems, &items, nil
//}
