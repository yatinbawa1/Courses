<script>
	import { TopBar } from '$lib/components/self/top-bar';
	import { onMount } from 'svelte';
	let { children } = $props();

	onMount(() => {});
</script>

<TopBar />
{@render children()}
