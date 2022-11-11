package db

type User struct {
	Id   int
	Name string
}

type UserCreateDTO struct {
	Name string `json:"name" binding:"required"`
}

type Balance struct {
	Id     int
	UserId int
	Amount int64
}

type Hold struct {
	Id            int
	TransactionId int
	BalanceId     int
}

type TransactionType string

type Transaction struct {
	Id     int
	UserId int
	Amount int
	Type   string
}

type TransactionCreateDTO struct {
	UserId int             `json:"user_id" binding:"required"`
	Amount int             `json:"amount" binding:"required"`
	Type   TransactionType `json:"type" binding:"required"`
}
