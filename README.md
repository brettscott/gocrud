# Go CRUD

A CRUD library

## Features

- Create record
- Update record
- Delete record
- Duplicate record
- List records
- Pagination
- Database agnostic
- Domain agnostic


## Run

Runs within Docker containers locally.  Assumes you have Docker installed.

* `./run.sh -l` - run locally
* `./run.sh -lx` - run locally without building containers
* `./run.sh` - run tests


## Example

* Code: [./examples/basic_examples.go](./examples/basic_examples.go)
* Execute: `run.sh -l`
* View: [http://localhost:8080/gocrud](http://localhost:8080/gocrud)  _(where localhost is IP of your Docker VM)_
    * [http://localhost:8080/gocrud/healtcheck](http://localhost:8080/gocrud/healthcheck)
    * [http://localhost:8080/gocrud/api/users](http://localhost:8080/gocrud/api/users)
    * [http://localhost:8080/gocrud/api/users/59ba5a8a0000000000000000](http://localhost:8080/gocrud/api/users/59ba5a8a0000000000000000) _(example record ID)
* Post:
    ```
    POST http://localhost:8080/gocrud/api/users
    Content-Type: application/json
    {
    	"keyValues": [
    		{
    			"key": "name",
    			"value": "Brett"
    		}
    	]
    }
    ```
* Put:
    ```
    PUT http://localhost:8080/gocrud/api/users/59ba5a8a0000000000000000
    Content-Type: application/json
    {
    	"keyValues": [
    		{
    			"key": "name",
    			"value": "Brett"
    		}
    	]
    }
    ```    
* View in MongoDB:  
    ```
    mongo
    > use gocrud;
    > db.users.find({});
    ```

## Todo

Everythings!

* Update Chi to V3 (replaces /:id with /{id})
* API
    * GET
        * Write Mongo integration test
    * POST
        * Write Mongo integration test
    * PUT
        * Write Mongo integration test
    * PATCH
        * **NEXT**
        * age is float64 but want int. //
    * DELETE
    * LIST
        * Write Mongo integration test
        * Pagination
    * Authentication and authorisation
* Front-end
    * Create
    * Results
    * Update
    * Delete
* Mongo
    * create indexes

## Packages

* crud
    * wire up CRUD for consumption by Developer
* entity
    * define a flexible schema
    * ability to model DB schema into a generic data model for crud-ing
* api
    * handles HTTP requests to REST interface
* store
    * database abstraction - allow Developer to choose database eg Mongo, MySQL
    
*Workflows:*

* POST:
    * receive request (/api) <- crud pkg
    * route request (post) <- api pkg 
    * req.body -> entity.record <- api pkg for now
    * validation <- api pkg for now
    * entity.record -> database query <- api pkg for now
    * database query -> mongo client <- store pkg
    



## Schema

- Entity eg User
    - Element eg Name
        - Identifier eg name
        - Form Attributes
            - Input Type eg text
            - Label eg Name
        - Results Attributes (individual element)
            - Column Order eg 0 (first)
            - Sort eg ascending
    - Results Attributes (all)
        - Responsive CSS?
    - Create Hook
    - Update Hook
    - Delete Hook
    - Plugins
    
- Store
    - list/get/post/put/delete
    - hooks


## Terminology

* Entity - the content being CRUDed eg users, permissions, intelligence.
* Element - a child of Entity.  An Entity will typically have more than one element eg id, name, description.
* Record - represents an Entity Item made up of populated Element `value`s. Think of a Record as a row in a database table.