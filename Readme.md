## TODO
-[x] Use session management package
- [x] Setup Middleware
- [x] Setup Middleware group with go-chi for authenticated routes
- [x] (Using a package) Setup session management - DB table - cookies - hashed token 
- [x] In login handler create a new session entry to the sessions table
- [x] Inlcude the session data in the context for the views
- [x] Link Middleware and session to get auth user

- Get the user details from context in the templ files
- Check django for any other userful data to be included in the context
- Create a logout handler
    - Delete the session when user logs out
    - If user is loggedin redirect them back with a mesasge

- CSRF token protection for the forms
