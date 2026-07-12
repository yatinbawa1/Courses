<script>
	import AvatarBadge from '$lib/components/ui/avatar/avatar-badge.svelte';
	import AvatarFallback from '$lib/components/ui/avatar/avatar-fallback.svelte';
	import AvatarImage from '$lib/components/ui/avatar/avatar-image.svelte';
	import Avatar from '$lib/components/ui/avatar/avatar.svelte';
	import CardContent from '$lib/components/ui/card/card-content.svelte';
	import CardDescription from '$lib/components/ui/card/card-description.svelte';
	import CardTitle from '$lib/components/ui/card/card-title.svelte';
	import Card from '$lib/components/ui/card/card.svelte';
	import InputGroupInput from '$lib/components/ui/input-group/input-group-input.svelte';
	import InputGroup from '$lib/components/ui/input-group/input-group.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import { UserCog, ImageUp, UserRoundPen, SquarePen, Image } from '@lucide/svelte';
	import { CardFooter } from '$lib/components/ui/card';
	import { api } from '$lib/api/api';
	import { auth } from '$lib/stores/authStore/authStore';
	import { toast } from 'svelte-sonner';
	import { invalidateAll } from '$app/navigation';

	let name = $state('');

	const HandleSubmit = async () => {
		const res = await api.post('/api/account/update-user', {
			user_id: $auth.user_id,
			name: name,
			profile_photo_url: '',
			email: $auth.email
		});

		if (res.status != 200) {
			toast.error('Unable to save user data');
		} else {
			$auth.name = name;
			// TODO
			// Should add the url as well when It has been created
			invalidateAll();
		}
	};
</script>

<div class="w-full h-screen flex items-center justify-center background">
	<Card class="p-5 z-10">
		<CardTitle class="text-3xl flex items-center gap-1 px-5">
			<UserCog size="25px" />Complete Your Profile</CardTitle
		>
		<CardDescription class="px-5 -mt-2.5 text-center"
			>Help us get to know you better</CardDescription
		>
		<CardContent class="flex justify-center items-center flex-col gap-5">
			<Avatar class="h-30 w-30">
				<AvatarImage src="" />
				<div
					class="transition-all duration-75 absolute bottom-0 right-0 h-10 w-10 rounded-full text-white hover:bg-gray-600 hover:cursor-pointer bg-gray-400 flex items-center justify-center p-2"
				>
					<ImageUp />
				</div>
				<AvatarFallback>
					<UserRoundPen />
				</AvatarFallback>
			</Avatar>
			<div class="w-full">
				<label for="username" class="playfair text-base text-gray-600">Username</label>
				<InputGroup>
					<InputGroupInput bind:value={name} id="username" placeholder="eg. username_12" />
				</InputGroup>
			</div>
			<Button
				type="submit"
				variant="black"
				onclick={HandleSubmit}
				size="lg"
				class="py-5 w-full hover:bg-mauve-700 hover:cursor-pointer"
				>Complete Onboarding <SquarePen size="1rem" /></Button
			>
		</CardContent>
		<CardFooter>
			<Button variant="ghost" class="w-full">Skip For Now</Button>
		</CardFooter>
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
