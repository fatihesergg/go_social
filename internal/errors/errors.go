package errors

import "net/http"

var InternalServerError = AppError{Message: "internal server error", Code: http.StatusInternalServerError}

// Not found
var ReplyNotFoundError = AppError{Message: "reply not found", Code: http.StatusNotFound}
var PostNotFoundError = AppError{Message: "post not found", Code: http.StatusNotFound}
var CommentNotFoundError = AppError{Message: "comment not found", Code: http.StatusNotFound}
var UserNotFoundError = AppError{Message: "user not found", Code: http.StatusNotFound}

// Empty slice
var NoFollowersFoundError = AppError{Message: "no followers found", Code: http.StatusOK}
var NoFollowingsFoundError = AppError{Message: "no followings found", Code: http.StatusOK}
var NoUsersFoundError = AppError{Message: "no users found", Code: http.StatusOK}
var NoRepliesFoundError = AppError{Message: "no replies found", Code: http.StatusOK}
var NoPostsFoundError = AppError{Message: "no posts found", Code: http.StatusOK}
var NoCommentsFoundError = AppError{Message: "no comments found", Code: http.StatusOK}

// Exist
var AlreadyReplyLikeError = AppError{Message: "reply already liked", Code: http.StatusBadRequest}
var AlreadyFollowingError = AppError{Message: "you are already following this user", Code: http.StatusBadRequest}
var AlreadyPostLikeError = AppError{Message: "post already liked", Code: http.StatusBadRequest}
var AlreadyCommentLikeError = AppError{Message: "comment already liked", Code: http.StatusBadRequest}
var EmailExistError = AppError{Message: "email already exist", Code: http.StatusBadRequest}
var UsernameExistError = AppError{Message: "username already exists", Code: http.StatusBadRequest}

// Not exists
var PostNotLikedError = AppError{Message: "post not liked yet", Code: http.StatusBadRequest}
var CommentNotLikedError = AppError{Message: "comment not liked yet", Code: http.StatusBadRequest}
var ReplyNotLikedError = AppError{Message: "reply not liked yet", Code: http.StatusBadRequest}
var NotFollowingError = AppError{Message: "you are not following this user", Code: http.StatusForbidden}

// Other
var InvalidCredentialsError = AppError{Message: "invalid credentials", Code: http.StatusBadRequest}
var UnauthorizedError = AppError{Message: "you are not authorized to perform this action", Code: http.StatusUnauthorized}
var IDRequiredError = AppError{Message: "id is required", Code: http.StatusBadRequest}
var InvalidIDFormatError = AppError{Message: "invalid ID format", Code: http.StatusBadRequest}
var InvalidPermissionError = AppError{Message: "you don't have enough permission to do this operation", Code: http.StatusForbidden}
var InvalidTag = AppError{Message: "tag query must only contain word and/or digit", Code: http.StatusBadRequest}
