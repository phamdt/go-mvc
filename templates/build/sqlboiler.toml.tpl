[psql]
  dbname = "{{dbName}}"
  host   = "0.0.0.0"
  port   = {{dbPort}}
  user   = "{{dbUser}}"
  pass   = "{{dbPassword}}"
  blacklist = [
    {{#blacklist}}{{blacklist}}{{/blacklist}}
  ]
  sslmode = "disable"
{{#templates}}
  templates = [
    {{templates}}
  ]
{{/templates}}
