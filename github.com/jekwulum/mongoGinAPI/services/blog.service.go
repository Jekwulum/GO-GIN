package services

import (
	// "fmt"

	"github.com/jekwulum/mongoGinAPI/models"

	"github.com/jinzhu/gorm"
)


func GetBlogs(db *gorm.DB) ([]models.Blog, error) {
	blogs := []models.Blog{}
	// query := db.Select("blogs.*")
	if err := db.Find(&blogs).Error; err != nil {
		return blogs, err
	}
	return blogs, nil
	// var blogs []models.Blog
	// db.Find(&blogs)
	// fmt.Println("blogs------------------------------------> ", blogs)
	// return blogs, nil
}

func GetBlog(id uint64, db *gorm.DB) (models.Blog, bool, error) {
	b := models.Blog{}
	query := db.Select("blogs.*")
	err := query.Where("blogs.id = ?", id).First(&b).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return b, false, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return b, false, nil
	}
	return b, true, nil
}

func UpdateBlog(db *gorm.DB, b *models.Blog) error {
	if err := db.Save(&b).Error; err != nil {
		return err
	}
	return nil
}

func DeleteBlog(id uint64, db *gorm.DB) error {
	var b models.Blog
	if err := db.Where("id = ? ", id).Delete(&b).Error; err != nil {
		return err
	}
	return nil
}

