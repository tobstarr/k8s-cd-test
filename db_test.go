package main

import "testing"

func TestDatabase(t *testing.T) {
	con, err := connect("postgres://localhost/k8s_cd_test?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	tx, err := con.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("CREATE TABLE users (id INTEGER NOT NULL PRIMARY KEY, email VARCHAR NOT NULL)")
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("INSERT INTO users (id, email) VALUES ($1, $2)", 1, "tobstarr@gmail.com")
	if err != nil {
		t.Fatal(err)
	}
	var cnt int
	err = tx.QueryRow("SELECT count(1) FROM users").Scan(&cnt)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct{ Has, Want interface{} }{
		{cnt, 1},
	}
	for i, tc := range tests {
		if tc.Has != tc.Want {
			t.Errorf("%d: want=%#v has=%#v", i+1, tc.Want, tc.Has)
		}
	}
}
