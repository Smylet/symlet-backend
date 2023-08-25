
# Handlers

This directory contains all the handlers for the API endpoints. The handlers are grouped into files based on the type of endpoint they handle.
The purpose of separating the handlers into a different directory is to make it easier to access the properties of the server such as the `database connection`, the `logger`, `cron`, `config` etc. from the handlers.

## USER ACCOUNT & AUTHENTICATION ENDPOINTS

`auth.go` contains all the handlers related to user account management and authentication.

1. **Register User:**
    - `POST /users/register`
    - Creates a new user profile with `username`, `email`, `password`, etc.
  
2. **Login User:**
    - `POST /users/login`
    - Authenticates a user with `username` or `email` and `password`.
  
3. **Confirm Email:**
    - `POST /users/confirm-email`
    - Confirms user email using a `confirmation_token`.
  
4. **Resend Email Confirmation:**
    - `POST /users/email/confirmation/resend`
    - Resends the email confirmation token using the user's `email`.
  
5. **Request Password Reset:**
    - `POST /users/password-reset`
    - Initiates password reset using the user's `email`.
  
6. **Change Password:**
    - `PUT /users/password-change`
    - Changes the logged-in user's password using old and new passwords.
  
7. **Setup Two-Factor Authentication:**
    - `POST /users/:username/2fa/setup`
    - Sets up two-factor authentication for a user.
  
8. **Verify Two-Factor Authentication:**
    - `POST /users/:username/2fa/verify`
    - Verifies the provided `2fa_code` for authentication.
  
9. **Logout:**
    - `POST /users/logout`
    - Logs out the user, invalidating their session or token.

## USER PROFILE MANAGEMENT ENDPOINTS

`profile.go` contains all the handlers related to user profile management.

10. **Get User Profile:**
    - `GET /users/:username/profile`
    - Retrieves a user's profile details.
  
11. **Edit User Profile:**
    - `PUT /users/:username/profile`
    - Updates the user's profile with provided details.
  
12. **Delete User Profile:**
    - `DELETE /users/:username`
    - Deletes the user's profile.
  
13. **View Profile Edit History:**
    - `GET /users/:username/history`
    - Shows changes made to the user's profile over time.
  
14. **Backup User Profile:**
    - `POST /users/:username/backup`
    - Creates a backup of the user's profile.
  
15. **List Profile Backups:**
    - `GET /users/:username/backups`
    - Lists all backups related to the user's profile.
  
16. **Restore User Profile:**
    - `PUT /users/:username/restore`
    - Restores the user's profile from a provided backup.
  
17. **Export User Profile:**
    - `GET /users/:username/export`
    - Exports the user's profile data in formats like JSON, CSV.

19. **Deactivate Account:**
    - `PUT /users/:username/deactivate`
    - Temporarily deactivates a user's account.
  
20. **Reactivate Account:**
    - `PUT /users/:username/reactivate`
    - Reactivates a previously deactivated account.

## USER PROFILE PICTURE ENDPOINTS

21. **Upload Profile Picture:**
    - `POST /users/:username/picture`
    - Uploads a new profile picture for the user.
  
22. **Update Profile Picture:**
    - `PUT /users/:username/picture`
    - Replaces the existing profile picture with a new one.
  
23. **Delete Profile Picture:**
    - `DELETE /users/:username/picture`
    - Removes the user's profile picture.

## USER SOCIAL & NOTIFICATIONS ENDPOINTS

`notification.go` contains all the handlers related to user notifications.

24. **Follow a User:**
    - `POST /users/:username/follow`
    - Allows a user to follow another user.
  
25. **Unfollow a User:**
    - `DELETE /users/:username/unfollow`
    - Allows a user to unfollow another user.
  
26. **View Followers:**
    - `GET /users/:username/followers`
    - Displays a list of users following the user.
  
27. **View Following:**
    - `GET /users/:username/following`
    - Displays a list of users the user is following.
  
28. **Send Friend Request:**
    - `POST /users/:username/friend-request`
    - Sends a friend request to another user.
  
29. **Accept Friend Request:**
    - `POST /users/:username/friend-request/accept`
    - Accepts a friend request from another user.
  
30. **Reject Friend Request:**
    - `DELETE /users/:username/friend-request/reject`
    - Rejects a friend request from another user.
  
31. **Get Notifications:**
    - `GET /users/:username/notifications`
    - Retrieves all notifications for the user.
  
32. **Mark Notification as Read:**
    - `PUT /users/:username/notifications/:notificationId/read`
    - Marks a specific notification as read for the user.

### USER PRIVACY SETTINGS ENDPOINTS

33. **Set Profile Privacy:**
    - `PUT /users/:username/privacy`
    - Sets the user's profile privacy settings.

### USER SEARCH ENDPOINTS

34. **Search Users:**
    - `GET /users/search`
    - Searches for users based on a query string.

### USER BLOCKING ENDPOINTS

35. **Block User:**
    - `POST /users/:username/block`
    - Blocks a user from interacting with the user.

36. **Unblock User:**
    - `DELETE /users/:username/unblock`
    - Unblocks a user from interacting with the user.