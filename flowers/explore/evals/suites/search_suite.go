// Package suites provides pre-built eval suites for the Exploration Bench.
package suites

import "twinflower/flowers/explore/evals"

// SearchSuite returns the eval suite for search_skill capability tracking.
// These cases track the evolution of ambiguity resolution over time.
func SearchSuite() evals.Suite {
	return evals.Suite{
		Name: "search_skill",
		Cases: []evals.Case{
			{
				ID:     1,
				Name:   "direct-search-chinese",
				Prompt: "搜索北京天气",
				Expectations: []string{
					"contains: 🔍",
					"contains: 北京",
				},
			},
			{
				ID:     2,
				Name:   "direct-search-english",
				Prompt: "search for Go language",
				Expectations: []string{
					"contains: 🔍",
					"contains: Go",
				},
			},
			{
				ID:     3,
				Name:   "empty-query-clarify",
				Prompt: "查一下",
				Expectations: []string{
					"contains: 搜索",
					"contains: 什么",
				},
			},
			{
				ID:     4,
				Name:   "ambiguous-brand-cn",
				Prompt: "查一下苹果",
				Expectations: []string{
					"contains: 苹果公司",
					"contains: 还是",
				},
			},
			{
				ID:     5,
				Name:   "ambiguous-disambiguated",
				Prompt: "搜索苹果公司",
				Expectations: []string{
					"contains: 🔍",
					"contains: 苹果",
				},
			},
			{
				ID:     6,
				Name:   "known-ambiguous-howabout",
				Prompt: "帮我查一下小米怎么样",
				Expectations: []string{
					"contains: 小米公司",
					"contains: 还是",
				},
			},
			{
				ID:     7,
				Name:   "non-search-not-routed",
				Prompt: "北京天气",
				Expectations: []string{
					"no error",
				},
			},
		},
	}
}
