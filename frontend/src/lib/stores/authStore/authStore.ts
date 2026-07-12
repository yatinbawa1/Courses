import { browser } from '$app/env';
import { api } from '$lib/api/api';
import { writable, get } from 'svelte/store';

const initialState = {
	isAuthenticated: false,
	email: '',
	name: '',
	profile_photo_url: '',
	user_id: ''
};

export const auth = writable(initialState);

export function logoutUser() {
	if (!browser) return;

	const currentAuthState = get(auth);
	api.post(`/auth/logout/${currentAuthState.user_id}`);
	auth.set(initialState);
}
