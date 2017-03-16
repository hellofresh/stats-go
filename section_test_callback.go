package stats

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

const (
	sectionsDelimiter = ":"

	SectionTestTrue       = "true"
	SectionTestIsNumeric  = "numeric"
	SectionTestIsNotEmpty = "not_empty"
)

var (
	ErrInvalidFormat      = errors.New("Invalid sections format")
	ErrUnknownSectionTest = errors.New("Unknown section test")
)

type PathSection string
type SectionTestCallback func(string) bool
type SectionTestDefinition struct {
	Name     string
	Callback SectionTestCallback
}
type SectionsTestsMap map[PathSection]SectionTestDefinition

func (m SectionsTestsMap) String() string {
	var sections []string
	for k, v := range m {
		sections = append(sections, fmt.Sprintf("%s: %s", k, v.Name))
	}

	return fmt.Sprintf("[%s]", strings.Join(sections, ", "))
}

func TestAlwaysTrue(string) bool {
	return true
}

func TestIsNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func TestIsNotEmpty(s string) bool {
	return s != MetricEmptyPlaceholder
}

var (
	sectionsTestSync     sync.Mutex
	sectionsTestRegistry = map[string]SectionTestCallback{
		SectionTestTrue:       TestAlwaysTrue,
		SectionTestIsNumeric:  TestIsNumeric,
		SectionTestIsNotEmpty: TestIsNotEmpty,
	}
)

func NewHasIDAtSecondLevelCallback(hasIDAtSecondLevel SectionsTestsMap) HttpMetricNameAlterCallback {
	return func(operation MetricOperation, r *http.Request) MetricOperation {
		firstFragment := "/"
		for _, fragment := range strings.Split(r.URL.Path, "/") {
			if fragment != "" {
				firstFragment = fragment
				break
			}
		}
		if testFunction, ok := hasIDAtSecondLevel[PathSection(firstFragment)]; ok {
			if testFunction.Callback(operation[2]) {
				operation[2] = MetricIDPlaceholder
			}
		}

		return operation
	}
}

func RegisterSectionTest(name string, callback SectionTestCallback) {
	sectionsTestSync.Lock()
	defer sectionsTestSync.Unlock()

	sectionsTestRegistry[name] = callback
}

func GetSectionTestCallback(name string) SectionTestCallback {
	sectionsTestSync.Lock()
	defer sectionsTestSync.Unlock()

	return sectionsTestRegistry[name]
}

func ParseSectionsTestsMap(s string) (SectionsTestsMap, error) {
	result := make(SectionsTestsMap)
	var parts []string

	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) != "" {
			for _, part := range strings.Split(strings.TrimSpace(line), sectionsDelimiter) {
				if strings.TrimSpace(part) != "" {
					parts = append(parts, part)
				}
			}
		}
	}
	if len(parts)%2 != 0 {
		return nil, ErrInvalidFormat
	}

	for i := 0; i < len(parts); i += 2 {
		pathSection := PathSection(parts[i])
		sectionTestName := parts[i+1]

		if sectionTestCallback := GetSectionTestCallback(sectionTestName); sectionTestCallback == nil {
			return nil, ErrUnknownSectionTest
		} else {
			result[pathSection] = SectionTestDefinition{sectionTestName, sectionTestCallback}
		}
	}

	return result, nil
}
