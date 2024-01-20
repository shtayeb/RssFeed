# [RSS FEED AGGREGATOR](https://rssfeed.cyclic.app/)
An RSS feed aggragtor built with GO.

View the site [Analytics](https://rssfeed.goatcounter.com/)

## Tech Stack
- GO
- Templ
- TailwindCSS
- v0.dev
- GoatCounter for Analytics


## TODO
- [x] Use session management package
- [x] Setup Middleware
- [x] Setup Middleware group with go-chi for authenticated routes
- [x] (Using a package) Setup session management - DB table - cookies - hashed token 
- [x] In login handler create a new session entry to the sessions table
- [x] Inlcude the session data in the context for the views
- [x] Link Middleware and session to get auth user
-[x] Get the user details from context in the templ files
- Create a logout handler
    - [x] Delete the session when user logs out
    - If user is loggedin redirect them to referrer
-[x] Build a dark theme based on color inverting check the extention for reference
- [x] /posts
- [x] Create page for post list, 
- [x] /feed/{feedId}/posts -> make Post templ customizable

- feed following pages
- Create feed follow CRUD
- /user/feeds scoped to the user
- User profile management  
    - [x] Create a user profile dropdown for auth user
    - [x] Error handling with htmx
    - Fix Change password error: No consistent hashing for passwords

- CSRF token protection for the forms
- Check django requres object for any other userful data to be included in the context
- Check the laravel readirect()->back() and implement it for the GO
