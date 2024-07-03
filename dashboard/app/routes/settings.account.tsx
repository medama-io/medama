import { useForm, zodResolver } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import {
	json,
	useLoaderData,
	useSubmit,
	type ClientActionFunctionArgs,
	type MetaFunction,
} from '@remix-run/react';
import { z } from 'zod';

import type { components } from '@/api/types';
import { userGet, userUpdate } from '@/api/user';
import { PasswordInput, TextInput } from '@/components/settings/Input';
import { Section } from '@/components/settings/Section';
import { getString, getType } from '@/utils/form';

interface LoaderData {
	user: components['schemas']['UserGet'];
}

export const meta: MetaFunction = () => {
	return [{ title: 'Account Settings | Medama' }];
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
	const type = getType(body);

	let res: Response | undefined;
	switch (type) {
		case 'account': {
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

	if (!res) {
		throw new Error('Failed to update user.');
	}

	if (!res.ok) {
		throw new Response(res.statusText, {
			status: res.status,
		});
	}

	const message = 'Successfully updated account details.';
	notifications.show({
		title: 'Success.',
		message,
		withBorder: true,
		color: '#17cd8c',
	});
	return json({ message });
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

	const handleSubmit = (values: typeof account.values) => {
		submit(values, { method: 'POST' });
	};

	return (
		<Section
			title="Account Details"
			description="Edit your username and password."
			onSubmit={account.onSubmit(handleSubmit)}
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
