import ky from 'ky';

const api = ky.create({
    prefixUrl: import.meta.env.VITE_APP_API_BASE_URL || 'http://localhost:8080',
    headers: {
        'Content-Type': 'application/json',
    },
    timeout: 10000,
})

export const apiClient = {
  // GET request
  get: async <T>(endpoint: string, searchParams?: Record<string, string | number>): Promise<T> => {
    return api.get(endpoint, { searchParams }).json<T>();
  },

  // POST request
  post: async <T, U = unknown>(endpoint: string, data?: U): Promise<T> => {
    return api.post(endpoint, { json: data }).json<T>();
  },

  // PUT request
  put: async <T, U = unknown>(endpoint: string, data?: U): Promise<T> => {
    return api.put(endpoint, { json: data }).json<T>();
  },

  // PATCH request
  patch: async <T, U = unknown>(endpoint: string, data?: U): Promise<T> => {
    return api.patch(endpoint, { json: data }).json<T>();
  },

  // DELETE request
  delete: async <T>(endpoint: string): Promise<T> => {
    return api.delete(endpoint).json<T>();
  },

  // Upload file
  upload: async <T>(endpoint: string, formData: FormData): Promise<T> => {
    return api.post(endpoint, { body: formData }).json<T>();
  }
}