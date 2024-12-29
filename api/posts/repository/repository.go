package repository

import (
	"context"
	"errors"
	"hornet/api/posts/model"
	"sync"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PostRepository defines the methods for interacting with the database
type PostRepository struct {
	Collection *mongo.Collection
}

// Declare a global variable for the singleton instance of PostRepository
var (
	postRepositoryInstance *PostRepository
	once                   sync.Once
)

// NewPostRepository creates a new PostRepository instance if it doesn't exist
func NewPostRepository(db *mongo.Database) *PostRepository {
	once.Do(func() {
		postRepositoryInstance = &PostRepository{
			Collection: db.Collection("posts"),
		}
	})
	return postRepositoryInstance
}

// FindPostByID retrieves a post by its ID
func (r *PostRepository) FindPostByID(ctx context.Context, id uuid.UUID) (model.Post, error) {
	// Filter for finding the post by its ID
	filter := bson.M{"_id": id}

	var post model.Post
	err := r.Collection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Post{}, errors.New("post not found")
		}
		return model.Post{}, err
	}

	return post, nil
}

// FindPostsByAuthorID retrieves all posts by a given author ID
func (r *PostRepository) FindPostsByAuthorID(ctx context.Context, authorID uuid.UUID) ([]model.Post, error) {
	// Filter for finding posts by author ID
	filter := bson.M{"author_id": authorID}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var posts []model.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	return posts, nil
}

// SavePost saves a new post to the database
func (r *PostRepository) SavePost(ctx context.Context, post model.Post) error {
	_, err := r.Collection.InsertOne(ctx, post, options.InsertOne())
	return err
}

// DeletePost deletes a post by its ID
func (r *PostRepository) DeletePost(ctx context.Context, id uuid.UUID) error {
	// Filter for finding the post by its ID
	filter := bson.M{"_id": id}

	// Attempt to delete the post from the collection
	_, err := r.Collection.DeleteOne(ctx, filter)
	return err
}

// UpdatePost updates a post in the database
func (r *PostRepository) UpdatePost(ctx context.Context, post model.Post) error {
	// Filter for finding the post by its ID
	filter := bson.M{"_id": post.ID}

	// Prepare the update data
	update := bson.M{
		"$set": bson.M{
			"replies_count": post.RepliesCount,
			"shares_count":  post.SharesCount,
		},
	}

	// Attempt to update the post in the collection
	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

// FindPostsByParentID retrieves posts by their parent ID (for replies)
func (r *PostRepository) FindPostsByParentID(ctx context.Context, id uuid.UUID) ([]model.Post, error) {
	// Filter for finding posts where the ParentPostID matches the given parent post ID
	filter := bson.M{"parent_post_id": id}

	// Fetch the replies from the database
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var replies []model.Post
	if err := cursor.All(ctx, &replies); err != nil {
		return nil, err
	}

	return replies, nil
}
