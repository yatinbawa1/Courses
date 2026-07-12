<script lang="ts">
	let { children } = $props();
	import '../app.css';
	import { Toaster } from '$lib/components/ui/sonner';
</script>

<Toaster />
{@render children()}

<style>
</style>
