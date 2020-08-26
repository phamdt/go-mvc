[psql]
  dbname = "{{dbName}}"
  host   = "0.0.0.0"
  port   = {{dbPort}}
  user   = "{{dbUser}}"
  pass   = "{{dbPassword}}"
  blacklist = [
    {{#denyList}}{{denyList}}{{/denyList}}
  ]
  sslmode = "disable"
{{#templates}}
  templates = [
    {{templates}}
  ]
{{/templates}}
