import { Flex, Group, Modal, Paper, SimpleGrid, Text } from '@mantine/core';
import {
	json,
	redirect,
	useLoaderData,
	type ClientActionFunctionArgs,
	type MetaFunction,
} from '@remix-run/react';

import type { components } from '@/api/types';
import { userLoggedIn } from '@/api/user';
import { websiteCreate, websiteList } from '@/api/websites';
import { ButtonDark } from '@/components/Button';
import { IconPlus } from '@/components/icons/plus';
import { Add } from '@/components/index/Add';
import { WebsiteCard } from '@/components/index/WebsiteCard';
import { InnerHeader } from '@/components/layout/InnerHeader';
import { useDisclosure } from '@mantine/hooks';
import { useState } from 'react';

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

export const clientAction = async ({ request }: ClientActionFunctionArgs) => {
	const body = await request.formData();

	const hostname = body.get('hostname')
		? String(body.get('hostname'))
		: undefined;
	const name = body.get('name') ? String(body.get('name')) : hostname;

	if (!hostname) {
		throw json('Missing hostname', {
			status: 400,
		});
	}

	const { data, res } = await websiteCreate({
		body: {
			hostname,
		},
	});

	if (!data) {
		throw json('Failed to create website.', {
			status: res.status,
		});
	}

	return redirect(`/${data.hostname}`);
};

export default function Index() {
	const { websites } = useLoaderData<LoaderData>();
	const [opened, { open, close }] = useDisclosure(false);

	// Add website state
	const [hostname, setHostname] = useState('');
	const closeModalWithReset = () => {
		setHostname('');
		close();
	};

	return (
		<>
			<InnerHeader>
				<Flex justify="space-between" align="center" py={8}>
					<h1>My Websites</h1>
					<ButtonDark onClick={open}>
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
				<Modal opened={opened} onClose={close} withCloseButton={false} centered>
					<Add
						hostname={hostname}
						setHostname={setHostname}
						close={closeModalWithReset}
					/>
				</Modal>
			</main>
		</>
	);
}
