import { Flex, Group, Paper, SimpleGrid, Text } from '@mantine/core';
import { json, type MetaFunction, useLoaderData } from '@remix-run/react';

import type { components } from '@/api/types';
import { websiteList } from '@/api/websites';
import { InnerHeader } from '@/components/layout/InnerHeader';
import { ButtonDark } from '@/components/Button';
import { IconPlus } from '@/components/icons/plus';
import { WebsiteCard } from '@/components/index/WebsiteCard';
import { userLoggedIn } from '@/api/user';

interface LoaderData {
	websites: Array<components['schemas']['WebsiteGet']>;
}

export const meta: MetaFunction = () => {
	return [
		{ title: 'Medama | Privacy Focused Web Analytics' },
		{ name: 'description', content: 'Privacy focused web analytics.' },
	];
};

export const clientLoader = async () => {
	await userLoggedIn();

	const { data, res } = await websiteList({ query: { summary: true } });

	if (!res.ok) {
		if (res.status === 404) {
			return json<LoaderData>({ websites: [] });
		}

		throw json('Failed to fetch websites.', {
			status: res.status,
		});
	}

	if (!data) {
		throw json('Failed to fetch websites.', {
			status: 500,
		});
	}

	return json<LoaderData>({ websites: data });
};

export default function Index() {
	const { websites } = useLoaderData<LoaderData>();

	return (
		<>
			<InnerHeader>
				<Flex justify="space-between" align="center" py={8}>
					<h1>My Websites</h1>
					<ButtonDark to="/add">
						<Group>
							<IconPlus />
							<span>Add Website</span>
						</Group>
					</ButtonDark>
				</Flex>
			</InnerHeader>
			<main>
				{websites.length === 0 && (
					<Paper w="100%" p={16} radius={8} withBorder>
						<Text ta="center">No websites found. Please add a website!</Text>
					</Paper>
				)}
				<SimpleGrid cols={3}>
					{websites.map((website) => (
						<WebsiteCard key={website.hostname} website={website} />
					))}
				</SimpleGrid>
			</main>
		</>
	);
}
