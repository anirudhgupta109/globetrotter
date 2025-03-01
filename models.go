package main

import (
	"github.com/google/uuid"
	"time"
)

// User model
type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Don't expose password in JSON responses
	AuthToken string    `json:"auth_token,omitempty"`
}

// Destination model
type Destination struct {
	ID      uuid.UUID `json:"id"`
	City    string    `json:"city"`
	Country string    `json:"country"`
}

// Clue model
type Clue struct {
	ID           uuid.UUID `json:"id"`
	DestinationID uuid.UUID `json:"destination_id"`
	ClueText     string    `json:"clue_text"`
}

// FunFact model
type FunFact struct {
	ID           uuid.UUID `json:"id"`
	DestinationID uuid.UUID `json:"destination_id"`
	FactText     string    `json:"fact_text"`
}

// Trivia model
type Trivia struct {
	ID           uuid.UUID `json:"id"`
	DestinationID uuid.UUID `json:"destination_id"`
	TriviaText   string    `json:"trivia_text"`
}

// Challenge model
type Challenge struct {
	ID              uuid.UUID   `json:"id"`
	Inviter         string      `json:"inviter"`
	Score           int         `json:"score"`
	CorrectAnswers  int         `json:"correct_answers"`
	IncorrectAnswers int        `json:"incorrect_answers"`
	CluesRevealed   int         `json:"clues_revealed"`
	IsActive        bool        `json:"is_active"`
	EndedAt         *time.Time  `json:"ended_at,omitempty"`
	QuestionIDs     []uuid.UUID `json:"question_ids"`
	CreatedAt       time.Time   `json:"created_at"`
}

// Question model
type Question struct {
	ID           uuid.UUID `json:"id"`
	DestinationID uuid.UUID `json:"destination_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// API Request and Response structures
type UserRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Username  string `json:"username"`
	AuthToken string `json:"auth_token"`
}

type AnswerRequest struct {
	QuestionId  string    `json:"question_id" binding:"required"`
	Username    string    `json:"username" binding:"required"`
	City        string    `json:"city" binding:"required"`
	ChallengeId uuid.UUID `json:"challenge_id"`
}

type CreateChallengeRequest struct {
	Username string `json:"username" binding:"required"`
}

type CreateChallengeResponse struct {
	ChallengeId uuid.UUID `json:"challenge_id"`
	Inviter     string    `json:"inviter"`
}

type EndChallengeRequest struct {
	ChallengeId      string      `json:"challenge_id" binding:"required"`
	Username         string      `json:"username" binding:"required"`
	Score            int         `json:"score"`
	CorrectAnswers   int         `json:"correct_answers"`
	IncorrectAnswers int         `json:"incorrect_answers"`
	CluesRevealed    int         `json:"clues_revealed"`
	QuestionIDs      []uuid.UUID `json:"question_ids"`
}

type RevealClueRequest struct {
	Username    string    `json:"username" binding:"required"`
	ChallengeId uuid.UUID `json:"challenge_id"`
}

type GetQuestionResponse struct {
	QuestionId uuid.UUID `json:"question_id"`
	Clues      []string  `json:"clues"`
	Choices    []string  `json:"choices"`
	Trivia     string    `json:"trivia"`
}

type SubmitAnswerResponse struct {
	Correct bool   `json:"correct"`
	FunFact string `json:"fun_fact"`
}