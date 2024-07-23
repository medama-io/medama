import { Group, Text, Title } from '@mantine/core';
import { useForm, zodResolver } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import {
	json,
	useLoaderData,
	useSearchParams,
	useSubmit,
	type ClientActionFunctionArgs,
	type MetaFunction,
} from '@remix-run/react';
import { useState } from 'react';
import { z } from 'zod';

import { userGet } from '@/api/user';
import { websiteDelete, websiteList } from '@/api/websites';
import { ModalChild, ModalInput, ModalWrapper } from '@/components/Modal';
import {
	SectionDanger,
	SectionTitle,
	SectionWrapper,
} from '@/components/settings/Section';
import { WebsiteSelector } from '@/components/settings/WebsiteSelector';
import { useDidUpdate } from '@/hooks/use-did-update';
import { useDisclosure } from '@/hooks/use-disclosure';
import { getString, getType } from '@/utils/form';

export const meta: MetaFunction = () => {
	return [{ title: 'Account Settings | Medama' }];
};

export const clientLoader = async () => {
	const [{ data: user }, { data: websites }] = await Promise.all([
		userGet(),
		websiteList(),
	]);

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
			const deleteWebsite = await websiteDelete({
				pathKey: getString(body, 'hostname'),
				noThrow: true,
			});

			res = deleteWebsite.res;
			break;
		}
		default:
			throw new Response('Invalid setting type.', {
				status: 400,
			});
	}

	if (!res || !res.ok) {
		throw new Response(res?.statusText || 'Failed to delete website.', {
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

	const [searchParams, setSearchParams] = useSearchParams();
	const [website, setWebsite] = useState<string>(
		searchParams.get('website') ?? websites[0] ?? '',
	);

	useDidUpdate(() => {
		setSearchParams((prevParams) => {
			const newParams = new URLSearchParams(prevParams);
			newParams.set('website', website);
			return newParams;
		});
	}, [website]);

	const deleteSchema = z.object({
		_setting: z.literal('delete'),
		hostname: z.string().refine((hostname) => hostname === website, {
			message: 'Domain name does not match.',
		}),
	});

	const form = useForm({
		mode: 'uncontrolled',
		initialValues: {
			_setting: 'delete',
			hostname: '',
		},
		validate: zodResolver(deleteSchema),
	});

	const resetAndClose = () => {
		close();
		form.reset();
	};

	const handleSubmit = (values: typeof form.values) => {
		submit(values, { method: 'POST' });
		resetAndClose();
		setWebsite(websites[0] ?? '');
	};

	if (!user) {
		return null;
	}

	const modalChildren = (
		<ModalWrapper opened={opened} onClose={close}>
			<ModalChild
				title="Delete website"
				closeAriaLabel="Close delete website modal"
				description="This website's analytics data will be permanently deleted."
				submitLabel="Delete Website"
				onSubmit={form.onSubmit(handleSubmit)}
				resetForm={resetAndClose}
				isDanger
			>
				<ModalInput
					label={
						<Text fz={13} mb={4}>
							Enter the domain name{' '}
							<Text fz={13} fw={600} component="span">
								{website}
							</Text>{' '}
							to continue:
						</Text>
					}
					key={form.key('hostname')}
					{...form.getInputProps('hostname')}
					mt="md"
					autoComplete="off"
					data-autofocus
					disabled={website === ''}
				/>
			</ModalChild>
		</ModalWrapper>
	);

	return (
		<>
			<SectionWrapper>
				<Group justify="space-between">
					<SectionTitle>
						<Title order={3}>Choose Website</Title>
					</SectionTitle>
					<WebsiteSelector
						websites={websites}
						website={website}
						setWebsite={setWebsite}
					/>
				</Group>
			</SectionWrapper>
			<SectionDanger
				title="Delete Website"
				description="The website's analytics data will be permanently deleted. This action is irreversible and can not be undone."
				modalChildren={modalChildren}
				open={open}
				disabled={websites.length === 0}
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
