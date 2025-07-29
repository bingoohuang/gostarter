package controllers

import (
	"github.com/bingoohuang/gostarter/pkg/db"
	"github.com/bingoohuang/gostarter/pkg/model"
)

func init() {
	registerController("POST", "/demo", demo)
	registerController("POST", "/users", addUser)
	registerController("GET", "/users", findUsers)
}

func demo(m *model.DemoReq) *model.DemoRsp {
	return &model.DemoRsp{Name: "Echo: " + m.Name}
}

type Result struct {
	Status  int
	Message string
}

type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Age   int    `db:"age"`
}

/*
	curl -X POST http://localhost:1234/users \
	  -H "Content-Type: application/json" \
	  -d '{"name":"zhangsan", "email":"zhangsan@gmail.com", "age":18}'
*/
func addUser(user User) string {
	sql := `INSERT INTO users (name, email, age) VALUES (?, ?, ?)`
	_, err := db.X.Exec(sql, user.Name, user.Email, user.Age)
	if err != nil {
		return err.Error()
	}
	return "success"
}

/*
	curl http://localhost:1234/users
*/

func findUsers() []User {
	var users []User
	err := db.X.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil
	}
	return users
}

/*
	$ gurl POST :1235/demo name=@name
	POST /demo HTTP/1.1
	Host: 127.0.0.1:1235
	Accept: application/json
	Accept-Encoding: gzip, deflate
	Content-Type: application/json
	Gurl-Date: Tue, 10 May 2022 06:20:41 GMT
	User-Agent: gurl/1.0.0

	{
	  "name": "Fairyink"
	}

	HTTP/1.1 200 OK
	Content-Type: application/json; charset=utf-8
	Date: Tue, 10 May 2022 06:20:41 GMT
	Content-Length: 25

	{
	  "name": "Echo: Fairyink"
	}
*/
