import { browser } from '$app/environment';
import { writable } from 'svelte/store';

const initialState = {
	isAuthenticated: false,
	email: '',
	name: '',
	profile_photo_exists: false,
	user_id: ''
};

export const auth = writable(initialState);

export function logoutUser() {
	if (!browser) return;
	auth.set(initialState);
}
