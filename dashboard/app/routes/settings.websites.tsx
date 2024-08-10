import { useForm } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import {
	type ClientActionFunctionArgs,
	type MetaFunction,
	json,
	useLoaderData,
	useSearchParams,
	useSubmit,
} from '@remix-run/react';
import { valibotResolver } from 'mantine-form-valibot-resolver';
import { useState } from 'react';
import * as v from 'valibot';

import { userGet } from '@/api/user';
import { websiteDelete, websiteList } from '@/api/websites';
import { TextInput } from '@/components/Input';
import { ModalChild, ModalWrapper } from '@/components/Modal';
import { Group } from '@/components/layout/Flex';
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
	return [{ title: 'Website Settings | Medama' }];
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

	const deleteSchema = v.object({
		_setting: v.literal('delete', 'Invalid setting type.'),
		hostname: v.pipe(
			v.string('Hostname is not string.'),
			v.check(
				(hostname) => hostname === website,
				'Domain name does not match.',
			),
		),
	});

	const form = useForm({
		mode: 'uncontrolled',
		initialValues: { _setting: 'delete', hostname: '' },
		validate: valibotResolver(deleteSchema),
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
		<ModalWrapper opened={opened} close={close}>
			<ModalChild
				title="Delete website"
				closeAriaLabel="Close delete website modal"
				description="This website's analytics data will be permanently deleted."
				submitLabel="Delete Website"
				onSubmit={form.onSubmit(handleSubmit)}
				close={resetAndClose}
				isDanger
			>
				<TextInput
					label="Enter the domain name"
					autoComplete="off"
					disabled={website === ''}
					key={form.key('hostname')}
					{...form.getInputProps('hostname')}
				/>
			</ModalChild>
		</ModalWrapper>
	);

	return (
		<>
			<SectionWrapper>
				<Group>
					<SectionTitle>
						<h3>Choose Website</h3>
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
