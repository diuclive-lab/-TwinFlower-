package suites

import "twinflower/flowers/explore/evals"

// FilesystemSuite returns the eval suite for filesystem_skill.
func FilesystemSuite() evals.Suite {
	return evals.Suite{
		Name: "filesystem_skill",
		Cases: []evals.Case{
			{
				ID:     1,
				Name:   "list-directory",
				Prompt: "列出当前目录",
				Expectations: []string{
					"contains: stem/",
					"contains: vascular/",
				},
			},
			{
				ID:     2,
				Name:   "read-file",
				Prompt: "读 main.go",
				Expectations: []string{
					"contains: package main",
				},
			},
			{
				ID:     3,
				Name:   "find-largest",
				Prompt: "看看 Downloads 里最大的文件",
				Expectations: []string{
					"no error",
				},
			},
			{
				ID:     4,
				Name:   "missing-path-recovery",
				Prompt: "读 main_copy.go",
				Expectations: []string{
					"contains: 不存在",
				},
			},
		},
	}
}
