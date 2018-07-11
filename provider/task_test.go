package provider

import (
	"reflect"
	"testing"
)

func Test_getDiscoverMovieByPage(t *testing.T) {
	type args struct {
		page int
	}
	tests := []struct {
		name    string
		args    args
		want    discoverResult
		wantErr bool
	}{
		{
			"ok",
			args{1},
			discoverResult{
				Page: 1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDiscoverMovieByPage(tt.args.page)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDiscoverMovieByPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDiscoverMovieByPage() = %v, want %v", got, tt.want)
			}
		})
	}
}
