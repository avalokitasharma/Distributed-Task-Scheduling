package repository

import "database/sql"

type Tenant struct {
	ID   string
	Name string
}

type TenantRepo struct {
	db *sql.DB
}

func NewTenantRepo(db *sql.DB) *TenantRepo {
	return &TenantRepo{db: db}
}

func (r *TenantRepo) CreateTenant(t *Tenant) error {
	_, err := r.db.Exec(`
		INSERT INTO tenants (id, name)
		VALUES ($1, $2)
	`, t.ID, t.Name)

	return err
}

func (r *TenantRepo) GetByID(id string) (*Tenant, error) {
	t := &Tenant{}

	row := r.db.QueryRow(`
		SELECT id, name
		FROM tenants
		WHERE id = $1
	`, id)

	err := row.Scan(&t.ID, &t.Name)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// useful for signup dedupe - tenant names are unique
func (r *TenantRepo) GetByName(name string) (*Tenant, error) {
	t := &Tenant{}

	row := r.db.QueryRow(`
		SELECT id, name
		FROM tenants
		WHERE name = $1
	`, name)

	err := row.Scan(&t.ID, &t.Name)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TenantRepo) Exists(id string) (bool, error) {
	var exists bool

	err := r.db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM tenants where id=$1
		)
	`, id).Scan(&exists)

	return exists, err
}
