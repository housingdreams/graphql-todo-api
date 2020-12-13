package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/leminhson2398/todo-api/internal/db"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"golang.org/x/crypto/bcrypt"
)

func (r *commentResolver) UpdatedAt(ctx context.Context, obj *db.Comment) (*time.Time, error) {
	return &obj.UpdatedAt.Time, nil
}

func (r *commentResolver) NumOfLikes(ctx context.Context, obj *db.Comment) (int, error) {
	likes, err := r.Repository.CountNumOfLikesOfCommentByID(ctx, obj.ID)
	return int(likes), err
}

func (r *commentResolver) Liked(ctx context.Context, obj *db.Comment) (bool, error) {
	userID, ok := GetCurrentUserID(ctx)
	if !ok {
		return false, &gqlerror.Error{
			Message: "Error fetching like",
		}
	}
	_, err := r.Repository.SelectLikeByOwnerIDAndCommentID(ctx, db.SelectLikeByOwnerIDAndCommentIDParams{
		CommentID: obj.ID,
		OwnerID:   userID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, err
}

func (r *mutationResolver) CreateComment(ctx context.Context, input CreateCommentInput) (*CreateCommentResponse, error) {
	userID, ok := GetCurrentUserID(ctx)
	if !ok {
		return &CreateCommentResponse{
				Ok:      false,
				Message: "error creating comment",
				Comment: &db.Comment{},
			}, &gqlerror.Error{
				Message: "Error creating comment",
			}
	}

	comment, err := r.Repository.CreateComment(ctx, db.CreateCommentParams{
		OwnerID: userID,
		TodoID:  input.TodoID,
		Content: input.Content,
	})
	if err != nil {
		return &CreateCommentResponse{
			Ok:      false,
			Message: "Error creating comment",
			Comment: &db.Comment{},
		}, err
	}
	return &CreateCommentResponse{
		Ok:      true,
		Message: "Successfully created comment",
		Comment: &comment,
	}, nil
}

func (r *mutationResolver) CreateChildComment(ctx context.Context, input CreateChildCommentInput) (*CreateChildCommentResponse, error) {
	userID, ok := GetCurrentUserID(ctx)
	if !ok {
		return &CreateChildCommentResponse{
				Ok:      false,
				Message: "error creating comment",
				Comment: &db.Comment{},
			}, &gqlerror.Error{
				Message: "Error creating comment",
			}
	}
	comment, err := r.Repository.CreateChildCOmment(ctx, db.CreateChildCOmmentParams{
		OwnerID:         userID,
		TodoID:          input.TodoID,
		Content:         input.Content,
		ParentCommentID: input.ParentID,
	})
	if err != nil {
		return &CreateChildCommentResponse{
			Ok:      false,
			Message: "Error creating comment",
			Comment: &db.Comment{},
		}, err
	}
	return &CreateChildCommentResponse{
		Ok:      true,
		Message: "Successfully created comment",
		Comment: &comment,
	}, nil
}

func (r *mutationResolver) LikeComment(ctx context.Context, input LikeCommentInput) (*LikeCommentResponse, error) {
	userID, ok := GetCurrentUserID(ctx)
	if !ok {
		return &LikeCommentResponse{
				Ok:      false,
				Message: "error like comment",
			}, &gqlerror.Error{
				Message: "Error like comment",
			}
	}

	err := r.Repository.LikeComment(ctx, db.LikeCommentParams{
		OwnerID:   userID,
		CommentID: input.CommentID,
	})
	if err != nil {
		return &LikeCommentResponse{
			Ok:      false,
			Message: "error like comment",
		}, err
	}
	return &LikeCommentResponse{
		Ok:      true,
		Message: "",
	}, nil
}

func (r *mutationResolver) UnlikeComment(ctx context.Context, input UnlikeCommentInput) (*UnLikeCommentResponse, error) {
	userID, ok := GetCurrentUserID(ctx)
	if !ok {
		return &UnLikeCommentResponse{
				Ok:      false,
				Message: "error unlike comment",
			}, &gqlerror.Error{
				Message: "Error unlike comment",
			}
	}
	err := r.Repository.UnlikeComment(ctx, db.UnlikeCommentParams{
		OwnerID:   userID,
		CommentID: input.CommentID,
	})
	if err != nil {
		return &UnLikeCommentResponse{
			Ok:      false,
			Message: "error unlike comment",
		}, err
	}
	return &UnLikeCommentResponse{
		Ok:      true,
		Message: "",
	}, nil
}

func (r *mutationResolver) CreateTodo(ctx context.Context, input CreateTodoInput) (*CreateTodoResponse, error) {
	userID, ok := GetCurrentUserID(ctx)
	if !ok {
		return &CreateTodoResponse{
				Ok:      false,
				Message: "",
			}, &gqlerror.Error{
				Message: "Error creating new todo",
			}
	}

	cloneInput := new(CreateTodoInput)
	*cloneInput = input
	cloneInput.Title = strings.TrimSpace(cloneInput.Title)
	// prevents injecting pure markdown or js into the database
	cloneInput.Content = template.HTMLEscapeString(cloneInput.Content)
	cloneInput.Content = template.JSEscapeString(cloneInput.Content)

	var dueDate time.Time
	if cloneInput.DueDate != nil {
		dueDate = *cloneInput.DueDate
	} else {
		dueDate = time.Now()
	}

	var background string
	if cloneInput.BgColor != nil {
		background = *cloneInput.BgColor
	} else {
		background = "#fff"
	}

	newTodo, err := r.Repository.CreateTodo(ctx, db.CreateTodoParams{
		Title:      cloneInput.Title,
		Content:    cloneInput.Content,
		Duedate:    dueDate.UTC(),
		OwnerID:    userID,
		Completed:  cloneInput.Done,
		Background: background,
	})

	if err != nil {
		return &CreateTodoResponse{
				Ok:      false,
				Message: "",
			},
			&gqlerror.Error{
				Message: err.Error(),
			}
	}
	return &CreateTodoResponse{
		Ok:      true,
		Message: "Created todo successfully",
		Todo:    &newTodo,
	}, nil
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, input UpdateTodoInput) (*UpdateTodoResponse, error) {
	userID, ok := GetCurrentUserID(ctx)
	if !ok {
		return &UpdateTodoResponse{
				Ok:      false,
				Message: "",
			}, &gqlerror.Error{
				Message: "Error updating todo",
			}
	}

	cloneInput := new(UpdateTodoInput)
	*cloneInput = input
	cloneInput.Title = strings.TrimSpace(cloneInput.Title)
	// prevents injecting pure markdown or js into the database
	cloneInput.Content = template.HTMLEscapeString(cloneInput.Content)
	cloneInput.Content = template.JSEscapeString(cloneInput.Content)

	var bgColor string
	if cloneInput.BgColor != nil {
		bgColor = *cloneInput.BgColor
	}

	var dueDate time.Time
	if cloneInput.DueDate == nil {
		dueDate = time.Now()
	} else {
		dueDate = *cloneInput.DueDate
	}

	todo, err := r.Repository.UpdateTodoByID(ctx, db.UpdateTodoByIDParams{
		ID:         cloneInput.ID,
		Title:      cloneInput.Title,
		Content:    cloneInput.Content,
		Duedate:    dueDate,
		Completed:  cloneInput.Done,
		OwnerID:    userID,
		Background: bgColor,
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		return &UpdateTodoResponse{
			Ok:      false,
			Message: "Error updating todo",
			Todo:    &db.Todo{},
		}, err
	}
	return &UpdateTodoResponse{
		Ok:      true,
		Message: "Successfully updated todo",
		Todo:    &todo,
	}, err
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, input DeleteTodoInput) (*DeleteTodoResponse, error) {
	userID, ok := GetCurrentUserID(ctx)
	if !ok {
		return &DeleteTodoResponse{
				Ok:      false,
				Message: "Error deleting todo",
			}, &gqlerror.Error{
				Message: "Error deleting todo",
			}
	}

	err := r.Repository.DeleteTodoByID(ctx, db.DeleteTodoByIDParams{
		ID:      input.ID,
		OwnerID: userID,
	})
	if err != nil {
		return &DeleteTodoResponse{
			Ok:      false,
			Message: "Error deleting todo",
		}, err
	}
	return &DeleteTodoResponse{
		Ok:      true,
		Message: "Successfully deleted todo",
	}, err
}

func (r *mutationResolver) LikeTodo(ctx context.Context, input LikeTodoInput) (*LikeTodoResponse, error) {
	userID, ok := GetCurrentUserID(ctx)
	if !ok {
		return &LikeTodoResponse{
				Ok:      false,
				Message: "You are not authenticated to like this todo",
			}, &gqlerror.Error{
				Message: "You are not authenticated to like this todo",
			}
	}

	err := r.Repository.CreateTodoLike(ctx, db.CreateTodoLikeParams{
		OwnerID: userID,
		TodoID:  input.TodoID,
	})
	if err != nil {
		return &LikeTodoResponse{
			Ok:      false,
			Message: "error like todo",
		}, err
	}
	return &LikeTodoResponse{
		Ok:      true,
		Message: "Successfully liked todo",
	}, err
}

func (r *mutationResolver) UnlikeTodo(ctx context.Context, input UnlikeTodoinput) (*UnlikeTodoResponse, error) {
	userID, ok := GetCurrentUserID(ctx)
	if !ok {
		return &UnlikeTodoResponse{
				Ok:      false,
				Message: "Error unlike todo",
			}, &gqlerror.Error{
				Message: "Error unlike todo",
			}
	}

	err := r.Repository.DeleteTodoLike(ctx, db.DeleteTodoLikeParams{
		OwnerID: userID,
		TodoID:  input.TodoID,
	})
	if err != nil {
		return &UnlikeTodoResponse{
			Ok:      false,
			Message: "Error unlike todo",
		}, err
	}
	return &UnlikeTodoResponse{
		Ok:      true,
		Message: "Successfully unlike todo",
	}, err
}

func (r *mutationResolver) CreateNewUser(ctx context.Context, input CreateNewUserInput) (*CreateNewUserResponse, error) {
	cloneInput := new(CreateNewUserInput)
	*cloneInput = input
	cloneInput.Password = strings.TrimSpace(cloneInput.Password)
	cloneInput.Username = strings.TrimSpace(cloneInput.Username)
	cloneInput.PasswordConfirm = strings.TrimSpace(cloneInput.PasswordConfirm)
	cloneInput.Email = strings.TrimSpace(cloneInput.Email)

	// check if an user with this email does exit:
	_, err := r.Repository.GetUserByEmailOrUsername(ctx, db.GetUserByEmailOrUsernameParams{
		Email:    cloneInput.Email,
		Username: cloneInput.Username,
	})
	if err != nil {
		if err == sql.ErrNoRows { // mean we can create new user
			// check if passwords match
			if cloneInput.Password != cloneInput.PasswordConfirm {
				return &CreateNewUserResponse{
					Ok:      false,
					Message: "Passwords don't match",
				}, nil
			}
			createdAt := time.Now().UTC()
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cloneInput.Password), 14)
			if err != nil {
				return &CreateNewUserResponse{
					Ok:      false,
					Message: "Server error generating hashed password",
				}, err
			}

			userAccount, err := r.Repository.CreateUserAccount(ctx, db.CreateUserAccountParams{
				FirstName:    cloneInput.FirstName,
				LastName:     cloneInput.LastName,
				Username:     cloneInput.Username,
				Email:        cloneInput.Email,
				PasswordHash: string(hashedPassword),
				IsOnline:     false,
				CreatedAt:    createdAt,
			})
			return &CreateNewUserResponse{
				Ok:   true,
				User: &userAccount,
			}, err
		}
	}
	message := fmt.Sprintf("User with email %s or username %s is already exist.", cloneInput.Email, cloneInput.Username)
	return &CreateNewUserResponse{
		Ok:      false,
		Message: message,
	}, err
}

func (r *queryResolver) Users(ctx context.Context) ([]db.UserAccount, error) {
	users, err := r.Repository.GetAllUserAccounts(ctx)
	return users, err
}

func (r *queryResolver) FindUser(ctx context.Context, input FindUserInput) (*db.UserAccount, error) {
	account, err := r.Repository.GetUserAccountByID(ctx, input.ID)
	if err == sql.ErrNoRows {
		return &db.UserAccount{}, &gqlerror.Error{
			Message: "User not found",
		}
	}
	return &account, err
}

func (r *queryResolver) TodoMainComments(ctx context.Context, input TodoMainCommentInput) ([]db.Comment, error) {
	comments, err := r.Repository.SelectMainCommentsByTodoId(ctx, input.TodoID)
	return comments, err
}

func (r *queryResolver) TodoSubComments(ctx context.Context, input TodoSubcommentsInput) ([]db.Comment, error) {
	comments, err := r.Repository.SelectSubcommentsByParentCommentId(ctx, input.ParentCommentID)
	return comments, err
}

func (r *queryResolver) Todos(ctx context.Context) ([]db.Todo, error) {
	return r.Repository.GetAllTodos(ctx)
}

func (r *queryResolver) FindTodo(ctx context.Context, input FindTodoInput) (*db.Todo, error) {
	// user can find any todo inside the system
	todo, err := r.Repository.GetTodoByID(ctx, input.ID)
	if err == sql.ErrNoRows {
		return &db.Todo{}, &gqlerror.Error{
			Message: "Todo not found",
		}
	}

	return &todo, err
}

func (r *queryResolver) FindMyTodos(ctx context.Context) ([]db.Todo, error) {
	userID, ok := GetCurrentUserID(ctx)
	if !ok {
		return []db.Todo{}, &gqlerror.Error{
			Message: "Cannot get your id",
		}
	}
	todos, err := r.Repository.SelectAllTodosOfUserByUserID(ctx, userID)
	return todos, err
}

func (r *queryResolver) FindTodoLikers(ctx context.Context, input FindTodoLikersInput) (*FindTodoLikersResponse, error) {
	users, err := r.Repository.SelectAllUsersWhoLikeTodoByTodoID(ctx, input.TodoID)
	if err == sql.ErrNoRows {
		return &FindTodoLikersResponse{
			Ok:     true,
			Likers: []db.UserAccount{},
		}, nil
	}
	return &FindTodoLikersResponse{
		Ok:     true,
		Likers: users,
	}, err
}

func (r *todoResolver) Status(ctx context.Context, obj *db.Todo) (Status, error) {
	now := time.Now().UTC()

	if obj.Completed {
		return StatusDone, nil
	} else {
		if obj.Duedate.UTC().Before(now) || obj.Duedate.UTC().Equal(now) {
			return StatusMissed, nil
		} else {
			return StatusInProgress, nil
		}
	}
}

func (r *todoResolver) UpdatedAt(ctx context.Context, obj *db.Todo) (*time.Time, error) {
	return &obj.UpdatedAt.Time, nil
}

func (r *todoResolver) NumOfLikes(ctx context.Context, obj *db.Todo) (int, error) {
	numOfLikes, err := r.Repository.SelectNumberOfLikesByTodoID(ctx, obj.ID)
	return int(numOfLikes), err
}

func (r *todoResolver) NumOfComments(ctx context.Context, obj *db.Todo) (int, error) {
	numOfComments, err := r.Repository.CountNumberOfCommentsByTodoID(ctx, obj.ID)
	return int(numOfComments), err
}

func (r *todoResolver) Liked(ctx context.Context, obj *db.Todo) (bool, error) {
	userID, ok := GetCurrentUserID(ctx)
	if !ok {
		return false, nil
	}
	_, err := r.Repository.SelectLikeByOwnerIDAndTodoID(ctx, db.SelectLikeByOwnerIDAndTodoIDParams{
		OwnerID: userID,
		TodoID:  obj.ID,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
	}
	return true, nil
}

func (r *userAccountResolver) ID(ctx context.Context, obj *db.UserAccount) (uuid.UUID, error) {
	return obj.UserID, nil
}

func (r *userAccountResolver) Role(ctx context.Context, obj *db.UserAccount) (*db.Role, error) {
	role, err := r.Repository.SelectRoleByCode(ctx, obj.RoleCode)
	if err != nil {
		return &db.Role{}, err
	}
	return &role, err
	// panic("not implemented")
}

func (r *userAccountResolver) Todos(ctx context.Context, obj *db.UserAccount) ([]db.Todo, error) {
	return r.Repository.SelectAllTodosOfUserByUserID(ctx, obj.UserID)
}

// Comment returns CommentResolver implementation.
func (r *Resolver) Comment() CommentResolver { return &commentResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Todo returns TodoResolver implementation.
func (r *Resolver) Todo() TodoResolver { return &todoResolver{r} }

// UserAccount returns UserAccountResolver implementation.
func (r *Resolver) UserAccount() UserAccountResolver { return &userAccountResolver{r} }

type commentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type todoResolver struct{ *Resolver }
type userAccountResolver struct{ *Resolver }
