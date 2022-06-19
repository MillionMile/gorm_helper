package gorm_helper

import (
	"testing"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id        int       `gorm:"colum:id"`
	NickName  string    `gorm:"column:nick_name"`
	Phone     string    `gorm:"column:phone"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (User) TableName() string {
	return "users"
}

func TestFilterWhere(t *testing.T) {
	db, err := getMockGormDB()
	if err != nil {
		t.Errorf("GetGormDB failed %v", err)
	}

	createdAt, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")

	type args struct {
		id        int
		nickName  string
		phone     string
		createdAt time.Time
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no filter",
			args: args{
				id:       0,
				nickName: "",
				phone:    "",
			},
			want: "SELECT * FROM `users`",
		},
		{
			name: "filter equal int",
			args: args{
				id:       1,
				nickName: "",
				phone:    "",
			},
			want: "SELECT * FROM `users` WHERE id = (1)",
		},
		{
			name: "filter equal string",
			args: args{
				id:       0,
				nickName: "",
				phone:    "13333333333",
			},
			want: "SELECT * FROM `users` WHERE phone = ('13333333333')",
		},
		{
			name: "filter like string",
			args: args{
				id:       0,
				nickName: "%mile%",
				phone:    "",
			},
			want: "SELECT * FROM `users` WHERE nick_name like ('%mile%')",
		},
		{
			name: "filter equal time",
			args: args{
				id:        0,
				nickName:  "",
				phone:     "",
				createdAt: createdAt,
			},
			want: "SELECT * FROM `users` WHERE created_at = ('2006-01-02 15:04:05')",
		},
		{
			name: "comprehensive filter",
			args: args{
				id:       1,
				nickName: "million%",
				phone:    "13333333333",
			},
			want: "SELECT * FROM `users` WHERE id = (1) AND nick_name like ('million%') AND phone = ('13333333333')",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL := getRawSql(db, func(tx *gorm.DB) *gorm.DB {
				return tx.
					Scopes(FilterWhere("id = ?", tt.args.id)).
					Scopes(FilterWhere("nick_name like ?", tt.args.nickName)).
					Scopes(FilterWhere("phone = ?", tt.args.phone)).
					Scopes(FilterWhere("created_at = ?", tt.args.createdAt)).
					Find(&User{})
			})
			if gotSQL != tt.want {
				t.Errorf("getRawSql() = %v, want %v", gotSQL, tt.want)
			}
		})
	}
}
