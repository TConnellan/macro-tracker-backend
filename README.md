# macro-tracker-backend








## Stream of conscious log


### On setting up automated integration tests for 

- Goal is to fully automate this testing in some way, setup of DB, running of migrations, and running of tests
- Some considerations 
    - GO by default runs tests in a test file in serial, but will run test files in parallel
        - There is some functionality with dependencies between models, that is the recipes model also interacts with the recipe_compenents and consumables tables
    - The pattern that I started working with originally had a new db connection established with the same existing DB for each test in a file and data set up and torn down after each subtest execution, oh no race conditions in my tests if I'm not careful
        - I was concerned about the teardown of tests in one file effecting the data of tests running in parallel in another, thought of a few options
            - Working around this inside the current framework felt wrong, forcing serial execution of everything felt like a copout, and other options such as forcing separate tests to rely on separate data would just create more issues in the future
            - Decided to try and solve by isolating integration tests
                - Could force each test to run in a transaction which is always rolled back on cleanup
                    - This would force any tests that are running in parallel to wait for rollback of other queries which might effect the data they need 
                    - Unfortunately this won't work with the design of the models, they all take a pgx pool, while this method would require them to take an already created transaction, this could be refactored, but:
                        1. would require changing the signature of everything that accepts Pool to an interface that specifies all the shared functionality of a pool and a transaction, such an interface already exists but it's only in use for reusable helpers where a transaction might want to be passed and it doesn't yet specify a method for beginning a transaction.
                        2. Everywhere a transaction is already being created would need to be refactored to use Begin() rather than BeginTX(), but this would not allow the passing of isolation levels that is currently being done
                            - Furthermore this would change the functionality between the application and testing. This is because Postgres (and most DBMS) don't actually implement subtransactions. When the tests run a savepoint would be created whenever a transaction was expected to be started, while these provide similar rollback capabilities with parts of a transaction they are not the same and the testing should mirror the application
                - Could force each test to run in a completely separate database, as part of the setup create a completely new database and tear it down when done.                
- So to avoid refactoring that just creates more potential problems (even if refactored correctly) decided to take the approach of generating new test databases for each test. Will see what the overhead looks like.