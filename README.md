1. Login - who can do all the operation.
2. Listing of releases - which will contains list of releases
3. Creation of release - will have creation form with below fields
    -> Release name -> unique
    -> Release target Date(Release date) 31-05-2021
    -> Project details (Hybris/Ecommerce/ReactNative/React)
    -> Type of Release (hot fix/feature/release)
    -> Owner Email
4. Project details, Which will have the following details
    -> Bitbucket url
    -> Reviewer List
    -> Release Date : T
    -> Beta release: T-1
    -> Regression signor: T-2 
    -> Code Freeze: T -4 
    -> Dev Completion: T-5
5. Third party integration Bitbucket and Jira
6. Schedular - for milestone updates
7. Create DLs - Will be seeder for now. 


To run the project 
1. copy .env_sample to .env file
2. run - go run ./cmd/web/main.go
