The schema file is what determines what the databse should look like. 
When you run `atlas schema apply`, it will create the tables and fields in the database based on this schema file.
It compared the database with the schema file and applies any changes needed to make them match.
### Create the atlas schema file: 
`
atlas schema inspect -u "postgres://postgres:cheese@localhost:5434/brandongarate?sslmode=disable" > modules/anon_sessions/migrations/atlas_schema.hcl
`
### Create the tables
`
atlas schema apply -u "postgres://postgres:cheese@localhost:5434/brandongarate?sslmode=disable" --to file://modules/anon_sessions/migrations/atlas_schema.hcl
`

