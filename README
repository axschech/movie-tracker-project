# Movie Tracker

## Run the application
- Make sure Docker is running
- Run `make build up migrate logs`
- To stop the application run `make down`
- Make sure to run `make migrate` after bringing the application back up after its been down to recreate tables.

## Flow
Use a base URL of `localhost:8080` if you are running in Docker.
1. Create a user with the register endpoint. Make a POST request to `/api/user` with the payload `{"username":"anothertest","email": "anothertest@test.com" }`
2. Search for media with the query media endpoint. Make a GET request to `/api/search/media?query=<query>&type=tv`
3. Using the ID retrieved, associate the ID with a user by making a post request to `/api/media/user` with a payload `{ "user_id": 2, "media_id": 13, "status": "watched"}`, where `user_id` is the user to associate the media with, `media_id` is the media item to associate the user with, and `status` is the status of the pair.
4. Retrieve the user media pairings for a user by making a GET request to `/api/media/users/<user id>`. This will return all media for each pair found, with the status. This would be served to the frontend in a user friendly manner.

## Issues
- There is no authorization or authentication. I would use a service like Auth0 or something self hosted like Keycloak running separately outside of this application. The I would write middleware to verify tokens and guard against users accessing data that doesn't belong to them.
- There are issues with differences between how lax The TVDB is with user input and what the "WHERE LIKE" clause returns. For example, "Wonderman" - whose correct name is "Wonder Man" - will return "Wonder Man" from the API but would not be found with the "WHERE LIKE" clause. I'm not sure how to get around this currently. I might end up using something like elastic search for this instead of storing all the cached data in the DB.
- There is inconsistency between what supports bulk operations and what does not. This is for time savings.
- Some of the structs have too many concerns, I would separate most of the repository functions and service handler functions by domain.
- I need to add more data sources.
- I need to add unit tests. I would use the httptest package to call the handler functions (with a mocked DB). This would give the most coverage the fastest.
- There is no frontend, that is coming soon! 

## Learnings
- This was actually more challenging than I originally realized. I am pretty rusty on writing queries and making sure I wasn't duplicating data was important.
- I originally thought I would have time for a frontend and configs for a k8s deploy. I don't think this would be too challenging to run since there be only two containers, the app container and the db container.

## Progress Log
You can see more of this in the git history of the repo. I worked about two hours a day, with about 4 hours on Saturday and Sunday.
- It took about a day to get the service running and saving user data.
- It took about two days to get to querying and saving media

## Structure

### cmd/app
- Represents "app" command
- App runs from here
- Make sure the DB is up and running before routing starts

### Internal

#### Config
- Uses the spf13/viper package to read config from environment
- When moving to k8s I would use a config file with values read in from Google Secret Manager

#### Database
- Responsible for all database functionality
- I decided against using an ORM because I'm only doing a few queries and it was faster to implement with PGX
- Probably could be improved using GORM and preload for nested structures

*Repository*
- This is the actual SQL to get data
- Probably could spread this out logically by domain, didn't have enough time

#### Entities
- Data structures to keep things reusable and consistent between repository, business logic, and handler

#### Routing
- Set up for Chi Router

#### Service
- Holds route handlers and routing logic
- Would eventually need to split this up by domain
- Tried to keep business logic as left out of this layer as possible
- Only validation and calling business logic functions

#### Business Logic
- Media, Media User, and User
- Even the functions that are just mostly wrappers for repository functions would allow for easy modification / replacement if I had to change how I was accessing the database (this is a lesson I have learned).

### External
- Sources for getting information
- I started with The TVDB
- Would expand to include more sources
- Not sure if I would expose sources to the user, might be necessary since, for example, The TVDB also returns movie information, so it's not about 

#### Sourcer Interface
- This is a simple interface to make sure that there is always a "fetch" function on a source
- I went back and forth with this a lot. I thought I could use generics to allow for a default return type but ended up just keeping it an HTTP response. I thought that was better anyway because it should just be responsible for making a request and receiving a response

#### Migration
- This is just an SQL file to create the tables and insert some sample data.