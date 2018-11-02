/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	duckv1alpha1 "github.com/knative/pkg/apis/duck/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func TestGitHubSourceStatusIsReady(t *testing.T) {
	tests := []struct {
		name string
		s    *GitHubSourceStatus
		want bool
	}{{
		name: "uninitialized",
		s:    &GitHubSourceStatus{},
		want: false,
	}, {
		name: "initialized",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			return s
		}(),
		want: false,
	}, {
		name: "mark sink",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			return s
		}(),
		want: false,
	}, {
		name: "mark valid",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			s.MarkValid()
			return s
		}(),
		want: false,
	}, {
		name: "mark sink and valid",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkValid()
			return s
		}(),
		want: true,
	}, {
		name: "mark sink and valid then no sink",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkValid()
			s.MarkNoSink("Testing", "")
			return s
		}(),
		want: false,
	}, {
		name: "mark sink and valid then not valid",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkValid()
			s.MarkNotValid("Testing", "")
			return s
		}(),
		want: false,
	}, {
		name: "mark sink empty and valid",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("")
			s.MarkValid()
			return s
		}(),
		want: false,
	}, {
		name: "mark sink empty and valid then sink",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("")
			s.MarkValid()
			s.MarkSink("uri://example")
			return s
		}(),
		want: true,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.s.IsReady()
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("%s: unexpected condition (-want, +got) = %v", test.name, diff)
			}
		})
	}
}

func TestGitHubSourceStatusGetCondition(t *testing.T) {
	tests := []struct {
		name      string
		s         *GitHubSourceStatus
		condQuery duckv1alpha1.ConditionType
		want      *duckv1alpha1.Condition
	}{{
		name:      "uninitialized",
		s:         &GitHubSourceStatus{},
		condQuery: GitHubSourceConditionReady,
		want:      nil,
	}, {
		name: "initialized",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			return s
		}(),
		condQuery: GitHubSourceConditionReady,
		want: &duckv1alpha1.Condition{
			Type:   GitHubSourceConditionReady,
			Status: corev1.ConditionUnknown,
		},
	}, {
		name: "mark sink",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			return s
		}(),
		condQuery: GitHubSourceConditionReady,
		want: &duckv1alpha1.Condition{
			Type:   GitHubSourceConditionReady,
			Status: corev1.ConditionUnknown,
		},
	}, {
		name: "mark valid",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			s.MarkValid()
			return s
		}(),
		condQuery: GitHubSourceConditionReady,
		want: &duckv1alpha1.Condition{
			Type:   GitHubSourceConditionReady,
			Status: corev1.ConditionUnknown,
		},
	}, {
		name: "mark sink and valid",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkValid()
			return s
		}(),
		condQuery: GitHubSourceConditionReady,
		want: &duckv1alpha1.Condition{
			Type:   GitHubSourceConditionReady,
			Status: corev1.ConditionTrue,
		},
	}, {
		name: "mark sink and valid then no sink",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkValid()
			s.MarkNoSink("Testing", "hi%s", "")
			return s
		}(),
		condQuery: GitHubSourceConditionReady,
		want: &duckv1alpha1.Condition{
			Type:    GitHubSourceConditionReady,
			Status:  corev1.ConditionFalse,
			Reason:  "Testing",
			Message: "hi",
		},
	}, {
		name: "mark sink and valid then not valid",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("uri://example")
			s.MarkValid()
			s.MarkNotValid("Testing", "hi%s", "")
			return s
		}(),
		condQuery: GitHubSourceConditionReady,
		want: &duckv1alpha1.Condition{
			Type:    GitHubSourceConditionReady,
			Status:  corev1.ConditionFalse,
			Reason:  "Testing",
			Message: "hi",
		},
	}, {
		name: "mark sink empty and valid",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("")
			s.MarkValid()
			return s
		}(),
		condQuery: GitHubSourceConditionReady,
		want: &duckv1alpha1.Condition{
			Type:    GitHubSourceConditionReady,
			Status:  corev1.ConditionUnknown,
			Reason:  "SinkEmpty",
			Message: "Sink has resolved to empty.",
		},
	}, {
		name: "mark sink empty and valid then sink",
		s: func() *GitHubSourceStatus {
			s := &GitHubSourceStatus{}
			s.InitializeConditions()
			s.MarkSink("")
			s.MarkValid()
			s.MarkSink("uri://example")
			return s
		}(),
		condQuery: GitHubSourceConditionReady,
		want: &duckv1alpha1.Condition{
			Type:   GitHubSourceConditionReady,
			Status: corev1.ConditionTrue,
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.s.GetCondition(test.condQuery)
			ignoreTime := cmpopts.IgnoreFields(duckv1alpha1.Condition{}, "LastTransitionTime")
			if diff := cmp.Diff(test.want, got, ignoreTime); diff != "" {
				t.Errorf("unexpected condition (-want, +got) = %v", diff)
			}
		})
	}
}
