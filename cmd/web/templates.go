package main

import "snippetbox.alexedwards.net/internal/models"

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
type templateData struct {
    Post            models.Post
    Comments        []models.Comment  // slice kuna mitu sID 240928_214152
    Categories      []models.Category // slice kuna mitu
    Posts           []models.Post     // slice kuna mitu
    CategoryID      int
    CategoryName    string            // Lisa see väli kategooria nime jaoks
    CommentID       int
    Form            userSignupForm    // Changed from interface{} to userSignupForm
    IsAuthenticated bool              // Add an IsAuthenticated field to the templateData struct.
    Reactions       []models.Reaction
    Likes           int               // Post likes
    Dislikes        int               // Post dislikes
    Title           string
}

