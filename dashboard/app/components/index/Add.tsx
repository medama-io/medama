import {
	Box,
	CloseButton,
	Group,
	Text,
	TextInput,
	UnstyledButton,
} from '@mantine/core';
import { Form } from '@remix-run/react';
import { TimezoneSelect } from './TimezoneSelect';

import { useState } from 'react';
import { useTimezoneSelect, type ITimezone } from 'react-timezone-select';
import classes from './Add.module.css';
import { IconArrowRight } from '../icons/arrow-right';

interface AddProps {
	close: () => void;
}

export const Add = ({ close }: AddProps) => {
	const [hostname, setHostname] = useState('');

	const { options: tzOptions, parseTimezone } = useTimezoneSelect({
		labelStyle: 'abbrev',
	});

	const [timezone, setTimezone] = useState<ITimezone>(
		parseTimezone(Intl.DateTimeFormat().resolvedOptions().timeZone),
	);

	const resetAndClose = () => {
		setHostname('');
		close();
	};

	return (
		<Box className={classes.wrapper}>
			<Group justify="space-between" align="center">
				<h2>Let's add your website</h2>
				<CloseButton
					size="lg"
					onClick={resetAndClose}
					aria-label="Close add website modal"
				/>
			</Group>
			<Text size="sm" mt="xs">
				Tell us more about your website so we can add it to your dashboard.
			</Text>
			<Form method="post">
				<TextInput
					name="hostname"
					label="Domain Name"
					placeholder="yourwebsite.com"
					description="The domain or subdomain name of your website."
					classNames={{ input: classes.input }}
					required
					mt="md"
					onChange={(e) => setHostname(e.currentTarget.value)}
					value={hostname}
					autoComplete="off"
					data-autofocus
				/>
				<TimezoneSelect
					timezone={timezone}
					setTimezone={setTimezone}
					parseTimezone={parseTimezone}
					tzOptions={tzOptions}
				/>
				<UnstyledButton className={classes.submit} mt="xl" type="submit">
					<span>Add Website</span>
					<IconArrowRight />
				</UnstyledButton>
			</Form>
		</Box>
	);
};
