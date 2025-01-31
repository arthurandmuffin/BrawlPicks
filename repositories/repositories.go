package repositories

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) RepositoryTest() string {
	return "Repository"
}
