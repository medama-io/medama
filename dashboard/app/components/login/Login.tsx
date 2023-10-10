import {
	Anchor,
	Button,
	Checkbox,
	Container,
	Group,
	Paper,
	PasswordInput,
	Text,
	TextInput,
	Title,
} from '@mantine/core';
import { Form, NavLink } from '@remix-run/react';

import classes from './Login.module.css';
import { email$, password$ } from './observables';

export function Login() {
	const email = email$.use();
	const password = password$.use();

	return (
		<Container className={classes.wrapper}>
			<Title ta="center" className={classes.title}>
				Login to Medama
			</Title>
			<Text c="dimmed" size="sm" ta="center" mt={5}>
				{"Don't have an account? "}
				<Anchor component={NavLink} to="/signup" size="sm">
					Create account
				</Anchor>
			</Text>

			<Paper withBorder shadow="md" p={30} mt={30} radius="md">
				<TextInput
					label="Email"
					required
					onChange={(e) => email$.set(e.currentTarget.value)}
					value={email}
				/>
				<PasswordInput
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
				<Button fullWidth mt="xl">
					Sign in
				</Button>
			</Paper>
		</Container>
	);
}

export function Signup() {
	const email = email$.use();
	const password = password$.use();

	return (
		<Container className={classes.wrapper}>
			<Title ta="center" className={classes.title}>
				Signup to Medama
			</Title>
			<Text c="dimmed" size="sm" ta="center" mt={5}>
				{'Already have an account? '}
				<Anchor component={NavLink} to="/login" size="sm">
					Login
				</Anchor>
			</Text>

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
}
