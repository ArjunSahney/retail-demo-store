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

const resource = "/carts";
export default {
    get() {
        return connection.get(`${resource}`)
    },
    getCartByID(cartID) {
        if (!cartID || cartID.length == 0)
            throw "cartID required"
        return connection.get(`${resource}/${cartID}`)
    },    
    createCart(username) {

        if (!username || username.length == 0)
            throw "username required"
        let payload = {
            username: username  
        }
        return connection.post(`${resource}`, payload)
    },    
    updateCart(cart) {
        if (!cart)
            throw "cart required"
        return connection.put(`${resource}/${cart.id}`, cart)
    },
    signAmazonPayPayload(payload) {
        return connection.post('/sign', payload)
    }
}