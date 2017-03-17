package stats

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBucketHttpRequest_getRequestOperations(t *testing.T) {
	dataProvider := []struct {
		Method     string
		Path       string
		Operations MetricOperation
	}{
		{"GET", "/", MetricOperation{"get", MetricEmptyPlaceholder, MetricEmptyPlaceholder}},
		{"TRACE", "/api", MetricOperation{"trace", "api", MetricEmptyPlaceholder}},
		{"TRACE", "/api/", MetricOperation{"trace", "api", MetricEmptyPlaceholder}},
		{"POST", "/api/recipes", MetricOperation{"post", "api", "recipes"}},
		{"POST", "/api/recipes/", MetricOperation{"post", "api", "recipes"}},
		{"DELETE", "/api/recipes/123", MetricOperation{"delete", "api", "recipes"}},
		{"DELETE", "/api/recipes.foo-bar/123", MetricOperation{"delete", "api", "recipes.foo-bar"}},
		{"DELETE", "/api/recipes.foo_bar/123", MetricOperation{"delete", "api", "recipes.foo_bar"}},
		// paths withs IDs at the path second level
		{"GET", "/user/qwerty", MetricOperation{"get", "user", MetricIDPlaceholder}},
		{"GET", "/users/qwerty", MetricOperation{"get", "users", MetricIDPlaceholder}},
		{"GET", "/allergens/foobarbaz", MetricOperation{"get", "allergens", MetricIDPlaceholder}},
		{"GET", "/cuisines/foobarbaz", MetricOperation{"get", "cuisines", MetricIDPlaceholder}},
		{"GET", "/favorites/foobarbaz", MetricOperation{"get", "favorites", MetricIDPlaceholder}},
		{"GET", "/ingredients/foobarbaz", MetricOperation{"get", "ingredients", MetricIDPlaceholder}},
		{"GET", "/menus/foobarbaz", MetricOperation{"get", "menus", MetricIDPlaceholder}},
		{"GET", "/ratings/foobarbaz", MetricOperation{"get", "ratings", MetricIDPlaceholder}},
		{"GET", "/recipes/foobarbaz", MetricOperation{"get", "recipes", MetricIDPlaceholder}},
		{"GET", "/addresses/foobarbaz", MetricOperation{"get", "addresses", MetricIDPlaceholder}},
		{"GET", "/boxes/foobarbaz", MetricOperation{"get", "boxes", MetricIDPlaceholder}},
		{"GET", "/coupons/foobarbaz", MetricOperation{"get", "coupons", MetricIDPlaceholder}},
		{"GET", "/customers/foobarbaz", MetricOperation{"get", "customers", MetricIDPlaceholder}},
		{"GET", "/delivery_options/foobarbaz", MetricOperation{"get", "delivery_options", MetricIDPlaceholder}},
		{"GET", "/product_families/foobarbaz", MetricOperation{"get", "product_families", MetricIDPlaceholder}},
		{"GET", "/products/foobarbaz", MetricOperation{"get", "products", MetricIDPlaceholder}},
		{"GET", "/recipients/foobarbaz", MetricOperation{"get", "recipients", MetricIDPlaceholder}},
		// path may have either numeric ID or non-numeric trackable path
		{"GET", "/subscriptions/12345", MetricOperation{"get", "subscriptions", MetricIDPlaceholder}},
		{"GET", "/subscriptions/search", MetricOperation{"get", "subscriptions", "search"}},
		{"GET", "/freebies/12345", MetricOperation{"get", "freebies", MetricIDPlaceholder}},
		{"GET", "/freebies/search", MetricOperation{"get", "freebies", "search"}},
		// path may be short or full
		{"GET", "/clients", MetricOperation{"get", "clients", MetricEmptyPlaceholder}},
		{"GET", "/clients/qwe123", MetricOperation{"get", "clients", MetricIDPlaceholder}},
	}

	callback := NewHasIDAtSecondLevelCallback(map[PathSection]SectionTestDefinition{
		"addresses":        {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"allergens":        {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"boxes":            {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"clients":          {SectionTestIsNotEmpty, GetSectionTestCallback(SectionTestIsNotEmpty)},
		"coupons":          {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"cuisines":         {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"customers":        {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"delivery_options": {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"favorites":        {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"freebies":         {SectionTestIsNumeric, GetSectionTestCallback(SectionTestIsNumeric)},
		"ingredients":      {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"menus":            {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"product_families": {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"products":         {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"ratings":          {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"recipes":          {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"recipients":       {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"subscriptions":    {SectionTestIsNumeric, GetSectionTestCallback(SectionTestIsNumeric)},
		"user":             {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
		"users":            {SectionTestTrue, GetSectionTestCallback(SectionTestTrue)},
	})

	for _, data := range dataProvider {
		r := &http.Request{Method: data.Method, URL: &url.URL{Path: data.Path}}
		assert.Equal(t, data.Operations, getRequestOperations(r, callback))
	}
}

func TestBucketHttpRequest_MetricNameAlterCallback(t *testing.T) {
	dataProvider := []struct {
		Method     string
		Path       string
		Operations MetricOperation
		Query      string
	}{
		{"GET", "/users/qwerty", MetricOperation{"get", "users", MetricIDPlaceholder}, ""},
		{"GET", "/clients", MetricOperation{"get", "clients", MetricEmptyPlaceholder}, ""},
		{"GET", "/clients/qwe123", MetricOperation{"get", "clients", MetricIDPlaceholder}, ""},
		{"GET", "/token/revoke", MetricOperation{"get", "token", "revoke"}, ""},
		{"GET", "/token/revoke", MetricOperation{"get", "token", "revoke"}, "foo=bar&grant_type=baz"},
		{"GET", "/token", MetricOperation{"get", "token", "baz"}, "foo=bar&grant_type=baz"},
		{"GET", "/token", MetricOperation{"get", "token", MetricEmptyPlaceholder}, "foo=bar"},
		{"GET", "/token/client_credentials", MetricOperation{"get", "token", "client_credentials"}, ""},
	}

	callback := func(metricFragments MetricOperation, r *http.Request) MetricOperation {
		if metricFragments[1] == "token" && metricFragments[2] != "revoke" {
			grantType := r.URL.Query().Get("grant_type")
			if grantType != "" {
				metricFragments[2] = grantType
			}
			return metricFragments
		}

		return NewHasIDAtSecondLevelCallback(map[PathSection]SectionTestDefinition{
			"users":   {SectionTestIsNotEmpty, GetSectionTestCallback(SectionTestIsNotEmpty)},
			"clients": {SectionTestIsNotEmpty, GetSectionTestCallback(SectionTestIsNotEmpty)},
		})(metricFragments, r)
	}

	for _, data := range dataProvider {
		r := &http.Request{Method: data.Method, URL: &url.URL{Path: data.Path, RawQuery: data.Query}}
		assert.Equal(t, data.Operations, getRequestOperations(r, callback))
	}
}

func TestBucketHttpRequest_Metric(t *testing.T) {
	dataProvider := []struct {
		Method  string
		Path    string
		Success bool
		Metric  string
	}{
		{"GET", "/foo/bar/baz", true, "request.get.foo.bar"},
		{"GET", "/foo/bar/baz", false, "request.get.foo.bar"},
		{"GET", "/token/client_credentials", false, "request.get.token.client__credentials"},
		{"GET", "/delivery_options/foobarbaz", true, "request.get.delivery__options.foobarbaz"},
		{"GET", "/product_families/foobarbaz", true, "request.get.product__families.foobarbaz"},
		{"DELETE", "/api/recipes.foo-bar/123", true, "request.delete.api.recipes_foo-bar"},
		{"DELETE", "/api/recipes.foo_bar/123", true, "request.delete.api.recipes_foo__bar"},
	}

	for _, data := range dataProvider {
		r := &http.Request{Method: data.Method, URL: &url.URL{Path: data.Path}}
		b := NewBucketHttpRequest(sectionRequest, r, data.Success, nil)
		assert.Equal(t, data.Metric, b.Metric())
	}
}

func TestBucketHttpRequest_MetricWithSuffix(t *testing.T) {
	dataProvider := []struct {
		Method  string
		Path    string
		Success bool
		Metric  string
	}{
		{"GET", "/foo/bar/baz", true, "request-ok.get.foo.bar"},
		{"GET", "/foo/bar/baz", false, "request-fail.get.foo.bar"},
	}

	for _, data := range dataProvider {
		r := &http.Request{Method: data.Method, URL: &url.URL{Path: data.Path}}
		b := NewBucketHttpRequest(sectionRequest, r, data.Success, nil)
		assert.Equal(t, data.Metric, b.MetricWithSuffix())
	}
}

func TestBucketHttpRequest_MetricTotal(t *testing.T) {
	dataProvider := []struct {
		Method  string
		Path    string
		Success bool
		Metric  string
	}{
		{"GET", "/foo/bar/baz", true, "total.request"},
		{"GET", "/foo/bar/baz", false, "total.request"},
	}

	for _, data := range dataProvider {
		r := &http.Request{Method: data.Method, URL: &url.URL{Path: data.Path}}
		b := NewBucketHttpRequest(sectionRequest, r, data.Success, nil)
		assert.Equal(t, data.Metric, b.MetricTotal())
	}
}

func TestBucketHttpRequest_TMetricTotalWithSuffix(t *testing.T) {
	dataProvider := []struct {
		Method  string
		Path    string
		Success bool
		Metric  string
	}{
		{"GET", "/foo/bar/baz", true, "total.request-ok"},
		{"GET", "/foo/bar/baz", false, "total.request-fail"},
	}

	for _, data := range dataProvider {
		r := &http.Request{Method: data.Method, URL: &url.URL{Path: data.Path}}
		b := NewBucketHttpRequest(sectionRequest, r, data.Success, nil)
		assert.Equal(t, data.Metric, b.MetricTotalWithSuffix())
	}
}

func TestBucketHttpRequest_Metric_customSection(t *testing.T) {
	section := "section111"
	dataProvider := []struct {
		Method  string
		Path    string
		Success bool
		Metric  string
	}{
		{"GET", "/foo/bar/baz", true, section + ".get.foo.bar"},
		{"GET", "/foo/bar/baz", false, section + ".get.foo.bar"},
	}

	for _, data := range dataProvider {
		r := &http.Request{Method: data.Method, URL: &url.URL{Path: data.Path}}
		b := NewBucketHttpRequest(section, r, data.Success, nil)
		assert.Equal(t, data.Metric, b.Metric())
	}
}

func TestBucketHttpRequest_MetricWithSuffix_customSection(t *testing.T) {
	section := "section111"
	dataProvider := []struct {
		Method  string
		Path    string
		Success bool
		Metric  string
	}{
		{"GET", "/foo/bar/baz", true, section + "-ok.get.foo.bar"},
		{"GET", "/foo/bar/baz", false, section + "-fail.get.foo.bar"},
	}

	for _, data := range dataProvider {
		r := &http.Request{Method: data.Method, URL: &url.URL{Path: data.Path}}
		b := NewBucketHttpRequest(section, r, data.Success, nil)
		assert.Equal(t, data.Metric, b.MetricWithSuffix())
	}
}
