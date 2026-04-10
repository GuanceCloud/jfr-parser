package v0

import (
	"regexp"
	"strings"
)

var (
	PathSplitRegex = regexp.MustCompile(`\\/`)

	HumanCategoryNames = map[string]string{
		"adf":        "ADF",
		"java":       "Java Application",
		"os":         "Operating System",
		"recordings": "Flight Recorder",
		"jbo":        "JBO",
		"vm":         "Java Virtual Machine",
		"dms":        "Dynamic Monitoring System",
		"gc":         "GC",
		"prof":       "Profiling",
		"class":      "Class Loading",
		"mds":        "Metadata Services",
		"wls":        "WebLogic Server",
	}
)

func GetHumanSegmentArray(path string) []string {
	pieces := PathSplitRegex.Split(path, -1)
	categories := make([]string, 0, len(pieces))
	for _, piece := range pieces {
		categories = append(categories, GetHumanSegmentName(strings.TrimSpace(piece)))
	}
	return categories
}

func GetHumanSegmentName(name string) string {
	if category, ok := HumanCategoryNames[name]; ok {
		return category
	}
	return humanifyName(name)
}

func humanifyName(name string) string {
	if name == "" {
		return "<Empty>"
	}

	sb := new(strings.Builder)
	firstLetter := true

	for _, u := range name {
		if u == '_' {
			u = ' '
		}
		if firstLetter {
			sb.WriteString(strings.ToUpper(string(u)))
		} else {
			sb.WriteRune(u)
		}
		firstLetter = u == ' '
	}
	return sb.String()
}
