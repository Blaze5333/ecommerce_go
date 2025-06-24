package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	First_Name      *string            `json:"first_name" bson:"first_name" validate:"required,min=2,max=30" `
	Last_Name       *string            `json:"last_name" bson:"last_name"  validate:"required,min=2,max=30"`
	Email           *string            `json:"email" bson:"email" validate:"required,email"`
	Phone           *string            `json:"phone" bson:"phone" validate:"required,min=10,max=15"`
	Password        *string            `json:"password" bson:"password" validate:"required,min=8,max=30"`
	Token           *string            `json:"token" bson:"token"`
	Refresh_Token   *string            `json:"refresh_token" bson:"refresh_token"`
	Created_At      time.Time          `json:"created_at" bson:"created_at"`
	Updated_At      time.Time          `json:"updated_at" bson:"updated_at"`
	User_ID         string             `json:"user_id" bson:"user_id"`
	UserCart        []ProductUser      `json:"user_cart" bson:"user_cart"`
	Address_Details []Address          `json:"address" bson:"address" validate:"max=3"`
	Order_Status    []Order            `json:"order_status" bson:"order_status"`
}
type LoginUser struct {
	Email    *string `json:"email" bson:"email" validate:"required,email"`
	Password *string `json:"password" bson:"password" validate:"required,min=8,max=30"`
}
type Product struct {
	Product_ID   primitive.ObjectID `json:"_id" bson:"_id"`
	Product_Name *string            `json:"product_name" bson:"product_name"`
	Price        uint64             `json:"price" bson:"price"`
	Rating       *uint8             `json:"rating" bson:"rating"`
	Image        *string            `json:"image" bson:"image"`
}
type ProductUser struct {
	Product_ID   primitive.ObjectID `bson:"_id" json:"_id"`
	Product_Name *string            `bson:"product_name" json:"product_name"`
	Price        uint64             `bson:"price" json:"price"`
	Rating       *float32           `bson:"rating" json:"rating"`
	Image        *string            `bson:"image" json:"image"`
}
type Address struct {
	Address_ID     primitive.ObjectID `bson:"_id" json:"_id" validate:"required" binding:"required"`
	Address_Line_1 *string            `bson:"address_line_1" json:"address_line_1" validate:"required" binding:"required"`
	Address_Line_2 *string            `bson:"address_line_2" json:"address_line_2"`
	Phone          *string            `bson:"phone" json:"phone" validate:"required" binding:"required"`
	House          *string            `bson:"house" json:"house" validate:"required" binding:"required"`
	City           *string            `bson:"city" json:"city" validate:"required" binding:"required"`
	Pincode        *string            `bson:"pincode" json:"pincode" validate:"required" binding:"required"`
}
type Order struct {
	Order_ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Order_Cart     []ProductUser      `bson:"order_list" json:"order_list"`
	Ordered_At     time.Time          `bson:"ordered_at" json:"ordered_at"`
	Price          *uint64            `bson:"price" json:"price"`
	Discount       *int               `bson:"discount" json:"discount"`
	Payment_Method Payment            `bson:"payment_method" json:"payment_method"`
	UserId         primitive.ObjectID `bson:"user_id" json:"user_id"`
}
type Payment struct {
	Digital bool `bson:"digital" json:"digital"`
	COD     bool `bson:"cod" json:"cod"`
}
