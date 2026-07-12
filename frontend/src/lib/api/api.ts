import axios from 'axios';
import { goto } from '$app/navigation';
import { browser } from '$app/environment';

export const api = axios.create({
	withCredentials: true,
	headers: {
		'Content-Type': 'application/json'
	}
});

api.interceptors.response.use(
	(response) => response,
	async (error) => {
		const originalRequest = error.config;

		if (browser && error.response && error.response.status === 401 && !originalRequest._retry) {
			originalRequest._retry = true;
			try {
				await axios.get('/api/auth/refresh', { withCredentials: true });
				return api(originalRequest);
			} catch (refreshError) {
				goto('/login');
				return Promise.reject(refreshError);
			}
		}

		return Promise.reject(error);
	}
);
