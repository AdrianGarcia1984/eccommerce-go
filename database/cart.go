package database

import "errors"

var (
	ErrCantFindProduct    = errors.New("can't find the products")
	ErrCantDecodeProducts = errors.New("can't find the products")
	ErrUserIdNotValid     = errors.New("this user not valid")
	ErrCantUpdateUser     = errors.New("cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove product to the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
	ErrCantGetItem        = errors.New("was unable to get the item from the cart")
)

func AddProductToCart()

func RemoveCartItem()

func BuyItemFromCart()

func InstantBuyer()
