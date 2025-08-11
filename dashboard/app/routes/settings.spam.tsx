import { Box, Code, Flex, Stack, Text } from '@mantine/core';
import { useForm } from '@mantine/form';
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
import { useCallback, useState } from 'react';
import * as v from 'valibot';

import { userGet, userUpdate } from '@/api/user';
import { Anchor } from '@/components/Anchor';
import { Button } from '@/components/Button';
import { Card } from '@/components/Card';
import { Checkbox } from '@/components/Checkbox';
import { InputWithButton } from '@/components/Input';
import { SectionStack, SectionSubtitle } from '@/components/settings/Section';
import { getBoolean, getString, getType } from '@/utils/form';

export const meta: MetaFunction = () => {
	return [{ title: 'Spam Settings | Medama' }];
};

const FIREHOL_LEVEL1_URL = 'https://iplists.firehol.org/?ipset=firehol_level1';
const TOR_PROJECT_URL = 'https://check.torproject.org/torbulkexitlist';

const spamSchema = v.object({
	_setting: v.literal('spam', 'Invalid setting type.'),
	blockAbusiveIPs: v.boolean(),
	blockTorExitNodes: v.boolean(),
	blockedIPs: v.array(v.string()),
});

export const clientLoader = async (_: ClientLoaderFunctionArgs) => {
	const { data } = await userGet();

	if (!data) {
		throw json('Failed to get user.', {
			status: 500,
		});
	}

	return {
		settings: {
			blockAbusiveIPs: data.settings?.blockAbusiveIPs ?? false,
			blockTorExitNodes: data.settings?.blockTorExitNodes ?? false,
			blockedIPs: data.settings?.blockedIPs ?? [],
		},
	};
};

export const clientAction = async ({ request }: ClientActionFunctionArgs) => {
	const body = await request.formData();
	const type = getType(body);

	let res: Response | undefined;
	switch (type) {
		case 'spam': {
			const blockedIPs = (body.getAll('blockedIPs') as string[]).filter(
				(ip) => ip.trim() !== '',
			);

			const update = await userUpdate({
				body: {
					settings: {
						blockAbusiveIPs: getBoolean(body, 'blockAbusiveIPs'),
						blockTorExitNodes: getBoolean(body, 'blockTorExitNodes'),
						blockedIPs: blockedIPs,
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
		throw new Response(res?.statusText || 'Failed to update settings.', {
			status: res?.status || 500,
		});
	}

	notifications.show({
		title: 'Success',
		message: 'Successfully updated spam protection settings.',
		withBorder: true,
		color: 'green',
	});

	return { ok: true };
};

export default function SpamPage() {
	const { settings } = useLoaderData<typeof clientLoader>();
	const submit = useSubmit();
	const [newIP, setNewIP] = useState('');

	const form = useForm({
		mode: 'uncontrolled',
		initialValues: {
			_setting: 'spam',
			blockAbusiveIPs: settings.blockAbusiveIPs,
			blockTorExitNodes: settings.blockTorExitNodes,
			blockedIPs: settings.blockedIPs || [],
		},
		validate: valibotResolver(spamSchema),
	});

	const handleSubmit = useCallback(
		(values: typeof form.values) => {
			submit(values, { method: 'POST' });
		},
		[submit],
	);

	const handleAddIP = useCallback(() => {
		const currentValues = form.getValues();
		const trimmedIP = newIP?.trim();

		if (trimmedIP && !currentValues.blockedIPs.includes(trimmedIP)) {
			form.insertListItem('blockedIPs', trimmedIP);
			setNewIP('');

			handleSubmit(currentValues);
		}
	}, [form, newIP, handleSubmit]);

	const handleDeleteIP = useCallback(
		(index: number) => {
			form.removeListItem('blockedIPs', index);

			handleSubmit(form.getValues());
		},
		[form, handleSubmit],
	);

	const blockedIPs = form.getValues().blockedIPs || [];

	return (
		<>
			<SectionStack
				title="Preset Filters"
				description="Enable our curated blocklists to filter out common spam sources."
				onSubmit={form.onSubmit(handleSubmit)}
			>
				<input
					type="hidden"
					key={form.key('_setting')}
					{...form.getInputProps('_setting')}
				/>

				<Flex justify="space-between" align="center">
					<Box flex={1} pr="md">
						<SectionSubtitle>Abusive IPs</SectionSubtitle>
						<p>
							Block known IP addresses that are responsible for attacks,
							malware, and spam. Sourced from{' '}
							<Anchor href={FIREHOL_LEVEL1_URL}>Firehol Level 1</Anchor>.
						</p>
					</Box>
					<Checkbox
						key={form.key('blockAbusiveIPs')}
						{...form.getInputProps('blockAbusiveIPs', { type: 'checkbox' })}
					/>
				</Flex>
				<Flex justify="space-between" align="center">
					<Box flex={1} pr="md">
						<SectionSubtitle>Tor Exit Nodes</SectionSubtitle>
						<p>
							Block traffic from known Tor exit nodes. Sourced from the{' '}
							<Anchor href={TOR_PROJECT_URL}>Tor Project</Anchor>.
						</p>
					</Box>
					<Checkbox
						value="tor"
						key={form.key('blockTorExitNodes')}
						{...form.getInputProps('blockTorExitNodes', { type: 'checkbox' })}
					/>
				</Flex>
			</SectionStack>
			<SectionStack
				title="Filter IP addresses"
				description="Block specific IP addresses from sending data to your dashboard."
				hasButton={false}
			>
				<input
					type="hidden"
					key={form.key('_setting')}
					{...form.getInputProps('_setting')}
				/>
				<input type="hidden" name="blockedIPs" value={blockedIPs} />
				<div>
					<SectionSubtitle>Add IP address</SectionSubtitle>
					<InputWithButton
						aria-label="IP Address"
						placeholder="192.168.1.1"
						value={newIP}
						onChange={(e) => setNewIP(e.currentTarget.value)}
						onKeyDown={(e) => {
							if (e.key === 'Enter') {
								e.preventDefault();
								handleAddIP();
							}
						}}
						buttonLabel="Add"
						onButtonClick={handleAddIP}
						autoComplete="off"
					/>
				</div>

				<Box mt="xs">
					<SectionSubtitle>Blocked IP addresses</SectionSubtitle>
					<Card>
						<Stack>
							{blockedIPs.length > 0 ? (
								blockedIPs.map((ip, index) => (
									<Flex key={ip} justify="space-between" align="center" p="xs">
										<Code pl="xs">{ip}</Code>
										<Button
											//	rightIcon={<IconTrash size={16} />}
											onClick={() => handleDeleteIP(index)}
										>
											Delete
										</Button>
									</Flex>
								))
							) : (
								<Flex p="md" justify="center">
									<Text size="sm" c="dimmed">
										No IP addresses have been blocked yet.
									</Text>
								</Flex>
							)}
						</Stack>
					</Card>
				</Box>
			</SectionStack>
		</>
	);
}
