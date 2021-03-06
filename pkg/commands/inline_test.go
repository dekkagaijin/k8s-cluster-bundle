// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"strings"
	"testing"

	"context"
	test "github.com/GoogleCloudPlatform/k8s-cluster-bundle/pkg/commands/testing"
	"github.com/GoogleCloudPlatform/k8s-cluster-bundle/pkg/transformer"

	bpb "github.com/GoogleCloudPlatform/k8s-cluster-bundle/pkg/apis/bundle/v1alpha1"
)

// Fake implementation of FileReader to create a fake Inliner for unit tests.
type fakeFileReader struct{}

func (f *fakeFileReader) ReadFile(ctx context.Context, file *bpb.File) ([]byte, error) {
	return nil, nil
}

func TestRunInline(t *testing.T) {
	validFile := "/bundle.yaml"
	invalidFile := "/invalid.yaml"

	var testcases = []struct {
		testName          string
		opts              *inlineOptions
		expectErrContains string
	}{
		{
			testName: "success case",
			opts: &inlineOptions{
				bundle: validFile,
				output: validFile,
			},
		},
		{
			testName: "bundle read error",
			opts: &inlineOptions{
				bundle: invalidFile,
				output: validFile,
			},
			expectErrContains: "error reading",
		},
		{
			testName: "bundle write error",
			opts: &inlineOptions{
				bundle: validFile,
				output: invalidFile,
			},
			expectErrContains: "error writing",
		},
	}

	// Override the createInlinerFn to return a fake Inliner.
	createInlinerFn = func(cwd string) *transformer.Inliner {
		return &transformer.Inliner{
			LocalReader: &fakeFileReader{},
		}
	}
	brw := test.NewFakeReaderWriter(validFile)
	ctx := context.Background()

	for _, tc := range testcases {
		t.Run(tc.testName, func(t *testing.T) {
			err := runInline(ctx, tc.opts, brw)
			if (tc.expectErrContains != "" && err == nil) || (tc.expectErrContains == "" && err != nil) {
				t.Errorf("runInline(opts: %+v) returned err: %v, Want Err: %v", tc.opts, err, tc.expectErrContains)
			}
			if err == nil {
				return
			}
			if !strings.Contains(err.Error(), tc.expectErrContains) {
				t.Errorf("runInline(opts: %+v) returned unexpected error message: %v, Should contain: %v", tc.opts, err, tc.expectErrContains)
			}
		})
	}
}
