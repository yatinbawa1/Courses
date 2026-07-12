import { browser } from '$app/env';
import type { LayoutLoad } from './$types';
import { get } from 'svelte/store';
import { auth } from '$lib/stores/authStore/authStore';
import { redirect } from '@sveltejs/kit';
import { page } from '$app/state';

export const load: LayoutLoad = async ({ fetch, url }) => {
	if (!browser) return;

	const path = url.pathname;
	let currentAuth = get(auth);

	if (!currentAuth.isAuthenticated) {
		try {
			const res = await fetch('/api/auth/refresh');
			if (res.status === 200) {
				const userData = await res.json();
				auth.set({
					isAuthenticated: true,
					email: userData.email,
					name: userData.name,
					profile_photo_url: userData.profile_photo_url,
					user_id: userData.user_id
				});

				// Refresh our local reference variable
				currentAuth = get(auth);
			}
		} catch (err) {
			console.error('Silent auth refresh failed:', err);
		}
	}

	const isAuthRoute =
		path.startsWith('/login') || path.startsWith('/signup') || path.startsWith('/verify-otp');
	const isOnboardingRoute = path.startsWith('/onboarding');
	if (!currentAuth.isAuthenticated) {
		if (!isAuthRoute && path !== '/') {
			throw redirect(303, '/login');
		}
		return;
	}

	if (currentAuth.name == null && !isOnboardingRoute) {
		throw redirect(303, '/onboarding');
	}

	if (currentAuth.name != null && isOnboardingRoute) {
		throw redirect(303, '/app');
	}

	if (isAuthRoute || path === '/') {
		throw redirect(303, '/app');
	}
};

export const ssr = false;
export const prerender = false;
