import {
	Anchor,
	Button,
	Checkbox,
	Container,
	Group,
	Paper,
	PasswordInput,
	TextInput,
	Title,
} from '@mantine/core';
import { Form } from '@remix-run/react';

import classes from './Login.module.css';
import { password$, username$ } from './observables';

export const Login = () => {
	const username = username$.use();
	const password = password$.use();

	return (
		<Container className={classes.wrapper}>
			<Title ta="center" className={classes.title}>
				Login to Medama
			</Title>

			<Paper withBorder shadow="md" p={30} mt={30} radius="md">
				<Form method="post">
					<TextInput
						name="username"
						label="Username"
						required
						onChange={(e) => username$.set(e.currentTarget.value)}
						value={username}
					/>
					<PasswordInput
						name="password"
						label="Password"
						required
						mt="md"
						onChange={(e) => password$.set(e.currentTarget.value)}
						value={password}
					/>
					<Group justify="space-between" mt="lg">
						<Checkbox label="Remember me" />
						<Anchor component="button" size="sm">
							Forgot password?
						</Anchor>
					</Group>
					<Button fullWidth mt="xl" type="submit">
						Sign in
					</Button>
				</Form>
			</Paper>
		</Container>
	);
};
