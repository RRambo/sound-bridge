import axios from 'axios';

// Create axios instance with base configuration
const api = axios.create({
  baseURL: '/api',
  auth: {
    username: 'kids_noisemeter_admin',
    password: 'passwordkids'
  }
  // Mind that this auth still needs futuer extensions.
});

// Add request interceptor to handle Content-Type header
api.interceptors.request.use(
  (config) => {
    // For POST, PUT, PATCH requests, always set Content-Type
    // Even if there's no body, the backend middleware expects it
    if (['post', 'put', 'patch'].includes(config.method.toLowerCase())) {
      config.headers['Content-Type'] = 'application/json';
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Add response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // Log error details for debugging
    if (error.response) {
      console.error('API Error:', {
        status: error.response.status,
        data: error.response.data,
        url: error.config.url,
        method: error.config.method
      });
    }
    return Promise.reject(error);
  }
);

// Location-related API calls
export const locationAPI = {
  // Get all locations
  getAll: () => api.get('/locations'),
  
  // Create a new location
  create: (name) => api.post('/locations', { 
    name, 
    chosen: true 
  }),
  
  // Set a location as chosen
  // Note: Even though we don't send a body, the backend expects Content-Type
  setChosen: (id) => api.put(`/locations/${id}`, null, {
    headers: {
      'Content-Type': 'application/json'
    }
  }),
  // Had to set headers manually here to satisfy backend middleware

  // Update threshold of location
  updateThreshold: (id, newThreshold) => api.put(`/locations/${id}`, null, {
    params: { newThreshold }
  }),
  
  // Delete a location
  delete: (id) => api.delete(`/locations/${id}`)
};

// Data-related API calls (for future use)
export const dataAPI = {
  getAll: () => api.get('/data'),
  getById: (id) => api.get(`/data/${id}`),
  getByRoom: (room) => api.get(`/data/weekly/${room}`),
  getDailySummary: (room, date) => api.get(`/data/daily/${room}`, {
    params: { date }
  }),
  create: (data) => api.post('/data', data),
  update: (data) => api.put('/data', data),
  delete: (id) => api.delete(`/data/${id}`)
};

export default api;