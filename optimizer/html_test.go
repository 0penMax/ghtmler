package optimizer

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetAllClasses(t *testing.T) {
	type args struct {
		htmlCode string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{"<p class='hello'><span class='test test2'> world</span> </p>"},
			want:    []string{"hello", "test", "test2"},
			wantErr: false,
		}, {
			name:    "test2",
			args:    args{"<p class='hello'><span class='test test2'> world</span> </p> <p class='outside'>text</p>"},
			want:    []string{"hello", "test", "test2", "outside"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllClasses(tt.args.htmlCode)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllClasses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllClasses() got = [%s], want [%s]", strings.Join(got, `','`), strings.Join(tt.want, `','`))
			}
		})
	}
}
