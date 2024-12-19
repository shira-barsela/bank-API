# bank-API

--- API ---
I added a var of loginUsersToken. It's a map that for every user that login the server, saves for the user id, it's token.
Now, every action that is allowed only by a specific user (getBalance, depositBalance, withdrawBalance),
I changed these function to check that the token given in the request, is the one of the user ID that asked for the action.


--- Main ---
The main function declare the handler function that should be read at each type of request.
It's ready to except and serve any request.
All the handler function are 'wraped' with a logging function, that will write the required data in a log file.
