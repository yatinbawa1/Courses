<script lang="ts">
	let { children } = $props();
	import CardContent from '$lib/components/ui/card/card-content.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import CardTitle from '$lib/components/ui/card/card-title.svelte';
	import CardFooter from '$lib/components/ui/card/card-footer.svelte';
</script>

<div class="background flex min-h-screen items-center justify-center">
	<Card class="w-full max-w-lg p-5 m-4 z-10">
		<CardTitle>
			<div class="text-center">
				<h1 class="text-2xl font-semibold tracking-widest">COURSES</h1>
			</div>
		</CardTitle>
		<CardContent>
			{@render children()}
		</CardContent>

		<CardFooter
			>Courses Project: <a
				class="hover:underline underline-offset-4 text-accent"
				href="https://github.com/yatinbawa1/">©Yatin Bawa</a
			></CardFooter
		>
	</Card>
</div>

<style>
	.background {
		width: 100%;
		background: url('https://images.unsplash.com/photo-1783228905491-7f25b72eabf7?q=80&w=2670&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D');
		background-size: cover;
		background-position: center;
	}
	.background::after {
		content: '';
		position: absolute;
		width: 100%;
		height: 100%;
		background-color: rgba(0, 0, 0, 0.4);
		z-index: 1;
	}
</style>
