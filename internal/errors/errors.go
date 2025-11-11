package errors

import "net/http"

var InternalServerError = AppError{Message: "Internal server error", Code: http.StatusInternalServerError}

// Not found
var ReplyNotFoundError = AppError{Message: "Reply not found", Code: http.StatusNotFound}
var PostNotFoundError = AppError{Message: "Post not found", Code: http.StatusNotFound}
var CommentNotFoundError = AppError{Message: "Comment not found", Code: http.StatusNotFound}
var UserNotFoundError = AppError{Message: "User not found", Code: http.StatusNotFound}

// Empty slice
var NoFollowersFoundError = AppError{Message: "No followers found", Code: http.StatusOK}
var NoFollowingsFoundError = AppError{Message: "No followings found", Code: http.StatusOK}
var NoUsersFoundError = AppError{Message: "No users found", Code: http.StatusOK}
var NoRepliesFoundError = AppError{Message: "No replies found", Code: http.StatusOK}
var NoPostsFoundError = AppError{Message: "No posts found", Code: http.StatusOK}
var NoCommentsFoundError = AppError{Message: "No comments found", Code: http.StatusOK}

// Exist
var AlreadyReplyLikeError = AppError{Message: "Reply already liked", Code: http.StatusBadRequest}
var AlreadyFollowingError = AppError{Message: "You are already following this user", Code: http.StatusBadRequest}
var AlreadyPostLikeError = AppError{Message: "Post already liked", Code: http.StatusBadRequest}
var AlreadyCommentLikeError = AppError{Message: "Comment already liked", Code: http.StatusBadRequest}
var EmailExistError = AppError{Message: "Email already exist", Code: http.StatusBadRequest}
var UsernameExistError = AppError{Message: "Username already exists", Code: http.StatusBadRequest}

// Not exists
var PostNotLikedError = AppError{Message: "Post not liked yet", Code: http.StatusBadRequest}
var CommentNotLikedError = AppError{Message: "Comment not liked yet", Code: http.StatusBadRequest}
var ReplyNotLikedError = AppError{Message: "Reply not liked yet", Code: http.StatusBadRequest}
var NotFollowingError = AppError{Message: "You are not following this user", Code: http.StatusForbidden}

// Other
var InvalidCredentialsError = AppError{Message: "Invalid credentials", Code: http.StatusBadRequest}
var UnauthorizedError = AppError{Message: "You are not authorized to perform this action", Code: http.StatusUnauthorized}
var IDRequiredError = AppError{Message: "ID is required", Code: http.StatusBadRequest}
var InvalidIDFormatError = AppError{Message: "Invalid ID format", Code: http.StatusBadRequest}
var InvalidPermissionError = AppError{Message: "You don't have enough permission to do this operation", Code: http.StatusForbidden}
var InvalidTag = AppError{Message: "Tag query must only contain word and/or digit", Code: http.StatusBadRequest}
