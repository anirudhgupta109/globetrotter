import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Confetti from 'react-confetti';
import { getQuestion, submitAnswer, revealClue } from '../services/api';

const QuestionPage = () => {
  const navigate = useNavigate();
  const [question, setQuestion] = useState(null);
  const [isSecondClueVisible, setIsSecondClueVisible] = useState(false);
  const [selectedAnswer, setSelectedAnswer] = useState('');
  const [answerResult, setAnswerResult] = useState(null);
  const [showConfetti, setShowConfetti] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  
  // Game stats
  const [score, setScore] = useState(0);
  const [correctAnswers, setCorrectAnswers] = useState(0);
  const [incorrectAnswers, setIncorrectAnswers] = useState(0);
  const [cluesRevealed, setCluesRevealed] = useState(0);
  
  const fetchQuestion = async () => {
    try {
      setIsLoading(true);
      setError('');
      setIsSecondClueVisible(false);
      setSelectedAnswer('');
      setAnswerResult(null);
      
      const data = await getQuestion();
      setQuestion(data);
    } catch (err) {
      setError('Failed to load question. Please try again.');
      console.error('Error fetching question:', err);
    } finally {
      setIsLoading(false);
    }
  };
  
  useEffect(() => {
    // Initialize game stats from localStorage or set to 0
    const savedScore = parseInt(localStorage.getItem('score') || '0');
    const savedCorrect = parseInt(localStorage.getItem('correctAnswers') || '0');
    const savedIncorrect = parseInt(localStorage.getItem('incorrectAnswers') || '0');
    const savedCluesRevealed = parseInt(localStorage.getItem('cluesRevealed') || '0');
    
    setScore(savedScore);
    setCorrectAnswers(savedCorrect);
    setIncorrectAnswers(savedIncorrect);
    setCluesRevealed(savedCluesRevealed);
    
    fetchQuestion();
  }, []);
  
  const handleRevealClue = async () => {
    try {
      await revealClue();
      setIsSecondClueVisible(true);
      
      // Update stats
      const newCluesRevealed = cluesRevealed + 1;
      const newScore = score - 1; // Subtract 1 point for revealing a clue
      
      setCluesRevealed(newCluesRevealed);
      setScore(newScore);
      
      // Update localStorage
      localStorage.setItem('cluesRevealed', newCluesRevealed);
      localStorage.setItem('score', newScore);
    } catch (err) {
      setError('Failed to reveal clue. Please try again.');
      console.error('Error revealing clue:', err);
    }
  };
  
  const handleSelectAnswer = async (city) => {
    if (answerResult || !question) return;
    
    setSelectedAnswer(city);
    
    try {
      const result = await submitAnswer(question.question_id, city);
      setAnswerResult(result);
      
      // Update stats based on answer
      if (result.correct) {
        const newCorrect = correctAnswers + 1;
        const newScore = score + 3; // Add 3 points for correct answer
        
        setCorrectAnswers(newCorrect);
        setScore(newScore);
        setShowConfetti(true);
        
        // Update localStorage
        localStorage.setItem('correctAnswers', newCorrect);
        localStorage.setItem('score', newScore);
      } else {
        const newIncorrect = incorrectAnswers + 1;
        setIncorrectAnswers(newIncorrect);
        
        // Update localStorage
        localStorage.setItem('incorrectAnswers', newIncorrect);
      }
    } catch (err) {
      setError('Failed to submit answer. Please try again.');
      console.error('Error submitting answer:', err);
    }
  };
  
  const handleNextQuestion = () => {
    setShowConfetti(false);
    fetchQuestion();
  };
  
  const handleBackToMenu = () => {
    navigate('/menu');
  };
  
  if (isLoading) {
    return (
      <div className="game-container">
        <div style={{ textAlign: 'center', padding: '40px' }}>
          <p>Loading question...</p>
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
          <p>Score: {score}</p>
          <p>Correct: {correctAnswers} | Incorrect: {incorrectAnswers}</p>
        </div>
      </div>
      
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
                <strong>Fun Fact:</strong> {answerResult.fun_fact}
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
    </div>
  );
};

export default QuestionPage;