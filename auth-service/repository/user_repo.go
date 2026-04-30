package repository

import "database/sql"

type User struct {
	ID       string
	Email    string
	Password string
	TenantID string
	Role     string
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *User) error {
	_, err := r.db.Exec(`
	INSERT INTO users (id, email, password, tenant_id, role)
	VALUES ($1, $2, $3, $4, $5)
	`, user.ID, user.Email, user.Password, user.TenantID, user.Role)

	return err
}

func (r *UserRepo) GetByMail() {

}
