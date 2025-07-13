import api from './auth'

export const adminAPI = {
  getDashboard: async () => {
    const response = await api.get('/admin/dashboard');
    return response.data;
  },
  
  getSystemInfo: async () => {
    const response = await api.get('/admin/system');
    return response.data;
  },
  
  purgeOldLogs: async (days = 30) => {
    const response = await api.delete('/admin/logs/purge', {
      params: { days },
    });
    return response.data;
  },
  
  getHealth: async () => {
    const response = await api.get('/health');
    return response.data;
  },
};