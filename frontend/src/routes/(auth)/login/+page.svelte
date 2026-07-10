<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import { LogIn, Mail, KeyIcon, User } from '@lucide/svelte';
	import InputGroup from '$lib/components/ui/input-group/input-group.svelte';
	import InputGroupAddon from '$lib/components/ui/input-group/input-group-addon.svelte';
	import InputGroupInput from '$lib/components/ui/input-group/input-group-input.svelte';

	let email = $state('');
	let password = $state('');

	let touches = $state({
		email: false,
		password: false
	});

	const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
	const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#$%^&*(),.?":{}|<>_\-\[\]]).{9,}$/;

	const errors = $derived({
		email: email.length > 0 && !emailRegex.test(email) ? 'Email is not correct' : '',
		password: password.length > 0 && !passwordRegex.test(password) ? 'Password Unsecure' : ''
	});

	const isFormInvalid = $derived(!emailRegex.test(email) || !passwordRegex.test(password));

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

<div>
	<div class="flex items-center">
		<User class="mr-2" />
		<p class="text-lg">Log In</p>
	</div>

	<p class="text-mauve-600 mb-4">Log into your Account</p>
</div>
<form onsubmit={handleLogin} class="space-y-4">
	<label for="email" class="playfair text-base text-mauve-800">Email</label>
	<InputGroup class="py-5" id="email">
		<InputGroupAddon>
			<Mail />
		</InputGroupAddon>

		<InputGroupInput
			type="email"
			bind:value={email}
			placeholder="Email"
			onblur={() => (touches.email = true)}
		></InputGroupInput>
	</InputGroup>
	{#if errors.email && touches.email}
		<p class="text-xs text-red-500 mt-1">{errors.email}</p>
	{/if}

	<label for="Password" class="playfair text-base text-mauve-800">Password</label>
	<InputGroup id="email" class="py-5">
		<InputGroupAddon>
			<KeyIcon />
		</InputGroupAddon>

		<InputGroupInput
			type="password"
			bind:value={password}
			placeholder="Password"
			onblur={() => (touches.password = true)}
		></InputGroupInput>
	</InputGroup>
	{#if errors.password && touches.password}
		<p class="text-xs text-red-500 mt-1">{errors.password}</p>
	{/if}

	<Button
		type="submit"
		variant="black"
		disabled={isFormInvalid}
		size="lg"
		class="py-5 w-full hover:bg-mauve-700 hover:cursor-pointer"
		>Sign in <LogIn size="1rem" /></Button
	>

	<p class="text-center text-base text-muted-foreground playfair">
		Don't have an account?
		<a href="/signup" class="text-primary underline-offset-4 hover:underline">Sign up</a>
	</p>
</form>
