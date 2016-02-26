package templates

import (
	"github.com/go-ant/ant/core/server/modules/startup"
	"gopkg.in/flosch/pongo2.v3"
)

const (
	keyNoUser    string = "noUser"
	keyNoPosts   string = "noPost"
	keyPostCount string = "postCount"
)

type page struct {
	Count uint32
	Prev  uint32
	Next  uint32
}

func init() {
	startup.Register(func() {
		pongo2.RegisterTag("asset", tagAssetParser)
		pongo2.RegisterTag("ant_head", tagHeadCodeInjectionParser)
		pongo2.RegisterTag("ant_foot", tagFootCodeInjectionParser)
		pongo2.RegisterTag("body_class", tagBodyClassParser)
		pongo2.RegisterTag("excerpt", tagExcerptParser)
		pongo2.RegisterTag("post_class", tagPostClassParser)
		pongo2.RegisterTag("page_url", tagPageUrlParser)
		pongo2.RegisterTag("prev_post", tagPrevPostParser)
		pongo2.RegisterTag("next_post", tagNextPostParser)
		pongo2.RegisterTag("posts", tagPostsParser)
		pongo2.RegisterTag("user", tagUserParser)
	})
}
