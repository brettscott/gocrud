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

* Code: [./sample_app.go](./sample_app.go)
* Execute: `run.sh -l`
* View: [http://localhost:8080](http://localhost:8080)  _(where localhost is IP of your Docker VM)_


## Todo

Everythings!

* API
    * GET
    * POST
    * PUT
    * PATCH
    * DELETE
* Front-end
    * Create
    * Results
    * Update
    * Delete
* Mongo
    * create indexes

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
