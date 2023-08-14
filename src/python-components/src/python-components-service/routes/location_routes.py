# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0

from flask import Flask
from flask import Blueprint
from flask import abort, jsonify, request
from flask_cors import CORS

import boto3
import json
import os
import pprint

location_blueprint = Blueprint('location', __name__)

RESOURCE_BUCKET = os.environ.get('RESOURCE_BUCKET')

s3 = boto3.resource('s3')

store_location = {}
customer_route = {}

cstore_location = {}
cstore_route = {}


def load_s3_data():
    global customer_route
    route_file_obj = s3.Object(RESOURCE_BUCKET, 'retaildemostore/location_services/customer_route.json')
    customer_route = json.loads(route_file_obj.get()['Body'].read().decode('utf-8'))

    global store_location
    location_file_obj = s3.Object(RESOURCE_BUCKET, 'retaildemostore/location_services/store_location.json')
    store_location = json.loads(location_file_obj.get()['Body'].read().decode('utf-8'))

    global cstore_route
    route_file_obj = s3.Object(RESOURCE_BUCKET, 'retaildemostore/location_services/cstore_route.json')
    cstore_route = json.loads(route_file_obj.get()['Body'].read().decode('utf-8'))

    global cstore_location
    route_file_obj = s3.Object(RESOURCE_BUCKET, 'retaildemostore/location_services/cstore_location.json')
    cstore_location = json.loads(route_file_obj.get()['Body'].read().decode('utf-8'))


@location_blueprint.route('/checkLocation')
def index():
    return 'Location Service Service'


@location_blueprint.route('/store_location')
def get_store_location():
    return jsonify(store_location)


@location_blueprint.route('/customer_route')
def get_customer_route():
    return jsonify(customer_route)


@location_blueprint.route('/cstore_location')
def get_cstore_location():
    return jsonify(cstore_location)


@location_blueprint.route('/cstore_route')
def get_cstore_route():
    return jsonify(cstore_route)
