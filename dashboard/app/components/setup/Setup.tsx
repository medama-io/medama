import {
	Button,
	Container,
	Paper,
	PasswordInput,
	TextInput,
	Title,
} from '@mantine/core';
import { Form } from '@remix-run/react';

import { email$, password$ } from './observables';
import classes from './Setup.module.css';

export const Setup = () => {
	const email = email$.use();
	const password = password$.use();

	return (
		<Container className={classes.wrapper}>
			<Title ta="center" className={classes.title}>
				Setup Medama
			</Title>
			<Paper withBorder shadow="md" p={30} mt={30} radius="md">
				<Form method="post">
					<TextInput
						name="email"
						label="Email"
						required
						onChange={(e) => email$.set(e.currentTarget.value)}
						value={email}
					/>
					<PasswordInput
						name="password"
						label="Password"
						required
						mt="md"
						onChange={(e) => password$.set(e.currentTarget.value)}
						value={password}
					/>
					<Button fullWidth mt="xl" type="submit">
						Create My Account
					</Button>
				</Form>
			</Paper>
		</Container>
	);
};
