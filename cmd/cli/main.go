package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/yauritux/cartsvc/pkg/adapter/repository/inmem"
	cartSvc "github.com/yauritux/cartsvc/pkg/usecase/carts"
	productSvc "github.com/yauritux/cartsvc/pkg/usecase/products"
)

var prodUsecase *productSvc.ProductUsecase
var cartRepository *inmem.CartRepository
var cartUsecase *cartSvc.CartUsecase

func init() {
	prodRepository := inmem.NewProductRepository()
	cartRepository = inmem.NewCartRepository("yauritux")
	prodUsecase = productSvc.NewProductUsecase(prodRepository)
	cartUsecase = cartSvc.NewCartUsecase(cartRepository, prodRepository)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	showMenu()

	for scanner.Scan() {
		text := scanner.Text()
		switch text {
		case "1":
			if err := addItemToCart(scanner); err != nil {
				fmt.Printf("failed to add item... %v\n", err)
			}
			fmt.Println()
			showMenu()
		case "2":
			if err := showCartItems(); err != nil {
				fmt.Printf("failed to show cart items...%v\n", err)
			}
			fmt.Println()
			showMenu()
		case "3":
			os.Exit(0)
		default:
			fmt.Println()
			showMenu()
		}
	}
}

func showMenu() {
	fmt.Println("1. Add item to cart")
	fmt.Println("2. Show items in cart")
	fmt.Println("3. Exit")
	fmt.Print("Your choice [1 or 3]: ")
}

func addItemToCart(r *bufio.Scanner) error {
	fmt.Print("Enter product ID: ")
	r.Scan()
	prodID := r.Text()
	prodFound, err := prodUsecase.FindByProductID(prodID)
	if err != nil {
		return err
	}
	prod, ok := prodFound.(*productSvc.Product)
	if !ok {
		return errors.New("failed to get product from the repository...invalid type of product usecase model")
	}

	fmt.Print("Enter amount: ")
	r.Scan()
	amt := r.Text()
	qty, err := strconv.Atoi(amt)
	if err != nil {
		return errors.New("invalid amount")
	}

	if err := cartUsecase.AddToCart("yauritux", &cartSvc.CartItem{
		ID:    prod.ID,
		Name:  prod.Name,
		Qty:   qty,
		Price: prod.Price,
		Disc:  prod.Disc,
	}); err != nil {
		fmt.Println("error = ", err)
		return err
	}
	fmt.Printf("added %d item of %s to cart\n", qty, prod.Name)
	return nil
}

func showCartItems() error {
	userCart, err := cartUsecase.FetchUserCart("yauritux")
	if err != nil {
		return err
	}
	cart, ok := userCart.(*cartSvc.Cart)
	if !ok {
		return errors.New("failed to show cart items, invalid type of cart usecase model")
	}
	fmt.Println("\nYour cart items:")
	for _, v := range cart.Items {
		fmt.Printf("%d pcs of %s (%s)\n", v.Qty, v.ID, v.Name)
	}
	return nil
}
