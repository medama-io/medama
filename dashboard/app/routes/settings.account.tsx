import { useForm, zodResolver } from '@mantine/form';
import {
	json,
	useActionData,
	useLoaderData,
	useSubmit,
	type ClientActionFunctionArgs,
	type MetaFunction,
} from '@remix-run/react';
import { z } from 'zod';

import type { components } from '@/api/types';
import { userGet, userLoggedIn, userUpdate } from '@/api/user';
import { PasswordInput, TextInput } from '@/components/settings/Input';
import { Section } from '@/components/settings/Section';
import { notifications } from '@mantine/notifications';

interface LoaderData {
	user: components['schemas']['UserGet'];
}

export const meta: MetaFunction = () => {
	return [
		{ title: 'Settings | Medama' },
		{ name: 'description', content: 'Privacy focused web analytics.' },
	];
};

const accountSchema = z.object({
	_setting: z.literal('account'),
	username: z
		.string()
		.max(48, {
			message: 'Username should not exceed 36 characters.',
		})
		.trim()
		.refine((value) => value.length === 0 || value.length >= 3, {
			message: 'Username should include at least 3 characters.',
		})
		.optional(),
	password: z
		.string()
		.max(128, {
			message: 'Password should not be larger than 128 characters.',
		})
		.trim()
		.refine((value) => value.length === 0 || value.length >= 5, {
			message: 'Password should include at least 5 characters.',
		})
		.optional(),
});

export const clientLoader = async () => {
	await userLoggedIn();

	const { data } = await userGet();

	if (!data) {
		throw json('Failed to get user.', {
			status: 500,
		});
	}

	return json<LoaderData>({
		user: data,
	});
};

export const clientAction = async ({ request }: ClientActionFunctionArgs) => {
	const body = await request.formData();
	const type = body.get('_setting') ? String(body.get('_setting')) : undefined;

	const getValue = (key: string) => {
		const value = body.get(key);
		return value && value !== '' ? String(value) : undefined;
	};

	let res: Response | undefined;
	switch (type) {
		case 'account': {
			const update = await userUpdate({
				body: {
					username: getValue('username'),
					password: getValue('password'),
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

	if (!res) {
		throw new Error('Failed to update user.');
	}

	if (!res.ok) {
		throw new Response(await res.text(), {
			status: res.status,
		});
	}

	notifications.show({
		title: 'Success.',
		message: 'Successfully updated account details.',
		withBorder: true,
		color: '#17cd8c',
	});
	return json({ message: 'Successfully updated user.' });
};

export default function Index() {
	const { user } = useLoaderData<LoaderData>();
	const submit = useSubmit();

	if (!user) {
		return;
	}

	const account = useForm({
		mode: 'uncontrolled',
		initialValues: {
			_setting: 'account',
			username: user.username,
			password: '',
		},
		validate: zodResolver(accountSchema),
	});

	const handleAccountSubmit = async (values: typeof account.values) => {
		submit(values, { method: 'POST' });
	};

	return (
		<Section
			title="Account details"
			description="Edit your username and password."
			onSubmit={account.onSubmit(handleAccountSubmit)}
		>
			<input
				type="hidden"
				key={account.key('_setting')}
				{...account.getInputProps('_setting')}
			/>
			<TextInput
				label="Username"
				placeholder="Username"
				key={account.key('username')}
				{...account.getInputProps('username')}
			/>
			<PasswordInput
				label="Password"
				placeholder="New password"
				type="password"
				key={account.key('password')}
				{...account.getInputProps('password')}
			/>
		</Section>
	);
}
