# [RSS FEED AGGREGATOR](https://rssfeed.cclic.app/)
An RSS feed aggragtor built with GO.

View the site [Analytics](https://rssfeed.goatcounter.com/)

## Tech Stack
- GO
- Templ
- TailwindCSS
- v0.dev
- GoatCounter for Analytics

```shell
tailwindcss --input tailwind.css --output test.css --content basic.templ
```

## TODO

### Not Started 

- Create feed follow/unfollow CRUD
- Email notification for your feed - when new post is added
- User followings feeds page
- /user/feeds scoped to the user

### Done
- [x] Pagination for posts list page
- [x] Use session management package
- [x] Setup Middleware
- [x] Setup Middleware group with go-chi for authenticated routes
- [x] (Using a package) Setup session management - DB table - cookies - hashed token 
- [x] In login handler create a new session entry to the sessions table
- [x] Inlcude the session data in the context for the views
- [x] Link Middleware and session to get auth user
- [x] Get the user details from context in the templ files
- [x] Create a logout handler
    - [x] Delete the session when user logs out
    - [x] If user is loggedin redirect them to referrer
- [x] Build a dark theme based on color inverting check the extention for reference
- [x] /posts
- [x] Create page for post list, 
- [x] /feed/{feedId}/posts -> make Post templ customizable
- [x] User profile management 
    - [x] Create a user profile dropdown for auth user
    - [x] Error handling with htmx
    - [x] Fix Change password error: No consistent hashing for passwords


### Backlog

- [] auth
    - [x] Testing reset-password 
        - Create a reset_password table and save the reset requests there. right now users can reset their passwords as much as they want 
        as long as their token is valid
    - [] Verify new users email

- CSRF token protection for the forms
- Check django requres object for any other userful data to be included in the context
- Check the laravel readirect()->back() and implement it for the GO
