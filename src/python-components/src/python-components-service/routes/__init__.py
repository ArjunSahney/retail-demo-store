from flask import Flask

STATIC_FOLDER = '/app/static'
STATIC_URL_PATH = '/static'

def create_app():
    app = Flask(__name__,
            static_folder=STATIC_FOLDER,
            static_url_path=STATIC_URL_PATH)

    # Import and register the blueprints
    from routes.offers_routes import offers_blueprint
    app.register_blueprint(offers_blueprint)

    from routes.location_routes import location_blueprint
    app.register_blueprint(location_blueprint)

    from routes.search_routes import search_blueprint
    app.register_blueprint(search_blueprint)

    from routes.videos_routes import videos_blueprint
    app.register_blueprint(videos_blueprint)    
    
    from routes.recommendations_routes import recommendations_blueprint
    app.register_blueprint(recommendations_blueprint)
    
    return app