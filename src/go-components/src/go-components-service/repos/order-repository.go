// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

package repos

import (
	"go-component-service/models"
	"strconv"
)

var currentID int

var Orders models.Orders = models.Orders{}

// Init
func init() {
}

// RepoFindOrderByID Function
func RepoFindOrderByID(id string) models.Order {
	for _, t := range Orders {
		if t.ID == id {
			return t
		}
	}
	// return empty Order if not found
	return models.Order{}
}

// RepoFindOrdersByUsername Function
func RepoFindOrdersByUsername(username string) models.Orders {

	var o models.Orders = models.Orders{}

	for _, t := range Orders {
		if t.Username == username {
			o = append(o, t)
		}
	}

	return o
}

func RepoUpdateOrder(t models.Order) models.Order {

	for i := 0; i < len(Orders); i++ {
		o := &Orders[i]
		if o.ID == t.ID {
			Orders[i] = t
			return RepoFindOrderByID(t.ID)
		}
	}

	// return empty Order if not found
	return models.Order{}
}

// RepoCreateOrder Function
func RepoCreateOrder(t models.Order) models.Order {
	currentID++
	t.ID = strconv.Itoa(currentID)
	Orders = append(Orders, t)
	return t
}
