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

Running tests within your IDE, start Mongo:

    docker-compose run --service-ports mongodb

Create an entry in your `hosts` file for `mongodb` to point to your Docker VM's IP address.


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

* API
    * LIST
        * Pagination
    * More validation rules
    * More examples (eg basic_example.go)
    * Authentication and authorisation
    * Support for search - ES?
    * Convert Entity setup into JSON which is unmarshalled into Entity
    * Return _crud information in KV to client
* Front-end
    * Create
    * Results (/list) 
        * edit link is not correct because KeyValues is an array and cannot easily get ID value.  `row.go` struct.  **NEXT** 
    * Update
    * Delete
* Mongo
    * create indexes
* SQL
    * Add a SQL adapter as a proof of concept to ensure they can be easily created
* Neo4j
    * Add adapter
* Move /crud/ folder into root (/), rename package as `gocrud` in order for consumers to use `gocrud.` instead of `crud.` (or `gocrud.crud.`?).
* Add statsd metrics
* Review packages and dependencies (cyclic) 
    * https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
    * http://idiomaticgo.com/post/best-practice/server-project-layout/
    
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



## Structs

```
  Browser        <-->         Application         <--->        Store/Database  
  -------                     -----------                      --------------
   
  API Route:
  
            ClientRecord                    StoreRecord
            
            
  UI Route:    
  
                   Row        ClientRecord       StoreRecord
                   
```                   
                 
