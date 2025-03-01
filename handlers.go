package main

import (
	"context"
	"encoding/base64"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

// Register a new user
func registerUser(c *gin.Context) {
	var req UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	var exists bool
	err := db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", req.Username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	// Generate auth token
	authToken := generateToken(32)

	// Insert user
	_, err = db.Exec(context.Background(),
		"INSERT INTO users (username, password, auth_token) VALUES ($1, $2, $3)",
		req.Username, string(hashedPassword), authToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "User registered successfully",
		"username":   req.Username,
		"auth_token": authToken,
	})
}

// Login user
func loginUser(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from database
	var user User
	var hashedPassword string
	err := db.QueryRow(context.Background(),
		"SELECT id, username, password, auth_token FROM users WHERE username = $1",
		req.Username).Scan(&user.ID, &user.Username, &hashedPassword, &user.AuthToken)
	if err != nil {
		if err == pgx.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate new auth token
	authToken := generateToken(32)
	_, err = db.Exec(context.Background(), "UPDATE users SET auth_token = $1 WHERE username = $2", authToken, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating auth token"})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Username:  user.Username,
		AuthToken: authToken,
	})
}

// Get a random question
func getRandomQuestion(c *gin.Context) {
	challengeIDStr := c.Query("challenge_id")
	username := c.Query("username")
	
	var questionID uuid.UUID
	var destinationID uuid.UUID
	var err error
	
	// Check if this is part of a challenge
	if challengeIDStr != "" {
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required for challenge mode"})
			return
		}
		
		challengeID, err := uuid.Parse(challengeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid challenge ID"})
			return
		}
		
		// Get challenge information
		var challenge Challenge
		err = db.QueryRow(context.Background(), 
			"SELECT id, inviter, is_active, question_ids FROM challenges WHERE id = $1", 
			challengeID).Scan(&challenge.ID, &challenge.Inviter, &challenge.IsActive, &challenge.QuestionIDs)
		
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
			} else {
				log.Printf("Database error retrieving challenge %s: %v", challengeIDStr, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving challenge"})
			}
			return
		}
		
		// Determine if this is player 1 (inviter) or player 2 (challenger)
		isInviter := username == challenge.Inviter
		
		// If there are no questions yet and this is the inviter, generate some
		if len(challenge.QuestionIDs) == 0 && isInviter {
			log.Printf("Player 1 (inviter) generating questions for challenge %s", challengeIDStr)
			
			// Get 5 random destinations for this challenge
			rows, err := db.Query(context.Background(), 
				"SELECT id FROM destinations ORDER BY RANDOM() LIMIT 5")
			if err != nil {
				log.Printf("Error retrieving random destinations: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving destinations"})
				return
			}
			defer rows.Close()
			
			var destIDs []uuid.UUID
			for rows.Next() {
				var id uuid.UUID
				if err := rows.Scan(&id); err != nil {
					log.Printf("Error scanning destination ID: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning destination ID"})
					return
				}
				destIDs = append(destIDs, id)
			}
			
			// Create questions for each destination
			for _, destID := range destIDs {
				var qID uuid.UUID
				err = db.QueryRow(context.Background(),
					"INSERT INTO questions (destination_id) VALUES ($1) RETURNING id",
					destID).Scan(&qID)
				if err != nil {
					log.Printf("Error creating question for destination %s: %v", destID, err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating question"})
					return
				}
				challenge.QuestionIDs = append(challenge.QuestionIDs, qID)
			}
			
			// Update challenge with question IDs
			_, err = db.Exec(context.Background(),
				"UPDATE challenges SET question_ids = $1 WHERE id = $2",
				challenge.QuestionIDs, challengeID)
			if err != nil {
				log.Printf("Error updating challenge %s with question IDs: %v", challengeIDStr, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating challenge"})
				return
			}
		} else if len(challenge.QuestionIDs) == 0 && !isInviter {
			// Player 2 has no more questions and can't generate new ones
			c.JSON(http.StatusOK, gin.H{"message": "Challenge complete! No more questions available."})
			return
		}
		
		// Make sure we have questions to return
		if len(challenge.QuestionIDs) == 0 {
			log.Printf("No questions available for challenge %s", challengeIDStr)
			c.JSON(http.StatusOK, gin.H{"message": "No questions available for this challenge"})
			return
		}
		
		// Pick the first question that hasn't been answered yet
		questionID = challenge.QuestionIDs[0]
		
		// Get the destination ID for this question
		err = db.QueryRow(context.Background(),
			"SELECT destination_id FROM questions WHERE id = $1",
			questionID).Scan(&destinationID)
		if err != nil {
			log.Printf("Error retrieving destination for question %s: %v", questionID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving question"})
			return
		}
	} else {
		// This is a regular game, not a challenge
		// Get a random destination
		err = db.QueryRow(context.Background(),
			"SELECT id FROM destinations ORDER BY RANDOM() LIMIT 1").Scan(&destinationID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving random destination"})
			return
		}
		
		// Create a new question for this destination
		err = db.QueryRow(context.Background(),
			"INSERT INTO questions (destination_id) VALUES ($1) RETURNING id",
			destinationID).Scan(&questionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating question"})
			return
		}
	}
	
	// Get clues for the destination
	rows, err := db.Query(context.Background(),
		"SELECT clue_text FROM clues WHERE destination_id = $1 ORDER BY RANDOM() LIMIT 2",
		destinationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving clues"})
		return
	}
	defer rows.Close()
	
	var clues []string
	for rows.Next() {
		var clue string
		if err := rows.Scan(&clue); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning clue"})
			return
		}
		clues = append(clues, clue)
	}
	
	// Get a random trivia for this destination
	var trivia string
	err = db.QueryRow(context.Background(),
		"SELECT trivia_text FROM trivia WHERE destination_id = $1 ORDER BY RANDOM() LIMIT 1",
		destinationID).Scan(&trivia)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving trivia"})
		return
	}
	
	// Get the correct destination name
	var correctCity string
	err = db.QueryRow(context.Background(),
		"SELECT city FROM destinations WHERE id = $1",
		destinationID).Scan(&correctCity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving destination"})
		return
	}
	
	// Get 5 random incorrect destinations
	rows, err = db.Query(context.Background(),
		"SELECT city FROM destinations WHERE id != $1 ORDER BY RANDOM() LIMIT 5",
		destinationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving incorrect destinations"})
		return
	}
	defer rows.Close()
	
	choices := []string{correctCity}
	for rows.Next() {
		var city string
		if err := rows.Scan(&city); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning city"})
			return
		}
		choices = append(choices, city)
	}
	
	// Shuffle the choices
	rand.Shuffle(len(choices), func(i, j int) {
		choices[i], choices[j] = choices[j], choices[i]
	})
	
	c.JSON(http.StatusOK, GetQuestionResponse{
		QuestionId: questionID,
		Clues:      clues,
		Choices:    choices,
		Trivia:     trivia,
	})
}

// Submit an answer
func submitAnswer(c *gin.Context) {
	var req AnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	questionID, err := uuid.Parse(req.QuestionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID"})
		return
	}
	
	// Get the destination ID for this question
	var destinationID uuid.UUID
	err = db.QueryRow(context.Background(),
		"SELECT destination_id FROM questions WHERE id = $1",
		questionID).Scan(&destinationID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Question %s not found", req.QuestionId)
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		} else {
			log.Printf("Error retrieving question %s: %v", req.QuestionId, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving question"})
		}
		return
	}
	
	// Get the correct city for this destination
	var correctCity string
	err = db.QueryRow(context.Background(),
		"SELECT city FROM destinations WHERE id = $1",
		destinationID).Scan(&correctCity)
	if err != nil {
		log.Printf("Error retrieving destination %s: %v", destinationID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving destination"})
		return
	}
	
	// Check if the answer is correct
	isCorrect := strings.EqualFold(req.City, correctCity)
	
	// Get a random fun fact for this destination
	var funFact string
	err = db.QueryRow(context.Background(),
		"SELECT fact_text FROM fun_facts WHERE destination_id = $1 ORDER BY RANDOM() LIMIT 1",
		destinationID).Scan(&funFact)
	if err != nil {
		log.Printf("Error retrieving fun fact for destination %s: %v", destinationID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving fun fact"})
		return
	}
	
	// If this is part of a challenge, update the challenge stats
	if req.ChallengeId != uuid.Nil {
		// Update score and correct/incorrect counts
		tx, err := db.Begin(context.Background())
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error starting transaction"})
			return
		}
		defer tx.Rollback(context.Background())
		
		// Get the current challenge data
		var challenge Challenge
		err = tx.QueryRow(context.Background(),
			"SELECT score, correct_answers, incorrect_answers, question_ids FROM challenges WHERE id = $1",
			req.ChallengeId).Scan(&challenge.Score, &challenge.CorrectAnswers, &challenge.IncorrectAnswers, &challenge.QuestionIDs)
		if err != nil {
			if err == pgx.ErrNoRows {
				log.Printf("Challenge %s not found", req.ChallengeId)
				c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
			} else {
				log.Printf("Error retrieving challenge %s: %v", req.ChallengeId, err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving challenge"})
			}
			return
		}
		
		// Update the score
		if isCorrect {
			challenge.Score += 3
			challenge.CorrectAnswers++
		} else {
			challenge.IncorrectAnswers++
		}
		
		// Remove the answered question from the list
		var updatedQuestionIDs []uuid.UUID
		for _, id := range challenge.QuestionIDs {
			if id != questionID {
				updatedQuestionIDs = append(updatedQuestionIDs, id)
			}
		}
		
		// Update the challenge
		_, err = tx.Exec(context.Background(),
			"UPDATE challenges SET score = $1, correct_answers = $2, incorrect_answers = $3, question_ids = $4 WHERE id = $5",
			challenge.Score, challenge.CorrectAnswers, challenge.IncorrectAnswers, updatedQuestionIDs, req.ChallengeId)
		if err != nil {
			log.Printf("Error updating challenge %s: %v", req.ChallengeId, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating challenge"})
			return
		}
		
		err = tx.Commit(context.Background())
		if err != nil {
			log.Printf("Error committing transaction: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error committing transaction"})
			return
		}
	}
	
	c.JSON(http.StatusOK, SubmitAnswerResponse{
		Correct: isCorrect,
		FunFact: funFact,
	})
}

// Reveal a clue
func revealClue(c *gin.Context) {
	var req RevealClueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// If this is part of a challenge, update the challenge stats
	if req.ChallengeId != uuid.Nil {
		// Update clues revealed count and subtract a point
		_, err := db.Exec(context.Background(),
			"UPDATE challenges SET clues_revealed = clues_revealed + 1, score = GREATEST(0, score - 1) WHERE id = $1",
			req.ChallengeId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating challenge"})
			return
		}
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Clue revealed"})
}

// Create a challenge
func createChallenge(c *gin.Context) {
	var req CreateChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Verify user exists
	var exists bool
	err := db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", req.Username).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking user existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	
	if !exists {
		// Create the user if they don't exist (for easier testing)
		log.Printf("User %s not found, creating...", req.Username)
		// Generate random password since we just need the user to exist
		password := generateToken(10)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}
		
		_, err = db.Exec(context.Background(),
			"INSERT INTO users (username, password) VALUES ($1, $2)",
			req.Username, string(hashedPassword))
		if err != nil {
			log.Printf("Error creating user %s: %v", req.Username, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}
	}
	
	// Create a new challenge with empty question_ids array
	var challengeID uuid.UUID
	err = db.QueryRow(context.Background(),
		"INSERT INTO challenges (inviter, is_active, question_ids) VALUES ($1, true, '{}') RETURNING id",
		req.Username).Scan(&challengeID)
	if err != nil {
		log.Printf("Error creating challenge for user %s: %v", req.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating challenge"})
		return
	}
	
	log.Printf("Challenge %s created for user %s", challengeID, req.Username)
	c.JSON(http.StatusCreated, CreateChallengeResponse{
		ChallengeId: challengeID,
		Inviter:     req.Username,
	})
}

// Get challenge details
func getChallenge(c *gin.Context) {
	challengeIDStr := c.Param("id")
	challengeID, err := uuid.Parse(challengeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid challenge ID"})
		return
	}
	
	// Get challenge details
	var challenge Challenge
	err = db.QueryRow(context.Background(),
		"SELECT id, inviter, score, correct_answers, incorrect_answers, clues_revealed, is_active, question_ids, created_at FROM challenges WHERE id = $1",
		challengeID).Scan(&challenge.ID, &challenge.Inviter, &challenge.Score, &challenge.CorrectAnswers, &challenge.IncorrectAnswers, &challenge.CluesRevealed, &challenge.IsActive, &challenge.QuestionIDs, &challenge.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("Challenge %s not found", challengeIDStr)
			c.JSON(http.StatusNotFound, gin.H{"error": "Challenge not found"})
			return
		}
		log.Printf("Error retrieving challenge %s: %v", challengeIDStr, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving challenge"})
		return
	}
	
	// Ensure challenge is active for joining
	if !challenge.IsActive {
		// Re-activate if needed
		_, err := db.Exec(context.Background(),
			"UPDATE challenges SET is_active = true WHERE id = $1",
			challengeID)
		if err != nil {
			log.Printf("Error reactivating challenge %s: %v", challengeIDStr, err)
			// Don't return an error, just log it
		} else {
			challenge.IsActive = true
			log.Printf("Challenge %s reactivated", challengeIDStr)
		}
	}
	
	c.JSON(http.StatusOK, challenge)
}

// End a challenge
func endChallenge(c *gin.Context) {
	var req EndChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	challengeID, err := uuid.Parse(req.ChallengeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid challenge ID"})
		return
	}
	
	// Update challenge status and stats
	now := time.Now()
	_, err = db.Exec(context.Background(),
		"UPDATE challenges SET is_active = false, ended_at = $1, score = $2, correct_answers = $3, incorrect_answers = $4, clues_revealed = $5, question_ids = $6 WHERE id = $7",
		now, req.Score, req.CorrectAnswers, req.IncorrectAnswers, req.CluesRevealed, req.QuestionIDs, challengeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating challenge"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Challenge ended successfully"})
}

// Helper function to generate random token
func generateToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Utility function for debugging
func logUUIDArray(arr []uuid.UUID) string {
	if len(arr) == 0 {
		return "[]"
	}
	
	result := "["
	for i, id := range arr {
		if i > 0 {
			result += ", "
		}
		result += id.String()
	}
	result += "]"
	return result
}