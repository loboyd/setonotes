Goal: Build a messaging services which can be interacted through via an API.

`messages` table
----------------
* `id`
* `from_user`
* `to_user`
* `timestamp`
* `body`

API functionality
-----------------
* send a message to a user (probably a POST request)
  * Sending client makes a request
  * The server authenticates the client+user
  * Does the receiver heve to be verified, etc?
  * Database is updated with message (and notification?)
  * [eventually] A success/error message is returned in reponse
* receive message(s) (probably a GET request)
  * Receiving client makes a request
  * The server authenticates the client+user
  * Database is queried for all relevant messages
  * Messages are returned in response
