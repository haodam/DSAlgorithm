package redis

type Option func(*redis)

func ConnPoolSize(poolSize int) Option {
	return func(r *redis) {
		r.poolSize = poolSize
	}
}
func ConnPassword(password string) Option {
	return func(r *redis) {
		r.password = password
	}
}
func ConnDataBase(db int) Option {
	return func(r *redis) {
		r.database = db
	}
}
