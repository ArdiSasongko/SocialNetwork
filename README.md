# Social Network Backend API

Welcome to the **Social Network Backend API**, a project built to implement my backend development skills in a real-world application. While this project is not yet perfect, I am committed to continuously improving and adding new features in the future.

## 🚀 Features

- **🔐 Authentication**: Secure user registration, login, and token-based authentication.
- **📝 Content Management**: Create, edit, delete posts, and interact with user-generated content.
- **🤝 Social Features**: Like, comment, follow, and track user activities.
- **📦 File Uploads**: Media uploads made effortless with Cloudinary integration.
- **📈 Scalability**: Fully Dockerized for hassle-free deployment and scaling.

## 🛠️ Tech Stack

- **Programming Language**: Go
- **Web Framework**: Chi
- **Database**: PostgreSQL
- **Cloud Storage**: Cloudinary
- **Deployment**: Docker

## 📖 API Routes

Explore the available endpoints in the API:

### Health Check

- **GET /v1/health**: Verify the API health status.

### Authentication

- **POST /v1/authentication/register**: Register a new user.
- **POST /v1/authentication/login**: Login and obtain a token.

### Profile Management

- **GET /v1/profile/**: Fetch the logged-in user's profile (requires authentication).
- **PATCH /v1/profile/**: Update user profile.
- **PUT /v1/profile/image**: Update user profile image.
- **GET /v1/profile/{postID}**: Get a specific post by the logged-in user (requires post context).

### Post Management

- **POST /v1/posts/**: Create a new post (requires authentication).
- **GET /v1/posts/{postID}/**: Retrieve a post by its ID.
- **PATCH /v1/posts/{postID}/**: Update a post if the user is a moderator.
- **DELETE /v1/posts/{postID}/**: Delete a post if the user is an admin.

### User Management

- **GET /v1/users/{userID}/**: Fetch another user's profile.
- **POST /v1/users/{userID}/follow**: Follow a user.
- **DELETE /v1/users/{userID}/unfollow**: Unfollow a user.

### Feeds

- **GET /v1/feeds/**: Retrieve a feed of posts (requires authentication).
- **GET /v1/feeds/{postID}/**: View a specific post in the feed.
- **POST /v1/feeds/{postID}/comment**: Add a comment to a post.
- **PUT /v1/feeds/{postID}/like**: Like a post.
- **PUT /v1/feeds/{postID}/dislike**: Dislike a post.

## 📚 Full Documentation

For a comprehensive guide to all endpoints and their usage, check out our Postman documentation:
👉 [API Documentation on Postman](https://documenter.getpostman.com/view/29238176/2sAYQZJsuW)

## 🌟 Why This Project?

This project embodies clean coding principles and modern backend development practices, making it an excellent foundation for:

- Aspiring backend developers showcasing their portfolio.
- Building scalable and feature-rich social networking platforms.

---

Feel free to explore, contribute, and build on top of this project. Happy coding! 😄

