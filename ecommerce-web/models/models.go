package models

import (
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Id         *int    `json:"id"`
	Email      *string `json:"email"`
	First_Name *string `json:"firstName"`
	Last_Name  *string `json:"lastName"`
	Password   *string `json:"Password"`
}
type Image struct {
	Url         *string `json:"imageURL"`
	Description *string `json:"description"`
}
type Category struct {
	Id_category   *int    `json:"idCategory"`
	Name_category *string `json:"nameCategory"`
	Desc_category *string `json:"descCategory"`
}
type Technical struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
type Product struct {
	ProductId       *int        `json:"id"`
	ProductName     *string     `json:"name"`
	Rate            *float64    `json:"rate"`
	Sold            *int        `json:"sold"`
	CurrentPrice    *float64    `json:"currentPrice"`
	LastPrice       *float64    `json:"lastPrice"`
	ProductImages   []Image     `json:"productImage"`
	ProuctTechnical []Technical `json:"productTechnical"`
	ProductGroup    []Group     `json:"productGroup"`
	ProductComment  []Comment   `json:"productComment"`
}
type CartItem struct {
	Id           *int     `json:"id"`
	Name         *string  `json:"name"`
	CurrentPrice *float64 `json:"currentPrice"`
	LastPrice    *float64 `json:"lastPrice"`
	ProductImage []Image  `json:"productImage"`
	Quantity     *int     `json:"quantity"`
}
type Group struct {
	Id    *int    `json:"id"`
	Title *string `json:"title"`
	Type  []Type  `json:"productType"`
}
type Type struct {
	Id          *int    `json:"id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
type Product_Type struct {
	IdType    *int `json:"idType"`
	IdProduct *int `json:"idProduct"`
}
type Comment struct {
	FirstName   *string `json:"firstName"`
	LastName    *string `json:"lastName"`
	Description *string `json:"description"`
}
type WishList struct {
	Id           *int     `json:"id"`
	Name         *string  `json:"name"`
	CurrentPrice *float64 `json:"currentPrice"`
	LastPrice    *float64 `json:"lastPrice"`
	ProductImage []Image  `json:"productImage"`
}

