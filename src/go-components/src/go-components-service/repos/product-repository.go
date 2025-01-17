// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

package repos

import (
	"fmt"
	"go-component-service/models"
	"go-component-service/util"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	guuuid "github.com/google/uuid"
)

// Root/base URL to use when building fully-qualified URLs to product detail view.
var webRootURL = os.Getenv("WEB_ROOT_URL")

var MAX_BATCH_GET_ITEM = 100

func setProductURL(p *models.Product) {
	if len(webRootURL) > 0 {
		p.URL = webRootURL + "/#/product/" + p.ID
	}
}

func setCategoryURL(c *models.Category) {
	if len(webRootURL) > 0 && len(c.Name) > 0 {
		c.URL = webRootURL + "/#/category/" + c.Name
	}
}

// RepoFindProduct Function
func RepoFindProduct(id string) models.Product {
	var product models.Product

	id = strings.ToLower(id)

	log.Println("RepoFindProduct: ", id, util.DbTableProducts)

	result, err := util.DynamoClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(util.DbTableProducts),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		log.Println("get item error " + string(err.Error()))
		return product
	}

	if result.Item != nil {
		err = dynamodbattribute.UnmarshalMap(result.Item, &product)

		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}

		setProductURL(&product)

		log.Println("RepoFindProduct returning: ", product.Name, product.Category)
	}

	return product
}

// RepoFindMultipleProducts Function
func RepoFindMultipleProducts(ids []string) models.Products {
	if len(ids) > MAX_BATCH_GET_ITEM {
		panic(fmt.Sprintf("Failed to unmarshal Record, %d", MAX_BATCH_GET_ITEM))
	}

	var products models.Products

	mapOfAttrKeys := []map[string]*dynamodb.AttributeValue{}

	for _, id := range ids {
		mapOfAttrKeys = append(mapOfAttrKeys, map[string]*dynamodb.AttributeValue{
			"id": &dynamodb.AttributeValue{
				S: aws.String(id),
			},
		})
		log.Println(string(id))
	}

	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			util.DbTableProducts: &dynamodb.KeysAndAttributes{
				Keys: mapOfAttrKeys,
			},
		},
	}

	result, err := util.DynamoClient.BatchGetItem(input)

	if err != nil {
		log.Println("BatchGetItem error " + string(err.Error()))

		return products
	}

	var itemCount = 0

	for _, table := range result.Responses {
		for _, item := range table {
			product := models.Product{}

			err = dynamodbattribute.UnmarshalMap(item, &product)

			if err != nil {
				log.Println("Got error unmarshalling:")
				log.Println(err.Error())
			} else {
				setProductURL(&product)
			}

			products = append(products, product)
			itemCount += 1
		}
	}

	if itemCount == 0 {
		products = make([]models.Product, 0)
	}

	return products
}

// RepoFindCategory Function
func RepoFindCategory(id string) models.Category {
	var category models.Category

	id = strings.ToLower(id)

	log.Println("RepoFindCategory: ", id, util.DbTableCategories)

	result, err := util.DynamoClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(util.DbTableCategories),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		log.Println("get item error " + string(err.Error()))
		return category
	}

	if result.Item != nil {
		err = dynamodbattribute.UnmarshalMap(result.Item, &category)

		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}

		setCategoryURL(&category)

		log.Println("RepoFindCategory returning: ", category.Name)
	}

	return category
}

// RepoFindCategoriesByName Function
func RepoFindCategoriesByName(name string) models.Categories {
	var categories models.Categories

	log.Println("RepoFindCategoriesByName: ", name, util.DbTableCategories)

	keycond := expression.Key("name").Equal(expression.Value(name))
	proj := expression.NamesList(expression.Name("id"),
		expression.Name("name"),
		expression.Name("image"))
	expr, err := expression.NewBuilder().WithKeyCondition(keycond).WithProjection(proj).Build()

	if err != nil {
		log.Println("Got error building expression:")
		log.Println(err.Error())
	}

	// Build the query input parameters
	params := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(util.DbTableCategories),
		IndexName:                 aws.String("name-index"),
	}
	// Make the DynamoDB Query API call
	result, err := util.DynamoClient.Query(params)

	if err != nil {
		log.Println("Got error QUERY expression:")
		log.Println(err.Error())
	}

	log.Println("RepoFindCategoriesByName / items found =  ", len(result.Items))

	for _, i := range result.Items {
		item := models.Category{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			log.Println("Got error unmarshalling:")
			log.Println(err.Error())
		} else {
			setCategoryURL(&item)
		}

		categories = append(categories, item)
	}

	if len(result.Items) == 0 {
		categories = make([]models.Category, 0)
	}

	return categories
}

// RepoFindProductByCategory Function
func RepoFindProductByCategory(category string) models.Products {

	log.Println("RepoFindProductByCategory: ", category)

	var f models.Products

	keycond := expression.Key("category").Equal(expression.Value(category))
	proj := expression.NamesList(expression.Name("id"),
		expression.Name("category"),
		expression.Name("name"),
		expression.Name("image"),
		expression.Name("style"),
		expression.Name("description"),
		expression.Name("price"),
		expression.Name("gender_affinity"),
		expression.Name("current_stock"),
		expression.Name("promoted"))
	expr, err := expression.NewBuilder().WithKeyCondition(keycond).WithProjection(proj).Build()

	if err != nil {
		log.Println("Got error building expression:")
		log.Println(err.Error())
	}

	// Build the query input parameters
	params := &dynamodb.QueryInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(util.DbTableProducts),
		IndexName:                 aws.String("category-index"),
	}
	// Make the DynamoDB Query API call
	result, err := util.DynamoClient.Query(params)

	if err != nil {
		log.Println("Got error QUERY expression:")
		log.Println(err.Error())
	}

	log.Println("RepoFindProductByCategory / items found =  ", len(result.Items))

	for _, i := range result.Items {
		item := models.Product{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			log.Println("Got error unmarshalling:")
			log.Println(err.Error())
		} else {
			setProductURL(&item)
		}

		f = append(f, item)
	}

	if len(result.Items) == 0 {
		f = make([]models.Product, 0)
	}

	return f
}

// RepoFindFeatured Function
func RepoFindFeatured() models.Products {

	log.Println("RepoFindFeatured | featured=true")

	var f models.Products

	filt := expression.Name("featured").Equal(expression.Value("true"))
	expr, err := expression.NewBuilder().WithFilter(filt).Build()

	if err != nil {
		log.Println("Got error building expression:")
		log.Println(err.Error())
	}

	// Build the query input
	// using index for performance (few items are featured)
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(util.DbTableProducts),
		IndexName:                 aws.String("featured-index"),
	}
	// Make the DynamoDB Query API call
	result, err := util.DynamoClient.Scan(params)

	if err != nil {
		log.Println("Got error scan expression:")
		log.Println(err.Error())
	}

	log.Println("RepoFindProductFeatured / items found =  ", len(result.Items))

	for _, i := range result.Items {
		item := models.Product{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			log.Println("Got error unmarshalling:")
			log.Println(err.Error())
		} else {
			setProductURL(&item)
		}

		f = append(f, item)
	}

	if len(result.Items) == 0 {
		f = make([]models.Product, 0)
	}

	return f
}

// RepoFindALLCategories - loads all categories
func RepoFindALLCategories() models.Categories {
	// TODO: implement some caching

	log.Println("RepoFindALLCategories: ")

	var f models.Categories

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		TableName: aws.String(util.DbTableCategories),
	}
	// Make the DynamoDB Query API call
	result, err := util.DynamoClient.Scan(params)

	if err != nil {
		log.Println("Got error scan expression:")
		log.Println(err.Error())
	}

	log.Println("RepoFindALLCategories / items found =  ", len(result.Items))

	for _, i := range result.Items {
		item := models.Category{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			log.Println("Got error unmarshalling:")
			log.Println(err.Error())
		} else {
			setCategoryURL(&item)
		}

		f = append(f, item)
	}

	if len(result.Items) == 0 {
		f = make([]models.Category, 0)
	}

	return f
}

// RepoFindALLProducts Function
func RepoFindALLProducts() models.Products {

	log.Println("RepoFindALLProducts")

	var f models.Products

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		TableName: aws.String(util.DbTableProducts),
	}
	// Make the DynamoDB Query API call
	result, err := util.DynamoClient.Scan(params)

	if err != nil {
		log.Println("Got error scan expression:")
		log.Println(err.Error())
	}

	log.Println("RepoFindALLProducts / items found =  ", len(result.Items))

	for _, i := range result.Items {
		item := models.Product{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			log.Println("Got error unmarshalling:")
			log.Println(err.Error())
		} else {
			setProductURL(&item)
		}

		f = append(f, item)
	}

	if len(result.Items) == 0 {
		f = make([]models.Product, 0)
	}

	return f
}

// RepoUpdateProduct - updates an existing product
func RepoUpdateProduct(existingProduct *models.Product, updatedProduct *models.Product) error {
	updatedProduct.ID = existingProduct.ID // Ensure we're not changing product ID.
	updatedProduct.URL = ""                // URL is generated so ignore if specified
	log.Printf("UpdateProduct from %#v to %#v", existingProduct, updatedProduct)

	av, err := dynamodbattribute.MarshalMap(updatedProduct)

	if err != nil {
		fmt.Println("Got error calling dynamodbattribute MarshalMap:")
		fmt.Println(err.Error())
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(util.DbTableProducts),
	}

	_, err = util.DynamoClient.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
	}

	setProductURL(updatedProduct)

	return err
}

// RepoUpdateInventoryDelta - updates a product's current inventory
func RepoUpdateInventoryDelta(product *models.Product, stockDelta int) error {

	log.Printf("RepoUpdateInventoryDelta for product %#v, delta: %v", product, stockDelta)

	if product.CurrentStock+stockDelta < 0 {
		// ensuring we don't get negative stocks, just down to zero stock
		// FUTURE: allow backorders via negative current stock?
		stockDelta = -product.CurrentStock
	}

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":stock_delta": {
				N: aws.String(strconv.Itoa(stockDelta)),
			},
			":currstock": {
				N: aws.String(strconv.Itoa(product.CurrentStock)),
			},
		},
		TableName: aws.String(util.DbTableProducts),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(product.ID),
			},
			"category": {
				S: aws.String(product.Category),
			},
		},
		ReturnValues:        aws.String("UPDATED_NEW"),
		UpdateExpression:    aws.String("set current_stock = current_stock + :stock_delta"),
		ConditionExpression: aws.String("current_stock = :currstock"),
	}

	_, util.Pro_err = util.DynamoClient.UpdateItem(input)
	if util.Pro_err != nil {
		fmt.Println("Got error calling UpdateItem:")
		fmt.Println(util.Pro_err.Error())
	} else {
		product.CurrentStock = product.CurrentStock + stockDelta
	}

	return util.Pro_err
}

// RepoNewProduct - initializes and persists new product
func RepoNewProduct(product *models.Product) error {
	log.Printf("RepoNewProduct --> %#v", product)

	if len(product.ID) == 0 {
		product.ID = strings.ToLower(guuuid.New().String())
	}
	av, err := dynamodbattribute.MarshalMap(product)

	if err != nil {
		fmt.Println("Got error calling dynamodbattribute MarshalMap:")
		fmt.Println(err.Error())
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(util.DbTableProducts),
	}

	_, err = util.DynamoClient.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
	}

	setProductURL(product)

	return err
}

// RepoDeleteProduct - deletes a single product
func RepoDeleteProduct(product *models.Product) error {
	log.Println("Deleting product: ", product)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(product.ID),
			},
			"category": {
				S: aws.String(product.Category),
			},
		},
		TableName: aws.String(util.DbTableProducts),
	}

	_, err := util.DynamoClient.DeleteItem(input)

	if err != nil {
		fmt.Println("Got error calling DeleteItem:")
		fmt.Println(err.Error())
	}

	return err
}
