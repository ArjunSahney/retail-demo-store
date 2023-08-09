// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

package main

import (
	"go-component-service/handlers"
	"net/http"
)

// Route Struct
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes Array
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		handlers.Index,
	},
	Route{
		"OrdersIndex",
		"GET",
		"/orders/all",
		handlers.OrderIndex,
	},
	Route{
		"OrderShowByID",
		"GET",
		"/orders/id/{orderID}",
		handlers.OrderShowByID,
	},
	Route{
		"OrderShowByUsername",
		"GET",
		"/orders/username/{username}",
		handlers.OrderIndexByUsername,
	},
	Route{
		"OrderCreate",
		"POST",
		"/orders",
		handlers.OrderCreate,
	},
	Route{
		"OrderCreate",
		"OPTIONS",
		"/orders",
		handlers.OrderCreate,
	},
	Route{
		"OrderUpdate",
		"PUT",
		"/orders/id/{orderID}",
		handlers.OrderUpdate,
	},
	Route{
		"OrderUpdate",
		"OPTIONS",
		"/orders/id/{orderID}",
		handlers.OrderUpdate,
	},
	Route{
		"ProductIndex",
		"GET",
		"/products/all",
		handlers.ProductIndex,
	},
	Route{
		"ProductShow",
		"GET",
		"/products/id/{productIDs}",
		handlers.ProductShow,
	},
	Route{
		"ProductFeatured",
		"GET",
		"/products/featured",
		handlers.ProductFeatured,
	},
	Route{
		"ProductInCategory",
		"GET",
		"/products/category/{categoryName}",
		handlers.ProductInCategory,
	},
	Route{
		"ProductUpdate",
		"PUT",
		"/products/id/{productID}",
		handlers.UpdateProduct,
	},
	Route{
		"ProductDelete",
		"DELETE",
		"/products/id/{productID}",
		handlers.DeleteProduct,
	},
	Route{
		"NewProduct",
		"POST",
		"/products",
		handlers.NewProduct,
	},
	Route{
		"InventoryUpdate",
		"PUT",
		"/products/id/{productID}/inventory",
		handlers.UpdateInventory,
	},
	Route{
		"CategoryIndex",
		"GET",
		"/categories/all",
		handlers.CategoryIndex,
	},
	Route{
		"CategoryShow",
		"GET",
		"/categories/id/{categoryID}",
		handlers.CategoryShow,
	},
	Route{
		"CartsIndex",
		"GET",
		"/carts",
		handlers.CartIndex,
	},
	Route{
		"CartShowByID",
		"GET",
		"/carts/{cartID}",
		handlers.CartShowByID,
	},
	Route{
		"CartCreate",
		"POST",
		"/carts",
		handlers.CartCreate,
	},
	Route{
		"CartCreate",
		"OPTIONS",
		"/carts",
		handlers.CartCreate,
	},
	Route{
		"CartUpdate",
		"PUT",
		"/carts/{cartID}",
		handlers.CartUpdate,
	},
	Route{
		"CartUpdate",
		"OPTIONS",
		"/carts/{cartID}",
		handlers.CartUpdate,
	},
	Route{
		"SignPayload",
		"POST",
		"/sign",
		handlers.SignAmazonPayPayload,
	},
	Route{
		"SignPayload",
		"OPTIONS",
		"/sign",
		handlers.SignAmazonPayPayload,
	},
}