import { SimpleGrid } from '@mantine/core';
import { schemaResolver, useForm } from '@mantine/form';
import { useDisclosure } from '@mantine/hooks';
import { Plus } from 'lucide-react';
import { data as json, redirect, useSubmit } from 'react-router';
import * as v from 'valibot';
import isFQDN from 'validator/lib/isFQDN';
import { websiteCreate, websiteList } from '@/api/websites';
import { Button } from '@/components/Button';
import { TextInput } from '@/components/Input';
import { WebsiteCard } from '@/components/index/WebsiteCard';
import { Group } from '@/components/layout/Flex';
import { InnerHeader } from '@/components/layout/InnerHeader';
import { ModalChild, ModalWrapper } from '@/components/Modal';
import type { Route } from './+types/_index';

export const meta: Route.MetaFunction = () => {
	return [
		{ title: 'Medama | Privacy Focused Web Analytics' },
		{ name: 'description', content: 'Privacy focused web analytics.' },
	];
};

export const clientLoader = async () => {
	const { data, res } = await websiteList({ query: { summary: true } });

	if (!res.ok || !data) {
		if (res.status === 404) {
			return { websites: [] };
		}

		throw json('Failed to fetch websites.', {
			status: res.status,
		});
	}

	return { websites: data };
};

export const clientAction = async ({ request }: Route.ClientActionArgs) => {
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

export default function Index({ loaderData }: Route.ComponentProps) {
	const { websites } = loaderData;
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
		validate: schemaResolver(addWebsiteSchema),
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
							<Plus size={16} />
							<span>Add Website</span>
						</Group>
					</Button>
				</Group>
			</InnerHeader>
			<main>
				{websites.length === 0 && (
					<div
						style={{
							width: '100%',
							padding: 16,
							borderRadius: 8,
							border: '1px solid var(--border-muted)',
						}}
					>
						<p style={{ textAlign: 'center' }}>
							No websites found. Please add a website!
						</p>
					</div>
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
