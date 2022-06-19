package gorm_helper

import (
	"testing"

	"gorm.io/gorm"
)

func TestPaginate(t *testing.T) {
	db, err := getMockGormDB()
	if err != nil {
		t.Errorf("GetGormDB failed %v", err)
	}

	type args struct {
		page     int64
		pageSize int64
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default page 1 when page zero",
			args: args{
				page:     0,
				pageSize: 10,
			},
			want: "SELECT * FROM `users` LIMIT 10",
		},
		{
			name: "default pageSize 10 when pageSize zero",
			args: args{
				page:     10,
				pageSize: 0,
			},
			want: "SELECT * FROM `users` LIMIT 10 OFFSET 90",
		},
		{
			name: "jump page 1",
			args: args{
				page:     1,
				pageSize: 10,
			},
			want: "SELECT * FROM `users` LIMIT 10",
		},
		{
			name: "jump page 3",
			args: args{
				page:     3,
				pageSize: 10,
			},
			want: "SELECT * FROM `users` LIMIT 10 OFFSET 20",
		},
		{
			name: "default page 1 when error page",
			args: args{
				page:     -10,
				pageSize: 10,
			},
			want: "SELECT * FROM `users` LIMIT 10",
		},
		{
			name: "default pageSize 1 when error pageSize",
			args: args{
				page:     1,
				pageSize: -10,
			},
			want: "SELECT * FROM `users` LIMIT 10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL := getRawSql(db, func(tx *gorm.DB) *gorm.DB {
				return tx.
					Scopes(Paginate(tt.args.page, tt.args.pageSize)).
					Find(&User{})
			})
			if gotSQL != tt.want {
				t.Errorf("getRawSql() = %v, want %v", gotSQL, tt.want)
			}
		})
	}
}
