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

import classes from './Login.module.css';
import { email$, password$ } from './observables';

export const Login = () => {
	const email = email$.use();
	const password = password$.use();

	return (
		<Container className={classes.wrapper}>
			<Title ta="center" className={classes.title}>
				Login to Medama
			</Title>

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
};
