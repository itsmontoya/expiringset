package expiringset

import (
	"testing"
	"time"
)

func TestExpiringSet_Set(t *testing.T) {
	type args struct {
		key string
	}

	type testcase struct {
		name       string
		prepare    func() *ExpiringSet
		args       args
		wantExists bool
	}

	tests := []testcase{
		{
			name: "no exists",
			prepare: func() *ExpiringSet {
				return New(time.Minute)
			},
			args: args{
				key: "foo",
			},
			wantExists: false,
		},
		{
			name: "exists",
			prepare: func() *ExpiringSet {
				e := New(time.Minute)
				e.Set("foo")
				return e
			},
			args: args{
				key: "foo",
			},
			wantExists: true,
		},
		{
			name: "expired, no exists",
			prepare: func() *ExpiringSet {
				e := New(time.Millisecond)
				e.Set("foo")
				time.Sleep(time.Millisecond * 2)
				return e
			},
			args: args{
				key: "foo",
			},
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.prepare()
			if gotExists := e.Set(tt.args.key); gotExists != tt.wantExists {
				t.Errorf("ExpiringSet.Set() = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}

func TestExpiringSet_set(t *testing.T) {
	type args struct {
		key string
	}

	type testcase struct {
		name       string
		prepare    func() *ExpiringSet
		args       args
		wantExists bool
	}

	tests := []testcase{
		{
			name: "no exists",
			prepare: func() *ExpiringSet {
				return New(time.Minute)
			},
			args: args{
				key: "foo",
			},
			wantExists: false,
		},
		{
			name: "exists",
			prepare: func() *ExpiringSet {
				e := New(time.Minute)
				e.Set("foo")
				return e
			},
			args: args{
				key: "foo",
			},
			wantExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.prepare()
			if gotExists := e.set(tt.args.key); gotExists != tt.wantExists {
				t.Errorf("ExpiringSet.set() = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}
