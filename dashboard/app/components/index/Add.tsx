import {
	Box,
	CloseButton,
	Group,
	Text,
	TextInput,
	UnstyledButton,
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { Form, useSubmit } from '@remix-run/react';
import { zodResolver } from 'mantine-form-zod-resolver';
import isFQDN from 'validator/lib/isFQDN';
import { z } from 'zod';

import { IconArrowRight } from '@/components/icons/arrow-right';

import classes from './Add.module.css';

interface AddProps {
	close: () => void;
}

const addWebsiteSchema = z.object({
	hostname: z
		.string()
		.refine((value) => value === 'localhost' || isFQDN(value), {
			message: 'Please enter a valid domain name.',
		}),
});

export const Add = ({ close }: AddProps) => {
	const submit = useSubmit();
	const form = useForm({
		mode: 'uncontrolled',
		initialValues: { hostname: '' },
		validate: zodResolver(addWebsiteSchema),
	});

	const resetAndClose = () => {
		form.reset();
		close();
	};

	const handleSubmit = (values: typeof form.values) => {
		submit(values, { method: 'POST' });
		resetAndClose();
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
			<Form onSubmit={form.onSubmit(handleSubmit)}>
				<TextInput
					label="Domain Name"
					placeholder="yourwebsite.com"
					description="The domain or subdomain name of your website."
					key={form.key('hostname')}
					{...form.getInputProps('hostname')}
					classNames={{ input: classes.input }}
					required
					mt="md"
					autoComplete="off"
					data-autofocus
				/>
				<UnstyledButton className={classes.submit} mt="xl" type="submit">
					<span>Add Website</span>
					<IconArrowRight />
				</UnstyledButton>
			</Form>
		</Box>
	);
};
