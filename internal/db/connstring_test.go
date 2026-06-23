package db

import "testing"

func TestEncodeUserinfo(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "hash in password (issue #64)",
			in:   "postgres://db_user:EwQxO#g#blawiuZUFGVukzw@host:5432/table_name",
			want: "postgres://db_user:EwQxO%23g%23blawiuZUFGVukzw@host:5432/table_name",
		},
		{
			name: "question mark in password",
			in:   "postgres://u:p?ss@host/db",
			want: "postgres://u:p%3Fss@host/db",
		},
		{
			name: "slash in password",
			in:   "postgres://u:p/ss@host/db",
			want: "postgres://u:p%2Fss@host/db",
		},
		{
			name: "percent-encoded sequence treated as literal",
			in:   "postgres://u:p%23ss@host/db",
			want: "postgres://u:p%2523ss@host/db",
		},
		{
			name: "mixed percent and reserved treated as literal",
			in:   "postgres://u:a%23b#c@host/db",
			want: "postgres://u:a%2523b%23c@host/db",
		},
		{
			name: "no password",
			in:   "postgres://user@host/db",
			want: "postgres://user@host/db",
		},
		{
			name: "no userinfo",
			in:   "postgres://host:5432/db",
			want: "postgres://host:5432/db",
		},
		{
			name: "keyword DSN untouched",
			in:   "host=localhost user=u password=p#ass",
			want: "host=localhost user=u password=p#ass",
		},
		{
			name: "sqlite path untouched",
			in:   "/home/u/data.db",
			want: "/home/u/data.db",
		},
		{
			name: "snowflake with query params and # in password",
			in:   "snowflake://u:p#ss@acct/db?warehouse=w&role=r",
			want: "snowflake://u:p%23ss@acct/db?warehouse=w&role=r",
		},
		{
			name: "clickhouse with # in password",
			in:   "clickhouse://u:p#ss@host:9000/db",
			want: "clickhouse://u:p%23ss@host:9000/db",
		},
		{
			name: "sqlserver URL with # in password",
			in:   "sqlserver://u:p#ss@host:1433?database=db",
			want: "sqlserver://u:p%23ss@host:1433?database=db",
		},
		{
			name: "oracle URL with # in password",
			in:   "oracle://u:p#ss@host:1521/sid",
			want: "oracle://u:p%23ss@host:1521/sid",
		},
		{
			name: "firebird URL with # in password",
			in:   "firebird://u:p#ss@host/db",
			want: "firebird://u:p%23ss@host/db",
		},
		{
			name: "@ in query is not mistaken for userinfo delimiter",
			in:   "postgres://u:p@host/db?role=a@b",
			want: "postgres://u:p@host/db?role=a@b",
		},
		{
			name: "@ in fragment is not mistaken for userinfo delimiter",
			in:   "postgres://u:p@host/db#sec@tion",
			want: "postgres://u:p@host/db#sec@tion",
		},
		{
			name: "raw percent in password encoded as %25",
			in:   "postgres://u:100%pass@host/db",
			want: "postgres://u:100%25pass@host/db",
		},
		{
			name: "truncated percent escape at end of password",
			in:   "postgres://u:abc%@host/db",
			want: "postgres://u:abc%25@host/db",
		},
		{
			name: "@ in password encoded as %40",
			in:   "postgres://u:p@ss@host/db",
			want: "postgres://u:p%40ss@host/db",
		},
		{
			name: "multiple @ in password all encoded",
			in:   "postgres://u:a@b@c@host/db",
			want: "postgres://u:a%40b%40c@host/db",
		},
		{
			name: "username-only with reserved char encoded",
			in:   "postgres://us#er@host/db",
			want: "postgres://us%23er@host/db",
		},
		{
			name: "username-only with no reserved char unchanged",
			in:   "postgres://user@host/db",
			want: "postgres://user@host/db",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := EncodeUserinfo(tc.in)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
