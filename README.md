# service-partition

The partition service repository contains code that enables multitenancy for ant investor services. 
This service will answer the question related to which partition does request x or user y get confined in?

### How do I get set up? ###

* Database fixtures and migrations are combined and will be run automatically before the container starts during deployments.

* Dependencies
    Running this project requires an sql database and access to the authorization service
    
* How to run tests execute the command :
    `go test ./...`
    
* Deployment :

    This is an automated process done via github actions.
    Once the code is ready for deployment, merging to develop deploys to the staging environment.
    Merging to master deploys to production so all changes must be thoroughly validated before merging.
    

