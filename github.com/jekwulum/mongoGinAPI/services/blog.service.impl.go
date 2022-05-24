package services

// import (
// 	"errors"
// 	"html"
// 	"strings"
// 	"time"

// 	"github.com/jinzhu/gorm"
// 	"github.com/jekwulum/mongoGinAPI/models"
// )

// type BlogServiceImpl struct {
// 	blogdb	*models.Blog
// }

// func NewBlogService(blogdb *gorm.DB) BlogService {
// 	return &BlogServiceImpl{ blogdb: blogdb, }
// }


// func (b *BlogServiceImpl) Prepare() {
// 	b.blogdb.Title = html.EscapeString(strings.TrimSpace(b.blogdb.Title))
// 	b.blogdb.Content = html.EscapeString(strings.TrimSpace(b.blogdb.Content))
// 	b.blogdb.CreatedAt = time.Now()
// 	b.blogdb.UpdatedAt = time.Now()
// }

// func (b *BlogServiceImpl) Validate() map[string]string {
// 	var err error
// 	var errorMessages = make(map[string]string)

// 	if b.blogdb.Title == "" {
// 		err = errors.New("Required Title")
// 		errorMessages["Required_title"] = err.Error()
// 	}

// 	if b.blogdb.Content == "" {
// 		err = errors.New("Required Content")
// 		errorMessages["Required_content"] = err.Error()
// 	}

// 	return errorMessages
// }

// func (b *BlogServiceImpl) SaveBlog() (*models.Blog, error) {
// 	var err error
// 	err = b.blogdb.Debug().Model(&models.Blog{}).Create(&b.blogdb).Error
// 	if err != nil {
// 		return &models.Blog{}, err
// 	}
// 	return b.blogdb, nil
// }

// func (b *models.Blog) GetBlogs(db *gorm.DB) (*[]models.Blog, error) {
// 	var err error
// 	blogs := []models.Blog{}
// 	err = db.Debug().Model(&models.Blog{}).Find(&blogs).Error
// 	if err != nil {
// 		return &[]models.Blog{}, err
// 	}

// 	return &blogs, nil
// }

// func (b *models.Blog) GetBlog(db *gorm.DB, blog_id uint64) (*models.Blog, error) {
// 	var err error
// 	err = db.Debug().Model(&models.Blog{}).Where("id = ?", blog_id).Take(&b).Error
// 	if err != nil {
// 		return &models.Blog{}, err
// 	}
	
// 	return b, nil
// }

// func (b *models.Blog) UpdateBlog(db *gorm.DB) (*models.Blog, error) {
// 	var err error
// 	err = db.Debug().Model(&models.Blog{}).Where("id = ?", b.ID).Updates(models.Blog{
// 		Title: b.Title, Content: b.Content, UpdatedAt: time.Now()}).Error
// 	if err != nil {
// 		return &models.Blog{}, err
// 	}
// 	return b, nil
// }

// func (b *models.Blog) DeleteBlog(db *gorm.DB) (int64, error) {
// 	db = db.Debug().Model(&models.Blog{}).Where("id = ?", b.ID).Take(&models.Blog{}).Delete(&models.Blog{})
// 	if db.Error != nil {
// 		return 0, db.Error
// 	}
// 	return db.RowsAffected, nil
// }