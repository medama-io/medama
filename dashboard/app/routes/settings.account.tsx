import { useForm } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import {
	type ClientActionFunctionArgs,
	type MetaFunction,
	data as json,
	useLoaderData,
	useSubmit,
} from '@remix-run/react';
import { valibotResolver } from 'mantine-form-valibot-resolver';
import * as v from 'valibot';

import type { components } from '@/api/types';
import { userGet, userUpdate } from '@/api/user';
import { PasswordInput, TextInput } from '@/components/Input';
import { Section } from '@/components/settings/Section';
import { getString, getType } from '@/utils/form';

interface LoaderData {
	user: components['schemas']['UserGet'];
}

export const meta: MetaFunction = () => {
	return [{ title: 'Account Settings | Medama' }];
};

const accountSchema = v.object({
	_setting: v.literal('account', 'Invalid setting type.'),
	username: v.optional(
		v.pipe(
			v.string(),
			v.trim(),
			v.minLength(3, 'Username should not exceed 36 characters.'),
			v.maxLength(36, 'Username should not exceed 36 characters.'),
		),
	),
	password: v.optional(
		v.pipe(
			v.string(),
			v.check(
				(value) => value === '' || value.length >= 5,
				'Password should include at least 5 characters.',
			),
			v.maxLength(128, 'Password should not be larger than 128 characters.'),
			v.trim(),
		),
	),
});

export const clientLoader = async () => {
	const { data } = await userGet();

	if (!data) {
		throw json('Failed to get user.', {
			status: 500,
		});
	}

	return {
		user: data,
	};
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
					settings: {
						language: 'en',
					},
				},
				shouldThrow: false,
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

	const message = 'Successfully updated account details.';
	notifications.show({
		title: 'Success.',
		message,
		withBorder: true,
		color: '#17cd8c',
	});

	return { message };
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
		validate: valibotResolver(accountSchema),
	});

	const handleSubmit = (values: typeof account.values) => {
		submit(values, { method: 'POST' });
		account.setFieldValue('password', '');
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
				key={account.key('password')}
				{...account.getInputProps('password')}
			/>
		</Section>
	);
}
