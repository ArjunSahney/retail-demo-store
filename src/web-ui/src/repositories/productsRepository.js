// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

import axios from "axios";
import resolveBaseURL from './resolveBaseURL'

const baseURL = resolveBaseURL(
    import.meta.env.VITE_GO_COMPONENTS_SERVICE_DOMAIN,
    import.meta.env.VITE_GO_COMPONENTS_SERVICE_PORT,
    import.meta.env.VITE_GO_COMPONENTS_SERVICE_PATH
)

const connection = axios.create({
    baseURL
})

const resource = "/products";
export default {
    get() {
        return connection.get(`${resource}/all`)
    },
    getFeatured() {
        return connection.get(`${resource}/featured`)
    },
    getProduct(productID) {
        if (!productID || productID.length == 0)
            throw "productID required"
        if (Array.isArray(productID))
            productID = productID.join()
        return connection.get(`${resource}/id/${productID}`)
    },
    getProductsByCategory(categoryName) {
        if (!categoryName || categoryName.length == 0)
            throw "categoryName required"
        return connection.get(`${resource}/category/${categoryName}`)
    },
    getCategories() {
        return connection.get(`categories/all`)
    }
}