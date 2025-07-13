import api from './auth'

export const usersAPI = {
  getUsers: async (page = 1, limit = 10, search = '') => {
    const response = await api.get('/users', {
      params: { page, limit, search },
    });
    return response.data;
  },
  
  getUser: async (id) => {
    const response = await api.get(`/users/${id}`);
    return response.data;
  },
  
  createUser: async (userData) => {
    const response = await api.post('/users', userData);
    return response.data;
  },
  
  updateUser: async (id, userData) => {
    const response = await api.put(`/users/${id}`, userData);
    return response.data;
  },
  
  deleteUser: async (id) => {
    const response = await api.delete(`/users/${id}`);
    return response.data;
  },
  
  changePassword: async (currentPassword, newPassword) => {
    const response = await api.post('/change-password', {
      current_password: currentPassword,
      new_password: newPassword,
    });
    return response.data;
  },
};