import { useForm } from '@mantine/form';
import { Form, useSubmit } from '@remix-run/react';
import { valibotResolver } from 'mantine-form-valibot-resolver';
import * as v from 'valibot';

import { Button } from '@/components/Button';
import { PasswordInput, TextInput } from '@/components/Input';
import { Flex, Group } from '@/components/layout/Flex';

import classes from './Login.module.css';

const loginSchema = v.object({
	username: v.pipe(
		v.string(),
		v.trim(),
		v.minLength(3, 'Username should include at least 3 characters.'),
	),
	password: v.pipe(
		v.string(),
		v.trim(),
		v.minLength(5, 'Password should include at least 5 characters.'),
	),
});

export const Login = () => {
	const submit = useSubmit();
	const form = useForm({
		mode: 'uncontrolled',
		initialValues: { username: '', password: '' },
		validate: valibotResolver(loginSchema),
	});

	const handleSubmit = (values: typeof form.values) => {
		submit(values, { method: 'POST' });
	};

	return (
		<div className={classes.wrapper}>
			<h3 className={classes.title}>Log in to your dashboard</h3>
			<Form onSubmit={form.onSubmit(handleSubmit)}>
				<Flex>
					<TextInput
						key={form.key('username')}
						required
						label="Username"
						placeholder="Your username"
						autoComplete="username"
						{...form.getInputProps('username')}
					/>
					<PasswordInput
						key={form.key('password')}
						required
						label="Password"
						placeholder="Your password"
						{...form.getInputProps('password')}
					/>
				</Flex>

				<Button className={classes.submit} type="submit">
					<Group style={{ gap: 4 }}>
						<span>Log In</span>
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
				</Button>
			</Form>
		</div>
	);
};
