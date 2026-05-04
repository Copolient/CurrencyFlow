package repository

import (
	"exchangeapp/internal/model"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(article *model.Article) error
	FindAll() ([]model.Article, error)
	FindByID(id string) (*model.Article, error)
}

type articleRepo struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepo{db: db}
}

func (r *articleRepo) Create(article *model.Article) error {
	return r.db.Create(article).Error
}

func (r *articleRepo) FindAll() ([]model.Article, error) {
	var articles []model.Article
	err := r.db.Find(&articles).Error
	return articles, err
}

func (r *articleRepo) FindByID(id string) (*model.Article, error) {
	var article model.Article
	err := r.db.Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}
