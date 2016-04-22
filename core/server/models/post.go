package models

import (
	"path"
	"time"
	"unicode/utf8"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"github.com/go-ant/ant/core/server/modules/setting"
	"github.com/go-ant/ant/core/server/modules/utils"
	"github.com/go-ant/ant/core/server/modules/utils/slug"
)

type Post struct {
	Id              uint32    `json:"id" gorm:"primary_key"`
	Title           string    `json:"title" sql:"not null;type:varchar(150);"`
	Slug            string    `json:"slug" sql:"not null;type:varchar(200);"`
	Markdown        string    `json:"markdown" sql:"not null;type:text;"`
	Html            string    `json:"html" sql:"not null;type:text;"`
	Cover           string    `json:"cover" sql:"type:varchar(150)"`
	Language        string    `json:"language" sql:"type:varchar(6)"`
	Page            bool      `json:"page" sql:"not null;type:tinyint(1)"`
	Featured        bool      `json:"featured" sql:"not null;type:tinyint(1)"`
	Status          string    `json:"status" sql:"not null;type:varchar(20)"`
	MetaTitle       string    `json:"meta_title" sql:"type:varchar(150)"`
	MetaDescription string    `json:"meta_description" sql:"type:varchar(200)"`
	AuthorId        uint32    `json:"-" sql:"not null;type:bigint unsigned;"`
	PublishedAt     time.Time `json:"published_at" sql:"not null;type:datetime;"`
	PublishedBy     uint32    `json:"-" sql:"not null;type:bigint unsigned;"`
	CreatedAt       time.Time `json:"-" sql:"not null;type:datetime;"`
	CreatedBy       uint32    `json:"-" sql:"not null;type:bigint unsigned;"`
	UpdatedAt       time.Time `json:"-" sql:"type:datetime;"`
	UpdatedBy       uint32    `json:"-" sql:"type:bigint unsigned;"`

	Author *User `json:"author,omitempty"`
}

func (c *Post) GetAuthor() {
	c.Author = new(User)
	db.Model(&c).Related(&c.Author, "author_id")
}

func (c *Post) GetURL() string {
	return path.Join(setting.Host.Path, "post/", c.Slug)
}

func (c *Post) GetPageMeta() map[string]string {
	metaTitle := c.MetaTitle
	metaDesc := c.MetaDescription

	if utils.IsEmpty(metaTitle) {
		metaTitle = c.Title
	}

	if utils.IsEmpty(metaDesc) {
		metaDesc = utils.StripTagsFromHtml(c.Html)
		metaDesc = utils.SubString(metaDesc, 150)

	}

	return map[string]string{"Title": metaTitle, "Description": metaDesc}
}

const (
	PostStatusPublished string = "published"
	PostStatusDraft     string = "draft"
)

func CreatePost(post *Post, opts *Options) ApiErr {
	if utf8.RuneCountInString(post.Title) > 150 {
		return ApiMsg.ErrPostTitleTooLong
	}
	if utf8.RuneCountInString(post.MetaTitle) > 150 {
		return ApiMsg.ErrPostTitleTooLong
	}
	if utf8.RuneCountInString(post.MetaDescription) > 200 {
		return ApiMsg.ErrPostTitleTooLong
	}

	if utils.IsEmpty(post.Title) {
		post.Title = "(Untitled)"
	}
	if utils.IsEmpty(post.Slug) {
		post.Slug = slug.Make(post.Title)
	}
	if utf8.RuneCountInString(post.Slug) > 200 {
		return ApiMsg.ErrPostSlugTooLong
	}

	duplicate := &Post{}
	db.Select("id, slug").Where("slug = ?", post.Slug).First(duplicate)
	if duplicate.Id > 0 {
		return ApiMsg.ErrPostSlugAlreadyExist
	}

	if !utils.IsEmpty(post.Markdown) {
		htmlPolicy := bluemonday.UGCPolicy()
		post.Html = htmlPolicy.Sanitize(string(blackfriday.MarkdownCommon([]byte(post.Markdown))))
	}

	if post.Status != PostStatusPublished {
		post.Status = PostStatusDraft
	}

	tx := db.Begin()
	if err := tx.Create(post).Error; err != nil {
		tx.Rollback()
		return UnknowError(err.Error())
	}
	tx.Commit()

	return ApiMsg.Created
}

func EditPost(post *Post, opts *Options) ApiErr {
	if utf8.RuneCountInString(post.Title) > 150 {
		return ApiMsg.ErrPostTitleTooLong
	}
	if utf8.RuneCountInString(post.MetaTitle) > 150 {
		return ApiMsg.ErrPostTitleTooLong
	}
	if utf8.RuneCountInString(post.MetaDescription) > 200 {
		return ApiMsg.ErrPostTitleTooLong
	}

	post.Html = ""
	post.UpdatedAt = time.Now()

	if utils.IsEmpty(post.Title) {
		post.Title = "(Untitled)"
	}
	if utils.IsEmpty(post.Slug) {
		post.Slug = slug.Make(post.Title)
	}
	if utf8.RuneCountInString(post.Slug) > 200 {
		return ApiMsg.ErrPostSlugTooLong
	}

	duplicate := &Post{}
	db.Select("id, slug").Where("slug = ?", post.Slug).First(duplicate)
	if duplicate.Id > 0 && post.Id != duplicate.Id {
		return ApiMsg.ErrPostSlugAlreadyExist
	}

	if !utils.IsEmpty(post.Markdown) {
		htmlPolicy := bluemonday.UGCPolicy()
		post.Html = htmlPolicy.Sanitize(string(blackfriday.MarkdownCommon([]byte(post.Markdown))))
	}

	tx := db.Begin()
	if err := tx.Select([]string{"title", "slug", "markdown", "html", "cover", "language", "page", "featured", "status", "meta_title", "meta_description", "author_id", "published_at", "published_by", "updated_at", "updated_by"}).Save(post).Error; err != nil {
		tx.Rollback()
		return UnknowError(err.Error())
	}
	tx.Commit()

	return ApiMsg.Saved
}

func DeletePost(id uint32, opts *Options) ApiErr {
	tx := db.Begin()
	tx.Where("id = ?", id).Delete(Post{})
	tx.Commit()

	return ApiMsg.Deleted
}

func GetPost(opts *Options) (*Post, ApiErr) {
	dbInit := initDb(opts)
	post := new(Post)
	if u := dbInit.First(post); u.Error != nil {
		return nil, ApiMsg.ErrPostNotFound
	}

	if opts != nil && opts.IsInclude("author") && post.Id > 0 {
		post.GetAuthor()
	}
	return post, ApiMsg.Success
}

func GetPostById(id uint32, opts *Options) (*Post, ApiErr) {
	if opts == nil {
		opts = &Options{}
	}
	opts.GormAdp = &GormAdapter{
		Query: "id = ?",
		Args:  []interface{}{id},
	}
	post, errApi := GetPost(opts)

	return post, errApi
}

func GetPosts(opts *Options) ([]*Post, uint32, ApiErr) {
	var recordCount uint32
	dbInit := initDb(opts)
	posts := make([]*Post, 0)
	if err := dbInit.Model(Post{}).Count(&recordCount).Limit(opts.Limit).Offset(opts.Offset).Find(&posts).Error; err != nil {
		return nil, 0, UnknowError(err.Error())
	}
	if len(posts) > 0 && opts.IsInclude("user") {
		for _, post := range posts {
			post.GetAuthor()
		}
	}
	return posts, recordCount, ApiMsg.Success
}
