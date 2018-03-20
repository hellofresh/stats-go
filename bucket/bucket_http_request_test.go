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
		{"GET", "/", NewMetricOperation([3]string{"get", MetricEmptyPlaceholder, MetricEmptyPlaceholder}, []string{})},
		{"TRACE", "/api", NewMetricOperation([3]string{"trace", "api", MetricEmptyPlaceholder}, []string{})},
		{"TRACE", "/api/", NewMetricOperation([3]string{"trace", "api", MetricEmptyPlaceholder}, []string{})},
		{"POST", "/api/recipes", NewMetricOperation([3]string{"post", "api", "recipes"}, []string{})},
		{"POST", "/api/recipes/", NewMetricOperation([3]string{"post", "api", "recipes"}, []string{})},
		{"DELETE", "/api/recipes/123", NewMetricOperation([3]string{"delete", "api", "recipes"}, []string{})},
		{"DELETE", "/api/recipes.foo-bar/123", NewMetricOperation([3]string{"delete", "api", "recipes.foo-bar"}, []string{})},
		{"DELETE", "/api/recipes.foo_bar/123", NewMetricOperation([3]string{"delete", "api", "recipes.foo_bar"}, []string{})},
		// paths withs IDs at the path second level
		{"GET", "/user/qwerty", NewMetricOperation([3]string{"get", "user", MetricIDPlaceholder}, []string{})},
		{"GET", "/users/qwerty", NewMetricOperation([3]string{"get", "users", MetricIDPlaceholder}, []string{})},
		{"GET", "/allergens/foobarbaz", NewMetricOperation([3]string{"get", "allergens", MetricIDPlaceholder}, []string{})},
		{"GET", "/cuisines/foobarbaz", NewMetricOperation([3]string{"get", "cuisines", MetricIDPlaceholder}, []string{})},
		{"GET", "/favorites/foobarbaz", NewMetricOperation([3]string{"get", "favorites", MetricIDPlaceholder}, []string{})},
		{"GET", "/ingredients/foobarbaz", NewMetricOperation([3]string{"get", "ingredients", MetricIDPlaceholder}, []string{})},
		{"GET", "/menus/foobarbaz", NewMetricOperation([3]string{"get", "menus", MetricIDPlaceholder}, []string{})},
		{"GET", "/ratings/foobarbaz", NewMetricOperation([3]string{"get", "ratings", MetricIDPlaceholder}, []string{})},
		{"GET", "/recipes/foobarbaz", NewMetricOperation([3]string{"get", "recipes", MetricIDPlaceholder}, []string{})},
		{"GET", "/addresses/foobarbaz", NewMetricOperation([3]string{"get", "addresses", MetricIDPlaceholder}, []string{})},
		{"GET", "/boxes/foobarbaz", NewMetricOperation([3]string{"get", "boxes", MetricIDPlaceholder}, []string{})},
		{"GET", "/coupons/foobarbaz", NewMetricOperation([3]string{"get", "coupons", MetricIDPlaceholder}, []string{})},
		{"GET", "/customers/foobarbaz", NewMetricOperation([3]string{"get", "customers", MetricIDPlaceholder}, []string{})},
		{"GET", "/delivery_options/foobarbaz", NewMetricOperation([3]string{"get", "delivery_options", MetricIDPlaceholder}, []string{})},
		{"GET", "/product_families/foobarbaz", NewMetricOperation([3]string{"get", "product_families", MetricIDPlaceholder}, []string{})},
		{"GET", "/products/foobarbaz", NewMetricOperation([3]string{"get", "products", MetricIDPlaceholder}, []string{})},
		{"GET", "/recipients/foobarbaz", NewMetricOperation([3]string{"get", "recipients", MetricIDPlaceholder}, []string{})},
		// path may have either numeric ID or non-numeric trackable path
		{"GET", "/subscriptions/12345", NewMetricOperation([3]string{"get", "subscriptions", MetricIDPlaceholder}, []string{})},
		{"GET", "/subscriptions/search", NewMetricOperation([3]string{"get", "subscriptions", "search"}, []string{})},
		{"GET", "/freebies/12345", NewMetricOperation([3]string{"get", "freebies", MetricIDPlaceholder}, []string{})},
		{"GET", "/freebies/search", NewMetricOperation([3]string{"get", "freebies", "search"}, []string{})},
		// path may be short or full
		{"GET", "/clients", NewMetricOperation([3]string{"get", "clients", MetricEmptyPlaceholder}, []string{})},
		{"GET", "/clients/qwe123", NewMetricOperation([3]string{"get", "clients", MetricIDPlaceholder}, []string{})},
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
		assert.Equal(t, NewMetricOperation([3]string{"get", firstSectionFoo, uItoA(i)}, []string{}), BuildHTTPRequestMetricOperation(rFoo, callback))

		rBar := &http.Request{Method: http.MethodGet, URL: &url.URL{Path: fmt.Sprintf("/%s/%v", firstSectionBar, i)}}
		assert.Equal(t, NewMetricOperation([3]string{"get", firstSectionBar, uItoA(i)}, []string{}), BuildHTTPRequestMetricOperation(rBar, callback))
	}
	assert.Equal(t, uint(0), logErrors+logMessages)

	for i := idConfig.AutoDiscoverThreshold; i < idConfig.AutoDiscoverThreshold+idConfig.AutoDiscoverThreshold; i++ {
		rFoo := &http.Request{Method: http.MethodGet, URL: &url.URL{Path: fmt.Sprintf("/%s/%v", firstSectionFoo, i)}}
		assert.Equal(t, NewMetricOperation([3]string{"get", firstSectionFoo, MetricIDPlaceholder}, []string{}), BuildHTTPRequestMetricOperation(rFoo, callback))

		rBar := &http.Request{Method: http.MethodGet, URL: &url.URL{Path: fmt.Sprintf("/%s/%v", firstSectionBar, i)}}
		assert.Equal(t, NewMetricOperation([3]string{"get", firstSectionBar, uItoA(i)}, []string{}), BuildHTTPRequestMetricOperation(rBar, callback))
	}
	require.Equal(t, idConfig.AutoDiscoverThreshold, logErrors+logMessages)
	require.Error(t, logLastError)
	assert.Equal(t, logSuspiciousMetric, logLastMessage)
	assert.Equal(t, map[string]interface{}{
		"method":    "GET",
		"path":      "/foo/49",
		"operation": NewMetricOperation([3]string{"get", "foo", "49"}, []string{}),
	}, logLastFields)
}

func TestHttpRequest_MetricNameAlterCallback(t *testing.T) {
	dataProvider := []struct {
		Method     string
		Path       string
		Operations MetricOperation
		Query      string
	}{
		{"GET", "/users/qwerty", NewMetricOperation([3]string{"get", "users", MetricIDPlaceholder}, []string{}), ""},
		{"GET", "/clients", NewMetricOperation([3]string{"get", "clients", MetricEmptyPlaceholder}, []string{}), ""},
		{"GET", "/clients/qwe123", NewMetricOperation([3]string{"get", "clients", MetricIDPlaceholder}, []string{}), ""},
		{"GET", "/token/revoke", NewMetricOperation([3]string{"get", "token", "revoke"}, []string{}), ""},
		{"GET", "/token/revoke", NewMetricOperation([3]string{"get", "token", "revoke"}, []string{}), "foo=bar&grant_type=baz"},
		{"GET", "/token", NewMetricOperation([3]string{"get", "token", "baz"}, []string{}), "foo=bar&grant_type=baz"},
		{"GET", "/token", NewMetricOperation([3]string{"get", "token", MetricEmptyPlaceholder}, []string{}), "foo=bar"},
		{"GET", "/token/client_credentials", NewMetricOperation([3]string{"get", "token", "client_credentials"}, []string{}), ""},
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
