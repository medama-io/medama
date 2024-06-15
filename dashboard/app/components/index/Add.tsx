import {
	CloseButton,
	Container,
	Group,
	Text,
	TextInput,
	UnstyledButton,
} from '@mantine/core';
import { Form } from '@remix-run/react';
import { TimezoneSelect } from './TimezoneSelect';

import classes from './Add.module.css';

interface AddProps {
	hostname: string;
	setHostname: (hostname: string) => void;
	close: () => void;
}

export const Add = ({ hostname, setHostname, close }: AddProps) => {
	return (
		<Container className={classes.wrapper}>
			<Group justify="space-between">
				<h2>Let's add your website</h2>
				<CloseButton onClick={close} />
			</Group>
			<Text size="sm">
				Tell us more about your website so we can add it to your dashboard.
			</Text>
			<Form method="post">
				<TextInput
					name="hostname"
					label="Hostname"
					required
					mt="md"
					onChange={(e) => setHostname(e.currentTarget.value)}
					value={hostname}
					data-autofocus
				/>
				<TimezoneSelect />
				<UnstyledButton className={classes.submit} mt="xl" type="submit">
					Add Website
				</UnstyledButton>
			</Form>
		</Container>
	);
};
