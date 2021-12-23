package csv

import (
	"reflect"
	"testing"
)

type Response struct {
	Code   int
	Msg    string
	Result interface{}
}

type Lists struct {
	List []*Content
}

type Content struct {
	Name string `csv:"name"`
	Age  int    `csv:"age,string"`
}

func Test_parse_Parse(t *testing.T) {
	type fields struct {
		tagName string
	}
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		{
			name: "OK",
			fields: fields{
				tagName: "csv",
			},
			args: args{
				obj: &Lists{List: []*Content{
					{
						Name: "lihua",
						Age:  26,
					},
					{
						Name: "zhangsan",
						Age:  20,
					},
				}},
			},
			want: []map[string]interface{}{
				{
					"List": []*Content{
						{
							Name: "lihua",
							Age:  26,
						},
						{
							Name: "zhangsan",
							Age:  20,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &parse{
				tagName: tt.fields.tagName,
			}
			got, err := p.Parse(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
