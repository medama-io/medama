import { notifications } from '@mantine/notifications';
import {
	type ClientActionFunctionArgs,
	type MetaFunction,
	json,
	useLoaderData,
	useSubmit,
} from '@remix-run/react';
import { useForm } from '@tanstack/react-form';
import { valibotValidator } from '@tanstack/valibot-form-adapter';
import * as v from 'valibot';

import type { components } from '@/api/types';
import { userGet, userUpdate } from '@/api/user';
import { CheckBox } from '@/components/Checkbox';
import { Flex } from '@/components/layout/Flex';
import { Section } from '@/components/settings/Section';
import { getType } from '@/utils/form';

interface LoaderData {
	user: components['schemas']['UserGet'];
}

export const meta: MetaFunction = () => {
	return [{ title: 'Tracker Settings | Medama' }];
};

const SCRIPT_TYPES = {
	Default: 'default',
	'Tagged Events': 'tagged-events',
} as const;

const trackerSchema = v.strictObject({
	_setting: v.literal('tracker', 'Invalid setting type.'),
	script_type: v.array(
		v.enum(SCRIPT_TYPES, 'Invalid script type.'),
		'Invalid script type array.',
	),
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
		case 'tracker': {
			const update = await userUpdate({
				body: {
					settings: {
						script_type: 'default',
					},
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

	const message = 'Successfully updated tracker details.';
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

	const { Field, handleSubmit } = useForm({
		defaultValues: {
			_setting: 'tracker',
			script_type: user.settings.script_type
				? [user.settings.script_type]
				: ['default'],
		},
		validatorAdapter: valibotValidator(),
		validators: {
			onSubmit: trackerSchema,
		},
		onSubmit: (values) => {
			submit(values.value, { method: 'POST' });
		},
	});

	return (
		<>
			<Section
				title="Tracker Configuration"
				description="Choose what features you want to enable in the tracker."
				onSubmit={handleSubmit}
			>
				<Field name="_setting">{() => <input type="hidden" />}</Field>
				<Flex style={{ gap: 8 }}>
					<Field name="script_type" mode="array">
						{(field) => (
							<>
								{Object.entries(SCRIPT_TYPES).map(([key, value]) => (
									<CheckBox
										key={value}
										label={key}
										value={value}
										checked={field.state.value.includes(value)}
										disabled={value === SCRIPT_TYPES.Default}
										onCheckedChange={(checked) => {
											if (checked) {
												field.setValue([...field.state.value, value]);
											} else {
												field.setValue(
													field.state.value.filter((v) => v !== value),
												);
											}
										}}
									/>
								))}
							</>
						)}
					</Field>
				</Flex>
			</Section>
			<Flex>fa</Flex>
		</>
	);
}
