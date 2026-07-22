<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import { api } from '$lib/api/api';
	import AvatarFallback from '$lib/components/ui/avatar/avatar-fallback.svelte';
	import AvatarImage from '$lib/components/ui/avatar/avatar-image.svelte';
	import Avatar from '$lib/components/ui/avatar/avatar.svelte';
	import DropdownMenuContent from '$lib/components/ui/dropdown-menu/dropdown-menu-content.svelte';
	import DropdownMenuGroup from '$lib/components/ui/dropdown-menu/dropdown-menu-group.svelte';
	import DropdownMenuItem from '$lib/components/ui/dropdown-menu/dropdown-menu-item.svelte';
	import DropdownMenuLabel from '$lib/components/ui/dropdown-menu/dropdown-menu-label.svelte';
	import DropdownMenuTrigger from '$lib/components/ui/dropdown-menu/dropdown-menu-trigger.svelte';
	import DropdownMenu from '$lib/components/ui/dropdown-menu/dropdown-menu.svelte';
	import InputGroupAddon from '$lib/components/ui/input-group/input-group-addon.svelte';
	import InputGroupInput from '$lib/components/ui/input-group/input-group-input.svelte';
	import InputGroup from '$lib/components/ui/input-group/input-group.svelte';
	import Spinner from '$lib/components/ui/spinner/spinner.svelte';
	import { auth, logoutUser } from '$lib/stores/authStore/authStore';
	import {
		Search,
		House,
		GraduationCap,
		ShoppingBasket,
		Pen,
		Wallet,
		LogOut,
		User,
		Menu,
		X
	} from '@lucide/svelte';
	import { onMount, type Component } from 'svelte';
	import { toast } from 'svelte-sonner';

	interface MenuItem {
		icon: Component;
		link: string;
		label: string;
	}

	let menuItems: MenuItem[] = [
		{
			icon: House,
			link: '/app',
			label: 'Home'
		},
		{
			icon: ShoppingBasket,
			link: '/app/market',
			label: 'Market'
		},
		{
			icon: GraduationCap,
			link: '/app/my-courses',
			label: 'My Courses'
		},
		{
			icon: Pen,
			link: '/app/creator',
			label: 'Creator'
		},
		{
			icon: Wallet,
			link: '/app/wallet',
			label: 'wallet'
		}
	];
	const standardImage =
		'https://images.unsplash.com/photo-1773332611476-6ec2ba68049f?q=80&w=1830&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDF8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D';

	let profileUrl = $derived.by(() => {
		if (!$auth.profile_photo_exists) {
			return standardImage;
		}

		return `https://courses-content-portfolio-go-next.s3.ap-south-1.amazonaws.com/users/${$auth.user_id}`;
	});

	let activeLink = $derived(page.url.pathname);

	const logout = async () => {
		const res = await api.post(`/api/auth/logout/${$auth.email}`);
		if (res.status == 200) {
			logoutUser();
			invalidateAll();
		} else {
			toast.error('Unable to logout!');
		}
	};
</script>

<div class="w-full bg-gray-900 px-5 py-4">
	<input type="checkbox" id="nav-toggle" class="hidden peer" />

	<div class="flex items-center justify-between md:grid md:grid-cols-3 md:items-center">
		<div class="playfair text-2xl tracking-[5px] text-white">COURSES</div>

		<div class="hidden md:flex md:justify-center">
			<div class="flex gap-2">
				{#each menuItems as item (item.link)}
					{@const Icon = item.icon}
					<div>
						<a
							href={item.link}
							aria-label={item.label}
							class="flex items-center justify-center w-10 h-10 transition-colors duration-200
									{activeLink === item.link
								? 'bg-gray-800 text-white border border-gray-700'
								: 'text-gray-400 hover:bg-gray-800/50 hover:text-gray-200'}"
						>
							<Icon class="w-5 h-5" />
						</a>
					</div>
				{/each}
			</div>
		</div>

		<div class="hidden md:flex w-full justify-end items-center gap-5">
			<InputGroup class="border w-full py-5 border-gray-800">
				<InputGroupAddon>
					<Search />
				</InputGroupAddon>
				<InputGroupInput placeholder="Search" class="text-gray-300"></InputGroupInput>
			</InputGroup>
			<DropdownMenu>
				<DropdownMenuTrigger>
					<Avatar>
						<AvatarImage src={profileUrl} />
						<AvatarFallback>
							<Spinner />
						</AvatarFallback>
					</Avatar>
				</DropdownMenuTrigger>
				<DropdownMenuContent class="bg-gray-800">
					<DropdownMenuGroup>
						<DropdownMenuLabel class="flex items-center">
							<User size="1rem" /> Profile</DropdownMenuLabel
						>
						<DropdownMenuItem onclick={logout} class="text-gray-500">
							<LogOut />
							<span>Log Out</span>
						</DropdownMenuItem>
					</DropdownMenuGroup>
				</DropdownMenuContent>
			</DropdownMenu>
		</div>

		<label
			for="nav-toggle"
			class="flex md:hidden peer-checked:hidden items-center justify-center w-10 h-10 cursor-pointer text-gray-400 hover:text-white"
		>
			<Menu class="w-6 h-6" />
		</label>

		<!-- Mobile: close icon (visible when menu is open) -->
		<label
			for="nav-toggle"
			class="hidden md:hidden peer-checked:flex items-center justify-center w-10 h-10 cursor-pointer text-gray-400 hover:text-white"
		>
			<X class="w-6 h-6" />
		</label>
	</div>

	<div class="hidden peer-checked:flex flex-col gap-6 mt-4 md:hidden">
		<div class="">
			<div class="w-full flex flex-col gap-1">
				{#each menuItems as item (item.link)}
					{@const Icon = item.icon}
					<div class="w-full">
						<a
							href={item.link}
							class="flex items-center gap-3 w-full px-3 py-2 transition-colors duration-200
								{activeLink === item.link
								? 'bg-gray-800 text-white border border-gray-700'
								: 'text-gray-400 hover:bg-gray-800/50 hover:text-gray-200'}"
						>
							<Icon class="w-5 h-5" />
							<span class="text-sm">{item.label}</span>
						</a>
					</div>
				{/each}
			</div>
		</div>

		<div class="flex gap-2">
			<InputGroup class="border py-5 border-gray-800">
				<InputGroupAddon>
					<Search />
				</InputGroupAddon>
				<InputGroupInput placeholder="Search" class="text-gray-300"></InputGroupInput>
			</InputGroup>

			<DropdownMenu>
				<DropdownMenuTrigger>
					<Avatar>
						<AvatarImage src={profileUrl} />
						<AvatarFallback>CN</AvatarFallback>
					</Avatar>
				</DropdownMenuTrigger>
				<DropdownMenuContent class="bg-gray-800">
					<DropdownMenuGroup>
						<DropdownMenuLabel class="flex items-center">
							<User size="1rem" /> Profile</DropdownMenuLabel
						>
						<DropdownMenuItem onclick={logout} class="text-gray-500">
							<LogOut />
							<span>Log Out</span>
						</DropdownMenuItem>
					</DropdownMenuGroup>
				</DropdownMenuContent>
			</DropdownMenu>
		</div>
	</div>
</div>
