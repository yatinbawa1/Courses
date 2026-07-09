<script>
	import { goto } from '$app/navigation';
	import Button from '$lib/components/ui/button/button.svelte';

	let email = $state('');
	let password = $state('');
	let output_text = $state('');

	const TryToLogIn = async () => {
		try {
			const response = await fetch('/api/auth/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include',
				body: JSON.stringify({
					email: email,
					password: password
				})
			});

			if (response.status != 200) {
				output_text = 'Unable to Login';
			} else {
				output_text = 'Success';
				goto('/');
			}
		} catch (err) {
			console.log(err);
		}
	};
</script>

<main>
	<h1>Login</h1>
	<input bind:value={email} placeholder="Email" />
	<input bind:value={password} type="password" placeholder="Password " />
	<Button variant="outline" onclick={TryToLogIn}>Submit</Button>

	<i>{output_text}</i>
</main>
