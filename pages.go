package crowi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/context"
)

// PagesService handles communication with the Pages related
// methods of the Crowi API.
type PagesService service

// Create makes a page in your Crowi. The request requires
// the path and page content used for the page name
func (s *PagesService) Create(ctx context.Context, path, body string) (*Page, error) {
	var page Page
	params := url.Values{}
	params.Set("access_token", s.client.config.Token)
	params.Set("path", path)
	params.Set("body", body)
	err := s.client.newRequest(ctx, http.MethodPost, "/_api/pages.create", params, &page)
	if err != nil {
		return nil, err
	}
	if err := fillPageInfo(&page.Page); err != nil {
		return nil, err
	}
	return &page, nil
}

// Update updates the page content. A page_id is necessary to know which
// page should be updated.
func (s *PagesService) Update(ctx context.Context, id, revisionID, body string) (*Page, error) {
	var page Page
	params := url.Values{}
	params.Set("access_token", s.client.config.Token)
	params.Set("page_id", id)
	params.Set("revision_id", revisionID)
	params.Set("body", body)
	err := s.client.newRequest(ctx, http.MethodPost, "/_api/pages.update", params, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// PagesListOptions specifies the optional parameters to the
// PagesService.List methods.
type PagesListOptions struct {
	ListOptions
}

// List shows the list of pages. A parameter of path or user is required.
func (s *PagesService) List(ctx context.Context, path, user string, opt *PagesListOptions) (*Pages, error) {
	var pages Pages
	params := url.Values{}
	params.Set("access_token", s.client.config.Token)
	params.Set("path", path)
	params.Set("user", user)
	err := s.client.newRequest(ctx, http.MethodGet, "/_api/pages.list", params, &pages)
	if err != nil {
		return nil, err
	}
	if opt != nil && opt.ListOptions.Pagenation {
		offset := 0
		var p []PageInfo
		for {
			params.Set("offset", fmt.Sprintf("%d", offset))
			err := s.client.newRequest(ctx, http.MethodGet, "/_api/pages.list", params, &pages)
			if err != nil {
				break
			}
			p = append(p, pages.Pages...)
			offset += 50
		}

		for i := range p {
			if err := fillPageInfo(&p[i]); err != nil {
				return nil, err
			}
		}

		pages.Pages = p
	}
	return &pages, nil
}

func (s *PagesService) Get(ctx context.Context, path string) (*Page, error) {
	var page Page
	params := url.Values{}
	params.Set("access_token", s.client.config.Token)
	params.Set("path", path)
	err := s.client.newRequest(ctx, http.MethodGet, "/_api/pages.get", params, &page)
	if err != nil {
		return nil, err
	}
	if err := fillPageInfo(&page.Page); err != nil {
		return &page, err
	}
	return &page, nil
}

func fillPageInfo(p *PageInfo) error {
	if p.RawRevision == nil {
		return nil
	}
	if err := json.Unmarshal(p.RawRevision, &p.Revision); err == nil {
		p.RevisionID = p.Revision.ID
		return nil
	}
	return json.Unmarshal(p.RawRevision, &p.RevisionID)
}

type Page struct {
	Page  PageInfo `json:"page"`
	OK    bool     `json:"ok"`
	Error string   `json:"error"`
}

type PageInfo struct {
	RawRevision    json.RawMessage `json:"revision"`
	RevisionID     string          `json:"-"`
	Revision       PageRevision    `json:"-"`
	V              int             `json:"__v"`
	RedirectTo     interface{}     `json:"redirectTo"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	LastUpdateUser interface{}     `json:"lastUpdateUser"`
	Creator        interface{}     `json:"creator"`
	Path           string          `json:"path"`
	ID             string          `json:"_id"`
	CreatedAt      time.Time       `json:"createdAt"`
	CommentCount   int             `json:"commentCount"`
	SeenUsers      []interface{}   `json:"seenUsers"`
	Liker          []interface{}   `json:"liker"`
	GrantedUsers   []interface{}   `json:"grantedUsers"`
	Grant          int             `json:"grant"`
	Status         string          `json:"status"`
	Extended       PageExtended    `json:"extended,omitempty"`
}

type PageRevision struct {
	ID        string      `json:"_id"`
	Author    interface{} `json:"author"`
	Body      string      `json:"body"`
	Path      string      `json:"path"`
	V         int         `json:"__v"`
	CreatedAt time.Time   `json:"createdAt"`
	Format    string      `json:"format"`
}

type PageExtended struct {
	Slack string `json:"slack"`
}

type Pages struct {
	Pages []PageInfo `json:"pages"`
	OK    bool       `json:"ok"`
	Error string     `json:"error"`
}
