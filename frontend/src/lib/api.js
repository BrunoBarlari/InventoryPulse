// API Client for InventoryPulse Backend

const API_URL = 'http://localhost:8080/api';

// Get stored token
function getToken() {
  return localStorage.getItem('access_token');
}

// Generic fetch wrapper with auth
async function fetchAPI(endpoint, options = {}) {
  const token = getToken();

  const config = {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  };

  if (token) {
    config.headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_URL}${endpoint}`, config);

  // Handle 401 - clear token and redirect
  if (response.status === 401) {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
    window.location.href = '/';
  }

  const data = await response.json().catch(() => null);

  if (!response.ok) {
    throw { status: response.status, ...data };
  }

  return data;
}

// Auth API
export const auth = {
  async login(email, password) {
    const data = await fetchAPI('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    });

    // Store tokens
    localStorage.setItem('access_token', data.access_token);
    localStorage.setItem('refresh_token', data.refresh_token);

    // Get user info
    const user = await this.me();
    localStorage.setItem('user', JSON.stringify(user));

    return { tokens: data, user };
  },

  async me() {
    return fetchAPI('/auth/me');
  },

  async register(email, password, role) {
    return fetchAPI('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, role }),
    });
  },

  async refresh() {
    const refreshToken = localStorage.getItem('refresh_token');
    if (!refreshToken) throw new Error('No refresh token');

    const data = await fetchAPI('/auth/refresh', {
      method: 'POST',
      body: JSON.stringify({ refresh_token: refreshToken }),
    });

    localStorage.setItem('access_token', data.access_token);
    localStorage.setItem('refresh_token', data.refresh_token);

    return data;
  },

  logout() {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('user');
  },

  isAuthenticated() {
    return !!getToken();
  },

  getUser() {
    const user = localStorage.getItem('user');
    return user ? JSON.parse(user) : null;
  },

  isAdmin() {
    const user = this.getUser();
    return user?.role === 'admin';
  }
};

// Categories API
export const categories = {
  async list(page = 1, pageSize = 10) {
    return fetchAPI(`/categories?page=${page}&page_size=${pageSize}`);
  },

  async get(id) {
    return fetchAPI(`/categories/${id}`);
  },

  async create(data) {
    return fetchAPI('/categories', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  async update(id, data) {
    return fetchAPI(`/categories/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  async delete(id) {
    return fetchAPI(`/categories/${id}`, {
      method: 'DELETE',
    });
  },
};

// Products API
export const products = {
  async list(page = 1, pageSize = 10, categoryId = null) {
    let url = `/products?page=${page}&page_size=${pageSize}`;
    if (categoryId) url += `&category_id=${categoryId}`;
    return fetchAPI(url);
  },

  async get(id) {
    return fetchAPI(`/products/${id}`);
  },

  async create(data) {
    return fetchAPI('/products', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  async update(id, data) {
    return fetchAPI(`/products/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  async updateStock(id, quantity) {
    return fetchAPI(`/products/${id}/stock`, {
      method: 'PATCH',
      body: JSON.stringify({ quantity }),
    });
  },

  async delete(id) {
    return fetchAPI(`/products/${id}`, {
      method: 'DELETE',
    });
  },
};

export default { auth, categories, products };

