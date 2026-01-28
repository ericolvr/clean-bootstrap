package domain

type Branch struct {
	ID           int    `json:"id" db:"id"`
	Client       string `json:"client" db:"client"`
	Name         string `json:"name" db:"name"`
	Uniorg       string `json:"uniorg" db:"uniorg"`
	Zipcode      string `json:"zipcode" db:"zipcode"`
	State        string `json:"state" db:"state"`
	City         string `json:"city" db:"city"`
	Neighborhood string `json:"neighborhood" db:"neighborhood"`
	Address      string `json:"address" db:"address"`
	Complement   string `json:"complement" db:"complement"`
}

type BranchRepository interface {
	Create(branch *Branch) error
	List(limit, offset int) ([]*Branch, error)
	GetByID(id int) (*Branch, error)
	GetByUniorg(uniorg string) (*Branch, error)
	Update(branch *Branch) error
	Delete(id int) error
}
