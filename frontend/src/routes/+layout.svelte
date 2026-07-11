<script lang="ts">
	import { goto } from '$app/navigation';
	let { children } = $props();
	import { onMount } from 'svelte';
	import '../app.css';
	import { Toaster } from '$lib/components/ui/sonner';
	import { page } from '$app/state';

	onMount(async () => {
		try {
			const res = await fetch('/api/auth/refresh');

			if (res.status === 200) {
				if (page.url.pathname === '/login' || page.url.pathname === '/') {
					goto('/app', { replaceState: true });
				}
			} else {
				goto('/login', { replaceState: true });
			}
		} catch (error) {
			goto('/login', { replaceState: true });
		}
	});
</script>

<Toaster />
{@render children()}

<style>
</style>
