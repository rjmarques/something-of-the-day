# Create a new Heroku app
resource "heroku_app" "default" {
  name = "something-of-the-day"
  region = "eu"
}

# Create a database, and configure the app to use it
resource "heroku_addon" "database" {
  app  = heroku_app.default.name
  plan = "heroku-postgresql:hobby-dev"
}

data "heroku_app" "sotd" {
  name = heroku_app.default.name
}