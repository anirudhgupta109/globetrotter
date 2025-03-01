import React from 'react';
import { useLocation, useNavigate } from 'react-router-dom';

const ChallengeComplete = () => {
  const location = useLocation();
  const navigate = useNavigate();
  
  // Get stats from location state
  const { 
    score = 0, 
    correct_answers = 0, 
    incorrect_answers = 0, 
    clues_revealed = 0, 
    created_challenge = false,
    inviter_info = null 
  } = location.state || {};
  
  const handlePlayAgain = () => {
    navigate('/menu');
  };
  
  return (
    <div className="game-container challenge-complete">
      <h1 className="form-title">Challenge Complete!</h1>
      
      {inviter_info && (
        <div style={{ marginBottom: '20px' }}>
          <p>You played a challenge from <strong>{inviter_info.inviter}</strong></p>
          <p>Their score: <strong>{inviter_info.score}</strong></p>
          <p style={{ marginTop: '10px', fontWeight: 'bold', fontSize: '18px' }}>
            {score > inviter_info.score 
              ? "🎉 Congratulations! You beat their score!" 
              : score === inviter_info.score 
                ? "It's a tie!" 
                : "You were close! Try again to beat their score!"}
          </p>
        </div>
      )}
      
      <h2 className="challenge-score">Your Score: {score}</h2>
      
      <div className="challenge-stats">
        <p>Correct Answers: {correct_answers}</p>
        <p>Incorrect Answers: {incorrect_answers}</p>
        <p>Clues Revealed: {clues_revealed}</p>
      </div>
      
      {created_challenge && (
        <p style={{ marginBottom: '20px' }}>
          Your challenge has been saved. Share it with friends to see if they can beat your score!
        </p>
      )}
      
      <button className="form-button" onClick={handlePlayAgain}>
        Back to Menu
      </button>
    </div>
  );
};

export default ChallengeComplete;