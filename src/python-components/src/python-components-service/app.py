# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: MIT-0
# AWS X-ray support
#

from aws_xray_sdk.core import xray_recorder
from aws_xray_sdk.ext.flask.middleware import XRayMiddleware
from aws_xray_sdk.core import patch_all

patch_all()

xray_recorder.begin_segment("Videos-init")

STATIC_FOLDER = '/app/static'
STATIC_URL_PATH = '/static'

from flask import Flask
from flask_cors import CORS
from common.common import LoggingMiddleware
from routes.offers_routes import load_offers 
from routes.location_routes import load_s3_data
from routes import create_app
from routes.videos_routes import start_streams
import logging

EXPERIMENTATION_LOGGING = True
DEBUG_LOGGING = True

app = create_app()
logger = app.logger
corps = CORS(app, expose_headers=['X-Experiment-Name', 'X-Experiment-Type', 'X-Experiment-Id', 'X-Personalize-Recipe'])

xray_recorder.configure(service='PythonComponent Service')
XRayMiddleware(app, xray_recorder)

if __name__ == '__main__':

    if DEBUG_LOGGING:
        level = logging.DEBUG
    else:
        level = logging.INFO
    app.logger.setLevel(level)
    if EXPERIMENTATION_LOGGING:
        logging.getLogger('experimentation').setLevel(level=level)
        logging.getLogger('experimentation.experiment_manager').setLevel(level=level)
        for handler in app.logger.handlers:
            logging.getLogger('experimentation').addHandler(handler)
            logging.getLogger('experimentation.experiment_manager').addHandler(handler)
            handler.setLevel(level)  # this will get the main app logs to CloudWatch

    app.wsgi_app = LoggingMiddleware(app.wsgi_app)

    load_offers()
    load_s3_data()

    app.logger.info("Starting video streams")
    start_streams(app)

    app.logger.info("Starting API")
    app.run(debug=True, host='0.0.0.0', port=80)