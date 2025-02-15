# Flicker - Where Sparks Turn Into Stories

<!--  ![Flicker Logo](https://github.com/wanton-idol/Flicker/blob/main/assets/logo.png) --->
<p align="center">
  <img src="https://github.com/wanton-idol/Flicker/blob/main/assets/logo.png" alt="Flicker Logo" width="200">
</p>

## Introduction

**Flicker** is a **next-generation dating app** designed for Gen Z, where sparks ignite meaningful stories by blending social vibes with real connections. Built with **scalability, performance, and a microservice-driven architecture**, Flicker ensures seamless interactions, real-time communication, and a personalized user experience.

## Features

- **User Registration & Authentication** (Email, OTP, Google Sign-in, JWT-based authentication)
- **Profile Management** (User profile creation, updates, media handling)
- **Matchmaking & Swiping System**
- **Chat System** (Real-time messaging, read receipts)
- **Stories Feature** (Share moments & experiences)
- **Event System** (Create & join events)
- **Elasticsearch-Powered Search & Indexing**
- **Push Notifications** (AWS SNS integration)
- **Optimized API Performance**

## Tech Stack

- **Backend:** Golang (Go Gin framework)
- **Database:** MySQL
- **Caching:** Redis
- **Search Engine:** Elasticsearch
- **Cloud Services:** AWS (SNS for notifications, S3 for storage)
- **Auth & Security:** JWT, Twilio (OTP verification)
- **Containerization:** Docker
- **Migrations:** Go Migrate

## Architecture

Flicker follows a **microservice architecture**, ensuring modularity, scalability, and better maintenance. Key services include:

1. **User Service:** Handles authentication, profiles, and matchmaking.
2. **Chat Service:** Manages real-time messaging and chat history.
3. **Stories Service:** Indexing and retrieving user stories.
4. **Event Service:** Allows users to create and explore events.
5. **Notification Service:** AWS SNS-powered push notifications.

## API Endpoints

### User Authentication & Registration

- `POST /user/register` - Register a new user.
- `POST /user/login` - User login.
- `POST /user/google/login` - Google-based login.
- `POST /user/send/otp` - Send OTP for verification.
- `POST /user/verify/otp` - Verify OTP.
- `POST /user/verify/email` - Request email verification.
- `GET /user/verify/email` - Verify email.

### User Profile Management

- `POST /user/upgradePremium` - Upgrade to premium.
- `POST /user/profile` - Create a user profile.
- `PUT /user/profile` - Update user profile.
- `GET /user/profile` - Get user profile.
- `PUT /user/updateSearchProfile` - Update search preferences.
- `GET /user/searchProfile` - Fetch search profile.
- `POST /user/updateLocation` - Update user location.
- `POST /user/delete` - Delete user account.

### Media Management

- `POST /user/profileMedia` - Upload media.
- `POST /user/update/profileMedia` - Update media.
- `DELETE /user/profileMedia` - Delete media.
- `GET /user/profileMedia` - Fetch user media.

### Matchmaking & Interactions

- `GET /user/matches` - Fetch matches.
- `GET /user/likes` - Fetch likes.
- `GET /searchProfile` - Search for user profiles.
- `POST /user/swipe` - Swipe on profiles.
- `GET /interests` - Fetch available interests.
- `POST /user/interests` - Add user interests.
- `GET /user/interests` - Get user interests.
- `PUT /user/interests` - Update user interests.

### Nudges (Icebreakers)

- `GET /nudges` - Fetch available nudges.
- `POST /user/nudge` - Send a nudge.
- `POST /user/nudge/media` - Send a nudge with media.
- `GET /user/nudges` - Fetch received nudges.
- `PUT /user/nudge` - Update a nudge.
- `DELETE /user/nudge` - Delete a nudge.
- `GET /filters` - Fetch available filters.

### Chat System

- `POST /chat/message` - Send a message.
- `GET /chat/user/chats` - Retrieve chat conversations.
- `GET /chat/user/list` - Fetch chat list.
- `PUT /chat/messages/status` - Update message status.
- `GET /chat/last/messages` - Fetch last messages.

### Stories Feature

- `POST /user/stories/index` - Index user stories.
- `GET /user/stories/search/profileID` - Get stories by user ID.
- `GET /user/stories/search/location` - Get stories by location.

### Events

- `POST /event/index` - Create event index.
- `POST /user/event` - Create an event.
- `GET /user/event` - Fetch user events.
- `PUT /user/event` - Update an event.
- `DELETE /user/event` - Delete an event.
- `GET /events/search` - Search for events.

### Notifications & Device Management

- `POST /user/device/token` - Register device token.
- `POST /user/send/notification` - Send notifications to users.


## Installation & Setup

### Prerequisites

Ensure you have the following installed:
- Golang 1.18+
- Docker & Docker Compose
- MySQL
- Elasticsearch
- Redis
- AWS CLI (for cloud services)

### Steps to Run Locally

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/flicker.git
   cd flicker
   ```
2. Create a `.env` file and configure environment variables.
3. Run database migrations:
   ```sh
   go run main.go migrate up
   ```
4. Start the services using Docker:
   ```sh
   docker-compose up --build
   ```
5. Run the application:
   ```sh
   go run main.go
   ```

## Unit Testing

Flicker includes unit tests to ensure API reliability and correctness. The test suite covers:

- **Authentication & Authorization** (Login, Registration, OTP verification)
- **User Profile Management** (Profile creation, updates, search queries)
- **Matchmaking & Swiping System** (Swiping, matching, user preferences)
- **Chat System** (Message handling, retrieval, and updates)
- **Stories & Events** (Indexing, retrieval, search functionality)

### Running Tests

To execute unit tests, run the following command:

```bash
go test ./...
```

This will run all test cases across different modules.


## License

This project is licensed under the MIT License. See the LICENSE file for details.

---
_Connect, match, and spark new stories with Flicker!_ 🚀
