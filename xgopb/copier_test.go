package xgopb_test

import (
	"testing"
	"time"

	"github.com/glassonion1/xgo/xgopb"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestDeepCopy_protobuf(t *testing.T) {

	// Model type
	type SampleType int64
	// Protobuf type
	type Sample_SampleType int32
	const (
		SampleTypeAllowStandby SampleType        = 1
		Sample_ALLOW_STANDBY   Sample_SampleType = 1
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

	type PbPos struct {
		Lat float64
		Lon float64
	}

	// Protobuf struct
	type PbSample struct {
		Id         string
		Order      int64
		CreatedAt  *timestamppb.Timestamp
		UpdatedAt  *timestamppb.Timestamp
		SampleType Sample_SampleType
		Pos        *PbPos
	}

	type TimestampField struct {
		FinishedAt *timestamppb.Timestamp
	}

	type TimeField struct {
		FinishedAt time.Time
	}

	type DurationPbField struct {
		Duration *durationpb.Duration
	}

	type DurationField struct {
		Duration time.Duration
	}

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
				CreatedAt:  &timestamppb.Timestamp{Seconds: now.Unix(), Nanos: 0},
				UpdatedAt:  &timestamppb.Timestamp{Seconds: now.Unix(), Nanos: 0},
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
				CreatedAt:  &timestamppb.Timestamp{Seconds: now.Unix(), Nanos: 0},
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
					CreatedAt:  &timestamppb.Timestamp{Seconds: now.Unix(), Nanos: 0},
					UpdatedAt:  &timestamppb.Timestamp{Seconds: now.Unix(), Nanos: 0},
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
					Duration: &durationpb.Duration{Seconds: 300},
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
				Duration: &durationpb.Duration{Seconds: 300},
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
			err := xgopb.DeepCopy(tt.in.src, tt.in.dest)
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
			opt := cmpopts.IgnoreUnexported(timestamppb.Timestamp{},
				durationpb.Duration{})
			if diff := cmp.Diff(tt.want, got, opt); diff != "" {
				t.Errorf("testing %s mismatch (-want +got):\n%s\n", tt.name, diff)
			}
		})
	}
}
