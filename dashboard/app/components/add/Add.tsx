import { Button, Container, Paper, TextInput, Title } from '@mantine/core';
import { Form } from '@remix-run/react';

import classes from './Add.module.css';
import { hostname$, name$ } from './observables';

export const Add = () => {
	const name = name$.use();
	const hostname = hostname$.use();

	return (
		<Container className={classes.wrapper}>
			<Title ta="center" className={classes.title}>
				Login to Medama
			</Title>

			<Paper withBorder shadow="md" p={30} mt={30} radius="md">
				<Form method="post">
					<TextInput
						name="name"
						label="Name"
						onChange={(e) => name$.set(e.currentTarget.value)}
						value={name}
					/>
					<TextInput
						name="hostname"
						label="Hostname"
						required
						mt="md"
						onChange={(e) => hostname$.set(e.currentTarget.value)}
						value={hostname}
					/>
					<Button fullWidth mt="xl" type="submit">
						Add Website
					</Button>
				</Form>
			</Paper>
		</Container>
	);
};
