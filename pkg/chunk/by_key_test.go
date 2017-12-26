package chunk

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func c(id string) Descriptor {
	return Descriptor{UserID: id}
}

func TestUnique(t *testing.T) {
	for _, tc := range []struct {
		in   DescByKey
		want DescByKey
	}{
		{nil, DescByKey{}},
		{DescByKey{c("a"), c("a")}, DescByKey{c("a")}},
		{DescByKey{c("a"), c("a"), c("b"), c("b"), c("c")}, DescByKey{c("a"), c("b"), c("c")}},
		{DescByKey{c("a"), c("b"), c("c")}, DescByKey{c("a"), c("b"), c("c")}},
	} {
		have := unique(tc.in)
		if !reflect.DeepEqual(tc.want, have) {
			assert.Equal(t, tc.want, have)
		}
	}
}

func TestMerge(t *testing.T) {
	type args struct {
		a DescByKey
		b DescByKey
	}
	for _, tc := range []struct {
		args args
		want DescByKey
	}{
		{args{DescByKey{}, DescByKey{}}, DescByKey{}},
		{args{DescByKey{c("a")}, DescByKey{}}, DescByKey{c("a")}},
		{args{DescByKey{}, DescByKey{c("b")}}, DescByKey{c("b")}},
		{args{DescByKey{c("a")}, DescByKey{c("b")}}, DescByKey{c("a"), c("b")}},
		{
			args{DescByKey{c("a"), c("c")}, DescByKey{c("a"), c("b"), c("d")}},
			DescByKey{c("a"), c("b"), c("c"), c("d")},
		},
	} {
		have := merge(tc.args.a, tc.args.b)
		if !reflect.DeepEqual(tc.want, have) {
			assert.Equal(t, tc.want, have)
		}
	}
}

func TestNWayUnion(t *testing.T) {
	for _, tc := range []struct {
		in   []DescByKey
		want DescByKey
	}{
		{nil, DescByKey{}},
		{[]DescByKey{{c("a")}}, DescByKey{c("a")}},
		{[]DescByKey{{c("a")}, {c("a")}}, DescByKey{c("a")}},
		{[]DescByKey{{c("a")}, {}}, DescByKey{c("a")}},
		{[]DescByKey{{}, {c("b")}}, DescByKey{c("b")}},
		{[]DescByKey{{c("a")}, {c("b")}}, DescByKey{c("a"), c("b")}},
		{
			[]DescByKey{{c("a"), c("c"), c("e")}, {c("c"), c("d")}, {c("b")}},
			DescByKey{c("a"), c("b"), c("c"), c("d"), c("e")},
		},
		{
			[]DescByKey{{c("c"), c("d")}, {c("b")}, {c("a"), c("c"), c("e")}},
			DescByKey{c("a"), c("b"), c("c"), c("d"), c("e")},
		},
	} {
		have := nWayUnion(tc.in)
		if !reflect.DeepEqual(tc.want, have) {
			assert.Equal(t, tc.want, have)
		}
	}
}

func TestNWayIntersect(t *testing.T) {
	for _, tc := range []struct {
		in   []DescByKey
		want DescByKey
	}{
		{nil, DescByKey{}},
		{[]DescByKey{{c("a"), c("b"), c("c")}}, []Descriptor{c("a"), c("b"), c("c")}},
		{[]DescByKey{{c("a"), c("b"), c("c")}, {c("a"), c("c")}}, DescByKey{c("a"), c("c")}},
		{[]DescByKey{{c("a"), c("b"), c("c")}, {c("a"), c("c")}, {c("b")}}, DescByKey{}},
		{[]DescByKey{{c("a"), c("b"), c("c")}, {c("a"), c("c")}, {c("a")}}, DescByKey{c("a")}},
	} {
		have := nWayIntersect(tc.in)
		if !reflect.DeepEqual(tc.want, have) {
			assert.Equal(t, tc.want, have)
		}
	}
}
