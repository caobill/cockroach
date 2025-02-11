// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package spanconfigkvaccessor

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/spanconfig"
	"github.com/cockroachdb/cockroach/pkg/testutils"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/stretchr/testify/require"
)

// TestValidateUpdateArgs ensures we validate arguments to
// UpdateSpanConfigRecords correctly.
func TestValidateUpdateArgs(t *testing.T) {
	defer leaktest.AfterTest(t)()

	for _, tc := range []struct {
		toDelete []spanconfig.Target
		toUpsert []spanconfig.Record
		expErr   string
	}{
		{
			toUpsert: nil, toDelete: nil,
			expErr: "",
		},
		{
			toDelete: []spanconfig.Target{
				spanconfig.MakeSpanTarget(
					roachpb.Span{Key: roachpb.Key("a")}, // empty end key in delete list
				),
			},
			expErr: "invalid span: a",
		},
		{
			toUpsert: []spanconfig.Record{
				{
					Target: spanconfig.MakeSpanTarget(
						roachpb.Span{Key: roachpb.Key("a")}, // empty end key in update list
					),
				},
			},
			expErr: "invalid span: a",
		},
		{
			toUpsert: []spanconfig.Record{
				{
					Target: spanconfig.MakeSpanTarget(
						roachpb.Span{Key: roachpb.Key("b"), EndKey: roachpb.Key("a")}, // invalid span; end < start
					),
				},
			},
			expErr: "invalid span: {b-a}",
		},
		{
			toDelete: []spanconfig.Target{
				spanconfig.MakeSpanTarget(
					roachpb.Span{Key: roachpb.Key("b"), EndKey: roachpb.Key("a")}, // invalid span; end < start
				),
			},
			expErr: "invalid span: {b-a}",
		},
		{
			toDelete: []spanconfig.Target{
				// overlapping spans in the same list.
				spanconfig.MakeSpanTarget(roachpb.Span{Key: roachpb.Key("a"), EndKey: roachpb.Key("c")}),
				spanconfig.MakeSpanTarget(roachpb.Span{Key: roachpb.Key("b"), EndKey: roachpb.Key("c")}),
			},
			expErr: "overlapping spans {a-c} and {b-c} in same list",
		},
		{
			toUpsert: []spanconfig.Record{ // overlapping spans in the same list
				{
					Target: spanconfig.MakeSpanTarget(
						roachpb.Span{Key: roachpb.Key("a"), EndKey: roachpb.Key("c")},
					),
				},
				{
					Target: spanconfig.MakeSpanTarget(
						roachpb.Span{Key: roachpb.Key("b"), EndKey: roachpb.Key("c")},
					),
				},
			},
			expErr: "overlapping spans {a-c} and {b-c} in same list",
		},
		{
			// Overlapping spans in different lists.
			toDelete: []spanconfig.Target{
				// Overlapping spans in the same list.
				spanconfig.MakeSpanTarget(roachpb.Span{Key: roachpb.Key("a"), EndKey: roachpb.Key("c")}),
			},
			toUpsert: []spanconfig.Record{
				{
					Target: spanconfig.MakeSpanTarget(
						roachpb.Span{Key: roachpb.Key("a"), EndKey: roachpb.Key("b")},
					),
				},
				{
					Target: spanconfig.MakeSpanTarget(
						roachpb.Span{Key: roachpb.Key("b"), EndKey: roachpb.Key("c")},
					),
				},
			},
			expErr: "",
		},
	} {
		require.True(t, testutils.IsError(validateUpdateArgs(tc.toDelete, tc.toUpsert), tc.expErr))
	}
}
