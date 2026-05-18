package suites

import "twinflower/flowers/explore/evals"

// ClarifySuite tracks clarify rate and threshold evolution.
// These cases should see fewer clarifies over time as preference learning improves.
func ClarifySuite() evals.Suite {
	return evals.Suite{
		Name: "clarify_behavior",
		Cases: []evals.Case{
			{
				ID:     1,
				Name:   "unambiguous-direct",
				Prompt: "北京天气",
				Expectations: []string{
					"contains: 北京",
					"contains: °C",
					"not contains: 你是想",
				},
			},
			{
				ID:     2,
				Name:   "ambiguous-should-clarify",
				Prompt: "查一下",
				Expectations: []string{
					"contains: 搜索",
				},
			},
			{
				ID:     3,
				Name:   "ambiguous-topic",
				Prompt: "看看苹果怎么样",
				Expectations: []string{
					"no error",
				},
			},
		},
	}
}
