<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '$lib/components/ui/button/button.svelte';
	import CardFooter from '$lib/components/ui/card/card-footer.svelte';
	import Input from '$lib/components/ui/input/input.svelte';
	import { toast } from 'svelte-sonner';
	import { LogIn, Mail, KeyIcon } from '@lucide/svelte';
	import InputGroup from '$lib/components/ui/input-group/input-group.svelte';
	import InputGroupAddon from '$lib/components/ui/input-group/input-group-addon.svelte';
	import InputGroupInput from '$lib/components/ui/input-group/input-group-input.svelte';

	let email = $state('');
	let password = $state('');

	const handleLogin = async () => {
		try {
			const response = await fetch('/api/auth/login', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				credentials: 'include',
				body: JSON.stringify({ email, password })
			});

			if (response.status !== 200) {
				toast.error('Unable to login');
			} else {
				goto('/');
			}
		} catch (err) {
			toast.error('An error occurred during login');
		}
	};
</script>

<form onsubmit={handleLogin} class="space-y-4">
	<label for="email" class="playfair text-base text-mauve-800">Email</label>
	<InputGroup class="py-5" id="email">
		<InputGroupAddon>
			<Mail />
		</InputGroupAddon>

		<InputGroupInput type="email" bind:value={email} placeholder="Email"></InputGroupInput>
	</InputGroup>

	<label for="Password" class="playfair text-base text-mauve-800">Password</label>
	<InputGroup id="email" class="py-5">
		<InputGroupAddon>
			<KeyIcon />
		</InputGroupAddon>

		<InputGroupInput type="password" bind:value={password} placeholder="Password"></InputGroupInput>
	</InputGroup>

	<Button
		type="submit"
		variant="black"
		size="lg"
		class="py-5 w-full hover:bg-mauve-700 hover:cursor-pointer"
		>Sign in <LogIn size="1rem" /></Button
	>

	<p class="text-center text-base text-muted-foreground playfair">
		Don't have an account?
		<a href="/signup" class="text-primary underline-offset-4 hover:underline">Sign up</a>
	</p>
</form>
