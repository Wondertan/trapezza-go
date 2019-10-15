package types

type Restaurant struct {
	ID string

	Address string
	Waiters []Waiter
}

type User struct {
	ID string

	Name  string
	Phone string
	Email string
}

type Client struct {
	User
}

type Waiter struct {
	User
}

type Order struct {
	Client Client
	Items  []MenuItem
}

type Trapeza struct {
	ID string

	RestaurantID string
	Waiter       Waiter
	Orders       []Order
}

type MenuItem struct {
	ID string

	Name string
}
