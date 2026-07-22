<script>
	import Button from '$lib/components/ui/button/button.svelte';
	import InputGroupAddon from '$lib/components/ui/input-group/input-group-addon.svelte';
	import InputGroupInput from '$lib/components/ui/input-group/input-group-input.svelte';
	import InputGroup from '$lib/components/ui/input-group/input-group.svelte';
	import { Search } from '@lucide/svelte';
</script>

<div class="w-full p-5 playfair flex">
	<div class="flex-.5">
		<h1 class="mt-5 text-4xl">Expand your intellectual horizon.</h1>
		<p class="mt-3 geomini text-gray-500 text-lg">
			Curated courses from industry leaders, designed for deep focus and tangible skill acquisition
		</p>
		<div class="flex mt-4">
			<InputGroup class="py-5 text-lg">
				<InputGroupInput class="geomini" placeholder="What do you want to learn?" />
				<InputGroupAddon>
					<Search />
				</InputGroupAddon>
			</InputGroup>
			<Button variant="black" class="py-5 px-7 ml-3 geomin capitalize">Explore</Button>
		</div>
	</div>
</div>
