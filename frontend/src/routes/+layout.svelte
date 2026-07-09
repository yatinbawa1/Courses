<script lang="ts">
	import { goto } from '$app/navigation';
	let { children } = $props();
	import { onMount } from 'svelte';

	onMount(async () => {
		const res = await fetch('/api/auth/refresh');
		if (res.status === 200) {
			goto('/');
		} else {
			goto('/login');
		}
	});
</script>

{@render children()}
