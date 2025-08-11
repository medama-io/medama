import { type TransformedValues, useForm } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import {
	type ClientActionFunctionArgs,
	type ClientLoaderFunctionArgs,
	type MetaFunction,
	data as json,
	useLoaderData,
	useSubmit,
} from '@remix-run/react';
import { valibotResolver } from 'mantine-form-valibot-resolver';
import * as v from 'valibot';

import type { components } from '@/api/types';
import { userGet, userUpdate } from '@/api/user';
import { Anchor } from '@/components/Anchor';
import { Checkbox } from '@/components/Checkbox';
import { Flex } from '@/components/layout/Flex';
import { CodeBlock } from '@/components/settings/Code';
import {
	Section,
	SectionTitle,
	SectionWrapper,
} from '@/components/settings/Section';
import { getString, getType } from '@/utils/form';

export const meta: MetaFunction = () => {
	return [{ title: 'Tracker Settings | Medama' }];
};

const trackerSchema = v.strictObject({
	_setting: v.literal('tracker', 'Invalid setting type.'),
	script_type: v.object({
		default: v.boolean(),
		'click-events': v.boolean(),
		'page-events': v.boolean(),
	}),
});

const getTrackingScript = (hostname: string) =>
	`<script defer src="https://${hostname}/script.js"></script>`;

export const clientLoader = async (_: ClientLoaderFunctionArgs) => {
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

	const message = 'Successfully updated tracker details.';
	notifications.show({
		title: 'Success.',
		message,
		withBorder: true,
		color: '#17cd8c',
	});

	return { message };
};

export default function Index() {
	const { user } = useLoaderData<typeof clientLoader>();
	const submit = useSubmit();

	if (!user) {
		return;
	}

	const form = useForm({
		mode: 'uncontrolled',
		initialValues: {
			_setting: 'tracker' as const,
			script_type: {
				default: true,
				'click-events': Boolean(
					user.settings.script_type?.includes('click-events'),
				),
				'page-events': Boolean(
					user.settings.script_type?.includes('page-events'),
				),
			},
		},
		validate: valibotResolver(trackerSchema),
		transformValues: (values) => {
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

	const code =
		location.hostname === 'localhost'
			? getTrackingScript('[your-analytics-server].com')
			: getTrackingScript(location.hostname);

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
					<Checkbox
						label="Default"
						value="default"
						tooltip={
							<>
								<p>The default page view tracking functionality.</p>
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
					<Checkbox
						label="Click Events"
						value="click-events"
						tooltip={
							<>
								<p>
									Enable custom properties tracking of click events on your
									website.
								</p>
								<br />
								<p>
									Read our{' '}
									<Anchor href="http://oss.medama.io/features/custom-properties/overview">
										Custom Properties
									</Anchor>{' '}
									and{' '}
									<Anchor href="http://oss.medama.io/features/custom-properties/click-events">
										Click Events
									</Anchor>{' '}
									guide for more information.
								</p>
							</>
						}
						key={form.key('script_type.click-events')}
						{...form.getInputProps('script_type.click-events', {
							type: 'checkbox',
						})}
					/>
					<Checkbox
						label="Page View Events"
						value="page-events"
						tooltip={
							<>
								<p>Enable tracking of page view events on your website.</p>
								<br />
								<p>
									Read our{' '}
									<Anchor href="http://oss.medama.io/features/custom-properties/overview">
										Custom Properties
									</Anchor>{' '}
									and{' '}
									<Anchor href="http://oss.medama.io/features/custom-properties/page-events">
										Page Events
									</Anchor>{' '}
									guide for more information.
								</p>
							</>
						}
						key={form.key('script_type.page-events')}
						{...form.getInputProps('script_type.page-events', {
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
