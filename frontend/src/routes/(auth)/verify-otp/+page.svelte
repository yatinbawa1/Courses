<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import Button from '$lib/components/ui/button/button.svelte';
	import InputGroupAddon from '$lib/components/ui/input-group/input-group-addon.svelte';
	import InputGroupInput from '$lib/components/ui/input-group/input-group-input.svelte';
	import InputGroup from '$lib/components/ui/input-group/input-group.svelte';
	import { BadgeCheck, FingerprintPattern, TicketCheck } from '@lucide/svelte';
	import { toast } from 'svelte-sonner';

	let otp = $state('');
	const email = (page.state as any).email || '';
	const password = (page.state as any).password || '';

	const VerifyOTP = async () => {
		const response = await fetch('/api/auth/send-otp/verify', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ email, password, otp })
		});

		switch (response.status) {
			case 201:
				toast.info('User Created! Login to continue with app');
				goto('/login');
				break;
			case 400:
				toast.error('Wrong OTP! Correct OTP has been sent on your email');
				break;
			default:
				toast.error('Unable To Verify');
		}
	};
</script>

<div>
	<div class="flex items-center">
		<BadgeCheck class="mr-2" />
		<p class="text-lg">Verify OTP</p>
	</div>

	<p class="text-mauve-600 mb-8">We have sent you a verification email!</p>
</div>

<InputGroup class="py-5">
	<InputGroupAddon>
		<FingerprintPattern />
	</InputGroupAddon>
	<InputGroupInput bind:value={otp} placeholder="Enter OTP"></InputGroupInput>
</InputGroup>

<Button class="mt-8 w-full py-5 hover:bg-green-900 hover:cursor-pointer" onclick={VerifyOTP}
	>Verify <TicketCheck /></Button
>
