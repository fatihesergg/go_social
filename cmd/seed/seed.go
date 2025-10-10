package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/fatihesergg/go_social/internal/model"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateUser(amount int) []model.User {
	users := make([]model.User, amount)

	for i := 0; i < amount; i++ {
		pass := gofakeit.Password(true, true, true, true, true, 20)
		hashPass, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		users[i] = model.User{
			ID:       uuid.New(),
			Name:     gofakeit.Name(),
			LastName: gofakeit.LastName(),
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: string(hashPass),
		}

	}
	return users
}

func GeneratePost(amount int, users []model.User) []model.Post {
	posts := make([]model.Post, amount)
	for i := 0; i < amount; i++ {
		user := users[gofakeit.Number(0, len(users)-1)]
		posts[i] = model.Post{
			ID:      uuid.New(),
			Content: gofakeit.Sentence(4),
			UserID:  user.ID,
		}

	}
	return posts

}

func GeneratePostLike(amount int, users []model.User, posts []model.Post) []model.PostLike {
	postLikes := make([]model.PostLike, 0, amount)
	used := make(map[string]bool)

	for len(postLikes) < amount {
		user := users[gofakeit.Number(0, len(users)-1)]
		post := posts[gofakeit.Number(0, len(posts)-1)]
		key := fmt.Sprintf("%s:%s", user.ID, post.ID)

		if used[key] {
			continue
		}
		used[key] = true

		postLikes = append(postLikes, model.PostLike{
			ID:     uuid.New(),
			UserID: user.ID,
			PostID: post.ID,
		})

	}
	return postLikes
}

func GenerateCommentLike(amount int, users []model.User, comments []model.Comment) []model.CommentLike {
	commentLikes := make([]model.CommentLike, 0, amount)
	used := make(map[string]bool)

	for len(commentLikes) < amount {
		user := users[gofakeit.Number(0, len(users)-1)]
		comment := comments[gofakeit.Number(0, len(comments)-1)]
		key := fmt.Sprintf("%s:%s", user.ID, comment.ID)

		if used[key] {
			continue
		}
		used[key] = true

		commentLikes = append(commentLikes, model.CommentLike{
			ID:        uuid.New(),
			UserID:    user.ID,
			CommentID: comment.ID,
		})

	}
	return commentLikes
}

func GenerateComments(amount int, users []model.User, posts []model.Post) []model.Comment {
	comments := make([]model.Comment, amount)
	for i := 0; i < amount; i++ {
		user := users[gofakeit.Number(0, len(users)-1)]
		post := posts[gofakeit.Number(0, len(posts)-1)]
		comments[i] = model.Comment{
			ID:      uuid.New(),
			Content: gofakeit.Sentence(4),
			UserID:  user.ID,
			PostID:  post.ID,
		}

	}
	return comments
}

func GenerateReplies(amount int, users []model.User, comments []model.Comment) []model.Reply {
	replies := make([]model.Reply, amount)
	for i := 0; i < amount; i++ {
		user := users[gofakeit.Number(0, len(users)-1)]
		comment := comments[gofakeit.Number(0, len(comments)-1)]
		replies[i] = model.Reply{
			ID:        uuid.New(),
			Message:   gofakeit.Sentence(4),
			CommentID: comment.ID,
			UserID:    user.ID,
		}
	}
	return replies
}

func GenerateFollows(amount int, users []model.User) []model.Follow {
	follows := make([]model.Follow, 0, amount)
	used := make(map[string]bool)

	for len(follows) < amount {
		user1 := users[gofakeit.Number(0, len(users)-1)]
		user2 := users[gofakeit.Number(0, len(users)-1)]
		key := fmt.Sprintf("%s:%s", user1.ID, user2.ID)

		if used[key] || user1.ID == user2.ID {
			continue
		}
		used[key] = true
		follows = append(follows, model.Follow{
			ID:       uuid.New(),
			UserID:   user1.ID,
			FollowID: user2.ID,
		})
	}
	return follows
}

func main() {
	gofakeit.Seed(0)

	amount := 100

	users := GenerateUser(amount)
	follows := GenerateFollows(amount, users)
	posts := GeneratePost(amount, users)
	postLikes := GeneratePostLike(amount, users, posts)
	comments := GenerateComments(amount, users, posts)
	commentLikes := GenerateCommentLike(amount, users, comments)
	replies := GenerateReplies(amount, users, comments)

	var sb strings.Builder
	sb.WriteString("-- Seed data using gofakeit/v7\n")
	sb.WriteString("BEGIN;\n\n")

	for _, user := range users {
		sb.WriteString(
			fmt.Sprintf("INSERT INTO users (id,name,last_name,username,email,password) VALUES  ('%s','%s','%s','%s','%s','%s');\n",
				user.ID,
				user.Name,
				user.LastName,
				user.Username,
				user.Email,
				user.Password))
	}
	sb.WriteString("\n")
	for _, follow := range follows {
		sb.WriteString(
			fmt.Sprintf("INSERT INTO follows (id,user_id,follow_id) VALUES ('%s','%s','%s');\n",
				follow.ID,
				follow.UserID,
				follow.FollowID))
	}
	sb.WriteString("\n")

	for _, post := range posts {
		sb.WriteString(
			fmt.Sprintf("INSERT INTO posts (id,user_id,content) VALUES ('%s','%s','%s');\n",
				post.ID,
				post.UserID,
				post.Content))
	}
	sb.WriteString("\n")

	for _, postLike := range postLikes {
		sb.WriteString(
			fmt.Sprintf("INSERT INTO post_likes (id,user_id,post_id) VALUES ('%s','%s','%s');\n",
				postLike.ID,
				postLike.UserID,
				postLike.PostID))
	}
	sb.WriteString("\n")

	for _, comment := range comments {
		sb.WriteString(
			fmt.Sprintf("INSERT INTO comments (id,user_id,post_id,content) VALUES ('%s','%s','%s','%s');\n",
				comment.ID,
				comment.UserID,
				comment.PostID,
				comment.Content))
	}

	sb.WriteString("\n")
	for _, commentLike := range commentLikes {
		sb.WriteString(
			fmt.Sprintf("INSERT INTO comment_likes (id,user_id,comment_id) VALUES ('%s','%s','%s');\n",
				commentLike.ID,
				commentLike.UserID,
				commentLike.CommentID))
	}
	sb.WriteString("\n")

	for _, reply := range replies {
		sb.WriteString(
			fmt.Sprintf("INSERT INTO replies (id,user_id,comment_id,message) VALUES ('%s','%s','%s','%s');\n",
				reply.ID,
				reply.UserID,
				reply.CommentID,
				reply.Message))
	}
	sb.WriteString("\n")

	sb.WriteString("COMMIT;\n")

	if err := os.WriteFile("seed.sql", []byte(sb.String()), 0644); err != nil {
		panic(err)
	}

	fmt.Println("Seed generated.Writed to internal/database/seed.sql")

}
