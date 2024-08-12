import { type TransformedValues, useForm } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import {
	type ClientActionFunctionArgs,
	type MetaFunction,
	json,
	useLoaderData,
	useSubmit,
} from '@remix-run/react';
import { valibotResolver } from 'mantine-form-valibot-resolver';
import { useState } from 'react';
import * as v from 'valibot';

import type { components } from '@/api/types';
import { userGet, userUpdate } from '@/api/user';
import { Anchor } from '@/components/Anchor';
import { CheckBox } from '@/components/Checkbox';
import { Flex } from '@/components/layout/Flex';
import { CodeBlock } from '@/components/settings/Code';
import {
	Section,
	SectionTitle,
	SectionWrapper,
} from '@/components/settings/Section';
import { getString, getType } from '@/utils/form';

interface LoaderData {
	user: components['schemas']['UserGet'];
}

export const meta: MetaFunction = () => {
	return [{ title: 'Tracker Settings | Medama' }];
};

const trackerSchema = v.strictObject({
	_setting: v.literal('tracker', 'Invalid setting type.'),
	script_type: v.object({
		default: v.boolean(),
		'tagged-events': v.boolean(),
	}),
});

const getTrackingScript = (hostname: string) =>
	`<script defer src="https://${hostname}/script.js"></script>`;

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
	const scriptType = getString(body, 'script_type');

	let res: Response | undefined;
	switch (type) {
		case 'tracker': {
			const update = await userUpdate({
				body: {
					settings: {
						script_type: scriptType?.split(
							',',
						) as components['schemas']['UserGet']['settings']['script_type'],
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

	const [taggedEvents, setTaggedEvents] = useState(
		Boolean(user.settings.script_type?.includes('tagged-events')),
	);

	const code =
		location.hostname === 'localhost'
			? getTrackingScript('[your-analytics-server].com')
			: getTrackingScript(location.hostname);

	const form = useForm({
		mode: 'uncontrolled',
		initialValues: {
			_setting: 'tracker',
			script_type: {
				default: true,
				'tagged-events': taggedEvents,
			},
		},
		validate: valibotResolver(trackerSchema),
		transformValues: (values) => {
			// It's difficult to get Radix checkboxes to work with @mantine/form for now
			values.script_type['tagged-events'] = taggedEvents;

			// Convert object to comma-separated string
			const scriptType = Object.entries(values.script_type)
				.filter(([, value]) => value)
				.map(([key]) => key)
				.join(',');

			return {
				...values,
				script_type: scriptType,
			};
		},
	});

	const handleSubmit = (values: TransformedValues<typeof form>) => {
		submit(values, { method: 'POST' });
	};

	return (
		<>
			<Section
				title="Tracker Configuration"
				description="Choose what features you want to enable in the tracker."
				onSubmit={form.onSubmit(handleSubmit)}
			>
				<input
					type="hidden"
					key={form.key('_setting')}
					{...form.getInputProps('_setting')}
				/>
				<Flex style={{ gap: 16, marginTop: 8 }}>
					<CheckBox
						label="Default"
						value="default"
						tooltip={
							<>
								<p>Enable default page view tracking functionality.</p>
								<br />
								<p>
									Read our{' '}
									<Anchor href="https://oss.medama.io/methodology/overview">
										Methodology
									</Anchor>{' '}
									for more information.
								</p>
							</>
						}
						disabled
						checked
						key={form.key('script_type.default')}
						{...form.getInputProps('script_type.default', { type: 'checkbox' })}
					/>
					<CheckBox
						label="Tagged Events"
						value="tagged-events"
						tooltip={
							<>
								<p>Enable tracking of tagged events on your website.</p>
								<br />
								<p>
									Read our <Anchor>Tagged Events</Anchor> guide for more
									information.
								</p>
							</>
						}
						checked={taggedEvents}
						onCheckedChange={() => setTaggedEvents(!taggedEvents)}
						key={form.key('script_type.tagged-events')}
						{...form.getInputProps('script_type.tagged-events', {
							type: 'checkbox',
						})}
					/>
				</Flex>
			</Section>
			<SectionWrapper>
				<SectionTitle>
					<h3>Copy Tracker Code</h3>
					<p style={{ marginTop: 4 }}>
						Paste the following code in the <code>&lt;head&gt;</code> component
						on your website:
					</p>
					<CodeBlock code={code} />
					<p>
						Learn more about configuring the tracker in our{' '}
						<Anchor
							href="https://oss.medama.io/config/tracking-snippet"
							aria-label="Visit tracking snippet documentation"
						>
							documentation
						</Anchor>
						.
					</p>
				</SectionTitle>
			</SectionWrapper>
		</>
	);
}
