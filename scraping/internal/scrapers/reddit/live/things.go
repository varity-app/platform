package live

import (
	"encoding/json"
	"fmt"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// Everything in this file is to implement fetching new comments for a given subreddit
// inside the RedditCommentScraper

const (
	kindComment           = "t1"
	kindUser              = "t2"
	kindPost              = "t3"
	kindMessage           = "t4"
	kindSubreddit         = "t5"
	kindTrophy            = "t6"
	kindListing           = "Listing"
	kindSubredditSettings = "subreddit_settings"
	kindKarmaList         = "KarmaList"
	kindTrophyList        = "TrophyList"
	kindUserList          = "UserList"
	kindMore              = "more"
	kindLiveThread        = "LiveUpdateEvent"
	kindLiveThreadUpdate  = "LiveUpdate"
	kindModAction         = "modaction"
	kindMulti             = "LabeledMulti"
	kindMultiDescription  = "LabeledMultiDescription"
	kindWikiPage          = "wikipage"
	kindWikiPageListing   = "wikipagelisting"
	kindWikiPageSettings  = "wikipagesettings"
	kindStyleSheet        = "stylesheet"
)

type anchor interface {
	After() string
}

// thing is an entity on Reddit.
// Its kind reprsents what it is and what is stored in the Data field.
// e.g. t1 = comment, t2 = user, t3 = post, etc.
type thing struct {
	Kind string      `json:"kind"`
	Data interface{} `json:"data"`
}

func (t *thing) After() string {
	if t == nil {
		return ""
	}
	a, ok := t.Data.(anchor)
	if !ok {
		return ""
	}
	return a.After()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *thing) UnmarshalJSON(b []byte) error {
	root := new(struct {
		Kind string          `json:"kind"`
		Data json.RawMessage `json:"data"`
	})

	err := json.Unmarshal(b, root)
	if err != nil {
		return err
	}

	t.Kind = root.Kind
	var v interface{}

	switch t.Kind {
	case kindListing:
		v = new(listing)
	case kindComment:
		v = new(reddit.Comment)
	case kindMore:
		v = new(reddit.More)
	case kindUser:
		v = new(reddit.User)
	case kindPost:
		v = new(reddit.Post)
	case kindSubreddit:
		v = new(reddit.Subreddit)
	default:
		return fmt.Errorf("unrecognized kind: %q", t.Kind)
	}

	err = json.Unmarshal(root.Data, v)
	if err != nil {
		return err
	}

	t.Data = v
	return nil
}

func (t *thing) Listing() (v *listing, ok bool) {
	v, ok = t.Data.(*listing)
	return
}

func (t *thing) Comment() (v *reddit.Comment, ok bool) {
	v, ok = t.Data.(*reddit.Comment)
	return
}

func (t *thing) More() (v *reddit.More, ok bool) {
	v, ok = t.Data.(*reddit.More)
	return
}

func (t *thing) User() (v *reddit.User, ok bool) {
	v, ok = t.Data.(*reddit.User)
	return
}

func (t *thing) Post() (v *reddit.Post, ok bool) {
	v, ok = t.Data.(*reddit.Post)
	return
}

func (t *thing) Subreddit() (v *reddit.Subreddit, ok bool) {
	v, ok = t.Data.(*reddit.Subreddit)
	return
}

// listing is a list of things coming from the Reddit API.
// It also contains the after anchor useful to get the next results via subsequent requests.
type listing struct {
	things things
	after  string
}

func (l *listing) After() string {
	return l.after
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *listing) UnmarshalJSON(b []byte) error {
	root := new(struct {
		Things things `json:"children"`
		After  string `json:"after"`
	})

	err := json.Unmarshal(b, root)
	if err != nil {
		return err
	}

	l.things = root.Things
	l.after = root.After

	return nil
}

func (l *listing) Comments() []*reddit.Comment {
	if l == nil {
		return nil
	}
	return l.things.Comments
}

func (l *listing) Mores() []*reddit.More {
	if l == nil {
		return nil
	}
	return l.things.Mores
}

func (l *listing) Users() []*reddit.User {
	if l == nil {
		return nil
	}
	return l.things.Users
}

func (l *listing) Posts() []*reddit.Post {
	if l == nil {
		return nil
	}
	return l.things.Posts
}

func (l *listing) Subreddits() []*reddit.Subreddit {
	if l == nil {
		return nil
	}
	return l.things.Subreddits
}

type things struct {
	Comments   []*reddit.Comment
	Mores      []*reddit.More
	Users      []*reddit.User
	Posts      []*reddit.Post
	Subreddits []*reddit.Subreddit
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *things) UnmarshalJSON(b []byte) error {
	var things []thing
	if err := json.Unmarshal(b, &things); err != nil {
		return err
	}

	t.add(things...)
	return nil
}

func (t *things) add(things ...thing) {
	for _, thing := range things {
		switch v := thing.Data.(type) {
		case *reddit.Comment:
			t.Comments = append(t.Comments, v)
		case *reddit.More:
			t.Mores = append(t.Mores, v)
		case *reddit.User:
			t.Users = append(t.Users, v)
		case *reddit.Post:
			t.Posts = append(t.Posts, v)
		case *reddit.Subreddit:
			t.Subreddits = append(t.Subreddits, v)
		}
	}
}
