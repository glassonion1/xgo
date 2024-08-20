package xgo_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/glassonion1/xgo"
)

func TestDeepCopy(t *testing.T) {

	type ModelC struct {
		Field string
	}

	type ModelD struct {
		Field string
	}

	type ModelA struct {
		Field     string
		CreatedAt time.Time
		UpdatedAt *time.Time
		Model     ModelC
	}

	type ModelB struct {
		Field     string
		CreatedAt time.Time
		UpdatedAt *time.Time
		Model     ModelD
	}

	type Int32Field struct {
		ID int32
	}

	type Int64Field struct {
		ID int64
	}

	type FromPtrTestdata struct {
		StructToStruct ModelA
		StructToPtr    ModelA
		PtrToStruct    *ModelA
		PtrToPtr       *ModelA
	}

	type ToPtrTestdata struct {
		StructToStruct ModelB
		StructToPtr    *ModelB
		PtrToStruct    ModelB
		PtrToPtr       *ModelB
	}

	type args struct {
		src  interface{}
		dest interface{}
	}

	now := time.Now()

	tests := []struct {
		name string
		in   args
		want interface{}
		err  error
	}{
		{
			name: "struct copy",
			in: args{
				src: FromPtrTestdata{
					StructToStruct: ModelA{
						Field:     "foo",
						CreatedAt: now,
						UpdatedAt: &now,
						Model: ModelC{
							Field: "bar",
						},
					},
					StructToPtr: ModelA{
						Field: "foo",
						Model: ModelC{
							Field: "bar",
						},
					},
					PtrToStruct: &ModelA{
						Field: "foo",
						Model: ModelC{
							Field: "bar",
						},
					},
					PtrToPtr: &ModelA{
						Field: "foo",
						Model: ModelC{
							Field: "bar",
						},
					},
				},
				dest: &ToPtrTestdata{},
			},
			want: &ToPtrTestdata{
				StructToStruct: ModelB{
					Field:     "foo",
					CreatedAt: now,
					UpdatedAt: &now,
					Model: ModelD{
						Field: "bar",
					},
				},
				StructToPtr: &ModelB{
					Field: "foo",
					Model: ModelD{
						Field: "bar",
					},
				},
				PtrToStruct: ModelB{
					Field: "foo",
					Model: ModelD{
						Field: "bar",
					},
				},
				PtrToPtr: &ModelB{
					Field: "foo",
					Model: ModelD{
						Field: "bar",
					},
				},
			},
			err: nil,
		},
		{
			name: "struct copy:field not found",
			in: args{
				src: FromPtrTestdata{
					StructToStruct: ModelA{
						Field: "foo",
					},
				},
				dest: &ModelA{},
			},
			want: &ModelA{},
			err:  nil,
		},
		{
			name: "int64 to int32 OK",
			in: args{
				src: Int64Field{
					ID: 2100000000,
				},
				dest: &Int32Field{},
			},
			want: &Int32Field{
				ID: 2100000000,
			},
			err: nil,
		},
		{
			name: "int64 to int32 overflow",
			in: args{
				src: Int64Field{
					ID: 2200000000,
				},
				dest: &Int32Field{},
			},
			want: &Int32Field{
				ID: -2094967296,
			},
			err: nil,
		},
		{
			name: "int32 to int64 OK",
			in: args{
				src: Int32Field{
					ID: 2100000000,
				},
				dest: &Int64Field{},
			},
			want: &Int64Field{
				ID: 2100000000,
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := xgo.DeepCopy(tt.in.src, tt.in.dest)
			got := tt.in.dest
			if tt.err == nil && err != nil {
				t.Errorf("testing %s: should not be error for %#v but: %v", tt.name, tt.in, err)
			}
			if tt.err != nil && err == nil {
				t.Errorf("testing %s: should be error for %#v but not:", tt.name, tt.in)
			}
			if tt.err != nil && err != tt.err {
				t.Errorf("testing %s: should be error of %v but got: %v", tt.name, tt.err, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("testing %s mismatch (-want +got):\n%s\n", tt.name, diff)
			}
		})
	}
}

func TestDeepCopy_time(t *testing.T) {

	type ModelA struct {
		CreatedAt *time.Time
		UpdatedAt time.Time
	}

	type ModelB struct {
		CreatedAt time.Time
		UpdatedAt *time.Time
	}

	type args struct {
		src  interface{}
		dest interface{}
	}

	now := time.Now()

	tests := []struct {
		name string
		in   args
		want interface{}
		err  error
	}{
		{
			name: "struct copy",
			in: args{
				src: ModelA{
					CreatedAt: &now,
					UpdatedAt: now,
				},
				dest: &ModelB{},
			},
			want: &ModelB{
				CreatedAt: now,
				UpdatedAt: &now,
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := xgo.DeepCopy(tt.in.src, tt.in.dest)
			got := tt.in.dest
			if tt.err == nil && err != nil {
				t.Errorf("testing %s: should not be error for %#v but: %v", tt.name, tt.in, err)
			}
			if tt.err != nil && err == nil {
				t.Errorf("testing %s: should be error for %#v but not:", tt.name, tt.in)
			}
			if tt.err != nil && err != tt.err {
				t.Errorf("testing %s: should be error of %v but got: %v", tt.name, tt.err, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("testing %s mismatch (-want +got):\n%s\n", tt.name, diff)
			}
		})
	}
}

func TestDeepCopy_customType(t *testing.T) {
	type StringField struct {
		Name string
	}

	type MyString string

	type MyStringField struct {
		Name MyString
	}

	type args struct {
		src  interface{}
		dest interface{}
	}

	tests := []struct {
		name string
		in   args
		want interface{}
		err  error
	}{
		{
			name: "string to myString value",
			in: args{
				src:  StringField{Name: "foo"},
				dest: &MyStringField{},
			},
			want: &MyStringField{Name: "foo"},
			err:  nil,
		},
		{
			name: "myString to string value",
			in: args{
				src:  MyStringField{Name: "bar"},
				dest: &StringField{},
			},
			want: &StringField{Name: "bar"},
			err:  nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := xgo.DeepCopy(tt.in.src, tt.in.dest)
			got := tt.in.dest
			if tt.err == nil && err != nil {
				t.Errorf("testing %s: should not be error for %#v but: %v", tt.name, tt.in, err)
			}
			if tt.err != nil && err == nil {
				t.Errorf("testing %s: should be error for %#v but not:", tt.name, tt.in)
			}
			if tt.err != nil && err != tt.err {
				t.Errorf("testing %s: should be error of %v but got: %v", tt.name, tt.err, err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("testing %s mismatch (-want +got):\n%s\n", tt.name, diff)
			}
		})
	}
}

func TestDeepCopy_private(t *testing.T) {
	type PrivateField struct {
		ID    string
		name  string
		state int
	}

	type args struct {
		src  PrivateField
		dest *PrivateField
	}

	tests := []struct {
		name string
		in   args
		want *PrivateField
		err  error
	}{
		{
			name: "private field value",
			in: args{
				src:  PrivateField{ID: "id", name: "hoge", state: 100},
				dest: &PrivateField{},
			},
			want: &PrivateField{ID: "id", name: ""},
			err:  nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := xgo.DeepCopy(tt.in.src, tt.in.dest)
			got := tt.in.dest
			if tt.err == nil && err != nil {
				t.Errorf("testing %s: should not be error for %#v but: %v", tt.name, tt.in, err)
			}
			if tt.err != nil && err == nil {
				t.Errorf("testing %s: should be error for %#v but not:", tt.name, tt.in)
			}
			if tt.err != nil && err != tt.err {
				t.Errorf("testing %s: should be error of %v but got: %v", tt.name, tt.err, err)
			}
			if ok := reflect.DeepEqual(tt.want, got); !ok {
				t.Errorf("testing %s mismatch (-want +got):\n%v\n%v", tt.name, tt.want, got)
			}
		})
	}
}

func TestDeepCopy_slice(t *testing.T) {
	type Model1 struct {
		Foo string
		Bar int
	}
	type Model2 struct {
		Foo string
		Bar int
	}
	type Example1 struct {
		ID         string
		Name       string
		State      int
		Tests      []string
		StructPtrs []*Model1
		Structs    []Model1
	}
	type Example2 struct {
		ID         string
		Name       string
		State      int
		Tests      []string
		StructPtrs []*Model2
		Structs    []Model2
	}

	type args struct {
		src  Example1
		dest *Example2
	}

	tests := []struct {
		name string
		in   args
		want *Example2
		err  error
	}{
		{
			name: "slice value",
			in: args{
				src: Example1{
					ID:         "id1",
					Name:       "hoge1",
					State:      100,
					Tests:      []string{"test1", "test2"},
					StructPtrs: []*Model1{{Foo: "foo1", Bar: 100}, {Foo: "foo2", Bar: 200}},
					Structs:    []Model1{{Foo: "foo1", Bar: 1000}, {Foo: "foo2", Bar: 2000}},
				},
				dest: &Example2{},
			},
			want: &Example2{
				ID:         "id1",
				Name:       "hoge1",
				State:      100,
				Tests:      []string{"test1", "test2"},
				StructPtrs: []*Model2{{Foo: "foo1", Bar: 100}, {Foo: "foo2", Bar: 200}},
				Structs:    []Model2{{Foo: "foo1", Bar: 1000}, {Foo: "foo2", Bar: 2000}},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := xgo.DeepCopy(tt.in.src, tt.in.dest)
			got := tt.in.dest
			if tt.err == nil && err != nil {
				t.Errorf("testing %s: should not be error for %#v but: %v", tt.name, tt.in, err)
			}
			if tt.err != nil && err == nil {
				t.Errorf("testing %s: should be error for %#v but not:", tt.name, tt.in)
			}
			if tt.err != nil && err != tt.err {
				t.Errorf("testing %s: should be error of %v but got: %v", tt.name, tt.err, err)
			}
			if ok := reflect.DeepEqual(tt.want, got); !ok {
				t.Errorf("testing %s mismatch (-want +got):\n%v\n%v", tt.name, tt.want, got)
			}
		})
	}
}

func TestDeepCopy_slice2(t *testing.T) {
	type Model1 struct {
		Foo string
		Bar int
	}
	type Model2 struct {
		Foo string
		Bar int
	}
	type Example1 struct {
		ID         string
		Name       string
		State      int
		Tests      []string
		StructPtrs []*Model1
		Structs    []Model1
	}
	type Example2 struct {
		ID         string
		Name       string
		State      int
		Tests      []string
		StructPtrs []*Model2
		Structs    []Model2
	}

	type args struct {
		src  []*Example1
		dest []*Example2
	}

	tests := []struct {
		name string
		in   args
		want []*Example2
		err  error
	}{
		{
			name: "slice value",
			in: args{
				src: []*Example1{
					{
						ID:         "id1",
						Name:       "hoge1",
						State:      100,
						Tests:      []string{"test1", "test2"},
						StructPtrs: []*Model1{{Foo: "foo1-1", Bar: 100}, {Foo: "foo1-2", Bar: 200}},
						Structs:    []Model1{{Foo: "foo1-1", Bar: 1000}, {Foo: "foo1-2", Bar: 2000}},
					},
					{
						ID:      "id2",
						Name:    "hoge2",
						State:   200,
						Tests:   []string{"test1", "test2"},
						Structs: []Model1{{Foo: "foo2-1", Bar: 1000}, {Foo: "foo2-2", Bar: 2000}},
					},
					{
						ID:         "id3",
						Name:       "hoge3",
						State:      300,
						Tests:      []string{"test1", "test2"},
						StructPtrs: []*Model1{{Foo: "foo3-1", Bar: 100}, {Foo: "foo3-2", Bar: 200}},
					},
				},
				dest: []*Example2{},
			},
			want: []*Example2{
				{
					ID:         "id1",
					Name:       "hoge1",
					State:      100,
					Tests:      []string{"test1", "test2"},
					StructPtrs: []*Model2{{Foo: "foo1-1", Bar: 100}, {Foo: "foo1-2", Bar: 200}},
					Structs:    []Model2{{Foo: "foo1-1", Bar: 1000}, {Foo: "foo1-2", Bar: 2000}},
				},
				{
					ID:      "id2",
					Name:    "hoge2",
					State:   200,
					Tests:   []string{"test1", "test2"},
					Structs: []Model2{{Foo: "foo2-1", Bar: 1000}, {Foo: "foo2-2", Bar: 2000}},
				},
				{
					ID:         "id3",
					Name:       "hoge3",
					State:      300,
					Tests:      []string{"test1", "test2"},
					StructPtrs: []*Model2{{Foo: "foo3-1", Bar: 100}, {Foo: "foo3-2", Bar: 200}},
				},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := xgo.DeepCopy(tt.in.src, &tt.in.dest)
			got := tt.in.dest
			if tt.err == nil && err != nil {
				t.Errorf("testing %s: should not be error for %#v but: %v", tt.name, tt.in, err)
			}
			if tt.err != nil && err == nil {
				t.Errorf("testing %s: should be error for %#v but not:", tt.name, tt.in)
			}
			if tt.err != nil && err != tt.err {
				t.Errorf("testing %s: should be error of %v but got: %v", tt.name, tt.err, err)
			}
			if ok := reflect.DeepEqual(tt.want, got); !ok {
				t.Errorf("testing %s mismatch (-want +got):\n%v\n%v", tt.name, tt.want, got)
			}
		})
	}
}

func TestDeepCopy_ptr(t *testing.T) {
	type Foo string
	type Bar string
	type FooNum int
	type BarNum int

	type Example1 struct {
		Tests   []Foo
		Test2s  []*Foo
		Test    *Foo
		Test2   Foo
		Nums    []FooNum
		NumPtrs []*FooNum
	}
	type Example2 struct {
		Tests   []Bar
		Test2s  []Bar
		Test    Bar
		Test2   *Bar
		Nums    []BarNum
		NumPtrs []*BarNum
	}

	type args struct {
		src  Example1
		dest *Example2
	}

	tests := []struct {
		name string
		in   args
		want *Example2
		err  error
	}{
		{
			name: "slice value",
			in: args{
				src: Example1{
					Tests: []Foo{"test1", "test2"},
					//Test2s:  []*Foo{xgo.ToPtr(Foo("test1")), xgo.ToPtr(Foo("test2"))},
					Test:    xgo.ToPtr(Foo("test3")),
					Test2:   "test4",
					Nums:    []FooNum{1, 2, 3, 4},
					NumPtrs: []*FooNum{xgo.ToPtr(FooNum(1)), xgo.ToPtr(FooNum(2))},
				},
				dest: &Example2{},
			},
			want: &Example2{
				Tests: []Bar{"test1", "test2"},
				//Test2s:  []Bar{"test1", "test2"},
				Test:    "test3",
				Test2:   xgo.ToPtr(Bar("test4")),
				Nums:    []BarNum{1, 2, 3, 4},
				NumPtrs: []*BarNum{xgo.ToPtr(BarNum(1)), xgo.ToPtr(BarNum(2))},
			},
			err: nil,
		},
		{
			name: "nil or zero value",
			in: args{
				src: Example1{
					Tests: []Foo{"test1", "test2"},
					Nums:  []FooNum{1, 2, 3, 4},
				},
				dest: &Example2{},
			},
			want: &Example2{
				Tests: []Bar{"test1", "test2"},
				Test:  "",
				Test2: xgo.ToPtr(Bar("")),
				Nums:  []BarNum{1, 2, 3, 4},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := xgo.DeepCopy(tt.in.src, tt.in.dest)
			got := tt.in.dest
			if tt.err == nil && err != nil {
				t.Errorf("testing %s: should not be error for %#v but: %v", tt.name, tt.in, err)
			}
			if tt.err != nil && err == nil {
				t.Errorf("testing %s: should be error for %#v but not:", tt.name, tt.in)
			}
			if tt.err != nil && err != tt.err {
				t.Errorf("testing %s: should be error of %v but got: %v", tt.name, tt.err, err)
			}
			if ok := reflect.DeepEqual(tt.want, got); !ok {
				t.Errorf("testing %s mismatch (-want +got):\n%v\n%v", tt.name, tt.want, got)
			}
		})
	}
}
