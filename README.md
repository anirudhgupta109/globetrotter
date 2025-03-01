# Globetrotter

Globetrotter is a full-stack web application where users get cryptic clues about famous places and must guess which destination they refer to. After guessing, users unlock fun facts and trivia about the destination.

## Features

- **Random Destination Clues**: Users are presented with one main clue and an optional hidden clue for a random destination
- **Multiple Choice Answers**: Select from 6 possible destinations
- **Score Tracking**: Earn 3 points for correct answers, lose 1 point for revealing hidden clues
- **Challenge Mode**: Challenge friends to beat your score with the same set of questions
- **Dynamic Feedback**: Confetti animations for correct answers and fun facts after each guess
- **Rich Content**: Discover fun facts and trivia about each destination upon answering



## Tech Stack

- **Frontend**: React.js
- **Backend**: Go with Gin framework
- **Database**: PostgreSQL

## Setup

### Prerequisites

- Node.js (v14+)
- npm or yarn
- Go 1.19+
- PostgreSQL

### Backend Installation
1. Clone the repository
```bash
git clone https://github.com/anirudhgupta109/globetrotter.git -b backend
cd globetrotter
```

2. Set up environment variables
```bash
cp .env.example .env
# Edit .env with your database credentials
```

3. Create the database
```bash
createdb globetrotter
```

4. Install dependencies
```bash
go mod tidy
```

5. Run the application
```bash
go run .
```

The server will start on port 8080 (or as configured in your .env file).

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:postgres@localhost:5432/globetrotter` |
| `PORT` | Server port | `8080` |
| `PENDING_MIGRATION` | If the initial migration is pending | `true` |

## Database Schema

The application uses PostgreSQL with the following tables:
- `users` - User accounts
- `destinations` - Travel destinations
- `clues` - Clues for destinations
- `fun_facts` - Fun facts about destinations
- `trivia` - Trivia about destinations
- `questions` - Generated questions
- `challenges` - Challenge invitations and stats



### Frontend Installation

1. Clone the repository
```bash
git clone https://github.com/anirudhgupta109/globetrotter.git -b frontend
cd globetrotter
```

2. Install dependencies
```bash
npm install
```

3. Start the development server
```bash
npm start
```

The frontend will run on http://localhost:3000 and connect to the backend on port 8080.

### Configuration Notes

- The frontend expects the backend to be running on `http://localhost:8080/api`
- CORS is already set up on the backend to allow requests from port 3000
- Authentication tokens are stored in localStorage and automatically included in API requests

## User Flow

1. **Registration/Login**: Create an account or log in with existing credentials
2. **Main Menu**: Choose between regular gameplay or challenge mode
3. **Gameplay**: 
   - View a clue about a destination
   - Optionally reveal a second clue (costs 1 point)
   - Select from 6 possible destinations
   - Receive immediate feedback and fun facts
   - Continue to the next random destination

4. **Challenge Mode**:
   - Create a challenge for friends
   - Share link via WhatsApp or copy to clipboard
   - Friends can play the same set of questions
   - Compare scores at the end

## API Endpoints

### User Authentication
- `POST /api/users/register` - Register a new user
- `POST /api/users/login` - Log in and receive auth token

### Game Functionality
- `GET /api/game/question` - Get a random question
- `POST /api/game/answer` - Submit an answer
- `POST /api/game/reveal-clue` - Reveal the hidden clue

### Challenge System
- `POST /api/challenges/create` - Create a new challenge
- `GET /api/challenges/:id` - Get challenge details
- `POST /api/challenges/end` - End the challenge and record stats

### Game Logic

#### Get Random Question
- **Endpoint:** `GET /api/question?challenge_id=optional_uuid&username=optional_string`
- **Response:**
```json
{
  "question_id": "uuid",
  "clues": ["string", "string"],
  "choices": ["string", "string", "string", "string", "string", "string"],
  "trivia": "string"
}
```
- **Notes:** 
  - When used with challenge_id, username is required
  - Only the inviter (player 1) can generate new questions
  - Player 2 cannot generate new questions and will receive a "Challenge complete" message when they've answered all questions

#### Submit Answer
- **Endpoint:** `POST /api/answer`
- **Body:**
```json
{
  "question_id": "uuid",
  "username": "string",
  "city": "string",
  "challenge_id": "optional_uuid"
}
```
- **Response:**
```json
{
  "correct": true|false,
  "fun_fact": "string"
}
```

#### Reveal Clue
- **Endpoint:** `POST /api/reveal-clue`
- **Body:**
```json
{
  "username": "string",
  "challenge_id": "optional_uuid"
}
```
- **Response:** 200 OK

### Challenge System

#### Create Challenge
- **Endpoint:** `POST /api/challenges/create`
- **Body:**
```json
{
  "username": "string"
}
```
- **Response:**
```json
{
  "challenge_id": "uuid",
  "inviter": "string"
}
```

#### Get Challenge Details
- **Endpoint:** `GET /api/challenges/:id`
- **Response:**
```json
{
  "id": "uuid",
  "inviter": "string",
  "score": number,
  "correct_answers": number,
  "incorrect_answers": number,
  "clues_revealed": number,
  "is_active": boolean,
  "question_ids": ["uuid"],
  "created_at": "timestamp"
}
```

#### End Challenge
- **Endpoint:** `POST /api/challenges/end`
- **Body:**
```json
{
  "challenge_id": "uuid",
  "username": "string",
  "score": number,
  "correct_answers": number,
  "incorrect_answers": number,
  "clues_revealed": number,
  "question_ids": ["uuid"]
}
```
- **Response:** 200 OK

## Challenge Flow

1. User creates a challenge with their username
2. The challenge ID is shared with a friend
3. Friend accesses challenge details via GET /api/challenges/:id
4. Friend gets questions via GET /api/question?challenge_id=:id&username=:username
5. Friend submits answers via POST /api/answer (includes challenge_id)
6. When done, client calls POST /api/challenges/end to record final stats

## Sample Data

The application includes sample data in `schemas/dataset.json`. You can add more destinations by following the format:

```json
{
  "city": "Paris",
  "country": "France",
  "clues": [
    "This city is home to a famous tower that sparkles every night.",
    "Known as the 'City of Love' and a hub for fashion and art."
  ],
  "fun_fact": [
    "The Eiffel Tower was supposed to be dismantled after 20 years!"
  ],
  "trivia": [
    "This city is famous for its croissants and macarons. Bon appétit!"
  ]
}
```

## License

MIT