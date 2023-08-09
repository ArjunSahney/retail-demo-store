// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

package repos

import (
	"go-component-service/models"
	"strconv"
)

var currentCartID int

var Carts = map[string]models.Cart{}

// RepoFindCartByID Function
func RepoFindCartByID(id string) models.Cart {
	cart, ok := Carts[id]
	if !ok {
		return models.Cart{}
	}
	return cart
}

// RepoUpdateCart Function
func RepoUpdateCart(id string, cart models.Cart) models.Cart {
	_, ok := Carts[id]

	if !ok {
		// return empty Cart if not found
		return models.Cart{}
	}

	cart.ID = id
	Carts[id] = cart

	return cart
}

// RepoCreateCart Function
func RepoCreateCart(t models.Cart) models.Cart {
	currentCartID++
	t.ID = strconv.Itoa(currentCartID)
	Carts[t.ID] = t
	return t
}
