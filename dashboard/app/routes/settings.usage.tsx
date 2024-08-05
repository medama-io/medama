import { SimpleGrid } from '@mantine/core';
import { useForm, zodResolver } from '@mantine/form';
import { notifications } from '@mantine/notifications';
import {
	type ClientActionFunctionArgs,
	type MetaFunction,
	json,
	useLoaderData,
	useRevalidator,
	useSubmit,
} from '@remix-run/react';
import { useEffect, useMemo } from 'react';
import { z } from 'zod';

import { userUsageGet, userUpdate } from '@/api/user';
import { TextInput, TextInputWithTooltip } from '@/components/settings/Input';
import {
	ResourcePanel,
	ResourcePanelCPU,
} from '@/components/settings/Resource';
import { Section, SectionWrapper } from '@/components/settings/Section';
import { useInterval } from '@/hooks/use-interval';
import { getNumber, getString, getType } from '@/utils/form';

export const meta: MetaFunction = () => {
	return [{ title: 'Usage Settings | Medama' }];
};

export const clientLoader = async () => {
	const { data } = await userUsageGet();

	if (!data) {
		throw json('Failed to get server usage metrics.', {
			status: 500,
		});
	}

	return json({
		usage: data,
	});
};

export const clientAction = async ({ request }: ClientActionFunctionArgs) => {
	const body = await request.formData();
	const type = getType(body);

	let res: Response | undefined;
	switch (type) {
		case 'usage': {
			const update = await userUpdate({
				body: {
					settings: {
						threads: getNumber(body, 'threads'),
						memory_limit: getString(body, 'memory_limit'),
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
		throw new Response(res?.statusText || 'Failed to update usage settings.', {
			status: res?.status || 500,
		});
	}

	const message = 'Successfully updated usage settings.';
	notifications.show({
		title: 'Success.',
		message,
		withBorder: true,
		color: '#17cd8c',
	});
	return json({ message });
};

export default function Index() {
	const { usage } = useLoaderData<typeof clientLoader>();
	const submit = useSubmit();
	const revalidator = useRevalidator();
	const interval = useInterval(revalidator.revalidate, 2500);
	const { cpu, memory, disk, metadata } = usage;

	useEffect(() => {
		interval.start();
		return interval.stop;
	}, [interval.start, interval.stop]);

	const usageSchema = useMemo(
		() =>
			z.object({
				_setting: z.literal('usage'),
				threads: z.preprocess(
					(x) => (x ? x : undefined),
					z.coerce
						.number()
						.int()
						.min(1, {
							message: 'Threads must be at least 1.',
						})
						.max(metadata.threads ?? 1, {
							message: `Threads must be less than or equal to ${metadata.threads}.`,
						})
						.optional(),
				),
				memory_limit: z
					.string()
					.regex(/^(\d+(?:\.\d+)?)(MB|GB|TB|MiB|GiB|TiB)$/, {
						message:
							'Invalid memory limit format. Supported formats: 1MB, 1MiB, 1GB, 1GiB, 1TB, 1TiB.',
					})
					.optional(),
			}),
		[metadata.threads],
	);

	const usageForm = useForm({
		mode: 'uncontrolled',
		initialValues: {
			_setting: 'usage',
			threads: String(metadata.threads),
			memory_limit: metadata.memory_limit,
		},
		validate: zodResolver(usageSchema),
	});

	const handleSubmit = (values: Partial<typeof usageForm.values>) => {
		submit(values, { method: 'POST' });
		usageForm.reset();
	};

	return (
		<>
			<SectionWrapper>
				<SimpleGrid cols={{ base: 1, lg: 2 }}>
					<ResourcePanelCPU title="CPU Usage" {...cpu} />
					<ResourcePanel title="Memory Usage" {...memory} />
					<ResourcePanel title="Disk Usage" {...disk} />
				</SimpleGrid>
			</SectionWrapper>
			<Section
				title="Resource Allocation"
				description="Manage the allocation of system resources."
				onSubmit={usageForm.onSubmit(handleSubmit)}
			>
				<input
					type="hidden"
					key={usageForm.key('_setting')}
					{...usageForm.getInputProps('_setting')}
				/>
				<TextInput
					label="Threads"
					description="Default is # of available CPU threads."
					placeholder="4"
					key={usageForm.key('threads')}
					{...usageForm.getInputProps('threads')}
				/>
				<TextInputWithTooltip
					label="Memory Limit"
					description="Default is 80% of total memory."
					placeholder="1GB"
					tooltip="The maximum memory usable by the database (e.g. 1GB)"
					key={usageForm.key('memory_limit')}
					{...usageForm.getInputProps('memory_limit')}
				/>
			</Section>
		</>
	);
}
