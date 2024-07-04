import { Flex, Group, Modal, Paper, SimpleGrid, Text } from '@mantine/core';
import { useForm } from '@mantine/form';
import {
	json,
	redirect,
	useLoaderData,
	useSubmit,
	type ClientActionFunctionArgs,
	type MetaFunction,
} from '@remix-run/react';
import { zodResolver } from 'mantine-form-zod-resolver';
import isFQDN from 'validator/lib/isFQDN';
import { z } from 'zod';

import type { components } from '@/api/types';
import { userLoggedIn } from '@/api/user';
import { websiteCreate, websiteList } from '@/api/websites';
import { ButtonDark } from '@/components/Button';
import { IconPlus } from '@/components/icons/plus';
import { WebsiteCard } from '@/components/index/WebsiteCard';
import { InnerHeader } from '@/components/layout/InnerHeader';
import { ModalChild, ModalInput } from '@/components/Modal';
import { useDisclosure } from '@mantine/hooks';

interface LoaderData {
	websites: Array<components['schemas']['WebsiteGet']>;
}

export const meta: MetaFunction = () => {
	return [
		{ title: 'Medama | Privacy Focused Web Analytics' },
		{ name: 'description', content: 'Privacy focused web analytics.' },
	];
};

const addWebsiteSchema = z.object({
	hostname: z
		.string()
		.refine((value) => value === 'localhost' || isFQDN(value), {
			message: 'Please enter a valid domain name.',
		}),
});

export const clientLoader = async () => {
	await userLoggedIn();

	const { data, res } = await websiteList({ query: { summary: true } });

	if (!res.ok || !data) {
		if (res.status === 404) {
			return json<LoaderData>({ websites: [] });
		}

		throw json('Failed to fetch websites.', {
			status: res.status,
		});
	}

	return json<LoaderData>({ websites: data });
};

export const clientAction = async ({ request }: ClientActionFunctionArgs) => {
	const body = await request.formData();

	const hostname = body.get('hostname')
		? String(body.get('hostname'))
		: undefined;

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
	const submit = useSubmit();

	const form = useForm({
		mode: 'uncontrolled',
		initialValues: { hostname: '' },
		validate: zodResolver(addWebsiteSchema),
	});

	const resetAndClose = () => {
		form.reset();
		close();
	};

	const handleSubmit = (values: typeof form.values) => {
		submit(values, { method: 'POST' });
		resetAndClose();
	};

	return (
		<>
			<InnerHeader>
				<Flex justify="space-between" align="center" py={8}>
					<h1>My Websites</h1>
					<ButtonDark onClick={open} visibleFrom="xs">
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
				<SimpleGrid cols={{ base: 1, xs: 2, md: 3 }}>
					{websites.map((website) => (
						<WebsiteCard key={website.hostname} website={website} />
					))}
				</SimpleGrid>
				<Modal
					opened={opened}
					onClose={close}
					withCloseButton={false}
					centered
					size="auto"
				>
					<ModalChild
						title="Let's add your website"
						closeAriaLabel="Close add website modal"
						description="Tell us more about your website so we can add it to your dashboard."
						submitLabel="Add Website"
						onSubmit={form.onSubmit(handleSubmit)}
						resetForm={resetAndClose}
					>
						<ModalInput
							label="Domain Name"
							placeholder="yourwebsite.com"
							description="The domain or subdomain name of your website."
							key={form.key('hostname')}
							{...form.getInputProps('hostname')}
							required
							mt="md"
							autoComplete="off"
							data-autofocus
						/>
					</ModalChild>
				</Modal>
			</main>
		</>
	);
}
