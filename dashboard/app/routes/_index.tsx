import { ModalChild, ModalWrapper } from '@/components/Modal';
import { Paper, SimpleGrid, Text } from '@mantine/core';
import { useForm } from '@mantine/form';
import { PlusIcon } from '@radix-ui/react-icons';
import {
	type ClientActionFunctionArgs,
	type MetaFunction,
	json,
	redirect,
	useLoaderData,
	useSubmit,
} from '@remix-run/react';
import { valibotResolver } from 'mantine-form-valibot-resolver';
import * as v from 'valibot';
import isFQDN from 'validator/lib/isFQDN';

import type { components } from '@/api/types';
import { userLoggedIn } from '@/api/user';
import { websiteCreate, websiteList } from '@/api/websites';
import { Button } from '@/components/Button';
import { TextInput } from '@/components/TextField';
import { WebsiteCard } from '@/components/index/WebsiteCard';
import { Group } from '@/components/layout/Flex';
import { InnerHeader } from '@/components/layout/InnerHeader';
import { useDisclosure } from '@/hooks/use-disclosure';

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

	const addWebsiteSchema = v.object({
		hostname: v.pipe(
			v.string('Hostname is not a string'),
			v.check(
				(value) => value === 'localhost' || isFQDN(value),
				'Please enter a valid domain name.',
			),
			v.check(
				(value) => !websites.find((website) => website.hostname === value),
				'Website already exists.',
			),
		),
	});

	const form = useForm({
		mode: 'uncontrolled',
		initialValues: { hostname: '' },
		validate: valibotResolver(addWebsiteSchema),
	});

	const resetAndClose = () => {
		close();
		form.reset();
	};

	const handleSubmit = (values: typeof form.values) => {
		submit(values, { method: 'POST' });
		resetAndClose();
	};

	return (
		<>
			<InnerHeader>
				<Group>
					<h1>My Websites</h1>
					<Button onClick={open} data-visible-from="xs">
						<Group style={{ gap: 8 }}>
							<PlusIcon />
							<span>Add Website</span>
						</Group>
					</Button>
				</Group>
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
				<ModalWrapper opened={opened} close={close}>
					<ModalChild
						title="Let's add your website"
						closeAriaLabel="Close add website modal"
						description="Tell us more about your website so we can add it to your dashboard."
						submitLabel="Add Website"
						onSubmit={form.onSubmit(handleSubmit)}
						close={resetAndClose}
					>
						<TextInput
							label="Domain Name"
							placeholder="yourwebsite.com"
							description="The domain or subdomain name of your website."
							required
							autoComplete="off"
							key={form.key('hostname')}
							{...form.getInputProps('hostname')}
						/>
					</ModalChild>
				</ModalWrapper>
			</main>
		</>
	);
}
