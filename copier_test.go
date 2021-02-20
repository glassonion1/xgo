package xgo_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/glassonion1/xgo"
)

type ModelSample struct {
	ID         string `copier:"Id"`
	Order      int64
	CreatedAt  time.Time
	UpdatedAt  *time.Time
	SampleType SampleType
	Pos        struct {
		Lat float64
		Lon float64
	}
}

type SampleType int64

const (
	SampleTypeAllowStandby SampleType = iota + 1
)

// Protobuf struct
type PbSample struct {
	Id         string
	Order      int64
	CreatedAt  *timestamp.Timestamp
	UpdatedAt  *timestamp.Timestamp
	SampleType Sample_SampleType
	Pos        *PbPos
}

type PbPos struct {
	Lat float64
	Lon float64
}

// Protobuf custom type
type Sample_SampleType int32

const (
	Sample_ALLOW_STANDBY Sample_SampleType = iota + 1
)

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

type ModelC struct {
	Field string
}

type ModelD struct {
	Field string
}

type Int32Field struct {
	ID int32
}

type Int64Field struct {
	ID int64
}

type TimestampField struct {
	FinishedAt *timestamp.Timestamp
}

type TimeField struct {
	FinishedAt time.Time
}

type DurationPbField struct {
	Duration *duration.Duration
}

type DurationField struct {
	Duration time.Duration
}

type StringField struct {
	Name string
}

// MyString is custom string type
type MyString string

type MyStringField struct {
	Name MyString
}

type PrivateField struct {
	ID    string
	name  string
	state int
}

func TestDeepCopy(t *testing.T) {
	type args struct {
		src  interface{}
		dest interface{}
	}

	now := time.Unix(time.Now().Unix(), 0)

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

func TestDeepCopy_CustomType(t *testing.T) {
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

func TestDeepCopy_Protobuf(t *testing.T) {
	type args struct {
		src  interface{}
		dest interface{}
	}

	now := time.Unix(time.Now().Unix(), 0)

	tests := []struct {
		name string
		in   args
		want interface{}
		err  error
	}{
		{
			name: "model to pb",
			in: args{
				src: ModelSample{
					ID:         "ididididid",
					Order:      1,
					CreatedAt:  now,
					UpdatedAt:  &now,
					SampleType: SampleTypeAllowStandby,
					Pos: struct {
						Lat float64
						Lon float64
					}{
						Lat: 1234.56,
						Lon: 6543.21,
					},
				},
				dest: &PbSample{},
			},
			want: &PbSample{
				Id:         "ididididid",
				Order:      1,
				CreatedAt:  &timestamp.Timestamp{Seconds: now.Unix(), Nanos: 0},
				UpdatedAt:  &timestamp.Timestamp{Seconds: now.Unix(), Nanos: 0},
				SampleType: Sample_ALLOW_STANDBY,
				Pos: &PbPos{
					Lat: 1234.56,
					Lon: 6543.21,
				},
			},
			err: nil,
		},
		{
			name: "model to pb",
			in: args{
				src: ModelSample{
					ID:         "ididididid",
					Order:      1,
					CreatedAt:  now,
					UpdatedAt:  nil,
					SampleType: SampleTypeAllowStandby,
					Pos: struct {
						Lat float64
						Lon float64
					}{
						Lat: 1234.56,
						Lon: 6543.21,
					},
				},
				dest: &PbSample{},
			},
			want: &PbSample{
				Id:         "ididididid",
				Order:      1,
				CreatedAt:  &timestamp.Timestamp{Seconds: now.Unix(), Nanos: 0},
				UpdatedAt:  nil,
				SampleType: Sample_ALLOW_STANDBY,
				Pos: &PbPos{
					Lat: 1234.56,
					Lon: 6543.21,
				},
			},
			err: nil,
		},
		{
			name: "pb to model",
			in: args{
				src: PbSample{
					Id:         "ididididid",
					Order:      1,
					CreatedAt:  &timestamp.Timestamp{Seconds: now.Unix(), Nanos: 0},
					UpdatedAt:  &timestamp.Timestamp{Seconds: now.Unix(), Nanos: 0},
					SampleType: Sample_ALLOW_STANDBY,
					Pos: &PbPos{
						Lat: 1234.56,
						Lon: 6543.21,
					},
				},
				dest: &ModelSample{},
			},
			want: &ModelSample{
				ID:         "ididididid",
				Order:      1,
				CreatedAt:  now,
				UpdatedAt:  &now,
				SampleType: SampleTypeAllowStandby,
				Pos: struct {
					Lat float64
					Lon float64
				}{
					Lat: 1234.56,
					Lon: 6543.21,
				},
			},
			err: nil,
		},
		{
			name: "timestamp to time zero value",
			in: args{
				src:  TimestampField{},
				dest: &TimeField{},
			},
			want: &TimeField{},
			err:  nil,
		},
		{
			name: "durationPbField to duration",
			in: args{
				src: &DurationPbField{
					Duration: &duration.Duration{Seconds: 300},
				},
				dest: &DurationField{},
			},
			want: &DurationField{Duration: 300 * time.Second},
			err:  nil,
		},
		{
			name: "duration to durationPbField",
			in: args{
				src:  &DurationField{Duration: 300 * time.Second},
				dest: &DurationPbField{},
			},
			want: &DurationPbField{
				Duration: &duration.Duration{Seconds: 300},
			},
			err: nil,
		},
		{
			name: "durationPb nil Field to duration",
			in: args{
				src:  &DurationPbField{},
				dest: &DurationField{},
			},
			want: &DurationField{},
			err:  nil,
		},
		{
			name: "duration zero value to durationPbField",
			in: args{
				src:  &DurationField{},
				dest: &DurationPbField{},
			},
			want: &DurationPbField{},
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
			opt := cmpopts.IgnoreUnexported(timestamp.Timestamp{},
				duration.Duration{})
			if diff := cmp.Diff(tt.want, got, opt); diff != "" {
				t.Errorf("testing %s mismatch (-want +got):\n%s\n", tt.name, diff)
			}
		})
	}
}

func TestDeepCopy_Private(t *testing.T) {
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
