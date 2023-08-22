# user service


## File structure overview

```sh
users/
├── doc.go             # Documentation related to the project
├── router/
│   ├── auth.go        # Authentication related endpoints   
│   ├── profile.go     # Profile related endpoints
│   ├── notification.go      # Notification related endpoints
│   ├── others.go      # Other user related endpoints - miscellanous
├── middleware.go    # All middleware functions and structures
├── model.go         # Data models and database related functions
├── router.go        # Routing and endpoint declarations
├── serializer.go    # Functions to convert data into serializable formats
├── unit_test.go      # Unit tests for various components of the project
└── validator.go     # Input validation logic and structures
└── task.go           # Background tasks
└── README.md           
```

1. **Documentation** (from `doc.go`):
   - **User Guide**: Step-by-step instructions on how to manage and edit profiles.
   - **FAQ**: Common questions related to account setup, privacy, and troubleshooting.
   - **Change Log**: Historical changes or updates related to profile management.
2. **Middleware Functions** (from `middlewares.go`):
   - **Authentication Middleware**: Ensure the user is authenticated before allowing profile changes.
   - **Error Handling**: Capture and handle errors during profile edits or updates.
   - **Rate Limiting**: Prevent excessive requests to avoid spam or abuse.
3. **Data Models & Database Operations** (from `models.go`):
   - **User Profile Data Structure**: Define fields like username, email, profile picture, bio, etc.
   - **Profile Edit History**: Maintain a history of changes made to a user profile.
   - **Backup and Restore**: Allow users to backup their profile data and restore if necessary.
4. **Routing and Endpoints** (from `routers.go`):
   - **Profile Viewing Endpoint**: Access a user's public profile.
   - **Profile Editing Endpoint**: Allows users to modify their personal details.
   - **Account Deletion Endpoint**: Facilitate users in account removal.
5. **Data Serialization** (from `serializers.go`):
   - **Profile Export**: Allow users to export their profile data in various formats (e.g., JSON, CSV).
   - **Privacy Filters**: Ensure sensitive information is not included in serialized data without permission.
   - **Profile Preview**: Preview how the profile appears to others post-serialization.
6. **Unit Testing** (from `unit_test.go`):
   - **Profile Edit Tests**: Ensure that changes to profiles work as expected.
   - **Data Integrity Tests**: Confirm that saved profile data remains consistent.
   - **Security Tests**: Verify that unauthorized users can't modify profiles.
7. **Input Validation** (from `validators.go`):
   - **Field Checks**: Ensure fields like email, password, and username meet criteria.
   - **Profile Picture Validation**: Confirm uploaded images are of the right format and size.
   - **Data Consistency Checks**: Ensure that user input is consistent and logical (e.g., date of birth is valid).


## ENDPOINTS

### USER ACCOUNT & AUTHENTICATION ENDPOINTS

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

### USER PROFILE MANAGEMENT ENDPOINTS

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

### USER PROFILE PICTURE ENDPOINTS

21. **Upload Profile Picture:**
    - `POST /users/:username/picture`
    - Uploads a new profile picture for the user.
  
22. **Update Profile Picture:**
    - `PUT /users/:username/picture`
    - Replaces the existing profile picture with a new one.
  
23. **Delete Profile Picture:**
    - `DELETE /users/:username/picture`
    - Removes the user's profile picture.

### USER SOCIAL & NOTIFICATIONS ENDPOINTS

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


## DATA MODELS

```sql
-- Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(512) NOT NULL,
    email_confirmed BOOLEAN DEFAULT FALSE,
    2fa_enabled BOOLEAN DEFAULT FALSE,
    account_status VARCHAR(50) DEFAULT 'active',
    profile_privacy VARCHAR(50) DEFAULT 'public',
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- User Profiles Table with additional details
CREATE TABLE user_profiles (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    bio TEXT,
    dob DATE,
    address TEXT,
    picture_url VARCHAR(2048),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- Profile Edit History Table with JSONB for change details
CREATE TABLE profile_edit_history (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    changes JSONB,
    edited_at TIMESTAMP DEFAULT NOW()
) PARTITION BY RANGE (edited_at);

-- Monthly partitions for the profile edit history
-- Adjust based on expected frequency of edits.
-- Additional partitions can be added later.
-- This helps in optimizing query performance.
CREATE TABLE profile_edit_history_2023_08 PARTITION OF profile_edit_history FOR VALUES FROM ('2023-08-01') TO ('2023-09-01');

-- ... Add more partitions as time goes on ...

-- Backups Table
CREATE TABLE profile_backups (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    backup_data JSONB,
    backup_at TIMESTAMP DEFAULT NOW()
) PARTITION BY RANGE (backup_at);

-- Monthly partitions for the profile backups
CREATE TABLE profile_backups_2023_08 PARTITION OF profile_backups FOR VALUES FROM ('2023-08-01') TO ('2023-09-01');

-- Social (Followers & Following) Table
CREATE TABLE user_social (
    follower_id INT REFERENCES users(id),
    following_id INT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id)
);

-- Friend Requests Table
CREATE TABLE friend_requests (
    requester_id INT REFERENCES users(id),
    requestee_id INT REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

-- Notifications Table
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    content TEXT,
    notification_type VARCHAR(50),
    status VARCHAR(50) DEFAULT 'unread',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
) PARTITION BY RANGE (created_at);

-- Monthly partitions for notifications
CREATE TABLE notifications_2023_08 PARTITION OF notifications FOR VALUES FROM ('2023-08-01') TO ('2023-09-01');

-- Blocked Users Table
CREATE TABLE blocked_users (
    user_id INT REFERENCES users(id),
    blocked_user_id INT REFERENCES users(id),
    reason TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, blocked_user_id)
);

-- Indices for faster search and access

-- Text-search might be frequent on usernames and emails
CREATE INDEX idx_username ON users USING gin(username gin_trgm_ops);
CREATE INDEX idx_email ON users USING gin(email gin_trgm_ops);

-- For frequent profile lookups and joins
CREATE INDEX idx_user_profiles_user_id ON user_profiles(user_id);

-- For social features and determining relationships
CREATE INDEX idx_user_social_follower ON user_social(follower_id);
CREATE INDEX idx_user_social_following ON user_social(following_id);

-- Depending on usage patterns, consider more indexes on the friend_requests and blocked_users tables

```

## DATA VALIDATION

### USER ACCOUNT & AUTHENTICATION ENDPOINTS VALIDATION

1. **Register User:**
    - `username`: Ensure it's unique, between 3-50 characters, alphanumeric with possible underscores but not starting with an underscore.
    - `email`: Check for a valid email format and ensure uniqueness.
    - `password`: Minimum 8 characters, at least one uppercase, one lowercase, one number, and one special character (like !, @, #, etc.).
  
2. **Login User:**
    - `username` or `email`: Validate based on respective formats.
    - `password`: Should not be empty.
  
3. **Confirm Email:**
    - `confirmation_token`: Ensure token exists, hasn't expired, and matches a user.
  
4. **Resend Email Confirmation:**
    - `email`: Check for a valid email format and ensure it exists in the system.
  
5. **Request Password Reset:**
    - `email`: Validate email format and existence in the system.
  
6. **Change Password:**
    - Old and new passwords should be validated against standard password rules.
  
7. **Setup Two-Factor Authentication & Verify Two-Factor Authentication:**
    - `2fa_code`: Ensure it's a valid format (usually 6 numeric digits) and hasn't expired.

### USER PROFILE MANAGEMENT ENDPOINTS VALIDATION

10. **Get, Edit, Delete User Profile:**
    - `:username`: Ensure the username exists and the logged-in user has permission to access/modify.
  
11. **Edit User Profile:**
    - Fields like `first_name`, `last_name`, `bio`, etc. should be checked for valid lengths and character types.
    - `dob`: Validate as a valid date format and logical range (e.g., not in the future or more than 120 years ago).
  
12. **Backup, List Backups, Restore, Export, Preview, Deactivate, Reactivate Account:**
    - `:username`: Ensure it exists and the logged-in user has permission.
  
### USER PROFILE PICTURE ENDPOINTS VALIDATION

21-23. **Manage Profile Picture:** \
    - `:username`: Ensure it exists. \
    - For uploads: Check file type, ensure it's an image, limit file size, and scan for malware.

### USER SOCIAL & NOTIFICATIONS ENDPOINTS VALIDATION

24-30. **Social Interactions:** \
    - `:username`: Ensure it exists.
    - Check that the action makes sense (e.g., not following someone you're already following).
  
31. **Get Notifications & Mark Notification as Read:**
    - Ensure the logged-in user has permission.
    - For marking notifications: Ensure the notification ID exists.

### USER PRIVACY SETTINGS & SEARCH ENDPOINTS VALIDATION

33. **Set Profile Privacy & Search Users:**
    - `:username`: Validate existence and permission.
    - For search: Ensure the query string isn't excessively long or malicious.

### USER BLOCKING ENDPOINTS VALIDATION

35-36. **Block/Unblock User:**
    - `:username`: Ensure both users exist and the action is logical (e.g., not blocking someone who is already blocked).


## Background tasks

1. **Email Notifications and Verifications**:
   - Sending verification emails (`POST /users/confirm-email`).
   - Resending verification emails (`POST /users/email/confirmation/resend`).
   - Sending password reset links (`POST /users/password-reset`).

    **Why**: Email operations can often be time-consuming. Making the user wait for the email to be sent before getting a response from the server can degrade the user experience. By delegating this task to a background worker, the server can immediately respond to the user while the email is being sent in the background.

2. **User Profile Operations**:
   - User profile backups (`POST /users/:username/backup`).
   - Restoring user profile from a backup (`PUT /users/:username/restore`).
   - Exporting user profile data (`GET /users/:username/export`).

    **Why**: Creating backups, restoring data, and exporting large sets of data can be resource-intensive operations. Executing them asynchronously ensures that the web server remains responsive.

3. **Notifications**:
   - Generating and storing notifications, especially if multiple notifications are generated simultaneously or if the generation process involves complex computations.

    **Why**: Notifications might be triggered by various actions. For instance, when a user gets a new follower or receives a friend request. Handling the creation and dispatch of these notifications in the background can help in decoupling the main application logic from notification management.

4. **Image Processing**:
   - Profile picture uploads and updates (`POST /users/:username/picture` and `PUT /users/:username/picture`).

    **Why**: Image processing (resizing, format conversions, compression) can be CPU intensive. Offloading these tasks ensures the main application remains snappy and responsive.

5. **Database Maintenance**:
   - Managing the partitions for tables like `profile_edit_history` and `profile_backups` as time goes on.

    **Why**: Database maintenance tasks can be scheduled to run during off-peak hours using background workers, ensuring minimal disruption to users.

8. **Housekeeping**:
   - Clearing out old or obsolete data, purging old backups, or archiving old profile edit histories.

    **Why**: These are routine tasks that don't need to run in real-time and can be scheduled to run when the system load is low.

9. **Security Scans**:
    - Scanning uploaded files for malware or viruses.
    
     **Why**: These scans can be time-consuming and might be better suited for background workers.
    
10. **Data Exports**:
    - Exporting data in various formats (like JSON, CSV, etc.).
    
     **Why**: Exporting large amounts of data can be resource-intensive and might be better suited for background workers.


## Middleware


1. **Authentication Middleware**:

   - **Description**: Ensure the user is authenticated before allowing access to certain endpoints.
   - **Applies to**:
     - Editing user profile.
     - Deleting user profile.
     - Profile backups and restore.
     - Two-Factor Authentication setup and verification.
     - Profile picture operations.
     - Social & notification operations.
     - Privacy settings operations.
     - Blocking operations.
     - ...and more.

2. **Authorization Middleware**:

   - **Description**: Verify if the authenticated user has the right permissions to perform a specific operation. For instance, a user shouldn't be able to modify another user's profile.
   - **Applies to**: Any endpoint that requires specific user permissions.

3. **Rate Limiting Middleware**:
   - **Description**: Limit the number of requests a user or IP can send in a specific time window. This prevents abuse and protects the application from potential DDoS attacks.
   - **Applies to**: All endpoints, especially those that are exposed publicly like registration, login, and password reset.

4. **Error Handling Middleware**:

   - **Description**: Capture, handle, and log errors gracefully, ensuring that the end-users get a consistent and informative error message.
   - **Applies to**: All endpoints.

5. **Data Validation Middleware**:

   - **Description**: Ensure that the incoming request data is valid before processing it.
   - **Applies to**: 
     - Registration.
     - Profile editing.
     - Password reset and change.
     - Two-Factor Authentication setup.
     - Profile picture upload and update.
     - ...and more.

6. **Two-Factor Authentication (2FA) Middleware**:

   - **Description**: Check if a user has 2FA enabled and, if so, verify the provided 2FA code before granting access.
   - **Applies to**: Login and any other sensitive operations if you choose to protect them with 2FA.

7. **CORS Middleware**:

   - **Description**: If you have a frontend that's separated from the backend (like in many modern web apps), you might need to handle Cross-Origin Resource Sharing (CORS) settings to allow your frontend to communicate with the backend.
   - **Applies to**: All endpoints accessed by a web frontend.

8. **Logging Middleware**:

   - **Description**: Log requests, responses, and other significant events. This is crucial for debugging and monitoring.
   - **Applies to**: All endpoints.

9. **Request & Response Sanitization Middleware**:

   - **Description**: Remove or escape potentially harmful characters from the request to prevent SQL injection, Cross-site Scripting (XSS), etc., and ensure sensitive data doesn't leak in responses.
   - **Applies to**: All endpoints where user input is involved.

10. **Session Management Middleware**:

    - **Description**: Handle user sessions, either through cookies, JWT, or any other method you prefer.
    - **Applies to**: Endpoints that require session management like login, logout, and others.

11. **Backup & Restore Middleware**: 🔺

    - **Description**: A middleware to handle backup initiation, management, and restoration processes.
    - **Applies to**: Backup and restore endpoints.

12. **Cache Management Middleware**:

    - **Description**: For frequent and consistent data like user profiles, caching mechanisms can reduce database load and improve response times.
    - **Applies to**: Frequently accessed data endpoints, especially profile viewing.

14. **Data Serialization & Privacy Middleware**:

    - **Description**: Serialize data based on user's privacy settings, ensuring no private information is exposed unintentionally.
    - **Applies to**: Endpoints that return user data, like profile exports and previews.

## Scheduled tasks - cron job

1. **Unconfirmed Emails Cleanup**:
   - **Purpose**: Remove or flag user profiles where the email hasn't been confirmed for an extended period (e.g., 72 hours).
   - **Endpoint Impacted**: `POST /users/confirm-email`, `POST /users/email/confirmation/resend`
   - **Suggested Schedule**: Daily

2. **Backup Cleanup**:
   - **Purpose**: Delete old backups to conserve storage space based on a retention policy.
   - **Endpoint Impacted**: `POST /users/:username/backup`, `GET /users/:username/backups`
   - **Suggested Schedule**: Weekly or Monthly, based on backup frequency and retention needs.

3. **Old Account Deactivation or Removal**:
   - **Purpose**: Deactivate or remove user accounts that have been inactive for a prolonged period (e.g., 6 months or 1 year).
   - **Endpoint Impacted**: `PUT /users/:username/deactivate`
   - **Suggested Schedule**: Monthly

4. **Profile Edit History Cleanup**:
   - **Purpose**: If you have a policy of not retaining edit histories indefinitely, periodically clean up old records.
   - **Endpoint Impacted**: `GET /users/:username/history`
   - **Suggested Schedule**: Monthly or Quarterly

5. **Notifications Cleanup**:
   - **Purpose**: Remove old notifications to keep the database optimized.
   - **Endpoint Impacted**: `GET /users/:username/notifications`, `PUT /users/:username/notifications/:notificationId/read`
   - **Suggested Schedule**: Weekly or Monthly

6. **Database Maintenance**:
   - **Purpose**: Reindexing, optimizing tables, updating statistics, etc., to ensure the database performs optimally.
   - **Data Models Impacted**: All tables
   - **Suggested Schedule**: Monthly or Quarterly, during off-peak hours.



7. **Reminder Notifications**:
    - **Purpose**: Send reminders or prompts to users based on various conditions. For instance, nudge users who have not logged in for a long time or remind users to set up two-factor authentication.
    - **Endpoint Impacted**: `POST /users/:username/2fa/setup`, `POST /users/logout`
    - **Suggested Schedule**: Weekly or Monthly



## User story

After collecting common fields for the 3 users, users select the type of user they are and the fields that are relevant to them are displayed. The user can then fill in the fields and submit the form. The form is then validated and the user is redirected to the profile page.

For a Student:

- university

For a PropertyManager:

- company
- properties
