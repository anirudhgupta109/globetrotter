import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import Confetti from 'react-confetti';
import SharePopup from './SharePopup';
import { getQuestion, submitAnswer, revealClue, createChallenge, getChallenge, endChallenge } from '../services/api';

const Challenge = () => {
  const { id: challengeId } = useParams();
  const navigate = useNavigate();
  const username = localStorage.getItem('username');

  const [question, setQuestion] = useState(null);
  const [isSecondClueVisible, setIsSecondClueVisible] = useState(false);
  const [selectedAnswer, setSelectedAnswer] = useState('');
  const [answerResult, setAnswerResult] = useState(null);
  const [showConfetti, setShowConfetti] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [showSharePopup, setShowSharePopup] = useState(false);
  const [inviterInfo, setInviterInfo] = useState(null);
  
  // Challenge stats
  const [score, setScore] = useState(0);
  const [correctAnswers, setCorrectAnswers] = useState(0);
  const [incorrectAnswers, setIncorrectAnswers] = useState(0);
  const [cluesRevealed, setCluesRevealed] = useState(0);
  const [questionIDs, setQuestionIDs] = useState([]);
  const [currentChallengeId, setCurrentChallengeId] = useState(null);
  
  // Initialize a new challenge or get existing one
  useEffect(() => {
    const initializeChallenge = async () => {
      try {
        setIsLoading(true);
        
        // If challengeId is provided, we're accepting a challenge
        if (challengeId) {
          const challengeData = await getChallenge(challengeId);
          setInviterInfo({
            inviter: challengeData.inviter,
            score: challengeData.score
          });
          setCurrentChallengeId(challengeId);
          
          if (challengeData.question_ids && challengeData.question_ids.length > 0) {
            setQuestionIDs(challengeData.question_ids);
          }
        } 
        // Otherwise, we're creating a new challenge
        else {
          const challengeData = await createChallenge();
          setCurrentChallengeId(challengeData.challenge_id);
          
          // Show share popup immediately for the challenger
          setShowSharePopup(true);
        }
        
        await fetchQuestion();
      } catch (err) {
        setError('Failed to initialize challenge. Please try again.');
        console.error('Error initializing challenge:', err);
      } finally {
        setIsLoading(false);
      }
    };
    
    initializeChallenge();
  }, []);
  
  const fetchQuestion = async () => {
    try {
      setIsLoading(true);
      setError('');
      setIsSecondClueVisible(false);
      setSelectedAnswer('');
      setAnswerResult(null);
      
      // Pass challenge_id when fetching questions if we're in a challenge
      const data = await getQuestion(currentChallengeId);
      setQuestion(data);
      
      // Store question ID for the challenge
      if (!challengeId) {
        setQuestionIDs(prevIds => [...prevIds, data.question_id]);
      }
    } catch (err) {
      setError('Failed to load question. Please try again.');
      console.error('Error fetching question:', err);
    } finally {
      setIsLoading(false);
    }
  };
  
  const handleRevealClue = async () => {
    try {
      await revealClue(currentChallengeId);
      setIsSecondClueVisible(true);
      
      // Update stats
      const newCluesRevealed = cluesRevealed + 1;
      const newScore = score - 1; // Subtract 1 point for revealing a clue
      
      setCluesRevealed(newCluesRevealed);
      setScore(newScore);
    } catch (err) {
      setError('Failed to reveal clue. Please try again.');
      console.error('Error revealing clue:', err);
    }
  };
  
  const handleSelectAnswer = async (city) => {
    if (answerResult || !question) return;
    
    setSelectedAnswer(city);
    
    try {
      const result = await submitAnswer(question.question_id, city, currentChallengeId);
      setAnswerResult(result);
      
      // Update stats based on answer
      if (result.correct) {
        const newCorrect = correctAnswers + 1;
        const newScore = score + 3; // Add 3 points for correct answer
        
        setCorrectAnswers(newCorrect);
        setScore(newScore);
        setShowConfetti(true);
      } else {
        const newIncorrect = incorrectAnswers + 1;
        setIncorrectAnswers(newIncorrect);
      }
    } catch (err) {
      setError('Failed to submit answer. Please try again.');
      console.error('Error submitting answer:', err);
    }
  };
  
  const handleNextQuestion = async () => {
    setShowConfetti(false);
    
    // If we're the challenger and this is the last question, end the challenge
    if (!challengeId && questionIDs.length >= 5) {
      try {
        await endChallenge(
          currentChallengeId, 
          score, 
          correctAnswers, 
          incorrectAnswers, 
          cluesRevealed,
          questionIDs
        );
        navigate('/challenge-complete', { 
          state: { 
            score, 
            correct_answers: correctAnswers, 
            incorrect_answers: incorrectAnswers, 
            clues_revealed: cluesRevealed, 
            created_challenge: true 
          } 
        });
      } catch (err) {
        setError('Failed to end challenge. Please try again.');
        console.error('Error ending challenge:', err);
      }
    } 
    // If we're accepting a challenge and this is the last question
    else if (challengeId && !questionIDs.length) {
      navigate('/challenge-complete', { 
        state: { 
          score, 
          correct_answers: correctAnswers, 
          incorrect_answers: incorrectAnswers, 
          clues_revealed: cluesRevealed,
          created_challenge: false,
          inviter_info: inviterInfo 
        } 
      });
    } 
    // Otherwise, fetch next question
    else {
      await fetchQuestion();
    }
  };
  
  const handleEndChallenge = async () => {
    try {
      await endChallenge(
        currentChallengeId, 
        score, 
        correctAnswers, 
        incorrectAnswers, 
        cluesRevealed,
        questionIDs
      );
      navigate('/challenge-complete', { 
        state: { 
          score, 
          correct_answers: correctAnswers, 
          incorrect_answers: incorrectAnswers, 
          clues_revealed: cluesRevealed, 
          created_challenge: true 
        } 
      });
    } catch (err) {
      setError('Failed to end challenge. Please try again.');
      console.error('Error ending challenge:', err);
    }
  };
  
  const handleBackToMenu = () => {
    navigate('/menu');
  };
  
  if (isLoading) {
    return (
      <div className="game-container">
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <p>Loading challenge...</p>
        </div>
      </div>
    );
  }
  
  if (error) {
    return (
      <div className="game-container">
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <p style={{ color: 'red' }}>{error}</p>
          <button className="form-button" onClick={fetchQuestion}>Try Again</button>
          <button className="form-button" style={{ marginTop: '10px' }} onClick={handleBackToMenu}>Back to Menu</button>
        </div>
      </div>
    );
  }
  
  return (
    <div className="game-container">
      {showConfetti && <Confetti recycle={false} numberOfPieces={200} />}
      
      <div className="game-header">
        <button className="reveal-button" onClick={handleBackToMenu}>
          Back to Menu
        </button>
        <div className="game-score">
          <p>Challenge Mode</p>
          <p>Score: {score}</p>
          <p>Correct: {correctAnswers} | Incorrect: {incorrectAnswers}</p>
        </div>
      </div>
      
      {inviterInfo && (
        <div style={{ marginBottom: '15px', padding: '10px', backgroundColor: '#f5f5f5', borderRadius: '4px' }}>
          <p>Challenge from: <strong>{inviterInfo.inviter}</strong></p>
          <p>Their score: <strong>{inviterInfo.score}</strong></p>
        </div>
      )}
      
      {question && (
        <>
          <div className="clue-container">
            <h2 style={{ marginBottom: '15px' }}>Guess the Destination</h2>
            <p className="clue-text">{question.clues[0]}</p>
            
            {!isSecondClueVisible ? (
              <div className="hidden-clue">
                <button className="reveal-button" onClick={handleRevealClue}>
                  Reveal Second Clue (-1 point)
                </button>
              </div>
            ) : (
              <div className="clue-text" style={{ marginTop: '15px', fontStyle: 'italic' }}>
                {question.clues[1]}
              </div>
            )}
          </div>
          
          {answerResult ? (
            <div className={`feedback-container ${answerResult.correct ? 'feedback-correct' : 'feedback-incorrect'}`}>
              <h3 className="feedback-text">
                {answerResult.correct ? '🎉 Correct Answer!' : '😢 Incorrect Answer'}
              </h3>
              <p>The answer was: <strong>{selectedAnswer}</strong></p>
              <p className="fun-fact">
                <strong>Fun Fact:</strong> {answerResult.funFact}
              </p>
              <button className="next-button" onClick={handleNextQuestion}>
                Next Destination
              </button>
            </div>
          ) : (
            <div className="choices-container">
              {question.choices.map((city, index) => (
                <button
                  key={index}
                  className="choice-button"
                  onClick={() => handleSelectAnswer(city)}
                >
                  {city}
                </button>
              ))}
            </div>
          )}
        </>
      )}
      
      {!challengeId && (
        <div style={{ textAlign: 'center', marginTop: '20px' }}>
          <button 
            className="form-button" 
            style={{ backgroundColor: '#ff9800' }} 
            onClick={handleEndChallenge}
          >
            End Challenge
          </button>
        </div>
      )}
      
      <SharePopup 
        challengeId={currentChallengeId}
        show={showSharePopup}
        onClose={() => setShowSharePopup(false)}
        inviter={username}
        score={score}
      />
    </div>
  );
};

export default Challenge;