from flask import Blueprint
from flask import abort, jsonify, request
import json

offers_blueprint = Blueprint('offers', __name__)

offers = []

def load_offers():
    global offers
    with open('data/offers.json') as f:
        offers = json.load(f)

@offers_blueprint.route('/')
def index():
    return 'Offers Service'


@offers_blueprint.route('/offers')
def get_offers():
    return jsonify({'tasks': offers})


@offers_blueprint.route('/offers/<offer_id>')
def get_offer(offer_id):
    for offer in offers:
        if offer['id'] == int(offer_id):
            return jsonify({'task': offer})
    abort(404)
