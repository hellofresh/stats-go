package bucket

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/hellofresh/stats-go/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpRequest_BuildHTTPRequestMetricOperation(t *testing.T) {
	dataProvider := []struct {
		Method     string
		Path       string
		Operations MetricOperation
	}{
		{"GET", "/", NewMetricOperation("get", MetricEmptyPlaceholder, MetricEmptyPlaceholder)},
		{"TRACE", "/api", NewMetricOperation("trace", "api", MetricEmptyPlaceholder)},
		{"TRACE", "/api/", NewMetricOperation("trace", "api", MetricEmptyPlaceholder)},
		{"POST", "/api/recipes", NewMetricOperation("post", "api", "recipes")},
		{"POST", "/api/recipes/", NewMetricOperation("post", "api", "recipes")},
		{"DELETE", "/api/recipes/123", NewMetricOperation("delete", "api", "recipes")},
		{"DELETE", "/api/recipes.foo-bar/123", NewMetricOperation("delete", "api", "recipes.foo-bar")},
		{"DELETE", "/api/recipes.foo_bar/123", NewMetricOperation("delete", "api", "recipes.foo_bar")},
		// paths withs IDs at the path second level
		{"GET", "/user/qwerty", NewMetricOperation("get", "user", MetricIDPlaceholder)},
		{"GET", "/users/qwerty", NewMetricOperation("get", "users", MetricIDPlaceholder)},
		{"GET", "/allergens/foobarbaz", NewMetricOperation("get", "allergens", MetricIDPlaceholder)},
		{"GET", "/cuisines/foobarbaz", NewMetricOperation("get", "cuisines", MetricIDPlaceholder)},
		{"GET", "/favorites/foobarbaz", NewMetricOperation("get", "favorites", MetricIDPlaceholder)},
		{"GET", "/ingredients/foobarbaz", NewMetricOperation("get", "ingredients", MetricIDPlaceholder)},
		{"GET", "/menus/foobarbaz", NewMetricOperation("get", "menus", MetricIDPlaceholder)},
		{"GET", "/ratings/foobarbaz", NewMetricOperation("get", "ratings", MetricIDPlaceholder)},
		{"GET", "/recipes/foobarbaz", NewMetricOperation("get", "recipes", MetricIDPlaceholder)},
		{"GET", "/addresses/foobarbaz", NewMetricOperation("get", "addresses", MetricIDPlaceholder)},
		{"GET", "/boxes/foobarbaz", NewMetricOperation("get", "boxes", MetricIDPlaceholder)},
		{"GET", "/coupons/foobarbaz", NewMetricOperation("get", "coupons", MetricIDPlaceholder)},
		{"GET", "/customers/foobarbaz", NewMetricOperation("get", "customers", MetricIDPlaceholder)},
		{"GET", "/delivery_options/foobarbaz", NewMetricOperation("get", "delivery_options", MetricIDPlaceholder)},
		{"GET", "/product_families/foobarbaz", NewMetricOperation("get", "product_families", MetricIDPlaceholder)},
		{"GET", "/products/foobarbaz", NewMetricOperation("get", "products", MetricIDPlaceholder)},
		{"GET", "/recipients/foobarbaz", NewMetricOperation("get", "recipients", MetricIDPlaceholder)},
		// path may have either numeric ID or non-numeric trackable path
		{"GET", "/subscriptions/12345", NewMetricOperation("get", "subscriptions", MetricIDPlaceholder)},
		{"GET", "/subscriptions/search", NewMetricOperation("get", "subscriptions", "search")},
		{"GET", "/freebies/12345", NewMetricOperation("get", "freebies", MetricIDPlaceholder)},
		{"GET", "/freebies/search", NewMetricOperation("get", "freebies", "search")},
		// path may be short or full
		{"GET", "/clients", NewMetricOperation("get", "clients", MetricEmptyPlaceholder)},
		{"GET", "/clients/qwe123", NewMetricOperation("get", "clients", MetricIDPlaceholder)},
	}

	idConfig := &SecondLevelIDConfig{
		HasIDAtSecondLevel: map[PathSection]SectionTestDefinition{
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
		},
		AutoDiscoverThreshold: 25,
		AutoDiscoverWhiteList: []string{"bar"},
	}
	callback := NewHasIDAtSecondLevelCallback(idConfig)

	for _, data := range dataProvider {
		r := &http.Request{Method: data.Method, URL: &url.URL{Path: data.Path}}
		assert.Equal(t, data.Operations, BuildHTTPRequestMetricOperation(r, callback))
	}

	var (
		logErrors      uint
		logMessages    uint
		logLastMessage string
		logLastFields  map[string]interface{}
		logLastError   error
	)
	log.SetHandler(func(msg string, fields map[string]interface{}, err error) {
		logLastMessage = msg
		logLastFields = fields
		logLastError = err
		if err != nil {
			logErrors++
		} else {
			logMessages++
		}
	})

	firstSectionFoo := "foo"
	firstSectionBar := "bar"

	uItoA := func(i uint) string {
		return strconv.FormatUint(uint64(i), 10)
	}

	for i := uint(0); i < idConfig.AutoDiscoverThreshold-1; i++ {
		rFoo := &http.Request{Method: http.MethodGet, URL: &url.URL{Path: fmt.Sprintf("/%s/%v", firstSectionFoo, i)}}
		assert.Equal(t, NewMetricOperation("get", firstSectionFoo, uItoA(i)), BuildHTTPRequestMetricOperation(rFoo, callback))

		rBar := &http.Request{Method: http.MethodGet, URL: &url.URL{Path: fmt.Sprintf("/%s/%v", firstSectionBar, i)}}
		assert.Equal(t, NewMetricOperation("get", firstSectionBar, uItoA(i)), BuildHTTPRequestMetricOperation(rBar, callback))
	}
	assert.Equal(t, uint(0), logErrors+logMessages)

	for i := idConfig.AutoDiscoverThreshold; i < idConfig.AutoDiscoverThreshold+idConfig.AutoDiscoverThreshold; i++ {
		rFoo := &http.Request{Method: http.MethodGet, URL: &url.URL{Path: fmt.Sprintf("/%s/%v", firstSectionFoo, i)}}
		assert.Equal(t, NewMetricOperation("get", firstSectionFoo, MetricIDPlaceholder), BuildHTTPRequestMetricOperation(rFoo, callback))

		rBar := &http.Request{Method: http.MethodGet, URL: &url.URL{Path: fmt.Sprintf("/%s/%v", firstSectionBar, i)}}
		assert.Equal(t, NewMetricOperation("get", firstSectionBar, uItoA(i)), BuildHTTPRequestMetricOperation(rBar, callback))
	}
	require.Equal(t, idConfig.AutoDiscoverThreshold, logErrors+logMessages)
	require.Error(t, logLastError)
	assert.Equal(t, logSuspiciousMetric, logLastMessage)
	assert.Equal(t, map[string]interface{}{
		"method":    "GET",
		"path":      "/foo/49",
		"operation": NewMetricOperation("get", "foo", "49"),
	}, logLastFields)
}

func TestHttpRequest_MetricNameAlterCallback(t *testing.T) {
	dataProvider := []struct {
		Method     string
		Path       string
		Operations MetricOperation
		Query      string
	}{
		{"GET", "/users/qwerty", NewMetricOperation("get", "users", MetricIDPlaceholder), ""},
		{"GET", "/clients", NewMetricOperation("get", "clients", MetricEmptyPlaceholder), ""},
		{"GET", "/clients/qwe123", NewMetricOperation("get", "clients", MetricIDPlaceholder), ""},
		{"GET", "/token/revoke", NewMetricOperation("get", "token", "revoke"), ""},
		{"GET", "/token/revoke", NewMetricOperation("get", "token", "revoke"), "foo=bar&grant_type=baz"},
		{"GET", "/token", NewMetricOperation("get", "token", "baz"), "foo=bar&grant_type=baz"},
		{"GET", "/token", NewMetricOperation("get", "token", MetricEmptyPlaceholder), "foo=bar"},
		{"GET", "/token/client_credentials", NewMetricOperation("get", "token", "client_credentials"), ""},
	}

	callback := func(metricFragments MetricOperation, r *http.Request) MetricOperation {
		if metricFragments.operations[1] == "token" && metricFragments.operations[2] != "revoke" {
			grantType := r.URL.Query().Get("grant_type")
			if grantType != "" {
				metricFragments.operations[2] = grantType
			}
			return metricFragments
		}

		return NewHasIDAtSecondLevelCallback(&SecondLevelIDConfig{
			HasIDAtSecondLevel: map[PathSection]SectionTestDefinition{
				"users":   {SectionTestIsNotEmpty, GetSectionTestCallback(SectionTestIsNotEmpty)},
				"clients": {SectionTestIsNotEmpty, GetSectionTestCallback(SectionTestIsNotEmpty)},
			}})(metricFragments, r)
	}

	for _, data := range dataProvider {
		r := &http.Request{Method: data.Method, URL: &url.URL{Path: data.Path, RawQuery: data.Query}}
		assert.Equal(t, data.Operations, BuildHTTPRequestMetricOperation(r, callback))
	}
}

func TestHttpRequest_Metric(t *testing.T) {
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
		b := NewHTTPRequest(SectionRequest, r, data.Success, nil, true)
		assert.Equal(t, data.Metric, b.Metric())
	}
}

func TestHttpRequest_MetricWithSuffix(t *testing.T) {
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
		b := NewHTTPRequest(SectionRequest, r, data.Success, nil, true)
		assert.Equal(t, data.Metric, b.MetricWithSuffix())
	}
}

func BenchmarkNewHTTPRequest(b *testing.B) {
	r := &http.Request{Method: "GET", URL: &url.URL{Path: "/foo/bar/baz"}}
	for n := 0; n < b.N; n++ {
		NewHTTPRequest(SectionRequest, r, true, nil, false)
	}
}

func BenchmarkNewHTTPRequest_unicode(b *testing.B) {
	r := &http.Request{Method: "GET", URL: &url.URL{Path: "/foo/bar/baz"}}
	for n := 0; n < b.N; n++ {
		NewHTTPRequest(SectionRequest, r, true, nil, true)
	}
}

func TestHttpRequest_MetricTotal(t *testing.T) {
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
		b := NewHTTPRequest(SectionRequest, r, data.Success, nil, true)
		assert.Equal(t, data.Metric, b.MetricTotal())
	}
}

func TestHttpRequest_TMetricTotalWithSuffix(t *testing.T) {
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
		b := NewHTTPRequest(SectionRequest, r, data.Success, nil, true)
		assert.Equal(t, data.Metric, b.MetricTotalWithSuffix())
	}
}

func TestHttpRequest_Metric_customSection(t *testing.T) {
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
		b := NewHTTPRequest(section, r, data.Success, nil, true)
		assert.Equal(t, data.Metric, b.Metric())
	}
}

func TestHttpRequest_MetricWithSuffix_customSection(t *testing.T) {
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
		b := NewHTTPRequest(section, r, data.Success, nil, true)
		assert.Equal(t, data.Metric, b.MetricWithSuffix())
	}
}
