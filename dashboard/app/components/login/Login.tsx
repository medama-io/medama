import {
	Group,
	Paper,
	PasswordInput,
	Stack,
	Text,
	TextInput,
	UnstyledButton,
} from '@mantine/core';
import { useForm, zodResolver } from '@mantine/form';
import { Form, useSubmit } from '@remix-run/react';
import { z } from 'zod';

import classes from './Login.module.css';

const loginSchema = z.object({
	username: z
		.string()
		.min(3, {
			message: 'Username should include at least 3 characters.',
		})
		.max(48, {
			message: 'Username should not exceed 36 characters.',
		})
		.trim(),
	password: z
		.string()
		.min(5, {
			message: 'Password should include at least 5 characters.',
		})
		.trim(),
});

export const Login = () => {
	const submit = useSubmit();
	const form = useForm({
		mode: 'uncontrolled',
		initialValues: { username: '', password: '' },
		validate: zodResolver(loginSchema),
	});

	const handleSubmit = (values: typeof form.values) => {
		submit(values, { method: 'POST' });
	};

	return (
		<Paper className={classes.wrapper} withBorder>
			<Text size="lg" fw={500} mb="lg">
				Log in to your dashboard
			</Text>
			<Form onSubmit={form.onSubmit(handleSubmit)}>
				<Stack gap="lg">
					<TextInput
						key={form.key('username')}
						required
						label="Username"
						placeholder="Your username"
						radius="md"
						{...form.getInputProps('username')}
					/>
					<PasswordInput
						key={form.key('password')}
						required
						label="Password"
						placeholder="Your password"
						radius="md"
						{...form.getInputProps('password')}
					/>
				</Stack>

				<UnstyledButton mt="xl" className={classes.submit} type="submit">
					<Group align="center" gap="xs">
						<Text fz={14}>Log In</Text>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							width="18"
							height="18"
							viewBox="0 0 24 24"
							fill="none"
							stroke="currentColor"
							strokeWidth="2"
							strokeLinecap="round"
							strokeLinejoin="round"
						>
							<title>Log In Button</title>
							<path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4" />
							<polyline points="10 17 15 12 10 7" />
							<line x1="15" x2="3" y1="12" y2="12" />
						</svg>
					</Group>
				</UnstyledButton>
			</Form>
		</Paper>
	);
};
