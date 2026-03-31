const API_BASE = localStorage.getItem('api_base') || '/api';

function getToken() {
  return localStorage.getItem('blog_token') || '';
}

function setToken(token) {
  if (token) localStorage.setItem('blog_token', token);
  else localStorage.removeItem('blog_token');
}

async function request(path, options = {}) {
  const headers = { ...(options.headers || {}) };
  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = headers['Content-Type'] || 'application/json';
  }

  const token = getToken();
  if (token) headers.Authorization = `Bearer ${token}`;

  const res = await fetch(`${API_BASE}${path}`, { ...options, headers });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.error || `HTTP ${res.status}`);
  return data;
}

const api = {
  login: (username, password) => request('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  }),

  listPosts: (params) => request(`/posts?${new URLSearchParams(params)}`),
  getPost: (slug) => request(`/posts/${slug}`),

  listCategories: () => request('/categories'),
  createCategory: (name) => request('/categories', { method: 'POST', body: JSON.stringify({ name }) }),
  deleteCategory: (id) => request(`/categories/${id}`, { method: 'DELETE' }),

  listTags: () => request('/tags'),
  createTag: (name) => request('/tags', { method: 'POST', body: JSON.stringify({ name }) }),
  deleteTag: (id) => request(`/tags/${id}`, { method: 'DELETE' }),

  listAdminPosts: (params) => request(`/admin/posts?${new URLSearchParams(params)}`),
  getAnyPost: (slug) => request(`/admin/posts/${slug}`),

  createPost: (payload) => request('/posts', { method: 'POST', body: JSON.stringify(payload) }),
  updatePost: (slug, payload) => request(`/posts/${slug}`, { method: 'PUT', body: JSON.stringify(payload) }),
  deletePost: (slug) => request(`/posts/${slug}`, { method: 'DELETE' }),

  uploadImage: (file) => {
    const fd = new FormData();
    fd.append('image', file);
    return request('/upload/image', { method: 'POST', body: fd, headers: {} });
  },
};

window.blogApi = api;
window.blogAuth = { getToken, setToken };
