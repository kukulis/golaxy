# game of galaxy in the web

## Small MVP: build a fleet, using given resources and compete other fleets.


## Scenario for creating fleet_build

API

1) create division and assign maxResources.
2) create new fleet build with a given division id
3) create new ship model
4) Update ship model with all parameters
5) update fleet build with the technologies parameters
6) assign ship model to a fleet build 
7) modify ship model to fleet build  assignment amount
8) get fleet build statistics

Need a data structure for (user + division) pair to be able to store a built fleet
User id is not provided in the request. User id will be loaded by authorization token; for development it will be fixed value.

9) build fleet of the given fleet build id, and upsert division+user -> fleet.
10) Get fleet by the given division id + user id 
 