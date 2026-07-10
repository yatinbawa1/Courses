<script lang="ts">
	import { goto } from '$app/navigation';
	let { children } = $props();
	import { onMount } from 'svelte';
	import '../app.css';
	import { Toaster } from '$lib/components/ui/sonner';

	onMount(async () => {
		const res = await fetch('/api/auth/refresh');
		if (res.status === 200) {
			goto('/', {
				replaceState: true
			});
		} else {
			goto('/login', {
				replaceState: true
			});
		}
	});
</script>

<Toaster />
{@render children()}

<style>
</style>
