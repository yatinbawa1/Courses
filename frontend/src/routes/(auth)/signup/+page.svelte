<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '$lib/components/ui/button/button.svelte';
	import InputGroup from '$lib/components/ui/input-group/input-group.svelte';
	import InputGroupAddon from '$lib/components/ui/input-group/input-group-addon.svelte';
	import InputGroupInput from '$lib/components/ui/input-group/input-group-input.svelte';
	import { LogIn, Mail, KeyIcon } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';

	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');

	let touched = $state({
		email: false,
		password: false,
		confirmPassword: false
	});

	const emailRegex = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
	const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[!@#$%^&*(),.?":{}|<>_\-\[\]]).{9,}$/;

	const errors = $derived({
		email: email.length > 0 && !emailRegex.test(email) ? 'Email is not correct' : '',
		password: password.length > 0 && !passwordRegex.test(password) ? 'Password Unsecure' : '',
		confirmPassword:
			password.length > 0 && confirmPassword.length > 0 && password !== confirmPassword
				? 'Confirm Password Does Not Match Password'
				: ''
	});

	const isFormInvalid = $derived(
		!emailRegex.test(email) || !passwordRegex.test(password) || password !== confirmPassword
	);

	const handleSignup = async (e: Event) => {
		e.preventDefault(); // Good practice to explicitly handle

		touched.email = true;
		touched.password = true;
		touched.confirmPassword = true;

		if (isFormInvalid) {
			toast.error('Please fix the errors before submitting.');
			return;
		}

		try {
			const response = await fetch('/api/auth/send-otp', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				credentials: 'include',
				body: JSON.stringify({ email, password })
			});

			if (response.status !== 200) {
				toast.error('Unable to Create Account');
			} else {
				goto('/verify-otp');
			}
		} catch (err) {
			console.error(err);
			toast.error('An unexpected error occurred');
		}
	};
</script>

<form onsubmit={handleSignup} class="space-y-4">
	<div>
		<label for="email" class="playfair text-base text-mauve-800">Email</label>
		<InputGroup class="py-5" id="email">
			<InputGroupAddon><Mail /></InputGroupAddon>
			<InputGroupInput
				type="email"
				bind:value={email}
				placeholder="Email"
				onblur={() => (touched.email = true)}
			/>
		</InputGroup>
		{#if touched.email && errors.email}
			<p class="text-xs text-red-500 mt-1">{errors.email}</p>
		{/if}
	</div>

	<div>
		<label for="password" class="playfair text-base text-mauve-800">Password</label>
		<InputGroup class="py-5" id="password">
			<InputGroupAddon><KeyIcon /></InputGroupAddon>
			<InputGroupInput
				type="password"
				bind:value={password}
				placeholder="Password"
				onblur={() => (touched.password = true)}
			/>
		</InputGroup>
		{#if touched.password && errors.password}
			<p class="text-xs text-red-500 mt-1">{errors.password}</p>
		{/if}
	</div>

	<div>
		<label for="confirmPassword" class="playfair text-base text-mauve-800">Confirm Password</label>
		<InputGroup class="py-5" id="confirmPassword">
			<InputGroupAddon><KeyIcon /></InputGroupAddon>
			<InputGroupInput
				type="password"
				bind:value={confirmPassword}
				placeholder="Confirm Password"
				onblur={() => (touched.confirmPassword = true)}
			/>
		</InputGroup>
		{#if touched.confirmPassword && errors.confirmPassword}
			<p class="text-xs text-red-500 mt-1">{errors.confirmPassword}</p>
		{/if}
	</div>

	<Button
		type="submit"
		variant="black"
		size="lg"
		disabled={isFormInvalid && (touched.email || touched.password || touched.confirmPassword)}
		class="py-5 w-full hover:bg-mauve-700 hover:cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
	>
		SEND OTP <LogIn size="1rem" />
	</Button>

	<p class="text-center text-base text-muted-foreground playfair">
		Already have an account?
		<a href="/login" class="text-primary underline-offset-4 hover:underline">Log In</a>
	</p>
</form>
