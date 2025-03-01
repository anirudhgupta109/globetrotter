import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

// Create axios instance with default config
const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  }
});

// Add request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers['Authorization'] = token;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// User API calls
export const registerUser = async (username, password) => {
  try {
    const response = await api.post('/users/register', { username, password });
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : error;
  }
};

export const loginUser = async (username, password) => {
  try {
    const response = await api.post('/users/login', { username, password });
    localStorage.setItem('auth_token', response.data.auth_token);
    localStorage.setItem('username', response.data.username);
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : error;
  }
};

// Game API calls
export const getQuestion = async (challengeId = null) => {
  try {
    let url = '/game/question';
    if (challengeId) {
      url += `?challenge_id=${challengeId}`;
    }
    const response = await api.get(url);
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : error;
  }
};

export const submitAnswer = async (questionId, city, challengeId = null) => {
  try {
    const username = localStorage.getItem('username');
    const payload = {
      question_id: questionId,
      username,
      city,
      ...(challengeId && { challenge_id: challengeId })
    };
    const response = await api.post('/game/answer', payload);
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : error;
  }
};

export const revealClue = async (challengeId = null) => {
  try {
    const username = localStorage.getItem('username');
    const payload = {
      username,
      ...(challengeId && { challenge_id: challengeId })
    };
    const response = await api.post('/game/reveal-clue', payload);
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : error;
  }
};

// Challenge API calls
export const createChallenge = async () => {
  try {
    const username = localStorage.getItem('username');
    const response = await api.post('/challenges/create', { username });
    return { 
      challenge_id: response.data.challenge_id,
      inviter: response.data.inviter
    };
  } catch (error) {
    throw error.response ? error.response.data : error;
  }
};

export const getChallenge = async (challengeId) => {
  try {
    const response = await api.get(`/challenges/${challengeId}`, {
      params: { challenge_id: challengeId }
    });
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : error;
  }
};

export const endChallenge = async (challengeId, score, correctAnswers, incorrectAnswers, cluesRevealed, questionIDs) => {
  try {
    const username = localStorage.getItem('username');
    const payload = {
      challenge_id: challengeId,
      username,
      score,
      correct_answers: correctAnswers,
      incorrect_answers: incorrectAnswers,
      clues_revealed: cluesRevealed,
      question_ids: questionIDs
    };
    const response = await api.post('/challenges/end', payload);
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : error;
  }
};