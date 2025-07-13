import api from './auth'

export const logsAPI = {
  getLogs: async (params = {}) => {
    const response = await api.get('/logs', { params });
    return response.data;
  },
  
  getLogStats: async (fromDate, toDate) => {
    const response = await api.get('/logs/stats', {
      params: { from_date: fromDate, to_date: toDate },
    });
    return response.data;
  },
};