import { Group, Modal, Title } from '@mantine/core';
import { useForm, zodResolver } from '@mantine/form';
import { useDisclosure } from '@mantine/hooks';
import { notifications } from '@mantine/notifications';
import {
	json,
	useLoaderData,
	useSubmit,
	type ClientActionFunctionArgs,
	type MetaFunction,
} from '@remix-run/react';
import { z } from 'zod';

import { userGet, userUpdate } from '@/api/user';
import { websiteList } from '@/api/websites';
import { ModalChild, ModalInput } from '@/components/Modal';
import {
	SectionDanger,
	SectionTitle,
	SectionWrapper,
} from '@/components/settings/Section';
import { WebsiteSelector } from '@/components/settings/WebsiteSelector';
import { getString, getType } from '@/utils/form';

export const meta: MetaFunction = () => {
	return [{ title: 'Account Settings | Medama' }];
};

const deleteSchema = z.object({
	_setting: z.literal('delete'),
	hostname: z.string().min(1),
});

export const clientLoader = async () => {
	const { data: user } = await userGet();
	const { data: websites } = await websiteList();

	if (!user || !websites) {
		throw json('Failed to get user.', {
			status: 500,
		});
	}

	return json({
		user,
		websites: websites.map((website) => website.hostname),
	});
};

export const clientAction = async ({ request }: ClientActionFunctionArgs) => {
	const body = await request.formData();
	const type = getType(body);

	let res: Response | undefined;
	switch (type) {
		case 'delete': {
			const update = await userUpdate({
				body: {
					username: getString(body, 'username'),
					password: getString(body, 'password'),
				},
				noThrow: true,
			});
			res = update.res;
			break;
		}
		default:
			throw new Response('Invalid setting type.', {
				status: 400,
			});
	}

	if (!res || !res.ok) {
		throw new Response(res?.statusText || 'Failed to update user.', {
			status: res?.status || 500,
		});
	}

	const message = 'Successfully deleted website.';
	notifications.show({
		title: 'Success.',
		message,
		withBorder: true,
		color: '#17cd8c',
	});
	return json({ message });
};

export default function Index() {
	const { user, websites } = useLoaderData<typeof clientLoader>();
	const submit = useSubmit();
	const [opened, { open, close }] = useDisclosure(false);

	if (!user) {
		return;
	}

	const form = useForm({
		mode: 'uncontrolled',
		initialValues: {
			_setting: 'delete',
			hostname: websites[0] ?? '',
		},
		validate: zodResolver(deleteSchema),
	});

	const resetAndClose = () => {
		form.reset();
		close();
	};

	const handleSubmit = (values: typeof form.values) => {
		submit(values, { method: 'POST' });
	};

	const modalChildren = (
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
	);

	return (
		<>
			<SectionWrapper>
				<Group justify="space-between">
					<SectionTitle>
						<Title order={3}>Choose Website</Title>
					</SectionTitle>
					<WebsiteSelector websites={websites} />
				</Group>
			</SectionWrapper>
			<SectionDanger
				title="Delete Website"
				description="The website will be permanently deleted, including its data. This action is irreversible and can not be undone."
				modalChildren={modalChildren}
				open={open}
			>
				<input
					type="hidden"
					key={form.key('_setting')}
					{...form.getInputProps('_setting')}
				/>
			</SectionDanger>
		</>
	);
}
