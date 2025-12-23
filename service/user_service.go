package service

type CreateUserReq struct {
	Name string `json:"name" binding:"required"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetUserById(id string) (*User, error) {
	return &User{
		Id:   id,
		Name: "Tom",
	}, nil
}

func CreateUser(req CreateUserReq) *User {
	return &User{
		Id:   "1",
		Name: req.Name,
	}
}
