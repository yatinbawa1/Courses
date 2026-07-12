<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { page } from '$app/state';
	import { api } from '$lib/api/api';
	import AvatarBadge from '$lib/components/ui/avatar/avatar-badge.svelte';
	import AvatarFallback from '$lib/components/ui/avatar/avatar-fallback.svelte';
	import AvatarImage from '$lib/components/ui/avatar/avatar-image.svelte';
	import Avatar from '$lib/components/ui/avatar/avatar.svelte';
	import Button from '$lib/components/ui/button/button.svelte';
	import DropdownMenuContent from '$lib/components/ui/dropdown-menu/dropdown-menu-content.svelte';
	import DropdownMenuGroup from '$lib/components/ui/dropdown-menu/dropdown-menu-group.svelte';
	import DropdownMenuItem from '$lib/components/ui/dropdown-menu/dropdown-menu-item.svelte';
	import DropdownMenuLabel from '$lib/components/ui/dropdown-menu/dropdown-menu-label.svelte';
	import DropdownMenuTrigger from '$lib/components/ui/dropdown-menu/dropdown-menu-trigger.svelte';
	import DropdownMenu from '$lib/components/ui/dropdown-menu/dropdown-menu.svelte';
	import InputGroupAddon from '$lib/components/ui/input-group/input-group-addon.svelte';
	import InputGroupInput from '$lib/components/ui/input-group/input-group-input.svelte';
	import InputGroup from '$lib/components/ui/input-group/input-group.svelte';
	import NavigationMenuItem from '$lib/components/ui/navigation-menu/navigation-menu-item.svelte';
	import NavigationMenuList from '$lib/components/ui/navigation-menu/navigation-menu-list.svelte';
	import NavigationMenu from '$lib/components/ui/navigation-menu/navigation-menu.svelte';
	import { auth, logoutUser } from '$lib/stores/authStore/authStore';
	import {
		Search,
		House,
		GraduationCap,
		ShoppingBasket,
		Pen,
		Wallet,
		LogOut,
		User
	} from '@lucide/svelte';
	import type { Component } from 'svelte';
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

	let userProfileLink =
		'https://images.unsplash.com/photo-1773332611476-6ec2ba68049f?q=80&w=1830&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDF8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D';

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

<div class="w-full bg-gray-900 grid grid-cols-3 items-center px-5 py-4">
	<div class="playfair text-2xl tracking-[5px] text-white justify-self-start">COURSES</div>
	<div class="justify-center flex">
		<NavigationMenu>
			<NavigationMenuList class="flex gap-2">
				{#each menuItems as item (item.link)}
					{@const Icon = item.icon}
					<NavigationMenuItem>
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
					</NavigationMenuItem>
				{/each}
			</NavigationMenuList>
		</NavigationMenu>
	</div>
	<div class="w-full justify-center items-center flex gap-5">
		<InputGroup class="border w-full py-5 border-gray-800">
			<InputGroupAddon>
				<Search />
			</InputGroupAddon>
			<InputGroupInput placeholder="Search" class="text-gray-300"></InputGroupInput>
		</InputGroup>
		<DropdownMenu>
			<DropdownMenuTrigger>
				<Avatar>
					<AvatarImage src={userProfileLink} />
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
